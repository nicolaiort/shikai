package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunChangelogPrintsCurrentReleaseNotes(t *testing.T) {
	tmpDir := t.TempDir()
	initGitRepo(t, tmpDir)

	writeFile(t, tmpDir, "README.md", "init")
	gitCmd(t, tmpDir, "add", "README.md")
	gitCmd(t, tmpDir, "commit", "-m", "chore: initial commit")
	gitCmd(t, tmpDir, "tag", "v0.1.0")

	writeFile(t, tmpDir, "feature.txt", "feat")
	gitCmd(t, tmpDir, "add", "feature.txt")
	gitCmd(t, tmpDir, "commit", "-m", "feat: add feature")

	writeFile(t, tmpDir, "fix.txt", "fix")
	gitCmd(t, tmpDir, "add", "fix.txt")
	gitCmd(t, tmpDir, "commit", "-m", "fix: patch bug")
	gitCmd(t, tmpDir, "tag", "v0.2.0")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	stdout, stderr := captureOutput(t, func() error {
		return runChangelog(&cobra.Command{}, nil)
	})

	if stderr != "" {
		t.Fatalf("expected no stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "add feature") || !strings.Contains(stdout, "patch bug") {
		t.Fatalf("stdout missing changelog entries: %q", stdout)
	}
	if strings.Contains(stdout, "## v0.2.0") || strings.Contains(stdout, "## Changes") || strings.Contains(stdout, "### [v") {
		t.Fatalf("stdout still contains version headers: %q", stdout)
	}
}

func captureOutput(t *testing.T, fn func() error) (string, string) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}

	os.Stdout = stdoutW
	os.Stderr = stderrW

	var stdoutBuf, stderrBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&stdoutBuf, stdoutR)
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&stderrBuf, stderrR)
	}()

	err = fn()

	_ = stdoutW.Close()
	_ = stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	wg.Wait()

	if err != nil {
		t.Fatalf("fn: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String()
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	gitCmd(t, dir, "init")
	gitCmd(t, dir, "config", "user.name", "Test User")
	gitCmd(t, dir, "config", "user.email", "test@example.com")
}

func gitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, string(out))
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
