package grpcapi

import (
	"fmt"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/types/cursor"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type tradeServer struct {
	s   service.InsiderTrade
	l   *logger.Logger
	cfg *config.Config
	api.UnimplementedTradeServiceServer
}

func NewTradeServer(service service.InsiderTrade, logger *logger.Logger, cfg *config.Config) *tradeServer {
	return &tradeServer{s: service, l: logger, cfg: cfg}
}

func (h *tradeServer) ListTransactions(request *api.TradeRequest, stream api.TradeService_ListTransactionsServer) error {
	const methodName = "(h *tradeServer) ListTransactions"

	c, err := cursor.NewFromEncodedString(request.GetCursor())
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	limit := request.GetLimit()
	if limit <= 0 {
		const defaultLimit = 20
		limit = defaultLimit
	}

	transactions, nextCursor, err := h.s.ListTransactions(stream.Context(), c, limit)
	if err != nil {
		return err
	}

	meta := metadata.New(map[string]string{"nextCursor": nextCursor.GetEncoded()})
	if err := stream.SetHeader(meta); err != nil {
		return err
	}

	for _, transaction := range transactions {

		transactionProto := &api.Transaction{
			ID:                  transaction.ID,
			SecFilingsID:        transaction.SecFilingsID,
			TransactionTypeName: transaction.TransactionTypeName,
			AveragePrice:        transaction.AveragePrice,
			TotalShares:         transaction.TotalShares,
			TotalValue:          transaction.TotalValue,
			CreatedAt:           timestamppb.New(transaction.CreatedAt),
		}

		if err := stream.Send(transactionProto); err != nil {
			return err
		}
	}

	return nil
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
