package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	binaryOnce sync.Once
	binaryPath string
	binaryErr  error
)

func TestDryRunDoesNotCreateArtifacts(t *testing.T) {
	repo := newGitRepo(t)
	initCommitAndTag(t, repo, "v0.1.0")
	writeFile(t, repo, "feature.txt", "new feature")
	git(t, repo, "add", "feature.txt")
	git(t, repo, "commit", "-m", "feat: add feature")
	initialHead := strings.TrimSpace(gitOutput(t, repo, "rev-parse", "HEAD"))
	initialTags := strings.TrimSpace(gitOutput(t, repo, "tag", "--list"))

	out := runRelease(t, repo, "--dry-run")
	if !strings.Contains(out, "[DRY RUN] Skipping actual release creation") {
		t.Fatalf("dry-run output missing marker:\n%s", out)
	}

	if got := strings.TrimSpace(gitOutput(t, repo, "rev-parse", "HEAD")); got != initialHead {
		t.Fatalf("HEAD changed on dry-run: got %s, want %s", got, strings.TrimSpace(initialHead))
	}
	if got := strings.TrimSpace(gitOutput(t, repo, "tag", "--list")); got != initialTags {
		t.Fatalf("tags changed on dry-run: got %q, want %q", got, strings.TrimSpace(initialTags))
	}
}

func TestFlagReleaseCreatesCommitAndTagWithoutPushPrompt(t *testing.T) {
	repo := newGitRepo(t)
	initCommitAndTag(t, repo, "v0.1.0")
	writeFile(t, repo, "feature.txt", "new feature")
	git(t, repo, "add", "feature.txt")
	git(t, repo, "commit", "-m", "feat: add feature")

	out := runRelease(t, repo, "--patch")
	if strings.Contains(out, "Push tag to remote?") {
		t.Fatalf("unexpected push prompt in output:\n%s", out)
	}
	if !strings.Contains(out, "✅ Created tag: v0.1.1") {
		t.Fatalf("missing tag creation output:\n%s", out)
	}

	if got := strings.TrimSpace(gitOutput(t, repo, "log", "-1", "--pretty=%s")); got != "chore(release): prepare v0.1.1" {
		t.Fatalf("release commit subject = %q", got)
	}
	if got := strings.TrimSpace(gitOutput(t, repo, "tag", "--list")); !strings.Contains(got, "v0.1.1") {
		t.Fatalf("release tag missing: %q", got)
	}
}

func runRelease(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.CommandContext(context.Background(), releaseBinary(t), args...)
	cmd.Dir = repo
	cmd.Stdin = strings.NewReader("\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("release %v: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func releaseBinary(t *testing.T) string {
	t.Helper()
	binaryOnce.Do(func() {
		root := repoRoot(t)
		tmp, err := os.CreateTemp("", "shikai-e2e-*")
		if err != nil {
			binaryErr = err
			return
		}
		_ = tmp.Close()
		binaryPath = tmp.Name()
		cmd := exec.Command("go", "build", "-o", binaryPath, ".")
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			binaryErr = err
			t.Logf("go build output:\n%s", string(out))
		}
	})
	if binaryErr != nil {
		t.Fatalf("build release binary: %v", binaryErr)
	}
	return binaryPath
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller failed")
	}
	return filepath.Dir(filepath.Dir(file))
}

func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.name", "Test User")
	git(t, dir, "config", "user.email", "test@example.com")
	return dir
}

func initCommitAndTag(t *testing.T, dir, tag string) {
	t.Helper()
	writeFile(t, dir, "README.md", "hello")
	git(t, dir, "add", "README.md")
	git(t, dir, "commit", "-m", "feat: initial")
	git(t, dir, "tag", tag)
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	return git(t, dir, args...)
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
