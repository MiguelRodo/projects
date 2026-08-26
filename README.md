# `projects` / `projectctl`

`projects` is the public home of a portable, agent-neutral protocol and Go administration core for GitHub issues and GitHub Projects.

## Status

This repository is in **pre-v1 architecture and contract design**.

There is currently no supported CLI, Go library, contract format or release to install. The discarded pre-v1 workspace manager had no users or releases and is not part of the product. Implementation begins only after the normative protocol issues and conformance fixtures are merged.

The authoritative roadmap and execution order are in [issue #1](../../issues/1).

## Product boundary

| Layer | Purpose |
| --- | --- |
| Shared repository contract | Collaborator-safe, versioned project topology and mappings |
| Private operator profile | Per-user privacy choices, private destinations, local mappings and optional provider configuration |
| Portable Go core | Parsing, normalisation, validation, resolution, planning, execution and reporting |
| Provider adapters | GitHub first, with optional private or external providers behind interfaces |
| `projectctl` | Thin non-interactive CLI over the reusable core |
| Future application | Graphical interface over the same core after v1 is proven |

GitHub issues and GitHub Projects are the authoritative shared task surface. Private control planes are optional consumers and must never become a prerequisite for ordinary collaborators.

## Participation paths

### Collaborator or external agent

Work through normal GitHub issues and pull requests. No private operator configuration or external provider access is required.

### Repository maintainer

Adopt or bootstrap a repository using the shared contract and `projectctl`, review plans before writes, and verify resulting GitHub state.

### Full-system operator

Add a private operator profile, private source routing, agenda rules or optional providers without placing private destinations or account-specific identifiers in shared files.

## Architecture

The normative layer, package and dependency boundaries are defined in [docs/architecture/v1-boundaries.md](docs/architecture/v1-boundaries.md).

The normative terminology, authority, identity and processing model is defined in [docs/spec/v1-conceptual-model.md](docs/spec/v1-conceptual-model.md).

The shared single-Project wire contract and schema are defined in [docs/spec/v1-single-project-contract.md](docs/spec/v1-single-project-contract.md).

The intended executable is `projectctl`. There is no `projects` compatibility binary and no multi-repository clone/sync workspace-manager scope.

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Work is issue-driven, one issue per PR, and foundational semantics must be merged before dependent implementation begins.

Security and privacy reports follow [SECURITY.md](SECURITY.md).

## Licence

This project is licensed under the [MIT Licence](LICENSE).
