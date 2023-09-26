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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// S3BucketGroupSpec defines the desired state of S3BucketGroup
type S3BucketGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DesiredBucketCount int `json:"desiredBucketCount,omitempty"`
}

// S3BucketGroupStatus defines the observed state of S3BucketGroup
type S3BucketGroupStatus struct {
	BucketCount int `json:"bucketCount,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// S3BucketGroup is the Schema for the s3bucketgroups API
type S3BucketGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   S3BucketGroupSpec   `json:"spec,omitempty"`
	Status S3BucketGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// S3BucketGroupList contains a list of S3BucketGroup
type S3BucketGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []S3BucketGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&S3BucketGroup{}, &S3BucketGroupList{})
}
