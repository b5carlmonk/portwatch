// Package export provides utilities for serialising port scan results
// to structured output formats such as JSON and CSV.
//
// Use export.New to create an Exporter bound to an io.Writer and a
// desired Format, then call Write with a slice of scanner.Result values.
package export
