// Package model defines the neutral v1 canonical domain models, contracts,
// capabilities, mappings, and routing structures for projects.
package model

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// PrivacyMode represents the privacy classification mode chosen by an acting user or supported by a repository.
type PrivacyMode string

const (
	// PrivacyModeShareableByDefault is the default mode where potentially sensitive
	// task context is kept off GitHub and may be routed to user-configured companion storage.
	PrivacyModeShareableByDefault PrivacyMode = "shareable_by_default"

	// PrivacyModeFullGitHubContext allows ordinary task context to reside on GitHub.
	PrivacyModeFullGitHubContext PrivacyMode = "full_github_context"
)

var (
	// ErrInvalidPrivacyMode indicates an unsupported privacy mode was supplied.
	ErrInvalidPrivacyMode = errors.New("invalid privacy mode")
	// ErrUnsupportedPrivacyMode indicates the requested privacy mode is not supported by the repository contract.
	ErrUnsupportedPrivacyMode = errors.New("unsupported privacy mode")
	// ErrPrivateCompanionNotAllowed indicates companion linkage was configured when disallowed by policy.
	ErrPrivateCompanionNotAllowed = errors.New("private companion linkage is not allowed by repository policy")
)

// RepositoryPrivacyPolicy defines the privacy modes and companion capabilities supported by a shared repository.
// Note: The actual effective privacy mode is an acting-user choice configured in the user's private operator profile.
type RepositoryPrivacyPolicy struct {
	SupportedModes        []PrivacyMode `json:"supported_modes,omitempty" yaml:"supported_modes,omitempty"`
	DefaultMode           PrivacyMode   `json:"default_mode,omitempty" yaml:"default_mode,omitempty"`
	AllowPrivateCompanion bool          `json:"allow_private_companion" yaml:"allow_private_companion"`
}

// DefaultRepositoryPrivacyPolicy returns a RepositoryPrivacyPolicy configured with standard defaults.
func DefaultRepositoryPrivacyPolicy() RepositoryPrivacyPolicy {
	return RepositoryPrivacyPolicy{
		SupportedModes: []PrivacyMode{
			PrivacyModeShareableByDefault,
			PrivacyModeFullGitHubContext,
		},
		DefaultMode:           PrivacyModeShareableByDefault,
		AllowPrivateCompanion: false,
	}
}

// Validate checks whether the shared repository privacy policy configuration is valid.
func (p *RepositoryPrivacyPolicy) Validate() error {
	if p.DefaultMode == "" {
		p.DefaultMode = PrivacyModeShareableByDefault
	}
	if err := validatePrivacyMode(p.DefaultMode); err != nil {
		return fmt.Errorf("default mode: %w", err)
	}

	if len(p.SupportedModes) == 0 {
		p.SupportedModes = []PrivacyMode{p.DefaultMode}
	}

	for i, m := range p.SupportedModes {
		if err := validatePrivacyMode(m); err != nil {
			return fmt.Errorf("supported mode %d: %w", i, err)
		}
	}

	if !slices.Contains(p.SupportedModes, p.DefaultMode) {
		p.SupportedModes = append(p.SupportedModes, p.DefaultMode)
	}

	return nil
}

// Supports checks whether a specific privacy mode is supported by the repository.
func (p *RepositoryPrivacyPolicy) Supports(mode PrivacyMode) bool {
	return slices.Contains(p.SupportedModes, mode)
}

// UserPrivacyPreference represents an acting user's private configuration for a repository session.
type UserPrivacyPreference struct {
	EffectiveMode       PrivacyMode `json:"effective_mode" yaml:"effective_mode"`
	PrivateCompanionRef string      `json:"private_companion_ref,omitempty" yaml:"private_companion_ref,omitempty"`
}

// Validate checks whether the user's privacy preference is valid against a repository policy.
func (u *UserPrivacyPreference) Validate(repoPolicy *RepositoryPrivacyPolicy) error {
	if u.EffectiveMode == "" {
		if repoPolicy != nil && repoPolicy.DefaultMode != "" {
			u.EffectiveMode = repoPolicy.DefaultMode
		} else {
			u.EffectiveMode = PrivacyModeShareableByDefault
		}
	}

	if err := validatePrivacyMode(u.EffectiveMode); err != nil {
		return err
	}

	if repoPolicy != nil {
		if !repoPolicy.Supports(u.EffectiveMode) {
			return fmt.Errorf("%w: repository does not support mode %q", ErrUnsupportedPrivacyMode, u.EffectiveMode)
		}
		if strings.TrimSpace(u.PrivateCompanionRef) != "" && !repoPolicy.AllowPrivateCompanion {
			return fmt.Errorf("%w: cannot configure private companion ref %q",
				ErrPrivateCompanionNotAllowed, u.PrivateCompanionRef)
		}
	}

	return nil
}

func validatePrivacyMode(mode PrivacyMode) error {
	trimmed := strings.TrimSpace(string(mode))
	switch PrivacyMode(trimmed) {
	case PrivacyModeShareableByDefault, PrivacyModeFullGitHubContext:
		return nil
	default:
		return fmt.Errorf("%w: %q (must be %q or %q)", ErrInvalidPrivacyMode, mode,
			PrivacyModeShareableByDefault, PrivacyModeFullGitHubContext)
	}
}
