package types

import "time"

// ContextPackItemType classifies the kind of content a ContextPackItem holds.
// This helps the AI system understand the nature of each piece of context.
type ContextPackItemType string

const (
	// ContextItemTypeResource indicates the item contains a description of an
	// infrastructure resource (e.g., an EC2 instance's type, state, and tags).
	ContextItemTypeResource ContextPackItemType = "resource"

	// ContextItemTypeRelationship indicates the item describes a relationship
	// between resources (e.g., a subnet is contained in a VPC).
	ContextItemTypeRelationship ContextPackItemType = "relationship"

	// ContextItemTypeArchitecture indicates the item provides an architectural
	// overview or diagram description of a system component.
	ContextItemTypeArchitecture ContextPackItemType = "architecture"

	// ContextItemTypeDocumentation indicates the item contains existing
	// documentation text that should be used as reference or updated.
	ContextItemTypeDocumentation ContextPackItemType = "documentation"

	// ContextItemTypeMetrics indicates the item contains operational metrics
	// or telemetry data about a resource or service.
	ContextItemTypeMetrics ContextPackItemType = "metrics"

	// ContextItemTypeConfig indicates the item contains configuration data
	// relevant to the documentation context.
	ContextItemTypeConfig ContextPackItemType = "config"
)

// ContextPack is a curated collection of context items assembled for AI-based
// documentation generation. The pack is built by querying the graph database
// and selecting the most relevant resources, relationships, and existing docs
// to fit within a target token budget.
type ContextPack struct {
	// Name is a descriptive name for this context pack
	// (e.g., "production-api-service-documentation").
	Name string `json:"name"`

	// TargetTokens is the maximum number of tokens the pack should contain
	// when sent to the AI model.
	TargetTokens int `json:"target_tokens"`

	// Items is the ordered list of context items included in this pack.
	// Items are typically sorted by relevance.
	Items []ContextPackItem `json:"items"`

	// TotalTokens is the sum of all item token counts after assembly.
	TotalTokens int `json:"total_tokens"`

	// LastUpdated is the timestamp when the pack was last (re)assembled.
	LastUpdated time.Time `json:"last_updated"`
}

// ContextPackItem represents a single piece of context to be included in an
// AI prompt. Each item has a type, the actual content string, and a precomputed
// token count for budget management.
type ContextPackItem struct {
	// Type classifies the content using one of the ContextPackItemType
	// constants.
	Type ContextPackItemType `json:"type"`

	// Content is the textual representation of this context item, formatted
	// for inclusion in an AI prompt.
	Content string `json:"content"`

	// Tokens is the estimated number of tokens this item consumes in the
	// AI model's context window.
	Tokens int `json:"tokens"`
}
