# v1 ordinary task-interaction conformance corpus

This directory is normative test data for [v1 ordinary task interactions](../../../docs/spec/v1-task-interactions.md).

`valid/` contains one schema-valid document for every request kind. `cases.json` derives schema-negative and semantic-invalid requests from those bases using RFC 6902 `add`, `remove` and `replace` operations. Each case starts from a fresh copy of its base.

A schema case MUST fail Draft 2020-12 validation at the declared instance path. A semantic case MUST pass schema validation first and then fail only its declared diagnostic against the synthetic context and `testdata/contracts/v1/interactions/valid/ordinary-full.yaml` where named.

The full repository contract and its mutation-vocabulary negative patches live under `testdata/contracts/v1/interactions/`.

`decision-cases.json` fixes exact read, plan, apply, no-op, stale, verification, create-recovery and private-companion outcomes. `operation-order-cases.json` fixes canonical order independently of request-array order. All identities and content are synthetic, and no case calls a live provider.
