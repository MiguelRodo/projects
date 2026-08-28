# `projects` / `projectctl`

`projects` is the public home of a portable, agent-neutral protocol for automated GitHub project setup and consistent ordinary interaction, plus a Go administration core and `projectctl` reference client.

## Status

This repository is in **pre-v1 architecture and contract design**.

There is currently no supported CLI, Go library, contract format or release to install. The discarded pre-v1 workspace manager had no users or releases and is not part of the product. Implementation begins only after the normative protocol issues and conformance fixtures are merged.

The authoritative roadmap and execution order are in [issue #1](../../issues/1).

## Product boundary

| Interface or layer | Purpose |
| --- | --- |
| Shared protocol | Collaborator-safe, versioned project topology, structure, mappings and interaction semantics |
| Private operator profile | Per-user privacy choices, private destinations, local mappings and optional provider configuration |
| Derived agent guidance | Concise deterministic project instructions derived from validated contracts |
| Environment-specific execution adapter | Direct provider tools, the portable client or an explicit optional bridge |
| Portable Go core | Parsing, normalisation, validation, resolution, planning, execution and reporting |
| Provider adapters | GitHub first, with optional private or external providers behind interfaces |
| `projectctl` | Portable reference implementation and setup/administration client over the reusable core |
| Future application | Graphical interface over the same core after v1 is proven |

GitHub issues and GitHub Projects are the authoritative shared task surface. Private control planes are optional consumers and must never become a prerequisite for ordinary collaborators. `projectctl` may be comprehensive, but it is not a mandatory runtime dependency for an agent that can apply the same protocol safely through provider tools.

The primary v1 product outcomes are:

1. automate creation or adoption of the declared repository/Project structure as far as provider APIs allow, with unsupported manual actions reported explicitly; and
2. let capable agents and operators read and manage ordinary project work consistently across collaborator/private, single/multi-Project and central/repository-backed scenarios.

## Participation paths

### Collaborator or external agent

Work through normal GitHub issues and pull requests. A capable agent may use suitable GitHub/provider tools directly from the shared contract and concise generated guidance. No private operator configuration, external provider access or local binary is required.

### Repository maintainer

Adopt or bootstrap a repository using the shared contract and `projectctl`, review plans before writes, and verify resulting GitHub state. Provider surfaces without a supported automation API are returned as explicit manual actions rather than false successes.

### Full-system operator

Add a private operator profile, private source routing, agenda rules or optional providers without placing private destinations or account-specific identifiers in shared files.

## Architecture

The normative layer, package and dependency boundaries are defined in [docs/architecture/v1-boundaries.md](docs/architecture/v1-boundaries.md).

The normative terminology, authority, identity and processing model is defined in [docs/spec/v1-conceptual-model.md](docs/spec/v1-conceptual-model.md).

The shared single-Project wire contract and schema are defined in [docs/spec/v1-single-project-contract.md](docs/spec/v1-single-project-contract.md).

The dispatcher extension and deterministic multi-Project routing semantics are defined in [docs/spec/v1-routing.md](docs/spec/v1-routing.md).

Collaborator-safe labels, overall-Project routing labels and optional label-based sub-projects are defined in [docs/spec/v1-labels-and-subprojects.md](docs/spec/v1-labels-and-subprojects.md).

Provider-neutral task dimensions and target-specific Project, Issue Type and Issue Field bindings are defined in [docs/spec/v1-project-dimensions.md](docs/spec/v1-project-dimensions.md).

The shared privacy advertisement and repository-safe companion linkage are defined in [docs/spec/v1-shared-privacy.md](docs/spec/v1-shared-privacy.md).

The private operator-profile wire format, provider-neutral bindings and companion-provider contract are defined in [docs/spec/v1-operator-profile.md](docs/spec/v1-operator-profile.md).

The intended executable is `projectctl`. There is no `projects` compatibility binary and no multi-repository clone/sync workspace-manager scope. Direct-provider agents and optional execution bridges remain separate adapters over the same protocol.

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Work is issue-driven, one issue per PR, and foundational semantics must be merged before dependent implementation begins.

Security and privacy reports follow [SECURITY.md](SECURITY.md).

## Licence

This project is licensed under the [MIT Licence](LICENSE).
