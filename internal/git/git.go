package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetLatestTag returns the most recent git tag as-is.
func GetLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// No tags found
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// GetPreviousTag returns the most recent tag before the given tag.
func GetPreviousTag(tag string) (string, error) {
	if tag == "" {
		return "", nil
	}

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", tag+"^")
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCommitsSinceTag returns commit messages since the given tag (not including the tag itself).
func GetCommitsSinceTag(tag string) ([]string, error) {
	var args []string
	if tag == "" {
		args = []string{"log", "--all", "--oneline", "--format=%s"}
	} else {
		args = []string{"log", tag + "..HEAD", "--oneline", "--format=%s"}
	}
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

// StageFiles adds files to the git staging area.
func StageFiles(paths ...string) error {
	args := []string{"add"}
	for _, p := range paths {
		if p != "" {
			args = append(args, p)
		}
	}
	if len(args) == 1 {
		return nil
	}
	cmd := exec.Command("git", args...)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	return nil
}

// CommitChanges creates a git commit from the staged changes.
func CommitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	return nil
}

// CreateAnnotatedTag creates an annotated git tag with the given message.
func CreateAnnotatedTag(tag string, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git tag: %w", err)
	}
	return nil
}

// PushTag pushes the tag to the remote origin.
func PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

// PushRelease pushes the current branch commit and the release tag to origin.
func PushRelease(tag string) error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "push", "origin", branch, tag)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

// GetCurrentBranch returns the checked out branch name.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch == "HEAD" || branch == "" {
		return "", fmt.Errorf("release push requires a checked-out branch")
	}

	return branch, nil
}
