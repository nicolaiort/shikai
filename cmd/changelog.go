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
	return &cobra.Command{
		Use:   "changelog",
		Short: "Print the current release changelog",
		Long:  "Generate the changelog for the current release and write it to stdout.",
		Args:  cobra.NoArgs,
		RunE:  runChangelog,
	}
}

func runChangelog(cmd *cobra.Command, args []string) error {
	tagPrefix := viper.GetString("tag-prefix")
	latestTag, err := git.GetLatestTag()
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}
	if latestTag == "" {
		return fmt.Errorf("no release tags found")
	}

	previousTag, err := git.GetPreviousTag(latestTag)
	if err != nil {
		return fmt.Errorf("failed to get previous tag: %w", err)
	}

	commitList, err := commits.ParseConventionalCommits(previousTag)
	if err != nil {
		return fmt.Errorf("failed to parse commits: %w", err)
	}

	changelogContent, err := changelog.GenerateReleaseNotes(latestTag, tagPrefix, viper.GetString("template"), commitList)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	fmt.Print(changelogContent)
	return nil
}
