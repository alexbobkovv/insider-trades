package service

import (
	"context"
	"fmt"
	"io"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
)

type gatewayService struct {
	receiver api.TradeServiceClient
}

func New(tradeService api.TradeServiceClient) *gatewayService {
	return &gatewayService{receiver: tradeService}
}

// TODO Cursor interface
func (s *gatewayService) ListTrades(ctx context.Context, crs *cursor.Cursor, limit uint32) ([]*api.TradeViewResponse, *cursor.Cursor, error) {
	const methodName = "(s *gatewayService) ListTrades"
	req := &api.TradeViewRequest{
		Cursor: crs.GetEncoded(),
		Limit:  limit,
	}
	views, err := s.receiver.ListViews(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", methodName, err)
	}

	viewsResponse := make([]*api.TradeViewResponse, limit)

	for idx := range viewsResponse {
		view, err := views.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("%s: failed to receive a view: %v", methodName, err)
		}

		viewsResponse[idx] = view

	}

	nextCursor := cursor.NewEmpty()

	if len(viewsResponse) > 0 {
		lastView := viewsResponse[len(viewsResponse)-1]
		createdAtTime := lastView.CreatedAt.AsTime()
		nextCursor = cursor.NewFromTime(&createdAtTime)
	}

	return viewsResponse, nextCursor, nil
}
