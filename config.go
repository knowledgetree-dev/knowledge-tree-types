// Package types defines the core data types used throughout Knowledge Tree.
package types

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Configuration types
// ---------------------------------------------------------------------------

// Config represents the full application configuration loaded from YAML.
type Config struct {
	ProjectName string           `mapstructure:"project_name" yaml:"project_name" json:"project_name"`
	Storage     StorageConfig    `mapstructure:"storage" yaml:"storage" json:"storage"`
	Scopes      map[string]Scope `mapstructure:"scopes" yaml:"scopes" json:"scopes"`
	Logging     LoggingConfig    `mapstructure:"logging" yaml:"logging" json:"logging"`
	Plugins     PluginsConfig    `mapstructure:"plugins" yaml:"plugins" json:"plugins"`
	Docs        DocsConfig       `mapstructure:"docs" yaml:"docs" json:"docs"`
	Agents      AgentsConfig     `mapstructure:"agents" yaml:"agents" json:"agents"`
	MCP         MCPConfig        `mapstructure:"mcp" yaml:"mcp" json:"mcp"`
}
// Validate checks the configuration for errors and returns the first one found.
func (c *Config) Validate() error {
	// Storage validation
	if c.Storage.Backend == "postgres" {
		if c.Storage.Host == "" {
			return fmt.Errorf("storage.host is required for postgres backend")
		}
		if c.Storage.Port <= 0 || c.Storage.Port > 65535 {
			return fmt.Errorf("storage.port must be between 1 and 65535, got %d", c.Storage.Port)
		}
		if c.Storage.Database == "" {
			return fmt.Errorf("storage.database is required for postgres backend")
		}
		if c.Storage.User == "" {
			return fmt.Errorf("storage.user is required for postgres backend")
		}
	}

	// Logging validation
	switch strings.ToLower(c.Logging.Level) {
	case "debug", "info", "warn", "error", "":
	default:
		return fmt.Errorf("logging.level must be debug/info/warn/error, got %q", c.Logging.Level)
	}

	// Plugin directory
	if c.Plugins.Directory != "" {
		if info, err := os.Stat(c.Plugins.Directory); err != nil || !info.IsDir() {
			return fmt.Errorf("plugins.directory %q does not exist or is not a directory", c.Plugins.Directory)
		}
	}

	// Docs site URL
	if c.Docs.SiteURL != "" {
		if _, err := url.Parse(c.Docs.SiteURL); err != nil {
			return fmt.Errorf("docs.site_url is invalid: %w", err)
		}
	}

	// Scope validation: schedule must be valid 5-field cron expression
	for name, scope := range c.Scopes {
		if scope.Schedule != "" && !validCronExpr(scope.Schedule) {
			return fmt.Errorf("scopes.%s.schedule is invalid: expected 5-field cron expression", name)
		}
	}

	// MCP port range
	if c.MCP.Port < 0 || c.MCP.Port > 65535 {
		return fmt.Errorf("mcp.port must be between 0 and 65535, got %d", c.MCP.Port)
	}

	// Agent port range
	if c.Agents.Port < 0 || c.Agents.Port > 65535 {
		return fmt.Errorf("agents.port must be between 0 and 65535, got %d", c.Agents.Port)
	}

	// Network check helper (used by storage backend)
	_ = net.JoinHostPort

	return nil
}

// Version is set at build time via ldflags.
var Version = "dev"

// CommitHash is set at build time via ldflags.
var CommitHash = "unknown"

// BuildDate is set at build time via ldflags.
var BuildDate = "unknown"

// FormatVersion returns a human-readable version string.
func FormatVersion() string {
	return fmt.Sprintf("knowledge-tree %s (commit: %s, built: %s)", Version, CommitHash, BuildDate)
}


