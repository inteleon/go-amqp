package amqp_test

import (
	"fmt"
	"github.com/inteleon/go-amqp/amqp"
	"github.com/inteleon/go-amqp/amqp/queue"
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
	ExpectedPublishExchange   string
	ExpectedPublishPayload    []byte
	ExpectedHeaders           map[string]interface{}
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
func (r *rabbitMQClientPublishSuccessful) PublishOnExchange(exchange string, payload []byte) (err error) {
	if exchange != r.ExpectedPublishExchange {
		r.t.Fatal("expected", r.ExpectedPublishExchange, "got", exchange)
	}

	if !reflect.DeepEqual(payload, r.ExpectedPublishPayload) {
		r.t.Fatal("expected", string(r.ExpectedPublishPayload), "got", string(payload))
	}

	return
}
func (r *rabbitMQClientPublishSuccessful) PublishWithHeaders(routingKey string, payload []byte, headers map[string]interface{}) (err error) {
	if routingKey != r.ExpectedPublishRoutingKey {
		r.t.Fatal("expected", r.ExpectedPublishRoutingKey, "got", routingKey)
	}

	if !reflect.DeepEqual(payload, r.ExpectedPublishPayload) {
		r.t.Fatal("expected", string(r.ExpectedPublishPayload), "got", string(payload))
	}

	return
}
func (r *rabbitMQClientPublishSuccessful) PublishOnExchangeWithHeaders(exchange string, payload []byte, headers map[string]interface{}) (err error) {
	if exchange != r.ExpectedPublishExchange {
		r.t.Fatal("expected", r.ExpectedPublishExchange, "got", exchange)
	}

	if !reflect.DeepEqual(payload, r.ExpectedPublishPayload) {
		r.t.Fatal("expected", string(r.ExpectedPublishPayload), "got", string(payload))
	}
	if fmt.Sprintf("%+v", headers) != fmt.Sprintf("%+v", r.ExpectedHeaders) {
		r.t.Fatal("expected headers", fmt.Sprintf("%+v", r.ExpectedHeaders), "got", fmt.Sprintf("%+v", headers))
	}

	return
}

type rabbitMQClientPublishError struct{}

func (r *rabbitMQClientPublishError) Publish(routingKey string, payload []byte) error {
	return fmt.Errorf("Publish error")
}
func (r *rabbitMQClientPublishError) PublishOnExchange(exchange string, payload []byte) error {
	return fmt.Errorf("Publish error")
}
func (r *rabbitMQClientPublishError) PublishWithHeaders(string, []byte, map[string]interface{}) error {
	return fmt.Errorf("Publish error")
}
func (r *rabbitMQClientPublishError) PublishOnExchangeWithHeaders(string, []byte, map[string]interface{}) error {
	return fmt.Errorf("Publish error")
}

type rabbitMQClientConsumeSuccessful struct{}

func (r *rabbitMQClientConsumeSuccessful) Consume() error {
	return nil
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

func simpleRabbitMQConfig() amqp.RabbitMQConfig {
	return amqp.RabbitMQConfig{
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
		logging.DebugLogLevel,
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
		rabbitMQClientConsumeSuccessful
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
		logging.DebugLogLevel,
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
		rabbitMQClientConsumeSuccessful
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
		logging.DebugLogLevel,
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
		rabbitMQClientConsumeSuccessful
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
		logging.DebugLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)
}

func TestRabbitMQClientReconnectReConsumeSuccess(t *testing.T) {
	l, w := helper.NewTestLogging()

	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
		rabbitMQClientConsumeSuccessful
	}

	cfg := simpleRabbitMQConfig()
	r := &amqp.RabbitMQ{
		Cfg: cfg,
		Dial: &rabbitMQDialSuccessful{
			Client: &client,
		},
		Log:                l,
		Connecting:         false,
		ReconnectConsumers: true,
	}

	reconnect := r.Reconnect()
	if reconnect != nil {
		t.Fatal("expected", nil, "got", reconnect)
	}

	if r.Connecting {
		t.Fatal("expected", false, "got", r.Connecting)
	}

	logsLen := len(w.Buffer)
	if logsLen != 3 {
		t.Fatal("expected", 3, "got", logsLen)
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
		logging.DebugLogLevel,
		fmt.Sprintf("Connecting to RabbitMQ at %s...", cfg.URL),
	)

	helper.ValidateLogEntry(
		t,
		w.Buffer[2],
		logging.InfoLogLevel,
		"Reconnect: Starting up consumers...",
	)
}

func TestRabbitMQClientPublishFailure(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishError
		rabbitMQClientConsumeSuccessful
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
		rabbitMQClientConsumeSuccessful
	}

	client.t = t
	client.ExpectedPublishRoutingKey = "hax"
	client.ExpectedPublishExchange = ""
	client.ExpectedPublishPayload = []byte("haxor")

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	publish := r.Publish("hax", []byte("haxor"))
	if publish != nil {
		t.Fatal("expected", nil, "got", publish)
	}
}

