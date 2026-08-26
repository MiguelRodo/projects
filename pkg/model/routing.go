package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyRouteKey indicates the routing key is empty.
	ErrEmptyRouteKey = errors.New("route key cannot be empty")
	// ErrEmptyRouteTarget indicates the route target reference is empty.
	ErrEmptyRouteTarget = errors.New("route target cannot be empty")
	// ErrInvalidFallbackBehavior indicates an unsupported fallback behavior.
	ErrInvalidFallbackBehavior = errors.New("invalid fallback behavior")
	// ErrDuplicateRouteRule indicates multiple route rules define the same key-value condition.
	ErrDuplicateRouteRule = errors.New("duplicate route rule condition")
)

// FallbackBehavior specifies the intended behavior when an item does not match any route.
type FallbackBehavior string

const (
	// FallbackBehaviorDefaultTarget routes to the designated default project.
	FallbackBehaviorDefaultTarget FallbackBehavior = "default_target"
	// FallbackBehaviorError indicates routing failure when no matching route exists.
	FallbackBehaviorError FallbackBehavior = "error"
	// FallbackBehaviorIgnore indicates skipping routing for unmatched items.
	FallbackBehaviorIgnore FallbackBehavior = "ignore"
)

// RouteRule represents a mapping condition from a key/value pair to a target project reference.
type RouteRule struct {
	Key           string            `json:"key" yaml:"key"` // e.g. "label", "repository", "component"
	Value         string            `json:"value" yaml:"value"`
	TargetRef     string            `json:"target_ref" yaml:"target_ref"` // Matches TargetProject.Ref
	DefaultLabels []string          `json:"default_labels,omitempty" yaml:"default_labels,omitempty"`
	DefaultFields map[string]string `json:"default_fields,omitempty" yaml:"default_fields,omitempty"`
}

// Validate checks route rule fields.
func (r *RouteRule) Validate() error {
	r.Key = strings.TrimSpace(r.Key)
	r.Value = strings.TrimSpace(r.Value)
	r.TargetRef = strings.TrimSpace(r.TargetRef)
	if r.Key == "" {
		return ErrEmptyRouteKey
	}
	if r.TargetRef == "" {
		return ErrEmptyRouteTarget
	}
	return nil
}

// DispatcherConfig configures routing targets and rules for multi-project topologies.
type DispatcherConfig struct {
	DefaultTargetRef string           `json:"default_target_ref,omitempty" yaml:"default_target_ref,omitempty"`
	Fallback         FallbackBehavior `json:"fallback,omitempty" yaml:"fallback,omitempty"`
	Routes           []RouteRule      `json:"routes" yaml:"routes"`
}

// Validate checks dispatcher configuration invariants.
func (d *DispatcherConfig) Validate() error {
	if d.Fallback == "" {
		d.Fallback = FallbackBehaviorDefaultTarget
	}
	switch d.Fallback {
	case FallbackBehaviorDefaultTarget, FallbackBehaviorError, FallbackBehaviorIgnore:
		// Valid
	default:
		return fmt.Errorf("%w: %q", ErrInvalidFallbackBehavior, d.Fallback)
	}

	if d.Fallback == FallbackBehaviorDefaultTarget && strings.TrimSpace(d.DefaultTargetRef) == "" && len(d.Routes) == 0 {
		return errors.New("dispatcher requires either a default_target_ref or explicit routes")
	}

	seenConditions := make(map[string]struct{})
	for i := range d.Routes {
		r := &d.Routes[i]
		if err := r.Validate(); err != nil {
			return fmt.Errorf("route %d: %w", i, err)
		}
		conditionKey := fmt.Sprintf("%s=%s", strings.ToLower(r.Key), strings.ToLower(r.Value))
		if _, exists := seenConditions[conditionKey]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateRouteRule, conditionKey)
		}
		seenConditions[conditionKey] = struct{}{}
	}
	return nil
}
