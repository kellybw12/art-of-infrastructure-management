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

package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	bucketv1 "art-of-infrastructure-management/api/bucket/v1"
	bucketgroupv1 "art-of-infrastructure-management/api/v1"
)

var bucket_id = 1
var DefaultRequeueInterval = time.Second * 10

// S3BucketGroupReconciler reconciles a S3BucketGroup object
type S3BucketGroupReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	S3Client *s3.S3
}

func colorCodeMessage(message string) string {
	colorCyan := "\033[36m"
	coloredMessage := colorCyan + message
	return coloredMessage
}

//+kubebuilder:rbac:groups=bucketgroup.my.domain,resources=s3bucketgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bucketgroup.my.domain,resources=s3bucketgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bucketgroup.my.domain,resources=s3bucketgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3BucketGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *S3BucketGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Uncomment to run part 2
	// result, err := DoPart2(r, ctx, req)
	// if err != nil {
	// 	log.Log.Error(err, colorCodeMessage("error occurred when running part 2"))
	// }

	// Uncomment to run part 3
	result, err := DoPart3(r, ctx, req)
	if err != nil {
		log.Log.Error(err, colorCodeMessage("error occurred when running part 3"))
	}

	return result, nil
}

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

// generateNewBucketName creates the new bucket name based on the current bucket_id
func generateNewBucketName() string {
	bucketName := "my-s3-bucket-" + strconv.Itoa(bucket_id)
	bucket_id += 1
	return bucketName
}

// GetBuckets retrieves the current number of S3 bucets in the S3BucketGroup
func (r *S3BucketGroupReconciler) GetBuckets() (*s3.ListBucketsOutput, error) {
	buckets, err := r.S3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Log.Error(err, colorCodeMessage("error listing S3 buckets"))
		return nil, err
	}
	return buckets, nil
}

// DoPart2 holds the logic to demonstrate the demo for Part 2
func DoPart2(r *S3BucketGroupReconciler, ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Retrieves the current number of buckets
	result, err := r.GetBuckets()
	if err != nil {
		log.Log.Error(err, colorCodeMessage("error while retrieving s3 buckets"))
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	// Retrieve the current state of the S3BucketGroup
	s3BucketGroup := &bucketgroupv1.S3BucketGroup{}
	err = r.Get(context.TODO(), req.NamespacedName, s3BucketGroup)
	if err != nil {
		log.Log.Error(err, colorCodeMessage("failed to retrieve current state of bucket group"))
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	// Update the S3BucketGroup status if the status differs from the actual state.
	// Force reconcile if status was updated.
	s3BucketGroup.Status.BucketCount = len(result.Buckets)
	if s3BucketGroup.Status.BucketCount != len(result.Buckets) {
		s3BucketGroup.Status.BucketCount = len(result.Buckets)
		if err := r.Status().Update(ctx, s3BucketGroup); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update status of bucket group"))
		}
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	log.Log.Info(
		colorCodeMessage(fmt.Sprintf("Current S3 bucket count: %d, Desired S3 Bucket count: %d",
			s3BucketGroup.Status.BucketCount, s3BucketGroup.Spec.DesiredBucketCount)),
	)

	// Create new S3 buckets if the current S3BucketGroup count < desired S3BucketGroup count
	if s3BucketGroup.Status.BucketCount < s3BucketGroup.Spec.DesiredBucketCount {
		deficit := s3BucketGroup.Spec.DesiredBucketCount - len(result.Buckets)
		for i := 0; i < deficit; i++ {
			bucketName := generateNewBucketName()
			err := createS3Bucket(r.S3Client, bucketName)
			if err != nil {
				log.Log.Error(err, colorCodeMessage("failed to create s3 bucket"))
			} else {
				s3BucketGroup.Status.BucketCount += 1
			}
			// Add sleep for easier traceability
			time.Sleep(time.Second * 10)
		}
		if err := r.Status().Update(ctx, s3BucketGroup); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update s3BucketGroup status"))
		}
	} else {
		log.Log.Info(colorCodeMessage("No creations needed. Desired S3 bucket count == Current S3 bucket count"))
	}
	return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
}

func createS3BucketCRD(r *S3BucketGroupReconciler, ctx context.Context, req ctrl.Request, bucketName string, bucketGroup *bucketgroupv1.S3BucketGroup) (*bucketv1.S3Bucket, error) {
	bucket := &bucketv1.S3Bucket{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bucketName,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"bucketGroupName": bucketGroup.Name,
			},
		},
		Spec: bucketv1.S3BucketSpec{
			Phase: bucketv1.PhaseOnline,
		},
	}
	err := r.Client.Create(ctx, bucket)
	if err != nil {
		return bucket, err
	}
	log.Log.Info(colorCodeMessage(fmt.Sprintf("S3 bucket '%s' created\n", bucket.Name)))
	return bucket, nil
}

// End of code for Part 2: Simple Example

