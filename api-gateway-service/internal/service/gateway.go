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
	cache    Cache
}

func New(tradeService api.TradeServiceClient, tradesCache Cache) *gatewayService {
	return &gatewayService{receiver: tradeService, cache: tradesCache}
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

	var viewsResponse []*api.TradeViewResponse

	for idx := 0; idx < int(limit); idx++ {
		view, err := views.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("%s: failed to receive a view: %w", methodName, err)
		}

		viewsResponse = append(viewsResponse, view)

	}

	nextCursor := cursor.NewEmpty()

	if len(viewsResponse) > 0 && viewsResponse[0] != nil {
		lastView := viewsResponse[len(viewsResponse)-1]
		if lastView != nil && lastView.CreatedAt != nil {
			createdAtTime := lastView.CreatedAt.AsTime()
			nextCursor = cursor.NewFromTime(&createdAtTime)
		}

		go s.cache.AddTrades(context.Background(), viewsResponse)

	} else {
		viewsResponse = []*api.TradeViewResponse{}
	}

	return viewsResponse, nextCursor, nil
}
