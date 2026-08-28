# v1 architecture boundaries

Status: **normative pre-v1 architecture decision**

This document fixes the top-level product interfaces, layers, package responsibilities and dependency direction for subsequent specification and implementation issues. Wire fields and detailed protocol semantics remain owned by their dedicated issues.

## Product decision

The pre-v1 workspace-management implementation is discarded. It has no releases or users and creates no compatibility obligation.

The v1 product has:

- a shared repository protocol that is the interoperability contract;
- one portable reference implementation and administration client named `projectctl`;
- environment-specific execution adapters, including direct provider tools, that may apply the same protocol without invoking `projectctl`;
- concise derived agent guidance for ordinary interaction;
- no `projects` alias;
- no clone, pull, status or arbitrary command runner across local repositories;
- no root-package re-export layer before the public API is stable;
- no requirement for a private control plane, Google provider, local shell or particular agent;
- no provider or CLI logic inside the canonical model.

Automated setup and consistent ordinary project interaction are co-equal primary outcomes. A comprehensive CLI is useful, but CLI completeness alone does not define product success.

## Product interfaces

| Interface | Role | Normative authority | Explicitly not |
| --- | --- | --- | --- |
| Shared protocol | Versioned repository contracts, private operator-profile contracts and their frozen semantics | Merged specifications, schemas and conformance fixtures | A particular program, prompt or provider SDK |
| Portable client and reference implementation | Validates, plans, sets up, administers and verifies through the reusable core | None beyond the shared protocol | A mandatory runtime hop for every agent |
| Environment-specific execution adapter | Converts an authorised request into conforming provider reads and writes | None beyond the shared protocol and operator request | Permission to invent defaults, routes or mutation authority |
| Derived agent guidance | Concise deterministic projection for ordinary interaction | None; every value traces to validated protocol input or explicit runtime fact | A competing configuration file or replacement specification |

`projectctl` is the portable client and reference implementation. A shell-capable agent, local operator or CI job MAY invoke it. An agent with suitable GitHub or other provider tools MAY instead implement the same pipeline directly. An explicitly configured bridge, such as a GitHub Actions workflow, MAY invoke the portable client for a bounded request. No execution path changes project meaning.

## Execution paths

### Direct provider adapter

An agent uses available provider tools directly. It reads the same shared contract and any private profile available to its operator, resolves one target, produces a complete plan, applies only authorised owned changes and verifies by readback. It does not need to materialise CLI JSON or invoke a local binary.

### Portable client

An operator, agent or CI job invokes `projectctl`. Non-interactive mutating commands are dry-run by default and require explicit apply selection. Structured CLI input and output are client interfaces, not the public protocol itself.

### Optional execution bridge

An environment that cannot perform a required provider operation MAY submit a bounded request to an explicitly configured bridge. The bridge does not become a new authority, broaden the request or hide its acting principal. A missing bridge is a capability limitation, not permission to fall back to unsafe behaviour.

## Authority layers

| Layer | Authority | Contains | Excludes |
| --- | --- | --- | --- |
| Shared repository contract | Participating repository | Version, topology, explicit target references, field/value mappings, collaborator-safe source references and supported shared policy | Effective user preferences, private destinations, local paths, credentials and private provider identifiers |
| Private operator profile | Acting user | Effective privacy choice, private companion destination, local workspace mapping, private sources and optional provider configuration | Changes to shared repository semantics or another user's profile |
| Observed provider state | Provider readback | Repository, issue, Project, field, relationship, ownership and capability facts | Guessed values or desired configuration |
| Canonical model | Pure semantics | Explicit normalised meanings of supported contracts and profiles, with their boundary preserved | Wire-format tags, network clients, CLI concerns and undiscovered provider facts |
| Resolved configuration | Runtime | The explicit operation target and effective policy for one request | Persistence as a competing authority |
| Plan | Pure semantics | Proposed operations, expected prior state and preconditions | Network mutation |
| Derived agent guidance | Projection only | Selected validated meanings and explicit runtime facts required for an audience | New defaults, policy, private leakage or authority |
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
9. **Operation lowering** turns the approved plan into typed owned changes and exact preconditions.
10. **Execution** re-reads stale-sensitive state and applies approved operations.
11. **Verification** reads back every attempted mutation and reports mismatches as failures.

Every execution adapter preserves these semantic stages even when its implementation combines transport calls or does not expose an intermediate serialisation to the operator.

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

The dependency direction is inward: CLI and provider adapters depend on public core packages. Core packages do not import CLI or concrete provider code. The package map describes the reference implementation, not a package requirement for independent conforming adapters.

## Identity and defaulting rules

- Contract references, not display titles, are identifiers.
- Local filesystem paths are operator-profile data, not shared repository-contract data.
- Repository owner kind, default branch, permissions and live field identifiers are discovered facts unless explicitly required by a contract field.
- Validators reject non-canonical identifiers rather than silently trimming or rewriting them.
- Validation never fills absent values.
- Case folding, wildcard behaviour, route precedence, fallback and multi-match behaviour require explicit protocol decisions.

## Privacy boundary

A shared contract may advertise supported privacy behaviour. The acting user's effective mode and any private destination belong to that user's private operator profile.

Private-companion support must have one authoritative representation. Capability discovery must not create a second conflicting source of truth. A direct provider adapter, portable client and execution bridge all apply the same boundary. Private content is never copied into GitHub merely because one execution path lacks the companion capability.

## Authorisation, permission and apply selection

Operator authorisation, provider permission and protocol mutation authority are independent:

- the operator request identifies what the operator authorises;
- the shared/private contracts limit what the requested workflow may mean and mutate;
- provider discovery establishes what the acting principal can currently read or write;
- the execution adapter selects whether to plan only or apply within the operator's authorisation.

For the non-interactive CLI, omission of `--apply` means plan only. An interactive agent MAY apply within an explicit current user request without asking for a redundant CLI-style confirmation. It MUST still make the exact target and owned change determinable, stop on ambiguity or missing authority and return verified readback. Broad, destructive or materially under-specified requests remain unresolved rather than being widened by the adapter.

## Mutation safety

Every mutating execution path must:

1. produce a complete plan before writes;
2. identify one exact target and the owned change set;
3. carry expected prior state or equivalent preconditions;
4. re-read stale-sensitive state before mutation;
5. refuse an unsafe stale write;
6. preserve unrelated live content;
7. read back each completed write;
8. report verification failure explicitly.

A client MAY present the plan as an explicit dry-run artefact. An interactive adapter MAY proceed under explicit operator authorisation, but it does not skip planning, stale protection or verification.

## Specialised specifications

Dedicated normative specifications own these decisions. Linked specifications are settled; unlinked issue references remain deferred:

- private operator-profile wire representation: [v1 operator profile](../spec/v1-operator-profile.md);
- label, routing-label and label-based sub-project declarations: [v1 labels and sub-projects](../spec/v1-labels-and-subprojects.md);
- semantic project dimensions and provider bindings: [v1 project dimensions](../spec/v1-project-dimensions.md);
- automated setup and manual-action outcomes: [v1 setup outcomes](../spec/v1-setup-outcomes.md);
- ordinary task-interaction requests and mutation scopes: [v1 ordinary task interactions](../spec/v1-task-interactions.md);
- contract-version vocabulary, presence and migration rules: #34;
- concise derived guidance projection and renderers: #51;
- final structured output schemas and exit codes: #26.

No implementation or execution adapter may contradict a settled specification or decide a remaining deferred item implicitly.
