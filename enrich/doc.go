// Package enrich provides a composable mechanism for attaching
// metadata to scanner results via pluggable Provider functions.
// Enriched results carry a Meta struct containing a timestamp and
// an arbitrary key/value map populated by the registered providers.
package enrich
