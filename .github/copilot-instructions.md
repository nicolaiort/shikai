# Copilot Instructions

## Commit Convention

Use Conventional Commits for all suggested or written commits.

- Format: `<type>(<scope>): <description>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `revert`
- Keep the description imperative, lowercase, and without a trailing period
- Group related commits with the same scope, or omit the scope if not applicable
- Don't just plainly commit all files since the last commit at once; instead, group changes into logical commits with clear messages

Examples:
- `feat(release): add prerelease flag`
- `fix(commits): treat feat as minor bump`
- `docs(readme): document release flow`

## Build, Test, and Lint Commands

- `go build -o release .` — build the CLI binary
- `go install github.com/shikai/release@latest` — install the CLI from source
- `task build` — build the CLI binary for the current platform
- `go test ./...` — run the full test suite
- `go test ./e2e` — run the CLI smoke/e2e tests
- `go test ./internal/commits -run TestCommitPatternMatchesConventionalCommits` — run one commit-parser test
- `go test ./internal/version -run TestBump` — run version bump tests
- `go test ./internal/changelog -run TestGenerateSimple` — run changelog tests

No dedicated lint command is currently defined in the repo.

## Architecture Overview

The CLI entrypoint is `main.go`, which calls `cmd.Execute()`. `cmd/root.go` defines the root `release` command, global flags with cobra/viper, and wires the release workflow directly. `cmd/release.go` owns the release logic.

The release flow is:
1. Read the latest git tag and strip any leading `v`
2. Collect commits since that tag
3. Parse Conventional Commits, skipping non-conforming messages with a warning
4. Recommend a semver bump from commit types and breaking changes
5. Let the user confirm or override the bump interactively, unless `--patch`, `--minor`, or `--major` is set
6. Generate a changelog with git-chglog, falling back to a simple internal generator
7. Update the manifest file version
8. Stage changes, create a release commit, create an annotated tag, and optionally push after confirmation

Supporting packages are split by responsibility:
- `internal/git` handles git commands
- `internal/commits` parses commit messages and computes bump severity
- `internal/version` bumps semver strings and prerelease suffixes
- `internal/changelog` uses the `github.com/git-chglog/git-chglog` library, resolves config/template paths, and writes changelog output
- `internal/manifest` detects and updates manifest files
- `internal/interactive` handles version selection prompts
- `.github/workflows/release.yml` builds tagged binaries for linux/windows/darwin on amd64 and arm64, then publishes a GitHub release

## Conventions

- Tags are written as `vX.Y.Z`; reading the current version strips the `v` prefix
- `--dry-run` is off by default
- `--push` bypasses the confirmation prompt before pushing the release commit and tag
- When `--patch`, `--minor`, or `--major` is used, skip the push prompt unless `--push` is also set
- When `--patch`, `--minor`, or `--major` is used without `--push`, do not push the release commit or tag
- `--prerelease` with `--prerelease-id alpha|beta|rc...` appends a prerelease suffix
- An optional `.shikai.yml` in the repo root can set defaults such as `push: true`
- `.shikai.yml` can also define hook commands under `hooks.before`, `hooks.after-changelog`, `hooks.after-tag`, and `hooks.after-done`
- Hook commands receive `SHIKAI_TAG` with the release tag being created
- The version chooser is an arrow-key select prompt and should fall back to the recommended bump when stdin/stdout are not terminals
- Invoke the CLI as `release`; there is no nested `release release` subcommand
- `shikai changelog` prints the current release notes body to stdout for reuse in release tooling
- The repo is expected to be run from the current working directory
- Non-conforming commits are ignored rather than coerced
- If more than one manifest is detected, the command should fail and ask for configuration instead of guessing
- Prefer small, package-local tests alongside the code they exercise
