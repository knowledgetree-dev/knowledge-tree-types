package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestResource_Key(t *testing.T) {
	tests := []struct {
		name     string
		resource Resource
		want     string
	}{
		{
			name: "AWS EC2 instance",
			resource: Resource{
				ID:       "i-1234567890abcdef0",
				Type:     TypeEC2Instance,
				Provider: "aws",
			},
			want: "aws:aws.ec2.instance:i-1234567890abcdef0",
		},
		{
			name: "Azure VM",
			resource: Resource{
				ID:       "my-vm",
				Type:     TypeAzureVM,
				Provider: "azure",
			},
			want: "azure:azure.compute.virtual_machine:my-vm",
		},
		{
			name: "GCP Compute Instance",
			resource: Resource{
				ID:       "my-instance",
				Type:     TypeGCPComputeInstance,
				Provider: "gcp",
			},
			want: "gcp:gcp.compute.instance:my-instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.resource.Key()
			if got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_ToNode(t *testing.T) {
	resource := Resource{
		ID:           "test-id",
		Type:         TypeEC2Instance,
		Name:         "test-instance",
		Provider:     "aws",
		Region:       "us-east-1",
		AccountID:    "123456789012",
		Properties:   map[string]string{"instance_type": "t3.medium"},
		Tags:         []string{"env:prod", "team:platform"},
		RawData:      json.RawMessage(`{"InstanceId": "test-id"}`),
		DiscoveredAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	node := resource.ToNode()

	if node["id"] != "test-id" {
		t.Errorf("ToNode() id = %v, want test-id", node["id"])
	}
	if node["type"] != TypeEC2Instance {
		t.Errorf("ToNode() type = %v, want %v", node["type"], TypeEC2Instance)
	}
	if node["name"] != "test-instance" {
		t.Errorf("ToNode() name = %v, want test-instance", node["name"])
	}
	if node["provider"] != "aws" {
		t.Errorf("ToNode() provider = %v, want aws", node["provider"])
	}
	if node["region"] != "us-east-1" {
		t.Errorf("ToNode() region = %v, want us-east-1", node["region"])
	}
	if node["account_id"] != "123456789012" {
		t.Errorf("ToNode() account_id = %v, want 123456789012", node["account_id"])
	}
	if node["prop_instance_type"] != "t3.medium" {
		t.Errorf("ToNode() prop_instance_type = %v, want t3.medium", node["prop_instance_type"])
	}
}

func TestResourceType_Constants(t *testing.T) {
	// Test that all resource type constants are defined and non-empty
	constants := []string{
		TypeEC2Instance,
		TypeEC2SecurityGroup,
		TypeS3Bucket,
		TypeRDSInstance,
		TypeVPC,
		TypeSubnet,
		TypeELB,
		TypeIAMRole,
		TypeAzureVM,
		TypeAzureVNet,
		TypeAzureAKS,
		TypeAzureSQLServer,
		TypeAzureStorageAcct,
		TypeGCPComputeInstance,
		TypeGCPVPC,
		TypeGCPGKECluster,
		TypeGCPStorageBucket,
		TypeK8sNamespace,
		TypeK8sDeployment,
		TypeK8sService,
		TypeK8sPod,
	}

	for _, c := range constants {
		if c == "" {
			t.Error("Resource type constant is empty")
		}
	}
}

func TestResource_JSONSerialization(t *testing.T) {
	resource := Resource{
		ID:        "test-id",
		Type:      TypeEC2Instance,
		Name:      "test-instance",
		Provider:  "aws",
		Region:    "us-east-1",
		AccountID: "123456789012",
		Properties: map[string]string{
			"instance_type": "t3.medium",
			"state":         "running",
		},
		Tags:         []string{"env:prod"},
		RawData:      json.RawMessage(`{"InstanceId":"test-id"}`),
		DiscoveredAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Test marshaling
	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal resource: %v", err)
	}

	// Test unmarshaling
	var decoded Resource
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal resource: %v", err)
	}

	if decoded.ID != resource.ID {
		t.Errorf("ID mismatch: got %v, want %v", decoded.ID, resource.ID)
	}
	if decoded.Type != resource.Type {
		t.Errorf("Type mismatch: got %v, want %v", decoded.Type, resource.Type)
	}
	if decoded.Name != resource.Name {
		t.Errorf("Name mismatch: got %v, want %v", decoded.Name, resource.Name)
	}
}
