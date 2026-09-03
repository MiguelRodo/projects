#!/usr/bin/env bash
set -Eeuo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skill_dir="$(cd "$test_dir/.." && pwd)"
repo_root="$(cd "$skill_dir/../.." && pwd)"
workflow="$repo_root/.github/workflows/project-admin-bridge.yml"
reference="$skill_dir/references/automation-bridge.md"

grep -Fq "github.event.label.name == 'automation:project-admin'" "$workflow"
grep -Fq 'copilot-requests: write' "$workflow"
grep -Fq 'npm install -g @github/copilot@1.0.82' "$workflow"
grep -Fq 'COPILOT_GITHUB_TOKEN: ${{ github.token }}' "$workflow"
grep -Fq 'GH_TOKEN: ${{ secrets.PROJECTS_TOKEN }}' "$workflow"
grep -Fq -- '--model=mai-code-1.1-flash' "$workflow"
grep -Fq -- "--allow-tool='read,shell(gh:*)'" "$workflow"
grep -Fq -- "--secret-env-vars='GH_TOKEN'" "$workflow"
grep -Fq 'PROJECT_ADMIN_RESULT: DONE' "$workflow"

if grep -Fq 'OPENAI_API_KEY' "$workflow"; then
  echo 'ERROR: bridge still requires OPENAI_API_KEY' >&2
  exit 1
fi
if grep -Fq 'openai/codex-action' "$workflow"; then
  echo 'ERROR: bridge still invokes openai/codex-action' >&2
  exit 1
fi
if grep -Fq 'apply-project-operation.mjs' "$workflow"; then
  echo 'ERROR: bridge still invokes the retired manifest executor' >&2
  exit 1
fi

grep -Fq 'Plain natural language is enough' "$reference"
grep -Fq 'gh secret set PROJECTS_TOKEN' "$reference"

echo 'project-admin bridge static tests passed'
