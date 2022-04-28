package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// New TODO refactor errors, configs
func NewServer(url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	err = ch.ExchangeDeclare(
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

	q, err := ch.QueueDeclare(
		"telegram_channel_queue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"trades",
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("amqp New: %w", err)
	}

	return ch, nil
}
