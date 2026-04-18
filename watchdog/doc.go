// Package watchdog orchestrates a complete port-watch cycle.
//
// A single call to Watchdog.Run will:
//  1. Scan the configured host and ports.
//  2. Apply inclusion/exclusion filter rules.
//  3. Diff the results against the previously persisted state.
//  4. Dispatch alerts for any opened or closed ports.
//  5. Append the cycle to the rolling history log.
//  6. Persist the new state for the next cycle.
//
// Watchdog is designed to be called from a schedule.Runner tick so that
// every field of the struct remains immutable after construction.
package watchdog
