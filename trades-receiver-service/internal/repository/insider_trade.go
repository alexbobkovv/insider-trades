package repository

import (
	"context"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/postgresql"
)

type InsiderTradeRepo struct {
	*postgresql.Postgres
}

func New(db *postgresql.Postgres) *InsiderTradeRepo {
	return &InsiderTradeRepo{db}
}

func (r *InsiderTradeRepo) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return []*entity.Transaction{{}}, nil
}

func (r *InsiderTradeRepo) Store(ctx context.Context, trade *entity.Transaction) error {
	return nil
}
