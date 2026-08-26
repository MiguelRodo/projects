# v1 dispatcher and routing specification

Status: **normative pre-v1 specification**

Issue: #14

This document extends the [v1 single-Project repository contract](v1-single-project-contract.md) with a dispatcher topology and fixes the pure resolution semantics used to select one GitHub Project. It refines the [v1 conceptual model](v1-conceptual-model.md) without defining Go types, provider discovery, private operator profiles or mutation planning.

The keywords **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** are normative.

## Frozen decisions

| Question | v1 decision |
| --- | --- |
| Topology discriminator | `kind: DispatcherRepository` |
| Issue stores | Exactly one semantic GitHub issue store, repeated explicitly in every reused base target |
| Project targets | Two to 64, each with a unique target reference and Project locator |
| Field mappings | Per target, because different Projects may expose different fields and options |
| Automatic predicate | A non-empty conjunction of exact issue label names, represented by `labelsAll` |
| Rule order | No semantic meaning and no precedence |
| Fallback | Required and explicit, either one target reference or a typed no-match error |
| Multi-match | Always `routing.ambiguous`, including when all matching rules name the same target |
| Resolution cardinality | Exactly one target or one typed failure per call |
| Explicit selection | A separate request mode that selects an exact target reference and does not evaluate rules |
| Label comparison | Exact, case-sensitive Unicode scalar sequence equality with no trimming or normalisation |
| GitHub locator comparison | Owner kind exact; owner login and repository name ASCII case-insensitive; Project number exact |
| Implicit behaviour | No wildcard, regular expression, negation, weight, priority, title lookup, default target or fan-out |

The canonical contract path, encoding, YAML input profile, API version, reference syntax, field model, source hooks and mutation vocabulary remain those frozen by the single-Project specification.

## Dispatcher wire shape

```yaml
apiVersion: projectctl/v1
kind: DispatcherRepository
spec:
  ref: shared-work
  targets:
    - target:
        ref: backend
        issueStore:
          owner:
            kind: organization
            login: example-org
          repository: shared-tasks
        project:
          owner:
            kind: organization
            login: example-org
          number: 11
      fields:
        work-state:
          name: Backend state
          dataType: singleSelect
          values:
            queued:
              name: Queued
              color: gray
              description: Ready for backend scheduling.
    - target:
        ref: frontend
        issueStore:
          owner:
            kind: organization
            login: example-org
          repository: shared-tasks
        project:
          owner:
            kind: organization
            login: example-org
          number: 12
      fields: {}
    - target:
        ref: triage
        issueStore:
          owner:
            kind: organization
            login: example-org
          repository: shared-tasks
        project:
          owner:
            kind: organization
            login: example-org
          number: 13
      fields: {}
  routing:
    rules:
      - ref: backend-route
        labelsAll:
          - area:backend
        targetRef: backend
      - ref: frontend-route
        labelsAll:
          - area:frontend
        targetRef: frontend
    fallback:
      kind: target
      targetRef: triage
  sourceRefs: []
  requirements:
    features:
      - issues
      - project-custom-fields
      - projects-v2
    mutations: []
```

Every shown container is required. `targets` has at least two entries. `routing.rules`, `sourceRefs`, `requirements.features` and `requirements.mutations` are explicit even when a permitted collection is empty. Null remains invalid everywhere.

## Topology and authority

The repository containing `.projectctl/project.yaml` is the participating repository. It owns the shared declaration but is not inferred to be an issue store, Project owner or Project member.

Each dispatcher entry wraps the exact base `target` representation from the single-Project contract and adds that target's `fields` map. Reusing the complete target deliberately repeats the issue-store locator. Semantic validation requires every repeated locator to denote the same GitHub repository.

These resources remain distinct:

| Resource | Meaning |
| --- | --- |
| Participating repository | Stores and versions this contract |
| Issue store | Stores the shared GitHub issue and its labels |
| Target Project | May contain that issue as a Project item and owns the target's field values |

A dispatcher does not transfer issue ownership to a Project. Project item membership is observed and mutated separately from issue storage. A target Project MAY have a different owner from the issue store and from every other target Project.

Target declaration order has no semantic meaning. A canonical writer sorts target entries by exact target reference.

