// Package version provides version and build information for the project.
package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current semantic version of the application.
	Version = "v0.1.0-dev"
	// GitCommit is the git SHA-1 commit hash of the build.
	GitCommit = "unknown"
	// BuildDate is the RFC3339 timestamp when the binary was built.
	BuildDate = "unknown"
)

// Info holds detailed build and version metadata.
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// GetInfo returns the current build and version metadata.
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a human-readable summary of the version information.
func String() string {
	return fmt.Sprintf("projectctl %s (commit: %s, built: %s, go: %s, platform: %s/%s)",
		Version, GitCommit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
