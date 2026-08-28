# v1 automated setup and manual-action outcomes

Status: **normative pre-v1 specification**

Issue: #49

This document defines create-once Project initialisation, bootstrap, non-destructive adoption, setup snapshots, outcome aggregation and manual-action reporting for the v1 shared protocol. It extends the [v1 conceptual model](v1-conceptual-model.md), [single-Project contract](v1-single-project-contract.md), [routing contract](v1-routing.md), [labels and sub-projects](v1-labels-and-subprojects.md) and [project dimensions](v1-project-dimensions.md).

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Product boundary

Automated setup is a primary product outcome. It is not continuous reconciliation.

V1 has three distinct setup modes:

| Mode | Input | Purpose | May change provider state |
| --- | --- | --- | --- |
| `initialise` | One setup blueprint | Create new Projects, read back their assigned identities and materialise a candidate repository contract | Yes, only the explicitly requested new Projects and initial repository links |
| `bootstrap` | One validated repository contract | Create missing contract-owned structure once | Yes, only missing resources authorised by the contract |
| `adopt` | One validated repository contract and explicit adoption choices | Preserve compatible live structure and add only missing contract-owned structure | Yes, only authorised missing additions |

An ordinary task operation is not setup. Migration is not setup. A later run does not acquire permission to restore mutable task state merely because setup originally created it.

## Frozen decisions

| Question | V1 decision |
| --- | --- |
| Existing Project identity | Exact owner and positive Project number from the repository contract |
| New Project identity | Assigned by GitHub and read back after a blueprint creation operation |
| Project title | Create-time display input only, never an identity selector |
| Blueprint role | One-off pre-contract intent, not an operational repository contract |
| Blueprint materialisation | Replace each create descriptor with its verified owner and assigned number; preserve the remaining validated intent |
| Bootstrap | Create missing exact declarations once; never rewrite a conflict |
| Adoption | Preserve explicit compatible matches and propose contract-selector amendments when display selectors differ |
| Setup declarations | Optional `spec.setup`; absence means no links, views, workflow expectations or seeds beyond the base contract surfaces |
| Setup target keys | Exact target references; their key set must equal the contract target-reference set when `setup` is present |
| Views | Create-only from exact declared configuration; existing conflicts are never updated |
| Built-in workflows | Exact shared expectation plus manual action, because current API discovery cannot read or create their configuration |
| Native organisation definitions | Existing Issue Types and Issue Fields may be used; creating or updating organisation-wide definitions is not authorised in v1 |
| Issue Field attachment | A declared organisation Issue Field may be attached to an exact Project under its own scope |
| Seeds | Create-once issues with a deterministic collaborator-safe marker |
| Existing seed content | Live shared state; title, body, labels, state and dimension values are preserved after creation |
| Manual action | A typed incomplete setup item, never provider-verified success |
| Apply success | Matching targeted provider readback only |
| Extra live structure | Preserved and reported as unmanaged observed state |

## Why Project initialisation is separate

The steady-state repository contract identifies a Project by exact owner and positive Project number. GitHub assigns that number only after creation. A title cannot fill the gap because titles are display values and may be duplicated or renamed.

V1 therefore does not weaken the repository contract to allow title lookup. A setup blueprint uses a create descriptor containing an owner and create-time display values. The initialiser creates each Project without searching for a same-title candidate, reads back its assigned number and emits a candidate steady-state contract.

The candidate is not authoritative until it is published through the repository's normal review path. A create response without readback, an interrupted create with unknown result, or a candidate that was not published is explicit incomplete state.

## Setup blueprint

The machine shape is defined by `schemas/v1/setup-blueprint.schema.json`. A single-Project blueprint has this form:

```yaml
apiVersion: projectctl/v1
kind: SingleProjectSetupBlueprint
spec:
  ref: sample-work
  target:
    ref: primary
    issueStore:
      owner: {kind: organization, login: sample-org}
      repository: tasks
    project:
      owner: {kind: organization, login: sample-org}
      title: Sample work
      shortDescription: Shared planning for the sample work.
      visibility: private
  fields: {}
  dimensions: {}
  dimensionBindings: {}
  labels:
    declarations: {}
    projectRouting: {kind: repository_scope}
  setup:
    targets:
      primary:
        repositoryLinks:
          - owner: {kind: organization, login: sample-org}
            repository: tasks
        views: {}
        workflows: {}
    seeds: {}
  sourceRefs: []
  requirements:
    features: [issues, project-repository-links, projects-v2]
    mutations: [project.repository.link]
```

