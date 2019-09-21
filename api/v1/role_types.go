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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RoleSpec defines the desired state of Role
type RoleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Options RoleOptions    `json:"options,omitempty"`
	// +optional
	Allow   RoleConditions `json:"allow,omitempty"`
	// +optional
	Deny    RoleConditions `json:"deny,omitempty"`
}

type RoleOptions struct {
	ForwardAgent          bool           `json:"forward_agent"`
	// +optional
	MaxSessionTTL         *time.Duration `json:"max_session_ttl,omitempty"`
	// +optional
	PortForwarding        *bool          `json:"port_forwarding,omitempty"`
	CertificateFormat     string         `json:"cert_format"`
	// +optional
	ClientIdleTimeout     int64          `json:"client_idle_timeout,omitempty"`
	// +optional
	DisconnectExpiredCert bool           `json:"disconnect_expired_cert,omitempty"`
}

type RoleConditions struct {
	Logins     []string            `json:"logins"`
	Namespaces []string            `json:"-"`
	// +optional
	NodeLabels map[string][]string `json:"node_labels,omitempty"`
	// +optional
	Rules      []Rule              `json:"rules,omitempty"`
	// +optional
	KubeGroups []string            `json:"kubernetes_groups,omitempty"`
}

type Rule struct {
	// +optional
	Resources []string `json:"resources,omitempty"`
	// +optional
	Verbs     []string `json:"verbs,omitempty"`
	// +optional
	Where     string   `json:"where,omitempty"`
	// +optional
	Actions   []string `json:"actions,omitempty"`
}

// RoleStatus defines the observed state of Role
type RoleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Condition string `json:"condition"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	LastTransitionTime *metav1.Time `json:"last_transition_time,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Condition",type=string,JSONPath=`.status.condition`

// Role is the Schema for the roles API
type Role struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleSpec   `json:"spec,omitempty"`
	Status RoleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleList contains a list of Role
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Role `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Role{}, &RoleList{})
}
