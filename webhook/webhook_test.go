package webhook_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/webhook"
)

func makeResult(port int, proto string) scanner.Result {
	return scanner.Result{Host: "localhost", Port: port, Protocol: proto, Open: true}
}

func TestSendPostsJSON(t *testing.T) {
	var got webhook.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := webhook.New(ts.URL)
	opened := []scanner.Result{makeResult(80, "tcp")}
	if err := s.Send("localhost", opened, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", got.Host)
	}
	if len(got.Opened) != 1 || got.Opened[0].Port != 80 {
		t.Errorf("expected opened port 80, got %+v", got.Opened)
	}
}

func TestSendSkipsWhenNoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := webhook.New(ts.URL)
	if err := s.Send("localhost", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call when there are no changes")
	}
}

func TestSendReturnsErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := webhook.New(ts.URL)
	opened := []scanner.Result{makeResult(443, "tcp")}
	if err := s.Send("localhost", opened, nil); err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestSendReturnsErrorOnBadURL(t *testing.T) {
	s := webhook.New("http://127.0.0.1:0/nowhere")
	opened := []scanner.Result{makeResult(22, "tcp")}
	if err := s.Send("localhost", opened, nil); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
