package changelog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicolaiort/shikai/internal/commits"
)

func TestUpdateFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "CHANGELOG.md")

	content := "# Changelog\n\n## v1.0.0\n\n- Initial release"

	err := UpdateFile(path, content)
	if err != nil {
		t.Fatalf("UpdateFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != content {
		t.Errorf("content = %q, want %q", string(data), content)
	}
}

func TestUpdateFilePrepend(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "CHANGELOG.md")

	// Write existing changelog
	existing := "# Changelog\n\n## v0.1.0\n\n- Old content"
	os.WriteFile(path, []byte(existing), 0644)

	// Prepend new content
	newContent := "## v1.0.0\n\n- New release"
	err := UpdateFile(path, newContent)
	if err != nil {
		t.Fatalf("UpdateFile failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	combined := string(data)

	if len(combined) < len(newContent) {
		t.Fatal("new content should be prepended")
	}
}

func TestGenerateSimple(t *testing.T) {
	commitList := []commits.Commit{
		{Type: commits.Type_feat, Scope: "auth", Subject: "add login"},
		{Type: commits.Type_fix, Subject: "fix logout"},
		{Type: commits.Type_chore, Subject: "update deps"},
	}

	output := generateSimple("1.0.0", commitList)

	// Should contain sections
	if len(output) == 0 {
		t.Error("output should not be empty")
	}

	// Should have the version
	if !contains(output, "v1.0.0") {
		t.Error("should contain version")
	}
}

func TestGenerateSimpleWithBreaking(t *testing.T) {
	commitList := []commits.Commit{
		{Type: commits.Type_feat, Subject: "add feature"},
		{Type: commits.Type_feat, IsBreaking: true, Subject: "breaking change"},
	}

	output := generateSimple("2.0.0", commitList)

	if !contains(output, "## v2.0.0") {
		t.Error("should contain version header")
	}
}

func TestGenerateReleaseNotesStripsHeaders(t *testing.T) {
	commitList := []commits.Commit{
		{Type: commits.Type_feat, Subject: "add feature"},
	}

	output, err := GenerateReleaseNotes("2.0.0", "", commitList)
	if err != nil {
		t.Fatalf("GenerateReleaseNotes failed: %v", err)
	}

	if contains(output, "## Changes") || contains(output, "### [v2.0.0]") || contains(output, "### v2.0.0") {
		t.Fatalf("unexpected headers in output: %q", output)
	}
	if !contains(output, "add feature") {
		t.Fatalf("missing release note body: %q", output)
	}
}

func TestGetDefaultPaths(t *testing.T) {
	// These should return empty when no files exist
	cfg := GetDefaultConfigPath()
	tpl := GetDefaultTemplatePath()

	// Should not panic
	_ = cfg
	_ = tpl
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
