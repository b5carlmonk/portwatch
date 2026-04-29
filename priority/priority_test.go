package priority_test

import (
	"testing"

	"github.com/user/portwatch/priority"
	"github.com/user/portwatch/scanner"
)

func makeResult(port int, proto string) scanner.Result {
	return scanner.Result{Host: "127.0.0.1", Port: port, Proto: proto, Open: true}
}

func TestScoreMatchesExactRule(t *testing.T) {
	s := priority.New([]priority.Rule{
		{Port: 22, Proto: "tcp", Level: priority.Critical},
	}, priority.Low)

	if got := s.Score(makeResult(22, "tcp")); got != priority.Critical {
		t.Fatalf("expected CRITICAL, got %s", got)
	}
}

func TestScoreFallsBackToDefault(t *testing.T) {
	s := priority.New(nil, priority.Medium)
	if got := s.Score(makeResult(9999, "tcp")); got != priority.Medium {
		t.Fatalf("expected MEDIUM, got %s", got)
	}
}

func TestScoreAnyProtoMatchesBoth(t *testing.T) {
	s := priority.New([]priority.Rule{
		{Port: 53, Proto: "", Level: priority.High},
	}, priority.Low)

	if got := s.Score(makeResult(53, "tcp")); got != priority.High {
		t.Fatalf("expected HIGH for tcp, got %s", got)
	}
	if got := s.Score(makeResult(53, "udp")); got != priority.High {
		t.Fatalf("expected HIGH for udp, got %s", got)
	}
}

func TestScoreAllReturnsMap(t *testing.T) {
	s := priority.New([]priority.Rule{
		{Port: 443, Proto: "tcp", Level: priority.High},
	}, priority.Low)

	results := []scanner.Result{
		makeResult(443, "tcp"),
		makeResult(80, "tcp"),
	}
	m := s.ScoreAll(results)
	if m["443/tcp"] != priority.High {
		t.Fatalf("expected HIGH for 443/tcp")
	}
	if m["80/tcp"] != priority.Low {
		t.Fatalf("expected LOW for 80/tcp")
	}
}

func TestLevelString(t *testing.T) {
	cases := map[priority.Level]string{
		priority.Low:      "LOW",
		priority.Medium:   "MEDIUM",
		priority.High:     "HIGH",
		priority.Critical: "CRITICAL",
	}
	for lvl, want := range cases {
		if lvl.String() != want {
			t.Errorf("Level(%d).String() = %q, want %q", lvl, lvl.String(), want)
		}
	}
}
