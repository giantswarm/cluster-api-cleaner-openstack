package cleaner

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/provider"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LoadBalancerCleaner struct {
}

func (lbc *LoadBalancerCleaner) Clean(cli client.Client, log logr.Logger, oc *capo.OpenStackCluster) error {
	clusterTag := getClusterTag(oc)
	if clusterTag == "" {
		return nil
	}

	providerClient, opts, err := provider.NewClientFromCluster(context.TODO(), cli, oc)
	if err != nil {
		return err
	}

	loadbalancerClient, err := openstack.NewLoadBalancerV2(providerClient, gophercloud.EndpointOpts{
		Region: opts.RegionName,
	})
	if err != nil {
		return err
	}

	allPages, err := loadbalancers.List(loadbalancerClient, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return err
	}

	lbList, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		return err
	}

	service, err := networking.NewService(providerClient, opts, log)
	if err != nil {
		return err
	}

	for _, lb := range lbList {

		if !mustBeDeleted(lb, clusterTag) {
			continue
		}

		deleteOpts := loadbalancers.DeleteOpts{
			Cascade: true,
		}

		if lb.VipPortID != "" {
			fip, err := service.GetFloatingIPByPortID(lb.VipPortID)
			if err != nil {
				return err
			}

			if fip != nil && fip.FloatingIP != "" {
				if err = service.DisassociateFloatingIP(oc, fip.FloatingIP); err != nil {
					return err
				}
				if err = service.DeleteFloatingIP(oc, fip.FloatingIP); err != nil {
					return err
				}
			}
		}

		log.Info("Deleting load balancer", "id", lb.ID)
		err = loadbalancers.Delete(loadbalancerClient, lb.ID, deleteOpts).ExtractErr()
		if err != nil {
			return err
		}
	}

	return nil
}

func getClusterTag(oc *capo.OpenStackCluster) string {
	if len(oc.Spec.Tags) > 0 {
		return oc.Spec.Tags[0]
	}
	return ""
}

func mustBeDeleted(lb loadbalancers.LoadBalancer, clusterTag string) bool {
	prefix := "kube_service_" + clusterTag
	for _, tag := range lb.Tags {
		if strings.HasPrefix(tag, prefix) {
			return true
		}
	}
	return false
}
