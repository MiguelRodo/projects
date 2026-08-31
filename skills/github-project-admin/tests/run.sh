#!/usr/bin/env bash

set -Eeuo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skill_dir="$(cd "$test_dir/.." && pwd)"
validator="$skill_dir/scripts/validate-contract.sh"
setup="$skill_dir/scripts/setup.sh"

bash -n "$validator"
bash -n "$setup"

bash "$validator" "$test_dir/fixtures/single"
bash "$validator" "$test_dir/fixtures/dispatcher"
if bash "$validator" "$test_dir/fixtures/invalid" >/dev/null 2>&1; then
  echo "ERROR: lossy Priority mapping unexpectedly validated" >&2
  exit 1
fi
if bash "$validator" "$test_dir/fixtures/invalid-dispatcher" >/dev/null 2>&1; then
  echo "ERROR: mismatched dispatcher leaf unexpectedly validated" >&2
  exit 1
fi

grep -Fqx 'name: github-project-admin' "$skill_dir/SKILL.md"
grep -Fq 'Set example#313 to P2.' "$test_dir/short-requests.md"
grep -Fq 'P3 | Low' "$skill_dir/SKILL.md"

test_tmp_dir="$(mktemp -d)"
cleanup() {
  if [[ -n "$test_tmp_dir" && "$test_tmp_dir" != "/" && -d "$test_tmp_dir" ]]; then
    rm -rf -- "$test_tmp_dir"
  fi
}
trap cleanup EXIT

mkdir -p "$test_tmp_dir/bin"
cat >"$test_tmp_dir/bin/gh" <<'EOF'
#!/usr/bin/env bash
set -Eeuo pipefail
case "${1:-}" in
  version)
    echo "gh version 2.98.0 (test)"
    ;;
  project)
    [[ "${2:-}" == "--help" ]]
    ;;
  api)
    if [[ "${2:-}" == "user" ]]; then
      echo "octocat"
    elif [[ "${2:-}" == "graphql" ]]; then
      if [[ "$*" == *"| length"* ]]; then
        echo "1"
      else
        echo "Example planning"
      fi
    else
      exit 2
    fi
    ;;
  repo)
    [[ "${2:-}" == "view" ]]
    echo "octo-org/example"
    ;;
  issue)
    [[ "${2:-}" == "list" ]]
    echo "[]"
    ;;
  skill)
    [[ "${2:-}" == "install" ]]
    ;;
  *)
    exit 2
    ;;
esac
EOF
chmod 0755 "$test_tmp_dir/bin/gh"

secret_value="test-token-must-not-appear"
PATH="$test_tmp_dir/bin:$PATH" GH_TOKEN="$secret_value" \
  bash -x "$setup" --skip-install --no-contract \
  --repository octo-org/example \
  --project-owner octo-org --project-number 12 \
  --project-title "Example planning" \
  --install-skill-from octo-org/project-skills --agent codex \
  >"$test_tmp_dir/setup.log" 2>&1

if grep -Fq "$secret_value" "$test_tmp_dir/setup.log"; then
  echo "ERROR: setup output exposed GH_TOKEN" >&2
  exit 1
fi
grep -Fq 'preflight passed' "$test_tmp_dir/setup.log"
grep -Fq 'Verified repository: octo-org/example.' "$test_tmp_dir/setup.log"
grep -Fq 'Verified Project: octo-org/12.' "$test_tmp_dir/setup.log"

echo "github-project-admin tests passed"
