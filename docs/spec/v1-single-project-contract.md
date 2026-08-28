# v1 single-Project repository contract

Status: **normative pre-v1 specification**

Issue: #13

This document fixes the shared wire contract for one participating repository, one explicit GitHub issue store and one explicit GitHub Project. It refines the [v1 conceptual model](v1-conceptual-model.md) without defining a Go loader, dispatcher routing, private operator profiles or version migration.

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Frozen decisions

| Question | v1 decision |
| --- | --- |
| Canonical path | `.projectctl/project.yaml`, relative to the participating repository root |
| Text encoding | UTF-8 without a byte-order mark |
| YAML version and schema | YAML 1.2.2, Core Schema, restricted to the JSON-compatible data model described below |
| JSON Schema | Draft 2020-12 at `schemas/v1/repository-contract.schema.json` |
| Version discriminator | `apiVersion: projectctl/v1` |
| Topology discriminator | `kind: SingleProjectRepository` |
| Contract references | Lowercase ASCII kebab-case, 1 to 63 characters, exact and case-sensitive |
| Project identity input | Explicit owner kind, owner login and positive Project number |
| Provider IDs | Discovered at runtime and never stored in this shared contract |
| Project title | Display-only observed state, never a selector or fallback |
| Custom fields | `date`, `number`, `text`, `singleSelect` and `multiSelect` |
| Iteration and built-ins | Iteration unsupported; standard built-ins use #48 semantic bindings, never the custom-field map |
| Unknown keys | Rejected at every object boundary |
| Null | Rejected everywhere |
| Defaults | None in this base schema |
| Labels | Explicit declarations and an explicit Project-routing mode |
| Mutations | Exhaustive opt-in scopes; omission of a scope forbids that write |
| Safety | Dry-run, stale checks, owned-field writes and readback are unconditional and not configurable |

The discriminator literal is part of the v1 wire format. #34 owns loader dispatch, compatibility and migration behaviour around versioned documents; it does not silently change this literal or add defaults to this schema.

## Complete shape

The root contains exactly three keys:

```yaml
apiVersion: projectctl/v1
kind: SingleProjectRepository
spec:
  ref: alpha
  target:
    ref: primary
    issueStore:
      owner:
        kind: organization
        login: example-org
      repository: alpha-tasks
    project:
      owner:
        kind: organization
        login: example-org
      number: 7
  fields: {}
  dimensions: {}
  dimensionBindings: {}
  labels:
    declarations: {}
    projectRouting:
      kind: repository_scope
  sourceRefs: []
  requirements:
    features:
      - issue-labels
      - issues
      - projects-v2
    mutations: []
```

No key in this example is defaulted. Every shown container is required, including an empty `fields`, `labels.declarations`, `sourceRefs` or `mutations` collection. Label declarations and routing modes are defined by [v1 labels and sub-projects](v1-labels-and-subprojects.md).

### `spec.ref`

`spec.ref` is the stable contract reference for this participating repository declaration. It is not a repository name, checkout directory, GitHub node ID or display title. A private registry may associate this reference with operator-local information, but that association cannot change the shared contract's meaning.

### `spec.target`

The singleton target has its own stable `ref`. Keeping a target reference in the single-Project form lets #14 reuse the target object without synthesising an identifier during migration to a dispatcher.

The target names two independent resources:

- `issueStore` is the GitHub repository that stores shared issues;
- `project` is the GitHub Project in which those issues may have membership and field values.

The repository containing `.projectctl/project.yaml`, the issue store and the Project owner MAY all be different. The contract therefore has no `self` shorthand and no inferred owner. Repetition is intentional because it removes environment-dependent meaning.

Both owners require:

- `kind`, exactly `organization` or `user`;
- `login`, supplied explicitly and passed to provider discovery without case folding.

`issueStore.repository` is the provider repository locator within the declared owner. `project.number` is a positive JSON integer no greater than `9007199254740991`, preserving exact interoperability across common JSON implementations.

The contract MUST NOT contain repository IDs, Project node IDs, field IDs, option IDs, Project URLs, Project titles, default branches or permissions. Discovery resolves and verifies those facts. A provider response that disagrees with the declared owner kind is a typed discovery conflict, not permission to reinterpret the contract.

