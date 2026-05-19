package types

import "time"

// DiscoveryStatus represents the current state of a discovery run.
type DiscoveryStatus string

const (
	// DiscoveryStatusPending indicates the run has been created but has not
	// started executing yet.
	DiscoveryStatusPending DiscoveryStatus = "PENDING"

	// DiscoveryStatusRunning indicates the run is currently in progress
	// and streaming discovery events.
	DiscoveryStatusRunning DiscoveryStatus = "RUNNING"

	// DiscoveryStatusCompleted indicates the run finished successfully.
	DiscoveryStatusCompleted DiscoveryStatus = "COMPLETED"

	// DiscoveryStatusFailed indicates the run terminated due to an error.
	DiscoveryStatusFailed DiscoveryStatus = "FAILED"
)

// DiscoveryRun tracks the lifecycle of a single discovery execution within
// a scope. It records timing, counts, and any errors encountered.
type DiscoveryRun struct {
	// ID is a unique identifier for this run, typically a UUID.
	ID string `json:"id"`

	// ScopeID references the DiscoveryScope that this run belongs to.
	ScopeID string `json:"scope_id"`

	// Status is the current state of the run.
	Status DiscoveryStatus `json:"status"`

	// StartedAt is the time when the run began executing. Zero value means
	// the run has not started yet.
	StartedAt time.Time `json:"started_at"`

	// CompletedAt is the time when the run finished (either successfully or
	// with a failure). Zero value means the run is still in progress or
	// has not started.
	CompletedAt time.Time `json:"completed_at"`

	// ResourcesFound is the total number of resources discovered during
	// this run.
	ResourcesFound int `json:"resources_found"`

	// RelationshipsFound is the total number of relationships discovered
	// during this run.
	RelationshipsFound int `json:"relationships_found"`

	// Errors collects all error messages encountered during the run. Non-fatal
	// errors are recorded here while the run continues; a fatal error will
	// set the Status to DiscoveryStatusFailed.
	Errors []string `json:"errors"`
}

// Duration returns the elapsed time of the run. If the run is still in
// progress, it returns the duration from start to now. If the run has not
// started, it returns zero.
func (r *DiscoveryRun) Duration() time.Duration {
	if r.StartedAt.IsZero() {
		return 0
	}
	if r.CompletedAt.IsZero() {
		return time.Since(r.StartedAt)
	}
	return r.CompletedAt.Sub(r.StartedAt)
}
