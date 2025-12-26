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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EC2InstanceSpec defines the desired state of EC2Instance
type EC2InstanceSpec struct {
	Region           string            `json:"region"`
	AmiID            string            `json:"amiId"`
	InstanceType     string            `json:"instanceType"`
	SubnetID         string            `json:"subnetId,omitempty"`
	SecurityGroupIDs []string          `json:"securityGroupIds,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type EC2InstanceStatus struct {
	InstanceID string `json:"instanceId,omitempty"`
	State      string `json:"state,omitempty"`
	PublicIP   string `json:"publicIp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EC2Instance is the Schema for the ec2instances API
type EC2Instance struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of EC2Instance
	// +required
	Spec EC2InstanceSpec `json:"spec"`

	// status defines the observed state of EC2Instance
	// +optional
	Status EC2InstanceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// EC2InstanceList contains a list of EC2Instance
type EC2InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []EC2Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EC2Instance{}, &EC2InstanceList{})
}
