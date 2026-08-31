# Contributing to `projects`

Thank you for contributing. The active product surface is the agent-neutral `github-project-admin` skill. A future CLI is deliberately deferred until repeated setup or command-generation friction justifies it.

## Start with the authoritative issue

Before working:

1. Read [issue #1](../../issues/1) for the execution order.
2. Read [issue #67](../../issues/67) for the post-dogfood architecture.
3. Select one open leaf issue whose dependencies are complete.
4. Treat the exact files and interfaces named by that issue as the change boundary.

If a required decision is absent or contradictory, stop and describe the conflict in the PR.

## Pull-request discipline

- Work on one issue per pull request.
- Create a branch and open a PR against `main`; do not push directly to `main`.
- Change only the files permitted by the issue.
- Do not include unrelated cleanup or anticipatory abstractions.
- Keep the PR description aligned with the actual diff.
- Use `Refs #N` until every completion condition is met; use `Closes #N` only when the PR genuinely completes the issue.

## Public boundary

Skills, code, schemas, documentation and fixtures must remain agent-neutral and operator-neutral.

Do not include:

- private repository names;
- private Drive, document or spreadsheet identifiers;
- credentials, tokens or secret material;
- personal names or account-specific configuration;
- private source locations or private destination URLs.

Use clearly synthetic fixtures and local test doubles.

## Skill changes

Keep the main `SKILL.md` concise and route detailed provider commands or conditional guidance to `references/`. Put deterministic repeated logic in `scripts/` and test it offline.

Validate the canonical skill with:

```bash
python /path/to/skill-creator/scripts/quick_validate.py \
  skills/github-project-admin
bash skills/github-project-admin/tests/run.sh
```

The first command uses the validator distributed with the Agent Skill authoring tools. CI runs the repository-owned offline checks.

Tests must not mutate live GitHub issues, Projects, repositories or private systems. Use synthetic fixtures and fake commands.

## Future CLI boundary

There is no CLI implementation yet. When a roadmap issue authorises one, prefer the executable name `projects`; use `projectscli` only if a collision review rejects it. The CLI must remain an optional backend for the skill and should first address setup and reliable command execution or generation.

Do not create a second executable or revive a mandatory-client architecture without revising the architecture issue first.
