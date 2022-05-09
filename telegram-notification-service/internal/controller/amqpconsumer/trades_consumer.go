package amqpconsumer

import (
	"fmt"

	"github.com/alexbobkovv/insider-trades/api"
	"github.com/alexbobkovv/insider-trades/pkg/logger"
	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	rmq    *rabbitmq.RabbitMQ
	rmqCfg *config.RabbitMQ
	s      service.Service
	l      *logger.Logger
}

func New(rabbitMQ *rabbitmq.RabbitMQ, rmqCfg *config.RabbitMQ, service service.Service, logger *logger.Logger) (*Consumer, error) {

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

	return &Consumer{rmq: rabbitMQ, rmqCfg: rmqCfg, s: service, l: logger}, nil
}

func (c *Consumer) Run() error {
	msgs, err := c.rmq.Channel.Consume(
		c.rmqCfg.QueueName,
		c.rmqCfg.ConsumerName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("consumer: run: %w", err)
	}

	for msg := range msgs {
		trade := &api.Trade{}
		if err := proto.Unmarshal(msg.Body, trade); err != nil {
			c.l.Errorf("consumer: Run: failed to unmarshal message to proto: %v", err)
		}

		if err := c.s.ProcessTrade(trade); err != nil {
			c.l.Errorf("consumer: Run: service error while processing trade: %v", err)
		}
	}

	return nil
}