## Reference syntax and namespaces

Every contract reference uses:

```text
^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$
```

The length is 1 to 63 Unicode code points; the pattern restricts actual values to lowercase ASCII. Comparison is exact byte-for-byte comparison of the decoded ASCII value. Implementations MUST NOT trim, case-fold, Unicode-normalise, pluralise or generate a reference from display text.

Reference namespaces are deliberately scoped:

| Location | Namespace |
| --- | --- |
| `spec.ref` | The contract selected by a caller or private registry |
| `spec.target.ref` | Targets inside this contract; singleton in this kind |
| Each key in `spec.fields` | Fields inside the target Project |
| Each key in a field's `values` | Values inside that field only |
| Each key in `spec.labels.declarations` | Labels inside this contract |
| Each item in `spec.sourceRefs` | Repository-safe source references exposed to integrations |

The same spelling in two different namespaces does not create a relationship. Map keys provide uniqueness for field and value references; duplicate YAML keys are rejected before schema validation.

## Custom-field mappings

`spec.fields` is a map from stable contract field references to explicit GitHub custom-field declarations. Every declared mapping is required structure for this contract. An empty map means the contract does not require or address any custom field. Standard semantic meaning is attached only through [v1 project dimensions](v1-project-dimensions.md); a field name never acquires Class, Priority, Workstream or Due-date meaning by itself.

### Scalar fields

```yaml
fields:
  target-date:
    name: Target date
    dataType: date
  estimate:
    name: Estimate
    dataType: number
  notes:
    name: Notes
    dataType: text
```

Scalar fields contain exactly `name` and `dataType`. A `values` key is invalid.

### Select fields and value mappings

```yaml
fields:
  work-state:
    name: Work state
    dataType: singleSelect
    values:
      queued:
        name: Queued
        color: gray
        description: Ready to be scheduled.
      active:
        name: Active
        color: blue
        description: Work is in progress.
  areas:
    name: Areas
    dataType: multiSelect
    values:
      api:
        name: API
        color: purple
        description: API-facing work.
```

Both select types require at least one value. A value contains exactly:

- `name`, the exact provider option selector and create-time name;
- `color`, one of `blue`, `gray`, `green`, `orange`, `pink`, `purple`, `red` or `yellow`;
- `description`, a required single-line string that MAY be empty.

The adapter maps lower-case contract colours to GitHub's corresponding upper-case enum. Colour and description are complete create-time input because GitHub requires both when creating select options. They are not option identity and are not continuously reconciled. If an exact-name option already exists, bootstrap and adoption preserve its current colour and description. A missing option may be created only when `project.field.option.create` is declared.

### Supported field types

The v1 schema supports the custom-field types that have fixed scalar or fixed option-map meaning:

| Contract value | Meaning |
| --- | --- |
| `date` | GitHub date custom field |
| `number` | GitHub number custom field |
| `text` | GitHub text custom field |
| `singleSelect` | One value reference resolves to one option |
| `multiSelect` | A set of distinct value references resolves to options |

Iteration fields are intentionally rejected. Their rolling schedule, iteration identifiers and update semantics are live configuration, not a fixed value map. Supporting them requires a later contract version with explicit discovery, comparison and mutation rules. Built-in and issue-owned fields are not custom-field mappings. #48 defines explicit bindings for built-in Status, assignees, hierarchy metadata, native Issue Types and organisation Issue Fields. Declared issue-label membership uses the separate label scopes and MUST NOT be written through `project.item.field.write`.

### Exact provider selection

The contract reference is protocol identity. `name` and `dataType` are explicitly designated provider selectors within the resolved Project, which is the narrow exception permitted by `CM-006`.

Discovery MUST resolve a declared field as follows:

1. fetch all supported field configurations in the exact Project scope;
2. retain fields whose returned name is exactly equal to `name` and whose type equals `dataType` after the adapter's fixed enum mapping;
3. return the one matching provider identity;
4. report not found for zero matches and ambiguity for more than one match.

