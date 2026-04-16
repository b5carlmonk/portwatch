package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	Service  string
}

// ScanResult holds the results of a full host scan.
type ScanResult struct {
	Host      string
	Timestamp time.Time
	Ports     []PortState
}

// Scanner performs port scanning on a host.
type Scanner struct {
	Timeout time.Duration
}

// New creates a new Scanner with the given timeout.
func New(timeout time.Duration) *Scanner {
	return &Scanner{Timeout: timeout}
}

// Scan scans the given ports on the host and returns a ScanResult.
func (s *Scanner) Scan(host string, ports []int) (*ScanResult, error) {
	result := &ScanResult{
		Host:      host,
		Timestamp: time.Now(),
	}

	for _, port := range ports {
		state := s.checkPort(host, port)
		result.Ports = append(result.Ports, state)
	}

	return result, nil
}

func (s *Scanner) checkPort(host string, port int) PortState {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	if err != nil {
		return PortState{Port: port, Protocol: "tcp", Open: false}
	}
	conn.Close()
	return PortState{
		Port:     port,
		Protocol: "tcp",
		Open:     true,
		Service:  knownServices[port],
	}
}

// knownServices maps well-known ports to service names.
var knownServices = map[int]string{
	21:   "ftp",
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
}
