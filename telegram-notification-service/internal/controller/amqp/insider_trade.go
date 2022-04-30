package amqp

import "github.com/alexbobkovv/insider-trades/pkg/rabbitmq"

type handler struct {
	rmq *rabbitmq.RabbitMQ
}

func New(rabbitMQ *rabbitmq.RabbitMQ) *handler {
	return &handler{rabbitMQ}
}
