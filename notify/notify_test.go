package notify_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/notify"
)

func fixedMsg(level notify.Level, title, body string) notify.Message {
	return notify.Message{
		Level: level,
		Title: title,
		Body:  body,
		Time:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestSendWritesLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(&buf)
	_ = n.Send(fixedMsg(notify.LevelAlert, "port change", "port 22 opened"))
	if !strings.Contains(buf.String(), "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", buf.String())
	}
}

func TestSendWritesTitle(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(&buf)
	_ = n.Send(fixedMsg(notify.LevelInfo, "scan complete", "no changes"))
	if !strings.Contains(buf.String(), "scan complete") {
		t.Errorf("expected title in output, got: %s", buf.String())
	}
}

func TestSendWritesBody(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(&buf)
	_ = n.Send(fixedMsg(notify.LevelWarn, "title", "port 443 closed"))
	if !strings.Contains(buf.String(), "port 443 closed") {
		t.Errorf("expected body in output, got: %s", buf.String())
	}
}

func TestSendUsesProvidedTime(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(&buf)
	_ = n.Send(fixedMsg(notify.LevelInfo, "t", "b"))
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Errorf("expected date in output, got: %s", buf.String())
	}
}

func TestSendFallsBackToNowWhenZeroTime(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(&buf)
	msg := notify.Message{Level: notify.LevelInfo, Title: "t", Body: "b"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty")
	}
}
