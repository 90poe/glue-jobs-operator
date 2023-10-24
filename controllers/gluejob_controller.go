/*
Copyright 2023.

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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	awsv1alpha1 "github.com/90poe/glue-jobs-operator/api/v1alpha1"
	"github.com/90poe/glue-jobs-operator/internal/config"
	"github.com/90poe/glue-jobs-operator/internal/consts"
	"github.com/90poe/glue-jobs-operator/internal/glue"
	"github.com/go-logr/logr"
)

const glueJobFinalizer = "gluejobs.aws.90poe.io/finalizer"

// GlueJobReconciler reconciles a GlueJob object
type GlueJobReconciler struct {
	ctx    context.Context
	config config.OperatorConfig
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aws.90poe.io,resources=gluejobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.90poe.io,resources=gluejobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.90poe.io,resources=gluejobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GlueJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	r.ctx = ctx
	reqLogger := log.FromContext(ctx).WithValues("gluejobs", req.NamespacedName)

	// Fetch the GlueJob K8S object instance
	glueJob := &awsv1alpha1.GlueJob{}
	err := r.Get(ctx, req.NamespacedName, glueJob)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.V(1).Info("GlueJob resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.V(0).Error(err, "Failed to get GlueJob.")
		return ctrl.Result{}, err
	}

	// 1. check if job exists on AWS
	// 2. if exists, update job up to date
	// 3. if not exists, create job on AWS

	// 1. check if job exists on AWS
	awsGlueJob, err := glue.NewJob(ctx, glueJob.Spec)
	if err != nil {
		return r.setLatestError(glueJob, err, "CreateNewGlueJobFailed")
	}

	// Check if the GlueJob instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isGlueJobdMarkedToBeDeleted := glueJob.GetDeletionTimestamp() != nil
	if isGlueJobdMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(glueJob, glueJobFinalizer) {
			// Run finalization logic for GlueJobFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeGlueJob(reqLogger, glueJob, awsGlueJob); err != nil {
				return ctrl.Result{
					// requeue after 5 seconds
					RequeueAfter: 5 * time.Second,
				}, err
			}

			// Remove GlueJobFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			err = r.addOrRemoveFinalizer(glueJob, false)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	message := "Successfully updated GlueJob"
	if awsGlueJob.JobExists() {
		// 2. if exists, update job up to date
		err = r.updateJob(awsGlueJob, reqLogger)
	} else {
		// 3. if not exists, create job on AWS
		err = r.createJob(awsGlueJob, reqLogger)
		message = "Successfully created GlueJob"
	}
	if err != nil {
		return r.setLatestError(glueJob, err, "GlueJobFailed")
	}

	// Add finalizer for this CR
	err = r.addOrRemoveFinalizer(glueJob, true)
	if err != nil {
		return r.setLatestError(glueJob, err, "GlueJobFailed")
	}

	return r.succReconcileRet(glueJob, reqLogger, message)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlueJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	var err error
	// get config from Env
	r.config, err = config.New()
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1alpha1.GlueJob{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: r.config.MaxConcurrentReconciles}).
		WithEventFilter(ignoreUpdateDeletePredicate()).
		Complete(r)
}

// ignoreUpdateDeletePredicater is brilliantly useful function, it will prevent multiple reconcile calls
func ignoreUpdateDeletePredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// addFinalizer will add finalizer to CR
func (r *GlueJobReconciler) addOrRemoveFinalizer(glueJob *awsv1alpha1.GlueJob, add bool) error {
	if add && controllerutil.ContainsFinalizer(glueJob, glueJobFinalizer) {
		// we are adding again and it already exists
		return nil
	}
	if add {
		// add finalizer
		controllerutil.AddFinalizer(glueJob, glueJobFinalizer)
	} else {
		// remove finalizer
		controllerutil.RemoveFinalizer(glueJob, glueJobFinalizer)
	}
	// update resource
	return r.Update(r.ctx, glueJob)
}

// setLatestError will set latest error on condition
func (r *GlueJobReconciler) setLatestError(
	gj *awsv1alpha1.GlueJob,
	err error,
	errType string,
) (reconcile.Result, error) {
	reterr := err
	condition := metav1.Condition{
		Type:               consts.StatusNotReady,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Status:             metav1.ConditionTrue,
		Reason:             errType,
		Message:            fmt.Sprintf("%v", err),
	}
	gj.Status.Conditions = append(gj.Status.Conditions, condition)
	err = r.Status().Update(r.ctx, gj)
	if err != nil {
		reterr = kerrors.NewAggregate([]error{reterr, err})
	}
	return ctrl.Result{}, reterr
}

func (r *GlueJobReconciler) createJob(awsGJ *glue.Job, reqLogger logr.Logger) error {
	reqLogger.V(1).Info("Create GlueJob")
	err := awsGJ.CreateJob()
	if err != nil {
		return err
	}
	return nil
}

func (r *GlueJobReconciler) updateJob(awsGJ *glue.Job, reqLogger logr.Logger) error {
	reqLogger.V(1).Info("Update GlueJob")
	err := awsGJ.UpateJob()
	if err != nil {
		return err
	}
	return nil
}

// Function would always return reconcile with requeue and time to requeue
func (r *GlueJobReconciler) succReconcileRet(gj *awsv1alpha1.GlueJob,
	reqLogger logr.Logger, message string) (reconcile.Result, error) {
	condition := metav1.Condition{
		Type:               consts.StatusReady,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Status:             metav1.ConditionTrue,
		Reason:             consts.SuccessReconcile,
		Message:            message,
	}
	gj.Status.Conditions = append(gj.Status.Conditions, condition)
	err := r.Status().Update(r.ctx, gj)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// finalizeGlueJob is part of finalizers logic and deletes the GlueJob on AWS
func (r *GlueJobReconciler) finalizeGlueJob(reqLogger logr.Logger, a *awsv1alpha1.GlueJob, awsGJ *glue.Job) error {
	// Delete the GlueJob instance
	reqLogger.V(0).Info("Deleting GlueJob on AWS")
	err := awsGJ.DeleteJob()
	if err != nil {
		return err
	}
	return nil
}
