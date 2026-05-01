# Copilot Instructions

## Commit Convention

Use Conventional Commits for all suggested or written commits.

- Format: `<type>(<scope>): <description>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `revert`
- Keep the description imperative, lowercase, and without a trailing period

Examples:
- `feat(release): add prerelease flag`
- `fix(commits): treat feat as minor bump`
- `docs(readme): document release flow`

## Build, Test, and Lint Commands

- `go build -o release .` — build the CLI binary
- `go test ./...` — run the full test suite
- `go test ./internal/commits -run TestCommitPatternMatchesConventionalCommits` — run one commit-parser test
- `go test ./internal/version -run TestBump` — run version bump tests
- `go test ./internal/changelog -run TestGenerateSimple` — run changelog tests

No dedicated lint command is currently defined in the repo.

## Architecture Overview

The CLI entrypoint is `main.go`, which calls `cmd.Execute()`. `cmd/root.go` defines global flags with cobra/viper, and `cmd/release.go` owns the release workflow.

The release flow is:
1. Read the latest git tag and strip any leading `v`
2. Collect commits since that tag
3. Parse Conventional Commits, skipping non-conforming messages with a warning
4. Recommend a semver bump from commit types and breaking changes
5. Let the user confirm or override the bump interactively, unless `--patch`, `--minor`, or `--major` is set
6. Generate a changelog with git-chglog, falling back to a simple internal generator
7. Update the manifest file version
8. Stage changes, create an annotated tag, and optionally push after confirmation

Supporting packages are split by responsibility:
- `internal/git` handles git commands
- `internal/commits` parses commit messages and computes bump severity
- `internal/version` bumps semver strings and prerelease suffixes
- `internal/changelog` writes changelog output and resolves git-chglog config/template paths
- `internal/manifest` detects and updates manifest files
- `internal/interactive` handles version selection prompts

## Conventions

- Tags are written as `vX.Y.Z`; reading the current version strips the `v` prefix
- `--dry-run` is off by default
- `--push` bypasses the confirmation prompt before pushing the release tag
- `--prerelease` with `--prerelease-id alpha|beta|rc...` appends a prerelease suffix
- The repo is expected to be run from the current working directory
- Non-conforming commits are ignored rather than coerced
- If more than one manifest is detected, the command should fail and ask for configuration instead of guessing
- Prefer small, package-local tests alongside the code they exercise
