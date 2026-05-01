package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadReleaseHooksReadsConfig(t *testing.T) {
	v := viper.New()
	v.Set("hooks.before", []string{"echo before"})
	v.Set("hooks.after-changelog", []string{"echo changelog"})
	v.Set("hooks.after-tag", []string{"echo tag"})
	v.Set("hooks.after-done", []string{"echo done"})

	hooks := loadReleaseHooks(v)
	if got := hooks.before; len(got) != 1 || got[0] != "echo before" {
		t.Fatalf("before hooks = %#v", got)
	}
	if got := hooks.afterChangelog; len(got) != 1 || got[0] != "echo changelog" {
		t.Fatalf("after changelog hooks = %#v", got)
	}
	if got := hooks.afterTag; len(got) != 1 || got[0] != "echo tag" {
		t.Fatalf("after tag hooks = %#v", got)
	}
	if got := hooks.afterDone; len(got) != 1 || got[0] != "echo done" {
		t.Fatalf("after done hooks = %#v", got)
	}
}

func TestPreviewReleaseHooksPrintsCommands(t *testing.T) {
	hooks := releaseHooks{
		before:         []string{"echo before"},
		afterChangelog: []string{"echo changelog"},
		afterTag:       []string{"echo tag"},
		afterDone:      []string{"echo done"},
	}

	stdout := captureStdout(t, func() {
		previewReleaseHooks("v1.2.3", hooks)
	})

	for _, want := range []string{
		"[DRY RUN] Would run hook before anything happens (SHIKAI_TAG=v1.2.3): echo before",
		"[DRY RUN] Would run hook after changelog generation (SHIKAI_TAG=v1.2.3): echo changelog",
		"[DRY RUN] Would run hook after tag creation (SHIKAI_TAG=v1.2.3): echo tag",
		"[DRY RUN] Would run hook after everything is done (SHIKAI_TAG=v1.2.3): echo done",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}

func TestRunHookPhaseRunsCommandsInOrder(t *testing.T) {
	oldRunner := executeHookCommand
	defer func() { executeHookCommand = oldRunner }()

	var ran []string
	var seenTag string
	executeHookCommand = func(command string, tag string) error {
		ran = append(ran, command)
		seenTag = tag
		return nil
	}

	if err := runHookPhase("v1.2.3", "before anything happens", []string{"one", "two"}, false); err != nil {
		t.Fatalf("runHookPhase: %v", err)
	}
	if strings.Join(ran, ",") != "one,two" {
		t.Fatalf("commands ran in wrong order: %#v", ran)
	}
	if seenTag != "v1.2.3" {
		t.Fatalf("seen tag = %q, want %q", seenTag, "v1.2.3")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	<-done
	return buf.String()
}
