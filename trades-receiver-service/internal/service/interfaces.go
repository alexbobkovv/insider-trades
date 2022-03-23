//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package service
package service

import (
	"context"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

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
		Publish(ctx context.Context, trade *entity.Trade) error
	}
)
