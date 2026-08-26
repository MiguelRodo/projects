package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MiguelRodo/projects/pkg/project"
)

func TestFindConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	path, kind := FindConfigFile(tmpDir)
	if path != "" || kind != "" {
		t.Fatalf("expected no config in empty dir, got %s (%s)", path, kind)
	}

	reposList := filepath.Join(tmpDir, "repos.list")
	if err := os.WriteFile(reposList, []byte("https://github.com/org/repo1.git\n"), 0644); err != nil {
		t.Fatalf("failed to write repos.list: %v", err)
	}

	path, kind = FindConfigFile(tmpDir)
	if path != reposList || kind != "list" {
		t.Fatalf("expected %s (list), got %s (%s)", reposList, path, kind)
	}

	projectsJSON := filepath.Join(tmpDir, "projects.json")
	if err := os.WriteFile(projectsJSON, []byte(`{"name":"test","repositories":[]}`), 0644); err != nil {
		t.Fatalf("failed to write projects.json: %v", err)
	}

	// projects.json should take priority over repos.list
	path, kind = FindConfigFile(tmpDir)
	if path != projectsJSON || kind != "json" {
		t.Fatalf("expected %s (json), got %s (%s)", projectsJSON, path, kind)
	}
}

func TestLoadAndSaveWorkspace(t *testing.T) {
	tmpDir := t.TempDir()

	ws := project.NewWorkspace("demo", tmpDir)
	_ = ws.AddRepository(project.Repository{
		Name: "svc1",
		URL:  "https://github.com/example/svc1.git",
		Path: "svc1",
	})

	jsonPath := filepath.Join(tmpDir, "projects.json")
	if err := SaveWorkspace(ws, jsonPath); err != nil {
		t.Fatalf("SaveWorkspace failed: %v", err)
	}

	loaded, loadedPath, err := LoadWorkspace(tmpDir, "")
	if err != nil {
		t.Fatalf("LoadWorkspace failed: %v", err)
	}
	if loadedPath != jsonPath {
		t.Errorf("expected loadedPath %q, got %q", jsonPath, loadedPath)
	}
	if loaded.Name != "demo" || len(loaded.Repositories) != 1 {
		t.Errorf("loaded workspace mismatch: %+v", loaded)
	}

	// Save as repos.list
	listPath := filepath.Join(tmpDir, "repos.list")
	if err := SaveWorkspace(ws, listPath); err != nil {
		t.Fatalf("SaveWorkspace as list failed: %v", err)
	}

	loadedList, _, err := LoadWorkspace(tmpDir, listPath)
	if err != nil {
		t.Fatalf("LoadWorkspace list failed: %v", err)
	}
	if len(loadedList.Repositories) != 1 || loadedList.Repositories[0].Name != "svc1" {
		t.Errorf("loaded repos.list mismatch: %+v", loadedList)
	}
}
