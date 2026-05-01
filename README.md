# shikai

`shikai` is a Conventional Commits-based release CLI. It inspects commits since the last tag, recommends a version bump, generates a changelog, creates a release commit and annotated tag, and can push the tag to the remote.

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

## Usage

Run the CLI from the root of the git repository you want to release:

```bash
shikai
```

Common flags:

- `--patch`, `--minor`, `--major` to pick the version bump directly
- `--dry-run` to show what would happen without changing anything
- `--push` to push the tag without prompting
- `--prerelease` and `--prerelease-id` to create a prerelease tag
- `--changelog-path` to write the changelog somewhere other than `CHANGELOG.md`
- `--template` to use a custom git-chglog template

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