For select values, discovery fetches the selected field's options, compares returned names exactly, and likewise requires one match. Provider search filters that are case-insensitive MUST NOT be treated as proof of identity. Implementations MUST NOT trim, case-fold, perform fuzzy matching, use display order or fall back to a same-name field of another type.

The following semantic conflicts use stable diagnostics even though JSON Schema cannot compare values across map entries:

| Diagnostic | Condition |
| --- | --- |
| `contract.semantic.duplicate-field-selector` | Two field references declare the same exact `(name, dataType)` selector |
| `contract.semantic.duplicate-value-selector` | Two value references within one field declare the same exact option `name` |

No two references may intentionally alias one provider field or option in v1.

## Repository-safe source hooks

`spec.sourceRefs` is a YAML sequence whose order has no semantic meaning. Its contract references are unique and it may be empty. #34 owns the canonical ordering used by normalised output.

A source reference only advertises a collaborator-safe symbolic hook. It does not include a path, URL, provider, document ID, sheet range, credential, access instruction or promise that every collaborator can retrieve the source. Repository-local source catalogues and private profile bindings are defined by later integration work. Unknown source locations and private destinations cannot be smuggled into this object because it accepts only reference strings.

## Requirements and mutation authority

`spec.requirements` always contains explicit `features` and `mutations` arrays. Both arrays are sets: duplicate values are invalid and declaration order has no semantic meaning. #34 owns their canonical output ordering.

### Features

Every contract requires `issues` and `projects-v2`. Other values are opt-in requirements:

| Feature | Meaning |
| --- | --- |
| `issues` | The declared issue store can be read as the shared issue authority |
| `projects-v2` | The declared Project can be discovered and read |
| `project-custom-fields` | Declared custom fields and options can be discovered |
| `project-item-membership` | Issue membership in the Project can be discovered |
| `issue-relationships` | Supported issue relationships can be discovered |
| `issue-labels` | Complete repository label catalogues and issue-label membership can be discovered |
| `issue-types` | Native Issue Types and an issue's type value can be discovered |
| `issue-fields` | Organisation Issue Fields, Project attachment and issue values can be discovered |

A feature declaration is a requirement to probe, not evidence that the acting principal has it. Capability and permission remain observed runtime facts. Missing, forbidden, unavailable and unsupported results remain distinct.

Any non-empty `fields` map requires `project-custom-fields`.

### Mutations

The mutation set is an exhaustive allowlist for plans made under this contract:

| Mutation | Maximum authorised effect |
| --- | --- |
| `issue.create` | Create a new issue in the exact issue store from separately validated intent |
| `issue.relationship.create` | Create one supported relationship between explicitly resolved issues |
| `issue.type.write` | Set or clear Class through one declared native Issue Type binding |
| `issue.field.write` | Set or clear Priority through one declared organisation Issue Field binding |
| `repository.label.create` | Create one missing declared label with its exact create-time attributes |
| `issue.label.add` | Add one exact declared label to one resolved issue |
| `issue.label.remove` | Remove one exact declared label from one resolved issue |
| `project.field.create` | Create one missing declared custom field with its declared select options |
| `project.field.option.create` | Add one missing declared option while preserving every observed existing option identity and value |
| `project.item.add` | Add one explicitly resolved issue to the exact Project |
| `project.item.field.write` | Set or explicitly clear one declared custom-field dimension or built-in Status value on one resolved Project item |

No v1 mutation authorises deleting resources, removing Project membership, renaming or replacing fields, rewriting existing issue title/body, updating or deleting repository labels, using a set-all-labels endpoint, editing option colour/description, or updating built-in metadata other than the explicitly bound Status value.

An absent mutation scope forbids planning or execution of that write. An empty mutation array is a valid read-only contract. Mutation scopes apply only to the declared target and declared mappings; they do not grant organisation-wide authority.

The following prerequisites are schema invariants:

