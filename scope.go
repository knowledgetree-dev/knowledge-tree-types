package types

import "time"

// ScopeStatus represents the current state of a discovery scope.
type ScopeStatus string

const (
	// ScopeStatusActive indicates the scope is enabled and can be used
	// for discovery runs.
	ScopeStatusActive ScopeStatus = "active"

	// ScopeStatusPaused indicates the scope is temporarily disabled and
	// will not be scheduled for automatic discovery.
	ScopeStatusPaused ScopeStatus = "paused"

	// ScopeStatusError indicates the scope has encountered a persistent
	// error (e.g., invalid credentials) and requires user intervention.
	ScopeStatusError ScopeStatus = "error"
)

// DiscoveryScope defines the boundary and configuration for a set of
// discovery operations. A scope typically maps to a cloud account, region,
// or Kubernetes cluster and specifies which plugin to use, what credentials
// to apply, and how often to run discovery.
type DiscoveryScope struct {
	// ID is a unique identifier for the scope, typically a UUID.
	ID string `json:"id"`

	// Name is a human-readable name for the scope
	// (e.g., "Production AWS us-east-1").
	Name string `json:"name"`

	// PluginName identifies which discovery plugin should be used for this
	// scope (e.g., "aws", "kubernetes", "azure").
	PluginName string `json:"plugin_name"`

	// Config contains plugin-specific configuration key-value pairs. The
	// schema of this config is defined by the plugin's InfoResponse.
	Config map[string]string `json:"config"`

	// CredentialSource specifies where the plugin should obtain credentials
	// for accessing the infrastructure.
	CredentialSource CredentialSource `json:"credential_source"`

	// Schedule is a cron expression that defines how often the discovery
	// should run automatically (e.g., "0 */6 * * *" for every 6 hours).
	// An empty string means the scope is only run manually.
	Schedule string `json:"schedule"`

	// Status is the current operational status of the scope.
	Status ScopeStatus `json:"status"`

	// LastRun is the timestamp of the most recently completed discovery run
	// for this scope. Zero value means the scope has never been run.
	LastRun time.Time `json:"last_run"`

	// ResourceCount is the total number of resources currently discovered
	// within this scope.
	ResourceCount int `json:"resource_count"`
}

// CredentialSource describes where a plugin should obtain credentials
// for accessing infrastructure. This mirrors the proto CredentialSource
// message as a native Go type for use within the core system.
type CredentialSource struct {
	// Type identifies the credential backend
	// (e.g., "vault", "aws_secrets_manager", "env", "file").
	Type string `json:"type"`

	// Path is the location within the credential backend to fetch credentials
	// from (e.g., "secret/data/aws/production").
	Path string `json:"path"`

	// Params contains additional parameters for the credential source
	// (e.g., {"role": "knowledge-tree"} for Vault authentication).
	Params map[string]string `json:"params,omitempty"`
}
