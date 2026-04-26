package scorecard_test

import (
	"testing"

	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/scorecard"
)

func makeResults(host string, ports []int, open bool) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Host: host, Port: p, Proto: "tcp", Open: open})
	}
	return out
}

func TestScoreLowForSafePorts(t *testing.T) {
	s := scorecard.New()
	results := makeResults("localhost", []int{80, 443}, true)
	rep := s.Score(results)
	if rep.Level != scorecard.Low {
		t.Errorf("expected low, got %s (score %d)", rep.Level, rep.Score)
	}
}

func TestScoreHighForRiskyPorts(t *testing.T) {
	s := scorecard.New()
	results := makeResults("localhost", []int{23, 445, 3389}, true)
	rep := s.Score(results)
	if rep.Level != scorecard.Critical && rep.Level != scorecard.High {
		t.Errorf("expected high/critical for risky ports, got %s (score %d)", rep.Level, rep.Score)
	}
}

func TestBreakdownListsRiskyPorts(t *testing.T) {
	s := scorecard.New()
	results := makeResults("localhost", []int{23}, true)
	rep := s.Score(results)
	if len(rep.Breakdown) == 0 {
		t.Error("expected at least one breakdown entry for port 23")
	}
}

func TestClosedPortsNotCounted(t *testing.T) {
	s := scorecard.New()
	results := makeResults("localhost", []int{23, 445}, false)
	rep := s.Score(results)
	if rep.OpenPorts != 0 {
		t.Errorf("expected 0 open ports, got %d", rep.OpenPorts)
	}
	if rep.Score != 0 {
		t.Errorf("expected score 0 for all-closed ports, got %d", rep.Score)
	}
}

func TestHistoryGrowsWithEachScore(t *testing.T) {
	s := scorecard.New()
	results := makeResults("10.0.0.1", []int{80}, true)
	s.Score(results)
	s.Score(results)
	h := s.History("10.0.0.1")
	if len(h) != 2 {
		t.Errorf("expected 2 history entries, got %d", len(h))
	}
}

func TestEmptyResultsReturnsLow(t *testing.T) {
	s := scorecard.New()
	rep := s.Score(nil)
	if rep.Level != scorecard.Low {
		t.Errorf("expected low for empty results, got %s", rep.Level)
	}
}

func TestHostSetCorrectly(t *testing.T) {
	s := scorecard.New()
	results := makeResults("192.168.1.1", []int{22}, true)
	rep := s.Score(results)
	if rep.Host != "192.168.1.1" {
		t.Errorf("expected host 192.168.1.1, got %s", rep.Host)
	}
}
