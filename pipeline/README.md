# pipeline

The `pipeline` package provides a composable, stage-based processing chain for
`scanner.Result` slices.

## Overview

A **Pipeline** is an ordered list of **Stages**. Each Stage is a pure function
that accepts a `[]scanner.Result` and returns a transformed `[]scanner.Result`.
Stages are applied sequentially; the output of one stage becomes the input of
the next.

## Usage

```go
import "github.com/user/portwatch/pipeline"

p := pipeline.New(
    func(rs []scanner.Result) []scanner.Result {
        // keep only TCP results
        var out []scanner.Result
        for _, r := range rs {
            if r.Proto == "tcp" {
                out = append(out, r)
            }
        }
        return out
    },
)

results, err := p.Run(ctx, scanResults)
```

## Builder

Use `Builder` for a fluent API that integrates common portwatch modules:

```go
p := pipeline.NewBuilder().
    OnlyOpen().
    WithFilter(myFilter).
    WithStage(customStage).
    Build()

results, err := p.Run(ctx, scanResults)
```

### Built-in builder methods

| Method | Description |
|---|---|
| `OnlyOpen()` | Removes closed-port results |
| `WithFilter(f)` | Applies a `filter.Filter` |
| `WithStage(s)` | Appends an arbitrary `Stage` |

## Context cancellation

`Run` checks `ctx` before each stage. If the context is cancelled, `Run`
returns immediately with the results produced so far and `ctx.Err()`.
