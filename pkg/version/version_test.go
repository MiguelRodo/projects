package version

import (
	"strings"
	"testing"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	if info.Version != Version {
		t.Errorf("expected Version %q, got %q", Version, info.Version)
	}
	if info.GitCommit != GitCommit {
		t.Errorf("expected GitCommit %q, got %q", GitCommit, info.GitCommit)
	}
	if info.BuildDate != BuildDate {
		t.Errorf("expected BuildDate %q, got %q", BuildDate, info.BuildDate)
	}
	if info.GoVersion == "" {
		t.Errorf("expected non-empty GoVersion")
	}
	if info.Compiler == "" {
		t.Errorf("expected non-empty Compiler")
	}
	if info.Platform == "" {
		t.Errorf("expected non-empty Platform")
	}
}

func TestString(t *testing.T) {
	str := String()
	if !strings.Contains(str, "projectctl") {
		t.Errorf("expected %q to contain 'projectctl'", str)
	}
	if !strings.Contains(str, Version) {
		t.Errorf("expected %q to contain Version %q", str, Version)
	}
}
