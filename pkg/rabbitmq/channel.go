package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func New(url string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{}
	err := rmq.connect(url)
	if err != nil {
		return nil, err
	}

	rmq.Channel, err = rmq.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	return rmq, nil
}

func (r *RabbitMQ) connect(url string) error {
	var connAttempts = 3
	const connTimeout = time.Second * 5
	var err error

	for connAttempts > 0 {
		r.Connection, err = amqp.Dial(url)
		if err == nil {
			break
		}

		log.Printf("amqp connect: failed to connect to rabbitmq, trying to reconnect.. error: %v", err)

		time.Sleep(connTimeout)
		connAttempts--
	}

	if err != nil {
		return fmt.Errorf("amqp connect: failed to connect to rabbitmq: %w", err)
	}

	return nil
}
