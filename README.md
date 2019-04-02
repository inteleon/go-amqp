# go-amqp
AMQP clients.

### Mocking this library
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