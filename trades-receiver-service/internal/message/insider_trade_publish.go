package message

import (
	"context"

	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/kafka"
)

type InsiderTradePublisher struct {
	broker *kafka.Kafka
}

func New(broker *kafka.Kafka) *InsiderTradePublisher {
	return &InsiderTradePublisher{broker}
}

func (p *InsiderTradePublisher) Publish(ctx context.Context, trade *entity.Transaction) error {
	return nil
}
