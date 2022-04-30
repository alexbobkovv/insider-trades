package message

import (
	"context"
	"fmt"

	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InsiderTradePublisher struct {
	rmq *rabbitmq.RabbitMQ
}

func New(rabbitMQ *rabbitmq.RabbitMQ) (*InsiderTradePublisher, error) {

	err := rabbitMQ.Channel.ExchangeDeclare(
		"trades",
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	q, err := rabbitMQ.Channel.QueueDeclare(
		"telegram_channel_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	err = rabbitMQ.Channel.QueueBind(
		q.Name,
		"",
		"trades",
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	return &InsiderTradePublisher{rabbitMQ}, nil
}

func (p *InsiderTradePublisher) Publish(ctx context.Context, trade *entity.Trade) error {
	// TODO fix
	err := p.rmq.Channel.Publish(
		"trades",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("new trade"),
		})

	if err != nil {
		return fmt.Errorf("message: Publish: failed to publish message: %w", err)
	}

	return nil
}
