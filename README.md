# cluster-api-cleaner-openstack

A helper operator for CAPO to delete resources created by apps in workload clusters.

## `OpenstackClusterController`

### Why?

`openstack-cloud-controller-manager` in workload cluster creates LoadBalancers in OpenStack for Services in the cluster. `openstack-cinder-csi` also creates some Volumes in Openstack for PersistentVolumes in the cluster. When the worklaod cluster is deleted, `cluster-api-provider-openstack` doesn't clean these resources. This controller helps for clean-up of workload clusters.

###  How does it work?

- It observes `OpenStackCluster` objects.
- It doesn't do anything in `reconcileNormal` other than adding finalizer.
- It respects `cluster.x-k8s.io/cluster-name` label in `OpenStackCluster` objects to get the actual cluster names.
- `clusterTag` is built as `giant_swarm_cluster_<management-cluster-name>_<workload_cluster-name>`.
- When an `OpenStackCluster` is deleted, it
  - cleans volumes ( whose metadata contains `cinder.csi.openstack.org/cluster: <clusterTag>` ) created by Cinder CSI 
  - cleans loadbalancers ( whose tags contain `kube_service_<clusterTag>.*` ) created by 
    openstack-cloud-controller-manager 

## `OpenstackMachineTemplateController`

### Why?

[cluster-openstack](https://github.com/giantswarm/cluster-openstack)  doesn't use `ClusterClass` and it creates new `OpenStackMachineTemplate` during upgrades if necessary. The old templates are necessary for CAPI controllers to complete machine updates so we cannot delete them during helm upgrades. We are adding a finalizers into templates to prevent their deletion by helm. This controller helps us to clean unused templates after upgrades.

###  How does it work?

- It observes `OpenstackMachineTemplate` objects.
- It doesn't do anything in `reconcileNormal`. Finalizers are added by [cluster-openstack](https://github.com/giantswarm/cluster-openstack) app during creation.
- When a OpenstackMachineTemplate is deleted, the controller checks all machinesets of the cluster and whether there is a machineset that consumes the OpenstackMachineTemplate or not. If there is no consumer, the controller deletes the OpenstackMachineTemplate.

> Note that machineset objects are not deleted by machinedeploymentcontroller just after rollout. revisionHistoryLimit is 1 by default. It means an unused machineset is deleted by machinedeploymentcontroller
after the next upgrade and it triggers deletion of OpenStackMachineTemplate deletion.
