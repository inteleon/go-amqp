package queue

import (
	aq "github.com/streadway/amqp"
)

type AMQPConsumer interface {
	Start() error
	Stop() error
}

type AMQPConsumerChannel interface {
	Start() (<-chan aq.Delivery, error)
	Stop() error
}
