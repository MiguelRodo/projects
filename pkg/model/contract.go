package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptySchemaVersion indicates the schema version header is missing.
	ErrEmptySchemaVersion = errors.New("schema version cannot be empty")
	// ErrDuplicateFieldMapping indicates multiple field mappings share the same canonical name.
	ErrDuplicateFieldMapping = errors.New("duplicate canonical field mapping name")
	// ErrTargetNotFound indicates a target reference could not be resolved.
	ErrTargetNotFound = errors.New("target project reference not found")
)

// SingleProjectContract represents the canonical in-memory model of a repository
// mapped to a single GitHub Project.
type SingleProjectContract struct {
	SchemaVersion string                  `json:"schema_version,omitempty" yaml:"schema_version,omitempty"`
	Project       ProjectIdentity         `json:"project" yaml:"project"`
	Repository    RepositoryRef           `json:"repository" yaml:"repository"`
	Target        TargetProject           `json:"target" yaml:"target"`
	Mappings      []FieldMapping          `json:"mappings,omitempty" yaml:"mappings,omitempty"`
	Capabilities  CapabilitySet           `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Mutation      MutationPolicy          `json:"mutation,omitempty" yaml:"mutation,omitempty"`
	Privacy       RepositoryPrivacyPolicy `json:"privacy,omitempty" yaml:"privacy,omitempty"`
	Sources       []SourceRef             `json:"sources,omitempty" yaml:"sources,omitempty"`
}

// NewSingleProjectContract initializes a SingleProjectContract with safe defaults.
func NewSingleProjectContract(name, repoURL string, target TargetProject) SingleProjectContract {
	_ = target.Validate()
	return SingleProjectContract{
		Project: ProjectIdentity{
			Name: name,
			Slug: generateSlug(name),
		},
		Repository: RepositoryRef{
			Name: name,
			URL:  repoURL,
		},
		Target:       target,
		Capabilities: DefaultCapabilities(),
		Mutation:     DefaultMutationPolicy(),
		Privacy:      DefaultRepositoryPrivacyPolicy(),
	}
}

// Validate verifies all fields and invariants for a SingleProjectContract.
func (c *SingleProjectContract) Validate() error {
	if err := c.Project.Validate(); err != nil {
		return fmt.Errorf("invalid project identity: %w", err)
	}
	if err := c.Repository.Validate(); err != nil {
		return fmt.Errorf("invalid repository reference: %w", err)
	}
	if err := c.Target.Validate(); err != nil {
		return fmt.Errorf("invalid target project: %w", err)
	}
	if err := c.Privacy.Validate(); err != nil {
		return fmt.Errorf("invalid privacy policy: %w", err)
	}
	if err := c.Mutation.Validate(); err != nil {
		return fmt.Errorf("invalid mutation policy: %w", err)
	}

	seenFields := make(map[string]struct{})
	for i := range c.Mappings {
		m := &c.Mappings[i]
		if err := m.Validate(); err != nil {
			return fmt.Errorf("mapping %d: %w", i, err)
		}
		lower := strings.ToLower(m.CanonicalName)
		if _, exists := seenFields[lower]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateFieldMapping, m.CanonicalName)
		}
		seenFields[lower] = struct{}{}
	}

	for i := range c.Sources {
		if err := c.Sources[i].Validate(); err != nil {
			return fmt.Errorf("source %d: %w", i, err)
		}
	}

	return nil
}

// MultiProjectContract represents the canonical in-memory model of a multi-project dispatcher.
type MultiProjectContract struct {
	SchemaVersion string                  `json:"schema_version,omitempty" yaml:"schema_version,omitempty"`
	IssueStore    RepositoryRef           `json:"issue_store" yaml:"issue_store"`
	Projects      []ProjectIdentity       `json:"projects,omitempty" yaml:"projects,omitempty"`
	Targets       []TargetProject         `json:"targets" yaml:"targets"`
	Dispatcher    DispatcherConfig        `json:"dispatcher" yaml:"dispatcher"`
	Mappings      []FieldMapping          `json:"mappings,omitempty" yaml:"mappings,omitempty"`
	Capabilities  CapabilitySet           `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Mutation      MutationPolicy          `json:"mutation,omitempty" yaml:"mutation,omitempty"`
	Privacy       RepositoryPrivacyPolicy `json:"privacy,omitempty" yaml:"privacy,omitempty"`
	Sources       []SourceRef             `json:"sources,omitempty" yaml:"sources,omitempty"`
}

