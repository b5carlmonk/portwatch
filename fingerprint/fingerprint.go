// Package fingerprint identifies services running on open ports by
// performing lightweight banner grabs and matching against known signatures.
package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/user/portwatch/scanner"
)

// Service holds the identified service name and banner for a port.
type Service struct {
	Port    int
	Proto   string
	Banner  string
	Service string
}

// Fingerprinter grabs banners from open ports and maps them to service names.
type Fingerprinter struct {
	timeout    time.Duration
	signatures map[string]string
}

// New returns a Fingerprinter with the given dial timeout.
func New(timeout time.Duration) *Fingerprinter {
	return &Fingerprinter{
		timeout:    timeout,
		signatures: defaultSignatures(),
	}
}

// Identify attempts a banner grab on each open result and returns a Service
// slice with any recognised service names attached.
func (f *Fingerprinter) Identify(results []scanner.Result) []Service {
	var services []Service
	for _, r := range results {
		if !r.Open {
			continue
		}
		banner := f.grab(r.Host, r.Port, r.Proto)
		name := f.match(banner)
		services = append(services, Service{
			Port:    r.Port,
			Proto:   r.Proto,
			Banner:  banner,
			Service: name,
		})
	}
	return services
}

func (f *Fingerprinter) grab(host string, port int, proto string) string {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(proto, addr, f.timeout)
	if err != nil {
		return ""
	}
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(f.timeout))
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	return strings.TrimSpace(string(buf[:n]))
}

func (f *Fingerprinter) match(banner string) string {
	upper := strings.ToUpper(banner)
	for sig, name := range f.signatures {
		if strings.Contains(upper, sig) {
			return name
		}
	}
	if banner != "" {
		return "unknown"
	}
	return ""
}

func defaultSignatures() map[string]string {
	return map[string]string{
		"SSH":   "ssh",
		"HTTP":  "http",
		"FTP":   "ftp",
		"SMTP":  "smtp",
		"IMAP":  "imap",
		"POP3":  "pop3",
		"MYSQL": "mysql",
		"REDIS": "redis",
	}
}
