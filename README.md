# GitHub Projects made easier

This repository provides one shared skill for ChatGPT, Codex and other agents that work with GitHub issues and Projects. Once a repository is set up, ordinary requests can stay short:

- `What should I work on next?`
- `Set example#313 to P2.`
- `Organise the existing issues and show me what changed.`

## Before you start

You need Git, Bash and [GitHub CLI](https://cli.github.com/) 2.96 or newer on the computer where the repository is checked out. Check GitHub CLI with:

```text
gh --version
```

The commands below are one line each. They work in Bash, Git Bash and WSL. They also work from PowerShell when `bash` is installed and available as a command.

Sign in to GitHub CLI if you have not already done so:

```text
gh auth login --web --scopes "project,read:org"
```

If you are already signed in but have not granted Project access, add the permissions instead:

```text
gh auth refresh --scopes "project,read:org"
```

Then check the active login:

```text
gh auth status
gh api user --jq .login
```

The initializer checks authentication and Project access again before asking any questions.

## 1. Create or find the GitHub Project

Open the **Projects** tab on your GitHub profile or organisation. Create the Project if it does not exist yet.

Keep the Project number handy. It is the number after `/projects/` in the web address:

```text
https://github.com/users/example/projects/40       Project number 40
https://github.com/orgs/example-org/projects/12    Project number 12
```

If one repository manages several Projects, have the first Project number ready. The initializer will add it and then ask whether you want to add another.

## 2. Install the skill in the repository

Open a terminal in the repository and run:

```text
gh skill install MiguelRodo/projects github-project-admin --agent universal --scope project
bash .agents/skills/github-project-admin/scripts/init-project.sh
```

Before the questions begin, the initializer explains what it will configure. It discovers the repository, its privacy, the Project title and whether its owner is a GitHub user or organisation. It asks only:

- whether other people work with you on the Project;
- whether the repository uses one Project or several;
- for each Project you add, its owner, number, Project key and routing label.

This covers personal or collaborative repositories with one Project or several. Repository and Project privacy are discovered separately from GitHub.

For one Project, it creates `.projects/project.md` and adds a small starting section to `AGENTS.md`. It does not print the whole contract or ask you to edit either file.

For several Projects, it creates a validated dispatcher first. It then offers to add one Project, discovers the live Project, asks for the routing identity, creates its child contract and asks whether to add another. If you stop before adding one, the empty dispatcher is saved safely but cannot yet resolve ordinary Project requests. Rerun the initializer to continue.

The initializer never changes live issues, fields or Project options. In particular, it leaves Priority exactly as it is. The initial contract marks its location and meaning as pending until an agent inspects the existing field and records a complete mapping. It will not add P3 or any other option during onboarding.

Finally, it asks whether it may stage, commit and push only these onboarding files. If committing or pushing fails, the local work is kept and the initializer prints the command to continue. Saying no simply leaves the files for you to review. Commit and push the installed skill and configuration before expecting a remote chat or agent to retrieve them.

## 3. Set up a chat interface

Create or open a [ChatGPT Project](https://chatgpt.com/projects), make the GitHub repository available to it, and paste this into the Project instructions:

> For work concerning a GitHub repository, especially reading or updating GitHub issues or Projects, first retrieve and follow the target repository's `AGENTS.md`. Follow the skill and configuration files it references. If the repository or `AGENTS.md` is unavailable, say so rather than guessing.
>
> Treat my prompt as the desired outcome. If this chat cannot make a required GitHub change, return the smallest safe command block for me to paste into a terminal, including a check of the result.

You can then ask for the result you want. The chat can inspect and propose the work. After you approve the proposal, it will do what its GitHub connection allows and give you terminal commands for anything it cannot do directly.

## 4. Set up an execution-capable agent

Codex cloud is one execution-capable option. Open [Codex environments](https://chatgpt.com/codex/settings/environments), create an environment and choose the repository.

Use this setup command:

```text
bash .agents/skills/github-project-admin/scripts/setup.sh
```

Create a [classic GitHub personal access token](https://github.com/settings/tokens/new), give it an expiry, and select the `repo`, `read:org` and `project` scopes. If the organisation uses SSO, authorise the token for that organisation.

In the Codex environment:

- add the token as an environment variable named `GH_TOKEN`, not a setup-only secret;
- enable internet access during the agent phase;
- allow `github.com` and `api.github.com`.

Environment variables remain available while the agent works, whereas setup-only secrets do not. The [official Codex environment guide](https://developers.openai.com/codex/environments/cloud-environment) explains these settings.

## 5. Give it a useful first request

After the chat interface or execution-capable agent is ready, the initializer gives you one proposal-only request tailored to a resolved Project.

The request asks the surface to inspect existing issues, confirm the pending Priority location and mapping without changing the live field, propose useful Issue Type or Class and Workstream values, preserve useful definitions and colours, organise the issues, build useful native parent/sub-issue relationships and suggest optional sub-project labels only where they add value. It explicitly forbids live changes until you approve the proposal.

After approval, an execution-capable agent can apply and verify the proposal. A chat interface that cannot write should return the smallest safe command block with independent readback. To add another Project later, rerun the initializer; it preserves every current route and contract.

## Update the installed skill

Run this in the repository and commit the changed files:

```text
gh skill update github-project-admin
```

Updating the shared skill does not replace `.projects/project.md` or a repository-specific `.projects/setup.sh`.

## Maintaining this repository

For `MiguelRodo/projects` itself, the canonical skill is under `skills/`. Its Codex setup command is:

```text
bash skills/github-project-admin/scripts/setup.sh --install-skill-from .
```

The shorter operator guide is in [the skill README](skills/github-project-admin/README.md). [Issue #1](../../issues/1) contains the roadmap and [issue #67](../../issues/67) records the current architecture.

Changes use GitHub issues and pull requests. See [CONTRIBUTING.md](CONTRIBUTING.md). This project uses the [MIT Licence](LICENSE).
