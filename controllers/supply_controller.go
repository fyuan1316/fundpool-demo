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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fundv1alpha1 "github.com/fyuan1316/fundpool/api/v1alpha1"
)

// SupplyReconciler reconciles a Supply object
type SupplyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=fund.demo.com,resources=supplies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fund.demo.com,resources=supplies/status,verbs=get;update;patch

func (r *SupplyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("supply", req.NamespacedName)

	// your logic here
	supply := fundv1alpha1.Supply{}
	if err := r.Get(context.Background(), req.NamespacedName, &supply); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if supply.Status.IsReachFinal() {
		return ctrl.Result{}, nil
	}
	var fundPools fundv1alpha1.FundPoolList
	if err := r.List(context.Background(), &fundPools); err != nil {
		r.Log.Error(err, fmt.Sprintf("list fundv1alpha1.FundPoolList error: %v", req.NamespacedName))
		return ctrl.Result{}, err
	}
	supplyCopy := supply.DeepCopy()
	if len(fundPools.Items) == 0 {
		supplyCopy.Status.MarkFailed()
		err := r.Update(context.Background(), supplyCopy)
		if err != nil {
			// TODO: process update failed.
			r.Log.Error(err, fmt.Sprintf("update Supply error: %v", req.NamespacedName))
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	resultPools := Calc(supply.Spec.Request, fundPools.Items)
	var allocations []fundv1alpha1.Allocation
	for i := range resultPools {
		resultPool := resultPools[i]
		allocation := fundv1alpha1.Allocation{}
		allocation.Pool = resultPool.Name
		allocation.Shortfalls = resultPool.GetShortfalls()
		allocations = append(allocations, allocation)
	}
	supplyCopy.Status.Allocations = allocations
	supplyCopy.Status.MarkSuccessed()

	if err := r.Update(context.Background(), supplyCopy); err != nil {
		// TODO: process update failed.
		r.Log.Error(err, fmt.Sprintf("update Supply error: %v", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SupplyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fundv1alpha1.Supply{}).
		Complete(r)
}
