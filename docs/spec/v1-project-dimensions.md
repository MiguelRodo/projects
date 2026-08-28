# v1 project dimensions and provider bindings

Status: **normative pre-v1 specification**

Issue: #48

This document defines the provider-neutral task dimensions and their target-specific GitHub storage bindings. It extends the repository contract from [v1 single-Project contract](v1-single-project-contract.md), [v1 routing](v1-routing.md) and [v1 labels and sub-projects](v1-labels-and-subprojects.md).

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Product boundary

A task has one canonical meaning even when two GitHub targets store that meaning differently. For example, canonical Class `bug` may be stored in a personal Project custom field or in an organisation Issue Type. Canonical Priority `p2` may be stored as Project option `P2` or organisation Issue Field values `Medium` and `Low`.

The shared contract declares the canonical vocabulary once. Each target declares exactly one provider binding per enabled dimension. Provider names select storage; they never become canonical identity.

This layer does not define ordinary mutation requests, setup workflows, views, caches or Go types. It defines the information those later layers consume.

## Frozen decisions

| Question | V1 decision |
| --- | --- |
| Canonical declaration | Required `spec.dimensions`, which may be empty |
| Single-Project bindings | Required `spec.dimensionBindings`, which may be empty |
| Dispatcher bindings | Required `dimensionBindings` on every target |
| Binding coverage | Every declared dimension is bound exactly once on every target |
| Canonical identity | Fixed dimension and value references |
| Provider identity | Discovered provider IDs, absent from shared files |
| Class values | Explicit non-empty subset of the finite v1 registry |
| Priority values | Exactly `p0`, `p1`, `p2` |
| Status values | Exactly `todo`, `in-progress`, `done` |
| Workstream values | Contract-specific stable references with canonical display metadata |
| Project, sub-project and labels | Owned by #47, not dimension bindings |
| Class storage | Project single-select field or native Issue Type |
| Priority storage | Project single-select field or organisation Issue Field |
| Status storage | Built-in Project Status only |
| Workstream storage | Project single-select field only |
| Due-date storage | Project date field only |
| Assignee storage | Native issue assignees only |
| Parent storage | Native parent relationship only |
| Progress storage | Read-only Project sub-issues progress projection |
| Matching | Exact selectors only; alternatives are listed explicitly |
| Clear | Supported for writable enum and date values |
| Native definition changes | Not authorised by #48 |

## Wire representation

### Canonical declarations

`spec.dimensions` is an object whose keys are the enabled dimensions:

```yaml
dimensions:
  class:
    values: [task, bug, enhancement, epic]
  priority:
    values: [p0, p1, p2]
  status:
    values: [todo, in-progress, done]
  workstream:
    values:
      api:
        name: API
        description: API-facing work.
      documentation:
        name: Documentation
        description: Durable project guidance.
  due-date: {}
  assignees: {}
  parent: {}
  sub-issues-progress: {}
```

An empty object explicitly declares that this contract does not expose a standard semantic dimension:

```yaml
dimensions: {}
dimensionBindings: {}
```

There are no implicit declarations or bindings.

### Target bindings

A single-Project contract binds dimensions at `spec.dimensionBindings`. A dispatcher binds them independently on every `spec.targets[]` entry:

```yaml
targets:
  - target: { ... }
    fields: { ... }
    dimensionBindings: { ... }
```

For every target, the exact key set of `dimensionBindings` MUST equal the key set of `spec.dimensions`. A missing binding is `contract.semantic.missing-dimension-binding`; an extra binding is `contract.semantic.undeclared-dimension-binding`.

Order of dimension, value, target and binding declarations has no semantic meaning.

## Canonical dimensions

### Class

Class answers what kind of work item the task is. A contract enables a non-empty subset of:

| Reference | Canonical display | Meaning |
| --- | --- | --- |
| `task` | Task | A specific piece of work |
| `bug` | Bug | An unexpected fault or incorrect behaviour |
| `enhancement` | Enhancement | An improvement to existing work or material |
| `raw-data` | Raw data | Acquisition, intake or stewardship of source data |
| `processed-data` | Processed data | Transformation, validation or publication of derived data |
| `analysis` | Analysis | Investigation, modelling, evaluation or interpretation |
| `report` | Report | A report, presentation, manuscript or similar communication |
| `documentation` | Documentation | Durable guidance, records or reference material |
| `epic` | Epic | A substantial outcome made up of child issues |

