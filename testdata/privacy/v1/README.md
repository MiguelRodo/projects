# v1 shared privacy decision corpus

`cases.json` is the normative machine-readable policy table. The profile objects are abstract resolver inputs, not the private wire format owned by #32. `privateSupplementWrittenToGitHub` is always false and is asserted in every row.

| Case | Shared policy | Profile | Private supplement | Expected |
| --- | --- | --- | --- | --- |
| `PV-C001` | absent baseline | absent | absent | shareable, repository default |
| `PV-C002` | shareable only | absent | absent | shareable, repository default |
| `PV-C003` | full only | absent | absent | full GitHub, repository default |
| `PV-C004` | both, companion | selects full | absent | full GitHub, profile choice |
| `PV-C005` | shareable only | selects full | absent | `privacy.mode.unsupported` |
| `PV-C006` | absent baseline | absent | present | `privacy.companion.unsupported` |
| `PV-C007` | both, companion | absent | present | `privacy.companion.profile-required` |
| `PV-C008` | both, companion | selects full | present | `privacy.companion.mode-conflict` |
| `PV-C009` | both, companion | selects shareable | present, known Issue ID | routed with canonical linkage |
| `PV-C010` | shareable only | selects shareable | present | `privacy.companion.unsupported` |
| `PV-C011` | both, companion | selects shareable | present, identity unknown | `privacy.linkage.issue-identity.unknown` |
| `PV-C012` | both, companion | selects shareable | present, identity unavailable | `privacy.linkage.issue-identity.unavailable` |
| `PV-C013` | both, companion | selects shareable | present, identity forbidden | `privacy.linkage.issue-identity.forbidden` |
| `PV-C014` | both, companion | selects shareable | present, identity not found | `privacy.linkage.issue-identity.not-found` |
| `PV-C015` | both, companion | selects shareable | present, PullRequest ID | `privacy.linkage.issue-identity.type-mismatch` |
| `PV-C016` | both, companion | absent | absent | shareable without private access |
| `PV-C017` | both, companion | operator A selects shareable | present | routed privately for A |
| `PV-C018` | same contract | operator B selects full | absent | full GitHub for B |
| `PV-C019` | full only | absent | present | `privacy.companion.mode-conflict` |
| `PV-C020` | both, default full | absent | absent | full GitHub, repository default |

Rows `PV-C017` and `PV-C018` intentionally use the same shared contract. They prove that two operators can resolve different supported behaviour without changing repository state or exposing either operator's destination.
