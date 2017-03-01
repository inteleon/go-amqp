package amqp_test

import (
	"fmt"
	"github.com/inteleon/go-amqp/amqp"
	"github.com/inteleon/go-logging/helper"
	"github.com/inteleon/go-logging/logging"
	"reflect"
	"testing"
)

type rabbitMQClientConnectSuccessful struct{}

func (r *rabbitMQClientConnectSuccessful) Connect() (err error) {
	return
}

type rabbitMQClientCloseSuccessful struct{}

func (r *rabbitMQClientCloseSuccessful) Close() (err error) {
	return
}

type rabbitMQClientConnectError struct{}

func (r *rabbitMQClientConnectError) Connect() error {
	return fmt.Errorf("Connect error")
}

type rabbitMQClientCloseError struct{}

func (r *rabbitMQClientCloseError) Close() error {
	return fmt.Errorf("Close error")
}

type rabbitMQClientPublishSuccessful struct {
	t                         *testing.T
	ExpectedPublishRoutingKey string
	ExpectedPublishPayload    []byte
}

func (r *rabbitMQClientPublishSuccessful) Publish(routingKey string, payload []byte) (err error) {
	if routingKey != r.ExpectedPublishRoutingKey {
		r.t.Fatal("expected", r.ExpectedPublishRoutingKey, "got", routingKey)
	}

	if !reflect.DeepEqual(payload, r.ExpectedPublishPayload) {
		r.t.Fatal("expected", string(r.ExpectedPublishPayload), "got", string(payload))
	}

	return
}

type rabbitMQClientPublishError struct{}

func (r *rabbitMQClientPublishError) Publish(routingKey string, payload []byte) error {
	return fmt.Errorf("Publish error")
}

type rabbitMQDialSuccessful struct {
	Client amqp.AMQPClient
}

func (d *rabbitMQDialSuccessful) Dial() (amqp.AMQPClient, error) {
	return d.Client, nil
}

type rabbitMQDialError struct{}

func (d *rabbitMQDialError) Dial() (amqp.AMQPClient, error) {
	return nil, fmt.Errorf("Dial error")
}

func simpleRabbitMQConfig() *amqp.RabbitMQConfig {
	return &amqp.RabbitMQConfig{
		URL: "hax://haxor.yo",
	}
}

