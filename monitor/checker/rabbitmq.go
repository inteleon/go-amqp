package checker

import (
	"fmt"
	"github.com/inteleon/go-amqp/amqp"
	"github.com/inteleon/go-logging/logging"
)

// RabbitMQChecker is the RabbitMQ healthchecker.
type RabbitMQChecker struct {
	Conn amqp.AMQP
	Log  logging.Logging
}

// RabbitMQCheckerOptions provides option values for a new instance of RabbitMQChecker to use.
type RabbitMQCheckerOptions struct {
	Conn amqp.AMQP
	Log  logging.Logging
}

// NewRabbitMQChecker creates and returns a new RabbitMQChecker instance.
func NewRabbitMQChecker(options RabbitMQCheckerOptions) Checker {
	return &RabbitMQChecker{
		Conn: options.Conn,
		Log:  options.Log,
	}
}

// Check is the healthcheck function. Handles the service probing.
func (r *RabbitMQChecker) Check() (err error) {
	if err = r.Conn.Ping(); err != nil {
		r.Log.Error(fmt.Sprintf("Unable to dispatch to RabbitMQ! Error: %s", err))
	}

	return
}
