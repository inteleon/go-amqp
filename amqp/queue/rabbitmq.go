package queue

import (
	"fmt"
	"github.com/inteleon/go-logging/logging"
	"github.com/satori/go.uuid"
	aq "github.com/streadway/amqp"
)

// RabbitMQContext is a RabbitMQ context.
// It holds the RabbitMQ consumer object reference and a delivery payload.
type RabbitMQContext struct {
	Consumer *RabbitMQConsumer
	Delivery AMQPDelivery
}

// RabbitMQDelivery implements the AMQPDelivery interface for returning RabbitMQ specific delivery data (payload, etc.).
type RabbitMQDelivery struct {
	Delivery aq.Delivery
}

// Payload returns the RabbitMQ delivery payload.
func (d *RabbitMQDelivery) Payload() ([]byte, error) {
	return d.Delivery.Body, nil
}

// Ack acks the message.
func (d *RabbitMQDelivery) Ack(multiple bool) error {
	return d.Delivery.Ack(multiple)
}

// Nack nacks the message.
func (d *RabbitMQDelivery) Nack(multiple, requeue bool) error {
	return d.Delivery.Nack(multiple, requeue)
}

// RabbitMQQueue defines a single queue that we will connect to (and declare if needed).
type RabbitMQQueue struct {
	Name        string
	Exchange    *RabbitMQExchange
	Durable     bool
	AutoDelete  bool
	Exclusive   bool
	NoWait      bool
	SkipDeclare bool
	ProcessFunc ProcessFunc
	AutoDLQ     bool
}

type RabbitMQExchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	AutoDLE    bool
}

type rabbitMQConsumerChannel struct {
	queue        *RabbitMQQueue
	rabbitMQChan *aq.Channel
	consumerName string
}

func (rc *rabbitMQConsumerChannel) Start() (<-chan aq.Delivery, error) {
	rc.consumerName = uuid.NewV4().String()

	return rc.rabbitMQChan.Consume(
		rc.queue.Name,
		rc.consumerName,
		false, // auto-ack (no-ack)
		rc.queue.Exclusive,
		false, // no-local
		rc.queue.NoWait,
		nil, // extra args
	)
}

func (rc *rabbitMQConsumerChannel) Stop() error {
	return rc.rabbitMQChan.Cancel(
		rc.consumerName,
		false, // no-wait
	)
}

// ProcessFunc is the function used to process every incoming message.
type ProcessFunc func(RabbitMQContext)

// RabbitMQConsumer is the consumer struct for RabbitMQ.
type RabbitMQConsumer struct {
	RabbitMQChan    *aq.Channel
	Queue           *RabbitMQQueue
	Log             logging.Logging
	DeliveryChan    chan aq.Delivery
	ConsumerChannel AMQPConsumerChannel
}

// RabbitMQConsumerOptions should be populated and passed to the NewRabbitMQConsumer function when called.
type RabbitMQConsumerOptions struct {
	RabbitMQChan *aq.Channel
	Queue        *RabbitMQQueue
	Log          logging.Logging
}

// NewRabbitMQConsumer is a helper function for creating a new RabbitMQConsumer.
func NewRabbitMQConsumer(options RabbitMQConsumerOptions) AMQPConsumer {
	return &RabbitMQConsumer{
		RabbitMQChan: options.RabbitMQChan,
		Queue:        options.Queue,
		Log:          options.Log,
		ConsumerChannel: &rabbitMQConsumerChannel{
			queue:        options.Queue,
			rabbitMQChan: options.RabbitMQChan,
		},
	}
}

// Start starts up the consumer and starts consuming the queue.
func (r *RabbitMQConsumer) Start() error {
	r.Log.Info(fmt.Sprintf("Starting up a consumer for the following queue: %s", r.Queue.Name))

	d, err := r.ConsumerChannel.Start()

	if err != nil {
		return err
	}

	go func(r *RabbitMQConsumer, d <-chan aq.Delivery) {
		for delivery := range d {
			go r.Queue.ProcessFunc(
				RabbitMQContext{
					Consumer: r,
					Delivery: &RabbitMQDelivery{
						Delivery: delivery,
					},
				},
			)
		}

		r.Log.Warn("The delivery channel has been closed! Exiting...")
	}(r, d)

	return nil
}

// Stop shuts down the consumer, stops consuming the queue.
func (r *RabbitMQConsumer) Stop() error {
	r.Log.Info(fmt.Sprintf("Shutting down a consumer for the following queue: %s", r.Queue.Name))

	return r.ConsumerChannel.Stop()
}
