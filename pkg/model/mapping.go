package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyCanonicalName indicates a mapping has an empty canonical name.
	ErrEmptyCanonicalName = errors.New("canonical name cannot be empty")
	// ErrEmptyGitHubField indicates the remote field name is empty.
	ErrEmptyGitHubField = errors.New("github field name cannot be empty")
	// ErrInvalidFieldKind indicates an unsupported field kind.
	ErrInvalidFieldKind = errors.New("invalid field kind")
	// ErrEmptyValueMapping indicates source or target value is empty.
	ErrEmptyValueMapping = errors.New("value mapping cannot have empty canonical or remote value")
)

// FieldKind represents the data type of a project custom field.
type FieldKind string

const (
	// FieldKindText is a freeform text custom field.
	FieldKindText FieldKind = "text"
	// FieldKindNumber is a numerical custom field.
	FieldKindNumber FieldKind = "number"
	// FieldKindSingleSelect is a single-select option field.
	FieldKindSingleSelect FieldKind = "single_select"
	// FieldKindIteration is an iteration/sprint custom field.
	FieldKindIteration FieldKind = "iteration"
	// FieldKindDate is a date custom field.
	FieldKindDate FieldKind = "date"
)

// ValueMapping maps a canonical status/priority value to a remote GitHub option name or ID.
type ValueMapping struct {
	Canonical string `json:"canonical" yaml:"canonical"`
	Remote    string `json:"remote" yaml:"remote"`
}

// Validate checks value mapping fields.
func (v *ValueMapping) Validate() error {
	v.Canonical = strings.TrimSpace(v.Canonical)
	v.Remote = strings.TrimSpace(v.Remote)
	if v.Canonical == "" || v.Remote == "" {
		return ErrEmptyValueMapping
	}
	return nil
}

// FieldMapping defines a mapping between a canonical field and a GitHub Projects custom field.
type FieldMapping struct {
	CanonicalName string         `json:"canonical_name" yaml:"canonical_name"`
	GitHubField   string         `json:"github_field" yaml:"github_field"`
	Kind          FieldKind      `json:"kind" yaml:"kind"`
	Required      bool           `json:"required,omitempty" yaml:"required,omitempty"`
	Values        []ValueMapping `json:"values,omitempty" yaml:"values,omitempty"`
}

// Validate verifies field mapping structure and options.
func (f *FieldMapping) Validate() error {
	f.CanonicalName = strings.TrimSpace(f.CanonicalName)
	f.GitHubField = strings.TrimSpace(f.GitHubField)
	if f.CanonicalName == "" {
		return ErrEmptyCanonicalName
	}
	if f.GitHubField == "" {
		return ErrEmptyGitHubField
	}
	if f.Kind == "" {
		f.Kind = FieldKindText
	}
	switch f.Kind {
	case FieldKindText, FieldKindNumber, FieldKindSingleSelect, FieldKindIteration, FieldKindDate:
		// Valid
	default:
		return fmt.Errorf("%w: %q", ErrInvalidFieldKind, f.Kind)
	}

	for i := range f.Values {
		if err := f.Values[i].Validate(); err != nil {
			return fmt.Errorf("value mapping %d in field %q: %w", i, f.CanonicalName, err)
		}
	}
	return nil
}
