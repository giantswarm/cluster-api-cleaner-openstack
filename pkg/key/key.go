package key

const (
	ClusterTagPrefix     = "giant_swarm_cluster"
	CapiClusterLabelKey  = "cluster.x-k8s.io/cluster-name"
	CleanerFinalizerName = "cluster-api-cleaner-openstack.finalizers.giantswarm.io"
	CinderCsiTag         = "cinder.csi.openstack.org/cluster"

	LoadBalancerProvisioningStatusActive = "ACTIVE"
	LoadBalancerProvisioningStatusError  = "ERROR"
	VolumeStatusDeleting                 = "deleting"
)
