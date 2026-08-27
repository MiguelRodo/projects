# v1 private operator profile and companion-provider contract

Status: **normative pre-v1 specification**

Issue: #32

This document defines the private operator-profile wire format, contract selection, provider-neutral private bindings and the abstract companion-provider interface used with the shared privacy policy in [v1-shared-privacy.md](v1-shared-privacy.md). It refines the [v1 conceptual model](v1-conceptual-model.md) and the execution-interface boundary in [v1-boundaries.md](../architecture/v1-boundaries.md).

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Purpose

One participating repository has one shared contract, but different operators may have different privacy choices, private destinations, source bindings and local workspaces. Those values are not collaborator-safe and do not belong in the repository contract.

The operator profile is a private canonical input. It is not:

- a shared project authority;
- a credential store;
- a provider capability claim;
- a required input for ordinary collaborators;
- a local-shell or `projectctl`-specific configuration file;
- permission to override shared topology, field semantics, routing or mutation authority.

A direct-provider agent, the portable client and an explicit execution bridge consume the same validated document and produce equivalent resolution results. The mechanism used to acquire the private bytes is outside the wire contract.

## Frozen decisions

| Question | V1 decision |
| --- | --- |
| Schema | Draft 2020-12 at `schemas/v1/operator-profile.schema.json` |
| Version discriminator | `apiVersion: projectctl/v1` |
| Kind | `OperatorProfile` |
| Text form | UTF-8 YAML 1.2.2 Core Schema restricted to the JSON-compatible profile used by the shared contract |
| Profile identity | None in the document; possession and storage authority identify the acting operator |
| Shared-contract selector | Participating GitHub repository locator plus exact shared `spec.ref` |
| Selector fallback | None: no wildcard, default binding, inheritance, title lookup or current-directory inference |
| Effective privacy choice | Optional `privacy.mode` on the exact contract binding |
| Missing privacy choice | Behaves as no private privacy selection and uses the shared repository default |
| Companion destination | Optional private reference, valid only with `shareable_by_default` |
| Providers | Logical adapter binding plus optional external configuration and credential references |
| Credentials | Never inline; `credentialRef` is only a logical handle resolved outside the profile |
| Provider resource identity | Private opaque `resourceRef`, compared exactly |
| Workspace paths | Typed absolute local paths only; never placed in shared state |
| Source hooks | Exact shared-source bindings plus separate private-only source references |
| Agenda and registry | Read-only extension bindings |
| Cache | Pull-only extension binding |
| Defaults | No wire defaults; every top-level collection is explicit |
| Profile acquisition | Explicit private input; never discovered from a participating repository |
| Companion linkage | Exact private-side `github_issue_node_id` key fixed by #15 |

## Document acquisition and storage

The canonical input to profile parsing is a byte sequence plus an adapter-private source identity. No absolute path, environment variable, account product or API is part of the protocol.

An adapter MAY acquire those bytes from:

- an explicitly selected local file;
- acting-user-only application storage;
- an explicitly configured private provider;
- an operator control plane.

It MUST NOT:

- search a participating repository or its parents for a profile;
- derive a profile location from shared contract content;
- accept a profile placed at `.projectctl/` in a participating repository;
- transmit the profile to GitHub or another shared task surface;
- turn absence into an implicit globally named profile.

When backed by a filesystem, the conventional filename is `operator-profile.yaml`, but there is no protocol-level absolute path. A filesystem adapter MUST require an absolute explicitly selected path or an implementation-owned user configuration location. It MUST reject a path inside the selected participating repository's tracked tree with `profile.storage.participating-repository`. It MUST reject storage that is readable or writable by principals beyond the acting user or explicitly authorised service principal with `profile.storage.permissions-unsafe`.

Equivalent access controls apply to non-filesystem stores. Public repositories may contain only the schema, specification and clearly synthetic fixtures.

