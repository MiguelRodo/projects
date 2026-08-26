# Specification: Neutral v1 Terminology & Canonical Data Model

**Status**: Baseline Draft / v1  
**Scope**: Issue [#12](https://github.com/MiguelRodo/projects/issues/12) (Child of Parent Epic [#7](https://github.com/MiguelRodo/projects/issues/7))

---

## 1. Overview & Objectives

This document defines the agent-neutral vocabulary, canonical in-memory domain structures, baseline validation invariants, and privacy boundary separation for the `projects` v1 protocol.

The core objectives are:
1. **GitHub as Shared Authority**: Repository manifests, issues, pull requests, and GitHub Projects are the canonical shared surface for task collaboration.
2. **Neutral Nomenclature**: Terminology is strictly decoupled from specific operators, personal names, account identifiers, proprietary tools, or specific AI agent runtimes.
3. **Decoupled Operator Extensions**: Personal agendas, private control repositories, and external companion backends (e.g. user-owned documents, private stores) are treated as optional downstream consumers without leaking into public repository contracts.
4. **Canonical Internal Representation**: Domain models represent the normalised in-memory structures into which supported repository contracts will be parsed. Wire schema formats and routing execution algorithms are formally specified in child issues [#13](https://github.com/MiguelRodo/projects/issues/13), [#14](https://github.com/MiguelRodo/projects/issues/14), and [#15](https://github.com/MiguelRodo/projects/issues/15).

---

## 2. Neutral Glossary

| Term | Definition |
| :--- | :--------- |
| **`ProjectIdentity`** | Identifies a logical project (`name`, `slug`, `owner`, `description`). |
| **`RepositoryRef`** | A reference to a git repository managed within a workspace (`name`, `url`, `path`, `default_branch`). |
| **`TargetProject`** | Identifies a target GitHub Project (Projects v2) by its owner kind (`organization` or `user`), owner name, project number, title, and local contract reference (`ref`). |
| **`IssueStore`** | A repository configured as a source or destination for tracking issues across one or more projects. |
| **`FieldMapping`** | A declaration mapping a canonical task field (e.g. `status`, `priority`, `estimate`) to a GitHub Project custom field name and kind. |
| **`ValueMapping`** | A declaration mapping canonical values (e.g. `todo`, `in_progress`, `done`) to specific option names or IDs in GitHub Projects. |
| **`FieldKind`** | Enumeration of supported custom field types: `text`, `number`, `single_select`, `iteration`, `date`. |
| **`RouteKey`** | The attribute used by a dispatcher to route issues to a target project (e.g. `label`, `repository`, `component`). |
| **`RouteRule`** | A rule mapping a specific key-value condition to a `TargetProject` reference with default fields and labels. |
| **`DispatcherConfig`** | Configuration declaring routing targets, default fallback targets, and rule lists for multi-project topologies. |
| **`SourceRef`** | A reference to a repository-safe documentation, catalogue, or asset source (`id`, `name`, `location`, `type`, `description`). |
| **`Capability`** | Protocol capabilities supported by the repository contract (e.g. `read_issues`, `write_issues`, `manage_projects`, `route_dispatcher`, `private_companion`). |
| **`MutationPolicy`** | Rules governing write operations (`allow_create`, `allow_update`, `allow_delete`, `stale_write_guard`). |
| **`PrivacyMode`** | The privacy mode for task context: `shareable_by_default` (default) or `full_github_context`. |
| **`RepositoryPrivacyPolicy`** | The shared repository contract's advertised privacy capabilities (`supported_modes`, `default_mode`, `allow_private_companion`). |
| **`UserPrivacyPreference`** | An acting user's private session configuration declaring their effective privacy mode and private companion location. |
| **`StableLinkageID`** | A non-sensitive unique identifier (UUID/slug) allowing private companion records to be associated with GitHub issues without publishing private URLs. |
| **`CompanionRecord`** | An optional operator-owned companion record stored in user-controlled private storage (e.g. personal document, private repository, local file). |

---

## 3. Privacy Policy & Acting-User Separation

### Shared Repository Capability vs. Acting-User Choice
A central requirement of the protocol is that **different collaborators using the same repository can choose different privacy modes and private destinations**:

1. **Shared Repository Level (`RepositoryPrivacyPolicy`)**:
   - The repository contract advertises which privacy modes it supports (`supported_modes`, e.g. `["shareable_by_default", "full_github_context"]`) and a `default_mode`.
   - It advertises whether client companion linkage is permitted (`allow_private_companion: true`). This is an explicit opt-in capability.
   - The repository manifest **never** contains private document URLs, personal file paths, or private destination endpoints.
2. **Acting-User Level (`UserPrivacyPreference`)**:
   - The individual user's private configuration (configured out-of-band in their personal operator profile or CLI flags) sets their `effective_mode`.
   - If using `shareable_by_default` and the repository allows companion linkage, the user's authorized agent routes sensitive context to their configured `private_companion_ref` (e.g. a personal document or private repository) using a non-sensitive `StableLinkageID`.
   - If the repository policy has `allow_private_companion: false`, the user cannot link a companion reference.
   - Another collaborator on the same repository may select `full_github_context` with no private companion destination, without modifying the shared repository contract.

### Privacy Modes
- **`shareable_by_default` (Default)**: Task descriptions and public fields are shareable; sensitive companion context is stored out-of-band.
- **`full_github_context`**: Ordinary task context is permitted to reside directly on GitHub.

---

## 4. Supported Topologies & Owner Kinds

GitHub Projects can be owned by organizations or individual user accounts:
- **`OwnerKindOrganization` (`"organization"`)**: Standard team and organization boards.
- **`OwnerKindUser` (`"user"`)**: Personal developer portfolios and solo-maintained project boards.
- **`OwnerKindUnspecified` (`""`)**: Permitted during draft contract authoring; inferred during authoritative GitHub discovery.

Target projects must not silently default to organization ownership.

---

## 5. Canonical In-Memory Hierarchy

```
SingleProjectContract
├── SchemaVersion (optional in-memory metadata)
├── Project (ProjectIdentity: Name, Slug, Owner, Description)
├── Repository (RepositoryRef: Name, URL, Path, DefaultBranch)
├── Target (TargetProject: Ref, Owner, OwnerKind, Number, Title, URL)
├── Mappings ([]FieldMapping: CanonicalName, GitHubField, Kind, Required, Values)
├── Capabilities (CapabilitySet: Items)
├── Mutation (MutationPolicy: AllowCreate, AllowUpdate, AllowDelete, StaleWriteGuard)
├── Privacy (RepositoryPrivacyPolicy: SupportedModes, DefaultMode, AllowPrivateCompanion)
└── Sources ([]SourceRef: ID, Name, Location, Type, Description)

MultiProjectContract (Dispatcher)
├── SchemaVersion (optional in-memory metadata)
├── IssueStore (RepositoryRef: Name, URL, Path, DefaultBranch)
├── Projects ([]ProjectIdentity)
├── Targets ([]TargetProject)
├── Dispatcher (DispatcherConfig: DefaultTargetRef, Fallback, Routes)
├── Mappings ([]FieldMapping)
├── Capabilities (CapabilitySet)
├── Mutation (MutationPolicy)
├── Privacy (RepositoryPrivacyPolicy)
└── Sources ([]SourceRef)
```

---

## 6. Validation Invariants & Scope Boundaries

### Baseline Model Validation Invariants
The canonical Go domain model enforces the following structural rules:
1. **Identities**: Project and repository names must be non-empty. Slugs and canonical field names must conform to lowercase kebab-case (`^[a-z0-9]+(-[a-z0-9]+)*$`).
2. **Targets**: Target project owners must be non-empty and numbers must be positive integers (> 0). If specified, `OwnerKind` must be `organization` or `user`. Local `Ref` lookups are exact matches on `TargetProject.Ref`.
3. **Field Mappings**: Canonical field names within a contract must be unique. `single_select` mappings must define at least one value mapping.
4. **Routing**: Route rules must declare non-empty keys and target references. Duplicate route conditions (`key=value`) are rejected using exact structural matching. All target references must resolve to declared targets in the contract.
5. **Privacy Policy**: Default and supported privacy modes must be recognized enum values. Private companion capability is opt-in, and companion preferences are rejected when companion linkage is not permitted.

### Deferred Wire & Routing Semantics
Issue [#12](https://github.com/MiguelRodo/projects/issues/12) specifies the terminology and canonical domain structures. Wire-level schema formats, exact matching semantics (e.g. wildcards, precedence, regexes, case sensitivity), and full conformance suites are deferred to the dependent issues:
- **Issue [#13](https://github.com/MiguelRodo/projects/issues/13)**: Single-project wire schema specification & fixtures.
- **Issue [#14](https://github.com/MiguelRodo/projects/issues/14)**: Multi-project dispatcher wire schema & routing resolution rules.
- **Issue [#15](https://github.com/MiguelRodo/projects/issues/15)**: Privacy modes & private extension schemas.
- **Issue [#16](https://github.com/MiguelRodo/projects/issues/16)**: Semantic conformance test fixtures.

---

## 7. Legacy Adaptation & Neutrality

Any legacy operator-specific terminology, private field aliases, or personal repository names belong strictly in private compatibility adapters or operator-level profiles. The public v1 protocol recognizes only generic, configurable domain concepts.
