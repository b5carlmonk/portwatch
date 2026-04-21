package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/filter"
	"github.com/user/portwatch/pipeline"
	"github.com/user/portwatch/scanner"
)

func TestOnlyOpenFiltersClosedPorts(t *testing.T) {
	results := []scanner.Result{
		{Host: "127.0.0.1", Port: 80, Proto: "tcp", Open: true},
		{Host: "127.0.0.1", Port: 81, Proto: "tcp", Open: false},
		{Host: "127.0.0.1", Port: 443, Proto: "tcp", Open: true},
	}
	p := pipeline.NewBuilder().OnlyOpen().Build()
	out, err := p.Run(context.Background(), results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 open results, got %d", len(out))
	}
	for _, r := range out {
		if !r.Open {
			t.Errorf("expected only open ports, got closed port %d", r.Port)
		}
	}
}

func TestWithFilterAppliesRules(t *testing.T) {
	f := filter.New(filter.Rule{Port: 22, Action: filter.Exclude})
	p := pipeline.NewBuilder().WithFilter(f).Build()

	input := makeResults(22, 80, 443)
	out, err := p.Run(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if r.Port == 22 {
			t.Error("port 22 should have been excluded")
		}
	}
}

func TestBuilderChainsMultipleStages(t *testing.T) {
	results := []scanner.Result{
		{Host: "127.0.0.1", Port: 22, Proto: "tcp", Open: false},
		{Host: "127.0.0.1", Port: 80, Proto: "tcp", Open: true},
		{Host: "127.0.0.1", Port: 443, Proto: "tcp", Open: true},
	}
	f := filter.New(filter.Rule{Port: 443, Action: filter.Exclude})
	p := pipeline.NewBuilder().OnlyOpen().WithFilter(f).Build()

	out, err := p.Run(context.Background(), results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Port != 80 {
		t.Fatalf("expected only port 80, got %v", out)
	}
}

func TestBuilderBuildReturnsCorrectLen(t *testing.T) {
	p := pipeline.NewBuilder().
		OnlyOpen().
		WithStage(pipeline.Stage(func(rs []scanner.Result) []scanner.Result { return rs })).
		Build()
	if p.Len() != 2 {
		t.Fatalf("expected pipeline length 2, got %d", p.Len())
	}
}
