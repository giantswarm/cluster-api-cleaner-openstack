package cleaner

import (
	"context"

	"github.com/go-logr/logr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha4"
)

type Cleaner interface {
	Clean(ctx context.Context, log logr.Logger, oc *capo.OpenStackCluster, clusterTag string) error
}
