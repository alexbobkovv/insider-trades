package service

import (
	"context"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package service_test
type (
	InsiderTrade interface {
		Receive(ctx context.Context, trade *entity.Trade) error
		GetAll(ctx context.Context, cursor string, limit int) ([]*entity.Transaction, string, error)
	}

	InsiderTradeRepo interface {
		StoreTrade(ctx context.Context, trade *entity.Trade) error
		GetAll(ctx context.Context, cursor string, limit int) ([]*entity.Transaction, string, error)
	}

	InsiderTradePublisher interface {
		PublishTrade(ctx context.Context, trade *entity.Trade) error
	}
)
