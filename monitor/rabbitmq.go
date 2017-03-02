package monitor

import (
	"github.com/inteleon/go-amqp/amqp"
	"github.com/inteleon/go-amqp/monitor/checker"
	"github.com/inteleon/go-logging/logging"
	"time"
)

func defaultMonitor(r *RabbitMQMonitor) {
	// Initiate the checker.
	c := checker.NewRabbitMQChecker(
		checker.RabbitMQCheckerOptions{
			Conn: r.Conn,
			Log:  r.Log,
		},
	)

	// Monitor the RabbitMQ connection.
	for {
		if r.Quit {
			return
		}

		// Probe the service and reconnect if needed.
		if err := c.Check(); err != nil {
			r.Log.Warn("Reconnecting to RabbitMQ...")

			r.Conn.Reconnect()
		}

		// Sleep 5 second between iterations.
		time.Sleep(5 * time.Second)
	}
}

// RabbitMQMonitor is the RabbitMQ connection monitor.
type RabbitMQMonitor struct {
	Conn        amqp.AMQP
	Log         logging.Logging
	Quit        bool
	MonitorFunc func(*RabbitMQMonitor)
}

// RabbitMQMonitorOptions provides option values for a new instance of RabbitMQMonitor to use.
type RabbitMQMonitorOptions struct {
	Conn amqp.AMQP
	Log  logging.Logging
}

// NewRabbitMQMonitor creates and returns a new RabbitMQMonitor instance.
func NewRabbitMQMonitor(options RabbitMQMonitorOptions) Monitor {
	return &RabbitMQMonitor{
		Conn:        options.Conn,
		Log:         options.Log,
		Quit:        false,
		MonitorFunc: defaultMonitor,
	}
}

// Start starts up the RabbitMQ monitoring.
func (r *RabbitMQMonitor) Start() {
	// Make sure we don't stop the monitor prematurely.
	r.Quit = false

	// Start up the monitor go routine.
	go r.MonitorFunc(r)
}

// Stop shuts down the RabbitMQ monitor.
func (r *RabbitMQMonitor) Stop() {
	r.Quit = true
}
