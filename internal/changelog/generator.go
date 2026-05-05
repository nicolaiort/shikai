package changelog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	chglog "github.com/git-chglog/git-chglog"
	"github.com/nicolaiort/shikai/internal/commits"
)

const bundledTemplatePath = "templates/release-changelog.tpl.md"

// Generate creates a changelog using the git-chglog Go library.
func Generate(tag string, tagPrefix string, templatePath string, commitList []commits.Commit) (string, error) {
	tagName := normalizeTagName(tag, tagPrefix)

	cfg, err := buildGeneratorConfig(tagPrefix, templatePath, tagName)
	if err != nil {
		return generateSimple(tagName, commitList), nil
	}

	logger := chglog.NewLogger(io.Discard, io.Discard, true, true)
	gen := chglog.NewGenerator(logger, cfg)

	var buf bytes.Buffer
	if err := gen.Generate(&buf, ""); err != nil {
		return generateSimple(tagName, commitList), nil
	}

	return buf.String(), nil
}

// GenerateFull creates a changelog with all tagged versions.
func GenerateFull(tagPrefix string, templatePath string) (string, error) {
	cfg, err := buildGeneratorConfig(tagPrefix, templatePath, "")
	if err != nil {
		return "", err
	}

	logger := chglog.NewLogger(io.Discard, io.Discard, true, true)
	gen := chglog.NewGenerator(logger, cfg)

	var buf bytes.Buffer
	if err := gen.Generate(&buf, ""); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateReleaseNotes creates release notes without the version history or version header.
func GenerateReleaseNotes(tag string, tagPrefix string, templatePath string, commitList []commits.Commit) (string, error) {
	content, err := Generate(tag, tagPrefix, templatePath, commitList)
	if err != nil {
		return "", err
	}
	return extractReleaseNotesBody(content), nil
}

// UpdateFile writes changelog content to a file.
func UpdateFile(path string, content string) error {
	if path == "" {
		path = "CHANGELOG.md"
	}

	if data, err := os.ReadFile(path); err == nil {
		content = content + "\n\n" + string(data)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func buildGeneratorConfig(tagPrefix string, templateOverride string, nextTag string) (*chglog.Config, error) {
	projectCfgPath := GetDefaultConfigPath()
	projectCfg, err := loadProjectConfig(projectCfgPath)
	if err != nil {
		return nil, err
	}

	templatePath := templateOverride
	if templatePath == "" {
		templatePath = projectCfg.Template
		if templatePath != "" && !filepath.IsAbs(templatePath) && projectCfgPath != "" {
			templatePath = filepath.Join(filepath.Dir(projectCfgPath), templatePath)
		}
	}
	if templatePath == "" {
		templatePath = bundledTemplatePath
	}

	tagFilterPattern := projectCfg.Options.TagFilterPattern
	if tagFilterPattern == "" {
		if tagPrefix == "" {
			tagFilterPattern = "^[0-9]"
		} else {
			tagFilterPattern = "^" + regexp.QuoteMeta(tagPrefix)
		}
	}

	cfg := &chglog.Config{
		Bin:        "git",
		WorkingDir: ".",
		Template:   templatePath,
		Info: &chglog.Info{
			Title:         "CHANGELOG",
			RepositoryURL: projectCfg.Info.RepositoryURL,
		},
		Options: &chglog.Options{
			NextTag:               nextTag,
			TagFilterPattern:      tagFilterPattern,
			Sort:                  "date",
			CommitSortBy:          "Scope",
			CommitGroupBy:         "Type",
			CommitGroupSortBy:     "Title",
			CommitGroupTitleOrder: []string{"feat", "fix", "refactor", "docs", "style", "test", "chore", "revert"},
			CommitGroupTitleMaps: map[string]string{
				"feat":     "Features",
				"fix":      "Bug Fixes",
				"refactor": "Refactoring",
				"docs":     "Documentation",
				"style":    "Styles",
				"test":     "Tests",
				"chore":    "Chores",
				"revert":   "Reverts",
			},
			HeaderPattern: `^(feat|fix|docs|style|refactor|test|chore|revert)(?:\(([\w\$\.\-\*\s]*)\))?(?:!)?:\s(.*)$`,
			HeaderPatternMaps: []string{
				"Type",
				"Scope",
				"Subject",
			},
			NoteKeywords: []string{"BREAKING CHANGE"},
		},
	}

	mergeProjectConfig(cfg, projectCfg)
	return cfg, nil
}

func mergeProjectConfig(cfg *chglog.Config, project projectConfig) {
	if project.Info.Title != "" {
		cfg.Info.Title = project.Info.Title
	}
	if project.Info.RepositoryURL != "" {
		cfg.Info.RepositoryURL = project.Info.RepositoryURL
	}
	if project.Options.TagFilterPattern != "" {
		cfg.Options.TagFilterPattern = project.Options.TagFilterPattern
	}
	if project.Options.Sort != "" {
		cfg.Options.Sort = project.Options.Sort
	}
	if project.Options.Commits.SortBy != "" {
		cfg.Options.CommitSortBy = project.Options.Commits.SortBy
	}
	if len(project.Options.Commits.Filters) > 0 {
		cfg.Options.CommitFilters = project.Options.Commits.Filters
	}
	if project.Options.CommitGroups.GroupBy != "" {
		cfg.Options.CommitGroupBy = project.Options.CommitGroups.GroupBy
	}
	if project.Options.CommitGroups.SortBy != "" {
		cfg.Options.CommitGroupSortBy = project.Options.CommitGroups.SortBy
	}
	if len(project.Options.CommitGroups.TitleOrder) > 0 {
		cfg.Options.CommitGroupTitleOrder = project.Options.CommitGroups.TitleOrder
	}
	if len(project.Options.CommitGroups.TitleMaps) > 0 {
		cfg.Options.CommitGroupTitleMaps = project.Options.CommitGroups.TitleMaps
	}
	if project.Options.Header.Pattern != "" {
		cfg.Options.HeaderPattern = project.Options.Header.Pattern
	}
	if len(project.Options.Header.PatternMaps) > 0 {
		cfg.Options.HeaderPatternMaps = project.Options.Header.PatternMaps
	}
	if len(project.Options.Notes.Keywords) > 0 {
		cfg.Options.NoteKeywords = project.Options.Notes.Keywords
	}
}

func generateSimple(tag string, commitList []commits.Commit) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("## %s\n\n", tag))

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

func extractReleaseNotesBody(content string) string {
	lines := strings.Split(content, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## Changes" {
			start = i + 1
			break
		}
	}

	if start == -1 {
		return stripFirstHeading(content)
	}

	body := lines[start:]
	for len(body) > 0 && strings.TrimSpace(body[0]) == "" {
		body = body[1:]
	}
	if len(body) > 0 && strings.HasPrefix(strings.TrimSpace(body[0]), "<a name=") {
		body = body[1:]
	}
	if len(body) > 0 && strings.HasPrefix(strings.TrimSpace(body[0]), "### ") {
		body = body[1:]
	}
	for len(body) > 0 && strings.TrimSpace(body[0]) == "" {
		body = body[1:]
	}
	if len(body) > 0 && strings.HasPrefix(strings.TrimSpace(body[0]), "> ") {
		body = body[1:]
	}
	for len(body) > 0 && strings.TrimSpace(body[0]) == "" {
		body = body[1:]
	}

	return strings.TrimSpace(strings.Join(body, "\n"))
}

func stripFirstHeading(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}
	if strings.HasPrefix(strings.TrimSpace(lines[0]), "## ") {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func normalizeTagName(tag string, prefix string) string {
	if prefix == "" {
		return tag
	}
	if strings.HasPrefix(tag, prefix) {
		return tag
	}
	return prefix + tag
}

// GetDefaultConfigPath returns the git-chglog config path.
func GetDefaultConfigPath() string {
	cfg := filepath.Join(".chglog", "config.yml")
	if _, err := os.Stat(cfg); err == nil {
		return cfg
	}
	return ""
}

// GetDefaultTemplatePath returns the default bundled template path.
func GetDefaultTemplatePath() string {
	if _, err := os.Stat(bundledTemplatePath); err == nil {
		return bundledTemplatePath
	}
	tpl := filepath.Join(".chglog", "CHANGELOG.tpl.md")
	if _, err := os.Stat(tpl); err == nil {
		return tpl
	}
	return ""
}
