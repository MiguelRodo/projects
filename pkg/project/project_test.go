package project

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRepositoryValidate(t *testing.T) {
	tests := []struct {
		name    string
		repo    Repository
		wantErr bool
		errType error
	}{
		{
			name: "valid repo",
			repo: Repository{
				Name: "my-repo",
				URL:  "https://github.com/example/my-repo.git",
				Path: "src/my-repo",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			repo: Repository{
				Name: "",
				URL:  "https://github.com/example/my-repo.git",
			},
			wantErr: true,
			errType: ErrEmptyRepoName,
		},
		{
			name: "empty url",
			repo: Repository{
				Name: "my-repo",
				URL:  "",
			},
			wantErr: true,
			errType: ErrEmptyRepoURL,
		},
		{
			name: "auto fills path if empty",
			repo: Repository{
				Name: "my-repo",
				URL:  "https://github.com/example/my-repo.git",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Fatalf("expected error type %v, got %v", tt.errType, err)
			}
			if !tt.wantErr && tt.repo.Path == "" {
				t.Errorf("expected Path to be defaulted to %q", tt.repo.Name)
			}
		})
	}
}

func TestWorkspaceOperations(t *testing.T) {
	ws := NewWorkspace("test-ws", "/tmp/test")

	if ws.Name != "test-ws" || ws.RootPath != "/tmp/test" {
		t.Fatalf("unexpected workspace initial state: %+v", ws)
	}

	repo1 := Repository{
		Name: "repo1",
		URL:  "https://github.com/example/repo1.git",
		Path: "repo1",
	}
	repo2 := Repository{
		Name: "repo2",
		URL:  "https://github.com/example/repo2.git",
		Path: "repo2",
	}

	if err := ws.AddRepository(repo1); err != nil {
		t.Fatalf("failed to add repo1: %v", err)
	}
	if err := ws.AddRepository(repo2); err != nil {
		t.Fatalf("failed to add repo2: %v", err)
	}

	if len(ws.Repositories) != 2 {
		t.Fatalf("expected 2 repositories, got %d", len(ws.Repositories))
	}

	// Duplicate add
	if err := ws.AddRepository(repo1); err == nil {
		t.Fatalf("expected duplicate error, got nil")
	}

	// Find repo
	found, ok := ws.FindRepository("repo1")
	if !ok || found.Name != "repo1" {
		t.Fatalf("failed to find repo1")
	}
	_, ok = ws.FindRepository("nonexistent")
	if ok {
		t.Fatalf("expected nonexistent repo not to be found")
	}

	// Remove repo
	if err := ws.RemoveRepository("repo1"); err != nil {
		t.Fatalf("failed to remove repo1: %v", err)
	}
	if len(ws.Repositories) != 1 {
		t.Fatalf("expected 1 repository after removal, got %d", len(ws.Repositories))
	}
	if err := ws.RemoveRepository("nonexistent"); err == nil {
		t.Fatalf("expected error removing nonexistent repo")
	}

	// Validation
	if err := ws.Validate(); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	emptyNameWs := &Workspace{Name: ""}
	if err := emptyNameWs.Validate(); !errors.Is(err, ErrEmptyWorkspaceName) {
		t.Fatalf("expected ErrEmptyWorkspaceName, got %v", err)
	}
}

func TestReposListParsingAndFormatting(t *testing.T) {
	input := `# Comment line
https://github.com/org/repo1.git
https://github.com/org/repo2.git custom/path
https://github.com/org/repo3.git src/repo3 main
`
	repos, err := ParseReposList(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseReposList failed: %v", err)
	}

	if len(repos) != 3 {
		t.Fatalf("expected 3 repos, got %d", len(repos))
	}

	if repos[0].Name != "repo1" || repos[0].Path != "repo1" {
		t.Errorf("unexpected repo 0: %+v", repos[0])
	}
	if repos[1].Name != "path" || repos[1].Path != "custom/path" {
		t.Errorf("unexpected repo 1: %+v", repos[1])
	}
	if repos[2].Name != "repo3" || repos[2].Path != "src/repo3" || repos[2].Branch != "main" {
		t.Errorf("unexpected repo 2: %+v", repos[2])
	}

	formatted := FormatReposList(repos)
	if !strings.Contains(formatted, "https://github.com/org/repo1.git") {
		t.Errorf("formatted text missing repo1: %s", formatted)
	}

	reParsed, err := ParseReposList(strings.NewReader(formatted))
	if err != nil {
		t.Fatalf("re-parsing formatted text failed: %v", err)
	}
	if len(reParsed) != 3 {
		t.Fatalf("expected 3 repos on reparsing, got %d", len(reParsed))
	}
}

func TestJSONWorkspaceSerialization(t *testing.T) {
	ws := NewWorkspace("my-projects", "/workspace")
	_ = ws.AddRepository(Repository{
		Name: "core",
		URL:  "https://github.com/example/core.git",
		Path: "core",
	})

	var buf bytes.Buffer
	if err := SaveWorkspaceToJSON(&buf, ws); err != nil {
		t.Fatalf("SaveWorkspaceToJSON failed: %v", err)
	}

	loaded, err := LoadWorkspaceFromJSON(&buf)
	if err != nil {
		t.Fatalf("LoadWorkspaceFromJSON failed: %v", err)
	}

	if loaded.Name != ws.Name || len(loaded.Repositories) != 1 {
		t.Fatalf("loaded workspace does not match: %+v", loaded)
	}
}
