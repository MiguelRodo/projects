// Package projects provides tools, libraries, and workspace management
// for multi-repository software development environments.
package projects

import (
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
)