The reference is canonical identity. The display and meaning above are fixed by v1 and are not repeated or overridden in the wire contract.

### Priority

Priority is exactly:

| Reference | Canonical display | Meaning |
| --- | --- | --- |
| `p0` | P0 | Urgent or deadline-critical |
| `p1` | P1 | Planned and important |
| `p2` | P2 | Useful but not currently urgent |

All three values are declared when Priority is enabled. Provider vocabularies may be many-to-one on read, but canonical output is always one of these three references.

### Status

Status is exactly `todo`, `in-progress` and `done`. The canonical displays are `Todo`, `In progress` and `Done`. All three are declared when Status is enabled.

Additional provider statuses are not canonical v1 values. Encountering one on a task is `dimension.value.unmapped`. A later protocol version may extend the vocabulary explicitly.

### Workstream

Workstream identifies the strand of the current overall Project that a task advances. It is single-valued and is separate from the #47 sub-project label.

Workstream references use the common contract-reference syntax. Every value has a required collaborator-safe canonical `name` and `description`. Provider option colour and provider-facing description remain target field metadata rather than canonical semantics.

### Due date

Due date is an optional item value in ISO 8601 full-date form, `YYYY-MM-DD`. The dimension declaration is an empty object because it has no value registry. Absence of an item value is a clear due date, not an inferred date.

### Assignees

Assignees are an exact set of issue actors. Provider IDs are observed identity; login is observed display/selector text. Canonical output sorts actors by provider identity and does not infer an assignee from contract owner, Project owner or current user.

### Parent

Parent is zero or one native issue relationship. The canonical value carries the observed parent issue identity and locator. A Project's `Parent issue` display is a projection, not a second authority.

### Sub-issues progress

Sub-issues progress is read-only derived state with non-negative integer `completed` and `total` values where `completed <= total`. It is never written independently of native relationships.

## Project custom-field binding

Class, Priority and Workstream may bind to a declared single-select Project field. Due date may bind to a declared date Project field.

```yaml
class:
  kind: project_field
  fieldRef: class
  values:
    task:
      readValueRefs: [task]
      writeValueRef: task
    bug:
      readValueRefs: [bug]
      writeValueRef: bug
```

`fieldRef` resolves in that target's existing `fields` map. For enum bindings:

- every canonical value has exactly one map entry;
- `readValueRefs` is a non-empty unique set of declared option references;
- `writeValueRef` is a member of `readValueRefs`;
- read sets for different canonical values are disjoint;
- the union of read references equals the declared option references of the bound field.

The bound field type MUST be `singleSelect` for Class, Priority and Workstream, and `date` for Due date. Type mismatch is `contract.semantic.dimension-field-type`.

One provider field reference may bind at most one semantic dimension. A custom field not referenced by any dimension binding remains ordinary non-semantic metadata under #13.

## Native Issue Type binding

Class may bind to an organisation's Issue Types:

```yaml
class:
  kind: issue_type
  values:
    task:
      readNames: [Task]
      writeName: Task
    bug:
      readNames: [Bug]
      writeName: Bug
```

The authority is the explicit organisation owner of the target issue store. The issue store MUST be organisation-owned. This locates the catalogue but does not prove the provider supports Issue Types or that the acting principal may read or write them.

Discovery lists the organisation's enabled Issue Types and resolves every exact name. Disabled exact types are reported as `dimension.issue-type.disabled`. Provider type IDs remain observed state.

