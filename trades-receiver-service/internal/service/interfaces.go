package service

import (
	"context"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

type (
	InsiderTrade interface {
		Receive(ctx context.Context, trade *entity.Trade) error
		GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
	}

	InsiderTradeRepo interface {
		StoreTrade(ctx context.Context, trade *entity.Trade) error
		GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
	}

	InsiderTradePublisher interface {
		Publish(ctx context.Context, trade *entity.Trade) error
	}
)
