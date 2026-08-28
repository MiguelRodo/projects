# V1 dimension decision corpus

These files are normative exact-result tables for issue #48.

- `discovery-cases.json` fixes provider configuration discovery and adoption.
- `read-cases.json` fixes canonical reads, clear values and fact-state failures.
- `write-cases.json` fixes set, clear, no-op, scope and readback behaviour.

Contract fixtures live under `testdata/contracts/v1/dimensions/`. All provider identities, owners, repositories, fields, actors and values are synthetic.

Case order is documentary only. Arrays named `operations`, `actorIds` and `availableNames` use the exact order returned in `expected`; contract map order never supplies that order.