// NewMultiProjectContract initializes a MultiProjectContract with safe defaults.
// It ensures target references are properly generated before populating the dispatcher.
func NewMultiProjectContract(issueStore RepositoryRef, defaultTarget TargetProject) MultiProjectContract {
	_ = defaultTarget.Validate()
	if defaultTarget.Ref == "" && defaultTarget.Owner != "" && defaultTarget.Number > 0 {
		defaultTarget.Ref = fmt.Sprintf("%s/%d", defaultTarget.Owner, defaultTarget.Number)
	}

	caps := DefaultCapabilities()
	caps.Add(CapabilityRouteDispatcher)

	return MultiProjectContract{
		IssueStore: issueStore,
		Targets:    []TargetProject{defaultTarget},
		Dispatcher: DispatcherConfig{
			DefaultTargetRef: defaultTarget.Ref,
			Fallback:         FallbackBehaviorDefaultTarget,
		},
		Capabilities: caps,
		Mutation:     DefaultMutationPolicy(),
		Privacy:      DefaultRepositoryPrivacyPolicy(),
	}
}

// FindTarget looks up a TargetProject by reference string.
func (m *MultiProjectContract) FindTarget(ref string) (*TargetProject, bool) {
	for i := range m.Targets {
		if m.Targets[i].Ref == ref || m.Targets[i].Title == ref {
			return &m.Targets[i], true
		}
	}
	return nil, false
}

// Validate verifies all fields and invariants for a MultiProjectContract.
func (m *MultiProjectContract) Validate() error {
	if err := m.IssueStore.Validate(); err != nil {
		return fmt.Errorf("invalid issue store: %w", err)
	}
	if len(m.Targets) == 0 {
		return errors.New("multi-project contract must declare at least one target project")
	}

	targetRefs := make(map[string]struct{})
	for i := range m.Targets {
		t := &m.Targets[i]
		if err := t.Validate(); err != nil {
			return fmt.Errorf("target %d: %w", i, err)
		}
		if _, exists := targetRefs[t.Ref]; exists {
			return fmt.Errorf("duplicate target reference %q", t.Ref)
		}
		targetRefs[t.Ref] = struct{}{}
	}

	if err := m.Dispatcher.Validate(); err != nil {
		return fmt.Errorf("invalid dispatcher config: %w", err)
	}

	// Verify routes point to declared targets
	for i, r := range m.Dispatcher.Routes {
		if _, exists := targetRefs[r.TargetRef]; !exists {
			return fmt.Errorf("route %d points to undeclared target %q", i, r.TargetRef)
		}
	}

	if m.Dispatcher.DefaultTargetRef != "" {
		if _, exists := targetRefs[m.Dispatcher.DefaultTargetRef]; !exists {
			return fmt.Errorf("default target %q is not declared in targets", m.Dispatcher.DefaultTargetRef)
		}
	}

	if err := m.Privacy.Validate(); err != nil {
		return fmt.Errorf("invalid privacy policy: %w", err)
	}
	if err := m.Mutation.Validate(); err != nil {
		return fmt.Errorf("invalid mutation policy: %w", err)
	}

	seenFields := make(map[string]struct{})
	for i := range m.Mappings {
		mapping := &m.Mappings[i]
		if err := mapping.Validate(); err != nil {
			return fmt.Errorf("mapping %d: %w", i, err)
		}
		lower := strings.ToLower(mapping.CanonicalName)
		if _, exists := seenFields[lower]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateFieldMapping, mapping.CanonicalName)
		}
		seenFields[lower] = struct{}{}
	}

	for i := range m.Sources {
		if err := m.Sources[i].Validate(); err != nil {
			return fmt.Errorf("source %d: %w", i, err)
		}
	}

	return nil
}
