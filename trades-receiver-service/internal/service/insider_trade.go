package service

import (
	"context"
	"errors"
	"math"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
)

type insiderTradeService struct {
	repo      InsiderTradeRepo
	publisher InsiderTradePublisher
	l         *logger.Logger
}

func New(r InsiderTradeRepo, p InsiderTradePublisher, logger *logger.Logger) *insiderTradeService {
	return &insiderTradeService{
		repo:      r,
		publisher: p,
		l:         logger,
	}
}

func (s *insiderTradeService) Receive(ctx context.Context, trade *entity.Trade) error {
	var err error
	trade.Trs, err = s.fillTransaction(trade.Sth)
	if err != nil {
		return err
	}

	err = s.store(ctx, trade)
	if err != nil {
		s.l.Error("failed to store trade: ", err)
		return err
	}

	err = s.publisher.Publish(ctx, trade)
	if err != nil {
		return err
	}

	return nil
}

func (s *insiderTradeService) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return []*entity.Transaction{{}}, nil
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
		return nil, errors.New("failed to match transaction code")
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
		return err
	}

	return nil
}
