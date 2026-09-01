# GitHub Project configuration

| Key | Value |
| --- | --- |
| Contract version | 1 |
| Mode | single |
| Issue repository | MiguelRodo/projects |
| Project owner | MiguelRodo |
| Project number | 40 |
| Project title | projects |
| Routing | Project membership; no routing label |
| Privacy | public repository with a private Project |

## Field locations

| Common dimension | Provider location | Provider field |
| --- | --- | --- |
| Class | project field | Class |
| Priority | pending live inspection | Priority |
| Status | project field | Status |
| Due date | project field | Target date |
| Parent | native issue relationship | Parent issue |

## Priority mapping

Priority mapping status: pending

The initializer left the Project's existing Priority field and options unchanged.
Before using Priority, an agent must confirm its provider location, inspect the
live options and replace the pending status with a complete one-to-one P0, P1,
P2 and P3 mapping.

## Class / Issue Type

The Class or Issue Type option set is intentionally not fixed by onboarding.
An agent may inspect the existing issues and suggest a useful vocabulary before
the live Project is changed. Workstream is not a standard semantic dimension.

## Governance

- This is a collaborative Project.
- Project membership determines Project scope; no routing label is required.
- Labels must not duplicate Class, Priority or Status.
- Assignment is explicit only unless a later repository decision says otherwise.
- Exact requested administration requires no separate scope-design source.