GitHub exposes organisation Issue Types and their name, enabled state, description and colour through its [Issue Types REST API](https://docs.github.com/en/rest/orgs/issue-types). This v1 binding adopts existing types only. It does not authorise creating, renaming, enabling, disabling or deleting them.

## Organisation Issue Field binding

Priority may bind to an organisation single-select Issue Field:

```yaml
priority:
  kind: issue_field
  name: Priority
  dataType: singleSelect
  values:
    p0:
      readNames: [Urgent]
      writeName: Urgent
    p1:
      readNames: [High]
      writeName: High
    p2:
      readNames: [Medium, Low]
      writeName: Medium
```

The authority is the explicit organisation owner of the target issue store. Discovery requires exactly one field whose name and data type match exactly. Its Project exposure, option catalogue and item values remain the identity discovered from GitHub.

GitHub's [Issue Fields REST API](https://docs.github.com/en/rest/orgs/issue-fields) exposes organisation fields and options. Updating a field's options replaces the supplied option collection and can destroy identity if existing IDs are omitted. #48 therefore authorises value writes only, not definition changes.

An Issue Field value remains issue-owned even when it is displayed in a Project. GitHub exposes Project issue-field configurations and values separately from Project custom-field values through the [Projects GraphQL API](https://docs.github.com/en/graphql/reference/projects). Attaching a missing Issue Field to a Project is a setup outcome owned by #49, not an item-value write.

## Built-in Project Status binding

Status binds only to the provider's built-in Project Status field:

```yaml
status:
  kind: project_status
  values:
    todo:
      readNames: [Todo]
      writeName: Todo
    in-progress:
      readNames: [In progress, In Progress]
      writeName: In progress
    done:
      readNames: [Done]
      writeName: Done
```

Alternative case or spelling is never inferred. A contract lists every accepted exact name. The discovered configuration must prove built-in Status identity; a custom single-select field merely named `Status` is not accepted.

The map permits deterministic compatibility with existing Project spelling without making display text canonical identity.

## Native and derived bindings

These bindings carry no selectors:

```yaml
assignees: {kind: issue_assignees}
parent: {kind: issue_parent}
sub-issues-progress: {kind: project_sub_issues_progress}
```

Assignees and parent are issue-owned. Sub-issues progress is a read-only Project projection. None may be replaced by a custom field with the same display name.

GitHub documents native parent and progress visibility in [Projects](https://docs.github.com/en/issues/planning-and-tracking-with-projects/understanding-fields/about-parent-issue-and-sub-issue-progress-fields).

## Exact native value maps

Every native enum map entry contains:

- `readNames`, one to 16 unique exact provider names;
- `writeName`, one exact provider name contained in `readNames`.

Across one binding, no exact read name may occur under two canonical values. Static duplicates use `contract.semantic.overlapping-native-value-map`. A write name absent from its read set is `contract.semantic.unreadable-write-value`.

At discovery time, each read name resolves to zero or one provider option identity. More than one exact provider result is ambiguous. At value-read time, the observed provider identity or exact configured selector resolves to exactly one canonical value.

The map is deliberately asymmetric. Canonical Priority `p2` may read `Medium` or `Low` and write `Medium`. Reading then writing the unchanged canonical value is idempotent and does not rewrite `Low` merely to prefer `Medium`; a write occurs only when requested canonical meaning differs or an explicit normalisation operation is later defined.

## Authority and duplicate prevention

For each target and dimension there is exactly one authoritative binding. The following are invalid:

- binding Class to both a Project field and Issue Type;
- binding Priority to both a Project field and Issue Field;
- binding two dimensions to one Project field reference;
- using labels for Class, Priority, Status or Workstream;
- representing sub-project as a Project field;
- treating a Project Parent or progress column as relationship storage;
- inferring a standard dimension from an unbound field's name.

The contract cannot determine the intended semantics of arbitrary live provider fields. Discovery reports same-name unbound provider configuration as observed unmanaged state. An adapter never silently adopts it as a semantic binding.

## Discovery states

Every binding produces one of these configuration states:

| State | Meaning |
| --- | --- |
| `known` | Exactly one supported binding and all required value identities were resolved |
| `missing` | A declared provider field, type, option or Project attachment does not exist |
| `unknown` | Required discovery has not run or is incomplete |
| `unavailable` | The provider or environment cannot currently return the facts |
| `forbidden` | The acting principal lacks required read permission |
| `unsupported` | The provider or adapter does not implement the declared binding kind |
| `ambiguous` | More than one provider identity satisfies an exact selector |

Owner kind, token shape, current Project contents and CLI availability MUST NOT be used to replace discovery.

Item values distinguish:

- `known` with one canonical value;
- `clear` with no provider value;
- `unknown`, `unavailable`, `forbidden` or `unsupported`;
- `unmapped` when a known provider value has no declared canonical map;
- `ambiguous` when configuration maps it to more than one canonical value.

Missing configuration is never reported as a clear item value.

## Read semantics

Reading a dimension proceeds in this order:

1. resolve the contract and target;
2. require the declared dimension and exact target binding;
3. require a known supported binding configuration;
4. read the authoritative provider value;
5. preserve non-known fact states;
6. return clear when the provider explicitly reports no value;
7. map one known provider value to one canonical value;
8. return a typed unmapped or ambiguous failure otherwise.

Canonical output contains `dimensionRef`, `state`, and the canonical value appropriate to the dimension. Provider names may be included only in an observed-state trace, never as replacement identity.

## Write and clear semantics

#48 defines storage authority and required scope, not the ordinary request envelope owned by #50.

For a requested canonical set or clear:

1. require a known complete before-state;
2. require the target's one binding;
3. validate the canonical value against the declaration;
4. require a known writable provider configuration;
5. lower through the binding's writable provider identity;
6. require the binding-family mutation scope;
7. return a no-op if canonical before and after values are equal, including clear to clear;
8. otherwise emit one owned write with a stale precondition;
9. read back and map through the same binding.

Scope mapping is exact:

| Binding | Required value-write scope |
| --- | --- |
| `project_field` | `project.item.field.write` |
| `project_status` | `project.item.field.write` |
| `issue_type` | `issue.type.write` |
| `issue_field` | `issue.field.write` |
| `issue_assignees` | Deferred to #50 |
| `issue_parent` | Deferred to #50 |
| `project_sub_issues_progress` | Never writable |

An idempotent no-op requires no mutation scope. Scope is checked only for an effective write.

## Feature requirements

The existing feature requirements retain their meanings. #48 adds:

| Feature | Meaning |
| --- | --- |
| `issue-types` | The issue store's native Issue Type catalogue and issue value can be discovered |
| `issue-fields` | The issue store's organisation Issue Field catalogue, Project attachment and issue value can be discovered |

Bindings require:

| Binding | Required features |
| --- | --- |
| Project custom field | `project-custom-fields`, `project-item-membership` |
| Project Status | `project-item-membership` |
| Issue Type | `issue-types` |
| Issue Field | `issue-fields` |
| Assignees | `issues` |
| Parent | `issue-relationships` |
| Sub-issues progress | `issue-relationships`, `project-item-membership` |

A feature declaration is a requirement to probe, not capability evidence.

## Mutation authority

#48 adds exactly:

| Scope | Maximum authorised effect |
| --- | --- |
| `issue.type.write` | Set or clear the declared Class through the resolved native Issue Type binding on one issue |
| `issue.field.write` | Set or clear the declared Priority through the resolved organisation Issue Field binding on one issue |

The existing `project.item.field.write` scope now includes bound Project custom fields and built-in Status. It remains restricted to declared target fields and bindings.

None of these scopes authorises definition creation, option replacement, rename, deletion, Project attachment, assignee changes or relationship changes.

## Static semantic diagnostics

| Diagnostic | Condition |
| --- | --- |
| `contract.semantic.missing-dimension-binding` | A declared dimension is absent on a target |
| `contract.semantic.undeclared-dimension-binding` | A target binds a dimension absent from `spec.dimensions` |
| `contract.semantic.unknown-dimension-field` | `fieldRef` does not resolve in the target |
| `contract.semantic.dimension-field-type` | Bound field type is incompatible with the dimension |
| `contract.semantic.dimension-value-map` | Canonical keys do not exactly match the declared value set |
| `contract.semantic.unknown-dimension-value-ref` | A project read/write value reference is absent from the field |
| `contract.semantic.unreadable-write-value` | A write value is absent from its canonical read set |
| `contract.semantic.overlapping-project-value-map` | One field value reference maps to two canonical values |
| `contract.semantic.overlapping-native-value-map` | One exact provider name maps to two canonical values |
| `contract.semantic.incomplete-project-value-map` | A declared bound field option is not mapped |
| `contract.semantic.duplicate-dimension-authority` | Two dimensions bind one provider field |
| `contract.semantic.native-binding-requires-organization` | Issue Type or Issue Field is bound on a non-organisation issue store |
| `contract.semantic.missing-dimension-feature` | A binding lacks a required feature declaration |
| `contract.semantic.issue-type-scope-without-binding` | `issue.type.write` has no Issue Type binding |
| `contract.semantic.issue-field-scope-without-binding` | `issue.field.write` has no Issue Field binding |
| `contract.semantic.project-field-scope-without-surface` | `project.item.field.write` has no declared Project field or Status binding |

Runtime diagnostics include:

| Diagnostic | Condition |
| --- | --- |
| `dimension.configuration.missing` | Required provider configuration is absent |
| `dimension.configuration.unknown` | Configuration discovery is incomplete |
| `dimension.configuration.unavailable` | Configuration cannot currently be retrieved |
| `dimension.configuration.forbidden` | Configuration read is not authorised |
| `dimension.configuration.unsupported` | Binding kind is unsupported |
| `dimension.configuration.ambiguous` | Exact configuration selector returns more than one identity |
| `dimension.issue-type.disabled` | The exact native type exists but is disabled |
| `dimension.value.unmapped` | Known provider value has no canonical mapping |
| `dimension.value.ambiguous` | Known provider value maps to multiple canonical values |
| `dimension.value.invalid` | Requested canonical value is not declared |
| `dimension.mutation.forbidden` | An effective write lacks its exact scope |
| `dimension.readback.mismatch` | Readback does not equal requested canonical state |

## Decision corpora

The #48 corpus contains:

- contract fixtures under `testdata/contracts/v1/dimensions/`;
- exact configuration discovery cases at `testdata/dimensions/v1/discovery-cases.json`;
- exact canonical read cases at `testdata/dimensions/v1/read-cases.json`;
- exact set, clear, no-op and forbidden write cases at `testdata/dimensions/v1/write-cases.json`;
- semantic contract diagnostics at `testdata/contracts/v1/dimensions/semantic-cases.json`.

Every provider locator, name and value is synthetic.

## Invariant register

| ID | Rule |
| --- | --- |
| `DM-001` | Dimension declaration and binding containers are explicit even when empty. |
| `DM-002` | Every declared dimension has one binding on every target. |
| `DM-003` | Canonical references are identity; provider names are selectors. |
| `DM-004` | Class uses a declared subset of the finite v1 registry. |
| `DM-005` | Priority is exactly p0, p1 and p2. |
| `DM-006` | Status is exactly todo, in-progress and done. |
| `DM-007` | Workstream values are contract-wide and target-independent. |
| `DM-008` | Each dimension has one closed set of permitted binding kinds. |
| `DM-009` | One provider field cannot authorise two dimensions. |
| `DM-010` | Native aliases are finite exact lists, never implicit case folding. |
| `DM-011` | Every writable selector is also readable. |
| `DM-012` | Read maps for different canonical values do not overlap. |
| `DM-013` | Sub-project remains label-based and outside dimension bindings. |
| `DM-014` | Parent is the native relationship authority. |
| `DM-015` | Sub-issues progress is derived and read-only. |
| `DM-016` | Native catalogue support is discovered, never inferred from owner kind. |
| `DM-017` | Provider IDs are observed state and absent from shared contracts. |
| `DM-018` | Missing configuration differs from a clear task value. |
| `DM-019` | Idempotent semantic writes are no-ops. |
| `DM-020` | Effective writes require the exact binding-family scope and readback. |

## Superseded earlier wording

This specification intentionally refines #13's statement that all built-in fields are unsupported as mappings. Built-in Status, native assignees, native hierarchy metadata, Issue Types and Issue Fields are now explicit semantic bindings. They are not added to the custom `fields` map.

The #13 exact custom-field selectors and create/adopt rules remain authoritative for `project_field` bindings. #47 remains authoritative for Project routing labels, sub-project labels and general labels.
