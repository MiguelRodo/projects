package model

import (
	"errors"
	"testing"
)

func TestPrivacyPolicyValidation(t *testing.T) {
	p := DefaultPrivacyPolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("default privacy policy invalid: %v", err)
	}
	if p.Mode != PrivacyModeShareableByDefault {
		t.Errorf("expected mode %q, got %q", PrivacyModeShareableByDefault, p.Mode)
	}

	p.Mode = PrivacyModeFullGitHubContext
	if err := p.Validate(); err != nil {
		t.Fatalf("full github context should be valid: %v", err)
	}

	p.Mode = "invalid_mode"
	if err := p.Validate(); !errors.Is(err, ErrInvalidPrivacyMode) {
		t.Fatalf("expected ErrInvalidPrivacyMode, got: %v", err)
	}
}

func TestProjectIdentityValidation(t *testing.T) {
	pi := ProjectIdentity{
		Name: "Core Platform",
	}
	if err := pi.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if pi.Slug != "core-platform" {
		t.Errorf("expected auto slug 'core-platform', got %q", pi.Slug)
	}

	piEmpty := ProjectIdentity{Name: ""}
	if err := piEmpty.Validate(); !errors.Is(err, ErrEmptyName) {
		t.Fatalf("expected ErrEmptyName, got: %v", err)
	}

	piInvalidSlug := ProjectIdentity{
		Name: "Test",
		Slug: "Invalid_Slug!",
	}
	if err := piInvalidSlug.Validate(); !errors.Is(err, ErrInvalidSlug) {
		t.Fatalf("expected ErrInvalidSlug, got: %v", err)
	}
}

func TestRepositoryRefValidation(t *testing.T) {
	ref := RepositoryRef{
		Name: "repo-alpha",
		URL:  "https://github.com/example-org/repo-alpha.git",
	}
	if err := ref.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if ref.Path != "repo-alpha" || ref.DefaultBranch != "main" {
		t.Errorf("unexpected defaults: path=%q, branch=%q", ref.Path, ref.DefaultBranch)
	}

	refNoURL := RepositoryRef{Name: "repo"}
	if err := refNoURL.Validate(); !errors.Is(err, ErrEmptyURL) {
		t.Fatalf("expected ErrEmptyURL, got: %v", err)
	}
}

func TestTargetProjectValidation(t *testing.T) {
	target := TargetProject{
		Owner:  "example-org",
		Number: 42,
		Title:  "Global Roadmap",
	}
	if err := target.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if target.OwnerKind != OwnerKindOrganization {
		t.Errorf("expected default owner kind 'organization', got %q", target.OwnerKind)
	}
	if target.Ref != "example-org/42" {
		t.Errorf("expected ref 'example-org/42', got %q", target.Ref)
	}

	targetBadNum := TargetProject{Owner: "example-org", Number: 0}
	if err := targetBadNum.Validate(); !errors.Is(err, ErrInvalidProjectNumber) {
		t.Fatalf("expected ErrInvalidProjectNumber, got: %v", err)
	}
}

func TestFieldMappingValidation(t *testing.T) {
	mapping := FieldMapping{
		CanonicalName: "status",
		GitHubField:   "Status",
		Kind:          FieldKindSingleSelect,
		Values: []ValueMapping{
			{Canonical: "todo", Remote: "Todo"},
			{Canonical: "done", Remote: "Done"},
		},
	}
	if err := mapping.Validate(); err != nil {
		t.Fatalf("mapping validation failed: %v", err)
	}

	badMapping := FieldMapping{CanonicalName: ""}
	if err := badMapping.Validate(); !errors.Is(err, ErrEmptyCanonicalName) {
		t.Fatalf("expected ErrEmptyCanonicalName, got: %v", err)
	}

	badValue := FieldMapping{
		CanonicalName: "status",
		GitHubField:   "Status",
		Kind:          FieldKindSingleSelect,
		Values: []ValueMapping{
			{Canonical: "todo", Remote: ""},
		},
	}
	if err := badValue.Validate(); !errors.Is(err, ErrEmptyValueMapping) {
		t.Fatalf("expected ErrEmptyValueMapping, got: %v", err)
	}
}

func TestSingleProjectContract(t *testing.T) {
	target := TargetProject{
		Owner:  "example-org",
		Number: 10,
		Title:  "Service Board",
	}
	contract := NewSingleProjectContract("service-a", "https://github.com/example-org/service-a.git", target)

	if err := contract.Validate(); err != nil {
		t.Fatalf("contract validation failed: %v", err)
	}

	contract.Mappings = append(contract.Mappings,
		FieldMapping{CanonicalName: "priority", GitHubField: "Priority", Kind: FieldKindSingleSelect},
		FieldMapping{CanonicalName: "priority", GitHubField: "Priority2", Kind: FieldKindText},
	)
	if err := contract.Validate(); !errors.Is(err, ErrDuplicateFieldMapping) {
		t.Fatalf("expected ErrDuplicateFieldMapping, got: %v", err)
	}
}

func TestMultiProjectContractRouting(t *testing.T) {
	store := RepositoryRef{
		Name: "issue-tracker",
		URL:  "https://github.com/example-org/issue-tracker.git",
	}
	targetBackend := TargetProject{
		Ref:    "backend-board",
		Owner:  "example-org",
		Number: 1,
		Title:  "Backend Tasks",
	}
	targetFrontend := TargetProject{
		Ref:    "frontend-board",
		Owner:  "example-org",
		Number: 2,
		Title:  "Frontend Tasks",
	}

	contract := NewMultiProjectContract(store, targetBackend)
	contract.Targets = append(contract.Targets, targetFrontend)
	contract.Dispatcher.Routes = []RouteRule{
		{Key: "label", Value: "area/backend", TargetRef: "backend-board"},
		{Key: "label", Value: "area/frontend", TargetRef: "frontend-board"},
	}

	if err := contract.Validate(); err != nil {
		t.Fatalf("multi-project validation failed: %v", err)
	}

	// Resolve routes
	res, err := contract.ResolveTarget("label", "area/backend")
	if err != nil || res.Ref != "backend-board" {
		t.Fatalf("expected backend-board, got %+v (err: %v)", res, err)
	}

	res, err = contract.ResolveTarget("label", "area/frontend")
	if err != nil || res.Ref != "frontend-board" {
		t.Fatalf("expected frontend-board, got %+v (err: %v)", res, err)
	}

	// Fallback to default
	res, err = contract.ResolveTarget("label", "unmatched")
	if err != nil || res.Ref != "backend-board" {
		t.Fatalf("expected default backend-board, got %+v (err: %v)", res, err)
	}
}

func TestSourceCatalogueValidation(t *testing.T) {
	cat := SourceCatalogue{
		Version: "v1",
		Sources: []SourceRef{
			{
				ID:          "api-spec",
				Name:        "API Specification",
				Location:    "docs/api/spec.yaml",
				Type:        SourceTypeDoc,
				Description: "Public API contract",
			},
		},
	}
	if err := cat.Validate(); err != nil {
		t.Fatalf("catalogue validation failed: %v", err)
	}
}
