# Project bootstrap

`MiguelRodo/projects` is the **cross-system bootstrapper for a project**. Use it when a research or lecturing project is new, only partly set up, or needs an audit against the standard project shape.

It is not the day-to-day task manager and it is not the implementation home for GitHub Project semantics or the local agent launcher.

## What this repository does

The `project-bootstrap` skill coordinates four surfaces without making any of them duplicate another:

| Surface | Role |
| --- | --- |
| Google Drive | Canonical project description and LLM-oriented source material |
| GitHub | Repository, GitHub Project and authoritative work items |
| Task-system registry | Durable mapping between the project, Drive and GitHub identities |
| ChatGPT Project | Conversational workspace, handed off through a verified manual checklist in version 1 |

A bootstrap can **create**, **resume** or **audit** a project. It first inspects existing state, then changes only what is missing or explicitly being corrected.

Typical requests are:

```text
Use $project-bootstrap to set up a new research project called Foobar.
```

```text
Use $project-bootstrap to finish setting up Foobar and verify the result.
```

```text
Use $project-bootstrap to audit Foobar and report anything missing from the standard setup.
```

## What it does not do

Do not use this repository for ordinary task administration inside an already bootstrapped project.

The surrounding repositories have deliberately separate responsibilities:

- [`MiguelRodo/github-projects-skill`](https://github.com/MiguelRodo/github-projects-skill) owns reusable GitHub issue and Project administration, the shared Project model and the deterministic `projects` CLI.
- [`MiguelRodo/pj`](https://github.com/MiguelRodo/pj) owns the local agent launcher, backend selection and queue execution.
- `MiguelRodo/projects` owns the **cross-system bootstrap** that joins Drive, GitHub, the task-system registry and the ChatGPT Project handoff.

That separation is intentional. This repository should not grow a second copy of GitHub Project semantics or `pj` behaviour.

## What a completed setup contains

For a standard research or lecturing project, the bootstrap aims to leave:

- a canonical Drive folder at `work/projects/<year>/<research|lecturing>/<project-slug>/`;
- a native Google Doc named `README` plus `references/llm/`;
- the requested GitHub repository and GitHub Project, when applicable;
- repository onboarding through `github-projects` when both a repository and Project are used;
- a bounded project-resources block in the repository README;
- exactly one verified row in the task system's Projects registry; and
- the exact manual checklist for creating or checking the private ChatGPT Project.

Every claimed write is independently read back. A partial setup remains resumable rather than being recreated from scratch.

## How the bootstrap runs

The workflow is intentionally conservative:

1. Inspect Drive, GitHub and the Projects registry for exact existing identities.
2. Resolve or create the requested GitHub repository and Project first, so their verified links are available.
3. Create only missing Drive resources and the minimal project README/reference structure.
4. Connect the repository to the canonical `github-projects` workflow when applicable.
5. Add or update exactly one registry row and record the material registry change.
6. Give the remaining ChatGPT Project setup as a manual, verifiable checklist.
7. Re-read each affected surface before reporting completion.

The full contract, including stop conditions and postconditions, is in [`PROJECT_BOOTSTRAP.md`](PROJECT_BOOTSTRAP.md). Agent-facing execution instructions live in [`skills/project-bootstrap/SKILL.md`](skills/project-bootstrap/SKILL.md).

## Install and use the skill

With a current GitHub CLI:

```bash
gh skill install MiguelRodo/projects project-bootstrap \
  --agent universal --scope user
```

Then invoke `$project-bootstrap` from a capable agent session. If skill installation is unavailable, point the agent at this repository and ask it to follow [`PROJECT_BOOTSTRAP.md`](PROJECT_BOOTSTRAP.md).

The shorter public guide is at [miguelrodo.github.io/projects](https://miguelrodo.github.io/projects/).

## Deterministic GitHub helper

The bootstrap skill includes `skills/project-bootstrap/scripts/github-resources.sh` for the GitHub resource-resolution step. It can resolve or create repositories and Projects, defaults to plan-only for missing resources, requires `--apply` before creation, and independently reads back the resulting GitHub identities.

Repeated GitHub Project administration belongs to the `projects` CLI in [`MiguelRodo/github-projects-skill`](https://github.com/MiguelRodo/github-projects-skill), not to a second CLI in this repository.

## Repository layout

| Path | Purpose |
| --- | --- |
| `PROJECT_BOOTSTRAP.md` | Human-readable bootstrap contract and postconditions |
| `skills/project-bootstrap/` | Installable Agent Skill, references, helper and offline tests |
| `templates/` | Minimal Drive README and repository-resource templates |
| `site/` | Public guide deployed through GitHub Pages |
| `.projects/` | This repository's own GitHub Project routing contract |

## Development

Repository changes go through a branch and pull request. The main offline checks are:

```bash
bash skills/project-bootstrap/tests/run.sh
bash tests/site.sh
```

The test suite must not mutate live GitHub, Drive or ChatGPT resources.

This repository is public. Do not add credentials, private Drive or spreadsheet identifiers, contact details, private source material or project-specific private URLs.
