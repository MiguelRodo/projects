# Provider Project instructions

Use this reference when configuring an ordinary ChatGPT Project or another provider's persistent instruction surface.

Keep the instructions small. They route to the canonical skill and repository contract; they do not repeat the procedure.

Suggested text:

```text
For GitHub issue or Project-administration requests, first look for
.projects/skills/github-project-admin/SKILL.md in the target repository. If it
exists, retrieve and follow that deliberate local replacement. Otherwise
retrieve and follow the canonical github-project-admin skill from
MiguelRodo/projects/skills/github-project-admin/. Always read the target
repository's .projects/project.md and any contract it resolves.

Treat the user's prompt as the desired outcome. If this surface cannot perform
the required GitHub mutation, apply the same skill and return the smallest
executable gh command block, including independent readback. Do not ask the
user to restate the skill's operating procedure.
```

Replace the canonical repository only when the organisation has intentionally forked the specification. Do not select a local skill merely because `.projects/` exists.

Routine prompts should remain short:

- `Set example#313 to P2.`
- `What are the highest-priority open items in example?`

If a routine prompt must repeat stale checks, narrow mutation, preservation or verification, move that missing behaviour into the skill instead of lengthening the standing instructions or prompt.
