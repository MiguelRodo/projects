# v1 routing decision corpus

`cases.json` is the normative machine-readable table for pure target resolution. Each row names a validated contract fixture, supplies the complete selection and fact state, and fixes the exact result and deterministic trace.

The compact index below is for review. `match-ci` means a known issue-store locator that matches only after the specified ASCII case-insensitive comparison. Label lists use exact, case-sensitive values.

| Case | Contract | Mode | Issue store | Labels | Expected |
| --- | --- | --- | --- | --- | --- |
| `RT-C001` | label routing | automatic | match | `area:backend` | target `backend`, rule |
| `RT-C002` | label routing | automatic | match | backend plus extras | target `backend`, rule |
| `RT-C003` | label routing | automatic | match | `area:frontend` | target `frontend`, rule |
| `RT-C004` | label routing | automatic | match | known empty | target `triage`, fallback |
| `RT-C005` | label routing | automatic | match | case-mismatched backend | target `triage`, fallback |
| `RT-C006` | label routing | automatic | match-ci | `area:backend` | target `backend`, rule |
| `RT-C007` | locator case | automatic | match-ci | literal `area:*` | target `alpha`, rule |
| `RT-C008` | locator case | automatic | match-ci | `area:backend` | target `beta`, fallback |
| `RT-C009` | error fallback | automatic | match | unrelated | `routing.no-match` |
| `RT-C010` | error fallback | automatic | match | incomplete conjunction | `routing.no-match` |
| `RT-C011` | error fallback | automatic | match | complete API conjunction | target `api`, rule |
| `RT-C012` | error fallback | automatic | match | API and web rules | `routing.ambiguous` |
| `RT-C013` | same-target overlap | automatic | match | two rules for `delivery` | `routing.ambiguous` |
| `RT-C014` | permuted routing | automatic | match | `area:backend` | target `backend`, rule |
| `RT-C015` | label routing | automatic | match | unknown | `routing.labels.unknown` |
| `RT-C016` | label routing | automatic | match | unavailable | `routing.labels.unavailable` |
| `RT-C017` | label routing | automatic | match | forbidden | `routing.labels.forbidden` |
| `RT-C018` | label routing | automatic | unknown | known backend | `routing.issue-store.unknown` |
| `RT-C019` | label routing | automatic | unavailable | known backend | `routing.issue-store.unavailable` |
| `RT-C020` | label routing | automatic | forbidden | known backend | `routing.issue-store.forbidden` |
| `RT-C021` | label routing | automatic | mismatch | known backend | `routing.issue-store.mismatch` |
| `RT-C022` | label routing | explicit `backend` | match | forbidden | target `backend`, explicit |
| `RT-C023` | label routing | explicit unknown | match | known empty | `routing.target.unknown` |
| `RT-C024` | label routing | explicit unknown | mismatch | unknown | `routing.issue-store.mismatch` |
| `RT-C025` | explicit only | automatic | match | unavailable | `routing.no-match` |
| `RT-C026` | single Project | automatic | match-ci | unavailable | target `primary`, singleton |
| `RT-C027` | single Project | explicit `primary` | match | forbidden | target `primary`, explicit |
| `RT-C028` | single Project | explicit unknown | match | known empty | `routing.target.unknown` |

Implementations MUST compare the complete `expected` object in `cases.json`, including ordered rule references, duplicated matched target references in same-target ambiguity, the fallback flag and the issue-store gate result.
