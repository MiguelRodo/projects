// Package git encapsulates git CLI operations and git interaction for workspaces.
package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client defines the interface for git operations.
type Client interface {
	Clone(ctx context.Context, url, targetDir, branch string) error
	Pull(ctx context.Context, repoDir string) (string, error)
	Fetch(ctx context.Context, repoDir string) error
	Status(ctx context.Context, repoDir string) (string, error)
	CurrentBranch(ctx context.Context, repoDir string) (string, error)
	IsClean(ctx context.Context, repoDir string) (bool, error)
	IsRepo(path string) bool
}

// ExecClient implements Client using the system's git binary.
type ExecClient struct {
	GitPath string
}

// NewExecClient returns a new ExecClient with git located in PATH.
func NewExecClient() (*ExecClient, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git binary not found in PATH: %w", err)
	}
	return &ExecClient{GitPath: gitPath}, nil
}

// IsRepo checks whether path contains a .git directory or is a git repo.
func (c *ExecClient) IsRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	fi, err := os.Stat(gitDir)
	return err == nil && (fi.IsDir() || fi.Mode().IsRegular())
}

// Clone clones a git repository from url into targetDir.
func (c *ExecClient) Clone(ctx context.Context, url, targetDir, branch string) error {
	args := []string{"clone"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, url, targetDir)

	cmd := exec.CommandContext(ctx, c.GitPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone %s failed: %s (%w)", url, strings.TrimSpace(stderr.String()), err)
	}
	return nil
}

// Pull runs git pull in repoDir.
func (c *ExecClient) Pull(ctx context.Context, repoDir string) (string, error) {
	cmd := exec.CommandContext(ctx, c.GitPath, "pull")
	cmd.Dir = repoDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git pull failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Fetch runs git fetch in repoDir.
func (c *ExecClient) Fetch(ctx context.Context, repoDir string) error {
	cmd := exec.CommandContext(ctx, c.GitPath, "fetch", "--all")
	cmd.Dir = repoDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}
	return nil
}

// Status runs git status --short in repoDir.
func (c *ExecClient) Status(ctx context.Context, repoDir string) (string, error) {
	cmd := exec.CommandContext(ctx, c.GitPath, "status", "--short")
	cmd.Dir = repoDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git status failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// CurrentBranch returns the currently checked-out branch name in repoDir.
func (c *ExecClient) CurrentBranch(ctx context.Context, repoDir string) (string, error) {
	cmd := exec.CommandContext(ctx, c.GitPath, "branch", "--show-current")
	cmd.Dir = repoDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git branch failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// IsClean returns true if working tree has no unstaged/staged modifications.
func (c *ExecClient) IsClean(ctx context.Context, repoDir string) (bool, error) {
	out, err := c.Status(ctx, repoDir)
	if err != nil {
		return false, err
	}
	return len(out) == 0, nil
}
