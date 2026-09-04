# `projects` CLI

`projects` handles the repeated, mechanical parts of GitHub Project
administration. The repository's `.projects/` contract still decides which
repository and Project are in scope. The `github-project-admin` skill still
interprets a user's request and applies its safety rules.

The CLI is optional. Keep using the repository scripts or direct `gh` and API
operations when it is not installed or does not support the requested change.

## Install on Ubuntu or Debian

Miguel's signed APT repository is already hosted on GitHub Pages. These are
one-time, human-run setup commands:

```bash
sudo apt-get install -y curl gpg
curl -fsSL https://miguelrodo.github.io/apt-miguelrodo/KEY.gpg \
  | sudo gpg --dearmor --yes -o /usr/share/keyrings/apt-miguelrodo.gpg
echo "deb [signed-by=/usr/share/keyrings/apt-miguelrodo.gpg] https://miguelrodo.github.io/apt-miguelrodo stable main" \
  | sudo tee /etc/apt/sources.list.d/apt-miguelrodo.list >/dev/null
sudo apt-get update
sudo apt-get install -y projects
```

`apt-get update` does not install anything, so it does not need `-y`.

To upgrade later:

```bash
sudo apt-get update
sudo apt-get install --only-upgrade -y projects
```

Agents should not run these `sudo` commands merely to check a version. Use one
of the read-only checks instead:

```bash
projects update check
apt-cache policy projects
```

Tagged releases also contain checksummed Linux, macOS and Windows archives.
If Go is already installed, the current source can be installed with:

```bash
go install github.com/MiguelRodo/projects/cmd/projects@latest
```

Maintainer setup and version-bump commands are in the
[`projects` release guide](releasing.md).

## Validate a repository contract

Run the command in the repository root:

```bash
projects contract validate
```

Or name another checkout without changing directory:

```bash
projects contract validate --root /path/to/repository
```

This validates the complete single-Project contract or dispatcher, including
its child contracts. It does not contact GitHub or change files. The existing
shell validator remains available:

```bash
bash .agents/skills/github-project-admin/scripts/validate-contract.sh .
```

Use `--json` when another program needs the resolved contract summary. Progress
is written to stderr, while JSON is written to stdout. `--quiet` hides progress.

## Read every Project item

For a single-Project repository:

```bash
projects project item-list --format json
```

The shorter `--json` flag is equivalent. Human-readable table output is the
default:

```bash
projects project item-list
```

The command checks the local contract, reads the live Project identity, asks
GitHub CLI for a deliberately large item set, then compares the returned items
with GitHub's reported total. It fails instead of describing a partial page as
the whole Project.

A dispatcher needs an exact route selector. Any identifiers supplied together
must agree:

```bash
projects project item-list --project-key personal --json
projects project item-list --routing-label project:personal --json
projects project item-list --project-number 40 --json
```

GitHub commands use the current authenticated `gh` account. Check it with:

```bash
gh auth status
```

## Version and update checks

```bash
projects version
projects version --json
projects update check
projects update check --json
```

`projects update check` only reads the latest GitHub Release. It does not run a
package manager, install a binary or require `sudo`.

## Output and failures

Normal commands print a small number of numbered stages to stderr so a person
or agent can see where work has reached. Command data goes to stdout. This keeps
JSON safe to pipe into another tool:

```bash
projects project item-list --json >project-items.json
```

Usage errors exit with status 2. Validation, GitHub and completeness failures
exit with status 1. A failed command names the stage that failed and includes
the underlying `gh` error when one exists.

The first release is read-only. Issue creation, Project membership and field
updates will be added only with plan-first behaviour, stale checks and separate
readback. Until then, the skill's documented direct GitHub path remains the
supported mutation route.
