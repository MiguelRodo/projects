# projects Project configuration

| Key | Value |
| --- | --- |
| Contract version | 1 |
| Mode | single |
| Issue repository | MiguelRodo/projects |
| Project owner | MiguelRodo |
| Owner type | user |
| Project number | 40 |
| Project title | projects |
| Routing | Project 40 membership; no routing label |
| Privacy | public repository with a private user Project; public issue content only |

## Field locations

| Common dimension | Provider location | Provider field |
| --- | --- | --- |
| Class | project field | Class |
| Priority | project field | Priority |
| Status | project field | Status |
| Workstream | project field | Workstream |
| Due date | project field | Target date |
| Parent | native issue relationship | Parent issue |

## Priority mapping

This Project uses the common names directly.

| Common value | Provider value |
| --- | --- |
| P0 | P0 |
| P1 | P1 |
| P2 | P2 |
| P3 | P3 |

## Class values

- Epic
- Task
- Bug
- Enhancement
- Research
- Documentation

## Status mapping

| Common value | Provider value |
| --- | --- |
| Todo | Todo |
| In progress | In progress |
| Done | Done |

## Workstream values

- Shared skill
- Repository contracts
- Agent distribution
- Provider integration
- CLI
- User interface
- Testing and pilots
- Governance and documentation

## Bootstrap target

- Class and Workstream are required single-select Project fields with the exact values above.
- Priority requires the P0, P1, P2 and P3 options. Add P3 if live discovery confirms that only P0 to P2 exist.
- Target date is optional and remains blank unless an issue has a real deadline.
- Use native parent and sub-issue relationships for the roadmap hierarchy.
- Bringing the Project into line with this contract may add the declared missing fields or options and populate them, but must not delete or rename unrelated Project fields.

## Governance

- The repository is public and the Project is private. Keep private material out of issues, pull requests, commits, logs and public reports.
- Project membership is the routing mechanism. Do not add a label merely to duplicate membership.
- Labels must not duplicate Class, Priority, Status or Workstream.
- Assignment is explicit only.
- Read issue #1 for the roadmap and issue #67 for the current architecture before inventing scope or changing roadmap meaning.
- Exact requested administration and organising existing issues to this declared shape require no external source.
