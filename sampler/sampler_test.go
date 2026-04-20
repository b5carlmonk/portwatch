package sampler_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/sampler"
	"github.com/user/portwatch/scanner"
)

func fakeScan(results []scanner.Result, err error) func(string, []int) ([]scanner.Result, error) {
	return func(_ string, _ []int) ([]scanner.Result, error) {
		return results, err
	}
}

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: "127.0.0.1", Port: p, Proto: "tcp", Open: true}
	}
	return out
}

func TestCaptureStoresSample(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(80, 443), nil))
	if err := s.Capture("baseline", "127.0.0.1", []int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sample, err := s.Get("baseline")
	if err != nil {
		t.Fatalf("expected sample: %v", err)
	}
	if len(sample.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(sample.Results))
	}
}

func TestCaptureEmptyNameErrors(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(80), nil))
	if err := s.Capture("", "127.0.0.1", []int{80}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestCaptureScanErrorPropagates(t *testing.T) {
	scanErr := errors.New("dial timeout")
	s := sampler.New(fakeScan(nil, scanErr))
	if err := s.Capture("test", "127.0.0.1", []int{80}); err == nil {
		t.Error("expected scan error to propagate")
	}
}

func TestGetMissingSampleErrors(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(80), nil))
	_, err := s.Get("missing")
	if err == nil {
		t.Error("expected error for missing sample")
	}
}

func TestDeleteRemovesSample(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(80), nil))
	_ = s.Capture("temp", "127.0.0.1", []int{80})
	s.Delete("temp")
	if _, err := s.Get("temp"); err == nil {
		t.Error("expected sample to be deleted")
	}
}

func TestNamesReturnsAllKeys(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(80), nil))
	_ = s.Capture("a", "127.0.0.1", []int{80})
	_ = s.Capture("b", "127.0.0.1", []int{443})
	names := s.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestSampleCapturedAtIsSet(t *testing.T) {
	s := sampler.New(fakeScan(makeResults(22), nil))
	_ = s.Capture("ts", "127.0.0.1", []int{22})
	sample, _ := s.Get("ts")
	if sample.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}
