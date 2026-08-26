# `projects` / `projectctl`

[![CI](https://github.com/MiguelRodo/projects/actions/workflows/ci.yml/badge.svg)](https://github.com/MiguelRodo/projects/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/MiguelRodo/projects.svg)](https://pkg.go.dev/github.com/MiguelRodo/projects)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/MiguelRodo/projects)](https://golang.org)

`projects` is a portable, agent-neutral project-administration specification, Go library, and CLI tool (`projectctl`). It enables multi-repository workspace management, structured project contracts, and multi-agent coordination across software repositories.

---

## Core Principles

- **GitHub is the Shared Authority**: Conforming repositories use GitHub (issues, pull requests, Projects, labels) as the single source of truth.
- **Private Control Planes are Optional Consumers**: Personal agendas, private control repositories, and external storage (e.g. Google Drive/Docs/Sheets) are strictly optional extensions. They operate as downstream consumers and must never be a hard dependency for ordinary contributors or automated bots.
- **Agent & Operator Neutrality**: All public schemas, contracts, and tools use neutral vocabulary. Configuration is decoupled from specific persons, cloud accounts, proprietary control systems, or particular AI agent runtimes.
- **Strict Privacy Separation**: Public repository manifests contain only public metadata and non-sensitive stable linkage identifiers. Sensitive companion data remains exclusively on user-controlled private destinations.

---

## Three Paths of Participation

`projects` supports three distinct user and agent journeys:

```
┌────────────────────────────────────────────────────────────────────────┐
│ Path 1: Collaborator / External Agent                                   │
│ • Interacts exclusively via standard GitHub issues & PRs               │
│ • Zero private control-plane or external toolchain dependencies         │
└────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Path 2: Repository Maintainer                                          │
│ • Adopts repositories & GitHub Projects into neutral schemas           │
│ • Maps existing fields non-destructively & configures validation CI     │
│ • Generates agent instructions and verification reports                │
└────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Path 3: Full-System Operator                                           │
│ • Configures optional private control repository & multi-project router │
│ • Manages personal agendas, inventories, and custom source providers    │
│ • Preserves clean boundary: private config never leaks into public repos│
└────────────────────────────────────────────────────────────────────────┘
```

### Path 1: Collaborator or External Agent
For developers, external contributors, and AI agents contributing to a conforming repository:
- **Read the Interface**: Understand the project layout through public manifests (`projects.json` or `repos.list`) and standard GitHub metadata.
- **Participate Naturally**: Create, discuss, and update standard GitHub issues and pull requests without requiring access to private operator infrastructure.
- **Zero Friction**: No proprietary credentials, private control repositories, or custom daemon setups are required.

### Path 2: Repository Maintainer
For repository owners organizing one or more repositories under structured project governance:
- **Adopt Existing Projects**: Non-destructively adopt existing repositories and GitHub Projects using `projectctl init` and `projectctl add`.
- **Map Fields**: Map existing repository labels and custom fields cleanly to neutral contracts rather than forcing disruptive renames.
- **Automate Validation**: Add reusable GitHub Actions workflows to validate project contracts on every pull request.
- **Generate Agent Guidance**: Produce neutral, deterministic instructions for AI coding assistants.

### Path 3: Full-System Operator
For operators managing large multi-repository portfolios, personal workflows, or multi-project routing:
- **Private Control Repositories**: Maintain an optional, standalone private control repository with cross-project inventories and personal agendas.
- **Multi-Project Dispatcher**: Route tasks and status updates across multiple GitHub Projects and organizations from a centralized operator profile.
- **Custom Source & Storage Providers**: Connect optional storage backends (e.g. Google Docs/Sheets for private companion records) while ensuring private URLs and IDs never touch public repositories.

---

## Architecture & Public Boundary

The repository is structured with strict public boundary separation:

```
├── cmd/
│   ├── projectctl/       # CLI executable: projectctl
│   └── projects/         # CLI executable alias: projects
├── pkg/
│   ├── project/          # [PUBLIC] Domain models, workspace logic, manifest parsers
│   └── version/          # [PUBLIC] Version metadata and runtime build inspection
├── internal/             # [PRIVATE] Encapsulated implementation details
│   ├── cli/              # CLI argument parser, formatting, subcommand handlers
│   ├── config/           # Manifest discovery, loading, and persistence logic
│   ├── git/              # Git operations wrapper and client abstractions
│   └── runner/           # Multi-repo task runner and concurrent orchestration
├── projects.go           # [PUBLIC] Root package convenience exports for Go consumers
├── .github/workflows/    # CI/CD and release automation
├── .devcontainer/        # Containerized development environment
├── CONTRIBUTING.md       # Contribution guidelines for developers and agents
├── SECURITY.md           # Security policy and privacy boundary guarantees
├── CODE_OF_CONDUCT.md    # Contributor Covenant Code of Conduct
├── Makefile              # Build, test, format, and lint automation
├── go.mod                # Go module definition: github.com/MiguelRodo/projects
└── LICENSE               # MIT License
```

### Boundary Guarantees
- **`github.com/MiguelRodo/projects`** and **`github.com/MiguelRodo/projects/pkg/*`** provide a stable, documented Go API that external tools and GUI applications can safely depend upon.
- **`github.com/MiguelRodo/projects/internal/*`** is private and cannot be imported by external packages, protecting consumers against internal refactoring.

---

## Supported Platforms & Release Policy

### Supported Platforms
- **Linux**: `amd64`, `arm64`
- **macOS (Darwin)**: `amd64` (Intel), `arm64` (Apple Silicon)
- **Windows**: `amd64`, `arm64`

### Release Policy
- **Semantic Versioning**: Releases follow [SemVer 2.0.0](https://semver.org) (`vMAJOR.MINOR.PATCH`).
- **Contract Compatibility**: Manifest and protocol contract versions are independent of tool binary versions; older contracts are normalised deterministically.
- **Automated Cross-Compilation**: Release binaries and checksums are built reproducibly via GitHub Actions and GoReleaser upon semantic version tags (`v*`).

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
# Compiled binaries will be located in ./bin/projectctl and ./bin/projects
```

---

## CLI Quickstart (`projectctl`)

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

### Common Commands

#### 1. Initialize a Workspace
```bash
projectctl init --name my-workspace
```

#### 2. Add Repositories
```bash
projectctl add https://github.com/example-org/repo-alpha.git
projectctl add https://github.com/example-org/service-b.git --name service-b --path services/service-b --branch main
```

#### 3. List Configured Repositories
```bash
projectctl list
# Format as JSON:
projectctl list --json
```

#### 4. Synchronize Workspace Repositories
```bash
projectctl sync --pull --concurrency 4
```

#### 5. Inspect Workspace Git Status
```bash
projectctl status
```

#### 6. Execute Commands in Parallel Across Repositories
```bash
projectctl exec -- git fetch --all
projectctl exec -- make test
```

---

## Go Library Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/MiguelRodo/projects"
)

func main() {
	ws := projects.NewWorkspace("demo-workspace", ".")

	err := ws.AddRepository(projects.Repository{
		Name: "core-service",
		URL:  "https://github.com/example-org/core-service.git",
		Path: "services/core-service",
	})
	if err != nil {
		log.Fatalf("failed to add repository: %v", err)
	}

	fmt.Printf("Workspace: %s, Repositories: %d\n", ws.Name, len(ws.Repositories))
}
```

---

## Development

```bash
make test       # Run unit tests with race detection & coverage
make lint       # Run golangci-lint
make fmt        # Format source code
make vet        # Run go vet
make build      # Build binaries into bin/
make all        # Run full verification pipeline
make clean      # Clean build artifacts
```

---

## Contributing & Governance

Please review our:
- [CONTRIBUTING.md](CONTRIBUTING.md) for contribution workflows and testing guidelines.
- [SECURITY.md](SECURITY.md) for vulnerability reporting and privacy boundary rules.
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for community standards.

---

## License

This project is licensed under the [MIT License](LICENSE).