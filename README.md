# shikai

`shikai` is a Conventional Commits-based release CLI. It inspects commits since the last tag, recommends a version bump, generates a changelog, creates a release commit and annotated tag, and can push the tag to the remote.

## Motivation

I created this to finally get rid of a bunch of simmilar tasks/scripts sprinkled throughout almost every repo I work in and written in a wild mixture of bash, make, taskfile and JS.
It ain't much, but it frees up time wasted by broken release scripts.

## Install

Download the binary for your platform from the [latest release](https://github.com/nicolaiort/shikai/releases/latest).

Choose the asset that matches your OS and architecture:

- `shikai-linux-amd64`
- `shikai-linux-arm64`
- `shikai-darwin-amd64`
- `shikai-darwin-arm64`
- `shikai-windows-amd64.exe`
- `shikai-windows-arm64.exe`

After downloading, make it executable on Unix-like systems:

```bash
chmod +x shikai-*
```

Then move it somewhere on your `PATH`.

### Helper script

On Linux and macOS, you can install the latest matching release binary with Bash:

```bash
bash ./scripts/install-latest-release.sh
```

By default, the script installs to `~/bin`. Pass `INSTALL_DIR=/some/path` or a path argument to use a different location.

## Configuration

### Repository config

You can add an optional `.shikai.yml` file in the repository root to set defaults such as always pushing tags:

```yaml
push: true
```

Any supported settings can be added later without changing the CLI shape.
Start from `shikai.sample.yml` in the repo root and copy it to `.shikai.yml`.

### Hooks

`shikai` can run shell hooks from `.shikai.yml`:

- `hooks.before`
- `hooks.after-changelog`
- `hooks.after-tag`
- `hooks.after-done`

Hooks are skipped during `--dry-run`; the CLI prints the commands it would have run instead.

## Usage

Run the CLI from the root of the git repository you want to release:

```bash
shikai
```

Generate the current release notes body for stdout:

```bash
shikai changelog > release-notes.md
```

Common flags:

- `--patch`, `--minor`, `--major` to pick the version bump directly
- `--dry-run` to show what would happen without changing anything
- `--push` to push the tag without prompting
- `--prerelease` and `--prerelease-id` to create a prerelease tag
- `--changelog-path` to write the changelog somewhere other than `CHANGELOG.md`
- `--template` to use a custom git-chglog template
- `--push` pushes the release commit and tag to the current branch remote

Example:

```bash
shikai --minor
```

## Development

Build locally:

```bash
go build -o shikai .
```

Or build for the current platform:

```bash
task build
```

Run tests:

```bash
go test ./...
```
