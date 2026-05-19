package types

import (
	"strings"
	"testing"
)

func TestStorageConfig_DSN(t *testing.T) {
	tests := []struct {
		name   string
		config StorageConfig
		want   string
	}{
		{
			name: "PostgreSQL DSN",
			config: StorageConfig{
				Backend:  "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "knowledge_tree",
				User:     "postgres",
				Password: "secret",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=postgres password=secret dbname=knowledge_tree sslmode=disable",
		},
		{
			name: "SQLite DSN",
			config: StorageConfig{
				Backend: "sqlite",
				Path:    "/path/to/db.sqlite",
			},
			want: "/path/to/db.sqlite",
		},
		{
			name: "PostgreSQL with SSL",
			config: StorageConfig{
				Backend:  "postgres",
				Host:     "db.example.com",
				Port:     5432,
				Database: "mydb",
				User:     "user",
				Password: "pass",
				SSLMode:  "require",
			},
			want: "host=db.example.com port=5432 user=user password=pass dbname=mydb sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.DSN()
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ProjectName != "my-infrastructure" {
		t.Errorf("ProjectName = %v, want my-infrastructure", cfg.ProjectName)
	}
	if cfg.Storage.Backend != "postgres" {
		t.Errorf("Storage.Backend = %v, want postgres", cfg.Storage.Backend)
	}
	if cfg.Storage.Host != "localhost" {
		t.Errorf("Storage.Host = %v, want localhost", cfg.Storage.Host)
	}
	if cfg.Storage.Port != 5432 {
		t.Errorf("Storage.Port = %v, want 5432", cfg.Storage.Port)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Logging.Level = %v, want info", cfg.Logging.Level)
	}
	if cfg.Plugins.Directory != "plugins" {
		t.Errorf("Plugins.Directory = %v, want plugins", cfg.Plugins.Directory)
	}
	if cfg.Docs.Theme != "material" {
		t.Errorf("Docs.Theme = %v, want material", cfg.Docs.Theme)
	}
	if cfg.Agents.Port != 9443 {
		t.Errorf("Agents.Port = %v, want 9443", cfg.Agents.Port)
	}
	if cfg.MCP.Port != 9090 {
		t.Errorf("MCP.Port = %v, want 9090", cfg.MCP.Port)
	}
}

func TestFormatVersion(t *testing.T) {
	// Save original values
	origVersion := Version
	origCommitHash := CommitHash
	origBuildDate := BuildDate
	defer func() {
		Version = origVersion
		CommitHash = origCommitHash
		BuildDate = origBuildDate
	}()

	// Set test values
	Version = "1.0.0"
	CommitHash = "abc123"
	BuildDate = "2024-01-01"

	result := FormatVersion()

	if !strings.Contains(result, "knowledge-tree") {
		t.Error("FormatVersion() should contain 'knowledge-tree'")
	}
	if !strings.Contains(result, "1.0.0") {
		t.Error("FormatVersion() should contain version")
	}
	if !strings.Contains(result, "abc123") {
		t.Error("FormatVersion() should contain commit hash")
	}
}

func TestCLIScopeStatus_Constants(t *testing.T) {
	statuses := []CLIScopeStatus{
		CLIScopeStatusActive,
		CLIScopeStatusInactive,
		CLIScopeStatusRunning,
		CLIScopeStatusError,
		CLIScopeStatusUnknown,
	}

	expected := []string{"active", "inactive", "running", "error", "unknown"}

	for i, status := range statuses {
		if string(status) != expected[i] {
			t.Errorf("CLIScopeStatus[%d] = %v, want %v", i, status, expected[i])
		}
	}
}

func TestRunStatus_Constants(t *testing.T) {
	statuses := []RunStatus{
		RunStatusRunning,
		RunStatusCompleted,
		RunStatusFailed,
		RunStatusPartial,
	}

	expected := []string{"running", "completed", "failed", "partial"}

	for i, status := range statuses {
		if string(status) != expected[i] {
			t.Errorf("RunStatus[%d] = %v, want %v", i, status, expected[i])
		}
	}
}

func TestGraphQueryResult(t *testing.T) {
	result := GraphQueryResult{
		Columns: []string{"id", "name", "type"},
		Rows: [][]interface{}{
			{"1", "resource1", "aws.ec2.instance"},
			{"2", "resource2", "aws.s3.bucket"},
		},
		Count: 2,
	}

	if len(result.Columns) != 3 {
		t.Errorf("Columns length = %v, want 3", len(result.Columns))
	}
	if len(result.Rows) != 2 {
		t.Errorf("Rows length = %v, want 2", len(result.Rows))
	}
	if result.Count != 2 {
		t.Errorf("Count = %v, want 2", result.Count)
	}
}

func TestDiscoveryError(t *testing.T) {
	err := DiscoveryError{
		Code:      "TEST_ERROR",
		Message:   "Something went wrong",
		Resource:  "resource-123",
		Retryable: true,
	}

	if err.Code != "TEST_ERROR" {
		t.Errorf("Code = %v, want TEST_ERROR", err.Code)
	}
	if err.Message != "Something went wrong" {
		t.Errorf("Message = %v, want Something went wrong", err.Message)
	}
	if err.Resource != "resource-123" {
		t.Errorf("Resource = %v, want resource-123", err.Resource)
	}
	if !err.Retryable {
		t.Error("Retryable should be true")
	}
}

func TestPluginInfo(t *testing.T) {
	info := PluginInfo{
		Name:            "aws",
		Version:         "1.0.0",
		Description:     "AWS discovery plugin",
		Capabilities:    []string{"compute", "storage"},
		ConfigSchema:    `{}`,
		CredentialTypes: []string{"aws_access_key", "aws_iam_role"},
		Path:            "/plugins/aws",
		Status:          "active",
	}

	if info.Name != "aws" {
		t.Errorf("Name = %v, want aws", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("Version = %v, want 1.0.0", info.Version)
	}
	if len(info.Capabilities) != 2 {
		t.Errorf("Capabilities length = %v, want 2", len(info.Capabilities))
	}
}

func TestAgentDeployment(t *testing.T) {
	deployment := AgentDeployment{
		Host:      "192.168.1.100",
		Status:    "online",
		Version:   "1.0.0",
		ScopeName: "production",
	}

	if deployment.Host != "192.168.1.100" {
		t.Errorf("Host = %v, want 192.168.1.100", deployment.Host)
	}
	if deployment.Status != "online" {
		t.Errorf("Status = %v, want online", deployment.Status)
	}
}

func TestGraphExportOptions(t *testing.T) {
	opts := GraphExportOptions{
		Format:      "json",
		ScopeFilter: "aws-prod",
		Depth:       3,
	}

	if opts.Format != "json" {
		t.Errorf("Format = %v, want json", opts.Format)
	}
	if opts.ScopeFilter != "aws-prod" {
		t.Errorf("ScopeFilter = %v, want aws-prod", opts.ScopeFilter)
	}
	if opts.Depth != 3 {
		t.Errorf("Depth = %v, want 3", opts.Depth)
	}
}
