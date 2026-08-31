# GitHub Projects made easier

This repository keeps one shared set of instructions for agents that work with GitHub issues and Projects.

The aim is simple: you should be able to ask for the result you want without explaining GitHub's commands every time.

For example:

- `Set stimgate#313 to P2.`
- `What are the highest-priority open items?`
- `Organise this Project and show me what changed.`

## Use it in a repository

The repository being managed keeps its own Project details in:

```text
.projects/project.md
```

Install the shared skill for Codex:

```bash
gh skill install MiguelRodo/projects github-project-admin --agent codex --scope user
```

Then ask for the outcome you want. The skill handles field lookup, safety checks and checking the result.

The short guide inside the skill explains [setup, updates and repository-specific setup](skills/github-project-admin/README.md).

## Set up Codex Cloud for this repository

Create an environment for `MiguelRodo/projects` and put this in its setup-script box:

```bash
bash skills/github-project-admin/scripts/setup.sh
```

Add `GH_TOKEN` as an environment variable with access to this repository and its GitHub Project. Do not add it as a setup-only secret because Codex needs it while doing the task. Allow agent internet access to `github.com` and `api.github.com`.

When the environment is ready, a normal request can be as short as:

> Bring the projects Project into line with `.projects/project.md`.

[OpenAI's Codex Cloud environment guide](https://developers.openai.com/codex/environments/cloud-environment) explains where setup scripts, environment variables and internet access are configured.

## What is here

- [`skills/github-project-admin/`](skills/github-project-admin/README.md): the reusable skill and its setup.
- [Issue #1](../../issues/1): the product roadmap.
- [Issue #67](../../issues/67): the design decisions from the first real trials.

A future command-line tool may be called `projects`, or `projectscli` if that name conflicts with another program. It will be added only when the trials show that it is useful.

## Contributing

Changes use GitHub issues and pull requests. See [CONTRIBUTING.md](CONTRIBUTING.md).

This project uses the [MIT Licence](LICENSE).
