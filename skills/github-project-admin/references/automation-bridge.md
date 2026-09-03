# Issue-triggered Project administration bridge

Use this bridge when an authorised GitHub Project mutation cannot be performed by the current chat/provider surface. The bridge deliberately mirrors the local short-prompt operator pattern: create one small operation issue, label it for execution, and let an execution-capable agent work on that issue using the repository's normal `AGENTS.md`, `.projects` contract and `github-project-admin` rules.

## Execution model

The operation issue is an execution queue item, not a mirror of the real task. Keep it concise and include only the exact authorised outcome, target identity, and any stale-sensitive state actually inspected by the originating surface.

The workflow triggers only when `automation:project-admin` is applied. It ignores issue comments as instructions. The model is told simply to work on that operation issue, just as a local launcher can receive a short natural-language task request and rely on repository guidance for the procedure.

GitHub Copilot CLI runs non-interactively in the checked-out repository. It automatically loads the repository's agent instructions unless explicitly disabled. The bridge additionally tells it to:

- treat the operation issue title and body as the exact mutation authority;
- read the repository's `AGENTS.md`, `.projects` contract and referenced `github-project-admin` skill;
- use live GitHub state and `gh`/API operations;
- re-read stale-sensitive state before writes;
- preserve unrelated state;
- independently verify every requested delta;
- avoid repository-file changes; and
- report `PROJECT_ADMIN_RESULT: DONE` only after verified success, otherwise `PROJECT_ADMIN_RESULT: BLOCKED`.

The first implementation uses MAI-Code-1.1-Flash through GitHub Copilot CLI. The workflow grants only repository reads plus `gh` shell commands to the model. The Project-capable token is exposed to that bounded Copilot CLI run as `GH_TOKEN`, so this is intentionally the same trust model as an unattended local Project-administration agent rather than the earlier manifest/executor split.

## Operation issue format

The issue does not need a machine-readable schema. Plain natural language is enough when the target and requested outcome are unambiguous. For example:

```markdown
## Requested change

For `example/repo#7` in Project 40:

- ensure the issue belongs to the Project;
- set Priority to `P1`;
- leave all unrelated issue and Project state unchanged.

## Observed state

The current chat surface could not inspect Project membership or Priority, so read them live immediately before any write.
```

Apply `automation:project-admin` only after the issue is complete. The workflow removes the trigger label when it claims the operation. A blocked retry requires removing `automation:blocked` and re-applying `automation:project-admin`.

## Repository setup

The repository must have Actions enabled and a label named `automation:project-admin`:

```text
gh label create automation:project-admin --color 5319e7 --description "Run the bounded Project administration bridge"
```

The workflow uses the built-in `GITHUB_TOKEN` for Copilot requests, so there is no separate model-provider API secret. It still needs one Actions secret containing a GitHub credential that can perform the required Project operations:

```text
gh secret set PROJECTS_TOKEN
```

For a user-owned Project, a classic personal access token with `project` and the repository access needed for the target issues is the least surprising `PROJECTS_TOKEN`. Add `read:org` when organisation discovery requires it. Prefer a shorter expiry and the minimum repository access that still covers the declared Project work. Never place the token in an issue, prompt, comment or log.

GitHub Copilot CLI authenticates separately with the workflow's built-in token through `COPILOT_GITHUB_TOKEN`. The Project credential is passed as `GH_TOKEN` only so the agent's `gh` commands can perform the authorised Project operation. The CLI is instructed to redact `GH_TOKEN` from output.

The reusable workflow automatically creates the reserved result labels `automation:running`, `automation:done` and `automation:blocked` when needed.

## Reusing the workflow from another repository

A repository can keep a tiny caller workflow and reuse the implementation from `MiguelRodo/projects`. Pin the reusable workflow to a reviewed commit SHA:

```yaml
name: Project administration bridge

on:
  issues:
    types: [labeled]

permissions:
  contents: read
  issues: write
  copilot-requests: write

jobs:
  bridge:
    if: github.event.label.name == 'automation:project-admin'
    uses: MiguelRodo/projects/.github/workflows/project-admin-bridge.yml@<reviewed-commit-sha>
    with:
      issue_number: ${{ github.event.issue.number }}
    secrets:
      PROJECTS_TOKEN: ${{ secrets.PROJECTS_TOKEN }}
```

The called workflow checks out the caller repository, so the agent sees that repository's own `AGENTS.md` and `.projects` contract.

## Chat/provider handoff

When the current surface cannot perform an authorised Project mutation:

1. inspect the target and stale-sensitive state as far as the current surface allows;
2. create one concise operation issue in the repository that owns the relevant Project contract;
3. apply `automation:project-admin` only when the workflow and `PROJECTS_TOKEN` are configured;
4. report the operation as queued, not completed;
5. treat the requested mutation as complete only after the operation issue closes with a verified `DONE` result.

If the bridge is not installed, the trigger label is unavailable, or the workflow is blocked, fall back to the smallest executable command/readback handoff. Creating the operation issue itself never proves that the requested Project change happened.
