package monitor_test

import (
	"github.com/inteleon/go-amqp/monitor"
	"github.com/inteleon/go-logging/helper"
	"github.com/inteleon/go-logging/logging"
	"testing"
	"time"
)

func testMonitor(r *monitor.RabbitMQMonitor) {
	for {
		if r.Quit {
			r.Log.Debug("Quitting")

			return
		}

		time.Sleep(1 * time.Second)
	}
}

func TestStartAndStop(t *testing.T) {
	l, w := helper.NewTestLogging()

	m := &monitor.RabbitMQMonitor{
		Log:         l,
		Quit:        false,
		MonitorFunc: testMonitor,
	}

	m.Start()

	if m.Quit {
		t.Fatal("expected", false, "got", true)
	}

	m.Stop()

	time.Sleep(2 * time.Second)

	if !m.Quit {
		t.Fatal("expected", true, "got", false)
	}

	bufLen := len(w.Buffer)
	if bufLen != 1 {
		t.Fatal("expected", 1, "got", bufLen)
	}

	helper.ValidateLogEntry(
		t,
		w.Buffer[0],
		logging.DebugLogLevel,
		"Quitting",
	)
}
