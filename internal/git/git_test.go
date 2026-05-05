package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommitChangesCreatesCommit(t *testing.T) {
	tmpDir := t.TempDir()
	runGit(t, tmpDir, "init")
	runGit(t, tmpDir, "config", "user.name", "Test User")
	runGit(t, tmpDir, "config", "user.email", "test@example.com")
	writeTestFile(t, tmpDir, "file.txt", "hello")
	runGit(t, tmpDir, "add", "file.txt")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := CommitChanges("chore(release): prepare v1.2.3"); err != nil {
		t.Fatalf("CommitChanges: %v", err)
	}

	out := runGitOutput(t, tmpDir, "log", "-1", "--pretty=%s")
	if got := strings.TrimSpace(out); got != "chore(release): prepare v1.2.3" {
		t.Fatalf("commit subject = %q", got)
	}
}

func TestPushReleasePushesBranchAndTag(t *testing.T) {
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	tmpDir := t.TempDir()
	runGit(t, tmpDir, "init")
	runGit(t, tmpDir, "config", "user.name", "Test User")
	runGit(t, tmpDir, "config", "user.email", "test@example.com")
	runGit(t, tmpDir, "branch", "-M", "main")
	runGit(t, tmpDir, "remote", "add", "origin", remoteDir)

	writeTestFile(t, tmpDir, "file.txt", "hello")
	runGit(t, tmpDir, "add", "file.txt")
	runGit(t, tmpDir, "commit", "-m", "chore: initial commit")
	runGit(t, tmpDir, "tag", "-a", "v1.2.3", "-m", "release")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := PushRelease("v1.2.3"); err != nil {
		t.Fatalf("PushRelease: %v", err)
	}

	branchRef := runGitOutput(t, remoteDir, "rev-parse", "refs/heads/main")
	tagRef := runGitOutput(t, remoteDir, "rev-parse", "refs/tags/v1.2.3^{commit}")
	if strings.TrimSpace(branchRef) == "" || strings.TrimSpace(tagRef) == "" {
		t.Fatal("expected branch and tag to be pushed")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, string(out))
	}
}

func runGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
