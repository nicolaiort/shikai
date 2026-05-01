package interactive

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/shikai/release/internal/commits"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// SelectVersion prompts the user to select a version bump or uses flag.
func SelectVersion(cmd *cobra.Command, recommended commits.BumpType) (string, error) {
	if cmd.Flags().Changed("major") {
		return "major", nil
	}
	if cmd.Flags().Changed("minor") {
		return "minor", nil
	}
	if cmd.Flags().Changed("patch") {
		return "patch", nil
	}

	choices := versionChoiceOptions()
	defaultChoice := versionChoiceDefault(recommended)
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return defaultChoice, nil
	}

	var choice string
	prompt := &survey.Select{
		Message: "Select version bump",
		Options: choices,
		Default: defaultChoice,
	}
	if err := survey.AskOne(prompt, &choice); err != nil {
		return "", fmt.Errorf("select version bump: %w", err)
	}
	return choice, nil
}

func versionChoiceOptions() []string {
	return []string{"major", "minor", "patch"}
}

func versionChoiceDefault(recommended commits.BumpType) string {
	if recommended == "" {
		return "patch"
	}
	return string(recommended)
}
