package cleaner

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/provider"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VolumeCleaner struct {
}

const cinderCsiTag = "cinder.csi.openstack.org/cluster"

func (vc *VolumeCleaner) Clean(cli client.Client, log logr.Logger, oc *capo.OpenStackCluster) error {
	clusterTag := getClusterTag(oc)
	if clusterTag == "" {
		return nil
	}

	providerClient, opts, err := provider.NewClientFromCluster(context.TODO(), cli, oc)
	if err != nil {
		return err
	}

	volumeClient, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region: opts.RegionName,
	})
	if err != nil {
		return err
	}

	allPages, err := volumes.List(volumeClient, volumes.ListOpts{Metadata: map[string]string{cinderCsiTag: clusterTag}}).AllPages()
	if err != nil {
		return err
	}

	volumeList, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return err
	}

	deleteOpts := volumes.DeleteOpts{
		Cascade: true,
	}

	for _, volume := range volumeList {
		log.Info("Deleting volume", "id", volume.ID)
		err = volumes.Delete(volumeClient, volume.ID, deleteOpts).ExtractErr()
		if err != nil {
			return err
		}
	}

	return nil
}
