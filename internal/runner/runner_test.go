package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/MiguelRodo/projects/internal/git"
	"github.com/MiguelRodo/projects/pkg/project"
)

func TestSyncWorkspace(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	ws := project.NewWorkspace("test-ws", tmpDir)
	_ = ws.AddRepository(project.Repository{
		Name: "repo-new",
		URL:  "https://github.com/example/repo-new.git",
		Path: "repo-new",
	})
	_ = ws.AddRepository(project.Repository{
		Name: "repo-existing",
		URL:  "https://github.com/example/repo-existing.git",
		Path: "repo-existing",
	})

	mock := &git.MockClient{
		IsRepoFunc: func(path string) bool {
			return filepath.Base(path) == "repo-existing"
		},
		CloneFunc: func(ctx context.Context, url, targetDir, branch string) error {
			return nil
		},
		PullFunc: func(ctx context.Context, repoDir string) (string, error) {
			return "Already up to date.", nil
		},
	}

	results := SyncWorkspace(ctx, ws, mock, SyncOptions{Concurrency: 2, Pull: true})
	if len(results) != 2 {
		t.Fatalf("expected 2 sync results, got %d", len(results))
	}

	foundCloned := false
	foundUpToDate := false
	for _, res := range results {
		if res.Repo.Name == "repo-new" && res.Action == "cloned" {
			foundCloned = true
		}
		if res.Repo.Name == "repo-existing" && res.Action == "up-to-date" {
			foundUpToDate = true
		}
	}

	if !foundCloned {
		t.Errorf("expected repo-new to be cloned")
	}
	if !foundUpToDate {
		t.Errorf("expected repo-existing to be up-to-date")
	}
}

func TestStatusWorkspace(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	ws := project.NewWorkspace("status-ws", tmpDir)
	_ = ws.AddRepository(project.Repository{
		Name: "repo1",
		URL:  "https://github.com/example/repo1.git",
		Path: "repo1",
	})

	mock := &git.MockClient{
		IsRepoFunc: func(path string) bool {
			return true
		},
		CurrentBranchFunc: func(ctx context.Context, repoDir string) (string, error) {
			return "main", nil
		},
		IsCleanFunc: func(ctx context.Context, repoDir string) (bool, error) {
			return true, nil
		},
	}

	results := StatusWorkspace(ctx, ws, mock)
	if len(results) != 1 {
		t.Fatalf("expected 1 status result, got %d", len(results))
	}
	if !results[0].Exists || results[0].Branch != "main" || !results[0].IsClean {
		t.Errorf("unexpected status result: %+v", results[0])
	}
}

func TestExecInWorkspace(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	repoDir := filepath.Join(tmpDir, "repo1")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	ws := project.NewWorkspace("exec-ws", tmpDir)
	_ = ws.AddRepository(project.Repository{
		Name: "repo1",
		URL:  "https://github.com/example/repo1.git",
		Path: "repo1",
	})

	results := ExecInWorkspace(ctx, ws, "echo", []string{"hello"}, 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 exec result, got %d", len(results))
	}
	if results[0].ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d (err: %v)", results[0].ExitCode, results[0].Err)
	}
}
