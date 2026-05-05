package cmd

import (
	"fmt"

	"github.com/nicolaiort/shikai/internal/changelog"
	"github.com/nicolaiort/shikai/internal/commits"
	"github.com/nicolaiort/shikai/internal/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newChangelogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Print the current release changelog",
		Long:  "Generate the current release notes by default, or the full changelog with --full.",
		Args:  cobra.NoArgs,
		RunE:  runChangelog,
	}
	cmd.Flags().Bool("full", false, "Print the full changelog for all versions")
	_ = viper.BindPFlag("full", cmd.Flags().Lookup("full"))
	return cmd
}

func runChangelog(cmd *cobra.Command, args []string) error {
	tagPrefix := viper.GetString("tag-prefix")
	full := viper.GetBool("full")
	latestTag, err := git.GetLatestTag()
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}
	if latestTag == "" {
		return fmt.Errorf("no release tags found")
	}

	var changelogContent string
	if full {
		changelogContent, err = changelog.GenerateFull(tagPrefix, viper.GetString("template"))
	} else {
		previousTag, prevErr := git.GetPreviousTag(latestTag)
		if prevErr != nil {
			return fmt.Errorf("failed to get previous tag: %w", prevErr)
		}

		commitList, parseErr := commits.ParseConventionalCommits(previousTag)
		if parseErr != nil {
			return fmt.Errorf("failed to parse commits: %w", parseErr)
		}

		changelogContent, err = changelog.GenerateReleaseNotes(latestTag, tagPrefix, viper.GetString("template"), commitList)
	}
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	fmt.Print(changelogContent)
	return nil
}
