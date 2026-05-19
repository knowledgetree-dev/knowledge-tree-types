package types

import "time"

type AssessmentType string

const (
	AssessmentTypeSecurity   AssessmentType = "security"
	AssessmentTypeCompliance AssessmentType = "compliance"
	AssessmentTypeCost       AssessmentType = "cost"
)

type AssessmentStatus string

const (
	AssessmentStatusQueued    AssessmentStatus = "queued"
	AssessmentStatusRunning   AssessmentStatus = "running"
	AssessmentStatusCompleted AssessmentStatus = "completed"
	AssessmentStatusFailed    AssessmentStatus = "failed"
)

type Finding struct {
	ID             string                 `json:"id"`
	ScopeID        string                 `json:"scope_id"`
	AssessmentID   string                 `json:"assessment_id"`
	AssessmentType AssessmentType         `json:"assessment_type"`
	Provider       string                 `json:"provider"`
	ResourceID     string                 `json:"resource_id"`
	ResourceType   string                 `json:"resource_type"`
	ResourceName   string                 `json:"resource_name,omitempty"`
	Region         string                 `json:"region,omitempty"`
	Severity       string                 `json:"severity"`
	Status         string                 `json:"status"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Remediation    string                 `json:"remediation,omitempty"`
	ControlID      string                 `json:"control_id,omitempty"`
	Framework      string                 `json:"framework,omitempty"`
	Standardized   bool                   `json:"standardized"`
	RawData        map[string]interface{} `json:"raw_data,omitempty"`
	DiscoveredAt   time.Time              `json:"discovered_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ResolvedAt     time.Time              `json:"resolved_at,omitempty"`
}

type CostRecord struct {
	ID            string    `json:"id"`
	ScopeID       string    `json:"scope_id"`
	AssessmentID  string    `json:"assessment_id"`
	Provider      string    `json:"provider"`
	ServiceName   string    `json:"service_name,omitempty"`
	ResourceID    string    `json:"resource_id"`
	ResourceType  string    `json:"resource_type"`
	ResourceName  string    `json:"resource_name,omitempty"`
	Region        string    `json:"region,omitempty"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	UsageAmount   float64   `json:"usage_amount"`
	UsageUnit     string    `json:"usage_unit,omitempty"`
	UnblendedCost float64   `json:"unblended_cost"`
	AmortizedCost float64   `json:"amortized_cost,omitempty"`
	Currency      string    `json:"currency,omitempty"`
	Granularity   string    `json:"granularity,omitempty"`
	DiscoveredAt  time.Time `json:"discovered_at,omitempty"`
}

type AssessmentRun struct {
	ID             string           `json:"id"`
	ScopeID        string           `json:"scope_id"`
	Type           AssessmentType   `json:"type"`
	Status         AssessmentStatus `json:"status"`
	Findings       int              `json:"findings"`
	FindingDetails []Finding        `json:"-"`
	CostRecords    []CostRecord     `json:"cost_records,omitempty"`
	Errors         []string         `json:"errors,omitempty"`
	StartedAt      time.Time        `json:"started_at,omitempty"`
	CompletedAt    time.Time        `json:"completed_at,omitempty"`
}
