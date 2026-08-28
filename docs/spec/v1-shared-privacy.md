# v1 shared privacy advertisement and companion linkage

Status: **normative pre-v1 specification**

Issue: #15

This document defines the collaborator-safe privacy policy that a participating repository may advertise in the v1 shared contract. It also defines the sole v1 association between a shared GitHub issue and an optional private companion record. It refines the [v1 conceptual model](v1-conceptual-model.md) without defining the private operator-profile wire format, a companion provider or content authoring. [V1 ordinary task interactions](v1-task-interactions.md) defines how a resolved private supplement composes with shared planning, execution, failure and recovery.

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Frozen decisions

| Question | V1 decision |
| --- | --- |
| Shared field | Optional `spec.privacy` in both `SingleProjectRepository` and `DispatcherRepository` |
| Absent field | Supports and defaults to `shareable_by_default`; private companion routing is not supported |
| Privacy modes | Exactly `shareable_by_default` and `full_github_context` |
| Supported modes | One or both exact literals, with no duplicates |
| Repository default | Explicit `defaultMode`, which must be a member of `supportedModes` |
| Effective operator choice | Private-profile data when present; never stored in the shared contract |
| No private profile | Use the repository default and perform no companion access unless a private payload makes a profile mandatory |
| Companion support | Represented only by presence of `privacy.companion` |
| Linkage advertisement | `companion.linkage.kind: github_issue_node_id` |
| Per-issue linkage | Exact observed GitHub global node ID for an `Issue`, wrapped in a typed provider-scoped key on the private side |
| Shared marker | None: no issue label, body marker, comment, Project field, URL or private-record identifier |
| Safety | Private payload never falls back to GitHub and missing private configuration never changes the selected mode |

The shared policy is repository authority. It declares supported behaviour and a default, not an acting user's preference, a provider capability result or a promise that a private provider is reachable.

## Wire representation

An explicit policy that supports both modes and private companions is:

```yaml
privacy:
  supportedModes:
    - shareable_by_default
    - full_github_context
  defaultMode: shareable_by_default
  companion:
    linkage:
      kind: github_issue_node_id
```

The object appears directly under `spec` for either repository kind. It contains exactly:

- `supportedModes`, a required non-empty set of one or both v1 mode literals;
- `defaultMode`, a required mode contained in that set;
- optional `companion`, whose presence is the one authoritative declaration that this repository supports private companion routing.

`companion` contains only the required `linkage` object. The sole v1 linkage kind is `github_issue_node_id`. There is no boolean, separate capability declaration, destination, provider account, path or credential.

An explicit shared-only policy is:

```yaml
privacy:
  supportedModes:
    - shareable_by_default
  defaultMode: shareable_by_default
```

A repository may explicitly require full GitHub context:

```yaml
privacy:
  supportedModes:
    - full_github_context
  defaultMode: full_github_context
```

That policy cannot advertise a companion because companion routing is meaningful only when `shareable_by_default` is supported.

## Absence and safe baseline

`spec.privacy` is optional so a repository that uses only the base shared workflow need not opt into a private integration surface. Its absence has one versioned v1 meaning:

```yaml
supportedModes:
  - shareable_by_default
defaultMode: shareable_by_default
# no companion advertisement
```

This is a semantic default, not a parser or schema mutation. Validation preserves that the field was absent. #34 owns whether a later canonical serialiser materialises the equivalent object and must reproduce this exact meaning rather than choose a new default.

Null is invalid. An empty object is invalid. An explicit policy never partially inherits from the absent baseline: both `supportedModes` and `defaultMode` are required when `privacy` is present.

## Privacy modes

Privacy resolution operates on an already separated request with:

- a collaborator-safe shared projection that may be proposed for GitHub; and
- either no private supplement or one private supplement that requires a private companion destination.

The caller or an earlier validated workflow owns this separation. A mode MUST NOT silently reclassify credentials, secrets or explicitly private content as shared.

### `shareable_by_default`

The shared projection may be planned for the selected GitHub issue and Project within the ordinary contract mutation rules. A private supplement MUST NOT be written to GitHub. It may be routed only when all of these are true:

1. the shared policy advertises `companion`;
2. the acting user's private profile is present and selects or accepts this mode;
3. that profile resolves one private destination under #32;
4. the shared issue's linkage identity is known and valid;
5. later capability, permission, planning and execution checks succeed.

If no private supplement exists, this mode requires no profile or companion provider.

### `full_github_context`

The request declares its complete task context as collaborator-visible shared context and may plan it for GitHub within ordinary contract authority. The mode does not declassify a private supplement. A request that combines this effective mode with a private supplement is contradictory and fails with `privacy.companion.mode-conflict`.