Dispatcher blueprints contain two or more target blueprints and retain the exact dispatcher routing rules. All normal contract references, labels, fields, dimensions, bindings, privacy and requirement rules still apply.

### Initialisation rules

1. Validate the complete blueprint before provider access.
2. Resolve every explicit owner and issue store without searching by Project title.
3. Require explicit operator authorisation for `initialise` and Project creation.
4. Create one target Project at a time.
5. Read back owner, assigned number, title, visibility and linked repository state.
6. Record each verified locator before starting the next target.
7. Materialise a candidate repository contract only from verified locators.
8. Emit a manual action to review and publish that candidate.

An interrupted operation whose create result is unknown MUST NOT retry blindly. It returns `setup.initialise.create-result-unknown` with an exact recovery checklist. A user may resume only by supplying the recovered provider identity or explicitly authorising a fresh creation after proving no prior Project was created.

### Materialisation

Materialisation changes only:

- `kind: SingleProjectSetupBlueprint` to `SingleProjectRepository`, or `DispatcherSetupBlueprint` to `DispatcherRepository`;
- each blueprint `project` object to `{owner, number}` using verified readback.

Create-time title, short description and visibility are not retained as steady-state selectors. Every other supplied field is copied byte-for-byte into the candidate semantic value. The serialised file may use the canonical formatting later defined by #34.

## Shared setup declarations

`spec.setup` is optional shared intent. When present it contains explicit `targets` and `seeds` containers. Empty containers are valid.

```yaml
setup:
  targets:
    primary:
      repositoryLinks:
        - owner: {kind: organization, login: sample-org}
          repository: tasks
      views:
        backlog:
          name: Backlog
          layout: table
          filter: "is:issue is:open"
          visibleFields:
            - {kind: builtin, name: title}
            - {kind: dimension, ref: priority}
            - {kind: dimension, ref: status}
          sortBy:
            - field: {kind: dimension, ref: priority}
              direction: asc
          groupBy: []
          verticalGroupBy: []
      workflows:
        add-open-issues:
          kind: auto_add
          name: Auto-add open issues
          enabled: true
          repositories:
            - owner: {kind: organization, login: sample-org}
              repository: tasks
          filter: "is:issue is:open"
  seeds:
    project-overview:
      targetRef: primary
      title: Project overview
      body: Track the main outcomes for this work.
      labelRefs: []
      dimensionValues:
        class: epic
        priority: p1
```

When `setup` is absent, base declarations such as labels, custom fields and dimension bindings remain valid contract intent and may still be bootstrapped under their existing scopes. Absence adds no repository links, views, workflows or seed issues.

## Target setup declarations

### Repository links

`repositoryLinks` is an explicit set of repositories that should be linked to the exact Project. The participating repository and issue store are not inferred as links. Extra observed links are preserved.

A missing declared link may be created only with `project.repository.link`. No v1 scope unlinks a repository.

