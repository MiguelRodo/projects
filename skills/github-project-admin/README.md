# GitHub Project admin

This skill lets an agent manage GitHub issues and Projects from ordinary requests.

You can ask things like:

- `Set example#313 to P2.`
- `What should I work on next?`
- `Add these issues to the Project and organise them.`

The skill checks the current Project before changing anything and checks the result afterwards. You do not need to put those instructions in every request.

## Set it up

The repository needs a `.projects/project.md` file describing its Project and field names.
Setup discovers automatically whether the Project owner is a person or an organisation.

Install the skill for compatible agents:

```bash
gh skill install MiguelRodo/projects github-project-admin --agent universal --scope user
```

Then run its setup from the repository you want to use:

```bash
bash "$HOME/.agents/skills/github-project-admin/scripts/setup.sh"
```

For the `MiguelRodo/projects` repository itself, use the checked-in copy:

```bash
bash skills/github-project-admin/scripts/setup.sh --install-skill-from .
```

In Codex Cloud, put the same command in the environment's setup-script box. Add `GH_TOKEN` as an environment variable so it is available while Codex is working, and allow access to `github.com` and `api.github.com`. Do not print the token or commit it to the repository.

## Add repository-specific setup

If a repository needs extra tools, add `.projects/setup.sh`. It runs automatically after the shared setup.

If the repository needs to replace the shared setup completely, put this near the top of that file:

```bash
# github-project-admin: override
```

Keep repository-specific files under `.projects/`. Updating the shared skill will not replace them.

## Update the skill

```bash
gh skill update github-project-admin
```

## If a command fails

Paste the terminal output back into the chat, or say which command failed. The agent should check what already succeeded and give you only the corrected or remaining commands to paste.

If the failure looks reusable, the agent may offer to improve this repository. It should open a pull request only after you agree.
