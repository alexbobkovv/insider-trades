package service

import (
	"context"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package service_test
type (
	InsiderTrade interface {
		Receive(ctx context.Context, trade *entity.Trade) error
		ListTransactions(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*entity.Transaction, *cursor.Cursor, error)
		ListViews(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error)
	}

	InsiderTradeRepo interface {
		StoreTrade(ctx context.Context, trade *entity.Trade) error
		ListTransactions(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*entity.Transaction, *cursor.Cursor, error)
		ListViews(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error)
	}

	InsiderTradePublisher interface {
		PublishTrade(ctx context.Context, trade *entity.Trade) error
	}
)
