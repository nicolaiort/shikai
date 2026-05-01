package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetLatestTag returns the most recent git tag, stripping the 'v' prefix.
func GetLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// No tags found
		return "", nil
	}
	tag := strings.TrimSpace(string(output))
	// Strip 'v' prefix
	if strings.HasPrefix(tag, "v") {
		tag = tag[1:]
	}
	return tag, nil
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

// CreateAnnotatedTag creates an annotated git tag with the given message.
func CreateAnnotatedTag(tag string, message string) error {
	cmd := exec.Command("git", "tag", "-a", "v"+tag, "-m", message)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git tag: %w", err)
	}
	return nil
}

// PushTag pushes the tag to the remote origin.
func PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", "v"+tag)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}