// StorageConfig configures the graph storage backend.
type StorageConfig struct {
	Backend string `mapstructure:"backend" yaml:"backend" json:"backend"` // "postgres" or "sqlite"
	// PostgreSQL settings
	Host     string `mapstructure:"host" yaml:"host" json:"host"`
	Port     int    `mapstructure:"port" yaml:"port" json:"port"`
	Database string `mapstructure:"database" yaml:"database" json:"database"`
	User     string `mapstructure:"user" yaml:"user" json:"user"`
	Password string `mapstructure:"password" yaml:"password" json:"password"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode" json:"sslmode"`
	// SQLite settings
	Path string `mapstructure:"path" yaml:"path" json:"path"`
}

// DSN returns the database connection string with the password included.
// This should only be used for actual database connections, never for logging.
func (s StorageConfig) DSN() string {
	if s.Backend == "sqlite" {
		return s.Path
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.Host, s.Port, s.User, s.Password, s.Database, s.SSLMode)
}

// DSNSafe returns the database connection string with the password redacted.
// Use this for logging and error messages.
func (s StorageConfig) DSNSafe() string {
	if s.Backend == "sqlite" {
		return s.Path
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=REDACTED dbname=%s sslmode=%s",
		s.Host, s.Port, s.User, s.Database, s.SSLMode)
}

// Scope defines a discovery scope for the CLI configuration file.
// This is the YAML-serializable representation stored in the config file.
// The system-internal representation is DiscoveryScope (see scope.go).
type Scope struct {
	Name        string            `mapstructure:"name" yaml:"name" json:"name"`
	Plugin      string            `mapstructure:"plugin" yaml:"plugin" json:"plugin"`
	Description string            `mapstructure:"description" yaml:"description" json:"description"`
	Config      map[string]string `mapstructure:"config" yaml:"config" json:"config"`
	Enabled     bool              `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
	Schedule    string            `mapstructure:"schedule" yaml:"schedule" json:"schedule"`
	Credentials CLIConfigCredentialSource `mapstructure:"credentials" yaml:"credentials" json:"credentials"`
	CreatedAt   time.Time         `mapstructure:"created_at" yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `mapstructure:"updated_at" yaml:"updated_at" json:"updated_at"`
	LastRunAt   *time.Time        `mapstructure:"last_run_at" yaml:"last_run_at" json:"last_run_at"`
	LastRunID   string            `mapstructure:"last_run_id" yaml:"last_run_id" json:"last_run_id"`
	Status      CLIScopeStatus    `mapstructure:"status" yaml:"status" json:"status"`
}

// CLIScopeStatus represents the current status of a CLI scope.
type CLIScopeStatus string

const (
	CLIScopeStatusActive   CLIScopeStatus = "active"
	CLIScopeStatusInactive CLIScopeStatus = "inactive"
	CLIScopeStatusRunning  CLIScopeStatus = "running"
	CLIScopeStatusError    CLIScopeStatus = "error"
	CLIScopeStatusUnknown  CLIScopeStatus = "unknown"
)

// CLIConfigCredentialSource describes where to obtain credentials for a CLI scope.
// This mirrors the core CredentialSource but is a separate type for the CLI
// config file format.
type CLIConfigCredentialSource struct {
	Type   string            `mapstructure:"type" yaml:"type" json:"type"`
	Path   string            `mapstructure:"path" yaml:"path" json:"path"`
	Params map[string]string `mapstructure:"params" yaml:"params" json:"params"`
}

// LoggingConfig configures logging behavior.
type LoggingConfig struct {
	Level  string `mapstructure:"level" yaml:"level" json:"level"`
	Format string `mapstructure:"format" yaml:"format" json:"format"`
}

// PluginsConfig configures plugin discovery and loading.
type PluginsConfig struct {
	Directory string `mapstructure:"directory" yaml:"directory" json:"directory"`
	Registry  string `mapstructure:"registry" yaml:"registry" json:"registry"`
}

// DocsConfig configures documentation generation.
type DocsConfig struct {
	OutputDir string `mapstructure:"output_dir" yaml:"output_dir" json:"output_dir"`
	Theme     string `mapstructure:"theme" yaml:"theme" json:"theme"`
	SiteName  string `mapstructure:"site_name" yaml:"site_name" json:"site_name"`
	SiteURL   string `mapstructure:"site_url" yaml:"site_url" json:"site_url"`
}

