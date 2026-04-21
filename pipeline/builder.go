package pipeline

import (
	"github.com/user/portwatch/filter"
	"github.com/user/portwatch/scanner"
)

// Builder provides a fluent API for assembling a Pipeline from
// common portwatch building blocks.
type Builder struct {
	stages []Stage
}

// NewBuilder returns an empty Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// WithStage appends an arbitrary Stage to the pipeline.
func (b *Builder) WithStage(s Stage) *Builder {
	b.stages = append(b.stages, s)
	return b
}

// WithFilter appends a stage that applies the given filter.Filter.
func (b *Builder) WithFilter(f *filter.Filter) *Builder {
	b.stages = append(b.stages, func(rs []scanner.Result) []scanner.Result {
		return f.Apply(rs)
	})
	return b
}

// OnlyOpen appends a stage that keeps only open ports.
func (b *Builder) OnlyOpen() *Builder {
	b.stages = append(b.stages, func(rs []scanner.Result) []scanner.Result {
		var out []scanner.Result
		for _, r := range rs {
			if r.Open {
				out = append(out, r)
			}
		}
		return out
	})
	return b
}

// Build returns the assembled Pipeline.
func (b *Builder) Build() *Pipeline {
	return New(b.stages...)
}
