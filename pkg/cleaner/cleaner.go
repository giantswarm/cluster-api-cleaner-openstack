package cleaner

import (
	"github.com/go-logr/logr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Cleaner interface {
	Clean(cli client.Client, log logr.Logger, oc *capo.OpenStackCluster) error
}
