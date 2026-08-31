#!/usr/bin/env bash

# Generic environment preflight for github-project-admin. The host supplies
# credentials and network access. This script never prints or stores a token.
set +x
set -Eeuo pipefail

readonly GH_CLI_VERSION="2.98.0"
readonly GH_CLI_MINIMUM="2.90.0"
readonly GH_CLI_SHA256_AMD64="3b8ac6b30336802fc1a858d7c084e11cdf24ac1a761ca90b68022d7d729208de"
readonly GH_CLI_SHA256_ARM64="cf689084f3a3618f7eae4a2420d335d74626d65f5e594b9828d125d69f800d86"

setup_tmp_dir=""
repository=""
project_owner=""
project_number=""
project_title=""
contract_root=""
discover_repository=true
skip_install=false

die() {
  echo "ERROR: $*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage: setup.sh [options]

Options:
  --repository OWNER/REPO    Verify access to this repository.
  --no-repository            Do not discover or verify a repository.
  --project-owner LOGIN      Verify this user or organisation Project owner.
  --project-number NUMBER    Verify this Project number.
  --project-title TITLE      Also verify the exact Project title.
  --contract-root DIRECTORY  Validate DIRECTORY/.projects/project.md.
  --no-contract              Skip automatic local contract validation.
  --skip-install             Do not install or upgrade gh.
  --help                     Show this help.
EOF
}

