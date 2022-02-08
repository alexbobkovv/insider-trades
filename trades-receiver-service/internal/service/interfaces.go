package service

import (
	"context"
	"insidertradesreceiver/internal/entity"
)

type (
	InsiderTrade interface {
		Receive(ctx context.Context, trade *entity.InsiderTrade) error
		GetAll(ctx context.Context, limit, offset int) ([]*entity.InsiderTrade, error)
	}

	InsiderTradeRepo interface {
		Store(ctx context.Context, trade *entity.InsiderTrade) error
		GetAll(ctx context.Context, limit, offset int) ([]*entity.InsiderTrade, error)
	}

	InsiderTradePublisher interface {
		Publish(ctx context.Context, trade *entity.InsiderTrade) error
	}
)