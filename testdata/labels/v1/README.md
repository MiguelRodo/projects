# V1 label decision corpora

The cases.json corpus fixes issue-store gating, exact label-name mapping, overall-Project routing, optional sub-project classification, general labels and preservation of undeclared live labels. It records the exact combined routing/classification outcome but deliberately omits routing traces, which are owned by testdata/routing/v1/cases.json.

The discovery-cases.json corpus fixes repository-label catalogue discovery, exact adoption, case conflicts, ambiguity and create authority. Provider IDs are synthetic observed values and never contract identity.

The transition-cases.json corpus fixes complete-before-state requirements, narrow mutation authority, idempotent no-ops, atomic route and sub-project transitions, after-state cardinality and preservation of undeclared labels. It is an abstract guard table, not the ordinary request wire format owned by #50.

Consumers MUST resolve each contractRef relative to the repository root. Arrays in expected results are already in canonical exact-reference or exact-name order. Implementations may use another internal representation but MUST produce an equivalent result for every row.
