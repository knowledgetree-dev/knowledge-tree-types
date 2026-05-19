// Package types defines the core domain types for Knowledge Tree.
//
// These types represent infrastructure resources, relationships between them,
// discovery scopes and runs, and context packs used for AI-generated
// documentation. They are used throughout the system by the orchestrator,
// storage layer, plugin system, and API.
package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// --------------------------------------------------------------------------
// Resource type constants
// --------------------------------------------------------------------------

// AWS resource types.
const (
	TypeEC2Instance       = "aws.ec2.instance"
	TypeEC2SecurityGroup  = "aws.ec2.security_group"
	TypeEC2EIP            = "aws.ec2.elastic_ip"
	TypeEC2Volume         = "aws.ec2.volume"
	TypeEC2LaunchTemplate = "aws.ec2.launch_template"
	TypeS3Bucket          = "aws.s3.bucket"
	TypeRDSInstance       = "aws.rds.instance"
	TypeRDSCluster        = "aws.rds.cluster"
	TypeRDSSubnetGroup    = "aws.rds.subnet_group"
	TypeLambdaFunction    = "aws.lambda.function"
	TypeLambdaLayer       = "aws.lambda.layer"
	TypeVPC               = "aws.vpc"
	TypeSubnet            = "aws.subnet"
	TypeRouteTable        = "aws.route_table"
	TypeInternetGateway   = "aws.internet_gateway"
	TypeNATGateway        = "aws.nat_gateway"
	TypeELB               = "aws.elb.load_balancer"
	TypeELBTargetGroup    = "aws.elb.target_group"
	TypeELBListener       = "aws.elb.listener"
	TypeIAMRole           = "aws.iam.role"
	TypeIAMPolicy         = "aws.iam.policy"
	TypeIAMUser           = "aws.iam.user"
	TypeSQSQueue          = "aws.sqs.queue"
	TypeSNSTopic          = "aws.sns.topic"
	TypeDynamoDBTable     = "aws.dynamodb.table"
	TypeElasticacheCluster = "aws.elasticache.cluster"
	TypeEKSCompute        = "aws.eks.cluster"
	TypeCloudFrontDist    = "aws.cloudfront.distribution"
	TypeAPIGateway        = "aws.apigateway.rest_api"
	TypeKinesisStream     = "aws.kinesis.stream"
	TypeRoute53Zone       = "aws.route53.hosted_zone"
	TypeRoute53Record     = "aws.route53.record_set"
	TypeKMSKey            = "aws.kms.key"
)

// Azure resource types.
const (
	TypeAzureVM            = "azure.compute.virtual_machine"
	TypeAzureVNet          = "azure.network.virtual_network"
	TypeAzureSubnet        = "azure.network.subnet"
	TypeAzureNSG           = "azure.network.security_group"
	TypeAzurePublicIP      = "azure.network.public_ip"
	TypeAzureAKS           = "azure.containerservice.managed_cluster"
	TypeAzureSQLServer     = "azure.sql.server"
	TypeAzureSQLDB         = "azure.sql.database"
	TypeAzureStorageAcct   = "azure.storage.account"
	TypeAzureAppService    = "azure.web.app_service"
	TypeAzureFunctionApp   = "azure.web.function_app"
	TypeAzureKeyVault      = "azure.keyvault.vault"
	TypeAzureCosmosDB      = "azure.documentdb.database_account"
	TypeAzureResourceGroup = "azure.resources.resource_group"
	TypeAzureDisk          = "azure.compute.disk"
	TypeAzureNIC           = "azure.network.network_interface"
	TypeAzureLoadBalancer  = "azure.network.load_balancer"
	TypeAzureAppServicePlan = "azure.web.app_service_plan"
	TypeAzureSQLPool       = "azure.sql.elastic_pool"
)

// GCP resource types.
const (
	TypeGCPComputeInstance = "gcp.compute.instance"
	TypeGCPDisk            = "gcp.compute.disk"
	TypeGCPVPC             = "gcp.compute.network"
	TypeGCPSubnet          = "gcp.compute.subnetwork"
	TypeGCPFirewall        = "gcp.compute.firewall"
	TypeGCPGKECluster      = "gcp.container.cluster"
	TypeGCPGKENodePool     = "gcp.container.node_pool"
	TypeGCPCloudSQL        = "gcp.cloudsql.instance"
	TypeGCPCloudSQLInstance = "gcp.cloudsql.instance"
	TypeGCPSpanner         = "gcp.spanner.instance"
	TypeGCPBigQuery        = "gcp.bigquery.dataset"
	TypeGCPStorageBucket   = "gcp.storage.bucket"
	TypeGCPPubSubTopic     = "gcp.pubsub.topic"
	TypeGCPPubSubSub       = "gcp.pubsub.subscription"
	TypeGCPCloudFunction   = "gcp.cloudfunctions.function"
	TypeGCPRunService      = "gcp.run.service"
	TypeGCPIAMServiceAccount = "gcp.iam.service_account"
	TypeGCPIAMRole         = "gcp.iam.role"
	TypeGCPDNSZone         = "gcp.dns.managed_zone"
)

