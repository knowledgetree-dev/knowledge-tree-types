package types

import (
	"time"
)

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	// Authentication
	JWT JWTConfig `yaml:"jwt"`
	
	// Session management
	Session SessionConfig `yaml:"session"`
	
	// Rate limiting
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	
	// Encryption
	Encryption EncryptionConfig `yaml:"encryption"`
	
	// CORS
	CORS CORSConfig `yaml:"cors"`
	
	// TLS
	TLS TLSConfig `yaml:"tls"`
	
	// Audit logging
	Audit AuditConfig `yaml:"audit"`
	
	// Password policy
	PasswordPolicy PasswordPolicyConfig `yaml:"password_policy"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string        `yaml:"secret" env:"KT_JWT_SECRET"`
	SecretFile       string        `yaml:"secret_file" env:"KT_JWT_SECRET_FILE"`
	AccessTokenTTL   time.Duration `yaml:"access_token_ttl" env:"KT_JWT_ACCESS_TTL"`
	RefreshTokenTTL  time.Duration `yaml:"refresh_token_ttl" env:"KT_JWT_REFRESH_TTL"`
	Issuer           string        `yaml:"issuer" env:"KT_JWT_ISSUER"`
	Audience         []string      `yaml:"audience" env:"KT_JWT_AUDIENCE"`
}

// SessionConfig holds session configuration
type SessionConfig struct {
	MaxAge       time.Duration `yaml:"max_age" env:"KT_SESSION_MAX_AGE"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"KT_SESSION_IDLE_TIMEOUT"`
	Secure       bool          `yaml:"secure" env:"KT_SESSION_SECURE"`
	HttpOnly     bool          `yaml:"http_only" env:"KT_SESSION_HTTP_ONLY"`
	SameSite     string        `yaml:"same_site" env:"KT_SESSION_SAME_SITE"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool          `yaml:"enabled" env:"KT_RATE_LIMIT_ENABLED"`
	Requests    int           `yaml:"requests" env:"KT_RATE_LIMIT_REQUESTS"`
	Window      time.Duration `yaml:"window" env:"KT_RATE_LIMIT_WINDOW"`
	Burst       int           `yaml:"burst" env:"KT_RATE_LIMIT_BURST"`
}

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	MasterKey      string `yaml:"master_key" env:"KT_ENCRYPTION_KEY"`
	MasterKeyFile  string `yaml:"master_key_file" env:"KT_ENCRYPTION_KEY_FILE"`
	KeyRotationDays int   `yaml:"key_rotation_days" env:"KT_KEY_ROTATION_DAYS"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins" env:"KT_CORS_ORIGINS"`
	AllowedMethods   []string `yaml:"allowed_methods" env:"KT_CORS_METHODS"`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"KT_CORS_HEADERS"`
	ExposedHeaders   []string `yaml:"exposed_headers" env:"KT_CORS_EXPOSED_HEADERS"`
	AllowCredentials bool     `yaml:"allow_credentials" env:"KT_CORS_CREDENTIALS"`
	MaxAge           int      `yaml:"max_age" env:"KT_CORS_MAX_AGE"`
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled     bool   `yaml:"enabled" env:"KT_TLS_ENABLED"`
	CertFile    string `yaml:"cert_file" env:"KT_TLS_CERT_FILE"`
	KeyFile     string `yaml:"key_file" env:"KT_TLS_KEY_FILE"`
	MinVersion  string `yaml:"min_version" env:"KT_TLS_MIN_VERSION"`
	RequireClientCert bool `yaml:"require_client_cert" env:"KT_TLS_REQUIRE_CLIENT_CERT"`
	ClientCAFile string `yaml:"client_ca_file" env:"KT_TLS_CLIENT_CA_FILE"`
}

// AuditConfig holds audit logging configuration
type AuditConfig struct {
	Enabled      bool   `yaml:"enabled" env:"KT_AUDIT_ENABLED"`
	LogFile      string `yaml:"log_file" env:"KT_AUDIT_LOG_FILE"`
	MaxSizeMB    int    `yaml:"max_size_mb" env:"KT_AUDIT_MAX_SIZE_MB"`
	MaxBackups   int    `yaml:"max_backups" env:"KT_AUDIT_MAX_BACKUPS"`
	MaxAgeDays   int    `yaml:"max_age_days" env:"KT_AUDIT_MAX_AGE_DAYS"`
	Compress     bool   `yaml:"compress" env:"KT_AUDIT_COMPRESS"`
}

// PasswordPolicyConfig holds password policy configuration
type PasswordPolicyConfig struct {
	MinLength      int  `yaml:"min_length" env:"KT_PASSWORD_MIN_LENGTH"`
	RequireUpper   bool `yaml:"require_upper" env:"KT_PASSWORD_REQUIRE_UPPER"`
	RequireLower   bool `yaml:"require_lower" env:"KT_PASSWORD_REQUIRE_LOWER"`
	RequireNumber  bool `yaml:"require_number" env:"KT_PASSWORD_REQUIRE_NUMBER"`
	RequireSpecial bool `yaml:"require_special" env:"KT_PASSWORD_REQUIRE_SPECIAL"`
	MaxAgeDays     int  `yaml:"max_age_days" env:"KT_PASSWORD_MAX_AGE_DAYS"`
	HistoryCount   int  `yaml:"history_count" env:"KT_PASSWORD_HISTORY_COUNT"`
}

// DefaultSecurityConfig returns secure defaults
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		JWT: JWTConfig{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
			Issuer:          "knowledge-tree",
			Audience:        []string{"knowledge-tree-api"},
		},
		Session: SessionConfig{
			MaxAge:      24 * time.Hour,
			IdleTimeout: 30 * time.Minute,
			Secure:      true,
			HttpOnly:    true,
			SameSite:    "Strict",
		},
		RateLimit: RateLimitConfig{
			Enabled:  true,
			Requests: 100,
			Window:   time.Minute,
			Burst:    150,
		},
		Encryption: EncryptionConfig{
			KeyRotationDays: 90,
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
			ExposedHeaders:   []string{"X-Request-ID"},
			AllowCredentials: false,
			MaxAge:           300,
		},
		TLS: TLSConfig{
			Enabled:    true,
			MinVersion: "1.3",
		},
		Audit: AuditConfig{
			Enabled:    true,
			LogFile:    "logs/audit.log",
			MaxSizeMB:  100,
			MaxBackups: 10,
			MaxAgeDays: 90,
			Compress:   true,
		},
		PasswordPolicy: PasswordPolicyConfig{
			MinLength:      12,
			RequireUpper:   true,
			RequireLower:   true,
			RequireNumber:  true,
			RequireSpecial: true,
			MaxAgeDays:     90,
			HistoryCount:   5,
		},
	}
}

// RBACConfig holds role-based access control configuration
type RBACConfig struct {
	Enabled     bool              `yaml:"enabled" env:"KT_RBAC_ENABLED"`
	DefaultRole string            `yaml:"default_role" env:"KT_RBAC_DEFAULT_ROLE"`
	Roles       map[string]Role   `yaml:"roles"`
}

// Role defines a role with permissions
type Role struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Permissions []string `yaml:"permissions"`
	Inherits    []string `yaml:"inherits"`
}

// Permission constants
const (
	PermissionReadResources   = "resources:read"
	PermissionWriteResources  = "resources:write"
	PermissionDeleteResources = "resources:delete"
	PermissionReadGraph       = "graph:read"
	PermissionWriteGraph      = "graph:write"
	PermissionReadDiscovery   = "discovery:read"
	PermissionRunDiscovery    = "discovery:run"
	PermissionReadAdmin       = "admin:read"
	PermissionWriteAdmin      = "admin:write"
	PermissionReadAudit       = "audit:read"
	PermissionExportData      = "data:export"
	PermissionImportData      = "data:import"
)

// DefaultRoles returns default RBAC roles
func DefaultRoles() map[string]Role {
	return map[string]Role{
		"viewer": {
			Name:        "viewer",
			Description: "Read-only access to resources",
			Permissions: []string{
				PermissionReadResources,
				PermissionReadGraph,
				PermissionReadDiscovery,
			},
		},
		"editor": {
			Name:        "editor",
			Description: "Can read and modify resources",
			Permissions: []string{
				PermissionReadResources,
				PermissionWriteResources,
				PermissionReadGraph,
				PermissionWriteGraph,
				PermissionReadDiscovery,
				PermissionRunDiscovery,
			},
		},
		"admin": {
			Name:        "admin",
			Description: "Full access to all features",
			Permissions: []string{
				PermissionReadResources,
				PermissionWriteResources,
				PermissionDeleteResources,
				PermissionReadGraph,
				PermissionWriteGraph,
				PermissionReadDiscovery,
				PermissionRunDiscovery,
				PermissionReadAdmin,
				PermissionWriteAdmin,
				PermissionReadAudit,
				PermissionExportData,
				PermissionImportData,
			},
		},
	}
}

// User represents a stored user account
type User struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`         // never serialized
	Role             string    `json:"role"`
	IdentityProvider string    `json:"identity_provider,omitempty"` // e.g. "google", "github"
	ExternalID       string    `json:"external_id,omitempty"`       // OIDC sub claim
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// SSOProvider holds OIDC provider configuration
type SSOProvider struct {
	Name        string            `json:"name"`
	Label       string            `json:"label"`
	IssuerURL   string            `json:"issuer_url"`
	ClientID    string            `json:"client_id"`
	ClientSecret string           `json:"client_secret,omitempty"`
	Scopes      []string          `json:"scopes,omitempty"`
	RoleMapping map[string]string `json:"role_mapping,omitempty"` // OIDC group → app role
	Enabled     bool              `json:"enabled"`
}

// SAMLProvider holds SAML identity provider configuration
type SAMLProvider struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	EntityID    string `json:"entity_id"`
	SSOURL      string `json:"sso_url"`
	Certificate string `json:"certificate,omitempty"` // IdP cert PEM
	SPCert      string `json:"sp_cert,omitempty"`     // SP cert PEM
	SPKey       string `json:"sp_key,omitempty"`      // SP private key PEM
	Enabled     bool   `json:"enabled"`
}

// CreateUserRequest is the request body for creating a user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UpdateUserRequest is the request body for updating a user
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

// Role constants for user management
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// ValidRole checks if a role string is one of the allowed roles
func ValidRole(r string) bool {
	switch r {
	case RoleAdmin, RoleEditor, RoleViewer:
		return true
	default:
		return false
	}
}
