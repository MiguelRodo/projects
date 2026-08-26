# Contributing to `projects`

Thank you for contributing. The project is currently establishing its v1 protocol and implementation architecture.

## Start with the authoritative issue

Before working:

1. Read [issue #1](../../issues/1) for the execution order.
2. Read [the v1 architecture boundaries](docs/architecture/v1-boundaries.md).
3. Select one open leaf issue whose dependencies are complete.
4. Treat the exact normative documents, fixtures and interfaces named by that issue as authoritative.

Do not implement a later layer because it seems convenient. If a required protocol decision is absent or contradictory, stop and describe the conflict in the PR.

## Pull-request discipline

- Work on one issue per pull request.
- Create a branch and open a PR against `main`; do not push directly to `main`.
- Change only the files and packages permitted by the issue.
- Do not include unrelated cleanup or anticipatory abstractions.
- Keep the PR description aligned with the actual diff.
- Use `Refs #N` until every completion condition is met; use `Closes #N` only when the PR genuinely completes the issue.

## Public boundary

Code, schemas, documentation and fixtures must remain agent-neutral and operator-neutral.

Do not include:

- private repository names;
- private Drive, document or spreadsheet identifiers;
- credentials, tokens or secret material;
- personal names or account-specific configuration;
- private source locations or private destination URLs.

Use clearly synthetic fixtures and local test servers.

## Testing and external systems

Tests must not mutate live GitHub issues, Projects, repositories or private systems. Use fixtures, recorded snapshots, mocks or `httptest`.

Each implementation issue defines the exact verification commands required for its scope. Do not invent a repository-wide build command before the Go foundation issue supplies one.

## Public and internal Go packages

When Go implementation begins, follow the approved package map:

- provider-neutral reusable packages under `pkg/`;
- GitHub API implementation under `internal/githubapi`;
- CLI adaptation under `internal/cli`;
- one executable at `cmd/projectctl`.

Do not add a root-package compatibility API or a second executable unless the architecture issue is explicitly revised first.
