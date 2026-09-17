# Project bootstrap

This repository is Miguel's public, end-to-end project bootstrap procedure. It
coordinates four systems without merging their responsibilities:

- Google Drive holds the canonical project description and LLM-oriented source
  material.
- GitHub holds the repository, Project and authoritative work items.
- `MiguelRodo/github-projects-skill` supplies the reusable `github-projects` workflow.
- A private ChatGPT Project provides the conversational workspace.

The bootstrapper creates resources directly when its current environment has
the required authenticated capability. When it cannot make a GitHub change, it
returns the smallest executable `gh` command block and verifies the result after
the user runs it. ChatGPT Project creation remains a verified manual checklist
in version 1.

The shorter public guide is at
[miguelrodo.github.io/projects](https://miguelrodo.github.io/projects/).
This repository remains the front door; it links to `MiguelRodo/github-projects-skill` for
the GitHub administration skill and CLI.

The website includes guides for the
[full Drive, registry and ChatGPT setup](https://miguelrodo.github.io/projects/#full-setup),
[adding another GitHub Project](https://miguelrodo.github.io/projects/#add-project)
and [setting up the local implementation queue](https://miguelrodo.github.io/projects/#implementation-queue).

## Optional `projects` CLI

`MiguelRodo/github-projects-skill` now has an optional Go CLI for repeated GitHub operations
that benefit from one tested implementation. It validates the existing
`.projects/` contract and can read a complete Project item set without relying
on `gh project item-list`'s default page.

There is no separate project-bootstrap CLI. This repository still coordinates
GitHub, Drive, the registry and the ChatGPT Project handoff through its skill.
It can use `projects` for a command the binary supports, while keeping the
current scripts and provider operations as fallbacks.

`pj` also stays separate. It chooses and launches a local agent; `projects`
performs deterministic administration. Installation and signed APT repository
instructions are in the
[`projects` CLI guide](https://github.com/MiguelRodo/github-projects-skill/blob/main/docs/cli.md).

## Install the skill

With a current GitHub CLI:

```bash
gh skill install MiguelRodo/projects project-bootstrap \
  --agent universal --scope user
```

Then ask, for example:

```text
Use $project-bootstrap to set up a new research project called Foobar.
```

The same skill can resume or audit a partial setup:

```text
Use $project-bootstrap to finish setting up Foobar and verify the result.
```

If skill installation is unavailable, point a capable ChatGPT or Codex chat to
this repository and ask it to follow [`PROJECT_BOOTSTRAP.md`](PROJECT_BOOTSTRAP.md).

## One-time local operator launcher

Operator bootstrap is separate from creating an individual project. The local
operator tooling (`pj`, `pja`, `pjcp`, `pjcd`, `pj-update-skills`) is maintained
in its standalone canonical repository [`MiguelRodo/pj`](https://github.com/MiguelRodo/pj).

Run the installer once per machine from the current `v0` release line:

```bash
git clone --depth 1 --branch v0 --single-branch https://github.com/MiguelRodo/pj.git
cd pj
bash install.sh
```

Alternatively, if `setupmjr` is installed:

```bash
setupmjr project --pj
```

It installs four launcher names into a sensible per-user executable directory:

- `pj` uses the configured default backend;
- `pja` always selects Google Antigravity;
- `pjcp` always selects GitHub Copilot CLI;
- `pjcd` always selects Codex.

It also installs `pj-update-skills`, the maintenance command for refreshing the
shared `github-projects` skill across managed repositories.

The installer prefers `~/.local/bin` when it is already on `PATH`, then `~/bin`
when that is the configured standard user bin directory. If neither is on
`PATH`, it prefers an existing `~/.local/bin` or `~/bin`, in that order, and
otherwise creates `~/.local/bin`. The selected directory is recorded in
`${XDG_CONFIG_HOME:-~/.config}/pj/install-bin-dir` so a later reinstall can
safely migrate an older installer-managed location rather than leave a stale
launcher shadowing the current one.

Set `PJ_BIN_DIR` to choose a different absolute or home-relative location for an
install, for example:

```bash
PJ_BIN_DIR='~/.local/bin' bash install.sh
PJ_BIN_DIR='~/tools/bin' bash install.sh
```

A new install starts with Codex as the `pj` default. Change or inspect that saved
default with:

```bash
pj --set-default antigravity
pj --set-default copilot
pj --set-default codex
pj --show-default
```

The saved choice lives at `${XDG_CONFIG_HOME:-~/.config}/pj/default-backend` and
is preserved when the installer is rerun. `PJ_DEFAULT_BACKEND` can override the
saved default for the current environment, while `PJ_BACKEND` remains the
one-run generic `pj` override. An explicit `--backend` flag has the highest
precedence and can also override a shorthand launcher.

### Per-backend model defaults

`pj` also keeps model selection separate for each backend. The built-in model
defaults are:

- Codex: `gpt-5.6-luna`, with the existing `xhigh` reasoning effort;
- GitHub Copilot CLI: `mai-code-1.1-flash` (MAI-Code-1.1-Flash);
- Google Antigravity: no `pj` model pin, so Antigravity follows its own current
  or configured default model instead of freezing the launcher to one Flash
  release.

Inspect all effective model defaults or one backend with:

```bash
pj --show-models
pj --show-model copilot
```

Save a different default for a backend with:

```bash
pj --set-model copilot gpt-5.6-terra
pj --set-model codex gpt-5.6-luna
pj --set-model antigravity gemini-3.8-flash-high
```

Saved model choices live under `${XDG_CONFIG_HOME:-~/.config}/pj/models/` and
are preserved across installer reruns. Reset one backend to the built-in `pj`
default with:

```bash
pj --reset-model copilot
pj --reset-model codex
pj --reset-model antigravity
```

Resetting Antigravity removes the launcher-level model pin and returns it to the
provider's own current/default model. This is intentional: when the desired
behaviour is simply "use the current Flash default", leaving Antigravity
unpinned allows provider updates to advance that choice without editing `pj`.

The per-backend environment variables `PJ_CODEX_MODEL`,
`PJ_ANTIGRAVITY_MODEL`, and `PJ_COPILOT_MODEL` override saved model choices for
the current environment. Backend-native model flags also remain available for
one run by placing options before a literal `--`, for example:

```bash
pja --model gemini-3.8-flash-medium -- "Review issue #12"
pjcp --model gpt-5.6-terra -- "Review issue #12"
pjcd -m gpt-5.6-sol -- "Review issue #12"
```

### Implement queued Chat handoffs

ChatGPT or another provider can leave a small `pj:implement-chat` issue when it
has authority for a GitHub change but cannot perform the final Project-side
mutation itself. Process those local handoffs with any of these equivalent
forms:

```bash
pj -i
pj --implement-issues
pj --implement-chat
```

The launcher itself does not contain the queue protocol. It converts the flag
into a short request for the configured backend, which then follows the
`github-projects` local implementation-queue guidance and each repository's
`AGENTS.md` and `.projects` contract.

Trusted queue items created by the currently authenticated GitHub user and
backed by the required unedited authority comment are handled without a routine
preview. Requests from another author, edited or missing authority comments,
and suspicious or genuinely ambiguous requests are reviewed first and require
local operator approval before execution. Project credentials remain local to
the operator rather than being stored in collaborator-controlled Actions
workflows.

Use a backend override normally when desired, for example:

```bash
pj --backend copilot -i
pjcp -i
```

### Conversational sessions

When `pj`, `pja`, `pjcp` or `pjcd` is run from an ordinary terminal with an
initial prompt, the default behaviour is conversational rather than one-shot.
The first request is executed and the selected agent remains open so follow-up
prompts can be typed directly into the same TUI. For example:

```bash
pja -- "Assess the stimgate project"
```

After Antigravity answers, stay in that Antigravity session and type a follow-up
such as:

```text
Tackle the last two things you suggested.
```

Do not type `pj` again while still inside the agent TUI. Exiting the TUI returns
to the shell.

The launcher uses each provider's supported conversation mechanism:

- Codex receives the prompt as the initial prompt to its ordinary interactive
  TUI.
- Copilot uses `-i` / `--interactive` so the initial prompt runs and the session
  remains open.
- Antigravity does not currently document an interactive initial-prompt flag,
  so `pj` seeds the request with one `agy -p` turn and immediately resumes that
  newly created workspace conversation with `agy --continue`.

In non-interactive contexts such as pipelines and command substitution, `pj`
uses the one-shot interface automatically. Force one-shot behaviour from a
terminal with `--oneshot` before any agent-specific options or prompt text:

```bash
pja --oneshot -- "Assess the stimgate project"
pjcp --oneshot -- "Assess the stimgate project"
pjcd --oneshot -- "Assess the stimgate project"
```

`PJ_SESSION_MODE` may be set to `auto`, `interactive` or `oneshot`. `auto` is the
default and chooses interactive mode only when stdin and stdout are attached to
a terminal. An explicit `--oneshot` overrides `PJ_SESSION_MODE=interactive` for
that invocation.

The installer treats `~/AGENTS.md` as the canonical cross-agent user guidance
and exposes its managed block through `~/.codex/AGENTS.md`,
`~/.copilot/copilot-instructions.md` and `~/.gemini/GEMINI.md`. An absent, empty
or installer-owned backend file becomes a symlink to the canonical file. A
regular file with genuine backend-specific content stays regular: its custom
content, file mode and block position are preserved while the canonical managed
block is refreshed in place. Unrelated symlinks and non-files are left
unchanged. A broken backend instruction symlink is reported as an installation
error rather than silently treated as configured.

The canonical home and workspace files, and the Codex rules file, may themselves
be symlinks into a dotfiles checkout. The installer updates their resolved file
targets without replacing the links or changing the target modes.

The same home guidance defines opt-in `agy` delegation: only Codex and Copilot
may use `agy`, only after explicit operator authorisation for the current task or
conversation, with `gemini-3.8-flash-high` as the preferred delegated model.
Conversation-level authorisation covers repeated calls, repository inspection
uses `--add-dir`, and the primary agent may satisfy a narrowly scoped headless
permission request on the operator's behalf. Broad wildcard permissions and
`--dangerously-skip-permissions` remain separately gated. The installer also
maintains a Codex exec-policy rule allowing the `agy` executable. See
[`MiguelRodo/pj`](https://github.com/MiguelRodo/pj) for migration and permission details.
The workspace defaults to `~/planning` and can be changed with `PJ_WORKSPACE`.
An absolute path or a quoted `~/...` path is accepted.

The same prompt-parsing rule applies to every backend. When the first argument
is ordinary text, all remaining arguments are combined into prompt text, so
later dash-prefixed fragments stay part of the request. A leading literal `--`
forces all following arguments to be prompt text. If the first argument is an
agent option, a later `--` separates agent options from one combined prompt.

For example, whatever backend `pj` currently defaults to receives this whole
request as prompt text:

```bash
pj "Create the issue - keep this dash as ordinary prompt text"
```

Google Antigravity is available through `pja` after `agy` is installed and
authenticated:

```bash
pja "Create the issue - keep this dash as text"
pja --effort medium -- "Review issue #12 - do not edit it"
```

Ordinary terminal `pja` runs stay conversational as described above. The
initial Antigravity turn uses high reasoning effort by default and automatic
tool approval. `pja` does not force the terminal sandbox because Project
administration needs live `gh` and provider API access; request `--sandbox`
explicitly when that restriction is wanted. Headless seed or one-shot turns use
a 15-minute response timeout by default, configurable with
`PJ_ANTIGRAVITY_TIMEOUT` or a one-run `--print-timeout` option. Automatic tool
approval is deliberately broad, so use this operator launcher only for requests
you intend the agent to execute.

GitHub Copilot CLI is available through `pjcp` after `copilot` is installed and
authenticated:

```bash
pjcp "Create the issue - keep this dash as text"
pjcp --model auto -- "Review issue #12 - do not edit it"
```

Ordinary terminal `pjcp` runs use Copilot's interactive initial-prompt mode and
grant all tool, path and URL permissions with `--allow-all`. One-shot or
non-interactive runs use the official `copilot -p` interface instead.

Codex remains directly available through `pjcd` regardless of the saved `pj`
default:

```bash
pjcd "Review the current project"
```

Ordinary terminal Codex runs use the interactive TUI. One-shot or
non-interactive runs use `codex exec`.

The long forms remain available when useful:

```bash
PJ_DEFAULT_BACKEND=antigravity pj "Review the current project"
PJ_BACKEND=copilot pj "Review the current project"
pj --backend codex "Review the current project"
```

Repository `AGENTS.md` and `.projects` contracts remain the source of
Project-administration behaviour for every backend rather than provider-specific
copies of the task model.

Check the public website with:

```bash
bash tests/site.sh
```

To preview it locally, run:

```bash
python3 -m http.server 8000 --directory site
```

then open `http://localhost:8000`. GitHub Pages deploys the same dependency-free
files in `site/` after a change reaches `main`.

## What a completed setup contains

- `work/projects/<year>/<research|lecturing>/<project-slug>/` in Google Drive;
- a native Google Doc named `README` and a `references/llm/` folder;
- an optional, verified GitHub repository and GitHub Project;
- repository onboarding through `MiguelRodo/github-projects-skill` when both a repository
  and Project are used;
- a bounded project-resources section in the repository README;
- one verified row in the task system's Projects registry; and
- exact manual steps for creating and checking the private ChatGPT Project.

The full human-readable contract is in
[`PROJECT_BOOTSTRAP.md`](PROJECT_BOOTSTRAP.md). Operational agent instructions
live in [`skills/project-bootstrap/SKILL.md`](skills/project-bootstrap/SKILL.md).

## Repository layout

| Path | Purpose |
| --- | --- |
| `PROJECT_BOOTSTRAP.md` | Human-readable workflow and postconditions. |
| `skills/project-bootstrap/` | Installable Agent Skill. |
| `site/` | Public guide deployed through GitHub Pages. |
| `templates/drive-readme.md` | Minimal native Google Doc README shape. |
| `templates/repository-resources.md` | Bounded repository README section. |

No credentials, private Drive identifiers, contact details or project-specific
URLs belong in this public repository.
