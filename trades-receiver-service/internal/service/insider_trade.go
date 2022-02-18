package service

import (
	"context"

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

func (s *insiderTradeService) store(ctx context.Context, trade *entity.Transaction) error {
	err := s.repo.Store(ctx, trade)
	if err != nil {
		return err
	}

	return nil
}
