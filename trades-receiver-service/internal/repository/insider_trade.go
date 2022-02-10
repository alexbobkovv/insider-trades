package repository

import (
	"context"
	"insidertradesreceiver/internal/entity"
	"insidertradesreceiver/pkg/postgresql"
)

type InsiderTradeRepo struct {
	*postgresql.Postgres
}

func New(db *postgresql.Postgres) *InsiderTradeRepo {
	return &InsiderTradeRepo{db}
}

func (r *InsiderTradeRepo) GetAll(ctx context.Context, limit, offset int) ([]*entity.InsiderTrade, error) {
	return []*entity.InsiderTrade{&entity.InsiderTrade{}}, nil
}

func (r *InsiderTradeRepo) Store(ctx context.Context, trade *entity.InsiderTrade) error {
	return nil
}
