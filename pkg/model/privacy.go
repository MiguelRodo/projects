// Package model defines the neutral v1 canonical domain models, contracts,
// capabilities, mappings, and routing structures for projects.
package model

import (
	"errors"
	"fmt"
	"strings"
)

// PrivacyMode represents the privacy classification mode for task context.
type PrivacyMode string

const (
	// PrivacyModeShareableByDefault is the default mode where potentially sensitive
	// context is kept off GitHub and may be routed to user-configured companion storage.
	PrivacyModeShareableByDefault PrivacyMode = "shareable_by_default"

	// PrivacyModeFullGitHubContext allows ordinary task context to reside on GitHub.
	PrivacyModeFullGitHubContext PrivacyMode = "full_github_context"
)

var (
	// ErrInvalidPrivacyMode indicates an unsupported privacy mode was supplied.
	ErrInvalidPrivacyMode = errors.New("invalid privacy mode")
)

// PrivacyPolicy defines privacy behavior and companion storage enablement.
type PrivacyPolicy struct {
	Mode                  PrivacyMode `json:"mode" yaml:"mode"`
	AllowPrivateCompanion bool        `json:"allow_private_companion" yaml:"allow_private_companion"`
	RetentionDays         int         `json:"retention_days,omitempty" yaml:"retention_days,omitempty"`
}

// DefaultPrivacyPolicy returns a PrivacyPolicy configured with shareable_by_default.
func DefaultPrivacyPolicy() PrivacyPolicy {
	return PrivacyPolicy{
		Mode:                  PrivacyModeShareableByDefault,
		AllowPrivateCompanion: true,
		RetentionDays:         0,
	}
}

// Validate checks whether the privacy policy configuration is valid.
func (p *PrivacyPolicy) Validate() error {
	mode := strings.TrimSpace(string(p.Mode))
	if mode == "" {
		p.Mode = PrivacyModeShareableByDefault
		return nil
	}
	switch PrivacyMode(mode) {
	case PrivacyModeShareableByDefault, PrivacyModeFullGitHubContext:
		return nil
	default:
		return fmt.Errorf("%w: %q (must be %q or %q)", ErrInvalidPrivacyMode, p.Mode,
			PrivacyModeShareableByDefault, PrivacyModeFullGitHubContext)
	}
}
