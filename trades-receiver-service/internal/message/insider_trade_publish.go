package message

import (
	"context"
	"fmt"
	"time"

	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/config"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/internal/entity"
	"github.com/gofrs/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InsiderTradePublisher struct {
	rmq    *rabbitmq.RabbitMQ
	rmqCfg config.RabbitMQ
}

func New(rabbitMQ *rabbitmq.RabbitMQ, rmqCfg config.RabbitMQ) (*InsiderTradePublisher, error) {

	err := rabbitMQ.Channel.ExchangeDeclare(
		rmqCfg.Exchange,
		amqp.ExchangeFanout,
		rmqCfg.Durable,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	q, err := rabbitMQ.Channel.QueueDeclare(
		rmqCfg.QueueName,
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
		rmqCfg.RoutingKey,
		rmqCfg.Exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	return &InsiderTradePublisher{rmq: rabbitMQ, rmqCfg: rmqCfg}, nil
}

func (p *InsiderTradePublisher) PublishTrade(ctx context.Context, trade *entity.Trade) error {
	// TODO fix
	// err := p.rmq.Channel.Publish(
	// 	p.rmqCfg.Exchange,
	// 	p.rmqCfg.RoutingKey,
	// 	false,
	// 	false,
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte("new trade"),
	// 	})
	//
	// if err != nil {
	// 	return fmt.Errorf("message: Publish: failed to publish message: %w", err)
	// }

	// TODO rewrite
	if err := p.publish([]byte("something")); err != nil {
		return fmt.Errorf("PublishTrade: %w", err)
	}

	return nil
}

func (p *InsiderTradePublisher) publish(body []byte) error {
	msgID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("publish: failed to generate a new uuid: %w", err)
	}

	msg := amqp.Publishing{
		MessageId:       msgID.String(),
		Timestamp:       time.Now(),
		DeliveryMode:    amqp.Persistent,
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            body,
	}

	if err := p.rmq.Channel.Publish(
		p.rmqCfg.Exchange,
		p.rmqCfg.RoutingKey,
		false, // mandatory
		false, // immediate
		msg,
	); err != nil {
		return fmt.Errorf("publish: failed to publish: %w", err)
	}

	return nil
}
