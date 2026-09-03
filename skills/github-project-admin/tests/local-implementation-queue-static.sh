#!/usr/bin/env bash
set -Eeuo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skill_dir="$(cd "$test_dir/.." && pwd)"
repo_root="$(cd "$skill_dir/../.." && pwd)"
reference="$skill_dir/references/local-implementation-queue.md"
contract="$repo_root/.projects/project.md"

[ -f "$reference" ]
[ ! -e "$repo_root/.github/workflows/project-admin-bridge.yml" ]

grep -Fq 'pj:implement-chat' "$reference"
grep -Fq 'PJ implementation authority:' "$reference"
grep -Fq 'created_at == updated_at' "$reference"
grep -Fq 'do not ask the operator for a routine preview or confirmation' "$reference"
grep -Fq 'Chat implementation label | pj:implement-chat' "$contract"

if grep -R -Fq 'PROJECTS_TOKEN' "$repo_root/.github/workflows" 2>/dev/null; then
  echo 'ERROR: repository workflow still references PROJECTS_TOKEN' >&2
  exit 1
fi

if grep -R -Fq 'OPENAI_API_KEY' "$repo_root/.github/workflows" 2>/dev/null; then
  echo 'ERROR: repository workflow still references OPENAI_API_KEY' >&2
  exit 1
fi

echo 'local implementation queue static tests passed'