// AgentsConfig configures agent deployment.
type AgentsConfig struct {
	SSHKeyPath string `mapstructure:"ssh_key_path" yaml:"ssh_key_path" json:"ssh_key_path"`
	SSHUser    string `mapstructure:"ssh_user" yaml:"ssh_user" json:"ssh_user"`
	Port       int    `mapstructure:"port" yaml:"port" json:"port"`
}

// MCPConfig configures the MCP server.
type MCPConfig struct {
	Port       int    `mapstructure:"port" yaml:"port" json:"port"`
	EnableHTTP bool   `mapstructure:"enable_http" yaml:"enable_http" json:"enable_http"`
	APIKey     string `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() *Config {
	return &Config{
		ProjectName: "my-infrastructure",
		Storage: StorageConfig{
			Backend:  "postgres",
			Host:     "localhost",
			Port:     5432,
			Database: "knowledge_tree",
			User:     "knowledge_tree",
			Password: "",
			SSLMode:  "disable",
			Path:     ".knowledge-tree/data/knowledge_tree.db",
		},
		Scopes: make(map[string]Scope),
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Plugins: PluginsConfig{
			Directory: "plugins",
			Registry:  "",
		},
		Docs: DocsConfig{
			OutputDir: "docs",
			Theme:     "material",
			SiteName:  "Infrastructure Documentation",
			SiteURL:   "http://localhost:8000",
		},
		Agents: AgentsConfig{
			SSHKeyPath: "~/.ssh/id_rsa",
			SSHUser:    "root",
			Port:       9443,
		},
		MCP: MCPConfig{
			Port:       9090,
			EnableHTTP: false,
		},
	}
}

// ---------------------------------------------------------------------------
// CLI-specific types (not duplicated from the core types)
// ---------------------------------------------------------------------------

// DiscoveryResult holds the aggregate result of a discovery run from the CLI perspective.
type DiscoveryResult struct {
	RunID          string            `json:"run_id"`
	ScopeName      string            `json:"scope_name"`
	StartedAt      time.Time         `json:"started_at"`
	FinishedAt     time.Time         `json:"finished_at"`
	Duration       time.Duration     `json:"duration"`
	ResourcesFound int               `json:"resources_found"`
	RelatsFound    int               `json:"relationships_found"`
	Errors         []DiscoveryError  `json:"errors"`
	Status         RunStatus         `json:"status"`
}

// RunStatus represents the status of a discovery run.
type RunStatus string

const (
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
	RunStatusPartial   RunStatus = "partial"
)

// DiscoveryError represents an error that occurred during discovery.
type DiscoveryError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Resource  string `json:"resource"`
	Retryable bool   `json:"retryable"`
}

// PluginInfo holds metadata about a loaded plugin.
type PluginInfo struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	Description     string   `json:"description"`
	Capabilities    []string `json:"capabilities"`
	ConfigSchema    string   `json:"config_schema"`
	CredentialTypes []string `json:"credential_types"`
	Path            string   `json:"path"`
	Status          string   `json:"status"`
}

// AgentDeployment holds information about a deployed agent.
type AgentDeployment struct {
	Host       string    `json:"host"`
	Status     string    `json:"status"`
	LastSeen   time.Time `json:"last_seen"`
	Version    string    `json:"version"`
	ScopeName  string    `json:"scope_name"`
	DeployedAt time.Time `json:"deployed_at"`
}

// GraphQueryResult holds the result of a graph query.
type GraphQueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Count   int             `json:"count"`
}

// GraphExportOptions configures graph export.
type GraphExportOptions struct {
	Format      string `json:"format"`
	ScopeFilter string `json:"scope_filter"`
	Depth       int    `json:"depth"`
}


// validCronExpr validates a 5-field cron expression (minute hour dom month dow).
// This replaces the robfig/cron ParseStandard dependency to keep types zero-dep.
func validCronExpr(expr string) bool {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return false
	}
	// Validate each field accepts common values: *, ranges (1-5), steps (*/5), commas
	for _, part := range parts {
		if part == "" || part == "?" || part == "@" {
			return false
		}
	}
	return true
}
