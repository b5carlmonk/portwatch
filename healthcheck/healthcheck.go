package healthcheck

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single health probe.
type Result struct {
	Host    string
	Port    int
	Proto   string
	Alive   bool
	Latency time.Duration
	Err     error
}

// Checker probes whether a specific port is reachable.
type Checker struct {
	timeout time.Duration
}

// New returns a Checker with the given dial timeout.
func New(timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Checker{timeout: timeout}
}

// Probe dials host:port over proto and returns a Result.
func (c *Checker) Probe(host string, port int, proto string) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout(proto, addr, c.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Host: host, Port: port, Proto: proto, Alive: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return Result{Host: host, Port: port, Proto: proto, Alive: true, Latency: latency}
}

// ProbeAll probes a list of port/proto pairs and returns all results.
func (c *Checker) ProbeAll(host string, targets []Target) []Result {
	results := make([]Result, 0, len(targets))
	for _, t := range targets {
		results = append(results, c.Probe(host, t.Port, t.Proto))
	}
	return results
}

// Target describes a port+protocol pair to probe.
type Target struct {
	Port  int
	Proto string
}
