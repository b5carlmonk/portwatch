package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/pipeline"
	"github.com/user/portwatch/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Host: "127.0.0.1", Port: p, Proto: "tcp", Open: true})
	}
	return out
}

func TestEmptyPipelineReturnsInput(t *testing.T) {
	p := pipeline.New()
	input := makeResults(80, 443)
	out, err := p.Run(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestStageTransformsResults(t *testing.T) {
	filterHTTPS := pipeline.Stage(func(rs []scanner.Result) []scanner.Result {
		var out []scanner.Result
		for _, r := range rs {
			if r.Port != 443 {
				out = append(out, r)
			}
		}
		return out
	})
	p := pipeline.New(filterHTTPS)
	out, err := p.Run(context.Background(), makeResults(80, 443, 8080))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results after filter, got %d", len(out))
	}
}

func TestStagesAppliedInOrder(t *testing.T) {
	var order []int
	stage1 := pipeline.Stage(func(rs []scanner.Result) []scanner.Result { order = append(order, 1); return rs })
	stage2 := pipeline.Stage(func(rs []scanner.Result) []scanner.Result { order = append(order, 2); return rs })
	p := pipeline.New(stage1, stage2)
	p.Run(context.Background(), makeResults(80)) //nolint
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("expected stages called in order [1,2], got %v", order)
	}
}

func TestCancelledContextStopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	stage := pipeline.Stage(func(rs []scanner.Result) []scanner.Result {
		called = true
		return rs
	})
	p := pipeline.New(stage)
	_, err := p.Run(ctx, makeResults(80))
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
	if called {
		t.Fatal("stage should not have been called after context cancel")
	}
}

func TestLenReturnsStageCount(t *testing.T) {
	p := pipeline.New(
		pipeline.Stage(func(rs []scanner.Result) []scanner.Result { return rs }),
		pipeline.Stage(func(rs []scanner.Result) []scanner.Result { return rs }),
	)
	if p.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", p.Len())
	}
}
