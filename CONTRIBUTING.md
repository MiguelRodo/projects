# Contributing to `projects`

Thank you for your interest in contributing to `projects`! We welcome contributions from community members, maintainers, and automated agents.

---

## Core Principles

1. **GitHub is the Shared Authority**: Public contracts and repository interactions must rely solely on standard GitHub capabilities (issues, pull requests, projects, labels). Private control planes and personal agendas are optional companion systems and must never be a hard dependency for contributing.
2. **Clean Public Boundary**: All domain abstractions, contracts, and interfaces in `pkg/` and root package must remain strictly neutral and reusable across diverse development workflows. Internal implementation details belong in `internal/`.
3. **Agent & Operator Neutrality**: Code, documentation, schemas, and fixtures must avoid hard-coded personal identifiers, specific accounts, proprietary tools, or operator-specific conventions.
4. **Safety & Non-Destructive Operations**: Workspace operations, syncs, and adoption workflows must be non-destructive and verifiable.

---

## Development Setup

### Prerequisites
- **Go 1.22+** (Go 1.24+ recommended)
- **Git 2.30+**
- **golangci-lint** (`latest` or `v2.x`)

### Building and Testing

```bash
# Clone the repository
git clone https://github.com/MiguelRodo/projects.git
cd projects

# Format code
make fmt

# Run static analysis / linters
make lint

# Run all unit tests with race detection and coverage
make test

# Build local binaries to ./bin/
make build

# Run the complete test and build pipeline
make all
```

---

## Development Guidelines

### 1. Public vs. Internal Packages
- Public API surfaces belong in `pkg/` or the root package (`projects.go`).
- Internal implementation details (CLI parsing, execution runners, git process calls, configuration IO) must reside in `internal/`.
- Exported functions, types, and constants must be thoroughly documented with godoc comments.

### 2. Testing & Test Doubles
- Use unit tests for all new features and bug fixes.
- Tests must pass with the `-race` detector enabled.
- **Do not make live network calls or mutate remote GitHub resources during tests.** Use fixtures, mocks (`git.MockClient`), or `httptest` servers.
- Use neutral fixture data:
  - Organizations: `example-org`, `acme-corp`
  - Repositories: `repo-alpha`, `service-a`, `core-api`
  - Users: `alice`, `bob`, `contributor`
  - URLs: `https://github.com/example-org/repo-alpha.git`

### 3. Commit Messages & Pull Requests
- Use clear, descriptive commit messages.
- Reference relevant GitHub issue numbers in commit descriptions and PR titles (e.g. `feat: add manifest validation (#11)`).
- Keep pull requests focused on a single issue or scope.

---

## Submitting a Pull Request

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
2. Implement your changes, add tests, and verify:
   ```bash
   make all
   ```
3. Commit and push your branch:
   ```bash
   git push origin feat/your-feature-name
   ```
4. Open a Pull Request on GitHub against the `main` branch.
5. Ensure CI checks pass on all supported Go matrix versions.
