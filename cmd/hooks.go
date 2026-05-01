package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type releaseHooks struct {
	before         []string
	afterChangelog []string
	afterTag       []string
	afterDone      []string
}

var executeHookCommand = runShellCommand

func loadReleaseHooks(v *viper.Viper) releaseHooks {
	return releaseHooks{
		before:         v.GetStringSlice("hooks.before"),
		afterChangelog: v.GetStringSlice("hooks.after-changelog"),
		afterTag:       v.GetStringSlice("hooks.after-tag"),
		afterDone:      v.GetStringSlice("hooks.after-done"),
	}
}

func previewReleaseHooks(tag string, h releaseHooks) {
	_ = runHookPhase(tag, "before anything happens", h.before, true)
	_ = runHookPhase(tag, "after changelog generation", h.afterChangelog, true)
	_ = runHookPhase(tag, "after tag creation", h.afterTag, true)
	_ = runHookPhase(tag, "after everything is done", h.afterDone, true)
}

func runHookPhase(tag string, phase string, commands []string, dryRun bool) error {
	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		if dryRun {
			fmt.Printf("[DRY RUN] Would run hook %s (SHIKAI_TAG=%s): %s\n", phase, tag, command)
			continue
		}

		if err := executeHookCommand(command, tag); err != nil {
			return fmt.Errorf("%s hook %q failed: %w", phase, command, err)
		}
	}

	return nil
}

func runShellCommand(command string, tag string) error {
	var shell string
	var shellArgs []string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		shellArgs = []string{"/C", command}
	} else {
		shell = "sh"
		shellArgs = []string{"-c", command}
	}

	cmd := exec.Command(shell, shellArgs...)
	cmd.Env = append(os.Environ(), "SHIKAI_TAG="+tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
