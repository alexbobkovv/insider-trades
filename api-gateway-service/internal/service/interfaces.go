package service

import (
	"context"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
)

type (
	Gateway interface {
		ListTrades(ctx context.Context, cursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error)
	}

	Cache interface {
		ListTrades(ctx context.Context, reqCursor *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error)
		AddTrades(ctx context.Context, trades []*api.TradeViewResponse)
	}
)
