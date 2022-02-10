package message

import (
	"context"
	"insidertradesreceiver/internal/entity"
	"insidertradesreceiver/pkg/kafka"
)

type InsiderTradePublisher struct {

}

func New(broker *kafka.Kafka) *InsiderTradePublisher {
	return &InsiderTradePublisher{}
}

func (p *InsiderTradePublisher) Publish(ctx context.Context, trade *entity.InsiderTrade) error {
	return nil
}
