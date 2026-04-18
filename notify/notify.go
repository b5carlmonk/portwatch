// Package notify provides notification backends for portwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds the data for a single notification.
type Message struct {
	Level   Level
	Title   string
	Body    string
	Time    time.Time
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(msg Message) error
}

// LogNotifier writes notifications to an io.Writer (default: os.Stderr).
type LogNotifier struct {
	Out io.Writer
}

// New returns a LogNotifier writing to os.Stderr.
func New() *LogNotifier {
	return &LogNotifier{Out: os.Stderr}
}

// NewWithWriter returns a LogNotifier writing to the given writer.
func NewWithWriter(w io.Writer) *LogNotifier {
	return &LogNotifier{Out: w}
}

// Send formats and writes the message to the configured writer.
func (l *LogNotifier) Send(msg Message) error {
	if msg.Time.IsZero() {
		msg.Time = time.Now()
	}
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s %s: %s\n",
		msg.Time.Format(time.RFC3339),
		msg.Level,
		msg.Title,
		msg.Body,
	)
	return err
}
