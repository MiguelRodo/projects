// Package projects provides tools, libraries, and workspace management
// for multi-repository software development environments and neutral v1 contracts.
package projects

import (
	"github.com/MiguelRodo/projects/pkg/model"
	"github.com/MiguelRodo/projects/pkg/project"
	"github.com/MiguelRodo/projects/pkg/version"
)

// Re-export core public types for clean and convenient top-level access.
type (
	// Repository represents a single git repository in a workspace.
	Repository = project.Repository
	// Workspace represents a collection of managed repositories.
	Workspace = project.Workspace
	// VersionInfo holds version and build metadata.
	VersionInfo = version.Info

	// SingleProjectContract represents a single-project repository interface contract.
	SingleProjectContract = model.SingleProjectContract
	// MultiProjectContract represents a multi-project dispatcher interface contract.
	MultiProjectContract = model.MultiProjectContract
	// TargetProject identifies a target GitHub Projects v2 board.
	TargetProject = model.TargetProject
	// OwnerKind represents the type of owner for a project (organization or user).
	OwnerKind = model.OwnerKind
	// RepositoryPrivacyPolicy defines privacy modes and companion capabilities supported by a shared repository.
	RepositoryPrivacyPolicy = model.RepositoryPrivacyPolicy
	// UserPrivacyPreference represents an acting user's private session configuration.
	UserPrivacyPreference = model.UserPrivacyPreference
	// PrivacyMode represents privacy classification modes.
	PrivacyMode = model.PrivacyMode
	// FieldMapping defines a mapping between a canonical field and a GitHub Projects custom field.
	FieldMapping = model.FieldMapping
	// RouteRule maps a key/value condition to a target project.
	RouteRule = model.RouteRule
	// DispatcherConfig configures routing for multi-project topologies.
	DispatcherConfig = model.DispatcherConfig
)

const (
	PrivacyModeShareableByDefault = model.PrivacyModeShareableByDefault
	PrivacyModeFullGitHubContext  = model.PrivacyModeFullGitHubContext

	OwnerKindOrganization = model.OwnerKindOrganization
	OwnerKindUser         = model.OwnerKindUser
	OwnerKindUnspecified  = model.OwnerKindUnspecified
)

// Common public constructors and helpers.
var (
	// NewWorkspace creates a new workspace instance.
	NewWorkspace = project.NewWorkspace
	// ParseReposList parses repository definitions from a repos.list formatted reader.
	ParseReposList = project.ParseReposList
	// FormatReposList formats repositories into repos.list string format.
	FormatReposList = project.FormatReposList
	// LoadWorkspaceFromJSON loads a workspace from JSON reader.
	LoadWorkspaceFromJSON = project.LoadWorkspaceFromJSON
	// SaveWorkspaceToJSON writes a workspace to a JSON writer.
	SaveWorkspaceToJSON = project.SaveWorkspaceToJSON
	// GetVersionInfo returns application version metadata.
	GetVersionInfo = version.GetInfo
	// VersionString returns human-readable version string.
	VersionString = version.String

	// NewSingleProjectContract creates a new single-project contract with safe defaults.
	NewSingleProjectContract = model.NewSingleProjectContract
	// NewMultiProjectContract creates a new multi-project contract with safe defaults.
	NewMultiProjectContract = model.NewMultiProjectContract
	// DefaultRepositoryPrivacyPolicy returns a RepositoryPrivacyPolicy with shareable_by_default.
	DefaultRepositoryPrivacyPolicy = model.DefaultRepositoryPrivacyPolicy
)
