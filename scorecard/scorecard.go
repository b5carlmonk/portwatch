// Package scorecard computes a risk score for a scanned host based on
// open ports, known vulnerabilities, and recent change activity.
package scorecard

import (
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/scanner"
)

// Risk levels returned by Score.
const (
	Low      = "low"
	Medium   = "medium"
	High     = "high"
	Critical = "critical"
)

// riskPorts maps well-known sensitive ports to a base penalty.
var riskPorts = map[int]int{
	21:   30, // FTP
	23:   40, // Telnet
	445:  35, // SMB
	3389: 35, // RDP
	5900: 25, // VNC
	6379: 20, // Redis
	27017: 20, // MongoDB
}

// Report holds the computed score and breakdown for a host.
type Report struct {
	Host       string
	Score      int
	Level      string
	Breakdown  []string
	OpenPorts  int
}

// Scorer computes risk scores.
type Scorer struct {
	mu       sync.Mutex
	history  map[string][]Report
}

// New returns a new Scorer.
func New() *Scorer {
	return &Scorer{history: make(map[string][]Report)}
}

// Score computes a risk report for the given scan results.
func (s *Scorer) Score(results []scanner.Result) Report {
	if len(results) == 0 {
		return Report{Level: Low}
	}

	host := results[0].Host
	total := 0
	var breakdown []string
	open := 0

	for _, r := range results {
		if !r.Open {
			continue
		}
		open++
		total += 2 // base cost per open port
		if penalty, ok := riskPorts[r.Port]; ok {
			total += penalty
			breakdown = append(breakdown, fmt.Sprintf("port %d (%s) +%d", r.Port, r.Proto, penalty))
		}
	}

	sort.Strings(breakdown)

	rep := Report{
		Host:      host,
		Score:     total,
		Level:     level(total),
		Breakdown: breakdown,
		OpenPorts: open,
	}

	s.mu.Lock()
	s.history[host] = append(s.history[host], rep)
	s.mu.Unlock()

	return rep
}

// History returns all past reports for a host.
func (s *Scorer) History(host string) []Report {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := make([]Report, len(s.history[host]))
	copy(copied, s.history[host])
	return copied
}

func level(score int) string {
	switch {
	case score >= 80:
		return Critical
	case score >= 50:
		return High
	case score >= 20:
		return Medium
	default:
		return Low
	}
}
