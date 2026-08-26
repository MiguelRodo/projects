// Package config provides workspace configuration discovery, loading, and persistence.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MiguelRodo/projects/pkg/project"
)

const (
	// DefaultProjectsFile is the primary workspace json config name.
	DefaultProjectsFile = "projects.json"
	// HiddenProjectsFile is an alternative dotfile workspace config.
	HiddenProjectsFile = ".projects.json"
	// DefaultReposListFile is the fallback repos.list text format.
	DefaultReposListFile = "repos.list"
)

// FindConfigFile scans rootDir for standard workspace manifest files in priority order.
func FindConfigFile(rootDir string) (string, string) {
	candidates := []struct {
		name string
		kind string
	}{
		{DefaultProjectsFile, "json"},
		{HiddenProjectsFile, "json"},
		{DefaultReposListFile, "list"},
	}

	for _, c := range candidates {
		fullPath := filepath.Join(rootDir, c.name)
		if fi, err := os.Stat(fullPath); err == nil && !fi.IsDir() {
			return fullPath, c.kind
		}
	}
	return "", ""
}

// LoadWorkspace loads a workspace from an explicit path or by auto-discovery in rootDir.
func LoadWorkspace(rootDir, configPath string) (*project.Workspace, string, error) {
	kind := ""
	if configPath != "" {
		ext := filepath.Ext(configPath)
		if ext == ".json" {
			kind = "json"
		} else {
			kind = "list"
		}
	} else {
		configPath, kind = FindConfigFile(rootDir)
	}

	if configPath == "" {
		// Default empty workspace if no config file found
		wsName := filepath.Base(rootDir)
		if wsName == "." || wsName == "/" {
			wsName = "projects"
		}
		ws := project.NewWorkspace(wsName, rootDir)
		return ws, "", nil
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, configPath, fmt.Errorf("opening config file %q: %w", configPath, err)
	}
	defer func() { _ = f.Close() }()

	if kind == "json" {
		ws, err := project.LoadWorkspaceFromJSON(f)
		if err != nil {
			return nil, configPath, err
		}
		if ws.RootPath == "" {
			ws.RootPath = rootDir
		}
		return ws, configPath, nil
	}

	// Parse repos.list
	repos, err := project.ParseReposList(f)
	if err != nil {
		return nil, configPath, err
	}

	wsName := filepath.Base(rootDir)
	ws := project.NewWorkspace(wsName, rootDir)
	ws.Repositories = repos

	return ws, configPath, nil
}

// SaveWorkspace saves the workspace back to disk at the specified path or default location.
func SaveWorkspace(ws *project.Workspace, targetPath string) error {
	if targetPath == "" {
		targetPath = filepath.Join(ws.RootPath, DefaultProjectsFile)
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("creating workspace file %q: %w", targetPath, err)
	}
	defer func() { _ = f.Close() }()

	if filepath.Ext(targetPath) == ".json" {
		return project.SaveWorkspaceToJSON(f, ws)
	}

	// Save as repos.list format
	content := project.FormatReposList(ws.Repositories)
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing repos.list: %w", err)
	}
	return nil
}
