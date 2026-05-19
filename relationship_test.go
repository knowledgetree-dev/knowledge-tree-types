package types

import (
	"encoding/json"
	"testing"
)

func TestRelationship_Key(t *testing.T) {
	tests := []struct {
		name         string
		relationship Relationship
		want         string
	}{
		{
			name: "VPC contains subnet",
			relationship: Relationship{
				SourceID: "vpc-123",
				TargetID: "subnet-456",
				Type:     RelContains,
			},
			want: "vpc-123-CONTAINS-subnet-456",
		},
		{
			name: "Instance connects to security group",
			relationship: Relationship{
				SourceID: "i-abc123",
				TargetID: "sg-def456",
				Type:     RelConnectsTo,
			},
			want: "i-abc123-CONNECTS_TO-sg-def456",
		},
		{
			name: "Lambda depends on IAM role",
			relationship: Relationship{
				SourceID: "lambda:my-function",
				TargetID: "role:lambda-role",
				Type:     RelDependsOn,
			},
			want: "lambda:my-function-DEPENDS_ON-role:lambda-role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.relationship.Key()
			if got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipType_Constants(t *testing.T) {
	// Test that all relationship type constants are defined correctly
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"RelContains", RelContains, "CONTAINS"},
		{"RelConnectsTo", RelConnectsTo, "CONNECTS_TO"},
		{"RelDependsOn", RelDependsOn, "DEPENDS_ON"},
		{"RelHasRole", RelHasRole, "HAS_ROLE"},
		{"RelHasPolicy", RelHasPolicy, "HAS_POLICY"},
		{"RelExposes", RelExposes, "EXPOSES"},
		{"RelRoutesTo", RelRoutesTo, "ROUTES_TO"},
		{"RelPeeredWith", RelPeeredWith, "PEERED_WITH"},
		{"RelRunsOn", RelRunsOn, "RUNS_ON"},
		{"RelBackedBy", RelBackedBy, "BACKED_BY"},
		{"RelManagedBy", RelManagedBy, "MANAGED_BY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.value, tt.want)
			}
		})
	}
}

func TestRelationship_JSONSerialization(t *testing.T) {
	rel := Relationship{
		SourceID:   "source-123",
		TargetID:   "target-456",
		Type:       RelContains,
		Properties: map[string]string{"cidr": "10.0.1.0/24"},
	}

	// Test marshaling
	data, err := json.Marshal(rel)
	if err != nil {
		t.Fatalf("Failed to marshal relationship: %v", err)
	}

	// Test unmarshaling
	var decoded Relationship
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal relationship: %v", err)
	}

	if decoded.SourceID != rel.SourceID {
		t.Errorf("SourceID mismatch: got %v, want %v", decoded.SourceID, rel.SourceID)
	}
	if decoded.TargetID != rel.TargetID {
		t.Errorf("TargetID mismatch: got %v, want %v", decoded.TargetID, rel.TargetID)
	}
	if decoded.Type != rel.Type {
		t.Errorf("Type mismatch: got %v, want %v", decoded.Type, rel.Type)
	}
	if decoded.Properties["cidr"] != "10.0.1.0/24" {
		t.Errorf("Properties mismatch: got %v, want 10.0.1.0/24", decoded.Properties["cidr"])
	}
}

func TestRelationship_KeyUniqueness(t *testing.T) {
	// Test that different relationships produce different keys
	rel1 := Relationship{SourceID: "a", TargetID: "b", Type: RelContains}
	rel2 := Relationship{SourceID: "b", TargetID: "a", Type: RelContains}
	rel3 := Relationship{SourceID: "a", TargetID: "b", Type: RelDependsOn}

	key1 := rel1.Key()
	key2 := rel2.Key()
	key3 := rel3.Key()

	if key1 == key2 {
		t.Error("Keys should be different for reversed source/target")
	}
	if key1 == key3 {
		t.Error("Keys should be different for different relationship types")
	}
}
