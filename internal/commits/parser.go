package commits

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shikai/release/internal/git"
)

// CommitType represents conventional commit types.
type CommitType string

const (
	Type_feat     CommitType = "feat"
	Type_fix      CommitType = "fix"
	Type_docs     CommitType = "docs"
	Type_style    CommitType = "style"
	Type_refactor CommitType = "refactor"
	Type_test     CommitType = "test"
	Type_chore    CommitType = "chore"
	Type_revert   CommitType = "revert"
	Type_unknown  CommitType = ""
)

// Commit represents a parsed conventional commit.
type Commit struct {
	Type       CommitType
	Scope      string
	Subject    string
	IsBreaking bool
	Raw        string
}

// ParseConventionalCommits parses commits since the last tag.
// Skips commits that don't follow conventional commit format.
func ParseConventionalCommits(tag string) ([]Commit, error) {
	rawCommits, err := git.GetCommitsSinceTag(tag)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	pattern := regexp.MustCompile(`^(feat|fix|docs|style|refactor|test|chore|revert)(?:\(([^)]*)\))?(!)?:\s*(.+)`)

	for _, msg := range rawCommits {
		parts := pattern.FindStringSubmatch(msg)
		if parts == nil {
			fmt.Fprintf(os.Stderr, "⚠️  Skipping non-conforming commit: %s\n", msg)
			continue
		}
		commits = append(commits, Commit{
			Type:       CommitType(parts[1]),
			Scope:      parts[2],
			Subject:    parts[4],
			IsBreaking: parts[3] == "!" || strings.Contains(parts[4], "BREAKING CHANGE"),
			Raw:        msg,
		})
	}

	return commits, nil
}

// BumpType represents version bump type.
type BumpType string

const (
	BumpMajor BumpType = "major"
	BumpMinor BumpType = "minor"
	BumpPatch BumpType = "patch"
)

// AnalyzeBumpType determines the version bump type based on commits.
func AnalyzeBumpType(commits []Commit) (BumpType, int) {
	hasBreaking := false
	hasFeat := false
	hasPatchLevel := false

	for _, c := range commits {
		if c.IsBreaking {
			hasBreaking = true
		}
		if c.Type == Type_feat {
			hasFeat = true
			continue
		}
		if c.Type == Type_fix || c.Type == Type_docs || c.Type == Type_style || c.Type == Type_refactor || c.Type == Type_test || c.Type == Type_chore || c.Type == Type_revert {
			hasPatchLevel = true
		}
	}

	if hasBreaking {
		return BumpMajor, countBreaking(commits)
	}
	if hasFeat {
		return BumpMinor, 0
	}
	if hasPatchLevel {
		return BumpPatch, 0
	}
	return BumpPatch, 0
}

func countBreaking(commits []Commit) int {
	count := 0
	for _, c := range commits {
		if c.IsBreaking {
			count++
		}
	}
	return count
}
