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

package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	bucketv1 "art-of-infrastructure-management/api/bucket/v1"
)

// S3BucketReconciler reconciles a S3Bucket object
type S3BucketReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	S3Client *s3.S3
}

var DefaultRequeueInterval = time.Second * 30

// createS3Bucket creates a new S3 bucket with the given bucket name
func createS3Bucket(svc *s3.S3, bucketName string) error {
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	log.Log.Info(colorCodeMessage(fmt.Sprintf("S3 bucket '%s' created\n", bucketName)))
	return nil
}

// GetBuckets retrieves the current number of S3 bucets in the S3BucketGroup
func (r *S3BucketReconciler) BucketExists(s3Bucket *bucketv1.S3Bucket) bool {
	result, err := r.S3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Log.Error(err, colorCodeMessage("error listing S3 buckets"))
		return false
	}
	for _, bucket := range result.Buckets {
		if *bucket.Name == s3Bucket.Name {
			return true
		}
	}
	return false
}

func colorCodeMessage(message string) string {
	colorYellow := "\033[33m"
	coloredMessage := colorYellow + message
	return coloredMessage
}

//+kubebuilder:rbac:groups=bucket.my.domain,resources=s3buckets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bucket.my.domain,resources=s3buckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bucket.my.domain,resources=s3buckets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3Bucket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *S3BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	s3Bucket := &bucketv1.S3Bucket{}
	err := r.Get(context.TODO(), req.NamespacedName, s3Bucket)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info(colorCodeMessage("Bucket was deleted...skipping reconcile"))
			return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
		}
		return reconcile.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	log.Log.Info(colorCodeMessage(fmt.Sprintf("Reconciling s3Bucket %s", s3Bucket.Name)))
	log.Log.Info(colorCodeMessage(fmt.Sprintf("Current Phase: %s, Desired Phase: %s", s3Bucket.Status.Phase, s3Bucket.Spec.Phase)))

	// If S3Bucket no longer exists, update status.Phase = "offline"
	if !r.BucketExists(s3Bucket) && s3Bucket.Status.Phase == bucketv1.PhaseOnline {
		s3Bucket.Status.Phase = bucketv1.PhaseOffline
		if err := r.Status().Update(ctx, s3Bucket); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update bucket status"))
			return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
		}
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
	}

	// If current phase = desired phase, skip reconcile
	if s3Bucket.Spec.Phase == s3Bucket.Status.Phase {
		log.Log.Info(colorCodeMessage("Skipping reconcile...Desired Phase == Current Phase"))
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
	}

	// If spec.Phase = "online" and status.Phase = "", this is a newly created bucket
	// Create a new s3 bucket and update status.Phase = "pending"
	if s3Bucket.Spec.Phase == bucketv1.PhaseOnline && s3Bucket.Status.Phase == "" {
		err := createS3Bucket(r.S3Client, s3Bucket.Name)
		if err != nil {
			log.Log.Error(err, colorCodeMessage("failed to create s3 bucket"))
			return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
		}
		s3Bucket.Status.Phase = bucketv1.PhasePending
		if err := r.Status().Update(ctx, s3Bucket); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update bucket status"))
			return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
		}
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
	}

	// If S3Bucket was pending online status and is now online, update status.Phase = "online"
	if s3Bucket.Status.Phase == bucketv1.PhasePending && r.BucketExists(s3Bucket) {
		s3Bucket.Status.Phase = bucketv1.PhaseOnline
		if err := r.Status().Update(ctx, s3Bucket); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update bucket status"))
			return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
		}
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
	}

	return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bucketv1.S3Bucket{}).
		Complete(r)
}