// Kubernetes resource types.
const (
	TypeK8sNamespace  = "k8s.namespace"
	TypeK8sDeployment = "k8s.deployment"
	TypeK8sService    = "k8s.service"
	TypeK8sPod        = "k8s.pod"
	TypeK8sConfigMap  = "k8s.config_map"
	TypeK8sSecret     = "k8s.secret"
	TypeK8sIngress    = "k8s.ingress"
	TypeK8sPVC        = "k8s.persistent_volume_claim"
	TypeK8sPV         = "k8s.persistent_volume"
	TypeK8sDaemonSet  = "k8s.daemon_set"
	TypeK8sStatefulSet = "k8s.stateful_set"
	TypeK8sCronJob    = "k8s.cron_job"
	TypeK8sJob        = "k8s.job"
	TypeK8sNetworkPolicy = "k8s.network_policy"
)

// ResourceType holds the metadata that describes a category of infrastructure
// resource. It is used primarily for display and categorization purposes.
type ResourceType struct {
	// ID is the unique dotted identifier for this resource type
	// (e.g., "aws.ec2.instance").
	ID string

	// Provider is the infrastructure provider (e.g., "aws", "azure", "gcp",
	// "kubernetes").
	Provider string

	// Label is a human-readable name for the resource type.
	Label string

	// Category groups related resource types (e.g., "compute", "networking",
	// "storage", "database").
	Category string
}

// Resource represents a single infrastructure resource discovered by a plugin.
// It is the fundamental unit of data in Knowledge Tree and is stored as a node
// in the graph database.
type Resource struct {
	// ID is a globally unique identifier for the resource. The format is
	// provider-specific but typically follows the pattern
	// "provider:resource_type:cloud_identifier".
	ID string `json:"id"`

	// Type is the resource type in dotted notation
	// (e.g., "aws.ec2.instance", "k8s.deployment").
	Type string `json:"type"`

	// Name is the human-readable name of the resource as configured by the
	// user or assigned by the provider.
	Name string `json:"name"`

	// Provider is the cloud or infrastructure provider
	// (e.g., "aws", "azure", "gcp", "kubernetes").
	Provider string `json:"provider"`

	// Region is the geographic region or zone where the resource is deployed.
	Region string `json:"region"`

	// AccountID is the cloud account, subscription, or project ID that owns
	// the resource.
	AccountID string `json:"account_id"`

	// Properties holds key-value pairs of resource-specific attributes such
	// as instance type, state, IP addresses, etc.
	Properties map[string]string `json:"properties"`

	// Tags are user-defined labels attached to the resource, typically used
	// for cost allocation, ownership, and organization.
	Tags []string `json:"tags"`

	// RawData contains the full JSON representation of the resource as
	// returned by the provider API. This is preserved for detailed inspection
	// and AI context generation.
	RawData json.RawMessage `json:"raw_data"`

	// DiscoveredAt is the timestamp when the resource was last discovered
	// by a plugin.
	DiscoveredAt time.Time `json:"discovered_at"`
}

// Key returns a unique identifier string for the resource, suitable for use
// as a graph node key. The key is composed of the provider, type, and ID,
// guaranteeing global uniqueness.
func (r *Resource) Key() string {
	return fmt.Sprintf("%s:%s:%s", r.Provider, r.Type, r.ID)
}

// ToNode converts the resource into a map suitable for insertion into a graph
// database. It flattens the resource fields and properties into a single
// map, making all data accessible for graph queries.
func (r *Resource) ToNode() map[string]interface{} {
	node := map[string]interface{}{
		"id":            r.ID,
		"type":          r.Type,
		"name":          r.Name,
		"provider":      r.Provider,
		"region":        r.Region,
		"account_id":    r.AccountID,
		"tags":          r.Tags,
		"discovered_at": r.DiscoveredAt.Unix(),
		"raw_data":      string(r.RawData),
	}
	for k, v := range r.Properties {
		node["prop_"+k] = v
	}
	return node
}