func TestRabbitMQClientAlreadyConnectedFailure(t *testing.T) {
	l, w := helper.NewTestLogging()

	r := &amqp.RabbitMQ{
		Log:        l,
		Connecting: true,
	}

	connect := r.Connect()
	if connect != nil {
		t.Fatal("expected", nil, "got", connect)
	}

	logsLen := len(w.Buffer)
	if logsLen != 1 {
		t.Fatal("expected", 1, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.WarningLogLevel,
		"Connect: Connecting is in progress",
	)
}

func TestRabbitMQDialDialFailure(t *testing.T) {
	l, w := helper.NewTestLogging()

	cfg := simpleRabbitMQConfig()
	r := &amqp.RabbitMQ{
		Cfg:        cfg,
		Dial:       &rabbitMQDialError{},
		Log:        l,
		Connecting: false,
	}

	connect := r.Connect()
	expErr := fmt.Errorf("Dial error")
	if !reflect.DeepEqual(connect, expErr) {
		t.Fatal("expected", expErr, "got", connect)
	}

	if r.Connecting {
		t.Fatal("expected", false, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 2 {
		t.Fatal("expected", 2, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.InfoLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[1],
		logging.ErrorLogLevel,
		"Unable to establish a connection to RabbitMQ: Dial error",
	)
}

func TestRabbitMQClientConnectFailure(t *testing.T) {
	l, w := helper.NewTestLogging()

	var client struct {
		rabbitMQClientConnectError
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	cfg := simpleRabbitMQConfig()
	r := &amqp.RabbitMQ{
		Cfg: cfg,
		Dial: &rabbitMQDialSuccessful{
			Client: &client,
		},
		Log:        l,
		Connecting: false,
	}

	connect := r.Connect()
	expErr := fmt.Errorf("Connect error")
	if !reflect.DeepEqual(connect, expErr) {
		t.Fatal("expected", expErr, "got", connect)
	}

	if r.Connecting {
		t.Fatal("expected", false, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 2 {
		t.Fatal("expected", 2, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.InfoLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[1],
		logging.ErrorLogLevel,
		"Unable to establish a connection to RabbitMQ: Connect error",
	)
}

func TestRabbitMQClientConnectSuccess(t *testing.T) {
	l, w := helper.NewTestLogging()

	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	cfg := simpleRabbitMQConfig()
	r := &amqp.RabbitMQ{
		Cfg: cfg,
		Dial: &rabbitMQDialSuccessful{
			Client: &client,
		},
		Log:        l,
		Connecting: false,
	}

	connect := r.Connect()
	if connect != nil {
		t.Fatal("expected", nil, "got", connect)
	}

	if r.Connecting {
		t.Fatal("expected", false, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 1 {
		t.Fatal("expected", 1, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.InfoLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)
}

func TestRabbitMQClientReconnectConnectionInProgressFailure(t *testing.T) {
	l, w := helper.NewTestLogging()

	r := &amqp.RabbitMQ{
		Log:        l,
		Connecting: true,
	}

	reconnect := r.Reconnect()
	if reconnect != nil {
		t.Fatal("expected", nil, "got", reconnect)
	}

	if !r.Connecting {
		t.Fatal("expected", true, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 1 {
		t.Fatal("expected", 1, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.WarningLogLevel,
		"Reconnect: Connecting is in progress",
	)
}

func TestRabbitMQClientReconnectSuccess(t *testing.T) {
	l, w := helper.NewTestLogging()

	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	cfg := simpleRabbitMQConfig()
	r := &amqp.RabbitMQ{
		Cfg: cfg,
		Dial: &rabbitMQDialSuccessful{
			Client: &client,
		},
		Log:        l,
		Connecting: false,
	}

	reconnect := r.Reconnect()
	if reconnect != nil {
		t.Fatal("expected", nil, "got", reconnect)
	}

	if r.Connecting {
		t.Fatal("expected", false, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 2 {
		t.Fatal("expected", 2, "got", logsLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.InfoLogLevel,
		"Reconnect: Connecting...",
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[1],
		logging.InfoLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)
}

func TestRabbitMQClientPublishFailure(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishError
	}

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	publish := r.Publish("hax", []byte("haxor"))
	expectedErr := fmt.Errorf("Publish error")
	if !reflect.DeepEqual(publish, expectedErr) {
		t.Fatal("expected", expectedErr, "got", publish)
	}
}

func TestRabbitMQClientPublishSuccess(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	client.t = t
	client.ExpectedPublishRoutingKey = "hax"
	client.ExpectedPublishPayload = []byte("haxor")

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	publish := r.Publish("hax", []byte("haxor"))
	if publish != nil {
		t.Fatal("expected", nil, "got", publish)
	}
}

func TestRabbitMQClientPingFailure(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishError
	}

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	ping := r.Ping()
	expectedErr := fmt.Errorf("Publish error")
	if !reflect.DeepEqual(ping, expectedErr) {
		t.Fatal("expected", expectedErr, "got", ping)
	}
}

func TestRabbitMQClientPingSuccess(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	client.t = t
	client.ExpectedPublishRoutingKey = "ping"
	client.ExpectedPublishPayload = []byte("ping")

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	ping := r.Ping()
	if ping != nil {
		t.Fatal("expected", nil, "got", ping)
	}
}

func TestRabbitMQClientCloseFailure(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseError
		rabbitMQClientPublishSuccessful
	}

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	close := r.Close()
	expectedErr := fmt.Errorf("Close error")
	if !reflect.DeepEqual(close, expectedErr) {
		t.Fatal("expected", expectedErr, "got", close)
	}
}

func TestRabbitMQClientCloseSuccess(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
	}

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	close := r.Close()
	if close != nil {
		t.Fatal("expected", nil, "got", close)
	}
}