Credentials and provider secrets remain prohibited regardless of mode. `full_github_context` is not authority to bypass repository visibility, issue permissions, mutation allowlists or normal content validation.

## Supported mode and default rules

Mode values compare as exact ASCII literals. Order in `supportedModes` has no semantic meaning. A canonical writer orders them as `full_github_context`, then `shareable_by_default` by exact lexical order.

The following are schema conflicts:

| Diagnostic | Condition |
| --- | --- |
| `contract.schema.privacy-mode` | A mode is not one of the two v1 literals |
| `contract.schema.privacy-supported-modes` | `supportedModes` is empty, null or not an array |
| `contract.schema.unique-items` | A supported mode appears more than once |
| `contract.schema.privacy-default-unsupported` | `defaultMode` is not contained in `supportedModes` |
| `contract.schema.companion-requires-shareable-mode` | `companion` is present without `shareable_by_default` support |
| `contract.schema.companion-linkage-kind` | The linkage kind is not exactly `github_issue_node_id` |
| `contract.schema.unknown-field` | A policy contains an effective mode, destination, account, credential, per-issue ID, support boolean or any other unknown key |

There is no semantic validator rule beyond the schema for this shared policy. Profile and observed-fact conflicts belong to pure privacy resolution after both inputs have been separately validated.

## Mode resolution

Privacy mode resolution consumes the validated shared policy and one abstract private-profile state. #32 owns the private wire representation that produces this state.

1. If `spec.privacy` is absent, use the frozen absent-policy meaning above.
2. If the profile is absent, resolve the repository's `defaultMode` with source `repository-default`.
3. If a present profile selects one exact supported mode, resolve it with source `operator-profile`.
4. If a present profile selects a mode not in `supportedModes`, fail with `privacy.mode.unsupported`.
5. Never infer a mode from repository visibility, acting principal, companion support, destination availability, provider capability or previous use.

The repository default is shared policy, not a stored user choice. When a profile is present, only its mode selection represents the operator's choice. Another collaborator may use the default or select a different supported mode without changing the shared file.

## No-profile behaviour

Absence of a private profile is not itself an error.

| Shared policy and request | Result |
| --- | --- |
| No profile, no private supplement | Resolve `defaultMode`; do not access a companion provider |
| No profile, `shareable_by_default`, private supplement, no companion advertisement | `privacy.companion.unsupported` |
| No profile, `shareable_by_default`, private supplement, companion advertised | `privacy.companion.profile-required` |
| No profile, `full_github_context`, private supplement | `privacy.companion.mode-conflict` |

In every failure row, the shared projection may still be shown in a dry-run plan, but no partial apply may occur. The private supplement is neither dropped nor copied into GitHub. A complete operation fails before mutation.

## One authoritative companion advertisement

Presence of `privacy.companion` is both necessary and sufficient shared-policy evidence that private companion routing is supported. No other shared field may repeat that fact.

In particular, v1 has no:

- `companionSupported`, `privateRouting`, `enabled` or similar boolean;
- private-companion entry in `requirements.features`;
- provider capability that can turn support on when the object is absent;
- destination type or provider account in the shared contract;
- profile reference, user reference or private record reference in shared state.

Runtime capability discovery may prove that an advertised operation is currently unavailable. It cannot add support that repository policy did not advertise or remove the policy declaration from the canonical model.

## Repository-safe linkage

### Shared advertisement

The contract stores only this scheme declaration:

```yaml
companion:
  linkage:
    kind: github_issue_node_id
```

It does not store an issue node ID because the contract describes a repository topology, not one issue. It does not store a destination because destinations are operator-private.

### Canonical per-issue key

After the shared issue has been resolved and observed, the private side uses this structured key:

```yaml
kind: github_issue_node_id
authority: github.com
issueNodeId: I_SYNTHETIC_001
```

The fields have these meanings:

| Field | Rule |
| --- | --- |
| `kind` | Exact literal `github_issue_node_id` |
| `authority` | Exact literal `github.com` for the v1 GitHub provider |
| `issueNodeId` | Exact GitHub global node ID observed for a GraphQL `Issue` object |

The tuple is the linkage identifier. Implementations MUST compare every component exactly and MUST NOT trim, case-fold, decode and re-encode, hash, shorten or reconstruct the node ID from owner, repository, issue number, URL or title.

