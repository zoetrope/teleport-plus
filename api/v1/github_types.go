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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitHubSpec defines the desired state of GitHub
type GitHubSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ClientID is the Github OAuth app client ID
	ClientID string `json:"client_id"`
	// ClientSecret is the Github OAuth app client secret
	ClientSecret string `json:"client_secret"`
	// RedirectURL is the authorization callback URL
	RedirectURL string `json:"redirect_url"`
	// TeamsToLogins maps Github team memberships onto allowed logins/roles
	TeamsToLogins []TeamMapping `json:"teams_to_logins"`
	// Display is the connector display name
	Display string `json:"display"`
}

// TeamMapping represents a single team membership mapping
type TeamMapping struct {
	// Organization is a Github organization a user belongs to
	Organization string `json:"organization"`
	// Team is a team within the organization a user belongs to
	Team string `json:"team"`
	// Logins is a list of allowed logins for this org/team
	// +optional
	Logins []string `json:"logins,omitempty"`
	// KubeGroups is a list of allowed kubernetes groups for this org/team
	// +optional
	KubeGroups []string `json:"kubernetes_groups,omitempty"`
}

// GitHubStatus defines the observed state of GitHub
type GitHubStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Registered bool `json:"registered"`
}

// +kubebuilder:object:root=true

// GitHub is the Schema for the githubs API
type GitHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubSpec   `json:"spec,omitempty"`
	Status GitHubStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GitHubList contains a list of GitHub
type GitHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHub{}, &GitHubList{})
}
