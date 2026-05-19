package types

// Relationship type constants describe the kinds of edges that can exist
// between resources in the Knowledge Tree graph.
const (
	// RelContains indicates that the source resource logically contains the
	// target (e.g., a VPC contains a subnet, a cluster contains a node pool).
	RelContains = "CONTAINS"

	// RelConnectsTo indicates a network-level connection between two resources
	// (e.g., a security group allows traffic to an instance, a load balancer
	// connects to a target group).
	RelConnectsTo = "CONNECTS_TO"

	// RelDependsOn indicates that the source resource depends on the target
	// for its operation (e.g., a Lambda function depends on an IAM role).
	RelDependsOn = "DEPENDS_ON"

	// RelHasRole indicates that the source resource has an associated IAM or
	// security role (e.g., an EC2 instance has an instance profile/role).
	RelHasRole = "HAS_ROLE"

	// RelHasPolicy indicates that the source resource has an attached policy
	// (e.g., an IAM role has an inline or managed policy).
	RelHasPolicy = "HAS_POLICY"

	// RelExposes indicates that the source resource exposes the target
	// (e.g., a load balancer exposes a service, an ingress exposes a pod).
	RelExposes = "EXPOSES"

	// RelRoutesTo indicates a routing relationship (e.g., a route table routes
	// to a NAT gateway, an ingress routes to a service).
	RelRoutesTo = "ROUTES_TO"

	// RelPeeredWith indicates a peering connection between two network
	// resources (e.g., VPC peering).
	RelPeeredWith = "PEERED_WITH"

	// RelRunsOn indicates that the source resource runs on the target
	// (e.g., a pod runs on a node, a container runs on an EC2 instance).
	RelRunsOn = "RUNS_ON"

	// RelBackedBy indicates that the source is backed by the target
	// (e.g., an RDS instance is backed by an EBS volume, a PVC is backed
	// by a PV).
	RelBackedBy = "BACKED_BY"

	// RelManagedBy indicates that the source resource is managed or
	// controlled by the target (e.g., a deployment manages replica pods,
	// CloudFormation manages a stack of resources).
	RelManagedBy = "MANAGED_BY"
)

// Relationship represents an edge between two resources in the Knowledge
// Tree graph. Relationships encode how infrastructure components are
// connected, dependent on, or grouped with each other.
type Relationship struct {
	// SourceID is the globally unique identifier of the source resource
	// (the origin of the directed edge).
	SourceID string `json:"source_id"`

	// TargetID is the globally unique identifier of the target resource
	// (the destination of the directed edge).
	TargetID string `json:"target_id"`

	// Type describes the kind of relationship using one of the Rel*
	// constants (e.g., RelContains, RelDependsOn).
	Type string `json:"type"`

	// Properties holds optional key-value metadata about this relationship
	// (e.g., port numbers for a CONNECTS_TO, access level for a HAS_ROLE).
	Properties map[string]string `json:"properties,omitempty"`
}

// Key returns a unique identifier for the relationship, composed of the
// source ID, type, and target ID. This is suitable for deduplication and
// as an edge key in the graph store.
func (r *Relationship) Key() string {
	return r.SourceID + "-" + r.Type + "-" + r.TargetID
}