## Target and mapping rules

Every target entry contains:

- `target.ref`, unique within the dispatcher target namespace;
- one explicit `target.issueStore` locator;
- one explicit `target.project` locator;
- `fields`, which MAY be empty and applies only to that Project.

The same semantic Project locator MUST NOT appear under two target references. This avoids aliases with different mappings or route names. Project equality uses exact owner kind, ASCII case-insensitive owner login and exact Project number.

All issue-store locators MUST be equal under exact owner kind plus ASCII case-insensitive owner login and repository-name comparison. Raw spelling is still passed to discovery as supplied. The equality rule only recognises provider-equivalent locators and does not rewrite the contract.

Field and option selector rules from the single-Project contract apply independently inside each target. A duplicate selector is invalid within one target but the same field name and type MAY occur in different Projects.

## Routing language

### Rules

`routing.rules` is a required array of zero to 64 rules. Each rule contains exactly:

- `ref`, a unique contract reference in the rule namespace;
- `labelsAll`, one to 64 unique label names;
- `targetRef`, an exact reference to a declared target.

`labelsAll` is a positive conjunction. Given the known issue-label set `L`, a rule with required set `R` matches exactly when `R` is a subset of `L`. Extra issue labels do not prevent a match.

The sequence order of rules and the sequence order inside `labelsAll` have no semantic meaning. A canonical writer sorts rules by `ref` and label names by exact Unicode scalar sequence order.

V1 defines no other predicate key. In particular, `label`, `labelsAny`, `labelsNone`, title text, repository name, assignee, milestone, Project membership, field value, wildcard, regular expression, prefix, suffix, negation, score, priority and weight are invalid. Characters such as `*`, `?`, `[` and `]` inside a label name are ordinary literal characters.

Two rules with identical `labelsAll` sets are a static semantic conflict, regardless of their references or targets. Other rule pairs are permitted even though an issue may carry the union of both predicates. Runtime overlap is handled as ambiguity rather than declaration-order precedence.

### Fallback

`routing.fallback` is required and has exactly one of these shapes:

```yaml
fallback:
  kind: target
  targetRef: triage
```

```yaml
fallback:
  kind: error
```

The target form must reference a declared target. The error form returns `routing.no-match` when a known label set matches no rule. No missing, first or last target becomes an implicit fallback.

An empty `rules` array is valid. With a target fallback, automatic selection always uses that target after issue-store validation. With an error fallback, automatic selection always returns `routing.no-match`. This permits an explicit-only dispatcher without inventing a predicate.

## Comparison rules

| Value | Equality and ordering |
| --- | --- |
| Contract, target and rule references | Exact decoded ASCII bytes; lowercase syntax is already required |
| Label names | Exact Unicode scalar sequences; case-sensitive; no trimming, case folding or Unicode normalisation |
| Label sets | Set equality or subset membership using exact label equality; declaration order ignored |
| Owner kind | Exact enum equality |
| GitHub owner login | ASCII case-insensitive equality, with no non-ASCII mapping |
| GitHub repository name | ASCII case-insensitive equality, with no non-ASCII mapping |
| Project number | Exact integer equality |
| Field and option selectors | The exact rules frozen by the single-Project specification |

An implementation MUST NOT use locale-sensitive lowercasing. ASCII case-insensitive equality maps only `A` through `Z` to `a` through `z` for comparison. It preserves supplied spelling for provider discovery and diagnostics.

## Static validation

Structural validation occurs through `schemas/v1/repository-contract.schema.json`. Cross-entry checks occur during pure semantic validation after parsing and schema validation.

Semantic checks are evaluated in canonical target or rule reference order so declaration order cannot change the first diagnostic. A validator MAY report more than one conflict, but it MUST include the applicable stable key and canonical JSON Pointer.