Actual profiles may reveal private destinations, source locations and local paths. They MUST NOT be committed to participating repositories even though the schema forbids inline credentials.

## Encoding and parser restrictions

The document uses the same parser boundary as the shared v1 contract:

- UTF-8 without a byte-order mark;
- YAML 1.2.2 Core Schema;
- one document;
- mappings with unique scalar string keys;
- sequences, strings, booleans, numbers and mappings only where the schema permits them;
- no aliases, anchors, merge keys, tags, duplicate keys, non-string mapping keys or cyclic graphs;
- no null values;
- no unknown keys at any object boundary.

JSON is accepted only when it represents exactly the same data model. Parsing preserves supplied presence. Parsing and schema validation do not materialise defaults.

## Complete shape

This synthetic profile exercises every v1 binding:

```yaml
apiVersion: projectctl/v1
kind: OperatorProfile
spec:
  providers:
    private-records:
      adapterRef: synthetic-record-store
      configurationRef: private-records-config
      credentialRef: private-records-credential
    private-resources:
      adapterRef: synthetic-resource-store
      configurationRef: private-resources-config
  companionDestinations:
    alpha-companions:
      providerRef: private-records
      resourceRef: synthetic/companion-container-alpha
  workspaces:
    alpha-workspace:
      path:
        kind: posix_absolute
        value: /synthetic/work/alpha
  integrations:
    alpha-roadmap:
      kind: source
      providerRef: private-resources
      resourceRef: synthetic/source-roadmap-alpha
      access: read_only
    alpha-private-notes:
      kind: source
      providerRef: private-resources
      resourceRef: synthetic/source-notes-alpha
      access: read_only
    alpha-agenda:
      kind: agenda
      providerRef: private-resources
      resourceRef: synthetic/agenda-alpha
      access: read_only
    alpha-registry:
      kind: registry
      providerRef: private-resources
      resourceRef: synthetic/registry-alpha
      access: read_only
    alpha-cache:
      kind: cache
      providerRef: private-resources
      resourceRef: synthetic/cache-alpha
      access: pull_only
  contractBindings:
    - contract:
        authority: github.com
        owner:
          kind: organization
          login: example-org
        repository: alpha-config
        ref: alpha
      privacy:
        mode: shareable_by_default
        companionDestinationRef: alpha-companions
      workspaceRef: alpha-workspace
      sources:
        shared:
          roadmap: alpha-roadmap
        private:
          - alpha-private-notes
      extensions:
        agendaRef: alpha-agenda
        registryRef: alpha-registry
        cacheRef: alpha-cache
```

The five `spec` collections are required even when a map is empty. `contractBindings` contains at least one entry. A binding is invalid when it has no privacy selection, workspace, source or extension use.

## Reference syntax and namespaces

Profile-owned symbolic references use the shared contract-reference syntax:

```text
^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$
```

They are 1 to 63 ASCII characters and compare exactly. The namespaces are independent:

| Location | Namespace |
| --- | --- |
| `spec.providers` keys | Provider bindings |
| `spec.companionDestinations` keys | Companion destinations |
| `spec.workspaces` keys | Local workspaces |
| `spec.integrations` keys | Source and extension resources |
| `contract.ref` | Shared contract reference within the selected participating repository |
| `sources.shared` keys | Shared `spec.sourceRefs` in the selected contract |

Provider `resourceRef` values are private opaque strings, 1 to 512 code points, with no line break or leading/trailing whitespace. They compare exactly. An adapter MUST NOT trim, case-fold, Unicode-normalise, parse as a URL unless its selected provider contract says so, or substitute a display name.

`adapterRef`, `configurationRef` and `credentialRef` are logical references, not inline configuration or secret values. Their external resolution is adapter-specific. The schema deliberately offers no token, password, key, header or arbitrary configuration object.

A `resourceRef` is an identifier, not authentication material. It MUST NOT contain an access token, password, private key, signed request or other reusable secret. Schema validation cannot recognise every secret shape, so adapters and review tooling retain this content-safety obligation.

