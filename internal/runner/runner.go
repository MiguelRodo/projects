// Package runner handles concurrent and sequential execution of tasks across workspace repositories.
package runner

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/MiguelRodo/projects/internal/git"
	"github.com/MiguelRodo/projects/pkg/project"
)

// SyncResult captures the result of syncing a single repository.
type SyncResult struct {
	Repo    project.Repository `json:"repo"`
	Action  string             `json:"action"` // "cloned", "updated", "up-to-date", "skipped", "error"
	Message string             `json:"message,omitempty"`
	Err     error              `json:"error,omitempty"`
}

// StatusResult captures status of a single repository.
type StatusResult struct {
	Repo    project.Repository `json:"repo"`
	Exists  bool               `json:"exists"`
	Branch  string             `json:"branch"`
	IsClean bool               `json:"is_clean"`
	Status  string             `json:"status,omitempty"`
	Err     error              `json:"error,omitempty"`
}

// ExecResult captures the result of executing a command in a repository.
type ExecResult struct {
	Repo     project.Repository `json:"repo"`
	ExitCode int                `json:"exit_code"`
	Stdout   string             `json:"stdout"`
	Stderr   string             `json:"stderr"`
	Err      error              `json:"error,omitempty"`
}

// SyncOptions provides configuration for repository synchronization.
type SyncOptions struct {
	Concurrency int
	Pull        bool
}

// SyncWorkspace synchronizes all repositories in the workspace (cloning missing ones, pulling existing ones).
func SyncWorkspace(ctx context.Context, ws *project.Workspace, gitClient git.Client, opts SyncOptions) []SyncResult {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 4
	}

	results := make([]SyncResult, len(ws.Repositories))
	sem := make(chan struct{}, opts.Concurrency)
	var wg sync.WaitGroup

	for i, r := range ws.Repositories {
		wg.Add(1)
		go func(idx int, repo project.Repository) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[idx] = SyncResult{
					Repo:   repo,
					Action: "canceled",
					Err:    ctx.Err(),
				}
				return
			}

			targetDir := filepath.Join(ws.RootPath, repo.Path)

			if !gitClient.IsRepo(targetDir) {
				// Clone repository
				err := gitClient.Clone(ctx, repo.URL, targetDir, repo.Branch)
				if err != nil {
					results[idx] = SyncResult{
						Repo:   repo,
						Action: "error",
						Err:    err,
					}
					return
				}
				results[idx] = SyncResult{
					Repo:   repo,
					Action: "cloned",
				}
				return
			}

			// Already cloned; check if pull requested
			if opts.Pull {
				msg, err := gitClient.Pull(ctx, targetDir)
				if err != nil {
					results[idx] = SyncResult{
						Repo:   repo,
						Action: "error",
						Err:    err,
					}
					return
				}
				action := "updated"
				if msg == "Already up to date." {
					action = "up-to-date"
				}
				results[idx] = SyncResult{
					Repo:    repo,
					Action:  action,
					Message: msg,
				}
				return
			}

			results[idx] = SyncResult{
				Repo:   repo,
				Action: "skipped",
			}
		}(i, r)
	}

	wg.Wait()
	return results
}

// StatusWorkspace inspects git status and branch across all repositories in the workspace.
func StatusWorkspace(ctx context.Context, ws *project.Workspace, gitClient git.Client) []StatusResult {
	results := make([]StatusResult, len(ws.Repositories))

	for i, repo := range ws.Repositories {
		targetDir := filepath.Join(ws.RootPath, repo.Path)

		if !gitClient.IsRepo(targetDir) {
			results[i] = StatusResult{
				Repo:   repo,
				Exists: false,
			}
			continue
		}

		branch, err := gitClient.CurrentBranch(ctx, targetDir)
		if err != nil {
			results[i] = StatusResult{
				Repo:   repo,
				Exists: true,
				Err:    err,
			}
			continue
		}

		clean, err := gitClient.IsClean(ctx, targetDir)
		if err != nil {
			results[i] = StatusResult{
				Repo:   repo,
				Exists: true,
				Branch: branch,
				Err:    err,
			}
			continue
		}

		statusOut, _ := gitClient.Status(ctx, targetDir)

		results[i] = StatusResult{
			Repo:    repo,
			Exists:  true,
			Branch:  branch,
			IsClean: clean,
			Status:  statusOut,
		}
	}

	return results
}

// ExecInWorkspace executes a command in each repository directory.
func ExecInWorkspace(ctx context.Context, ws *project.Workspace, command string, args []string, concurrency int) []ExecResult {
	if concurrency <= 0 {
		concurrency = 4
	}

	results := make([]ExecResult, len(ws.Repositories))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, r := range ws.Repositories {
		wg.Add(1)
		go func(idx int, repo project.Repository) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[idx] = ExecResult{
					Repo: repo,
					Err:  ctx.Err(),
				}
				return
			}

			repoDir := filepath.Join(ws.RootPath, repo.Path)
			cmd := exec.CommandContext(ctx, command, args...)
			cmd.Dir = repoDir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					exitCode = 1
				}
			}

			results[idx] = ExecResult{
				Repo:     repo,
				ExitCode: exitCode,
				Stdout:   stdout.String(),
				Stderr:   stderr.String(),
				Err:      err,
			}
		}(i, r)
	}

	wg.Wait()
	return results
}
