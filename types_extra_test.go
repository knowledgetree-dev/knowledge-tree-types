package types

import (
	"strings"
	"testing"
	"time"
)

func TestDiscoveryRun_Duration(t *testing.T) {
	tests := []struct {
		name       string
		run        DiscoveryRun
		wantZero   bool
		wantGtZero bool
	}{
		{
			name:     "not started",
			run:      DiscoveryRun{},
			wantZero: true,
		},
		{
			name: "completed",
			run: DiscoveryRun{
				StartedAt:   time.Now().Add(-5 * time.Minute),
				CompletedAt: time.Now(),
			},
			wantGtZero: true,
		},
		{
			name: "still running",
			run: DiscoveryRun{
				StartedAt: time.Now().Add(-2 * time.Minute),
			},
			wantGtZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.run.Duration()
			if tt.wantZero && d != 0 {
				t.Errorf("Duration() = %v, want zero", d)
			}
			if tt.wantGtZero && d <= 0 {
				t.Errorf("Duration() = %v, want > 0", d)
			}
		})
	}
}

func TestDiscoveryStatus_Constants(t *testing.T) {
	statuses := map[string]DiscoveryStatus{
		"PENDING":   DiscoveryStatusPending,
		"RUNNING":   DiscoveryStatusRunning,
		"COMPLETED": DiscoveryStatusCompleted,
		"FAILED":    DiscoveryStatusFailed,
	}
	for want, got := range statuses {
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestScopeStatus_Constants(t *testing.T) {
	if ScopeStatusActive != "active" {
		t.Errorf("ScopeStatusActive = %q", ScopeStatusActive)
	}
	if ScopeStatusPaused != "paused" {
		t.Errorf("ScopeStatusPaused = %q", ScopeStatusPaused)
	}
	if ScopeStatusError != "error" {
		t.Errorf("ScopeStatusError = %q", ScopeStatusError)
	}
}

func TestContextPackItemTypes(t *testing.T) {
	types := map[string]ContextPackItemType{
		"resource":      ContextItemTypeResource,
		"relationship":  ContextItemTypeRelationship,
		"architecture":  ContextItemTypeArchitecture,
		"documentation": ContextItemTypeDocumentation,
		"metrics":       ContextItemTypeMetrics,
		"config":        ContextItemTypeConfig,
	}
	for want, got := range types {
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

// testConfig returns a valid config for testing, clearing the plugin directory
// to avoid filesystem dependency.
func testConfig() *Config {
	cfg := DefaultConfig()
	cfg.Plugins.Directory = "" // clear to avoid fs check
	return cfg
}

func TestConfig_Validate_StorageErrors(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr string
	}{
		{
			name: "missing host",
			modify: func(c *Config) {
				c.Storage.Backend = "postgres"
				c.Storage.Host = ""
			},
			wantErr: "storage.host is required",
		},
		{
			name: "invalid port zero",
			modify: func(c *Config) {
				c.Storage.Backend = "postgres"
				c.Storage.Host = "localhost"
				c.Storage.Port = 0
			},
			wantErr: "storage.port must be between",
		},
		{
			name: "invalid port too large",
			modify: func(c *Config) {
				c.Storage.Backend = "postgres"
				c.Storage.Host = "localhost"
				c.Storage.Port = 70000
			},
			wantErr: "storage.port must be between",
		},
		{
			name: "missing database",
			modify: func(c *Config) {
				c.Storage.Backend = "postgres"
				c.Storage.Host = "localhost"
				c.Storage.Port = 5432
				c.Storage.Database = ""
			},
			wantErr: "storage.database is required",
		},
		{
			name: "missing user",
			modify: func(c *Config) {
				c.Storage.Backend = "postgres"
				c.Storage.Host = "localhost"
				c.Storage.Port = 5432
				c.Storage.Database = "db"
				c.Storage.User = ""
			},
			wantErr: "storage.user is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfig()
			tt.modify(cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestConfig_Validate_LoggingLevel(t *testing.T) {
	tests := []struct {
		level   string
		wantErr bool
	}{
		{"debug", false},
		{"info", false},
		{"warn", false},
		{"error", false},
		{"", false},
		{"trace", true},
		{"verbose", true},
	}
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			cfg := testConfig()
			cfg.Logging.Level = tt.level
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Validate_MCPPort(t *testing.T) {
	cfg := testConfig()
	cfg.MCP.Port = -1
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "mcp.port") {
		t.Errorf("expected mcp.port error, got %v", err)
	}

	cfg.MCP.Port = 70000
	err = cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "mcp.port") {
		t.Errorf("expected mcp.port error, got %v", err)
	}

	cfg.MCP.Port = 9090
	err = cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfig_Validate_AgentsPort(t *testing.T) {
	cfg := testConfig()
	cfg.Agents.Port = -1
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "agents.port") {
		t.Errorf("expected agents.port error, got %v", err)
	}
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	cfg := testConfig()
	err := cfg.Validate()
	if err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}

func TestConfig_Validate_InvalidScopeSchedule(t *testing.T) {
	cfg := testConfig()
	cfg.Scopes["test"] = Scope{
		Name:     "test",
		Plugin:   "aws",
		Schedule: "not-valid-cron",
	}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "schedule is invalid") {
		t.Errorf("expected schedule error, got %v", err)
	}
}

func TestConfig_Validate_ValidScopeSchedule(t *testing.T) {
	cfg := testConfig()
	cfg.Scopes["test"] = Scope{
		Name:     "test",
		Plugin:   "aws",
		Schedule: "0 */6 * * *",
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("valid schedule should pass: %v", err)
	}
}

func TestConfig_Validate_SqliteBackend(t *testing.T) {
	cfg := testConfig()
	cfg.Storage.Backend = "sqlite"
	cfg.Storage.Path = "/tmp/test.db"
	// sqlite doesn't require host/port/user/database
	err := cfg.Validate()
	if err != nil {
		t.Errorf("sqlite backend should not require postgres fields: %v", err)
	}
}

func TestStorageConfig_DSNSafe(t *testing.T) {
	s := StorageConfig{
		Backend:  "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "super-secret",
		Database: "mydb",
		SSLMode:  "require",
	}
	dsn := s.DSNSafe()
	if strings.Contains(dsn, "super-secret") {
		t.Error("DSNSafe should not contain the password")
	}
	if !strings.Contains(dsn, "REDACTED") {
		t.Error("DSNSafe should contain REDACTED")
	}
	if !strings.Contains(dsn, "admin") {
		t.Error("DSNSafe should contain user")
	}
}

func TestStorageConfig_DSNSafe_Sqlite(t *testing.T) {
	s := StorageConfig{Backend: "sqlite", Path: "/tmp/test.db"}
	dsn := s.DSNSafe()
	if dsn != "/tmp/test.db" {
		t.Errorf("DSNSafe() = %q, want %q", dsn, "/tmp/test.db")
	}
}

func TestDefaultSecurityConfig(t *testing.T) {
	cfg := DefaultSecurityConfig()

	// JWT defaults
	if cfg.JWT.AccessTokenTTL != 15*time.Minute {
		t.Errorf("JWT.AccessTokenTTL = %v", cfg.JWT.AccessTokenTTL)
	}
	if cfg.JWT.Issuer != "knowledge-tree" {
		t.Errorf("JWT.Issuer = %v", cfg.JWT.Issuer)
	}

	// Session defaults
	if !cfg.Session.Secure {
		t.Error("Session.Secure should be true")
	}
	if !cfg.Session.HttpOnly {
		t.Error("Session.HttpOnly should be true")
	}
	if cfg.Session.SameSite != "Strict" {
		t.Errorf("Session.SameSite = %v", cfg.Session.SameSite)
	}

	// Rate limit defaults
	if !cfg.RateLimit.Enabled {
		t.Error("RateLimit.Enabled should be true")
	}
	if cfg.RateLimit.Requests != 100 {
		t.Errorf("RateLimit.Requests = %v", cfg.RateLimit.Requests)
	}

	// Encryption defaults
	if cfg.Encryption.KeyRotationDays != 90 {
		t.Errorf("Encryption.KeyRotationDays = %v", cfg.Encryption.KeyRotationDays)
	}

	// TLS defaults
	if !cfg.TLS.Enabled {
		t.Error("TLS.Enabled should be true")
	}
	if cfg.TLS.MinVersion != "1.3" {
		t.Errorf("TLS.MinVersion = %v", cfg.TLS.MinVersion)
	}

	// Audit defaults
	if !cfg.Audit.Enabled {
		t.Error("Audit.Enabled should be true")
	}

	// Password policy defaults
	if cfg.PasswordPolicy.MinLength != 12 {
		t.Errorf("PasswordPolicy.MinLength = %v", cfg.PasswordPolicy.MinLength)
	}
	if !cfg.PasswordPolicy.RequireUpper {
		t.Error("PasswordPolicy.RequireUpper should be true")
	}
}

func TestDefaultRoles(t *testing.T) {
	roles := DefaultRoles()
	if len(roles) != 3 {
		t.Fatalf("expected 3 default roles, got %d", len(roles))
	}

	viewer, ok := roles["viewer"]
	if !ok {
		t.Fatal("viewer role missing")
	}
	if len(viewer.Permissions) == 0 {
		t.Error("viewer should have permissions")
	}

	editor, ok := roles["editor"]
	if !ok {
		t.Fatal("editor role missing")
	}
	if len(editor.Permissions) <= len(viewer.Permissions) {
		t.Error("editor should have more permissions than viewer")
	}

	admin, ok := roles["admin"]
	if !ok {
		t.Fatal("admin role missing")
	}
	if len(admin.Permissions) <= len(editor.Permissions) {
		t.Error("admin should have more permissions than editor")
	}
}

func TestPermissionConstants(t *testing.T) {
	perms := []string{
		PermissionReadResources, PermissionWriteResources, PermissionDeleteResources,
		PermissionReadGraph, PermissionWriteGraph,
		PermissionReadDiscovery, PermissionRunDiscovery,
		PermissionReadAdmin, PermissionWriteAdmin,
		PermissionReadAudit,
		PermissionExportData, PermissionImportData,
	}
	for _, p := range perms {
		if p == "" {
			t.Error("permission constant should not be empty")
		}
		if !strings.Contains(p, ":") {
			t.Errorf("permission %q should contain ':'", p)
		}
	}
}

func TestDiscoveryScope_Fields(t *testing.T) {
	scope := DiscoveryScope{
		ID:           "scope-1",
		Name:         "Production AWS",
		PluginName:   "aws",
		Config:       map[string]string{"region": "us-east-1"},
		Schedule:     "0 */6 * * *",
		Status:       ScopeStatusActive,
		ResourceCount: 42,
	}
	if scope.ID != "scope-1" {
		t.Errorf("ID = %q", scope.ID)
	}
	if scope.PluginName != "aws" {
		t.Errorf("PluginName = %q", scope.PluginName)
	}
	if scope.Status != ScopeStatusActive {
		t.Errorf("Status = %q", scope.Status)
	}
	if scope.ResourceCount != 42 {
		t.Errorf("ResourceCount = %d", scope.ResourceCount)
	}
}

func TestCredentialSource_Fields(t *testing.T) {
	cs := CredentialSource{
		Type:   "vault",
		Path:   "secret/data/aws",
		Params: map[string]string{"role": "knowledge-tree"},
	}
	if cs.Type != "vault" {
		t.Errorf("Type = %q", cs.Type)
	}
	if cs.Path != "secret/data/aws" {
		t.Errorf("Path = %q", cs.Path)
	}
	if cs.Params["role"] != "knowledge-tree" {
		t.Errorf("Params[role] = %q", cs.Params["role"])
	}
}

func TestDiscoveryRun_Fields(t *testing.T) {
	now := time.Now()
	run := DiscoveryRun{
		ID:                 "run-1",
		ScopeID:            "scope-1",
		Status:             DiscoveryStatusCompleted,
		StartedAt:          now.Add(-5 * time.Minute),
		CompletedAt:        now,
		ResourcesFound:     100,
		RelationshipsFound: 50,
		Errors:             []string{"warning: timeout on node-3"},
	}
	if run.ResourcesFound != 100 {
		t.Errorf("ResourcesFound = %d", run.ResourcesFound)
	}
	if len(run.Errors) != 1 {
		t.Errorf("Errors length = %d", len(run.Errors))
	}
}

func TestContextPack_Fields(t *testing.T) {
	pack := ContextPack{
		Name:         "test-pack",
		TargetTokens: 4000,
		Items: []ContextPackItem{
			{Type: ContextItemTypeResource, Content: "test", Tokens: 10},
		},
		TotalTokens: 10,
	}
	if pack.Name != "test-pack" {
		t.Errorf("Name = %q", pack.Name)
	}
	if len(pack.Items) != 1 {
		t.Errorf("Items length = %d", len(pack.Items))
	}
	if pack.Items[0].Type != ContextItemTypeResource {
		t.Errorf("Item Type = %q", pack.Items[0].Type)
	}
}
