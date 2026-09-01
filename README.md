# GitHub Projects made easier

This repository provides one shared skill for agents that work with GitHub issues and Projects. Once a repository is set up, ordinary requests can stay short:

- `Set example#313 to P2.`
- `What should I work on next?`
- `Organise this Project and show me what changed.`

## Set up a repository

### 1. Create the GitHub Project

Create the Project in GitHub, or open the Project you already use. You can find Projects under the **Projects** tab on your GitHub profile or organisation.

Keep the Project number handy. It is the number after `/projects/` in the Project's web address:

```text
https://github.com/users/example/projects/40       Project number 40
https://github.com/orgs/example-org/projects/12    Project number 12
```

If this repository manages several Projects, have each Project number available.

### 2. Install the skill and answer the questions

Open a terminal in the repository and run:

```bash
gh skill install MiguelRodo/projects github-project-admin \
  --agent universal \
  --scope project

bash .agents/skills/github-project-admin/scripts/init-project.sh
```

The questions cover who works on the Project, whether the repository manages one or several Projects, routing and field mappings. GitHub supplies facts such as repository visibility, Project title and whether the owner is a person or organisation.

For one Project, the script creates `.projects/project.md` and adds a small section to `AGENTS.md`. For several Projects, it gives you a tailored request for ChatGPT or Codex to create the routing configuration safely. You do not need to write these files by hand.

Review and commit the installed skill, `.projects/` and the `AGENTS.md` change. The script does not change live issues or Project fields.

### 3. Choose how you want to use it

#### ChatGPT or another conversational service

Create a ChatGPT Project, or the equivalent workspace in another service, and give it access to the relevant GitHub repository. Put this in its Project instructions:

> For work concerning a GitHub repository, especially reading or updating GitHub issues or Projects, first retrieve and follow the target repository's `AGENTS.md`. Follow the skill and configuration files it references. If the repository or `AGENTS.md` is unavailable, say so rather than guessing.

You can then ask ordinary questions or request changes. If that chat cannot make a GitHub change itself, it should give you a short command block to paste into a terminal.

#### Codex Cloud or another shell-capable agent

Create an agent environment for the repository and use this setup command:

```bash
bash .agents/skills/github-project-admin/scripts/setup.sh
```

Provide `GH_TOKEN` as an environment variable with access to the repository and Project, and allow access to `github.com` and `api.github.com`. The agent can then inspect, execute and verify the GitHub commands itself.

## Which repository shapes are supported?

The initializer asks about two independent choices:

- **Personal or collaborative:** whether only you use the Project or other people work on it too.
- **One or several Projects:** whether one repository maps to one Project or routes issues among several Projects.

Repository and Project privacy are discovered from GitHub. Single-Project configuration is generated directly. Multi-Project routing is completed through the tailored agent request because the correct routing rule is a genuine user decision.

## Maintaining this repository

For `MiguelRodo/projects` itself, the canonical skill is already under `skills/`. Its Codex Cloud setup command is:

```bash
bash skills/github-project-admin/scripts/setup.sh --install-skill-from .
```

The detailed guide is in [the skill README](skills/github-project-admin/README.md). [Issue #1](../../issues/1) contains the roadmap and [issue #67](../../issues/67) records the current architecture.

Changes use GitHub issues and pull requests. See [CONTRIBUTING.md](CONTRIBUTING.md). This project uses the [MIT Licence](LICENSE).
