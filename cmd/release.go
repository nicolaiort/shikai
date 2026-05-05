package cmd

import (
	"fmt"
	"strings"

	"github.com/nicolaiort/shikai/internal/changelog"
	"github.com/nicolaiort/shikai/internal/commits"
	"github.com/nicolaiort/shikai/internal/git"
	"github.com/nicolaiort/shikai/internal/interactive"
	"github.com/nicolaiort/shikai/internal/manifest"
	"github.com/nicolaiort/shikai/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runRelease(cmd *cobra.Command, args []string) error {
	dryRun := viper.GetBool("dry-run")
	hooks := loadReleaseHooks(viper.GetViper())
	tagPrefix := viper.GetString("tag-prefix")

	// 1. Find current version from latest tag
	currentVersion, err := git.GetLatestTag()
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}
	currentVersion, commitBaseRef := normalizeReleaseRefs(currentVersion, tagPrefix)
	if commitBaseRef == "" {
		fmt.Printf("No existing tags found. Starting from %s0.0.0\n", tagPrefix)
	}

	// 2. Get commits since last tag
	commitList, err := commits.ParseConventionalCommits(commitBaseRef)
	if err != nil {
		return fmt.Errorf("failed to parse commits: %w", err)
	}

	// 3. Analyze commits and recommend version bump
	bumpType, breakingChanges := commits.AnalyzeBumpType(commitList)

	// 4. Interactive version selection
	selectedBump, err := interactive.SelectVersion(cmd, bumpType)
	if err != nil {
		return fmt.Errorf("interactive selection: %w", err)
	}
	versionSelectedByFlag := releaseVersionSelectedByFlag(cmd)

	newVersion := version.Bump(currentVersion, selectedBump)
	releaseTag := formatReleaseTag(tagPrefix, newVersion)

	// Add prerelease suffix if needed
	if viper.GetBool("prerelease") {
		prereleaseID := viper.GetString("prerelease-id")
		newVersion = version.AddPrerelease(newVersion, prereleaseID)
		releaseTag = formatReleaseTag(tagPrefix, newVersion)
	}

	if dryRun {
		previewReleaseHooks(releaseTag, hooks)
		fmt.Println("\n[DRY RUN] Skipping actual release creation")
		return nil
	}

	if err := runHookPhase(releaseTag, "before anything happens", hooks.before, false); err != nil {
		return err
	}

	fmt.Printf("Releasing %s → %s\n", currentVersion, newVersion)
	if breakingChanges > 0 {
		fmt.Printf("⚠️  Contains %d breaking change(s)\n", breakingChanges)
	}

	// 5. Generate changelog
	changelogPath := viper.GetString("changelog-path")
	templatePath := viper.GetString("template")

	changelogContent, err := buildReleaseChangelog(releaseTag, tagPrefix, templatePath, commitList)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	// 6. Update CHANGELOG.md
	if err := changelog.UpdateFile(changelogPath, changelogContent); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}
	if err := runHookPhase(releaseTag, "after changelog generation", hooks.afterChangelog, false); err != nil {
		return err
	}
	effectiveChangelogPath := changelogPath
	if effectiveChangelogPath == "" {
		effectiveChangelogPath = "CHANGELOG.md"
	}

	// 7. Update manifest files
	if err := manifest.UpdateAll(newVersion); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// 8. Stage changes
	files := []string{effectiveChangelogPath}
	files = append(files, manifest.GetManifestFiles()...)
	if err := git.StageFiles(files...); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	// 9. Commit release changes
	if err := git.CommitChanges(fmt.Sprintf("chore(release): prepare %s", releaseTag)); err != nil {
		return fmt.Errorf("failed to commit release changes: %w", err)
	}

	// 10. Create annotated tag
	if err := git.CreateAnnotatedTag(releaseTag, changelogContent); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	if err := runHookPhase(releaseTag, "after tag creation", hooks.afterTag, false); err != nil {
		return err
	}

	fmt.Printf("\n✅ Created tag: %s\n", releaseTag)
	finishRelease := func() error {
		if err := runHookPhase(releaseTag, "after everything is done", hooks.afterDone, false); err != nil {
			return err
		}
		return nil
	}

	// 11. Push or prompt for push
	pushEnabled := viper.GetBool("push")
	if shouldSkipPush(versionSelectedByFlag, pushEnabled) {
		return finishRelease()
	}
	if !pushEnabled {
		fmt.Print("\nPush tag to remote? [y/N] ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			if branch, branchErr := git.GetCurrentBranch(); branchErr == nil {
				fmt.Printf("Push manually with: git push origin %s %s\n", branch, releaseTag)
			} else {
				fmt.Println("Push manually with: git push origin " + releaseTag)
			}
			return finishRelease()
		}
	}

	if err := git.PushRelease(releaseTag); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	fmt.Println("🚀 Pushed to remote")
	return finishRelease()
}

func normalizeReleaseRefs(latestTag string, tagPrefix string) (string, string) {
	if latestTag == "" {
		return "0.0.0", ""
	}
	version := latestTag
	if tagPrefix != "" && strings.HasPrefix(latestTag, tagPrefix) {
		version = strings.TrimPrefix(latestTag, tagPrefix)
	} else if tagPrefix == "" && strings.HasPrefix(latestTag, "v") {
		version = strings.TrimPrefix(latestTag, "v")
	}
	return version, latestTag
}

func effectiveChangelogPath(path string) string {
	if path == "" {
		return "CHANGELOG.md"
	}
	return path
}

func releaseVersionSelectedByFlag(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("patch") || cmd.Flags().Changed("minor") || cmd.Flags().Changed("major")
}

func shouldSkipPush(versionSelectedByFlag, push bool) bool {
	return versionSelectedByFlag && !push
}

func buildReleaseChangelog(tag string, tagPrefix string, templatePath string, commitList []commits.Commit) (string, error) {
	return changelog.Generate(tag, tagPrefix, templatePath, commitList)
}

func formatReleaseTag(tagPrefix string, version string) string {
	return tagPrefix + version
}
