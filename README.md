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

If one repository manages several Projects, have each Project number available for the later routing step.

## 2. Install the skill in the repository

Open a terminal in the repository and run:

```text
gh skill install MiguelRodo/projects github-project-admin --agent universal --scope project
bash .agents/skills/github-project-admin/scripts/init-project.sh
```

Before the questions begin, the initializer explains what it will configure. It discovers the repository, its privacy, the Project title and whether its owner is a GitHub user or organisation. It asks only:

- whether other people work with you on the Project;
- whether the repository uses one Project or several;
- for a single Project, its owner and number.

This covers personal or collaborative repositories with one Project or several. Repository and Project privacy are discovered separately from GitHub.

For one Project, it creates `.projects/project.md` and adds a small starting section to `AGENTS.md`. It does not print the whole contract or ask you to edit either file.

For several Projects, it installs the skill and adds the `AGENTS.md` starting point, then gives you a request for an agent to finish the routing. Project routing is a real repository decision and is not guessed.

The initializer never changes live issues, fields or Project options. In particular, it leaves Priority exactly as it is. The initial contract marks its location and meaning as pending until an agent inspects the existing field and records a complete mapping. It will not add P3 or any other option during onboarding.

Finally, it asks whether it may stage, commit and push only these onboarding files. If committing or pushing fails, the local work is kept and the initializer prints the command to continue. Saying no simply leaves the files for you to review.

## 3. Set up ChatGPT

Create or open a [ChatGPT Project](https://chatgpt.com/projects), make the GitHub repository available to it, and paste this into the Project instructions:

> For work concerning a GitHub repository, especially reading or updating GitHub issues or Projects, first retrieve and follow the target repository's `AGENTS.md`. Follow the skill and configuration files it references. If the repository or `AGENTS.md` is unavailable, say so rather than guessing.
>
> Treat my prompt as the desired outcome. If this chat cannot make a required GitHub change, return the smallest safe command block for me to paste into a terminal, including a check of the result.

You can then ask for the result you want. ChatGPT will do what its GitHub connection allows and give you terminal commands for anything it cannot do directly.

## 4. Set up Codex cloud

Open [Codex environments](https://chatgpt.com/codex/settings/environments), create an environment and choose the repository.

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

After ChatGPT or Codex is ready, the initializer gives you a request tailored to the repository.

For a single Project, it can ask the agent to inspect existing issues, confirm the pending Priority location and mapping without changing the live field, propose useful Issue Type or Class and Workstream values, add sensible colours, organise the issues and build useful native parent/sub-issue relationships.

You choose whether the request ends with:

> Give me an overview of what you plan to do based on this request. Do not make changes until I approve the plan.

or:

> Do this now, then independently verify and summarise the changes.

For several Projects, the first request completes the dispatcher and routing contracts. Once that is committed, ask the agent to organise the issues across the resolved Projects.

To add another Project later, use a request such as:

> Start from `AGENTS.md`. Add Project 12 owned by `example-org` to this repository's existing multi-Project configuration. Preserve current routes, ask me only for the new routing decision, and validate the result.

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
