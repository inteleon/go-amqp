package checker_test

import (
	"fmt"
	"github.com/inteleon/go-amqp/monitor/checker"
	"github.com/inteleon/go-logging/helper"
	"github.com/inteleon/go-logging/logging"
	"testing"
)

type testConn struct {
	pingRetVal error
}

func (c *testConn) Connect() error {
	return nil
}

func (c *testConn) Publish(routingKey string, payload []byte) error {
	return nil
}

func (c *testConn) PublishOnExchange(exchange string, payload []byte) error {
	return nil
}

func (c *testConn) PublishWithHeaders(string, []byte, map[string]interface{}) error {
	return nil
}

func (c *testConn) PublishOnExchangeWithHeaders(string, []byte, map[string]interface{}) error {
	return nil
}

func (c *testConn) Consume() error {
	return nil
}

func (c *testConn) Ping() error {
	return c.pingRetVal
}

func (c *testConn) Close() error {
	return nil
}

func (c *testConn) Reconnect() error {
	return nil
}

func TestCheckPingSuccess(t *testing.T) {
	l, w := helper.NewTestLogging()

	c := &checker.RabbitMQChecker{
		Conn: &testConn{
			pingRetVal: nil,
		},
		Log: l,
	}

	ping := c.Check()

	if ping != nil {
		t.Fatal("expected", nil, "got", ping)
	}

	bufLen := len(w.Buffer)
	if bufLen != 0 {
		t.Fatal("expected", 0, "got", bufLen)
	}
}

func TestCheckPingFailure(t *testing.T) {
	l, w := helper.NewTestLogging()

	errMsg := "ZOMG ERROR YOLO"

	c := &checker.RabbitMQChecker{
		Conn: &testConn{
			pingRetVal: fmt.Errorf(errMsg),
		},
		Log: l,
	}

	ping := c.Check()

	if ping == nil {
		t.Fatal("expected", "error", "got", ping)
	}

	if ping.Error() != errMsg {
		t.Fatal("expected", errMsg, "got", ping.Error())
	}

	bufLen := len(w.Buffer)
	if bufLen != 1 {
		t.Fatal("expected", 1, "got", bufLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.ErrorLogLevel,
		fmt.Sprintf("Unable to dispatch to RabbitMQ! Error: %s", errMsg),
	)
}
