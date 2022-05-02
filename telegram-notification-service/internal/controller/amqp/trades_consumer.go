package amqp

import (
	"fmt"

	"github.com/alexbobkovv/insider-trades/pkg/rabbitmq"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/config"
	"github.com/alexbobkovv/insider-trades/telegram-notification-service/internal/service"
	"github.com/alexbobkovv/insider-trades/trades-receiver-service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	rmq    *rabbitmq.RabbitMQ
	rmqCfg *config.RabbitMQ
	s      *service.Service
	l      *logger.Logger
}

func New(rabbitMQ *rabbitmq.RabbitMQ, rmqCfg *config.RabbitMQ, service *service.Service, logger *logger.Logger) (*Consumer, error) {

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
		c.l.Info(msg.Body)
	}

	return nil
}
