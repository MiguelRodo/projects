package model

import (
	"slices"
)

// Capability defines supported protocol capabilities for a repository or dispatcher.
type Capability string

const (
	// CapabilityReadIssues indicates issues can be fetched.
	CapabilityReadIssues Capability = "read_issues"
	// CapabilityWriteIssues indicates issues can be created or updated.
	CapabilityWriteIssues Capability = "write_issues"
	// CapabilityManageProjects indicates Project item cards can be moved or updated.
	CapabilityManageProjects Capability = "manage_projects"
	// CapabilityRouteDispatcher indicates multi-project routing is supported.
	CapabilityRouteDispatcher Capability = "route_dispatcher"
	// CapabilityPrivateCompanion indicates private companion record linkage is supported.
	CapabilityPrivateCompanion Capability = "private_companion"
)

// CapabilitySet encapsulates enabled capabilities.
type CapabilitySet struct {
	Items []Capability `json:"items,omitempty" yaml:"items,omitempty"`
}

// DefaultCapabilities returns baseline core capabilities without opt-in extensions.
func DefaultCapabilities() CapabilitySet {
	return CapabilitySet{
		Items: []Capability{
			CapabilityReadIssues,
			CapabilityWriteIssues,
			CapabilityManageProjects,
		},
	}
}

// Has checks whether a capability is enabled.
func (c *CapabilitySet) Has(cap Capability) bool {
	return slices.Contains(c.Items, cap)
}

// Add appends a capability if not already present.
func (c *CapabilitySet) Add(cap Capability) {
	if !c.Has(cap) {
		c.Items = append(c.Items, cap)
	}
}

// MutationPolicy defines permitted write operations and stale-write guards.
type MutationPolicy struct {
	AllowCreate     bool `json:"allow_create" yaml:"allow_create"`
	AllowUpdate     bool `json:"allow_update" yaml:"allow_update"`
	AllowDelete     bool `json:"allow_delete" yaml:"allow_delete"`
	StaleWriteGuard bool `json:"stale_write_guard" yaml:"stale_write_guard"`
}

// DefaultMutationPolicy returns a standard safe mutation policy.
func DefaultMutationPolicy() MutationPolicy {
	return MutationPolicy{
		AllowCreate:     true,
		AllowUpdate:     true,
		AllowDelete:     false,
		StaleWriteGuard: true,
	}
}

// Validate checks mutation policy settings.
func (m *MutationPolicy) Validate() error {
	return nil
}
