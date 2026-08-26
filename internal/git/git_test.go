package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMockClient(t *testing.T) {
	ctx := context.Background()
	mock := &MockClient{
		CurrentBranchFunc: func(ctx context.Context, repoDir string) (string, error) {
			return "feature/test", nil
		},
	}

	branch, err := mock.CurrentBranch(ctx, "/some/path")
	if err != nil || branch != "feature/test" {
		t.Fatalf("expected feature/test, got %q, err: %v", branch, err)
	}

	clean, err := mock.IsClean(ctx, "/some/path")
	if err != nil || !clean {
		t.Fatalf("expected clean, got %v, err: %v", clean, err)
	}

	if err := mock.Clone(ctx, "https://example.com/repo.git", "/target", ""); err != nil {
		t.Fatalf("expected nil error on clone")
	}

	if err := mock.Fetch(ctx, "/target"); err != nil {
		t.Fatalf("expected nil error on fetch")
	}

	if _, err := mock.Pull(ctx, "/target"); err != nil {
		t.Fatalf("expected nil error on pull")
	}

	if _, err := mock.Status(ctx, "/target"); err != nil {
		t.Fatalf("expected nil error on status")
	}

	if !mock.IsRepo("/target") {
		t.Fatalf("expected is repo to be true")
	}
}

func TestExecClientIsRepo(t *testing.T) {
	client, err := NewExecClient()
	if err != nil {
		t.Skipf("git not available in environment: %v", err)
	}

	tmpDir := t.TempDir()
	if client.IsRepo(tmpDir) {
		t.Errorf("empty dir should not be recognized as git repo")
	}

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("failed to create fake .git dir: %v", err)
	}

	if !client.IsRepo(tmpDir) {
		t.Errorf("dir with .git should be recognized as git repo")
	}
}
