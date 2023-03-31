package cleaner

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha6"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/provider"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cluster-api-cleaner-openstack/pkg/key"
)

type LoadBalancerCleaner struct {
	cli client.Client
}

func NewLoadBalancerCleaner(cli client.Client) *LoadBalancerCleaner {
	return &LoadBalancerCleaner{cli: cli}
}

// force implementing Cleaner interface
var _ Cleaner = &LoadBalancerCleaner{}

func (lbc *LoadBalancerCleaner) Clean(ctx context.Context, log logr.Logger, oc *capo.OpenStackCluster, clusterTag string) (bool, error) {
	log = log.WithName("LoadBalancerCleaner")

	providerClient, opts, projectID, err := provider.NewClientFromCluster(ctx, lbc.cli, oc)
	if err != nil {
		return true, microerror.Mask(err)
	}

	scope := &scope.Scope{
		ProviderClient:     providerClient,
		ProviderClientOpts: opts,
		ProjectID:          projectID,
		Logger:             log,
	}

	loadbalancerClient, err := openstack.NewLoadBalancerV2(providerClient, gophercloud.EndpointOpts{
		Region: opts.RegionName,
	})
	if err != nil {
		return true, microerror.Mask(err)
	}

	allPages, err := loadbalancers.List(loadbalancerClient, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return true, microerror.Mask(err)
	}

	lbList, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		return true, microerror.Mask(err)
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return true, microerror.Mask(err)
	}

	requeue := false
	for _, lb := range lbList {
		if !mustBeDeleted(lb, clusterTag) {
			continue
		}

		log.Info("Cleaning load balancer", "id", lb.ID, "status", lb.ProvisioningStatus, "project", projectID)

		if !isOkForDeletion(lb.ProvisioningStatus) {
			log.V(1).Info("Will requeue openstackcluster because of the loadbalancer",
				"id", lb.ID, "provisioningStatus", lb.ProvisioningStatus)
			requeue = true
			continue
		}

		deleteOpts := loadbalancers.DeleteOpts{
			Cascade: true,
		}

		if lb.VipPortID != "" {
			fip, err := networkingService.GetFloatingIPByPortID(lb.VipPortID)
			if err != nil {
				return true, microerror.Mask(err)
			}

			if fip != nil && fip.FloatingIP != "" {
				log.Info("Cleaning floating IP", "ip", fip.FloatingIP, "loadbalancer", lb.ID, "project", projectID)
				err = lbc.cleanFloatingIP(networkingService, oc, fip)
				if err != nil {
					return true, microerror.Mask(err)
				}
			}
		}

		err = loadbalancers.Delete(loadbalancerClient, lb.ID, deleteOpts).ExtractErr()
		if err != nil {
			return true, microerror.Mask(err)
		} else {
			requeue = true
		}
	}

	log.V(1).Info("", "Requeue", requeue)
	if requeue {
		return true, nil
	} else {
		return false, nil
	}
}

// isOkForDeletion allows deletion only when the status is ACTIVE or ERROR
// See https://docs.openstack.org/api-ref/load-balancer/v2/index.html#prov-status
func isOkForDeletion(status string) bool {
	return status == key.LoadBalancerProvisioningStatusActive ||
		status == key.LoadBalancerProvisioningStatusError
}

func (lbc *LoadBalancerCleaner) cleanFloatingIP(ns *networking.Service, oc *capo.OpenStackCluster, fip *floatingips.FloatingIP) error {
	if err := ns.DisassociateFloatingIP(oc, fip.FloatingIP); err != nil {
		return microerror.Mask(err)
	}
	if err := ns.DeleteFloatingIP(oc, fip.FloatingIP); err != nil {
		return microerror.Mask(err)
	}
	return nil
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
