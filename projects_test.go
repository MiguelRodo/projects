package projects_test

import (
	"testing"

	"github.com/MiguelRodo/projects"
)

func TestPublicAPIExports(t *testing.T) {
	ws := projects.NewWorkspace("public-ws", ".")
	if ws.Name != "public-ws" {
		t.Fatalf("expected workspace name %q, got %q", "public-ws", ws.Name)
	}

	repo := projects.Repository{
		Name: "test-repo",
		URL:  "https://github.com/example/test-repo.git",
	}
	if err := ws.AddRepository(repo); err != nil {
		t.Fatalf("failed to add repo via public API: %v", err)
	}

	info := projects.GetVersionInfo()
	if info.Version == "" {
		t.Fatalf("expected non-empty version from public API")
	}

	verStr := projects.VersionString()
	if verStr == "" {
		t.Fatalf("expected non-empty version string")
	}

	// Model exports test
	target := projects.TargetProject{
		Owner:  "example-org",
		Number: 1,
		Title:  "Roadmap",
	}
	contract := projects.NewSingleProjectContract("service-a", "https://github.com/example-org/service-a.git", target)
	if err := contract.Validate(); err != nil {
		t.Fatalf("contract validation failed via root export: %v", err)
	}
}