| Diagnostic | Condition and canonical path |
| --- | --- |
| `contract.semantic.duplicate-target-ref` | A later canonical target has an exact target reference already declared; point to its `/target/ref` |
| `contract.semantic.duplicate-project-selector` | A later canonical target denotes an already declared Project; point to its `/target/project` |
| `contract.semantic.mixed-issue-stores` | A target does not denote the dispatcher's canonical first issue store; point to its `/target/issueStore` |
| `contract.semantic.duplicate-route-ref` | A later canonical rule has an exact rule reference already declared; point to its `/ref` |
| `contract.semantic.duplicate-route-predicate` | A later canonical rule has the same exact `labelsAll` set as an earlier rule; point to its `/labelsAll` |
| `contract.semantic.unknown-route-target` | A rule's `targetRef` does not resolve exactly; point to its `/targetRef` |
| `contract.semantic.unknown-fallback-target` | A target fallback does not resolve exactly; point to `/spec/routing/fallback/targetRef` |
| `contract.semantic.duplicate-field-selector` | Two mappings in one target use the same field selector; point to the later canonical field entry |
| `contract.semantic.duplicate-value-selector` | Two mappings in one target field use the same option selector; point to the later canonical value entry |

The canonical first issue store is the issue store of the target with the lexicographically smallest exact target reference. It is used only to make diagnostic selection stable. It has no runtime precedence.

## Resolver input

Resolution is a pure operation over one validated canonical contract and explicitly supplied request facts. It performs no API lookup and never guesses a missing fact.

The request selects one mode:

```yaml
selection:
  kind: automatic
```

or:

```yaml
selection:
  kind: explicit
  targetRef: backend
```

The issue-store fact and label fact each carry an explicit state. A known issue store carries its locator. A known label fact carries a set of exact names.

```yaml
issueStore:
  state: known
  value:
    owner:
      kind: organization
      login: example-org
    repository: shared-tasks
labels:
  state: known
  values:
    - area:backend
```

For either fact, `state` is exactly one of:

| State | Meaning |
| --- | --- |
| `known` | The complete value required by routing is supplied |
| `unknown` | No authoritative observation was supplied |
| `unavailable` | An attempted observation failed for a reason other than permission |
| `forbidden` | The acting principal was denied access to the observation |

The canonical model always contains a state, so omission is not another state. An adapter that lacks a fact supplies `unknown`; it MUST NOT substitute an empty label set. Known label values form a set and therefore contain no exact duplicates. Known empty labels and unknown labels have different meanings.

## Resolution algorithm

The resolver executes these steps in order.

### Common issue-store gate

1. If the issue-store state is `unknown`, fail with `routing.issue-store.unknown`.
2. If it is `unavailable`, fail with `routing.issue-store.unavailable`.
3. If it is `forbidden`, fail with `routing.issue-store.forbidden`.
4. If the known locator does not equal the contract's shared issue store under the frozen locator comparison, fail with `routing.issue-store.mismatch`.

No target, rule or fallback is considered before this gate succeeds.

### Explicit mode

1. Resolve `selection.targetRef` by exact target-reference equality.
2. If it does not resolve, fail with `routing.target.unknown`.
3. Return that one target with source `explicit`.

Explicit mode does not evaluate rules, inspect label facts or use fallback. It can therefore succeed when label state is unknown, unavailable or forbidden. Explicit selection is not a tie-breaker for automatic routing and cannot be supplied simultaneously with automatic mode.

### Automatic mode

1. If `routing.rules` is empty, skip label-state checks and continue with an empty match set.
2. Otherwise, map label states `unknown`, `unavailable` and `forbidden` to `routing.labels.unknown`, `routing.labels.unavailable` and `routing.labels.forbidden` respectively.
3. For a known label set, evaluate every rule in exact rule-reference order and retain all matches.
4. If more than one rule matches, fail with `routing.ambiguous` and report every matching rule reference and target reference in rule-reference order.
5. If exactly one rule matches, return its one target with source `rule`.
6. If no rule matches and fallback kind is `target`, return that target with source `fallback`.
7. If no rule matches and fallback kind is `error`, fail with `routing.no-match`.

The resolver never stops at the first match. Rule declaration order, target declaration order and fallback target placement cannot change the outcome.

## Result and trace

Every result contains `status`, the common issue-store gate outcome and a deterministic trace. A resolved result also contains exactly one `targetRef` and one source from `singleton`, `explicit`, `rule` or `fallback`. A failed result contains exactly one stable diagnostic.

Automatic traces contain:

