package service

import (
	"context"
	"errors"
	"fmt"
	"math"

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

	err = s.publisher.Publish(ctx, trade)
	if err != nil {
		return fmt.Errorf("%v: %w", methodName, err)
	}

	return nil
}

// TODO refactor and test
func (s *insiderTradeService) GetAll(ctx context.Context, cursor string, limit int) ([]*entity.Transaction, string, error) {
	transactions, nextCursor, err := s.repo.GetAll(ctx, cursor, limit)
	if err != nil {
		return nil, "", fmt.Errorf("(s *insiderTradeService) GetAll: %w", err)
	}

	return transactions, nextCursor, nil
}

// TODO cover cases with derivatives
func (s *insiderTradeService) fillTransaction(securityHoldings []*entity.SecurityTransactionHoldings) (*entity.Transaction, error) {

	const (
		purchaseCode = 0
		saleCode     = 1
	)

	var totalNonDerivative, totalPrice, averagePrice, totalValue float64
	var totalShares int

	for _, sth := range securityHoldings {
		switch sth.TransactionCode {
		case purchaseCode:
			totalValue += float64(sth.Quantity) * sth.PricePerSecurity
			totalPrice += sth.PricePerSecurity
			totalShares += sth.Quantity

			totalNonDerivative++

		case saleCode:
			totalValue -= float64(sth.Quantity) * sth.PricePerSecurity
			totalPrice += sth.PricePerSecurity
			totalShares += sth.Quantity

			totalNonDerivative++
		}
	}

	if totalNonDerivative == 0 {
		return nil, errors.New("fill transaction: failed to match transaction code: derivative only transactions")
	}

	averagePrice = totalPrice / totalNonDerivative

	transactionNames := map[int]string{
		purchaseCode: "BUY",
		saleCode:     "SELL",
	}

	var transactionName string

	if totalValue > 0.0 {
		transactionName = transactionNames[purchaseCode]
	} else {
		transactionName = transactionNames[saleCode]
		totalValue = math.Abs(totalValue)
	}

	transaction := &entity.Transaction{
		TransactionTypeName: transactionName,
		AveragePrice:        averagePrice,
		TotalShares:         totalShares,
		TotalValue:          totalValue,
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