while (($#)); do
  case "$1" in
    --repository)
      (($# >= 2)) || die "--repository requires OWNER/REPO"
      repository="$2"
      discover_repository=false
      shift 2
      ;;
    --no-repository)
      discover_repository=false
      shift
      ;;
    --project-owner)
      (($# >= 2)) || die "--project-owner requires a login"
      project_owner="$2"
      shift 2
      ;;
    --project-number)
      (($# >= 2)) || die "--project-number requires a number"
      project_number="$2"
      shift 2
      ;;
    --project-title)
      (($# >= 2)) || die "--project-title requires a title"
      project_title="$2"
      shift 2
      ;;
    --contract-root)
      (($# >= 2)) || die "--contract-root requires a directory"
      contract_root="$2"
      shift 2
      ;;
    --no-contract)
      contract_root="-"
      shift
      ;;
    --skip-install)
      skip_install=true
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      die "unknown option: $1"
      ;;
  esac
done

if [[ -n "$project_owner" || -n "$project_number" || -n "$project_title" ]]; then
  [[ -n "$project_owner" && -n "$project_number" ]] ||
    die "--project-owner and --project-number must be supplied together"
fi
[[ -z "$project_number" || "$project_number" =~ ^[1-9][0-9]*$ ]] ||
  die "--project-number must be a positive integer"
[[ -z "$repository" || "$repository" =~ ^[^/[:space:]]+/[^/[:space:]]+$ ]] ||
  die "--repository must use OWNER/REPO form"

cleanup() {
  if [[ -n "$setup_tmp_dir" && "$setup_tmp_dir" != "/" && -d "$setup_tmp_dir" ]]; then
    rm -rf -- "$setup_tmp_dir"
  fi
}
trap cleanup EXIT

version_at_least() {
  local actual="${1%%-*}" required="$2" index
  local -a actual_parts required_parts
  IFS=. read -r -a actual_parts <<<"$actual"
  IFS=. read -r -a required_parts <<<"$required"
  for index in 0 1 2; do
    local actual_part="${actual_parts[$index]:-0}"
    local required_part="${required_parts[$index]:-0}"
    [[ "$actual_part" =~ ^[0-9]+$ && "$required_part" =~ ^[0-9]+$ ]] || return 1
    ((actual_part > required_part)) && return 0
    ((actual_part < required_part)) && return 1
  done
  return 0
}

install_gh() {
  local architecture archive_name archive_path expected_sha extracted_binary
  local install_dir user_directory path_line shell_rc

  case "$(uname -m)" in
    x86_64)
      architecture="amd64"
      expected_sha="$GH_CLI_SHA256_AMD64"
      ;;
    aarch64|arm64)
      architecture="arm64"
      expected_sha="$GH_CLI_SHA256_ARM64"
      ;;
    *)
      die "unsupported architecture for automatic gh installation: $(uname -m)"
      ;;
  esac

  for prerequisite in curl install tar sha256sum; do
    command -v "$prerequisite" >/dev/null 2>&1 ||
      die "$prerequisite is required to install gh"
  done

  setup_tmp_dir="$(mktemp -d)"
  archive_name="gh_${GH_CLI_VERSION}_linux_${architecture}.tar.gz"
  archive_path="$setup_tmp_dir/$archive_name"
  curl --fail --location --retry 3 --silent --show-error \
    --output "$archive_path" \
    "https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/${archive_name}"
  printf '%s  %s\n' "$expected_sha" "$archive_path" |
    sha256sum --check --status ||
    die "downloaded gh archive failed SHA-256 verification"

  tar --no-same-owner -xzf "$archive_path" -C "$setup_tmp_dir"
  extracted_binary="$setup_tmp_dir/gh_${GH_CLI_VERSION}_linux_${architecture}/bin/gh"
  [[ -x "$extracted_binary" ]] ||
    die "downloaded gh archive did not contain the expected binary"

  user_directory=""
  if [[ -d /usr/local/bin && -w /usr/local/bin ]]; then
    install_dir="/usr/local/bin"
  elif [[ -d /usr/bin && -w /usr/bin ]]; then
    install_dir="/usr/bin"
  else
    command -v getent >/dev/null 2>&1 ||
      die "getent is required to locate a user tool directory"
    user_directory="$(getent passwd "$(id -u)" | cut -d: -f6)"
    [[ -n "$user_directory" && -d "$user_directory" ]] ||
      die "could not resolve a writable user directory for gh"
    install_dir="$user_directory/.local/bin"
  fi

  mkdir -p "$install_dir"
  install -m 0755 "$extracted_binary" "$install_dir/gh"
  export PATH="$install_dir:$PATH"

  if [[ -n "$user_directory" ]]; then
    path_line="$(printf 'export PATH=%q:$PATH' "$install_dir")"
    for shell_rc in "$user_directory/.profile" "$user_directory/.bashrc"; do
      touch "$shell_rc"
      grep -Fqx "$path_line" "$shell_rc" || printf '%s\n' "$path_line" >>"$shell_rc"
    done
  fi
}

gh_is_usable=false
if command -v gh >/dev/null 2>&1; then
  installed_version="$(gh version 2>/dev/null | awk 'NR == 1 { print $3 }')"
  if version_at_least "$installed_version" "$GH_CLI_MINIMUM"; then
    gh_is_usable=true
  fi
fi

if [[ "$gh_is_usable" != true ]]; then
  [[ "$skip_install" == false ]] ||
    die "gh $GH_CLI_MINIMUM or newer is required"
  install_gh
fi

command -v gh >/dev/null 2>&1 || die "gh installation failed"
gh project --help >/dev/null 2>&1 || die "this gh build does not provide Project commands"

login="$(gh api user --jq .login 2>/dev/null)" ||
  die "GitHub authentication failed; configure a task-scoped GH_TOKEN or authenticate gh"
[[ -n "$login" ]] || die "GitHub authentication returned no login"

if [[ "$discover_repository" == true && -z "$repository" ]] &&
   command -v git >/dev/null 2>&1 && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  repository="$(gh repo view --json nameWithOwner --jq .nameWithOwner 2>/dev/null || true)"
fi

if [[ -n "$repository" ]]; then
  observed_repository="$(gh repo view "$repository" --json nameWithOwner --jq .nameWithOwner)"
  [[ "$observed_repository" == "$repository" ]] ||
    die "repository identity mismatch"
  gh issue list --repo "$repository" --limit 1 --json number >/dev/null
fi

if [[ -n "$project_owner" ]]; then
  readonly project_query='query($login: String!, $number: Int!) {
    user(login: $login) { projectV2(number: $number) { id number title } }
    organization(login: $login) { projectV2(number: $number) { id number title } }
  }'
  match_count="$(gh api graphql -f query="$project_query" \
    -F login="$project_owner" -F number="$project_number" \
    --jq '[.data.user.projectV2, .data.organization.projectV2] | map(select(. != null)) | length')"
  [[ "$match_count" == "1" ]] ||
    die "expected exactly one readable Project for the supplied owner and number"

  if [[ -n "$project_title" ]]; then
    observed_title="$(gh api graphql -f query="$project_query" \
      -F login="$project_owner" -F number="$project_number" \
      --jq '[.data.user.projectV2, .data.organization.projectV2] | map(select(. != null)) | .[0].title')"
    [[ "$observed_title" == "$project_title" ]] ||
      die "Project title mismatch"
  fi
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ "$contract_root" != "-" ]]; then
  if [[ -n "$contract_root" ]]; then
    bash "$script_dir/validate-contract.sh" "$contract_root"
  elif [[ -f .projects/project.md ]]; then
    bash "$script_dir/validate-contract.sh" .
  fi
fi

echo "GitHub Project-administration preflight passed for $login."
[[ -z "$repository" ]] || echo "Verified repository: $repository."
[[ -z "$project_owner" ]] || echo "Verified Project: $project_owner/$project_number."
