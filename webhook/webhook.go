// Package webhook sends port change notifications to an HTTP endpoint.
package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/scanner"
)

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp time.Time        `json:"timestamp"`
	Host      string           `json:"host"`
	Opened    []scanner.Result `json:"opened"`
	Closed    []scanner.Result `json:"closed"`
}

// Sender posts diff events to a remote HTTP endpoint.
type Sender struct {
	url    string
	client *http.Client
}

// New creates a Sender that posts to the given URL.
func New(url string) *Sender {
	return &Sender{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send posts the opened/closed results as JSON to the configured endpoint.
// It returns an error if the request fails or the server responds with a
// non-2xx status code.
func (s *Sender) Send(host string, opened, closed []scanner.Result) error {
	if len(opened) == 0 && len(closed) == 0 {
		return nil
	}

	p := Payload{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Opened:    opened,
		Closed:    closed,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}

	resp, err := s.client.Post(s.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	return nil
}
