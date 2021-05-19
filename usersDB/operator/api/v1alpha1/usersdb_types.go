/*
Copyright 2021.

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

// UsersDBSpec defines the desired state of UsersDB
type UsersDBSpec struct {
	//+kubebuilder:validation:Minimum=0
	// Size is the size of the usersDB deployment
	Size int32 `json:"size"`
}

// UsersDBStatus defines the observed state of UsersDB
type UsersDBStatus struct {
	// Nodes are the names of the users-db pods
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UsersDB is the Schema for the usersdbs API
//+kubebuilder:subresource:status
type UsersDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UsersDBSpec   `json:"spec,omitempty"`
	Status UsersDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UsersDBList contains a list of UsersDB
type UsersDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UsersDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UsersDB{}, &UsersDBList{})
}
