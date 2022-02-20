package service

import (
	"context"
	"errors"
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

func (s *insiderTradeService) Receive(ctx context.Context, trade *entity.Transaction) error {
	err := s.store(ctx, trade)
	if err != nil {
		return err
	}
	return nil
}

func (s *insiderTradeService) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return []*entity.Transaction{{}}, nil
}

// TODO change to private, cover cases with derivatives
func (s *insiderTradeService) FillTransaction(secFilings *entity.SecFiling, securityHoldings []*entity.SecurityTransactionHoldings) (*entity.Transaction, error) {

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
		SecFilingsID:        secFilings.ID,
		TransactionTypeName: transactionName,
		AveragePrice:        averagePrice,
		TotalShares:         totalShares,
		TotalValue:          totalValue,
	}

	return transaction, nil
}

func (s *insiderTradeService) store(ctx context.Context, trade *entity.Transaction) error {
	err := s.repo.Store(ctx, trade)
	if err != nil {
		return err
	}

	return nil
}