## Shared-contract selection

Each `contractBindings[].contract` contains:

- `authority`, exactly `github.com` in v1;
- explicit owner `kind` and `login`;
- participating repository name;
- exact shared contract `spec.ref`.

The selector identifies the repository that contains `.projectctl/project.yaml`. It does not identify the issue store or a target Project, which may be elsewhere.

Selection uses:

- exact authority;
- exact owner kind;
- ASCII case-insensitive owner login and repository name;
- exact case-sensitive contract reference.

Supplied spelling is preserved for provider discovery and diagnostics. Only ASCII `A` to `Z` are folded for locator equality. An implementation MUST NOT use locale-sensitive comparison.

There is no wildcard, pattern, owner-wide default, first entry, current checkout or title-based fallback. Zero matches means this profile has no binding for the selected contract. That is not an error. More than one semantically equal selector is invalid with `profile.semantic.duplicate-contract-binding`.

## Privacy selection

`privacy` is optional on a contract binding.

When absent, the privacy resolver receives the abstract profile state `absent` and applies the shared `defaultMode` under #15. The presence of a profile document or a non-privacy workspace/source binding does not itself create a privacy choice.

When present:

- `mode` is exactly `shareable_by_default` or `full_github_context`;
- the selected mode is authoritative for this operator only;
- the selected mode must be supported by the shared policy or resolution fails `privacy.mode.unsupported`;
- `companionDestinationRef` is optional with `shareable_by_default`;
- `companionDestinationRef` is prohibited with `full_github_context`.

A shareable selection without a destination is valid for requests with no private supplement. If a private supplement later requires companion routing, resolution fails `privacy.companion.destination.required`. A destination never changes the mode and never turns on shared companion support.

The profile cannot override the shared default, supported mode set, companion advertisement, issue store, target Project, routing rules, fields or mutation scopes.

## Provider and destination bindings

Each provider binding contains required `adapterRef` and optional `configurationRef` and `credentialRef`. These values are desired private configuration, not evidence that the adapter is installed, configured, authenticated, reachable or authorised. Runtime discovery preserves those states separately.

Each companion destination contains:

- `providerRef`, which must resolve to one declared provider binding;
- `resourceRef`, the exact provider-private destination selector.

Contract bindings refer to destinations only by their profile reference. A shared contract, issue, Project item or generated shared guidance MUST NOT contain the destination reference or resource selector.

Profile semantic validation rejects dangling destination and provider references before privacy resolution. Runtime destination discovery may still report unknown, unavailable, forbidden, not found or ambiguous state.

## Companion-provider interface

A conforming companion adapter supplies provider-neutral equivalents of these logical operations:

| Operation | Input | Output | Required property |
| --- | --- | --- | --- |
| Resolve provider | Valid provider binding and acting principal | Exact provider/configuration state | No capability or credential inference |
| Resolve destination | Resolved provider plus exact `resourceRef` | One destination identity or typed fact state | No first-match or title fallback |
| Find companion | Destination plus canonical linkage key | Zero, one or ambiguous matching records | Exact comparison of every key component |
| Plan companion | Private supplement plus zero/one observed record | Typed create/update delta and expected prior state | Private fields only |
| Apply companion | Approved operation and fresh precondition | Provider mutation result | Refuse stale or broadened writes |
| Read back | Exact destination and resulting record identity | Observed private record | Required before success |

The interface does not require Go, a local process, CLI JSON or one transport. It requires equivalent meanings and failure states.

The canonical linkage key remains:

```yaml
kind: github_issue_node_id
authority: github.com
issueNodeId: I_SYNTHETIC_001
```

The private provider indexes or looks up this exact tuple. It MUST NOT derive linkage from a repository name, issue number, URL, title, shared marker or private record title. More than one record for the tuple is `privacy.companion.record.ambiguous`; an adapter does not choose one.

