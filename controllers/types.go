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
	Kind     string              `json:"kind"`
	SubKind  string              `json:"sub_kind,omitempty"`
	Version  string              `json:"version"`
	Metadata Metadata            `json:"metadata"`
	Spec     teleportv1.RoleSpec `json:"spec"`
}
