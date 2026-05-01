package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shikai/release/internal/commits"
)

// Generate creates a changelog using git-chglog.
func Generate(tag string, templatePath string, commitList []commits.Commit) (string, error) {
	// Determine template path
	if templatePath == "" {
		templatePath = GetDefaultTemplatePath()
	}

	args := []string{"chglog", "output", "v" + tag}
	if templatePath != "" {
		args = append(args, "-t", templatePath)
	}

	// Check for project config override
	configPath := GetDefaultConfigPath()
	if configPath != "" {
		args = append(args, "-c", configPath)
	}

	cmd := exec.Command("git-chglog", args...)
	output, err := cmd.Output()
	if err != nil {
		// Fallback: generate simple changelog
		return generateSimple(tag, commitList), nil
	}
	return string(output), nil
}

// UpdateFile writes changelog content to a file.
func UpdateFile(path string, content string) error {
	if path == "" {
		path = "CHANGELOG.md"
	}

	// Check if file exists and prepend new content
	if data, err := os.ReadFile(path); err == nil {
		content = content + "\n\n" + string(data)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func generateSimple(tag string, commitList []commits.Commit) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Changelog\n\n## v%s\n\n", tag))

	grouped := make(map[commits.CommitType][]commits.Commit)
	for _, c := range commitList {
		grouped[c.Type] = append(grouped[c.Type], c)
	}

	order := []commits.CommitType{commits.Type_feat, commits.Type_fix, commits.Type_refactor, commits.Type_docs, commits.Type_style, commits.Type_test, commits.Type_chore}
	for _, t := range order {
		if cs, ok := grouped[t]; ok {
			builder.WriteString(fmt.Sprintf("### %s\n\n", t))
			for _, c := range cs {
				builder.WriteString(fmt.Sprintf("- %s\n", c.Subject))
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// GetDefaultConfigPath returns the git-chglog config path.
func GetDefaultConfigPath() string {
	cfg := filepath.Join(".chglog", "config.yml")
	if _, err := os.Stat(cfg); err == nil {
		return cfg
	}
	return ""
}

// GetDefaultTemplatePath returns the git-chglog template path.
func GetDefaultTemplatePath() string {
	tpl := filepath.Join(".chglog", "CHANGELOG.tpl.md")
	if _, err := os.Stat(tpl); err == nil {
		return tpl
	}
	return ""
}
