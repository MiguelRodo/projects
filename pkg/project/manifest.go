package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// ParseReposList reads repository definitions from a repos.list formatted stream.
// Lines starting with '#' or empty lines are ignored.
// Each line has the format:
//
//	<URL> [path] [branch]
//
// If path is omitted, the base repository name is extracted from the URL.
func ParseReposList(r io.Reader) ([]Repository, error) {
	var repos []Repository
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		url := fields[0]
		name := extractRepoNameFromURL(url)
		path := name
		branch := ""

		if len(fields) >= 2 {
			path = fields[1]
			name = filepath.Base(path)
		}
		if len(fields) >= 3 {
			branch = fields[2]
		}

		repo := Repository{
			Name:   name,
			URL:    url,
			Path:   path,
			Branch: branch,
		}

		if err := repo.Validate(); err != nil {
			return nil, fmt.Errorf("line %d: invalid repository definition: %w", lineNum, err)
		}

		repos = append(repos, repo)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading repos.list: %w", err)
	}

	return repos, nil
}

// FormatReposList formats a slice of repositories into the repos.list text format.
func FormatReposList(repos []Repository) string {
	var sb strings.Builder
	sb.WriteString("# Managed repositories list\n")
	sb.WriteString("# Format: <url> [path] [branch]\n\n")

	for _, repo := range repos {
		sb.WriteString(repo.URL)
		if repo.Path != "" && repo.Path != repo.Name {
			sb.WriteString(" ")
			sb.WriteString(repo.Path)
			if repo.Branch != "" {
				sb.WriteString(" ")
				sb.WriteString(repo.Branch)
			}
		} else if repo.Branch != "" {
			sb.WriteString(" ")
			sb.WriteString(repo.Path)
			sb.WriteString(" ")
			sb.WriteString(repo.Branch)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// LoadWorkspaceFromJSON reads and parses a Workspace from a JSON stream.
func LoadWorkspaceFromJSON(r io.Reader) (*Workspace, error) {
	var ws Workspace
	dec := json.NewDecoder(r)
	if err := dec.Decode(&ws); err != nil {
		return nil, fmt.Errorf("decoding workspace json: %w", err)
	}
	if err := ws.Validate(); err != nil {
		return nil, fmt.Errorf("validating workspace: %w", err)
	}
	return &ws, nil
}

// SaveWorkspaceToJSON encodes a Workspace into formatted JSON and writes to w.
func SaveWorkspaceToJSON(w io.Writer, ws *Workspace) error {
	if err := ws.Validate(); err != nil {
		return fmt.Errorf("validating workspace before saving: %w", err)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ws); err != nil {
		return fmt.Errorf("encoding workspace json: %w", err)
	}
	return nil
}

// extractRepoNameFromURL extracts the repository name from a git URL.
func extractRepoNameFromURL(rawURL string) string {
	trimmed := strings.TrimRight(rawURL, "/")
	trimmed = strings.TrimSuffix(trimmed, ".git")

	// Handles git@github.com:org/repo or https://github.com/org/repo
	parts := strings.Split(trimmed, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if colonIdx := strings.LastIndex(lastPart, ":"); colonIdx != -1 {
			lastPart = lastPart[colonIdx+1:]
		}
		if lastPart != "" {
			return lastPart
		}
	}
	return "repository"
}