| Declaration | Required declaration |
| --- | --- |
| Non-empty `fields` | Feature `project-custom-fields` |
| `issue.relationship.create` | Feature `issue-relationships` |
| Any non-empty label declaration set | Feature `issue-labels` |
| `repository.label.create` | Non-empty label declarations and feature `issue-labels` |
| `issue.label.add` | Non-empty label declarations and feature `issue-labels` |
| `issue.label.remove` | Non-empty label declarations and feature `issue-labels` |
| `project.field.create` | Non-empty `fields` and feature `project-custom-fields` |
| `project.field.option.create` | Non-empty `fields` and feature `project-custom-fields` |
| `project.item.add` | Feature `project-item-membership` |
| `project.item.field.write` | Non-empty `fields`, feature `project-custom-fields` and feature `project-item-membership` |

`issue.create` needs no feature beyond the mandatory `issues` requirement.

### Safety is not configurable

Mutation declarations can reduce authority but cannot reduce safety. There are deliberately no `dryRun`, `force`, `staleGuard`, `verify`, `continueOnError` or equivalent fields. Every implementation still MUST:

- produce the complete plan before writes;
- require an explicit apply request;
- identify owned fields and expected prior state;
- freshly re-read stale-sensitive state;
- refuse mismatch, unknown or unavailable state;
- stop after the first unsafe refusal or provider failure;
- read back every accepted mutation;
- report success only for matching observed after-state.

Permission to use a mutation scope is still limited by the acting principal's discovered provider permissions and the requested workflow's policy.

## Presence, null and empty rules

There are no implicit values in this base document.

| Path or category | Presence | Null | Empty behaviour |
| --- | --- | --- | --- |
| Root `apiVersion`, `kind`, `spec` | Required | Rejected | Not applicable |
| `spec.ref`, `target` | Required | Rejected | Empty reference rejected |
| Owners, repository and Project number | Required | Rejected | Empty strings and number zero rejected |
| `fields` | Required | Rejected | Empty map valid and means no custom-field mappings |
| `labels.declarations` | Required | Rejected | Empty map valid only when the selected routing mode permits it |
| `labels.projectRouting` | Required | Rejected | Exactly one explicit routing mode |
| Field `name`, `dataType` | Required | Rejected | Empty string rejected |
| Select `values` | Required for select, forbidden otherwise | Rejected | Empty map rejected |
| Option `name`, `color`, `description` | Required | Rejected | Only `description` may be the empty string |
| `sourceRefs` | Required | Rejected | Empty array valid and means no advertised source hooks |
| `requirements.features` | Required | Rejected | Must contain at least `issues` and `projects-v2` |
| `requirements.mutations` | Required | Rejected | Empty array valid and means read-only |
| Any unknown key | Rejected | Not applicable | Not applicable |

Omission is not equivalent to an empty value because required collections must be present. #34 may define normalisation for later extensions, but it cannot treat a schema-invalid omitted or null base field as valid v1 input.

## YAML input profile

