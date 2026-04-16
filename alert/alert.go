package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends alerts somewhere (stdout, file, etc.).
type Notifier struct {
	out io.Writer
}

// New creates a Notifier writing to the given writer.
// Pass nil to default to os.Stdout.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes alerts derived from a Diff result.
func (n *Notifier) Notify(diff scanner.DiffResult) {
	for _, r := range diff.Opened {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   fmt.Sprintf("Port OPENED: %s:%d (%s)", r.Host, r.Port, r.Service),
		}
		n.write(a)
	}
	for _, r := range diff.Closed {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   fmt.Sprintf("Port CLOSED: %s:%d (%s)", r.Host, r.Port, r.Service),
		}
		n.write(a)
	}
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
}