// Start code for Part 3: Complex Example

// listBucketsInBucketGroup returns a list of buckets that belong to a specific S3bucketGroup
func listBucketsInBucketGroup(r *S3BucketGroupReconciler, ctx context.Context, s3BucketGroup *bucketgroupv1.S3BucketGroup) ([]bucketv1.S3Bucket, error) {
	var buckets bucketv1.S3BucketList
	selector := labels.NewSelector()
	bgReq, _ := labels.NewRequirement("bucketGroupName", selection.Equals, []string{s3BucketGroup.Name})
	selector = selector.Add(*bgReq)

	listOptions := &client.ListOptions{
		LabelSelector: selector,
	}

	if err := r.Client.List(ctx, &buckets, client.InNamespace(s3BucketGroup.Namespace), listOptions); err != nil {
		return buckets.Items, err
	}
	return buckets.Items, nil
}

// clearOfflineBuckets deletes buckets that are completely offline (spec.Phase = "online" && status.Phase = "offline")
func clearOfflineBuckets(r *S3BucketGroupReconciler, ctx context.Context, buckets []bucketv1.S3Bucket) (bool, error) {
	isBucketsCleared := false
	for _, bucket := range buckets {
		if bucket.Spec.Phase == bucketv1.PhaseOnline && bucket.Status.Phase == bucketv1.PhaseOffline {
			log.Log.Info(colorCodeMessage(fmt.Sprintf("Deleting bucket %s", bucket.Name)))
			if err := r.Client.Delete(ctx, &bucket); err != nil {
				log.Log.Error(err, colorCodeMessage(fmt.Sprintf("failed to delete offline bucket %s", bucket.Name)))
			} else {
				isBucketsCleared = true
			}
		}
	}
	return isBucketsCleared, nil
}

func DoPart3(r *S3BucketGroupReconciler, ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Retrieve the current state of the S3BucketGroup
	s3BucketGroup := &bucketgroupv1.S3BucketGroup{}
	err := r.Get(context.TODO(), req.NamespacedName, s3BucketGroup)
	if err != nil {
		log.Log.Error(err, colorCodeMessage("failed to retrieve current state of s3BucketGroup"))
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	// Retrieve the list of buckets in the S3BucketGroup
	bucketsInBG, err := listBucketsInBucketGroup(r, ctx, s3BucketGroup)
	if err != nil {
		log.Log.Error(err, colorCodeMessage("failed to retrieve buckets in bucket group"), "Bucket Group", s3BucketGroup.Name)
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	// Clear buckets that are in an unrecoverable state (spec.Phase = "online" && status.Phase = "offline")
	isBucketsDeleted, err := clearOfflineBuckets(r, ctx, bucketsInBG)
	// If buckets were cleared, force reconcile to retrieve updated list of buckets
	if isBucketsDeleted {
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil
	}
	if err != nil {
		log.Log.Error(err, colorCodeMessage("failed to clear offline buckets in bucket group"), "Bucket Group", s3BucketGroup.Name)
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	// Update the S3BucketGroup status if the status differs from the actual state.
	// Force reconcile if status was updated.
	s3BucketGroup.Status.BucketCount = len(bucketsInBG)
	if s3BucketGroup.Status.BucketCount != len(bucketsInBG) {
		s3BucketGroup.Status.BucketCount = len(bucketsInBG)
		if err := r.Status().Update(ctx, s3BucketGroup); err != nil {
			log.Log.Error(err, colorCodeMessage("failed to update s3BucketGroup status"))
		}
		return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, err
	}

	log.Log.Info(
		colorCodeMessage(fmt.Sprintf("Current S3 bucket count: %d, Desired S3 Bucket count: %d",
			s3BucketGroup.Status.BucketCount, s3BucketGroup.Spec.DesiredBucketCount)),
	)

	if s3BucketGroup.Status.BucketCount < s3BucketGroup.Spec.DesiredBucketCount {
		deficit := s3BucketGroup.Spec.DesiredBucketCount - len(bucketsInBG)
		for i := 0; i < deficit; i++ {
			bucketName := generateNewBucketName()
			_, err = createS3BucketCRD(r, ctx, req, bucketName, s3BucketGroup)
			if err != nil {
				log.Log.Error(err, colorCodeMessage("failed to create s3 bucket"), "Bucket Group", s3BucketGroup.Name)
			}
			// Add sleep for better traceability
			time.Sleep(time.Second * 5)

		}
	} else {
		log.Log.Info(colorCodeMessage("No creations needed. Desired S3 bucket count == Current S3 bucket count"))
	}

	return ctrl.Result{RequeueAfter: DefaultRequeueInterval}, nil

}

// ignoreDeletionPredicate prevents the MachineGroupReconciler from listening to deletion events on desired Kinds
func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Object.GetObjectKind().GroupVersionKind().Kind != "S3Bucket"
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BucketGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bucketgroupv1.S3BucketGroup{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}
