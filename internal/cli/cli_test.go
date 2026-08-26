package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MiguelRodo/projects/internal/git"
)

func TestAppVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Execute(context.Background(), []string{"version"})
	if code != 0 {
		t.Fatalf("expected code 0, got %d. stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "projectctl") {
		t.Errorf("expected version output, got: %s", stdout.String())
	}

	stdout.Reset()
	code = app.Execute(context.Background(), []string{"version", "--json"})
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), `"version":`) {
		t.Errorf("expected JSON version output, got: %s", stdout.String())
	}
}

func TestAppHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Execute(context.Background(), []string{"help"})
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Errorf("expected help text, got: %s", stdout.String())
	}

	stdout.Reset()
	code = app.Execute(context.Background(), []string{})
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Errorf("expected help text, got: %s", stdout.String())
	}
}

func TestAppInitAndListAndAddAndRemove(t *testing.T) {
	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)
	app.WorkingDir = tmpDir

	// Init
	code := app.Execute(context.Background(), []string{"init", "--name", "test-project"})
	if code != 0 {
		t.Fatalf("init failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Initialized empty workspace") {
		t.Errorf("unexpected init output: %s", stdout.String())
	}

	// Double init should fail
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"init", "--name", "test-project"})
	if code != 1 {
		t.Fatalf("expected code 1 on double init, got %d", code)
	}

	// Add repository
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"add", "https://github.com/org/repo-a.git"})
	if code != 0 {
		t.Fatalf("add failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Added repository") {
		t.Errorf("unexpected add output: %s", stdout.String())
	}

	// Add second repository with options
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"add", "https://github.com/org/repo-b.git", "--name", "custom-b", "--path", "pkg/b", "--branch", "develop"})
	if code != 0 {
		t.Fatalf("add second failed with code %d: %s", code, stderr.String())
	}

	// List repositories
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"list"})
	if code != 0 {
		t.Fatalf("list failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "repo-a") || !strings.Contains(stdout.String(), "custom-b") {
		t.Errorf("list output missing repos: %s", stdout.String())
	}

	// List in JSON
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"list", "--json"})
	if code != 0 {
		t.Fatalf("list json failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"name": "repo-a"`) {
		t.Errorf("json list missing repo-a: %s", stdout.String())
	}

	// Remove repository
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"remove", "repo-a"})
	if code != 0 {
		t.Fatalf("remove failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Removed repository") {
		t.Errorf("unexpected remove output: %s", stdout.String())
	}

	// Verify removal in list
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"list"})
	if code != 0 {
		t.Fatalf("list failed: %s", stderr.String())
	}
	if strings.Contains(stdout.String(), "repo-a") {
		t.Errorf("expected repo-a to be removed from list: %s", stdout.String())
	}
}

func TestAppSyncAndStatus(t *testing.T) {
	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)
	app.WorkingDir = tmpDir
	app.GitClient = &git.MockClient{
		IsRepoFunc: func(path string) bool {
			return strings.Contains(path, "repo1")
		},
		CloneFunc: func(ctx context.Context, url, targetDir, branch string) error {
			return nil
		},
		PullFunc: func(ctx context.Context, repoDir string) (string, error) {
			return "Already up to date.", nil
		},
		CurrentBranchFunc: func(ctx context.Context, repoDir string) (string, error) {
			return "main", nil
		},
		IsCleanFunc: func(ctx context.Context, repoDir string) (bool, error) {
			return true, nil
		},
	}

	// Add a repo
	code := app.Execute(context.Background(), []string{"add", "https://github.com/org/repo1.git"})
	if code != 0 {
		t.Fatalf("add failed: %s", stderr.String())
	}

	// Sync
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"sync"})
	if code != 0 {
		t.Fatalf("sync failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Syncing 1 repositories") {
		t.Errorf("unexpected sync output: %s", stdout.String())
	}

	// Status
	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"status"})
	if code != 0 {
		t.Fatalf("status failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "repo1") {
		t.Errorf("status output missing repo1: %s", stdout.String())
	}
}

func TestAppExec(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo1")
	_ = os.MkdirAll(repoDir, 0755)

	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)
	app.WorkingDir = tmpDir

	code := app.Execute(context.Background(), []string{"add", "https://github.com/org/repo1.git"})
	if code != 0 {
		t.Fatalf("add failed: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = app.Execute(context.Background(), []string{"exec", "--", "echo", "hello-world"})
	if code != 0 {
		t.Fatalf("exec failed with code %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "hello-world") {
		t.Errorf("exec missing output: %s", stdout.String())
	}
}

func TestRunHelper(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected Run to return 0, got %d", code)
	}
}
