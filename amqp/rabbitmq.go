package amqp

import (
	"fmt"

	"github.com/inteleon/go-amqp/amqp/queue"
	"github.com/inteleon/go-logging/logging"
	aq "github.com/streadway/amqp"
)

// RabbitMQConfig is a configuration structure for the RabbitMQ client.
type RabbitMQConfig struct {
	URL    string
	Queues []queue.RabbitMQQueue
	Log    logging.Logging
}

// RabbitMQClient is the low level client used when talking with the RabbitMQ service.
type RabbitMQClient struct {
	Cfg     RabbitMQConfig
	Conn    *aq.Connection
	Channel *aq.Channel
	Log     logging.Logging
}

// Connect takes care of "on connect" specific tasks. Queue(s) and Exchange(s) may be created on the fly given the queues
// specified in the slice of RabbitMQQueue.
//
// q.Exchange - if not nil and SkipDeclare is false, an exchange will be created. If AutoDLE is true, a corresponding
// dead-letter exchange is created as well.
// q.AutoDLQ - if true, a DLQ with bindings will be created for the queue specified in q.Name. If AutoDLQ is true, a
// corresponding dead-letter queue is created as well.
func (c *RabbitMQClient) Connect() error {
	for _, q := range c.Cfg.Queues {
		if q.SkipDeclare {
			c.Log.Info(fmt.Sprintf("Skipping declaration of queue: %s", q.Name))
			continue
		}

		// Declare exchange and DLE if an exchange is specified and SkipDeclare is false
		if q.Exchange != nil && !q.Exchange.SkipDeclare {
			err := c.declareExchange(q.Exchange)
			if err != nil {
				return err
			}
		}

		// If AutoDLQ is true, provision DLQ for this queue
		if q.AutoDLQ {
			err := c.declareQueueWithDLQ(q)
			if err != nil {
				return err
			}
		} else {
			err := c.declareQueue(q)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// declareExchange creates a direct exchange with the name specified by the "Exchange" field as well as a dead-letter
// exchange using the exchange name as prefix.
func (c *RabbitMQClient) declareExchange(e *queue.RabbitMQExchange) error {
	err := c.Channel.ExchangeDeclare(e.Name, e.Kind, e.Durable, e.AutoDelete, e.Internal, e.NoWait, nil)
	if err != nil {
		return err
	}

	if e.AutoDLE {
		// Exchange for dead letters.
		dlx := e.Name + ".dead-letter"
		err = c.Channel.ExchangeDeclare(dlx, e.Kind, e.Durable, e.AutoDelete, e.Internal, e.NoWait, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// declareQueue declares a Queue with no arguments, optionally binding the new queue to an exchange if applicable.
func (c *RabbitMQClient) declareQueue(q queue.RabbitMQQueue) error {
	cq, err := c.Channel.QueueDeclare(
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

	// Bind to exchange if such is specified. Otherwise the queue declared above will bind to AMQP Default
	if q.Exchange != nil {
		err = c.Channel.QueueBind(cq.Name, q.Name, q.Exchange.Name, true, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// declareQueueWithDLQ creates a Queue with DLQ arguments as well as creating the DLQ and the necessary QueueBind(s).
func (c *RabbitMQClient) declareQueueWithDLQ(q queue.RabbitMQQueue) error {
	if q.Exchange == nil {
		return fmt.Errorf("exchange is a required parameter when AutoDLQ is true")
	}
	dlx := q.Exchange.Name + ".dead-letter"
	dlxq := q.Name + ".dead-letter"

	// Queue for dead letters.
	dlxQueue, err := c.Channel.QueueDeclare(dlxq, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Bind dead letter queue.
	err = c.Channel.QueueBind(dlxQueue.Name, dlxq, dlx, false, nil)
	if err != nil {
		return err
	}

	// Queue to consume messages from (with dlx).
	cq, err := c.Channel.QueueDeclare(q.Name, true, false, false, false, aq.Table{
		"x-dead-letter-exchange":    dlx,
		"x-dead-letter-routing-key": dlxQueue.Name,
	})
	if err != nil {
		return err
	}

	err = c.Channel.QueueBind(cq.Name, q.Name, q.Exchange.Name, true, nil)
	if err != nil {
		return err
	}
	return nil
}

// Close takes care of closing the connection to the RabbitMQ service.
func (c *RabbitMQClient) Close() error {
	return c.Conn.Close()
}

// Publish takes care of publishing a message to a single RabbitMQ queue.
func (c *RabbitMQClient) Publish(routingKey string, payload []byte) error {
	return c.publish(routingKey, "", payload, aq.Table{})
}

// PublishWithHeaders takes care of publishing a message to a single RabbitMQ queue.
func (c *RabbitMQClient) PublishWithHeaders(routingKey string, payload []byte, headers map[string]interface{}) error {
	return c.publish(routingKey, "", payload, headers)
}

// PublishOnExchange takes care of publishing a message to a single RabbitMQ exchange.
func (c *RabbitMQClient) PublishOnExchange(exchange string, payload []byte) error {
	return c.publish("", exchange, payload, aq.Table{})
}

// PublishOnExchangeWithHeaders takes care of publishing a message to a single RabbitMQ exchange.
func (c *RabbitMQClient) PublishOnExchangeWithHeaders(exchange string, payload []byte, headers map[string]interface{}) error {
	return c.publish("", exchange, payload, headers)
}

func (c *RabbitMQClient) publish(routingKey, exchange string, payload []byte, headers aq.Table) error {
	return c.Channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		aq.Publishing{
			Headers:         headers,
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            payload,
			DeliveryMode:    aq.Persistent,
			Priority:        0,
		},
	)
}

// Consume takes care of starting up queue consumers.
func (c *RabbitMQClient) Consume() error {
	for i, _ := range c.Cfg.Queues {
		cons := queue.NewRabbitMQConsumer(
			queue.RabbitMQConsumerOptions{
				RabbitMQChan: c.Channel,
				Queue:        &c.Cfg.Queues[i],
				Log:          c.Log,
			},
		)

		if err := cons.Start(); err != nil {
			return err
		}
	}

	return nil
}

// RabbitMQDial is used for establishing a connection to the RabbitMQ service.
type RabbitMQDial struct {
	Cfg RabbitMQConfig
	Log logging.Logging
}

// Dial establishes a connection to the RabbitMQ service and returns a RabbitMQClient.
func (d *RabbitMQDial) Dial() (AMQPClient, error) {
	conn, err := aq.Dial(d.Cfg.URL)
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
		Log:     d.Log,
	}, nil
}

// RabbitMQ is a RabbitMQ implementation of the AMQP interface.
type RabbitMQ struct {
	Cfg                RabbitMQConfig
	Log                logging.Logging
	Client             AMQPClient
	Dial               AMQPDial
	Connecting         bool
	ReconnectConsumers bool
}

// NewRabbitMQ creates and returns a new RabbitMQ object.
func NewRabbitMQ(cfg RabbitMQConfig) *RabbitMQ {
	var l logging.Logging

	if cfg.Log == nil {
		l, _ := logging.NewLogrusLogging(logging.LogrusLoggingOptions{})
		l.SetLogLevel(logging.DebugLogLevel)
	} else {
		l = cfg.Log
	}

	return &RabbitMQ{
		Cfg: cfg,
		Log: l,
		Dial: &RabbitMQDial{
			Cfg: cfg,
			Log: l,
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

	a.Log.Debug(fmt.Sprintf("Connecting to RabbitMQ at %s...", a.Cfg.URL))

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

// Publish takes care of dispatching messages to RabbitMQ using a routing key.
func (a *RabbitMQ) Publish(routingKey string, payload []byte) error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	return a.Client.Publish(routingKey, payload)
}

// PublishOnExchange takes care of dispatching messages to a named exchange on RabbitMQ.
func (a *RabbitMQ) PublishOnExchange(exchange string, payload []byte) error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	return a.Client.PublishOnExchange(exchange, payload)
}

// PublishWithHeaders takes care of dispatching messages to RabbitMQ using a routing key with custom headers.
func (a *RabbitMQ) PublishWithHeaders(routingKey string, payload []byte, headers map[string]interface{}) error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	return a.Client.PublishWithHeaders(routingKey, payload, headers)
}

// PublishOnExchangeWithHeaders takes care of dispatching messages to a named exchange on RabbitMQ with custom headers.
func (a *RabbitMQ) PublishOnExchangeWithHeaders(exchange string, payload []byte, headers map[string]interface{}) error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	return a.Client.PublishOnExchangeWithHeaders(exchange, payload, headers)
}

// Consume takes care of starting up queue consumers.
func (a *RabbitMQ) Consume() error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	a.ReconnectConsumers = true

	return a.Client.Consume()
}

// Ping is used for detecting dead RabbitMQ connections.
func (a *RabbitMQ) Ping() error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

	return a.Client.Publish("ping", []byte("ping"))
}

// Close shuts down the AMQP connection.
func (a *RabbitMQ) Close() error {
	if a.Client == nil {
		return a.clientDoesNotExist()
	}

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

	err = a.Connect()
	if err != nil {
		return
	}

	if a.ReconnectConsumers {
		a.Log.Info("Reconnect: Starting up consumers...")

		return a.Consume()
	}

	return
}

func (a *RabbitMQ) clientDoesNotExist() error {
	errMsg := "No available client!"
	a.Log.Error(errMsg)

	return fmt.Errorf(errMsg)
}
