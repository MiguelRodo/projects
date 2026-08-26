package model

import (
	"errors"
	"testing"
)

func TestPrivacyPolicyAndUserPreference(t *testing.T) {
	repoPolicy := DefaultRepositoryPrivacyPolicy()
	if err := repoPolicy.Validate(); err != nil {
		t.Fatalf("default repository privacy policy invalid: %v", err)
	}
	if repoPolicy.DefaultMode != PrivacyModeShareableByDefault {
		t.Errorf("expected default mode %q, got %q", PrivacyModeShareableByDefault, repoPolicy.DefaultMode)
	}
	if repoPolicy.AllowPrivateCompanion {
		t.Errorf("expected AllowPrivateCompanion to default to false (opt-in)")
	}

	// Regression test for Review Note 2: Explicit SupportedModes not containing DefaultMode fails
	mismatchedPolicy := RepositoryPrivacyPolicy{
		SupportedModes: []PrivacyMode{PrivacyModeFullGitHubContext},
		DefaultMode:    PrivacyModeShareableByDefault,
	}
	if err := mismatchedPolicy.Validate(); !errors.Is(err, ErrInvalidPrivacyPolicy) {
		t.Fatalf("expected ErrInvalidPrivacyPolicy for explicit mismatched SupportedModes, got: %v", err)
	}

	// User preference validation without companion
	userPref := UserPrivacyPreference{
		EffectiveMode: PrivacyModeFullGitHubContext,
	}
	if err := userPref.Validate(&repoPolicy); err != nil {
		t.Fatalf("valid user preference failed: %v", err)
	}

	// User preference attempting companion when repo policy has AllowPrivateCompanion=false
	userPrefWithCompanion := UserPrivacyPreference{
		EffectiveMode:       PrivacyModeShareableByDefault,
		PrivateCompanionRef: "personal-notes-doc",
	}
	if err := userPrefWithCompanion.Validate(&repoPolicy); !errors.Is(err, ErrPrivateCompanionNotAllowed) {
		t.Fatalf("expected ErrPrivateCompanionNotAllowed when repo policy disallows companion, got: %v", err)
	}

	// User preference with companion when repo policy enables it
	repoPolicy.AllowPrivateCompanion = true
	if err := userPrefWithCompanion.Validate(&repoPolicy); err != nil {
		t.Fatalf("user preference with companion should succeed when enabled: %v", err)
	}

	// User preference default fallback
	userPrefDefault := UserPrivacyPreference{}
	if err := userPrefDefault.Validate(&repoPolicy); err != nil {
		t.Fatalf("empty user preference should inherit default: %v", err)
	}
	if userPrefDefault.EffectiveMode != PrivacyModeShareableByDefault {
		t.Errorf("expected effective mode %q, got %q", PrivacyModeShareableByDefault, userPrefDefault.EffectiveMode)
	}

	// User preference with unsupported mode
	strictPolicy := RepositoryPrivacyPolicy{
		SupportedModes: []PrivacyMode{PrivacyModeShareableByDefault},
		DefaultMode:    PrivacyModeShareableByDefault,
	}
	_ = strictPolicy.Validate()
	userPrefUnsupported := UserPrivacyPreference{
		EffectiveMode: PrivacyModeFullGitHubContext,
	}
	if err := userPrefUnsupported.Validate(&strictPolicy); !errors.Is(err, ErrUnsupportedPrivacyMode) {
		t.Fatalf("expected ErrUnsupportedPrivacyMode, got: %v", err)
	}
}