[GitHub documents](https://docs.github.com/en/graphql/guides/using-global-node-ids) global node IDs through REST as `node_id` and GraphQL as `id`, supports direct node lookup, and recommends persisting them when an integration needs to reference objects across API versions. Persisting that provider identity on the private side avoids a mutable repository-name or issue-URL selector. The linkage remains a protocol identity distinct from display text and from the private destination.

The synthetic value above is not a valid production identifier and exists only as a fixture value.

### No shared per-issue marker

V1 writes no linkage marker to the shared issue. It does not create a label, body token, comment, Project field or backlink. A private provider indexes or looks up the private record by the canonical key. This keeps the existence and location of a particular operator's companion private and avoids consuming shared mutation authority.

Different operators may store records with the same canonical shared-issue key in different private destinations. Those records do not conflict because destination selection is operator-private and none is authoritative for shared task state.

### Linkage fact failures

When companion routing is otherwise required and allowed, the observed issue identity state maps exactly as follows:

| State | Diagnostic |
| --- | --- |
| `unknown` | `privacy.linkage.issue-identity.unknown` |
| `unavailable` | `privacy.linkage.issue-identity.unavailable` |
| `forbidden` | `privacy.linkage.issue-identity.forbidden` |
| known ID resolves to no node | `privacy.linkage.issue-identity.not-found` |
| known ID resolves to a node other than GraphQL `Issue` | `privacy.linkage.issue-identity.type-mismatch` |

No linkage lookup is required when there is no private supplement. Provider discovery and transport details remain outside this specification.

## Resolution order for a private supplement

The pure policy resolver uses this order so diagnostics are stable:

1. validate and resolve the effective mode;
2. if there is no private supplement, return without companion work;
3. if the mode is not `shareable_by_default`, fail `privacy.companion.mode-conflict`;
4. if the shared policy has no `companion`, fail `privacy.companion.unsupported`;
5. if the private profile is absent, fail `privacy.companion.profile-required`;
6. resolve the profile-owned destination and provider under #32;
7. require a known, valid `Issue` global node ID and construct the exact canonical key;
8. return the shared projection and one private companion intent as separate outputs.

The resolver performs no network write. Planning must keep the shared and private intents distinct. Execution must not report the shared operation as complete when a required companion operation failed, and must never use shared GitHub as a fallback destination.

## Decision corpus

`testdata/privacy/v1/cases.json` is the machine-readable policy table. It uses an abstract profile state rather than a private wire document so #32 can define that document without changing these shared decisions.

The corpus fixes:

- absent-policy semantics;
- repository-default and profile-selected mode sources;
- unsupported choices;
- all no-profile outcomes;
- companion mode and policy gates;
- linkage fact failures and exact key construction;
- two operators resolving differently against the same shared contract;
- the invariant that no private supplement is written to GitHub on any result.

## Invariant register

| ID | Rule |
| --- | --- |
| `PV-001` | `spec.privacy` is the sole shared v1 privacy advertisement and is valid on both repository kinds. |
| `PV-002` | Absence means shareable-only, shareable default and no companion support. |
| `PV-003` | The only v1 modes are `shareable_by_default` and `full_github_context`. |
| `PV-004` | The explicit default is a member of the exact supported-mode set. |
| `PV-005` | Shared policy never stores an acting user's effective choice. |
| `PV-006` | A missing profile uses the repository default and is not an error unless private companion work is required. |
| `PV-007` | `privacy.companion` presence is the one authoritative companion-support declaration. |
| `PV-008` | Companion advertisement requires support for `shareable_by_default`. |
| `PV-009` | Shared policy stores no destination, private provider, credential, local path, user or private record identifier. |
| `PV-010` | A private supplement is never dropped, reclassified or copied to GitHub as fallback. |
| `PV-011` | `full_github_context` cannot accompany a private supplement. |
| `PV-012` | The sole linkage scheme uses the exact observed GitHub global node ID of an `Issue`. |
| `PV-013` | The canonical linkage key is typed and scoped to `github.com`. |
| `PV-014` | No per-issue linkage marker is written to shared issue or Project state. |
| `PV-015` | Companion capability evidence cannot override shared support policy. |
| `PV-016` | Different operators may resolve different supported modes and private destinations without changing shared state. |
| `PV-017` | Shared issues remain authoritative; private companion records are optional and supplemental. |
| `PV-018` | Unknown, unavailable, forbidden, not-found and wrong-type linkage facts remain distinct. |

## Explicitly deferred

This specification does not define:

- the private operator-profile file, storage location, permissions or destination references, owned by #32;
- provider credentials or companion-provider APIs;
- content authoring, classification UI or secret detection;
- provider-specific GitHub or companion mutation payloads;
- canonical normalised serialisation and privacy-safe migrations, owned by #34;
- final cross-contract conformance packaging, owned by #16;
- Go model or resolver APIs.

Later work MUST consume this policy and decision corpus rather than add a second support flag, infer an effective choice or invent a shared linkage marker.
