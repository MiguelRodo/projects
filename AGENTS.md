# Agent guidance

This public repository contains Miguel's personal project-bootstrap contract and
its installable Agent Skill.

For skill changes, follow the Skill Creator guidance. Keep the main `SKILL.md`
focused and route conditional detail to `references/`. Put repeated,
deterministic GitHub logic in `scripts/` and test it without live mutations.

Do not add private Drive IDs, spreadsheet IDs, contact details, access tokens or
project-specific private URLs. Synthetic test data only.

Work on a branch and open a pull request rather than pushing directly to
`main`. Before proposing a change, run:

```bash
python3 /path/to/skill-creator/scripts/quick_validate.py \
  skills/project-bootstrap
bash skills/project-bootstrap/tests/run.sh
```

Tests must not mutate live GitHub, Drive or ChatGPT resources.

## Ponytail

For coding, refactoring, bug-fixing, review and implementation design, read
`.agents/skills/ponytail/SKILL.md` and apply Ponytail in **full** mode by
default. Do not use **ultra** unless the operator explicitly requests it.
Vendoring/provenance notes are in `.agents/skills/ponytail/README.md`.

Repository-specific instructions and the project-bootstrap contract take
precedence over generic Ponytail guidance.

<!-- github-projects:start -->
## GitHub issues and Projects

For GitHub issue or Project administration, use
`.agents/skills/github-projects/SKILL.md` and read
`.projects/project.md` before acting.
<!-- github-projects:end -->
