# Specification: Neutral v1 Terminology & Canonical Data Model

**Status**: Draft / v1  
**Related Issues**: [#7](https://github.com/MiguelRodo/projects/issues/7), [#12](https://github.com/MiguelRodo/projects/issues/12)

---

## 1. Overview & Objectives

This document defines the agent-neutral vocabulary, canonical in-memory data structures, invariants, and validation rules for the `projects` v1 protocol.

The core design objectives are:
1. **GitHub as Shared Authority**: Repository manifests, issues, pull requests, and GitHub Projects are the canonical shared surface for task collaboration.
2. **Neutral Nomenclature**: Terminology is decoupled from specific operators, personal names, account identifiers, proprietary tools, or specific AI agent runtimes.
3. **Decoupled Operator Extensions**: Personal agendas, private control repositories, and external companion backends (e.g. Google Drive/Docs/Sheets) are treated as optional downstream consumers without leaking into public repository contracts.
4. **Deterministic Normalisation**: All supported contract versions are deserialised and normalised into a canonical internal domain model with well-defined invariants.

---

## 2. Neutral Glossary

| Term | Definition |
| :--- | :--------- |
| **`ProjectIdentity`** | Identifies a project within an organization or user namespace (`name`, `slug`, `owner`, `description`). |
| **`RepositoryRef`** | A reference to a git repository managed within a project workspace (`name`, `url`, `path`, `default_branch`). |
| **`TargetProject`** | Represents a target GitHub Project (Projects v2) by its owner kind (`organization` or `user`), owner name, project number, and title. |
| **`IssueStore`** | A repository configured as a source or destination for tracking issues across one or more projects. |
| **`FieldMapping`** | A declaration mapping a canonical task field (e.g. `status`, `priority`, `estimate`) to a GitHub Project custom field name and kind. |
| **`ValueMapping`** | A declaration mapping canonical values (e.g. `todo`, `in_progress`, `done`) to specific option names or IDs in GitHub Projects. |
| **`FieldKind`** | Enumeration of supported field types: `text`, `number`, `single_select`, `iteration`, `date`. |
| **`RouteKey`** | The attribute used by a dispatcher to route issues to a target project (e.g. `label`, `repository`, `field_value`). |
| **`RouteRule`** | A rule mapping a specific key-value condition to a `TargetProject` with default fields and labels. |
| **`Dispatcher`** | Routing logic that determines target GitHub Project assignments for multi-project topologies. |
| **`SourceRef`** | A reference to a documentation, catalogue, or asset source (`id`, `name`, `location`, `type`, `description`). |
| **`Capability`** | Feature capabilities supported by the repository contract (e.g. `read_issues`, `write_issues`, `manage_projects`, `route_dispatcher`, `private_companion`). |
| **`MutationPolicy`** | Rules governing write operations (`allow_create`, `allow_update`, `allow_delete`, `stale_write_guard`). |
| **`PrivacyMode`** | The privacy policy mode for task context: `shareable_by_default` (default) or `full_github_context`. |
| **`StableLinkageID`** | A non-sensitive unique identifier (UUID/slug) allowing private companion records to be associated with GitHub issues without publishing private URLs. |
| **`CompanionRecord`** | An optional operator-owned companion record stored in user-controlled private storage (e.g. personal Google Doc, private repository, local file). |

---

## 3. Privacy Policy & Companion Separation

### Privacy Modes
The v1 protocol supports explicit, versioned privacy modes:

1. **`shareable_by_default` (Default)**:
   - Repository content, task descriptions, and field values are treated as potentially public/shareable.
   - Sensitive personal context is not stored on GitHub; it may be routed to an acting user's privately configured companion storage using a `StableLinkageID`.
2. **`full_github_context`**:
   - The repository or acting user explicitly permits full task context to be stored directly in GitHub issues and project fields.

### Privacy Invariants
- Public repository manifests must **never** contain personal Google Drive/Docs/Sheets URLs, private repository paths, or personal credentials.
- A shared repository may advertise `CapabilityPrivateCompanion: true`, indicating that clients may link companion records via `StableLinkageID`, but the destination configuration remains in the user's private operator profile.
- Different collaborators working on the same repository may use distinct private destinations without mutual interference.

---

## 4. Canonical Data Model Hierarchy

```
SingleProjectContract
├── SchemaVersion (e.g. "projects.dev/v1")
├── Project (ProjectIdentity)
├── Repository (RepositoryRef)
├── Target (TargetProject)
├── Mappings ([]FieldMapping)
│   └── Values ([]ValueMapping)
├── Capabilities (CapabilitySet)
├── Mutation (MutationPolicy)
├── Privacy (PrivacyPolicy)
└── Sources ([]SourceRef)

MultiProjectContract (Dispatcher)
├── SchemaVersion (e.g. "projects.dev/v1")
├── Dispatcher (DispatcherConfig)
│   ├── DefaultTarget (TargetProject)
│   ├── Routes ([]RouteRule)
│   └── FallbackRule (FallbackBehavior)
├── IssueStore (RepositoryRef)
├── Projects ([]ProjectIdentity)
├── Targets ([]TargetProject)
├── Mappings ([]FieldMapping)
├── Capabilities (CapabilitySet)
├── Mutation (MutationPolicy)
└── Privacy (PrivacyPolicy)
```

---

## 5. Invariants & Normalisation Rules

All schema parsers and Go structures must enforce the following invariants:

1. **Schema Version**: Must be non-empty and start with a known namespace (e.g. `projects.dev/v1`).
2. **Deterministic Slugs**: Slugs and canonical field names must consist of lowercase alphanumeric characters and hyphens (`^[a-z0-9]+(-[a-z0-9]+)*$`).
3. **Unique Field Names**: Canonical field names within a contract must be unique.
4. **Unique Route Keys**: Route rules must not contain duplicate key-value pairs that create ambiguous target resolution.
5. **Target Project Resolution**: Every route in a multi-project contract must resolve to a valid target declared in `Targets` or `DefaultTarget`.
6. **Value Mapping Completeness**: Single-select field mappings must declare valid canonical values and their corresponding remote option representations.
7. **Stale Write Guard**: When `MutationPolicy.StaleWriteGuard` is enabled, write operations must verify etags / timestamps before mutating GitHub resources.

---

## 6. Legacy Terms Relegation & Migration Mapping

To maintain backward compatibility with legacy prototypes while ensuring public cleanliness:

| Legacy Operator Term | Canonical Neutral Term | Migration Notes |
| :------------------- | :--------------------- | :-------------- |
| `miguel_private` | `PrivacyModeShareableByDefault` | Configured via `PrivacyPolicy.Mode` |
| `assigned_to_miguel` | `AssigneeFilter` / `RouteRule` | Generic route rule matching assignee field |
| `MiguelRodo/issues` | `IssueStoreRef` | Explicit repository reference |
| `Global tasks` | `TargetProject` | Configurable project title/number |
| Personal Drive links | `StableLinkageID` | Private companion backend configured out-of-band |