func TestStableLinkageIDAndCompanion(t *testing.T) {
	validID := StableLinkageID("task-42-companion-notes")
	if err := validID.Validate(); err != nil {
		t.Fatalf("expected valid linkage ID: %v", err)
	}

	emptyID := StableLinkageID("")
	if err := emptyID.Validate(); !errors.Is(err, ErrEmptyLinkageID) {
		t.Fatalf("expected ErrEmptyLinkageID, got: %v", err)
	}

	invalidID := StableLinkageID("Invalid_Slug!")
	if err := invalidID.Validate(); !errors.Is(err, ErrInvalidSlug) {
		t.Fatalf("expected ErrInvalidSlug, got: %v", err)
	}

	companion := CompanionLinkage{
		LinkageID:   validID,
		IssueNumber: 42,
		Repository:  "repo-alpha",
	}
	if err := companion.Validate(); err != nil {
		t.Fatalf("expected valid companion linkage: %v", err)
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
	// User-owned project
	targetUser := TargetProject{
		Owner:     "alice",
		OwnerKind: OwnerKindUser,
		Number:    7,
		Title:     "Personal Tasks",
	}
	if err := targetUser.Validate(); err != nil {
		t.Fatalf("user target validation failed: %v", err)
	}
	if targetUser.OwnerKind != OwnerKindUser {
		t.Errorf("expected owner kind 'user', got %q", targetUser.OwnerKind)
	}
	if targetUser.Ref != "alice/7" {
		t.Errorf("expected ref 'alice/7', got %q", targetUser.Ref)
	}

	// Organization-owned project
	targetOrg := TargetProject{
		Owner:     "example-org",
		OwnerKind: OwnerKindOrganization,
		Number:    42,
		Title:     "Global Roadmap",
	}
	if err := targetOrg.Validate(); err != nil {
		t.Fatalf("org target validation failed: %v", err)
	}

	// Target without explicit owner kind (remains unspecified, not defaulted)
	targetUnspecified := TargetProject{
		Owner:  "example-org",
		Number: 10,
	}
	if err := targetUnspecified.Validate(); err != nil {
		t.Fatalf("unspecified target validation failed: %v", err)
	}
	if targetUnspecified.OwnerKind != OwnerKindUnspecified {
		t.Errorf("expected OwnerKindUnspecified, got %q", targetUnspecified.OwnerKind)
	}

	// Invalid owner kind
	targetBadKind := TargetProject{
		Owner:     "example-org",
		OwnerKind: "invalid_kind",
		Number:    1,
	}
	if err := targetBadKind.Validate(); !errors.Is(err, ErrInvalidOwnerKind) {
		t.Fatalf("expected ErrInvalidOwnerKind, got: %v", err)
	}

	targetBadNum := TargetProject{Owner: "example-org", Number: 0}
	if err := targetBadNum.Validate(); !errors.Is(err, ErrInvalidProjectNumber) {
		t.Fatalf("expected ErrInvalidProjectNumber, got: %v", err)
	}
}

func TestFieldMappingValidation(t *testing.T) {
	mapping := FieldMapping{
		CanonicalName: "task-status",
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

	badSlugMapping := FieldMapping{
		CanonicalName: "Invalid_Name!",
		GitHubField:   "Status",
	}
	if err := badSlugMapping.Validate(); !errors.Is(err, ErrInvalidSlug) {
		t.Fatalf("expected ErrInvalidSlug for non-kebab-case canonical name, got: %v", err)
	}

	emptySingleSelect := FieldMapping{
		CanonicalName: "priority",
		GitHubField:   "Priority",
		Kind:          FieldKindSingleSelect,
		Values:        nil,
	}
	if err := emptySingleSelect.Validate(); !errors.Is(err, ErrEmptySingleSelectValues) {
		t.Fatalf("expected ErrEmptySingleSelectValues, got: %v", err)
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

	// Test private companion capability vs policy mismatch:
	// 1. AllowPrivateCompanion=true but capability absent
	contract.Privacy.AllowPrivateCompanion = true
	if err := contract.Validate(); !errors.Is(err, ErrPrivateCompanionMismatch) {
		t.Fatalf("expected ErrPrivateCompanionMismatch when AllowPrivateCompanion=true without capability, got: %v", err)
	}

	// 2. Capability present but AllowPrivateCompanion=false
	contract.Privacy.AllowPrivateCompanion = false
	contract.Capabilities.Add(CapabilityPrivateCompanion)
	if err := contract.Validate(); !errors.Is(err, ErrPrivateCompanionMismatch) {
		t.Fatalf("expected ErrPrivateCompanionMismatch when CapabilityPrivateCompanion is present but AllowPrivateCompanion=false, got: %v", err)
	}

	// 3. Both present -> valid
	contract.Privacy.AllowPrivateCompanion = true
	if err := contract.Validate(); err != nil {
		t.Fatalf("contract should be valid when both are aligned: %v", err)
	}

	// Reset companion settings
	contract.Privacy.AllowPrivateCompanion = false
	contract.Capabilities = DefaultCapabilities()

	contract.Mappings = append(contract.Mappings,
		FieldMapping{
			CanonicalName: "priority",
			GitHubField:   "Priority",
			Kind:          FieldKindSingleSelect,
			Values:        []ValueMapping{{Canonical: "high", Remote: "High"}},
		},
		FieldMapping{
			CanonicalName: "priority",
			GitHubField:   "Priority2",
			Kind:          FieldKindText,
		},
	)
	if err := contract.Validate(); !errors.Is(err, ErrDuplicateFieldMapping) {
		t.Fatalf("expected ErrDuplicateFieldMapping, got: %v", err)
	}
}

func TestNewMultiProjectContract_NoExplicitRefRegression(t *testing.T) {
	store := RepositoryRef{
		Name: "issues",
		URL:  "https://github.com/example-org/issues.git",
	}
	defaultTarget := TargetProject{
		Owner:  "example-org",
		Number: 1,
		Title:  "Default Board",
	}

	contract := NewMultiProjectContract(store, defaultTarget)
	if err := contract.Validate(); err != nil {
		t.Fatalf("NewMultiProjectContract failed validation with auto-generated target ref: %v", err)
	}
	if contract.Dispatcher.DefaultTargetRef != "example-org/1" {
		t.Errorf("expected Dispatcher.DefaultTargetRef to be 'example-org/1', got %q", contract.Dispatcher.DefaultTargetRef)
	}

	// Ref-only FindTarget lookup test
	found, exists := contract.FindTarget("example-org/1")
	if !exists || found == nil {
		t.Fatalf("expected FindTarget by ref 'example-org/1' to succeed")
	}

	_, existsByTitle := contract.FindTarget("Default Board")
	if existsByTitle {
		t.Fatalf("FindTarget should be ref-only and not match by title")
	}

	// Test companion capability mismatch in multi-project contract
	contract.Privacy.AllowPrivateCompanion = true
	if err := contract.Validate(); !errors.Is(err, ErrPrivateCompanionMismatch) {
		t.Fatalf("expected ErrPrivateCompanionMismatch in MultiProjectContract, got: %v", err)
	}
}

func TestMultiProjectContractDuplicateRoutes(t *testing.T) {
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

	contract := NewMultiProjectContract(store, targetBackend)
	contract.Dispatcher.Routes = []RouteRule{
		{Key: "label", Value: "backend", TargetRef: "backend-board"},
		{Key: "label", Value: "backend", TargetRef: "backend-board"},
	}

	if err := contract.Validate(); !errors.Is(err, ErrDuplicateRouteRule) {
		t.Fatalf("expected ErrDuplicateRouteRule for duplicate route condition, got: %v", err)
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
