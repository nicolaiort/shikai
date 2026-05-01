package cmd

import (
	"fmt"
	"strings"

	"github.com/shikai/release/internal/changelog"
	"github.com/shikai/release/internal/commits"
	"github.com/shikai/release/internal/git"
	"github.com/shikai/release/internal/interactive"
	"github.com/shikai/release/internal/manifest"
	"github.com/shikai/release/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runRelease(cmd *cobra.Command, args []string) error {
	dryRun := viper.GetBool("dry-run")

	// 1. Find current version from latest tag
	currentVersion, err := git.GetLatestTag()
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}
	currentVersion, commitBaseRef := normalizeReleaseRefs(currentVersion)
	if commitBaseRef == "" {
		fmt.Println("No existing tags found. Starting from v0.0.0")
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

	newVersion := version.Bump(currentVersion, selectedBump)

	// Add prerelease suffix if needed
	if viper.GetBool("prerelease") {
		prereleaseID := viper.GetString("prerelease-id")
		newVersion = version.AddPrerelease(newVersion, prereleaseID)
	}

	fmt.Printf("Releasing %s → %s\n", currentVersion, newVersion)
	if breakingChanges > 0 {
		fmt.Printf("⚠️  Contains %d breaking change(s)\n", breakingChanges)
	}
	if dryRun {
		fmt.Println("\n[DRY RUN] Skipping actual release creation")
		return nil
	}

	// 5. Generate changelog
	changelogPath := viper.GetString("changelog-path")
	templatePath := viper.GetString("template")

	changelogContent, err := changelog.Generate(newVersion, templatePath, commitList)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	// 6. Update CHANGELOG.md
	if err := changelog.UpdateFile(changelogPath, changelogContent); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	// 7. Update manifest files
	if err := manifest.UpdateAll(newVersion); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// 8. Stage changes
	files := []string{changelogPath}
	files = append(files, manifest.GetManifestFiles()...)
	if err := git.StageFiles(files...); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	// 9. Commit release changes
	if err := git.CommitChanges(fmt.Sprintf("chore(release): prepare v%s", newVersion)); err != nil {
		return fmt.Errorf("failed to commit release changes: %w", err)
	}

	// 10. Create annotated tag
	if err := git.CreateAnnotatedTag(newVersion, changelogContent); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Printf("\n✅ Created tag: v%s\n", newVersion)

	// 11. Prompt for push
	if !viper.GetBool("push") {
		fmt.Print("\nPush tag to remote? [y/N] ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Aborted. Push manually with: git push origin v" + newVersion)
			return nil
		}
	}

	if err := git.PushTag(newVersion); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	fmt.Println("🚀 Pushed to remote")
	return nil
}

func normalizeReleaseRefs(latestTag string) (string, string) {
	if latestTag == "" {
		return "0.0.0", ""
	}
	version := strings.TrimPrefix(latestTag, "v")
	return version, latestTag
}
