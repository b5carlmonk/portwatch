// Package pipeline chains scan results through a series of processing stages.
package pipeline

import (
	"context"

	"github.com/user/portwatch/scanner"
)

// Stage is a function that transforms a slice of ScanResults.
type Stage func([]scanner.Result) []scanner.Result

// Pipeline executes a sequence of stages against scan results.
type Pipeline struct {
	stages []Stage
}

// New returns a new Pipeline with the given stages applied in order.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run passes results through each registered stage in order.
// If ctx is cancelled before completion, Run returns the results
// produced up to that point along with ctx.Err().
func (p *Pipeline) Run(ctx context.Context, results []scanner.Result) ([]scanner.Result, error) {
	current := results
	for _, stage := range p.stages {
		select {
		case <-ctx.Done():
			return current, ctx.Err()
		default:
		}
		current = stage(current)
	}
	return current, nil
}

// Len returns the number of stages registered in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
