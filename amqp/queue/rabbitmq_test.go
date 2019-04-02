package queue_test

import (
	"fmt"
	"github.com/inteleon/go-amqp/amqp/queue"
	"github.com/inteleon/go-logging/helper"
	"github.com/inteleon/go-logging/logging"
	aq "github.com/streadway/amqp"
	"testing"
	"time"
)

type testConsumerChannel struct {
	deliveryChan chan aq.Delivery
}

func (cc *testConsumerChannel) Start() (<-chan aq.Delivery, error) {
	return cc.deliveryChan, nil
}

func (cc *testConsumerChannel) Stop() error {
	close(cc.deliveryChan)

	return nil
}

func testProcessFunc(p queue.RabbitMQContext) {
	payload, _ := p.Delivery.Payload()

	p.Consumer.Log.Debug(fmt.Sprintf("Consumer test process executed. Body: %s", string(payload)))
}

func TestFullFlow(t *testing.T) {
	l, w := helper.NewTestLogging()

	dChan := make(chan aq.Delivery)

	c := &queue.RabbitMQConsumer{
		Queue: &queue.RabbitMQQueue{
			Name:        "test",
			ProcessFunc: testProcessFunc,
		},
		Log: l,
		ConsumerChannel: &testConsumerChannel{
			deliveryChan: dChan,
		},
	}

	// Note ugly use of sleeps, makes sure log messages are deterministically ordered for the asserts.
	c.Start()
	time.Sleep(10 * time.Millisecond)
	dChan <- aq.Delivery{
		Body: []byte("1337"),
	}
	time.Sleep(10 * time.Millisecond)
	c.Stop()

	time.Sleep(500 * time.Millisecond)

	bufLen := len(w.Buffer)
	if bufLen != 4 {
		t.Fatal("expected", 4, "got", bufLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.InfoLogLevel,
		"Starting up a consumer for the following queue: test",
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[1],
		logging.DebugLogLevel,
		"Consumer test process executed. Body: 1337",
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[2],
		logging.InfoLogLevel,
		"Shutting down a consumer for the following queue: test",
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[3],
		logging.WarningLogLevel,
		"The delivery channel has been closed! Exiting...",
	)
}
