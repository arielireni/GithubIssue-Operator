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

// GitHubIssueSpec defines the desired state of GitHubIssue
type GitHubIssueSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Repo represents a clients repo url
	/* Validation in the CRD level - an attempt to create a CRD with malformed 'repo' will fail */
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9\_.-]+/[a-zA-Z0-9\_.-]+$
	Repo string `json:"repo"`

	// Title represents the title of the issue
	Title string `json:"title"`

	// Body represents the description of the issue
	Description string `json:"description,omitempty"`
}

// GitHubIssueStatus defines the observed state of GitHubIssue
type GitHubIssueStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// State represents the state of the real clients issue
	State string `json:"state,omitempty"`

	// LastUpdateTimestamp represents a timestamp of the last time the state was updated
	LastUpdateTimestamp string `json:"updated_at,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitHubIssue is the Schema for the githubissues API
type GitHubIssue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubIssueSpec   `json:"spec,omitempty"`
	Status GitHubIssueStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GitHubIssueList contains a list of GitHubIssue
type GitHubIssueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHubIssue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHubIssue{}, &GitHubIssueList{})
}
