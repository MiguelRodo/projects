# `projects`

`projects` is the canonical home of an agent-neutral skill for safe GitHub issue and Project administration, plus the roadmap for a future setup and command-reliability CLI.

## Status

The first supported surface is the [`github-project-admin` Agent Skill](skills/github-project-admin/SKILL.md). It lets capable agents act on short outcome requests while applying one shared procedure for exact resolution, fresh inspection, narrow mutations, preservation, stale-state refusal and verified readback.

There is no supported CLI or Go library yet. The authoritative roadmap and execution order are in [issue #1](../../issues/1), with the post-dogfood architecture decision in [issue #67](../../issues/67).

## Install the skill

With GitHub CLI 2.90.0 or newer:

```bash
gh skill install MiguelRodo/projects github-project-admin --agent codex --scope user
```

To install every skill published by this repository:

```bash
gh skill install MiguelRodo/projects --all --agent codex --scope user
```

The skill is agent-neutral. `--agent codex` is one installation adapter, not part of its semantic contract.

Run its generic environment preflight from an installed or checked-out copy:

```bash
bash skills/github-project-admin/scripts/setup.sh --repository OWNER/REPO
```

The host supplies credentials and network access. The setup installs or verifies `gh`, checks authentication and optional runtime context, and does not require a repository's language toolchain.

## Repository contract

Each participating repository keeps only its local topology and mappings under `.projects/`:

```text
.projects/project.md
```

A multi-Project issue store may also use `.projects/projects/*.md`. Repository `AGENTS.md` files and provider Project instructions should only route Project-administration work to the shared skill and the local `.projects/` contract.

The default cross-provider Priority vocabulary is lossless:

| Common value | Default provider value |
| --- | --- |
| P0 | Urgent |
| P1 | High |
| P2 | Medium |
| P3 | Low |

Repositories may declare a complete one-to-one override.

## Execution surfaces

The skill currently uses proven native `gh`, versioned REST and GraphQL operations. If a conversational surface cannot mutate GitHub, it returns the smallest executable command block governed by the same procedure.

A future CLI may be named `projects`, or `projectscli` if a collision review requires it. Its initial jobs are setup assistance and reliable execution or generation of Project-administration commands. It will be an optional backend for the skill, not a prerequisite for capable agents.

## Earlier design material

The documents under `docs/spec/v1/`, their schemas and fixtures record the more elaborate pre-dogfood protocol design. They remain useful evidence but are not the active ordinary-administration contract. Issue #67 governs which parts are retained, narrowed or retired.

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Work is issue-driven and uses one issue per PR.

Security and privacy reports follow [SECURITY.md](SECURITY.md).

## Licence

This project is licensed under the [MIT Licence](LICENSE).
