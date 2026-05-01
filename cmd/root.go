package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "shikai",
	Short: "Release management tool with Conventional Commits",
	Long:  `A CLI tool for managing releases. Analyzes commits since the last tag, suggests a semver bump, generates a changelog, and creates an annotated tag.`,
	Args:  cobra.NoArgs,
	RunE:  runRelease,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().Bool("dry-run", false, "Show what would happen without making changes")
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().BoolP("push", "p", false, "Automatically push tag to remote (skip confirmation)")
	viper.BindPFlag("push", rootCmd.PersistentFlags().Lookup("push"))

	rootCmd.PersistentFlags().String("changelog-path", "", "Custom path for the changelog file (default: CHANGELOG.md in repo root)")
	viper.BindPFlag("changelog-path", rootCmd.PersistentFlags().Lookup("changelog-path"))

	rootCmd.PersistentFlags().String("template", "", "Custom git-chglog template path")
	viper.BindPFlag("template", rootCmd.PersistentFlags().Lookup("template"))

	rootCmd.PersistentFlags().Bool("prerelease", false, "Create a prerelease tag")
	viper.BindPFlag("prerelease", rootCmd.PersistentFlags().Lookup("prerelease"))

	rootCmd.PersistentFlags().String("prerelease-id", "alpha", "Prerelease identifier (e.g., alpha, beta, rc)")
	viper.BindPFlag("prerelease-id", rootCmd.PersistentFlags().Lookup("prerelease-id"))

	rootCmd.Flags().Bool("patch", false, "Bump patch version (x.y.Z)")
	rootCmd.Flags().Bool("minor", false, "Bump minor version (x.Y.z)")
	rootCmd.Flags().Bool("major", false, "Bump major version (X.y.z)")
}
