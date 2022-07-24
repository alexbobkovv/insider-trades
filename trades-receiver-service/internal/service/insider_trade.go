package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

type insiderTradeService struct {
	repo      InsiderTradeRepo
	publisher InsiderTradePublisher
}

func New(r InsiderTradeRepo, p InsiderTradePublisher) *insiderTradeService {
	return &insiderTradeService{
		repo:      r,
		publisher: p,
	}
}

func (s *insiderTradeService) Receive(ctx context.Context, trade *entity.Trade) error {
	var err error
	const methodName = "(s *insiderTradeService) Receive"

	trade.Trs, err = s.fillTransaction(trade.Sth)
	if err != nil {
		return fmt.Errorf("%v: %w", methodName, err)
	}

	err = s.store(ctx, trade)
	if err != nil {
		return fmt.Errorf("%v: failed to store trade: %w", methodName, err)
	}

	err = s.publisher.PublishTrade(ctx, trade)
	if err != nil {
		return fmt.Errorf("%v: %w", methodName, err)
	}

	return nil
}

// TODO refactor and test
func (s *insiderTradeService) ListTransactions(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*entity.Transaction, *cursor.Cursor, error) {
	transactions, nextCursor, err := s.repo.ListTransactions(ctx, reqCursor, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("(s *insiderTradeService) ListTrades: %w", err)
	}

	return transactions, nextCursor, nil
}

// TODO implement refresh mat view
func (s *insiderTradeService) ListViews(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error) {
	tradeViews, nextCursor, err := s.repo.ListViews(ctx, reqCursor, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("ListViews: %w", err)
	}

	return tradeViews, nextCursor, nil

}

func (s *insiderTradeService) fillTransaction(securityHoldings []*entity.SecurityTransactionHoldings) (*entity.Transaction, error) {

	if securityHoldings == nil {
		return nil, errors.New("fill transaction: empty trade.SecurityTransactionHoldings")
	}

	const (
		purchaseCode = 0
		saleCode     = 1
	)

	var totalNonDerivative uint8
	var totalValue, value, averagePrice, totalShares big.Rat

	// TODO fix big.Rat
	for _, sth := range securityHoldings {
		quantity := new(big.Rat).SetInt64(sth.Quantity)
		switch sth.TransactionCode {
		case purchaseCode:
			value.Mul(quantity, &sth.PricePerSecurity)
			totalValue.Add(&totalValue, &value)

			totalShares.Add(&totalShares, quantity)

			totalNonDerivative++

		case saleCode:
			value.Mul(quantity, &sth.PricePerSecurity)
			totalValue.Sub(&totalValue, &value)

			totalShares.Add(&totalShares, quantity)

			totalNonDerivative++
		}
	}

	if totalNonDerivative == 0 {
		return nil, errors.New("fill transaction: failed to match transaction code: derivative only transactions")
	}

	averagePrice.Quo(&totalValue, &totalShares)

	transactionNames := map[int]string{
		purchaseCode: "BUY",
		saleCode:     "SELL",
	}

	var transactionName string

	totalVal, _ := totalValue.Float64()
	avgPrice, _ := averagePrice.Float64()
	tlShares, _ := totalShares.Float64()
	if totalVal > 0.0 {
		transactionName = transactionNames[purchaseCode]
	} else {
		transactionName = transactionNames[saleCode]
		totalVal = math.Abs(totalVal)
		avgPrice = math.Abs(avgPrice)
	}

	transaction := &entity.Transaction{
		TransactionTypeName: transactionName,
		AveragePrice:        avgPrice,
		TotalShares:         int64(tlShares),
		TotalValue:          totalVal,
	}

	return transaction, nil
}

func (s *insiderTradeService) store(ctx context.Context, trade *entity.Trade) error {
	err := s.repo.StoreTrade(ctx, trade)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	return nil
}
