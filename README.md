# go-amqp
AMQP clients.

## Usage

### Setup example
Simple usage example (typically on main thread) where we declare a queue with a companion exchange, including creation of DLQ and DLE:


    // Declare a list of queues, example here inlines a queue.RabbitMQQueue with an exchange
    queues := []queue.RabbitMQQueue{{
        Name: "parktura.deliveryreport.queue",
        Exchange: &queue.RabbitMQExchange{
            Name:       "parktura.exchange",
            Kind:       "direct",
            Durable:    true,
            AutoDLE:    true,
        },
        SkipDeclare: false,
        ProcessFunc: func(ctx queue.RabbitMQContext) {
            // Function body goes here, do stuff with the ctx.Delivery
        },
        AutoDLQ:     true,
    }}

    // Creates the logger transport instance with our config and queues
    amqpLogger, _ := logging.NewLogrusLogging(logging.LogrusLoggingOptions{})
    amqpTransport := amqp.NewRabbitMQ(
        amqp.RabbitMQConfig{
            URL:    "amqp://username:password@hostnametorabbitmq:5672/",
            Queues: queues,
            Log:    amqpLogger,
        },
    )
    
    // Connect to the broker
    if err := amqpTransport.Connect(); err != nil {
        panic(err)
    }
    // Make sure we clean up when main exits.
    defer amqpTransport.Close()
    
    // Start the consumer. A new go routine is started up for each and every queue.
    if err := amqpTransport.Consume(); err != nil {
        panic(err)
    }
    
    // Pass the amqpTransport to structs, write message to it etc...

The example above creates the following entries in RabbitMQ:

##### Queues
![queues](docs/images/queues.png)

##### Exchanges

![exchanges](docs/images/exchanges.png)

### Producer example
Producing messages is very simple, just use the _Publish(routingKey string, payload []byte)_ method:

    err := amqpTransport.Publish("parktura.deliveryreport.queue", []byte("{\"msg\":\"My JSON message!\"})
    if err != nil {
        return errors.Wrap("unable to publish message", err)
    }

### Consumer example
Consuming is also easy. We passed a _ProcessFunc_ to the queue declaration in the setup stage. Any func with the signature

    func(ctx queue.RabbitMQContext)

can be used as a consume function. Simple example with ACK / NACK:

    func DumpMessageToStdout(ctx queue.RabbitMQContext) {
        
        // Extract payload
        payload, err := ctx.Delivery.Payload()
        if err != nil {
            _ := ctx.Delivery.Nack(false, false)
            fmt.Printf("unable to extract payload from Delivery: %v", err)
            return
        }
        fmt.Printf("Got payload: %v\n", string(payload))
        
        // Ack it...
        if err = ctx.Delivery.Ack(false); err != nil {
            fmt.Printf("unable to ACK delivery from queue %v: %v", ctx.Consumer.Queue.Name, err)
        }
    }

## Mocking this library
The project provides mocks for its interfaces using [gomock](https://github.com/gomock).

Gomock and its mockgen tool can be installed using the following commands:

    go get github.com/golang/mock/gomock
    go install github.com/golang/mock/mockgen

Mocks can be re-generated using the mockgen tool.

    mockgen -source amqp/amqp.go -destination amqp/mock_amqp/mock_amqp.go -package mock_amqp
    mockgen -source amqp/queue/queue.go -destination amqp/queue/mock_queue/mock_queue.go -package mock_queue
    
Until a bug in gomock is fixed, one needs to manually modify mock_amqp.go:

Replace the relative import for _x "."_ with an absolute one:

    x "github.com/inteleon/go-amqp/amqp"
    
    
Code using this library can then easily mock our interfaces, for example the AMQPClient:

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockedClient := mock_amqp.NewMockAMQPClient(ctrl)
    mockedClient.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(1)
    
    // rest of test code...
    
The mock will now expect a single call to Publish with any arguments and it will return nil (no error) when this call occurs.