GitHub exposes Project repository discovery and the `linkProjectV2ToRepository` mutation through the [Projects GraphQL API](https://docs.github.com/en/graphql/reference/projects).

### Views

Views are a map from stable contract reference to exact create-once configuration. A view declares:

- exact `name`;
- `layout`, one of `table`, `board` or `roadmap`;
- exact GitHub Project `filter` text, which may be empty;
- ordered visible field selectors;
- ordered sort clauses;
- zero or one horizontal group field;
- zero or one board column field.

Field selectors are explicit:

- `builtin` uses one finite built-in token;
- `field` references a declared custom field on that target;
- `dimension` references a declared semantic dimension and resolves through that target's binding.

The adapter resolves selectors to observed field identities before planning. It does not match by a rendered field title outside the declared binding. A roadmap view has no `visibleFields` configuration. A non-board view has no `verticalGroupBy` configuration.

GitHub's current GraphQL and REST APIs support creating Project views, including layout, filters and field configuration. V1 therefore treats view creation as automatable when `project-views` is supported and `project.view.create` is authorised. See the [Projects GraphQL reference](https://docs.github.com/en/graphql/reference/projects) and [Project views REST API](https://docs.github.com/en/rest/projects/views).

Parent-oriented views resolve the declared `parent` dimension to GitHub's Parent issue field and may group by that field. A filter for one already known parent may use GitHub's exact `parent-issue:OWNER/REPOSITORY#NUMBER` syntax, but a contract never substitutes a seed title or unknown future issue number into a filter. Sub-project views use the exact declared label name in the opaque GitHub filter string. See GitHub's [parent issue fields](https://docs.github.com/en/issues/planning-and-tracking-with-projects/understanding-fields/about-parent-issue-and-sub-issue-progress-fields) and [Project filter syntax](https://docs.github.com/en/issues/planning-and-tracking-with-projects/customizing-views-in-your-project/filtering-projects).

Views remain create-once:

- exact existing configuration is `already_conformant`;
- a compatible exact-selector match that includes extra provider display state is `preserved_compatible`;
- a same-name view with different layout, filter, fields, sort or grouping is `conflict`;
- a missing view is planned for creation when authorised;
- no v1 operation updates, renames or deletes a view.

### Built-in workflows

V1 declares only an `auto_add` expectation. It contains an exact human-facing name, enabled state, one or more explicit repositories and an exact filter.

GitHub documents built-in auto-add workflows, but the current Projects GraphQL surface exposes only workflow identity, name and enabled state. It exposes no mutation to create or configure an auto-add workflow and no readable repository/filter configuration. The only workflow mutation currently documented is deletion. See [Automating your project](https://docs.github.com/issues/planning-and-tracking-with-projects/automating-your-project) and the [Projects GraphQL reference](https://docs.github.com/en/graphql/reference/projects).

Every declared auto-add workflow therefore produces `manual_action_required` or `manual_action_acknowledged`. Name and enabled state alone never prove configuration. The checklist includes the exact Project, repositories, filter and desired enabled state.

An operator acknowledgement is evidence that the checklist was performed. It is not provider readback and MUST NOT be labelled verified.

## Seed issues

Seeds are contract-reference keyed create-once issue definitions. Each seed declares one target reference, create-time title and body, declared label references, optional initial writable dimension values and an optional parent seed reference.

The deterministic marker is:

```text
<!-- projectctl-seed:v1:<contract-ref>:<seed-ref> -->
```

The marker is appended to the create-time body with one blank line before it. It is collaborator-safe setup metadata and is the sole provider-side seed selector. Titles never identify seeds.

Discovery scans only the explicit issue store and requires zero or one exact marker match. More than one is `setup.seed.duplicate-marker`. A marker on a pull request or in another repository does not match.

On creation, the planner may create the issue, add its target Project membership, apply declared create-time labels and writable dimension values, and create the declared parent relationship, but only under the relevant mutation scopes. Dependencies are explicit and ordered.

After the marker-bearing issue exists:

- title, body text other than the marker, open/closed state, labels, assignees, Project values and ordinary edits are live shared state;
- bootstrap and adoption preserve them even when they differ from the seed definition;
- exact Project membership and the declared parent relationship remain separate create-once structural items and may be added when known missing on an explicit later setup run;
- removal of the marker relinquishes automatic seed identity, so a later explicitly authorised bootstrap may propose a new seed issue.

The setup snapshot and report schemas may carry observed provider identities. They are runtime observed state, not shared repository-contract authority. V1 has no hidden local or provider-side setup receipt.

## Derived agent guidance

Generated concise guidance is a setup outcome but not shared contract authority. #51 owns its intermediate model, renderer semantics, output paths and privacy proofs. #49 fixes only the setup boundary:

- the stable item key is `guidance/projection`;
- desired state is the exact ordered output set and SHA-256 input/content digests produced under #51;
- observed state reads the named repository paths at one exact revision;
- matching bytes are `already_conformant`;
- missing or stale bytes are `contract_patch_required` with a candidate repository patch and review checklist;
- an unavailable renderer is `skipped_dependency`, never an invitation to invent guidance;
- setup never commits, pushes or merges guidance directly.

The repository review that publishes guidance is a manual action with reason `contract_patch_required`. Acknowledging the checklist does not prove that the repository bytes match. A later exact repository read is required before the item becomes conformant.

## Base setup surfaces

The repository contract already declares several create-once surfaces. #49 fixes their setup interpretation without changing their identity rules.

| Surface | Exact selector | Missing | Compatible existing | Conflict |
| --- | --- | --- | --- | --- |
| Project | Owner and positive number | Blocked; use initialisation and publish a new locator | Preserve | Wrong owner kind or ambiguous observation |
| Repository label | Exact declared name | Create with `repository.label.create` | Preserve current colour/description | Duplicate exact identities or inaccessible catalogue |
| Project custom field | Exact name and data type | Create with `project.field.create` | Preserve extra state | Same name with wrong type or ambiguous exact match |
| Select option | Exact name inside one field | Create with `project.field.option.create` while preserving observed identities | Preserve current colour/description | Ambiguous exact match or unsafe incomplete option snapshot |
| Native Issue Type | Exact #48 name map | Manual action if missing | Preserve | Disabled, ambiguous or overlapping mapping |
| Organisation Issue Field | Exact #48 name and type | Manual action if missing | Preserve | Wrong type, ambiguity or missing mapped values |
| Issue Field Project attachment | Exact observed Issue Field identity and Project | Attach with `project.issue-field.attach` | Preserve | Ambiguous field or forbidden Project scope |
| Built-in Status options | Exact #48 map | Manual action if mapped values are missing | Preserve aliases and extra provider states as observed | Mapping ambiguity |
| Project item membership | Exact issue and Project identities | Add with `project.item.add` | Preserve | Ambiguous issue or Project identity |
| Parent relationship | Exact child and parent issue identities | Add with `issue.parent.set` | Preserve | Different existing parent or cycle |
| Derived agent guidance | #51 output path plus exact input/content digests | Produce candidate repository patch | Preserve exact bytes | Stale bytes require a reviewed patch; renderer unavailable skips the dependency |

Organisation Issue Type and Issue Field APIs can create definitions, but those operations change organisation-wide configuration and require create-time metadata not owned by the semantic bindings. V1 deliberately omits organisation-wide definition mutation scopes. Missing definitions receive a manual administrator action rather than silently broadening one repository contract's authority. See GitHub's [Issue Types](https://docs.github.com/en/rest/orgs/issue-types) and [Issue Fields](https://docs.github.com/en/rest/orgs/issue-fields) REST APIs.

## Feature and mutation vocabulary

Add these feature requirements:

- `project-repository-links`;
- `project-views`;
- `project-workflows`.

Add these repository-contract mutation scopes:

- `project.issue-field.attach`;
- `project.repository.link`;
- `project.view.create`.

The blueprint's Project creation is authorised by the explicit initialisation request and operator action. It is not a steady-state repository-contract mutation scope.

No v1 setup scope authorises Project deletion, unlinking, view update/delete, workflow deletion, organisation definition mutation, Project metadata reconciliation, issue content reconciliation or cleanup of extra live state.

## Setup item model

Planning expands desired intent into stable item keys:

```text
target/<target-ref>/project
target/<target-ref>/repository-link/<owner>/<repository>
target/<target-ref>/label/<label-ref>
target/<target-ref>/field/<field-ref>
target/<target-ref>/field/<field-ref>/option/<value-ref>
target/<target-ref>/issue-field-attachment/<dimension-ref>
target/<target-ref>/view/<view-ref>
target/<target-ref>/workflow/<workflow-ref>
seed/<seed-ref>/issue
seed/<seed-ref>/membership
seed/<seed-ref>/label/<label-ref>
seed/<seed-ref>/dimension/<dimension-ref>
seed/<seed-ref>/parent
guidance/projection
contract/publication
```

Each item has one desired source, observed fact state, dependency list, outcome, blocking flag and verification class. Ordering of input maps never changes item keys or report order.

## Observed setup snapshot

`schemas/v1/setup-snapshot.schema.json` defines a deterministic bounded snapshot. Discovery reads only owners, repositories, Projects and issue stores explicitly named by the blueprint or contract.

Every resource fact has one state:

- `known`;
- `missing`;
- `unknown`;
- `unavailable`;
- `forbidden`;
- `unsupported`;
- `ambiguous`.

Only `known` carries typed facts. `missing` is a successful bounded lookup with no resource. Unknown and failure states never become missing. Provider identities in a snapshot are observed values and never flow into the repository contract except the verified Project number produced by blueprint materialisation.

## Per-item outcomes

Every setup item resolves to exactly one outcome:

| Outcome | Meaning | Provider verified |
| --- | --- | --- |
| `already_conformant` | Known live state already satisfies the create-once requirement | Yes |
| `planned_create` | A dry-run contains one authorised missing-resource creation | No |
| `applied_verified` | The authorised creation was applied and matching state was read back | Yes |
| `preserved_compatible` | Compatible existing state is retained without mutation | Yes for identity/compatibility only |
| `contract_patch_required` | Adoption or initialisation produced a candidate selector change that needs repository review | No |
| `conflict` | Known live state cannot safely satisfy the declaration without replacement or ambiguity | No |
| `unavailable` | The environment could not obtain a required fact | No |
| `forbidden` | Provider permission blocked required discovery or mutation | No |
| `unsupported` | The adapter/provider cannot represent the required surface | No |
| `scope_forbidden` | The contract does not authorise an otherwise possible write | No |
| `manual_action_required` | An exact human/provider action remains outstanding | No |
| `manual_action_acknowledged` | The operator attests that the checklist was performed | No; operator attestation only |
| `skipped_dependency` | An earlier required item prevented safe evaluation | No |

An accepted API response is not an outcome. It becomes `applied_verified` only after exact readback.

## Bootstrap comparison

Bootstrap uses the contract's exact selectors:

1. known exact compatible state becomes `already_conformant` or `preserved_compatible`;
2. known missing state becomes `planned_create` only when a matching mutation scope exists;
3. missing state without scope becomes `scope_forbidden`;
4. same-selector incompatible state becomes `conflict`;
5. extra state is preserved and reported outside desired items;
6. unavailable, forbidden, unsupported and ambiguous facts remain distinct;
7. manual-only surfaces emit exact manual actions;
8. dependent items are skipped after a blocking result, while independent targets may continue.

The second run against unchanged verified after-state produces no creation operations.

## Adoption comparison

Adoption never fuzzy-matches. It accepts only one of:

- the contract's existing exact selector;
- one explicit operator-supplied candidate identity within the already bounded target scope.

The planner proves the candidate's type and compatibility. If the provider display selector differs from the contract, adoption emits `contract_patch_required` containing a proposed selector amendment. It does not store a private alias, insert a provider ID into the contract or mutate the provider merely to match preferred text.

Until the patch is reviewed and published, ordinary contract resolution still uses the old selector and setup remains incomplete. Ambiguous candidates or incompatible types are conflicts. Adoption does not rewrite the contract automatically.

## Manual actions

A manual action contains:

- a stable action ID derived from the item key;
- exact target and surface;
- one reason code;
- concise desired state;
- ordered steps;
- the evidence needed to acknowledge completion;
- when acknowledged, a timestamped operator attestation statement;
- `providerVerified: false`.

Reason codes are:

- `api_unsupported`;
- `configuration_unobservable`;
- `organisation_wide_authority`;
- `scope_absent`;
- `contract_patch_required`;
- `ambiguous_recovery`.

Acknowledgement changes only the evidence class. It never creates a fabricated provider observation.

## Plan, apply and verification

Setup is dry-run by default. Planning is pure and contains the complete item graph before any write.

Apply requires explicit operator authorisation and follows this order:

1. re-read the exact stale-sensitive resource;
2. require the expected prior state;
3. perform one owned additive/create operation;
4. read back the exact resource;
5. compare only owned create-time state;
6. record `applied_verified` or a typed failure;
7. stop dependants after failure;
8. continue independent targets only when their plans do not rely on the failed item.

Executors do not re-plan, search a broader owner scope, repair conflicts or perform hidden manual actions.

## Overall report state

`schemas/v1/setup-report.schema.json` defines the machine report. Overall state is derived in this precedence order:

1. `blocked` if any non-manual blocking item is `conflict` (including ambiguity), `unavailable`, `forbidden`, `unsupported`, `scope_forbidden` or `skipped_dependency`;
2. `incomplete_manual` if any required manual action is unacknowledged or a contract patch is unpublished;
3. `changes_planned` for a dry-run containing authorised creations and no blockers;
4. `complete_with_manual_attestations` when automatic items are conformant/verified and every manual item is acknowledged;
5. `applied_verified` when at least one change was applied, all items are provider-verified or compatible and there are no manual items;
6. `already_conformant` when no mutation or manual action is required.

`complete_with_manual_attestations` is deliberately not named verified.

## Stable diagnostics

| Diagnostic | Condition |
| --- | --- |
| `setup.semantic.target-key-set` | Setup target keys differ from declared target references |
| `setup.semantic.unknown-target` | A seed or setup item references an unknown target |
| `setup.semantic.unknown-field` | A view field selector references an undeclared target field |
| `setup.semantic.unknown-dimension` | A view or seed references an undeclared dimension |
| `setup.semantic.unknown-label` | A seed references an undeclared label |
| `setup.semantic.unknown-seed-parent` | A seed parent reference is absent |
| `setup.semantic.seed-cycle` | Seed parent references contain a cycle |
| `setup.semantic.duplicate-view-selector` | Two view refs on one target use the same exact name |
| `setup.semantic.duplicate-workflow-selector` | Two workflow refs on one target use the same exact name |
| `setup.semantic.duplicate-repository-link` | One target declares the same repository link more than once |
| `setup.semantic.view-field-not-on-target` | A dimension or field cannot resolve to that Project |
| `setup.semantic.seed-value` | A seed value is not valid under the canonical dimension declaration |
| `setup.semantic.seed-not-writable` | A seed supplies a read-only or unsupported create-time dimension |
| `setup.semantic.mutation-prerequisite` | A setup mutation scope lacks its feature or desired surface |
| `setup.project.missing` | A numbered steady-state Project was not found |
| `setup.initialise.create-result-unknown` | Project creation may have succeeded but verified identity is unavailable |
| `setup.seed.duplicate-marker` | More than one issue contains one exact seed marker |
| `setup.view.conflict` | An exact-name view differs in create-once configuration |
| `setup.workflow.manual` | Workflow configuration requires a manual checklist |
| `setup.verification.mismatch` | Readback differs from the owned desired after-state |

## Conformance corpora

The repository contains:

- structural contract fixtures under `testdata/contracts/v1/setup/`;
- blueprint fixtures under `testdata/setup/v1/blueprints/`;
- setup snapshot fixtures under `testdata/setup/v1/snapshots/`;
- setup report fixtures under `testdata/setup/v1/reports/`;
- exact planning rows in `testdata/setup/v1/decision-cases.json`;
- overall-state precedence rows in `testdata/setup/v1/aggregate-cases.json`.

Compact negative corpora derive exact documents from a named valid base using only RFC 6902 `add`, `remove` and `replace` patches. Every case starts from a fresh base. Schema-negative cases fail at schema validation; semantic-negative cases first pass schema validation and then fail only their declared semantic rule.

Every decision row names its mode, desired item, observed state, scopes, apply state, outcome, blocking status, verification class and diagnostic.

## Invariants

| Rule | Invariant |
| --- | --- |
| `SU-001` | A Project title never identifies an existing Project. |
| `SU-002` | Initialisation and steady-state bootstrap are different request kinds. |
| `SU-003` | Only verified assigned Project numbers enter a materialised contract. |
| `SU-004` | Bootstrap and adoption are create-once, not reconciliation. |
| `SU-005` | Existing compatible human/provider state is preserved. |
| `SU-006` | Extra live state is never deleted merely because it is undeclared. |
| `SU-007` | Every setup item has one stable key and one outcome. |
| `SU-008` | Missing, unknown, unavailable, forbidden, unsupported and ambiguous remain distinct. |
| `SU-009` | Mutation scopes are additive and surface-specific. |
| `SU-010` | Organisation-wide Issue Type/Field definition changes remain manual in v1. |
| `SU-011` | Project views are automatable create-only surfaces under current APIs. |
| `SU-012` | Built-in workflow configuration is manual and never inferred from name/enabled alone. |
| `SU-013` | Seed titles and bodies never act as seed identity. |
| `SU-014` | Ordinary edits to an existing seed are live shared state and are not restored. |
| `SU-015` | Adoption selector changes require a reviewed contract patch. |
| `SU-016` | Dry-run performs no provider mutation. |
| `SU-017` | Apply uses fresh preconditions and exact readback. |
| `SU-018` | Manual acknowledgement is not provider verification. |
| `SU-019` | Independent targets may continue only when no failed dependency is shared. |
| `SU-020` | Setup reports are deterministic and independent of provider enumeration order. |
| `SU-021` | Derived guidance is reviewed repository output, never competing contract authority or a direct setup push. |

## Superseded and deferred wording

This specification supersedes earlier wording that treated Project views as necessarily manual-only. GitHub now exposes supported view creation/configuration APIs, so v1 automates create-only view setup when runtime capability and scope allow it.

It retains the earlier prohibition on title fallback, continuous reconciliation and unverified success.

Ordinary task-update scopes and requests are now defined by [v1 ordinary task interactions](v1-task-interactions.md). Setup continues to own only create-once seed effects and never reconciles later task edits.

Explicitly deferred:

- normalisation, version compatibility and migration, owned by #34 and #31;
- guidance content/rendering, owned by #51;
- Go initialisation, bootstrap and adoption workflows, owned by #59, #23 and #24 respectively;
- deletion, cleanup and continuous drift repair;
- organisation-wide Issue Type and Issue Field administration.
