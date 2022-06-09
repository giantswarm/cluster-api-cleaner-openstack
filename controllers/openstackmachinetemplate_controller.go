/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha5"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cluster-api-cleaner-openstack/pkg/key"
)

// OpenstackMachineTemplateReconciler reconciles a openstackCluster object
type OpenstackMachineTemplateReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplate,verbs=get;list;watch

func (r *OpenstackMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("openstackmachinetemplate", req.NamespacedName)
	log.V(1).Info("Reconciling")
	var template capo.OpenStackMachineTemplate
	err := r.Get(ctx, req.NamespacedName, &template)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, microerror.Mask(err)
	}

	// Handle deleted templates
	if !template.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, log, &template)
	}

	// Do nothing for non-deleted templates
	return ctrl.Result{}, nil
}

func (r *OpenstackMachineTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&capo.OpenStackMachineTemplate{}).
		Complete(r)
}

func (r *OpenstackMachineTemplateReconciler) reconcileDelete(ctx context.Context, log logr.Logger, template *capo.OpenStackMachineTemplate) (reconcile.Result, error) {
	if !controllerutil.ContainsFinalizer(template, key.CleanerFinalizerName) {
		return ctrl.Result{}, nil
	}

	exists, err := r.consumerMachineSetExists(ctx, log, template)
	if err != nil {
		return reconcile.Result{}, microerror.Mask(err)
	}

	if exists {
		log.V(1).Info("There are still some machinesets using this template.")
		return reconcile.Result{RequeueAfter: time.Minute * 5}, nil
	} else {
		log.Info("There is no machineset using this template. Removing finalizer.")
		controllerutil.RemoveFinalizer(template, key.CleanerFinalizerName)
		err = r.Update(ctx, template)
		return reconcile.Result{}, microerror.Mask(err)
	}
}

func (r *OpenstackMachineTemplateReconciler) consumerMachineSetExists(ctx context.Context, log logr.Logger, template *capo.OpenStackMachineTemplate) (bool, error) {
	clusterName, ok := template.Labels[key.CapiClusterLabelKey]
	if !ok {
		log.V(1).Info("Template don't have cluster name label",
			"expectedLabelKey", key.CapiClusterLabelKey,
			"existingLabels", template.Labels)
		return false, microerror.Maskf(invalidObjectError, "template don't have cluster name label")
	}

	var machineSetList capi.MachineSetList
	opts := client.ListOptions{Namespace: template.Namespace,
		LabelSelector: labels.SelectorFromSet(map[string]string{key.CapiClusterLabelKey: clusterName})}
	err := r.Client.List(ctx, &machineSetList, &opts)
	if err != nil {
		return false, microerror.Mask(err)
	}

	for _, machineSet := range machineSetList.Items {
		infraRef := machineSet.Spec.Template.Spec.InfrastructureRef
		if infraRef.Name == template.Name && infraRef.Kind == template.Kind {
			log.V(1).Info("There is a machineset using the template", "machineset", machineSet.Name)
			return true, nil
		}
	}

	return false, nil
}
