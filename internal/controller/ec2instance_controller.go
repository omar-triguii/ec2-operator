/*
Copyright 2025.

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

package controller

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrav1alpha1 "github.com/omar--triguii/ec2-operator/api/v1alpha1"
)

// EC2InstanceReconciler reconciles a EC2Instance object
type EC2InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infra.trigui.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.trigui.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infra.trigui.com,resources=ec2instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EC2Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *EC2InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var cr infrav1alpha1.EC2Instance
	const ec2Finalizer = "infra.trigui.com/ec2instance-finalizer"

	if err := r.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// ✅ ADD THIS BLOCK RIGHT HERE (after Get)
	if !cr.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(&cr, ec2Finalizer) {
			if cr.Status.InstanceID != "" {
				cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cr.Spec.Region))
				if err != nil {
					return ctrl.Result{}, err
				}
				ec2c := ec2.NewFromConfig(cfg)

				_, err = ec2c.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
					InstanceIds: []string{cr.Status.InstanceID},
				})
				if err != nil {
					return ctrl.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(&cr, ec2Finalizer)
			if err := r.Update(ctx, &cr); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	// ✅ ADD THIS BLOCK RIGHT HERE (before creation logic)
	if !controllerutil.ContainsFinalizer(&cr, ec2Finalizer) {
		controllerutil.AddFinalizer(&cr, ec2Finalizer)
		if err := r.Update(ctx, &cr); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Creation-only: if already created, do nothing
	if cr.Status.InstanceID != "" {
		return ctrl.Result{}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cr.Spec.Region))
	if err != nil {
		return ctrl.Result{}, err
	}
	ec2c := ec2.NewFromConfig(cfg)

	out, err := ec2c.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      aws.String(cr.Spec.AmiID),
		InstanceType: types.InstanceType(cr.Spec.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	cr.Status.InstanceID = aws.ToString(out.Instances[0].InstanceId)
	if err := r.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EC2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha1.EC2Instance{}).
		Named("ec2instance").
		Complete(r)
}
