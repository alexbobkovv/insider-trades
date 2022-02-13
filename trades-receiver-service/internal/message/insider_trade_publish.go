package message

import (
	"context"
	"insidertradesreceiver/internal/entity"
	"insidertradesreceiver/pkg/kafka"
)

type InsiderTradePublisher struct {
	broker *kafka.Kafka
}

func New(broker *kafka.Kafka) *InsiderTradePublisher {
	return &InsiderTradePublisher{broker}
}

func (p *InsiderTradePublisher) Publish(ctx context.Context, trade *entity.InsiderTrade) error {
	return nil
}
