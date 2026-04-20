// Package classify labels scan results by known service names based on port and protocol.
package classify

import (
	"fmt"

	"github.com/user/portwatch/scanner"
)

// wellKnown maps "port/proto" to a human-readable service name.
var wellKnown = map[string]string{
	"22/tcp":   "SSH",
	"80/tcp":   "HTTP",
	"443/tcp":  "HTTPS",
	"3306/tcp": "MySQL",
	"5432/tcp": "PostgreSQL",
	"6379/tcp": "Redis",
	"27017/tcp": "MongoDB",
	"8080/tcp": "HTTP-Alt",
	"8443/tcp": "HTTPS-Alt",
	"53/udp":   "DNS",
	"53/tcp":   "DNS",
	"25/tcp":   "SMTP",
	"587/tcp":  "SMTP-Submission",
	"110/tcp":  "POP3",
	"143/tcp":  "IMAP",
	"21/tcp":   "FTP",
	"23/tcp":   "Telnet",
}

// Classifier attaches service labels to scan results.
type Classifier struct {
	custom map[string]string
}

// New returns a Classifier with optional extra mappings merged over the defaults.
func New(extra map[string]string) *Classifier {
	custom := make(map[string]string, len(extra))
	for k, v := range extra {
		custom[k] = v
	}
	return &Classifier{custom: custom}
}

// Label returns the service name for the given port and protocol.
// It checks custom mappings first, then the built-in table.
// If no match is found it returns "unknown".
func (c *Classifier) Label(port int, proto string) string {
	k := key(port, proto)
	if v, ok := c.custom[k]; ok {
		return v
	}
	if v, ok := wellKnown[k]; ok {
		return v
	}
	return "unknown"
}

// Enrich annotates each result's Service field with the resolved label.
// Results are modified in place and the same slice is returned.
func (c *Classifier) Enrich(results []scanner.Result) []scanner.Result {
	for i := range results {
		results[i].Service = c.Label(results[i].Port, results[i].Proto)
	}
	return results
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
