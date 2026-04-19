// Package groupby groups scan results by a chosen dimension (port, protocol, host).
package groupby

import (
	"fmt"
	"strconv"

	"github.com/user/portwatch/scanner"
)

// Dimension controls how results are grouped.
type Dimension string

const (
	ByPort     Dimension = "port"
	ByProtocol Dimension = "protocol"
	ByHost     Dimension = "host"
)

// Grouper groups scan results by a dimension.
type Grouper struct {
	dim Dimension
}

// New returns a Grouper for the given dimension.
func New(dim Dimension) *Grouper {
	return &Grouper{dim: dim}
}

// Group partitions results into named buckets.
func (g *Grouper) Group(results []scanner.Result) map[string][]scanner.Result {
	out := make(map[string][]scanner.Result)
	for _, r := range results {
		k := g.keyFor(r)
		out[k] = append(out[k], r)
	}
	return out
}

func (g *Grouper) keyFor(r scanner.Result) string {
	switch g.dim {
	case ByPort:
		return strconv.Itoa(r.Port)
	case ByProtocol:
		return r.Protocol
	case ByHost:
		return r.Host
	default:
		return fmt.Sprintf("%s:%d", r.Protocol, r.Port)
	}
}
