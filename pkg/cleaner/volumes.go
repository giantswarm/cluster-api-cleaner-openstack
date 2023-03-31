package cleaner

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha6"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/provider"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cluster-api-cleaner-openstack/pkg/key"
)

type VolumeCleaner struct {
	cli client.Client
}

func NewVolumeCleaner(cli client.Client) *VolumeCleaner {
	return &VolumeCleaner{cli: cli}
}

// force implementing Cleaner interface
var _ Cleaner = &VolumeCleaner{}

func (vc *VolumeCleaner) Clean(ctx context.Context, log logr.Logger, oc *capo.OpenStackCluster, clusterTag string) (bool, error) {
	log = log.WithName("VolumeCleaner")

	providerClient, opts, projectID, err := provider.NewClientFromCluster(ctx, vc.cli, oc)
	if err != nil {
		return true, microerror.Mask(err)
	}

	volumeClient, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region: opts.RegionName,
	})
	if err != nil {
		return true, microerror.Mask(err)
	}

	listOpts := volumes.ListOpts{Metadata: map[string]string{key.CinderCsiTag: clusterTag}}
	allPages, err := volumes.List(volumeClient, listOpts).AllPages()
	if err != nil {
		return true, microerror.Mask(err)
	}

	volumeList, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return true, microerror.Mask(err)
	}

	deleteOpts := volumes.DeleteOpts{
		Cascade: true,
	}

	requeue := false
	for _, volume := range volumeList {
		log.Info("Cleaning volume.", "id", volume.ID, "status", volume.Status, "project", projectID)
		if volume.Status == key.VolumeStatusDeleting {
			log.V(1).Info("Volume deletion already triggered", "id", volume.ID)
			requeue = true
			continue
		}

		err = volumes.Delete(volumeClient, volume.ID, deleteOpts).ExtractErr()
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
