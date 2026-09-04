# Release `projects`

Releases use the Go-specific action in `MiguelRodo/actions`. The workflow creates
the semantic version tag and its floating aliases, runs GoReleaser, attaches
checksummed archives and Debian packages to the GitHub Release, then publishes
the `.deb` files to `MiguelRodo/apt-miguelrodo`.

## One-time repository setup

The `projects` repository currently needs two Actions secrets before its first
release:

- `APT_REPO_TOKEN`: a fine-grained token with Contents read and write access to
  `MiguelRodo/apt-miguelrodo`. The ordinary `GITHUB_TOKEN` cannot push to a
  different repository.
- `GPG_PRIVATE_KEY`: the ASCII-armoured private key whose public half is already
  published as `apt-miguelrodo/KEY.gpg`.

If the private key has a passphrase, also add
`GPG_PRIVATE_KEY_PASSPHRASE`. Do not paste any of these values into an issue,
commit, workflow file or agent prompt.

GitHub CLI can store the secrets without putting their values in the command
line. Run the first and third commands interactively; replace the key path in
the second command with the local private-key file:

```bash
gh secret set APT_REPO_TOKEN --repo MiguelRodo/projects
gh secret set GPG_PRIVATE_KEY --repo MiguelRodo/projects < /path/to/private-key.asc
gh secret set GPG_PRIVATE_KEY_PASSPHRASE --repo MiguelRodo/projects
```

The passphrase secret may be omitted for an unprotected signing key. Check only
the configured secret names with:

```bash
gh secret list --repo MiguelRodo/projects
```

## Publish a release

The first release has no earlier semantic version tag, so give it an exact
version:

```bash
gh workflow run go-version-release.yml \
  --repo MiguelRodo/projects \
  -f version=0.1.0
```

Later releases can request a bump:

```bash
gh workflow run go-version-release.yml \
  --repo MiguelRodo/projects \
  -f bump_type=patch
```

Use `minor` or `major` in the same way. The workflow also accepts the inputs in
the GitHub Actions web interface. Watch a run with:

```bash
gh run watch --repo MiguelRodo/projects
```

Do not set both `version` and `bump_type`. `version_force` is an escape hatch
for an intentional version jump or downgrade, and should normally remain
false.

## Verify the result

Check the GitHub Release, then verify that the signed APT repository contains
the same package version:

```bash
gh release view --repo MiguelRodo/projects
curl -fsSL https://miguelrodo.github.io/apt-miguelrodo/dists/stable/main/binary-amd64/Packages \
  | awk '/^Package: projects$/{show=1} show{print} show && /^$/{exit}'
```

On a machine that has added the APT source, this read-only command shows the
installed and candidate versions without `sudo`:

```bash
apt-cache policy projects
```

GoReleaser configuration can be checked locally with:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

Snapshot mode builds packages but does not publish them.
