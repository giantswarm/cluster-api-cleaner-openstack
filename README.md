# cluster-api-cleaner-openstack

An helper operator for CAPO to delete resources by workload clusters

## How does it work?

- It observes `OpenStackCluster` objects.
- It doesn't do anything in `reconcileNormal` other than adding finalizer.
- It respects `cluster.x-k8s.io/cluster-name` label in `OpenStackCluster` objects to get the actual cluster names.
- `clusterTag` is built as `giant_swarm_cluster_<management-cluster-name>_<workload_cluster-name>`.
- When an `OpenStackCluster` is deleted, it
  - cleans volumes ( whose metadata contains `cinder.csi.openstack.org/cluster: <clusterTag>` ) created by Cinder CSI 
  - cleans loadbalancers ( whose tags contain `kube_service_<clusterTag>.*` ) created by 
    openstack-cloud-controller-manager 