// Package schedule provides an interval-based task runner used by portwatch
// to periodically scan ports, diff against saved state, and trigger alerts.
//
// The runner executes the task immediately on start and then on each interval
// tick until the provided context is cancelled.
package schedule
