# projects Project configuration

| Key | Value |
| --- | --- |
| Contract version | 1 |
| Mode | single |
| Issue repository | MiguelRodo/projects |
| Project owner | MiguelRodo |
| Project number | 40 |
| Project title | projects |
| Routing | Project 40 membership; no routing label |
| Privacy | public repository with a private user Project; public issue content only |
| Issue write-up style | unrestricted |

## Field locations

| Common dimension | Provider location | Provider field |
| --- | --- | --- |
| Class | project field | Class |
| Priority | project field | Priority |
| Status | project field | Status |
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

| Option | Colour |
| --- | --- |
| Task | YELLOW |
| Bug | RED |
| Enhancement | GREEN |
| Data | PINK |
| Analysis | PURPLE |
| Deliverable | ORANGE |
| Documentation | GRAY |
| Epic | BLUE |

## Status mapping

| Common value | Provider value |
| --- | --- |
| Todo | Todo |
| In progress | In progress |
| Done | Done |

## Bootstrap target

- Class is a required single-select Project field with the exact values and colours above.
- Priority requires the P0, P1, P2 and P3 options. Add P3 if live discovery confirms that only P0 to P2 exist.
- Target date is optional and remains blank unless an issue has a real deadline.
- Use native parent and sub-issue relationships for roadmap hierarchy. Parenthood does not imply Epic.
- In an authorised live migration, map both retired `Raw data` and `Processed data` values to `Data` before removing the old options.
- Workstream is not part of the active contract. A legacy field with that name may remain live until an explicitly authorised migration removes it.
- Bringing the Project into line with this contract may add the declared missing Class or Priority options and populate them, but must not delete or rename unrelated Project fields unless that migration is separately authorised.

## Governance

- The repository is public and the Project is private. Keep private material out of issues, pull requests, commits, logs and public reports.
- Project membership is the routing mechanism. Do not add a label merely to duplicate membership.
- Labels must not duplicate Class, Priority or Status.
- Assignment is explicit only.
- Read issue #1 for the roadmap and issue #67 for the current architecture before inventing scope or changing roadmap meaning.
- Exact requested administration and organising existing issues to this declared shape require no external source.
