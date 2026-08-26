// Package project defines core domain models, interfaces, and manifest parsers
// for managing multi-repository projects and workspaces.
package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	// ErrEmptyRepoName indicates a repository was provided without a valid name.
	ErrEmptyRepoName = errors.New("repository name cannot be empty")
	// ErrEmptyRepoURL indicates a repository was provided without a URL.
	ErrEmptyRepoURL = errors.New("repository URL cannot be empty")
	// ErrDuplicateRepo indicates a repository with the same name or path already exists in the workspace.
	ErrDuplicateRepo = errors.New("repository already exists in workspace")
	// ErrRepoNotFound indicates the specified repository does not exist in the workspace.
	ErrRepoNotFound = errors.New("repository not found in workspace")
	// ErrEmptyWorkspaceName indicates the workspace name is empty.
	ErrEmptyWorkspaceName = errors.New("workspace name cannot be empty")
)

// Repository represents a git repository managed within a project workspace.
type Repository struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Path        string `json:"path" yaml:"path"`
	Branch      string `json:"branch,omitempty" yaml:"branch,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Validate checks whether the repository definition is valid.
func (r *Repository) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrEmptyRepoName
	}
	if strings.TrimSpace(r.URL) == "" {
		return ErrEmptyRepoURL
	}
	if strings.TrimSpace(r.Path) == "" {
		r.Path = r.Name
	}
	return nil
}

// Workspace represents a collection of repositories under a common workspace root.
type Workspace struct {
	Name         string       `json:"name" yaml:"name"`
	RootPath     string       `json:"root_path,omitempty" yaml:"root_path,omitempty"`
	Repositories []Repository `json:"repositories" yaml:"repositories"`
}

// NewWorkspace creates a new Workspace instance with the given name and root path.
func NewWorkspace(name, rootPath string) *Workspace {
	if rootPath == "" {
		rootPath = "."
	}
	return &Workspace{
		Name:         name,
		RootPath:     rootPath,
		Repositories: make([]Repository, 0),
	}
}

// Validate verifies the workspace structure and its repository definitions.
func (w *Workspace) Validate() error {
	if strings.TrimSpace(w.Name) == "" {
		return ErrEmptyWorkspaceName
	}
	names := make(map[string]struct{})
	paths := make(map[string]struct{})

	for i := range w.Repositories {
		repo := &w.Repositories[i]
		if err := repo.Validate(); err != nil {
			return fmt.Errorf("repository %d is invalid: %w", i, err)
		}
		if _, exists := names[repo.Name]; exists {
			return fmt.Errorf("%w: duplicate name %q", ErrDuplicateRepo, repo.Name)
		}
		cleanPath := filepath.Clean(repo.Path)
		if _, exists := paths[cleanPath]; exists {
			return fmt.Errorf("%w: duplicate path %q", ErrDuplicateRepo, repo.Path)
		}
		names[repo.Name] = struct{}{}
		paths[cleanPath] = struct{}{}
	}
	return nil
}

// FindRepository returns a pointer to the repository with the given name, if found.
func (w *Workspace) FindRepository(name string) (*Repository, bool) {
	for i := range w.Repositories {
		if strings.EqualFold(w.Repositories[i].Name, name) {
			return &w.Repositories[i], true
		}
	}
	return nil, false
}

// AddRepository adds a new repository to the workspace after validation.
func (w *Workspace) AddRepository(repo Repository) error {
	if err := repo.Validate(); err != nil {
		return err
	}
	if _, found := w.FindRepository(repo.Name); found {
		return fmt.Errorf("%w: %q", ErrDuplicateRepo, repo.Name)
	}
	cleanPath := filepath.Clean(repo.Path)
	for _, r := range w.Repositories {
		if filepath.Clean(r.Path) == cleanPath {
			return fmt.Errorf("%w: path %q already in use by %q", ErrDuplicateRepo, repo.Path, r.Name)
		}
	}
	w.Repositories = append(w.Repositories, repo)
	return nil
}

// RemoveRepository removes a repository from the workspace by name.
func (w *Workspace) RemoveRepository(name string) error {
	for i, r := range w.Repositories {
		if strings.EqualFold(r.Name, name) {
			w.Repositories = append(w.Repositories[:i], w.Repositories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %q", ErrRepoNotFound, name)
}
