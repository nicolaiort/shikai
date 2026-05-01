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

func previewReleaseHooks(h releaseHooks) {
	_ = runHookPhase("before anything happens", h.before, true)
	_ = runHookPhase("after changelog generation", h.afterChangelog, true)
	_ = runHookPhase("after tag creation", h.afterTag, true)
	_ = runHookPhase("after everything is done", h.afterDone, true)
}

func runHookPhase(phase string, commands []string, dryRun bool) error {
	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		if dryRun {
			fmt.Printf("[DRY RUN] Would run hook %s: %s\n", phase, command)
			continue
		}

		if err := executeHookCommand(command); err != nil {
			return fmt.Errorf("%s hook %q failed: %w", phase, command, err)
		}
	}

	return nil
}

func runShellCommand(command string) error {
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
