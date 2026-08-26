package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	// ErrEmptyName indicates a required name field is missing.
	ErrEmptyName = errors.New("name cannot be empty")
	// ErrInvalidSlug indicates a slug does not conform to lowercase kebab-case.
	ErrInvalidSlug = errors.New("invalid slug: must be lowercase alphanumeric and hyphens")
	// ErrEmptyOwner indicates owner is missing.
	ErrEmptyOwner = errors.New("owner cannot be empty")
	// ErrEmptyURL indicates URL is missing.
	ErrEmptyURL = errors.New("url cannot be empty")
	// ErrInvalidProjectNumber indicates a project number <= 0.
	ErrInvalidProjectNumber = errors.New("project number must be greater than zero")
	// ErrInvalidOwnerKind indicates an unsupported owner kind value was provided.
	ErrInvalidOwnerKind = errors.New("invalid owner kind: must be 'organization', 'user', or empty")
)

// OwnerKind represents the type of owner for a project or repository.
type OwnerKind string

const (
	// OwnerKindUnspecified indicates the owner kind is not yet determined.
	OwnerKindUnspecified OwnerKind = ""
	// OwnerKindOrganization represents an organization-owned resource.
	OwnerKindOrganization OwnerKind = "organization"
	// OwnerKindUser represents a user-owned resource.
	OwnerKindUser OwnerKind = "user"
)

// ValidateSlug verifies if a string matches lowercase alphanumeric kebab-case.
func ValidateSlug(s string) error {
	if !slugRegex.MatchString(s) {
		return fmt.Errorf("%w: %q", ErrInvalidSlug, s)
	}
	return nil
}

// ProjectIdentity uniquely identifies a project.
type ProjectIdentity struct {
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
	Owner       string `json:"owner,omitempty" yaml:"owner,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Validate checks whether the project identity is valid.
func (p *ProjectIdentity) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.Slug == "" {
		p.Slug = generateSlug(p.Name)
	}
	if err := ValidateSlug(p.Slug); err != nil {
		return err
	}
	return nil
}

// RepositoryRef points to a git repository.
type RepositoryRef struct {
	Name          string `json:"name" yaml:"name"`
	URL           string `json:"url" yaml:"url"`
	Path          string `json:"path,omitempty" yaml:"path,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty" yaml:"default_branch,omitempty"`
}

// Validate verifies repository reference fields.
func (r *RepositoryRef) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.URL = strings.TrimSpace(r.URL)
	if r.Name == "" {
		return ErrEmptyName
	}
	if r.URL == "" {
		return ErrEmptyURL
	}
	if r.Path == "" {
		r.Path = r.Name
	}
	if r.DefaultBranch == "" {
		r.DefaultBranch = "main"
	}
	return nil
}

// TargetProject identifies a target GitHub Projects v2 board.
type TargetProject struct {
	Ref       string    `json:"ref,omitempty" yaml:"ref,omitempty"` // Local identifier within contract
	Owner     string    `json:"owner" yaml:"owner"`
	OwnerKind OwnerKind `json:"owner_kind,omitempty" yaml:"owner_kind,omitempty"`
	Number    int       `json:"number" yaml:"number"`
	Title     string    `json:"title,omitempty" yaml:"title,omitempty"`
	URL       string    `json:"url,omitempty" yaml:"url,omitempty"`
}

// Validate verifies target project fields.
func (t *TargetProject) Validate() error {
	t.Owner = strings.TrimSpace(t.Owner)
	if t.Owner == "" {
		return ErrEmptyOwner
	}
	if t.Number <= 0 {
		return ErrInvalidProjectNumber
	}

	// Validate OwnerKind if specified; do NOT silently default to organization
	switch t.OwnerKind {
	case OwnerKindUnspecified, OwnerKindOrganization, OwnerKindUser:
		// Valid
	default:
		return fmt.Errorf("%w: %q", ErrInvalidOwnerKind, t.OwnerKind)
	}

	if t.Ref == "" {
		t.Ref = fmt.Sprintf("%s/%d", t.Owner, t.Number)
	}
	return nil
}

func generateSlug(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	var sb strings.Builder
	for _, ch := range lower {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			sb.WriteRune(ch)
		} else if ch == ' ' || ch == '-' || ch == '_' {
			if sb.Len() > 0 && !strings.HasSuffix(sb.String(), "-") {
				sb.WriteRune('-')
			}
		}
	}
	res := strings.Trim(sb.String(), "-")
	if res == "" {
		return "project"
	}
	return res
}
