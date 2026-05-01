package interactive

import (
	"fmt"
	"strings"

	"github.com/shikai/release/internal/commits"
	"github.com/spf13/cobra"
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

	fmt.Printf("\n📊 Recommended bump: %s\n", recommended)
	fmt.Println("Select version bump:")
	fmt.Println("  [M]ajor  - Breaking changes")
	fmt.Println("  [m]inor  - New features (backward compatible)")
	fmt.Println("  [p]atch  - Bug fixes")

	var choice string
	for {
		fmt.Print("\nChoice [M/m/p]: ")
		fmt.Scanln(&choice)
		switch strings.ToLower(choice) {
		case "m":
			return "major", nil
		case "":
			return string(recommended), nil
		case "p":
			return "patch", nil
		default:
			fmt.Println("Invalid selection")
		}
	}
}