- `evaluatedRuleRefs`, every evaluated rule reference in exact sorted order;
- `matchedRuleRefs`, every matching rule reference in the same order;
- `matchedTargetRefs`, present on ambiguity and aligned one-for-one with `matchedRuleRefs`, including repeated target references;
- `fallbackUsed`, true only for a resolved target fallback.

Explicit traces have empty evaluated and matched rule lists and `fallbackUsed: false`. Failure wording is not stable; the diagnostic key, ordered references and selected target are stable protocol data.

The routing decision corpus at `testdata/routing/v1/cases.json` is the executable table of inputs and exact expected outputs. #19 MUST implement results equivalent to every row without adding provider discovery or tie-breaking.

## Single-Project compatibility

The same resolver accepts a validated `SingleProjectRepository` as the degenerate one-target topology:

- the common issue-store gate is unchanged;
- automatic mode returns the sole target with source `singleton` and does not require label facts;
- explicit mode succeeds only for the sole exact target reference;
- rule and fallback traces are empty.

This behaviour does not alter the single-Project wire schema or add routing keys to it.

## Exactly one operation target

One resolver call returns exactly one target or one failure. It never returns a target set and never adds an issue to several Projects as a side effect.

An issue may already be a member of zero, one or several declared Projects. That observed membership does not choose a route, break a tie or invalidate the resolver result. A caller that deliberately wants operations against several Projects must make separate explicit resolution and planning requests. Automatic fan-out is outside v1.

## Requirements and mutation scope

`sourceRefs` and `requirements` remain contract-wide. The base feature and mutation meanings are unchanged, with these dispatcher applications:

- `projects-v2` applies to every declared target Project;
- any non-empty target `fields` map requires `project-custom-fields`;
- field creation, option creation or field writing requires at least one non-empty target field map;
- `project.item.add` requires `project-item-membership`;
- `project.item.field.write` requires both `project-custom-fields` and `project-item-membership`;
- mutation declarations authorise planning only against the one selected target and, for field writes, only that target's mappings.

Selection does not prove capability, permission, membership or mutation authority. Those remain later discovery, planning and execution gates.

## Invariant register

| ID | Rule |
| --- | --- |
| `RT-001` | `DispatcherRepository` is the sole v1 multi-Project discriminator. |
| `RT-002` | A dispatcher declares two to 64 targets and exactly one semantic issue store. |
| `RT-003` | Every entry reuses the complete base target and owns a separate field map. |
| `RT-004` | Target references and semantic Project locators are unique. |
| `RT-005` | Target and rule declaration order has no semantic meaning. |
| `RT-006` | A rule is only a non-empty `labelsAll` conjunction plus one target reference. |
| `RT-007` | Label equality is exact, case-sensitive and neither trimmed nor normalised. |
| `RT-008` | Locator equality uses exact owner kind, ASCII case-insensitive login and repository, and exact Project number. |
| `RT-009` | Duplicate predicate sets are invalid; other overlap is resolved only from runtime facts. |
| `RT-010` | More than one matching rule is always `routing.ambiguous`; order and same-target matches do not break ties. |
| `RT-011` | Fallback is mandatory and is exactly a target or a typed no-match error. |
| `RT-012` | Unknown, unavailable, forbidden, known empty and known non-empty fact states remain distinct. |
| `RT-013` | The issue-store gate precedes every selection mode. |
| `RT-014` | Explicit and automatic selection are distinct modes; explicit selection does not inspect labels. |
| `RT-015` | One resolver call returns exactly one target or one stable failure and never fans out. |
| `RT-016` | Existing Project membership neither selects a route nor resolves an ambiguity. |
| `RT-017` | Resolver traces and diagnostics are deterministic under input sequence permutation. |
| `RT-018` | Dispatcher mutation authority is global but applies only to the selected target and its mappings. |

## Explicitly deferred

This specification does not define:

- Go model or resolver APIs, owned by #19 and its prerequisites;
- YAML loading, diagnostic rendering or normalised serialisation, owned by #34;
- GitHub label, repository or Project discovery;
- private profile routing or privacy-mode resolution;
- fan-out, weighted routes, priorities, negative predicates or non-label predicates;
- mutation planning, execution or verification.

Later work MUST consume these rules and corpus rather than choose different precedence, matching, fallback or cardinality behaviour.
