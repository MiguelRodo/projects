# v1 architecture boundaries

Status: **normative pre-v1 architecture decision**

This document fixes the top-level layers, package responsibilities and dependency direction for subsequent specification and implementation issues. Wire fields and detailed protocol semantics remain owned by their dedicated issues.

## Product decision

The pre-v1 workspace-management implementation is discarded. It has no releases or users and creates no compatibility obligation.

The v1 product has:

- one executable named `projectctl`;
- no `projects` alias;
- no clone, pull, status or arbitrary command runner across local repositories;
- no root-package re-export layer before the public API is stable;
- no requirement for a private control plane or Google provider;
- no provider or CLI logic inside the canonical model.

## Authority layers

| Layer | Authority | Contains | Excludes |
| --- | --- | --- | --- |
| Shared repository contract | Participating repository | Version, topology, explicit target references, field/value mappings, collaborator-safe source references and supported shared policy | Effective user preferences, private destinations, local paths, credentials and private provider identifiers |
| Private operator profile | Acting user | Effective privacy choice, private companion destination, local workspace mapping, private sources and optional provider configuration | Changes to shared repository semantics or another user's profile |
| Observed provider state | Provider readback | Repository, issue, Project, field, relationship, ownership and capability facts | Guessed values or desired configuration |
| Canonical model | Pure core | Explicit normalised meanings of supported contracts and profiles, with their boundary preserved | Wire-format tags, network clients, CLI concerns and undiscovered provider facts |
| Resolved configuration | Runtime | The explicit operation target and effective policy for one request | Persistence as a competing authority |
| Plan | Pure core | Proposed operations, expected prior state and preconditions | Network mutation |
| Verification report | Readback result | Applied, skipped and failed operations with observed evidence | Unverified success claims |

## Behavioural pipeline

1. **Parse** preserves supplied versus absent values and reports syntax failures.
2. **Schema validation** checks the selected versioned wire shape.
3. **Normalisation** returns a new canonical value and applies only version-defined defaults.
4. **Semantic validation** is pure and never mutates input.
5. **Discovery** obtains provider facts and capabilities.
6. **Resolution** selects explicit targets using frozen routing semantics.
7. **Snapshot** records the relevant observed state deterministically.
8. **Planning** computes a deterministic proposed delta without writes.
9. **Execution** re-reads stale-sensitive state and applies approved operations.
10. **Verification** reads back every attempted mutation and reports mismatches as failures.

## Package map

The reusable Go core will be divided as follows:

| Package | Responsibility |
| --- | --- |
| `pkg/model` | Provider-neutral canonical types and invariants, without JSON/YAML tags |
| `pkg/contract` | Versioned wire types, parsing, schema validation and normalisation into `model` |
| `pkg/resolve` | Pure target and effective-policy resolution |
| `pkg/state` | Deterministic provider-neutral observed snapshots |
| `pkg/plan` | Pure desired-versus-observed comparison, operations and preconditions |
| `pkg/execute` | Execution orchestration through provider interfaces, stale checks and readback |
| `pkg/report` | Stable structured results, human rendering and exit-code mapping |
| `pkg/provider` | Provider interfaces and capability vocabulary |
| `internal/githubapi` | GitHub REST/GraphQL transport, authentication and provider implementation |
| `internal/cli` | Argument parsing, command wiring and terminal I/O |
| `cmd/projectctl` | Minimal executable entry point |

The dependency direction is inward: CLI and provider adapters depend on public core packages. Core packages do not import CLI or concrete provider code.

## Identity and defaulting rules

- Contract references, not display titles, are identifiers.
- Local filesystem paths are operator-profile data, not shared repository-contract data.
- Repository owner kind, default branch, permissions and live field identifiers are discovered facts unless explicitly required by a contract field.
- Validators reject non-canonical identifiers rather than silently trimming or rewriting them.
- Validation never fills absent values.
- Case folding, wildcard behaviour, route precedence, fallback and multi-match behaviour require explicit protocol decisions.

## Privacy boundary

A shared contract may advertise supported privacy behaviour. The acting user's effective mode and any private destination belong to that user's private operator profile.

Private-companion support must have one authoritative representation. Capability discovery must not create a second conflicting source of truth.

## Mutation safety

Every mutating command must:

1. provide a complete dry-run plan;
2. carry expected prior state or equivalent preconditions;
3. re-read stale-sensitive state before mutation;
4. refuse an unsafe stale write;
5. preserve unrelated live content;
6. read back each completed write;
7. report verification failure explicitly.

## Deferred decisions

The following remain deliberately deferred to dedicated normative issues:

- exact YAML or JSON filenames and schemas;
- contract-version vocabulary and migration rules;
- routing key syntax and matching semantics;
- privacy-policy wire representation;
- private operator-profile wire representation;
- stable linkage identifier representation;
- structured output schemas and exit codes.

No implementation issue may decide these implicitly.
