package service

import (
	"context"
	"insidertradesreceiver/internal/entity"
)

type insiderTradeService struct {
	repo InsiderTradeRepo
	publisher InsiderTradePublisher
}

func New(r InsiderTradeRepo, p InsiderTradePublisher) *insiderTradeService {
	return &insiderTradeService{
		repo: r,
		publisher: p,
	}
}

func (s *insiderTradeService) Receive(ctx context.Context, trade *entity.InsiderTrade) error {
	return nil
}

func (s *insiderTradeService) GetAll(ctx context.Context, limit, offset int) ([]*entity.InsiderTrade, error) {
	return []*entity.InsiderTrade{&entity.InsiderTrade{}}, nil
}

func (s *insiderTradeService) store(ctx context.Context, insiderTrade *entity.InsiderTrade) error {
	err := s.repo.Store(ctx, insiderTrade)
	if err != nil {
		return err
	}

	return nil
}
