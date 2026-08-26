# projects / projectctl

[![CI](https://github.com/MiguelRodo/projects/actions/workflows/ci.yml/badge.svg)](https://github.com/MiguelRodo/projects/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/MiguelRodo/projects.svg)](https://pkg.go.dev/github.com/MiguelRodo/projects)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/MiguelRodo/projects)](https://golang.org)

`projects` is a Go library and CLI tool (`projectctl`) designed for managing multi-repository workspaces, distributed project collections, and containerized development setups.

---

## Architecture & Public Boundary

The repository is organized following idiomatic Go package boundaries with strict separation between the public API surface and private internal implementation details:

```
├── cmd/
│   ├── projectctl/       # Entrypoint binary: projectctl CLI
│   └── projects/         # Entrypoint binary: projects CLI alias
├── pkg/
│   ├── project/          # [PUBLIC] Core domain models, interfaces, and manifest parsers
│   └── version/          # [PUBLIC] Version metadata and build inspection
├── internal/             # [PRIVATE] Encapsulated implementation details (unimportable externally)
│   ├── cli/              # CLI argument parser, formatting, and subcommand routing
│   ├── config/           # Manifest discovery, loading, and persistence logic
│   ├── git/              # Git operations wrapper and client abstractions
│   └── runner/           # Multi-repo task runner and concurrent orchestration
├── projects.go           # [PUBLIC] Root package convenience exports for Go consumers
├── .github/workflows/    # CI/CD and release automation
├── .devcontainer/        # Containerized development environment
├── Makefile              # Build, test, format, and lint automation
├── go.mod                # Module definition: github.com/MiguelRodo/projects
└── LICENSE               # MIT License
```

### Public Boundary Guarantees
- **`github.com/MiguelRodo/projects`** and **`github.com/MiguelRodo/projects/pkg/*`** expose stable types (`Workspace`, `Repository`), constructors, and manifest parsers (`ParseReposList`, `LoadWorkspaceFromJSON`, etc.).
- **`github.com/MiguelRodo/projects/internal/*`** is protected by Go compiler rules, ensuring internal CLI mechanics, git execution details, and configuration IO cannot be inadvertently coupled by external consumers.

---

## Installation

### Using `go install`
```bash
go install github.com/MiguelRodo/projects/cmd/projectctl@latest
go install github.com/MiguelRodo/projects/cmd/projects@latest
```

### From Source
```bash
git clone https://github.com/MiguelRodo/projects.git
cd projects
make build
# Binaries will be placed in ./bin/projectctl and ./bin/projects
```

---

## CLI Usage (`projectctl`)

`projectctl` (or `projects`) provides commands to manage workspace repositories:

```text
Usage:
  projectctl [flags] <command> [command-flags] [arguments]

Commands:
  init      Initialize a new workspace manifest (projects.json or repos.list)
  list      List repositories configured in the workspace
  add       Add a new repository to the workspace manifest
  remove    Remove a repository from the workspace manifest
  sync      Clone missing repositories and pull updates for existing ones
  status    Check git status and branch across all repositories
  exec      Execute a command across all repositories in the workspace
  version   Show version and build information
  help      Show help for projectctl

Flags:
  -C, --workspace <path>   Set workspace root directory (default: current directory)
      --config <path>      Path to workspace manifest file (default: auto-detect)
  -v, --version            Show version information
  -h, --help               Show help information
```

### Examples

#### 1. Initialize a Workspace
```bash
projectctl init --name my-workspace
```

#### 2. Add Repositories
```bash
projectctl add https://github.com/org/repo1.git
projectctl add https://github.com/org/repo2.git --name repo2 --path custom/path --branch develop
```

#### 3. List Configured Repositories
```bash
projectctl list
# Output in JSON format:
projectctl list --json
```

#### 4. Synchronize Repositories
Clones missing repositories and pulls latest commits concurrently:
```bash
projectctl sync --pull --concurrency 4
```

#### 5. Check Git Status Across All Repositories
```bash
projectctl status
```

#### 6. Execute Commands Across All Repositories
```bash
projectctl exec -- git fetch --all
projectctl exec -- make test
```

---

## Go Library Usage

To import and use `projects` in your own Go applications:

```go
package main

import (
	"fmt"
	"log"

	"github.com/MiguelRodo/projects"
)

func main() {
	ws := projects.NewWorkspace("my-project", ".")

	err := ws.AddRepository(projects.Repository{
		Name: "core-service",
		URL:  "https://github.com/org/core-service.git",
		Path: "services/core",
	})
	if err != nil {
		log.Fatalf("failed to add repository: %v", err)
	}

	fmt.Printf("Workspace: %s, Repositories: %d\n", ws.Name, len(ws.Repositories))
}
```

---

## Development

Requirements:
- Go 1.22+
- `golangci-lint`

### Common Makefile Targets
```bash
make test       # Run unit tests with race detection & coverage
make lint       # Run golangci-lint
make fmt        # Format source code
make vet        # Run go vet
make build      # Build binaries into bin/
make all        # Run fmt, vet, lint, test, and build
make clean      # Remove build artifacts
```

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.