func TestRabbitMQClientPublishOnExchangeSuccess(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishSuccessful
		rabbitMQClientConsumeSuccessful
	}

	client.t = t
	client.ExpectedPublishRoutingKey = ""
	client.ExpectedPublishExchange = "some.exchange"
	client.ExpectedPublishPayload = []byte("haxor")

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	publish := r.PublishOnExchange("some.exchange", []byte("haxor"))
	if publish != nil {
		t.Fatal("expected", nil, "got", publish)
	}
}

func TestRabbitMQClientPingFailure(t *testing.T) {
	var client struct {
		rabbitMQClientConnectSuccessful
		rabbitMQClientCloseSuccessful
		rabbitMQClientPublishError
		rabbitMQClientConsumeSuccessful
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
		rabbitMQClientConsumeSuccessful
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
		rabbitMQClientConsumeSuccessful
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
		rabbitMQClientConsumeSuccessful
	}

	r := &amqp.RabbitMQ{
		Client: &client,
	}

	close := r.Close()
	if close != nil {
		t.Fatal("expected", nil, "got", close)
	}
}

func TestClientDoesNotExist(t *testing.T) {
	l, w := helper.NewTestLogging()

	r := &amqp.RabbitMQ{
		Log: l,
	}

	cons := r.Consume()

	if cons == nil {
		t.Fatal("expected", "error", "got", cons)
	}

	expectedErrMsg := "No available client!"
	if cons.Error() != expectedErrMsg {
		t.Fatal("expected", expectedErrMsg, "got", cons.Error())
	}

	bufLen := len(w.Buffer)
	if bufLen != 1 {
		t.Fatal("expected", 1, "got", bufLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.ErrorLogLevel,
		expectedErrMsg,
	)
}

// TDD tests. :-)

func TestQueueConsuming(t *testing.T) {
	l, w := helper.NewTestLogging()

	defer func(t *testing.T, l logging.Logging, w *helper.TestLogWriter) {
		if r := recover(); r == nil {
			t.Fatal("expected", "panic", "got", nil)
		}

		// The program panics before reaching the init of the "hax" queue consumer.
		// Therefore we only expect 1 log entry to have been written.
		bufLen := len(w.Buffer)
		if bufLen != 1 {
			t.Fatal("expected", 1, "got", bufLen)
		}

		helper.ValidateLogEntry(
			t,
			w.Buffer[0],
			logging.InfoLogLevel,
			"Starting up a consumer for the following queue: test",
		)

	}(t, l, w)

	client := &amqp.RabbitMQClient{
		Cfg: amqp.RabbitMQConfig{
			Queues: []queue.RabbitMQQueue{
				{
					Name: "test",
				},
				{
					Name: "hax",
				},
			},
		},
		Log: l,
	}

	client.Consume()
}

func TestRabbitMQClientConnectQueueSkipDeclareSuccess(t *testing.T) {
	l, w := helper.NewTestLogging()

	queues := []queue.RabbitMQQueue{
		{
			Name:        "test",
			SkipDeclare: true,
		},
		{
			Name:        "hax",
			SkipDeclare: true,
		},
	}

	client := &amqp.RabbitMQClient{
		Cfg: amqp.RabbitMQConfig{
			Queues: queues,
		},
		Log: l,
	}

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	logsLen := len(w.Buffer)
	if logsLen != 2 {
		t.Fatal("expected", 2, "got", logsLen)
	}

	for k, q := range queues {
		helper.ValidateLogEntry(
			t,
			w.Buffer[k],
			logging.InfoLogLevel,
			fmt.Sprintf("Skipping declaration of queue: %s", q.Name),
		)
	}
}
