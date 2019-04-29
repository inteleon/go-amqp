package amqp

// AMQP is the interface that every AMQP service needs to implement.
type AMQP interface {
	Connect() error
	Publish(string, []byte) error
	PublishOnExchange(string, []byte) error
	Consume() error
	Ping() error
	Close() error
	Reconnect() error
}

// AMQPClient is for low level AMQP client operations.
type AMQPClient interface {
	Connect() error
	Close() error
	Publish(string, []byte) error
	PublishOnExchange(string, []byte) error
	Consume() error
}

// AMQPDial is used for dialing in to an AMQP service.
type AMQPDial interface {
	Dial() (AMQPClient, error)
}