Companion planning and execution preserve the #46 safety pipeline. A successful provider mutation response is not success until readback matches. Failure never copies the private supplement to GitHub.

## Workspaces

`workspaces` maps a profile reference to one typed absolute path:

- `posix_absolute`;
- `windows_drive_absolute`;
- `windows_unc_absolute`.

Relative paths, `~`, environment variables and URI inference are invalid. The path value is preserved exactly and is not shared contract identity. Before local use, an adapter MUST reject `.` or `..` path segments, an operating-system filesystem root and any path whose observed type is not a directory. It MUST distinguish not found, unavailable and forbidden observations.

A contract binding may omit `workspaceRef`. A direct-provider adapter that does not need a checkout ignores workspace bindings. An operation that explicitly requires a workspace fails with a typed unavailable result rather than asking another adapter to run `projectctl` or treating the current directory as the workspace.

## Sources and private extension points

Each integration contains:

- `kind`, one of `source`, `agenda`, `registry` or `cache`;
- a declared `providerRef`;
- exact private `resourceRef`;
- fixed v1 access direction.

`source`, `agenda` and `registry` are `read_only`. A `cache` is `pull_only`: it may receive a projection from its source of truth, but the cache never pushes authority back. V1 profile integrations do not grant generic write authority.

On a contract binding:

- `sources.shared` maps an exact shared `spec.sourceRefs` hook to a source integration;
- `sources.private` lists additional operator-private source integrations;
- optional `extensions.agendaRef`, `registryRef` and `cacheRef` each select an integration of the corresponding kind.

`sources.shared` does not change shared contract meaning. A key absent from the selected shared contract is invalid with `profile.contract.source-ref-unknown`. Private-only sources remain supplemental and cannot replace GitHub issues or Projects as shared task authority.

The profile defines only provider-neutral binding points. Source catalogue structure, attention rules, agenda contents and cache projection formats remain owned by their dedicated integration work.

## Semantic validation

Schema-valid profiles undergo pure semantic validation before use.

| Diagnostic | Condition |
| --- | --- |
| `profile.semantic.duplicate-contract-binding` | Two selectors are equal under the frozen contract-selector comparison |
| `profile.semantic.unknown-provider-ref` | A destination or integration refers to no declared provider |
| `profile.semantic.unknown-destination-ref` | A privacy selection refers to no declared companion destination |
| `profile.semantic.unknown-workspace-ref` | A binding refers to no declared workspace |
| `profile.semantic.unknown-integration-ref` | A source or extension refers to no declared integration |
| `profile.semantic.integration-kind` | A source or typed extension refers to an integration of the wrong kind |
| `profile.semantic.local-path-segment` | A path contains `.` or `..` as a complete segment |
| `profile.semantic.local-path-root` | A workspace denotes an operating-system filesystem root |
| `profile.contract.source-ref-unknown` | A `sources.shared` key is not in the selected shared contract's `sourceRefs` |

The first eight rules validate the profile alone. The final rule validates a profile binding against the separately validated selected shared contract. Validation does not discover a provider, mutate the profile or fill missing privacy choices.

## Deterministic privacy and destination resolution

The combined resolver follows this order:

1. validate the shared contract and private profile separately;
2. select the exact contract binding;
3. extract either no privacy choice or one selected mode;
4. resolve the effective mode under #15;
5. if there is no private supplement, return without provider work;
6. require `shareable_by_default`;
7. require the shared companion advertisement;
8. require a matching profile privacy selection;
9. require its `companionDestinationRef`;
10. resolve the declared destination and provider;
11. require the exact observed GitHub `Issue` node ID;
12. return separate shared and private intents.

Steps 4, 6, 7 and 11 retain the exact #15 diagnostics and ordering. #32 adds step 9 and the destination/provider states in step 10:

