package enrich_test

import (
	"testing"

	"github.com/user/portwatch/enrich"
	"github.com/user/portwatch/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "127.0.0.1", Port: 80, Proto: "tcp", Open: true},
		{Host: "127.0.0.1", Port: 443, Proto: "tcp", Open: true},
		{Host: "127.0.0.1", Port: 53, Proto: "udp", Open: true},
	}
}

func TestEnrichReturnsAllResults(t *testing.T) {
	e := enrich.New()
	out := e.Enrich(makeResults())
	if len(out) != 3 {
		t.Fatalf("expected 3 enriched results, got %d", len(out))
	}
}

func TestEnrichSetsScannedAt(t *testing.T) {
	e := enrich.New()
	out := e.Enrich(makeResults())
	for _, r := range out {
		if r.Meta.ScannedAt.IsZero() {
			t.Error("expected ScannedAt to be set")
		}
	}
}

func TestEnrichAppliesPortLabel(t *testing.T) {
	e := enrich.New(enrich.PortLabel)
	out := e.Enrich(makeResults())
	if got := out[0].Meta.Extra["port_label"]; got != "port-80" {
		t.Errorf("expected port-80, got %s", got)
	}
}

func TestEnrichAppliesProtoLabel(t *testing.T) {
	e := enrich.New(enrich.ProtoLabel)
	out := e.Enrich(makeResults())
	if got := out[2].Meta.Extra["proto_label"]; got != "udp" {
		t.Errorf("expected udp, got %s", got)
	}
}

func TestEnrichMultipleProviders(t *testing.T) {
	e := enrich.New(enrich.PortLabel, enrich.ProtoLabel)
	out := e.Enrich(makeResults())
	if _, ok := out[0].Meta.Extra["port_label"]; !ok {
		t.Error("expected port_label key")
	}
	if _, ok := out[0].Meta.Extra["proto_label"]; !ok {
		t.Error("expected proto_label key")
	}
}

func TestEnrichCustomProvider(t *testing.T) {
	custom := func(r scanner.Result) (string, string) {
		if r.Port == 443 {
			return "secure", "true"
		}
		return "", ""
	}
	e := enrich.New(custom)
	out := e.Enrich(makeResults())
	if got := out[1].Meta.Extra["secure"]; got != "true" {
		t.Errorf("expected true, got %s", got)
	}
	if _, ok := out[0].Meta.Extra["secure"]; ok {
		t.Error("expected no secure key for port 80")
	}
}

func TestEnrichEmptyInput(t *testing.T) {
	e := enrich.New(enrich.PortLabel)
	out := e.Enrich([]scanner.Result{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d", len(out))
	}
}
