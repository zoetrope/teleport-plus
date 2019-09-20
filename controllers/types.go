package controllers

import (
	"time"

	teleportv1 "github.com/zoetrope/teleport-plus/api/v1"
)

type Metadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"-"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Expires     *time.Time        `json:"expires,omitempty"`
}

// GithubConnectorV3 represents a Github connector
type GithubConnectorV3 struct {
	// Kind is a resource kind, for Github connector it is "github"
	Kind string `json:"kind"`
	// SubKind is a resource sub kind
	SubKind string `json:"sub_kind,omitempty"`
	// Version is resource version
	Version string `json:"version"`
	// Metadata is resource metadata
	Metadata Metadata `json:"metadata"`
	// Spec contains connector specification
	Spec teleportv1.GitHubSpec `json:"spec"`
}

type RoleV3 struct {
	Kind     string     `json:"kind"`
	SubKind  string     `json:"sub_kind,omitempty"`
	Version  string     `json:"version"`
	Metadata Metadata   `json:"metadata"`
	Spec     RoleSpecV3 `json:"spec"`
}

type RoleSpecV3 struct {
	Options RoleOptions    `json:"options,omitempty"`
	Allow   RoleConditions `json:"allow,omitempty"`
	Deny    RoleConditions `json:"deny,omitempty"`
}

type RoleOptions struct {
	ForwardAgent          bool           `json:"forward_agent"`
	MaxSessionTTL         *time.Duration `json:"max_session_ttl,omitempty"`
	PortForwarding        *bool          `json:"port_forwarding,omitempty"`
	CertificateFormat     string         `json:"cert_format"`
	ClientIdleTimeout     int64          `json:"client_idle_timeout,omitempty"`
	DisconnectExpiredCert bool           `json:"disconnect_expired_cert,omitempty"`
}

type RoleConditions struct {
	Logins     []string            `json:"logins"`
	Namespaces []string            `json:"-"`
	NodeLabels map[string][]string `json:"node_labels,omitempty"`
	Rules      []Rule              `json:"rules,omitempty"`
	KubeGroups []string            `json:"kubernetes_groups,omitempty"`
}

type Rule struct {
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
	Where     string   `json:"where,omitempty"`
	Actions   []string `json:"actions,omitempty"`
}
