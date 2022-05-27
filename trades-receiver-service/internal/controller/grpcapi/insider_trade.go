package grpcapi

import (
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"google.golang.org/grpc/metadata"
)

type tradeServer struct {
	s            service.InsiderTrade
	l            *logger.Logger
	cfg          *config.Config
	tradesStream map[api.TradeService_ListTradesServer][]*api.TradeRequest
	api.UnimplementedTradeServiceServer
}

func NewTradeServer(service service.InsiderTrade, logger *logger.Logger, cfg *config.Config) *tradeServer {
	return &tradeServer{s: service, l: logger, cfg: cfg, tradesStream: make(map[api.TradeService_ListTradesServer][]*api.TradeRequest)}
}

// TODO implement
func (h *tradeServer) ListTrades(request *api.TradeRequest, stream api.TradeService_ListTradesServer) error {
	for {
		if err := stream.Send(&api.Trade{
			Ins:  nil,
			Cmp:  nil,
			SecF: nil,
			Trs:  nil,
			Sth:  nil,
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 5)
	}
}

func (h *tradeServer) ListViews(request *api.TradeViewRequest, stream api.TradeService_ListViewsServer) error {
	const methodName = "(h *tradeServer) ListViews"
	c, err := cursor.NewFromEncodedString(request.GetCursor())
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	limit := request.GetLimit()
	if limit <= 0 {
		const defaultLimit = 20
		limit = defaultLimit
	}

	views, nextCursor, err := h.s.ListViews(stream.Context(), c, limit)
	if err != nil {
		return err
	}

	meta := metadata.New(map[string]string{"nextCursor": nextCursor.GetEncoded()})
	if err := stream.SetHeader(meta); err != nil {
		return err
	}

	for _, view := range views {

		if err := stream.Send(view); err != nil {
			return err
		}
	}

	return nil
}
