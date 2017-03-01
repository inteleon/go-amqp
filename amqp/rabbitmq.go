package amqp

import (
	"fmt"
	"github.com/inteleon/go-logging/logging"
	aq "github.com/streadway/amqp"
)

// RabbitMQConfig is a configuration structure for the RabbitMQ client.
type RabbitMQConfig struct {
	URL    string
	Queues []RabbitMQQueue
}

// RabbitMQQueue defines a single queue that we will connect to (and declare if needed).
type RabbitMQQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

// RabbitMQClient is the low level client used when talking with the RabbitMQ service.
type RabbitMQClient struct {
	Cfg     *RabbitMQConfig
	Conn    *aq.Connection
	Channel *aq.Channel
}

// Connect takes care of "on connect" specific tasks.
func (c *RabbitMQClient) Connect() error {
	for _, q := range c.Cfg.Queues {
		// See https://www.rabbitmq.com/amqp-0-9-1-reference.html for
		// more information about the arguments.
		_, err := c.Channel.QueueDeclare(
			q.Name,
			q.Durable,    // durable
			q.AutoDelete, // auto-delete
			q.Exclusive,  // exclusive
			q.NoWait,     // no-wait
			nil,          // arguments
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// Close takes care of closing the connection to the RabbitMQ service.
func (c *RabbitMQClient) Close() error {
	return c.Conn.Close()
}

// Publish takes care of publishing a message to a single RabbitMQ queue.
func (c *RabbitMQClient) Publish(routingKey string, payload []byte) error {
	return c.Channel.Publish(
		"",
		routingKey,
		false, // mandatory
		false, // immediate
		aq.Publishing{
			Headers:         aq.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            payload,
			DeliveryMode:    aq.Persistent,
			Priority:        0,
		},
	)
}

// RabbitMQDial is used for establishing a connection to the RabbitMQ service.
type RabbitMQDial struct {
	Cfg *RabbitMQConfig
	URL string
}

// Dial establishes a connection to the RabbitMQ service and returns a RabbitMQClient.
func (d *RabbitMQDial) Dial() (AMQPClient, error) {
	conn, err := aq.Dial(d.URL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		Conn:    conn,
		Channel: channel,
		Cfg:     d.Cfg,
	}, nil
}

// RabbitMQ is a RabbitMQ implementation of the AMQP interface.
type RabbitMQ struct {
	Cfg        *RabbitMQConfig
	Log        logging.Logging
	Client     AMQPClient
	Dial       AMQPDial
	Connecting bool
}

// NewRabbitMQ creates and returns a new RabbitMQ object.
func NewRabbitMQ(url string, cfg *RabbitMQConfig) *RabbitMQ {
	l := logging.NewLogrusLogging()
	l.SetLogLevel(logging.DebugLogLevel)

	return &RabbitMQ{
		Cfg: cfg,
		Log: l,
		Dial: &RabbitMQDial{
			Cfg: cfg,
			URL: url,
		},
		Connecting: false,
	}
}

// Connect takes care of creating and setting up an AMQP connection.
func (a *RabbitMQ) Connect() (err error) {
	// Connecting is in progress...
	if a.Connecting {
		a.Log.Warn("Connect: Connecting is in progress")

		return
	}

	// Just in case we don't try to reconnect while connecting.
	a.Connecting = true
	defer func(a *RabbitMQ) {
		a.Connecting = false
	}(a)

	a.Log.Info(fmt.Sprintf("Connecting to RabbitMQ at %s...", a.Cfg.URL))

	a.Client, err = a.Dial.Dial()
	if err != nil {
		a.Log.Error(fmt.Sprintf("Unable to establish a connection to RabbitMQ: %s", err))

		return
	}

	if err = a.Client.Connect(); err != nil {
		a.Log.Error(fmt.Sprintf("Unable to establish a connection to RabbitMQ: %s", err))

		return
	}

	return
}

// Publish takes care of dispatching messages to RabbitMQ.
func (a *RabbitMQ) Publish(routingKey string, payload []byte) error {
	return a.Client.Publish(routingKey, payload)
}

// Ping is used for detecting dead RabbitMQ connections.
func (a *RabbitMQ) Ping() error {
	return a.Client.Publish("ping", []byte("ping"))
}

// Close shuts down the AMQP connection.
func (a *RabbitMQ) Close() error {
	return a.Client.Close()
}

// Reconnect takes care of reconnecting to RabbitMQ.
func (a *RabbitMQ) Reconnect() (err error) {
	// Connecting is in progress...
	if a.Connecting {
		a.Log.Warn("Reconnect: Connecting is in progress")

		return
	}

	a.Log.Info("Reconnect: Connecting...")

	return a.Connect()
}