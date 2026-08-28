# v1 setup conformance corpus

This directory is normative test data for [v1 setup outcomes](../../../docs/spec/v1-setup-outcomes.md).

`blueprints/`, `snapshots/` and `reports/` contain complete valid documents plus schema-negative patch cases. Repository-contract setup cases live under `../../contracts/v1/setup/`.

Patch corpora derive one exact document from their named `base` file. Patches use the RFC 6902 `add`, `remove` and `replace` operations and JSON Pointer paths. `add` with a final `-` appends to an array. Each case starts from a fresh copy of the base. A schema case MUST fail Draft 2020-12 validation at the declared diagnostic. A semantic case MUST pass schema validation first and then fail only the declared semantic rule.

`decision-cases.json` fixes per-item planning outcomes. `aggregate-cases.json` fixes overall report-state precedence. Provider identifiers and timestamps are synthetic.