The canonical file uses [YAML 1.2.2](https://yaml.org/spec/1.2.2/) with these additional restrictions for portable and safe parsing:

1. the byte stream is UTF-8 and has no byte-order mark;
2. the stream contains exactly one document;
3. mappings have string keys that are unique in their mapping;
4. only Core Schema null, Boolean, integer, floating-point and string scalar types are recognised;
5. every value must be representable in the JSON data model accepted by the schema;
6. directives, explicit custom tags, anchors, aliases and merge keys are rejected;
7. non-finite floating-point values are rejected;
8. comments are permitted and have no semantic meaning.

These checks occur during parsing before JSON Schema validation. A parser MUST retain enough presence and source-location information to report the offending document. It MUST NOT silently keep the first or last duplicate mapping value, expand an alias, merge mappings or coerce an unsupported tagged value.

## JSON Schema contract

The machine schema uses [JSON Schema Draft 2020-12](https://json-schema.org/draft/2020-12). It is the structural authority for required keys, types, constants, reference syntax, enum values, collection limits, mutation prerequisites and unknown-key rejection.

Validation operates on the JSON-compatible representation produced by the restricted YAML parse. Schema validation MUST NOT mutate that value, insert defaults, discover provider state or perform the cross-entry selector checks listed above.

Non-Go implementations may validate the same representation directly. An implementation may compile the schema into code, but the compiled form is derived and MUST produce equivalent acceptance results for the frozen fixtures.

## Invalid-fixture manifest

`testdata/contracts/v1/single-project/invalid/manifest.json` assigns each invalid fixture:

- a stage, one of `yaml`, `schema` or `semantic`;
- one stable diagnostic key;
- a canonical JSON Pointer instance path when parsing reached an instance;
- the conceptual or single-Project rule exercised.

Fixtures are minimal. A validator MAY report additional low-level details, but it MUST include the listed stable diagnostic and path. Diagnostics are ordered first by stage, then path, then key. Human wording is not stable protocol.

## Invariant register

| ID | Rule |
| --- | --- |
| `SP-001` | The sole canonical shared contract path is `.projectctl/project.yaml`. |
| `SP-002` | The wire stream follows the restricted single-document UTF-8 YAML profile. |
| `SP-003` | `apiVersion` is exactly `projectctl/v1`. |
| `SP-004` | `kind` is exactly `SingleProjectRepository`. |
| `SP-005` | Every base container is explicit; this schema defines no defaults. |
| `SP-006` | Contract references use the frozen lowercase ASCII syntax and exact comparison. |
| `SP-007` | Participating contract, target, field, value and source references occupy scoped namespaces. |
| `SP-008` | Issue store and Project are independently and explicitly identified. |
| `SP-009` | Owner kind and login are explicit for both issue store and Project. |
| `SP-010` | Project number is positive and JSON-interoperable. |
| `SP-011` | Provider IDs, Project titles, local paths, permissions and private destinations are absent from the shared contract. |
| `SP-012` | Declared field references map to exact name and supported type selectors inside one resolved Project. |
| `SP-013` | Declared select value references map to exact option names inside one resolved field. |
| `SP-014` | Provider selection does not trim, case-fold, fuzzy-match or use display order. |
| `SP-015` | Duplicate field or option selectors are semantic conflicts, not aliases. |
| `SP-016` | Select option colour and description are create-time attributes, not identity or continuous reconciliation intent. |
| `SP-017` | Iteration is unsupported; built-in and issue-owned semantics use #48 bindings rather than custom-field mappings. |
| `SP-018` | Source hooks contain symbolic collaborator-safe references only. |
| `SP-019` | Feature declarations are requirements to probe, never observed capability claims. |
| `SP-020` | Mutation declarations are an exhaustive allowlist scoped to declared targets and mappings. |
| `SP-021` | Delete, destructive replacement and broad update scopes do not exist in v1. |
| `SP-022` | Dry-run, stale protection, owned-field writes and readback verification remain mandatory regardless of declarations. |
| `SP-023` | Null and unknown keys are rejected everywhere. |
| `SP-024` | Mapping order and set-member declaration order carry no semantic meaning; #34 owns canonical output ordering. |
| `SP-025` | Label declarations use their own reference namespace and explicit role. |
| `SP-026` | Issue-label writes use narrow create, add and remove scopes, never Project field writes or set-all replacement. |

## Explicitly deferred

This base specification is extended, without hidden defaults, by:

- dispatcher topology and selection in [v1 routing](v1-routing.md);
- label roles, route-label cardinality and sub-projects in [v1 labels and sub-projects](v1-labels-and-subprojects.md);
- shared privacy advertisement in [v1 shared privacy](v1-shared-privacy.md);
- private destinations, local paths and effective choices in [v1 operator profile](v1-operator-profile.md).

Still deferred are loader compatibility, canonical output, migrations and tool-version vocabulary (#34); final cross-contract conformance fixtures (#16); Go wire and canonical types (#18 and #33); and provider discovery, planning and execution APIs.

No downstream issue may add a fallback selector, hidden mutation scope or safety toggle to fill a deferred gap.

## Conformance checklist

A v1 single-Project document conforms only if:

- it exists at the canonical path and satisfies the restricted YAML profile;
- it validates against the Draft 2020-12 schema without mutation;
- its two semantic selector-uniqueness rules pass;
- it carries no provider IDs, private values or local paths;
- all field and value matching can be performed exactly within the declared target;
- its label declarations, roles and routing mode satisfy the label specification;
- all required feature and mutation prerequisites are explicit;
- a read-only document can use an empty mutation set;
- no document value can disable mandatory planning, stale checking or verification.