| State | Diagnostic |
| --- | --- |
| No selected destination | `privacy.companion.destination.required` |
| Provider state unknown | `privacy.companion.provider.unknown` |
| Provider unavailable | `privacy.companion.provider.unavailable` |
| Provider access forbidden | `privacy.companion.provider.forbidden` |
| Adapter lacks the companion interface | `privacy.companion.provider.unsupported` |
| Destination state unknown | `privacy.companion.destination.unknown` |
| Destination unavailable | `privacy.companion.destination.unavailable` |
| Destination access forbidden | `privacy.companion.destination.forbidden` |
| Destination not found | `privacy.companion.destination.not-found` |
| Destination selector ambiguous | `privacy.companion.destination.ambiguous` |

Configuration and credential failures are reported as provider failures with structured cause details, but do not change these top-level policy results. An adapter MUST NOT try another provider or destination.

## Equivalent multi-user result

`testdata/operator-profiles/v1/valid/user-shareable.yaml` and `user-full-context.yaml` select different supported modes for the same synthetic shared contract.

- The first resolves `shareable_by_default` and the private destination `alpha-companions`.
- The second resolves `full_github_context` and has no companion destination.

Neither profile changes the shared file. With no profile, the repository default applies. With no private supplement, provider availability is irrelevant in every case.

## Decision corpus

The machine-readable corpus consists of:

- valid profile fixtures;
- structurally invalid fixtures with expected schema diagnostics;
- semantic invalid cases in `testdata/operator-profiles/v1/semantic-cases.json`;
- exact profile/privacy/provider resolution rows in `testdata/operator-profiles/v1/resolution/cases.json`.

Every public value is synthetic. Conforming implementations produce equivalent status, mode source, destination result and diagnostic without requiring one execution adapter to invoke another.

## Invariant register

| ID | Rule |
| --- | --- |
| `OP-001` | The operator profile is private canonical input and never shared project authority. |
| `OP-002` | `projectctl/v1` plus `OperatorProfile` is the sole v1 profile discriminator. |
| `OP-003` | Profile acquisition is explicit and independent of participating-repository discovery. |
| `OP-004` | A contract binding uses the participating repository locator plus exact shared contract reference. |
| `OP-005` | Selectors have no wildcard, default, inheritance, title or current-directory fallback. |
| `OP-006` | Missing profile or privacy selection uses the shared repository default. |
| `OP-007` | An effective private privacy choice has one authoritative location. |
| `OP-008` | A full-context selection cannot name a companion destination. |
| `OP-009` | A destination cannot enable shared companion support or change privacy mode. |
| `OP-010` | Providers are private bindings, not capability or permission claims. |
| `OP-011` | Credentials are external logical references and are never stored inline. |
| `OP-012` | Provider resource references compare exactly and never fall back to display text. |
| `OP-013` | Local workspace paths are typed, absolute and never shared. |
| `OP-014` | Source, agenda and registry bindings are read-only; cache direction is pull-only. |
| `OP-015` | Shared source hooks and private-only sources remain distinguishable. |
| `OP-016` | Companion lookup uses the exact private-side GitHub issue node key fixed by #15. |
| `OP-017` | Companion ambiguity, stale state and readback mismatch fail closed. |
| `OP-018` | Private content never falls back to GitHub. |
| `OP-019` | Direct-provider, portable-client and bridge adapters resolve the same inputs identically. |
| `OP-020` | A workspace or provider is not required for an operation that does not use it. |

## Explicitly deferred

This specification does not define:

- a Google-specific provider or any real operator configuration;
- a default `projectctl` command-line path or environment variable;
- Go profile types or provider interfaces;
- provider transport, authentication or credential-store implementation;
- source catalogue, agenda, registry or cache payload formats;
- companion content fields or authoring UI;
- canonical normalised serialisation and profile migrations, owned by #34;
- final cross-contract conformance packaging, owned by #16.

Later work MUST consume this profile and decision corpus rather than add an effective privacy value to shared state, embed credentials, infer a local workspace or make the portable client a required hop.
