package queue

import (
	aq "github.com/streadway/amqp"
)

type AMQPDelivery interface {
	Payload() ([]byte, error)
	Ack(bool) error
	Nack(bool, bool) error
	SetHeader(key, value string) error
	GetHeader(key string) (string, error)
}

type AMQPConsumer interface {
	Start() error
	Stop() error
}

type AMQPConsumerChannel interface {
	Start() (<-chan aq.Delivery, error)
	Stop() error
}
