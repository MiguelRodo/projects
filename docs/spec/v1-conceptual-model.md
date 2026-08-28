# v1 conceptual model

Status: **normative pre-v1 specification**

Issue: #12; architecture clarification: #46

This document fixes the vocabulary, authority model, identity model and processing boundaries used by every v1 contract and implementation. Later specifications may choose wire representations and exact algorithms only within these boundaries.

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Purpose

The system allows a participating repository to describe collaborator-safe project administration intent while each acting operator keeps their own private configuration outside that repository. The same protocol supports automated setup and ordinary project interaction. A portable reference implementation or another conforming execution adapter validates those inputs, observes provider state, plans explicit changes and verifies any applied mutation.

The model is designed so that:

- ordinary collaborators can understand and use the shared GitHub surface without an operator's private systems;
- capable agents can operate through direct provider tools without a local custom binary;
- two operators can use the same repository with different private preferences;
- provider facts are discovered rather than guessed;
- display-name changes do not silently change identity;
- planning is reviewable and execution is stale-safe;
- no agent, account, provider extension or local checkout layout is required by the public protocol.

This document does not define YAML or JSON field names, routing match syntax, default values, schema versions or Go APIs. Those decisions belong to the issues listed in [Deferred decisions](#deferred-decisions).

## Normative vocabulary

### Actors and systems

| Term | Meaning | Explicitly not |
| --- | --- | --- |
| **Operator** | A person who invokes the system or authorises an agent to invoke it. | A repository owner, assignee or provider account merely because those identities happen to match. |
| **Acting principal** | The provider-authenticated identity whose permissions are used for one operation. | The authority for shared contract meaning or private operator preferences. |
| **Agent** | Software acting on an operator's behalf. Agents are interchangeable consumers of the same protocol. | A protocol authority or required product-specific identity. |
| **Collaborator** | Any person or agent that works through the shared repository and provider surfaces. | Someone who must have access to another operator's private profile. |
| **Provider** | An external system that stores or exposes repositories, issues, Projects or optional companion data. GitHub is the first shared-task provider. | A source of desired policy merely because it exposes current state. |
| **Provider adapter** | Code that translates provider-neutral reads and operations to one provider API. | A location for canonical business rules, defaults or routing policy. |
| **Portable client** | The `projectctl` reference implementation used by local operators, shell-capable agents and CI for setup, administration and ordinary interaction. | The public protocol itself or a mandatory hop for every agent. |
| **Environment-specific execution adapter** | A component that applies an authorised request using the capabilities of its environment while preserving the frozen protocol semantics. | A source of defaults, routes, privacy choices or extra mutation authority. |
| **Direct provider adapter** | An execution adapter, commonly an agent with suitable provider tools, that reads and writes the provider without invoking the portable client. | A weaker safety mode or permission to bypass planning and readback. |
| **Execution bridge** | An optional explicitly configured workflow or service that accepts a bounded request and invokes a conforming implementation. | A hidden control plane or protocol prerequisite. |
| **Operator control plane** | An optional private system that helps one operator organise profiles, sources, attention or integrations. | A shared task authority or prerequisite for collaborators. |

An agent has only the authority delegated by its operator and acting principal. The protocol MUST NOT grant an agent additional authority based on its product name or implementation.

### Shared topology

| Term | Meaning | Explicitly not |
| --- | --- | --- |
| **Repository** | A provider repository. In v1, an unqualified repository means a GitHub repository, not a local checkout. | A directory on disk. |
| **Participating repository** | A repository that publishes a supported shared repository contract. | Necessarily the issue store, Project owner or operator control repository. |
| **Issue** | A provider issue used as a shared task record. | A private companion record. |
| **Issue store** | The repository whose Issues collection contains a selected shared task. | Necessarily the participating repository or Project owner. |
| **Project** | A GitHub Projects v2 resource. This document capitalises the term when that product is meant. | A repository, an issue store or a generic body of work. |
| **Project owner** | The provider user or organisation namespace in which a Project exists. | The participating repository, issue store or task assignee. |
| **Project membership** | The relationship by which an issue or draft item appears in a Project. | Ownership of the issue or repository. |
| **Single-Project topology** | A shared topology in which the participating repository declares one Project target for the relevant operation class. | A claim that every issue can belong to only one Project. |
| **Dispatcher topology** | A shared topology that declares multiple named routes or targets and requires resolution before an operation. | A private control plane, a hidden registry or permission to guess a target. |
| **Multi-Project topology** | A topology in which a participating repository can address more than one Project. | Permission for one operation to mutate every named Project. |
| **Operation target** | The one explicit provider-neutral scope selected for an operation, including the relevant repository, issue store and Project references when applicable. | A display title, an unresolved candidate set or an organisation-wide search result. |

Repository ownership, issue storage, Project ownership and Project membership are independent relationships. Implementations MUST preserve those distinctions even when one GitHub repository or account occupies several roles.

A multi-Project declaration describes possible targets. Each mutating request MUST resolve to one explicit operation target before planning. Whether a shared issue may be a member of several Projects is an observed provider fact and MUST NOT be simplified away.

### Configuration and state

| Term | Meaning | Explicitly not |
| --- | --- | --- |
| **Shared repository contract** | A versioned, collaborator-safe declaration published by a participating repository. It describes public topology, stable references, mappings and supported shared policy. | Live provider state, an operator's private preferences or a continuously enforced task template. |
| **Private operator profile** | Configuration controlled by one operator, containing that operator's effective private choices, destinations, local mappings and optional provider settings. | A shared protocol extension that changes repository meaning. |
| **Shared task state** | Mutable issue and Project data that collaborators see through the shared provider. | Seed content that the contract may continually overwrite. |
| **Private companion record** | Optional operator-controlled content associated with a shared issue by a non-sensitive linkage identifier. | The authoritative shared task record or a requirement for collaborators. |
| **Repository-safe source reference** | A non-secret shared reference to material that collaborators are permitted to know exists. | A private source location, credential or promise that every collaborator can access the material. |
| **Private source** | A source location or identifier configured for one operator outside shared repository files. | A shared dependency or collaborator-visible catalogue entry. |
| **Observed provider state** | Provider facts returned by a read for a defined scope. | Desired state, a default or a guarantee that the fact is still current. |
| **Snapshot** | A deterministic provider-neutral representation of the relevant observed provider state at a point in an execution. | A cache that overrides a fresh pre-mutation read. |
| **Canonical model** | A provider-neutral, explicit in-memory meaning produced from a supported shared contract or private profile after normalisation. | A wire document, provider response, resolved target or bag of inferred facts. |
| **Resolved configuration** | A runtime-only result containing one explicit operation target and the effective policy for one request. | A persistent configuration source or a competing authority. |
| **Desired state** | The provider-neutral state requested by validated and resolved intent for fields owned by the requested operation. | A whole-resource replacement or permission to alter unrelated live state. |
| **Plan** | A deterministic comparison of desired state with a snapshot, producing proposed actions, preservation requirements and conflicts. | A mutation, a promise of success or permission to change work omitted from the plan. |
| **Operation** | One typed executable change lowered from an approved plan, with a target, owned fields, expected prior state and desired after-state. | An unbounded callback or arbitrary provider request. |
| **Verification result** | Structured evidence produced by reading back an attempted operation and comparing the observation with its desired after-state. | The provider's mutation response alone. |
| **Verification report** | The ordered aggregate of operation outcomes, readback evidence, skips, refusals and failures for one request. | A success claim without observed evidence. |
| **Derived agent guidance** | A concise deterministic projection of validated protocol meaning and explicit runtime facts for an agent audience. | A configuration authority, private profile or replacement specification. |

The shared repository contract and private operator profile are separate canonical inputs. Normalisation MUST NOT merge them into a single persisted object or erase which authority supplied a value.

### Identity and references

| Term | Meaning | Required property |
| --- | --- | --- |
| **Contract reference** | An identifier declared in a supported contract namespace and used to connect declarations without relying on display text. | Stable within its defined scope and treated according to versioned equality rules. |
| **Provider identity** | An identifier assigned or accepted as authoritative by a provider for a live resource. | Obtained from explicit input or provider discovery, never reconstructed from a display title. |
| **Stable linkage identifier** | A collaborator-safe identifier that associates a shared issue with an optional private companion record. | Contains no private destination, credentials or sensitive content. |
| **Display value** | Human-facing text such as a repository name, Project title, field label or option label. | Never acts as an alternative identity unless a later normative contract explicitly designates that exact field as an identifier. |
| **Local path** | A filesystem location meaningful to one operator's environment. | Never a shared repository fact. |

Contract references and provider identities occupy different namespaces. A contract reference may resolve to a provider identity, but it MUST NOT pretend to be one. A display value may be shown in diagnostics but MUST NOT silently participate in identity matching.

The exact character set, case sensitivity, normal form and namespace syntax for contract references are deferred. Until those rules are specified, implementations MUST preserve supplied identifier bytes and MUST NOT trim, fold case or rewrite them.

## Authority model

Authority answers “which source is allowed to decide this value?”. Ownership, provider permission and runtime capability are separate questions.

### Sources of authority

| Concept | Sole authority | Derived or observed representations | Forbidden competing authority |
| --- | --- | --- | --- |
| Protocol vocabulary and invariants | Merged normative specifications in this repository | Schemas, fixtures and implementation tests | Issue discussion, prototype behaviour or one agent's assumptions |
| Shared contract version and declared topology | Shared repository contract | Canonical shared model and resolved configuration | Private profile, provider display layout or local cache |
| Shared contract references and public mappings | Shared repository contract | Canonical shared model | Display titles or private registry aliases |
| Repository-safe privacy support advertisement | Shared repository contract | Canonical shared model | Capability discovery or another operator's profile |
| Acting operator's effective privacy choice | Private operator profile, subject to the shared contract's supported policy | Canonical private model and resolved configuration | Repository visibility, acting principal identity or shared defaults invented by code |
| Private destination and private source locations | Private operator profile | Provider-specific runtime configuration | Shared files, issue bodies or Project fields |
| Shared issue content and mutable task status | Live shared provider resource | Observed state, snapshot and verification report | Seed definitions, private companion content or stale cache |
| Project fields, values, membership and relationships | Live shared provider resource | Observed state, snapshot and verification report | Desired contract intent until a verified operation changes the provider |
| Private companion content | The operator-controlled companion provider/resource | Private-provider readback | Shared issue content or another operator's profile |
| Repository owner kind, default branch and provider permissions | Live provider response | Observed state and capability result | Naming conventions, repository URL shape or contract defaults |
| Selected operation target | Resolution result for the current request | Plan and report | CLI display order or first provider search result |
| Proposed change | Plan computed from canonical desired state and a snapshot | Typed operations | Executor improvisation |
| Mutation success | Matching provider readback | Verification result and report | HTTP/GraphQL acceptance, absence of an error or optimistic local state |
| Agent guidance | Validated contracts plus explicit runtime capability facts | Deterministic audience-specific projection | Prompt wording, agent brand or renderer defaults |

### Conflict rules

1. A private profile MUST NOT override shared repository semantics.
2. A shared contract MUST NOT publish or select a private destination.
3. A provider observation may show that desired state is not currently realised; it does not change the contract's meaning.
4. A plan MAY propose bringing owned provider state towards desired state, but it MUST preserve unrelated live state.
5. A stale snapshot never wins over a fresh pre-mutation read.
6. A private companion MAY supplement a shared issue but MUST NOT become the sole location of shared identity, status or collaboration-critical instructions.
7. If two authoritative inputs conflict within their permitted domains, validation or resolution MUST fail with an explicit diagnostic. Precedence MUST NOT be invented at execution time.

## Ownership, permission, support and capability

The following terms MUST remain distinct:

| Term | Question answered | Example |
| --- | --- | --- |
| **Authority** | Which source is allowed to decide a value? | The private operator profile decides that operator's private destination. |
| **Resource owner** | Which provider account owns a live resource? | A GitHub organisation owns a repository or Project. |
| **Assignee** | Which provider identity is assigned responsibility for a task? | A GitHub user is assigned to an issue. |
| **Feature support** | Can the provider and protocol represent this feature for the selected topology? | Projects v2 fields are supported by the GitHub adapter. |
| **Permission** | Does the acting principal have provider authorisation for a specific action on a specific resource? | The token may update an issue but not a Project field. |
| **Capability result** | Is a requested operation currently supportable, based on explicit feature and permission observations? | Reads are possible, but a required Project mutation is unavailable. |
| **Policy** | Is the action allowed or required by validated shared/private intent? | The effective privacy policy prohibits placing companion content in GitHub. |

A capability result is derived runtime evidence. It MUST NOT become a second declaration of policy, topology or privacy support. Authentication success alone does not imply permission. A missing or failed probe does not prove absence. Exact capability result categories belong to #35.

## Privacy and companion boundary

### Shared advertisement

The shared repository contract MAY advertise collaborator-safe privacy behaviour supported by the repository. That advertisement defines what the shared workflow can accommodate. It MUST NOT contain:

- an operator's effective choice;
- a private document, folder, repository or account identifier;
- credentials or access instructions;
- local filesystem paths;
- private source locations;
- another operator's settings.

The exact modes, wire representation and repository-safe linkage hook belong to #15.

### Operator choice

The private operator profile owns the acting operator's effective choice and any destination required by it. The choice applies only to operations performed for that operator. Repository visibility, account identity and the presence of a configured private provider MUST NOT be used to infer the choice.

The profile MAY be absent for workflows that require no private configuration. Exact absence and default behaviour belong to #32 and #34.

### Private companion

A private companion record is optional, operator-scoped and supplemental. When used:

- the shared issue retains the collaborator-visible task identity and shared status;
- the shared side contains at most the stable non-sensitive linkage representation defined by #15;
- the private destination and provider identity remain in the private profile;
- another operator may use a different destination, use no companion or choose a different supported mode;
- failure of the private provider MUST be explicit and MUST NOT be reported as verified shared success;
- private content MUST NOT be copied into shared state as a fallback.

## Conceptual topologies

### Single-Project repository

A participating repository declares one Project target for the relevant operation class. The repository containing the contract, the issue store and the Project owner may be different provider resources. Resolution makes those roles explicit before provider reads or planning.

### Dispatcher and multi-Project repository

A participating repository declares several named candidate routes or targets. A request supplies or derives only the routing inputs allowed by the future routing contract. Resolution must return exactly one operation target or an explicit failure before planning.

The dispatcher is shared protocol data. It is not an operator's private repository registry. Route key syntax, precedence, fallback, wildcard and multiple-match behaviour belong exclusively to #14.

### Shared issue with private companions

One shared issue may have zero or more operator-specific private companion records, one per participating operator/provider configuration as permitted later. Those private records do not compete with each other because none is the authority for shared task state. Shared files remain identical for all operators.

## Processing pipeline

Each stage has one responsibility and a typed boundary. Later stages MUST NOT repair or silently reinterpret earlier-stage failures. A portable client and a direct provider adapter MAY package stages differently, but both MUST preserve the same inputs, outputs, authority and failure meaning. CLI serialisation is not a required intermediate protocol for another adapter.

| Stage | Input | Output | Owns | MUST NOT |
| --- | --- | --- | --- | --- |
| **1. Parse** | Raw bytes and selected source | Version-discriminated wire value with presence information | Syntax and source-location errors | Fill defaults, trim identifiers, discover provider facts or mutate the source |
| **2. Schema validation** | Parsed wire value | Schema-valid wire value | Field presence, type, shape and version-specific allowed fields | Apply defaults, perform cross-resource semantics or access providers |
| **3. Normalisation** | Schema-valid wire value | New canonical shared or private model | Explicit version-defined defaults and transformations | Mutate input, guess facts, merge authorities or perform routing |
| **4. Semantic validation** | Canonical model | Ordered diagnostics or valid canonical model | Cross-field and conceptual invariants | Fill defaults, rewrite identity, query providers or select targets |
| **5. Discovery** | Requested operation and explicit provider scope | Provider facts, identities, permissions and feature observations | Read-only external facts | Decide policy, alter configuration or turn unknown into absent |
| **6. Resolution** | Valid canonical inputs plus permitted discovered facts | One resolved configuration or typed failure | Target selection and effective-policy selection under frozen rules | Search unbounded scope, choose by display order or mutate state |
| **7. Snapshot** | Resolved scope and provider reads | Deterministic observed-state snapshot | Canonical ordering and stale-sensitive observed values | Compare desired state, hide unavailable facts or mutate providers |
| **8. Planning** | Valid resolved intent and snapshot | Complete plan, diff, conflicts and preservation requirements | Pure desired-versus-observed comparison | Access the network, add hidden apply-time work or mutate inputs |
| **9. Operation lowering** | Approved plan | Ordered typed operations and exact preconditions | Owned changes, expected prior state and desired after-state | Weaken preconditions or encode arbitrary provider calls |
| **10. Execution** | Typed operations and provider interface | Attempt results | Fresh targeted reread, stale refusal and owned mutation | Re-plan, broaden the change, continue after unsafe failure or claim verification |
| **11. Verification** | Attempted operation and provider readback | Verification result/report | Desired-after comparison and preserved-state evidence | Treat mutation acceptance as success or hide mismatches |

### Determinism

Parsing diagnostics, normalised models, semantic diagnostics, resolution results, snapshots, plans, operation order and structured reports MUST be deterministic for equivalent inputs. Provider response ordering, map iteration order and wall-clock timestamps MUST NOT change canonical comparison output.

### Failure discipline

- Unknown, unavailable, unauthorised and absent are distinct states.
- A stage MUST expose failure rather than manufacture a value needed by the next stage.
- Mutation MUST fail closed when current stale-sensitive state cannot be read or differs from the approved precondition.
- Verification MUST fail when readback is unavailable or differs from desired after-state.
- Failure of optional private work MUST be represented separately from shared-provider results.

## Mutation ownership and preservation

A requested operation owns only the fields and relationships explicitly identified by its validated intent and plan.

Every mutating workflow MUST:

1. produce a complete plan before writes;
2. identify the exact target and owned change set;
3. carry expected prior state or expected absence;
4. re-read stale-sensitive state immediately before each write;
5. refuse unknown or changed prior state;
6. modify only owned fields;
7. stop after the first unsafe refusal or provider failure;
8. read back each accepted write;
9. compare desired and observed after-state;
10. report unverified or mismatched state as failure.

The non-interactive portable client exposes the plan as the default dry run and requires explicit apply selection. An interactive agent MAY execute within an explicit current operator request without asking for a redundant CLI-style confirmation. That authorisation does not remove planning, ambiguity refusal, stale checks, owned-field limits or readback verification.

Seed and bootstrap definitions are create-once intent. Once a shared issue or equivalent live record exists, its ordinary human or bot edits are live shared state. A bootstrap workflow MUST NOT continually restore the original seed text.

## Synthetic conceptual examples

These examples define boundaries only. They deliberately omit final wire names and routing algorithms.

### One repository and one Project

The synthetic participating repository `example-org/alpha` publishes a shared contract that references one issue store and one Project. The Project title is later renamed in GitHub.

- The contract reference remains the protocol identity.
- Discovery observes the current provider identity and display title.
- The rename alone does not change routing identity.
- A plan may display the new title but cannot substitute it for the reference.

### One repository and several Projects

The synthetic repository `example-org/beta` declares several named targets. A request does not resolve uniquely under the routing decision table.

- Resolution fails before snapshotting or planning.
- The resolver does not choose the first Project returned by GitHub.
- The CLI does not prompt interactively in a non-interactive workflow unless a later CLI contract explicitly defines such a mode.

### Two operators, different private choices

Operator A and Operator B use the same shared contract and issue. A's private profile selects a private companion provider. B's profile selects a supported shared-only mode.

- The repository contract and shared issue do not contain A's private destination.
- B does not need access to A's profile or companion.
- Both see the same shared task identity and status.
- A private-provider failure is reported for A's operation and does not cause private content to fall back into the shared issue.

### Stale provider state

A plan proposes adding an issue to a Project based on snapshot S. Before apply, another collaborator changes the relevant membership.

- Execution re-reads that membership.
- The operation refuses the stale precondition.
- The executor does not recompute a new plan or continue with later writes.
- A new snapshot and plan require a new explicit invocation.

## Settled decisions and deferred decisions

### Settled here

- Shared contract, private profile, live provider state and runtime artefacts are distinct authority layers.
- The shared contract never owns an operator's effective private choice or destination.
- Repository, issue store, Project owner and Project membership are distinct roles.
- Every mutating request resolves one explicit operation target.
- Contract references, provider identities and display values are distinct.
- Provider facts and permissions are discovered, not guessed.
- Capability is runtime evidence, not policy or configuration.
- The pipeline stages and their side-effect boundaries are fixed.
- Plans are complete, pure and non-authoritative.
- Mutations change only owned fields, fail closed on stale/unknown state and require matching readback.
- Shared issues remain the authoritative shared task records; private companions are optional and supplemental.
- Public protocol artefacts remain operator-neutral and agent-neutral.
- The shared protocol, not CLI JSON or one implementation, is the interoperability contract.
- The portable client, direct provider adapters and optional bridges are interchangeable execution paths only when they preserve the same semantics.
- Derived agent guidance adds no authority and remains traceable to validated inputs.
- Automated setup and ordinary project interaction are both first-class product outcomes.

### Specialised decision owners

| Issue | Decision owner |
| --- | --- |
| #13 | Shared v1 file location, encoding, version discriminator, single-Project schema, field presence and base fixtures |
| #14 | Dispatcher/multi-Project wire additions, routing key syntax, equality, case behaviour, wildcard rules, precedence, fallback and ambiguity outcomes |
| #15 | Shared privacy advertisement modes and repository-safe stable linkage representation |
| #46 | Protocol, portable-client, execution-adapter and interactive-authorisation boundary |
| #32 | Private operator-profile wire format, storage expectations, effective choice, private destination and companion-provider contract |
| #47 | Label, routing-label and label-based sub-project declarations |
| #48 | Semantic project dimensions and provider bindings |
| #49 | Automated setup and manual-action outcomes |
| [#50](v1-task-interactions.md) | Ordinary task-interaction semantics and granular mutation scopes, now settled |
| #34 | Version vocabulary, absent/null/empty/zero semantics, all defaults, normalisation, unknown-version and migration behaviour |
| #51 | Concise derived agent-guidance projection and renderer contract |
| #16 | Frozen v1 conformance corpus, semantic decision rows and expected diagnostics |
| #31 | Supported legacy inventory and exact compatibility/migration fixtures |
| #26 | Final structured report schema and CLI exit codes |

No downstream issue may contradict a settled specialised specification or decide a remaining deferred item implicitly. If implementation requires an answer that its named normative issue does not provide, work stops at that boundary and the normative issue is corrected first.

## Normative invariant register

Later specifications, fixtures and reviews MUST cite these stable requirement IDs where they exercise or refine a conceptual invariant. An ID MUST NOT be reused for a different rule.

| ID | Invariant |
| --- | --- |
| `CM-001` | Every persisted or observed value has one identifiable authority and provenance. |
| `CM-002` | Shared repository intent and private operator intent remain separate inputs and representations. |
| `CM-003` | Shared files and shared task state contain no private destination, credential, local path or private source identifier. |
| `CM-004` | Participating repository, issue store, Project owner and Project membership remain distinct roles. |
| `CM-005` | Contract references, provider identities, stable linkage identifiers and display values remain distinct identity classes. |
| `CM-006` | Display values do not silently act as identity or fallback selectors. |
| `CM-007` | Provider facts, ownership and permissions are supplied explicitly or discovered; they are never guessed. |
| `CM-008` | A mutating request resolves one explicit operation target before snapshotting or planning. |
| `CM-009` | Parse, schema validation, normalisation, semantic validation, discovery, resolution, snapshot, planning, lowering, execution and verification retain separate contracts. |
| `CM-010` | Normalisation returns a new canonical value and applies only version-defined transformations. |
| `CM-011` | Semantic validation is pure and does not normalise, discover or resolve. |
| `CM-012` | Discovery is read-only evidence and cannot become policy or desired configuration. |
| `CM-013` | Planning is pure, deterministic and complete; apply introduces no hidden work. |
| `CM-014` | Every executable operation identifies its exact target, owned change, expected prior state and desired after-state. |
| `CM-015` | Execution freshly reads stale-sensitive state and refuses a mismatch or unavailable fact before mutation. |
| `CM-016` | Mutation success requires matching provider readback; provider acceptance alone is not success. |
| `CM-017` | Seed intent is create-once and does not continually overwrite live shared edits. |
| `CM-018` | A private companion is optional and supplemental; shared issues remain the authoritative shared task records. |
| `CM-019` | Public protocol artefacts do not require or privilege a particular operator, agent, private control plane or optional provider. |
| `CM-020` | Canonical outputs are deterministic and unaffected by unordered provider responses or incidental timestamps. |
| `CM-021` | Unknown, unavailable, unauthorised and observed-absent states remain distinguishable. |
| `CM-022` | A stage reports missing or conflicting required knowledge instead of manufacturing a value for the next stage. |
| `CM-023` | The shared protocol is the interoperability contract; no conforming agent is required to invoke the portable client or consume CLI JSON. |
| `CM-024` | Portable clients, direct provider adapters and execution bridges preserve identical authority, target, ownership, stale-state and verification semantics. |
| `CM-025` | Derived agent guidance is deterministic, traceable and non-authoritative. |
| `CM-026` | Explicit interactive operator authorisation may select apply without a redundant CLI confirmation, but never weakens planning, ambiguity refusal, stale protection or readback. |
| `CM-027` | Automated setup and ordinary shared-task interaction remain first-class product outcomes across supported topologies. |

## Conformance checklist

A later specification or implementation conforms to this conceptual model only if all answers below are “yes”:

- Can every persisted value be traced to exactly one authority?
- Are shared contract data and private operator data represented separately?
- Are repository, issue store, Project and membership roles explicit?
- Are contract references, provider identities and display values distinguishable?
- Are external facts marked as supplied or discovered rather than inferred?
- Is every request resolved to one target before planning?
- Are validation, normalisation and discovery separate operations?
- Can planning run without provider access?
- Does every operation name its owned fields and exact prior-state expectation?
- Can stale or unavailable state prevent a write?
- Is successful mutation reported only after matching readback?
- Can a collaborator operate on shared state without access to any private profile or companion provider?
- Can a capable agent perform the same conforming operation through provider tools without invoking the portable client?
- Does concise agent guidance trace every value to validated protocol input or explicit runtime fact?
- Are automated setup and ordinary issue/Project interaction both represented in the conformance surface?
- Can two operators use different private choices without changing shared files?
- Are deferred wire and algorithm choices left to their owning issue?

## Relationship to the architecture boundary

This specification refines [the v1 architecture boundaries](../architecture/v1-boundaries.md). If the two documents appear to conflict, the architecture boundary controls package and layer direction, and this document controls terminology and conceptual authority. A genuine contradiction must be resolved in the normative documents before dependent work proceeds.
