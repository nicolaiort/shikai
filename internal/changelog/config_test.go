package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildGeneratorConfigUsesTemplateOverride(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "release.tpl.md")
	if err := os.WriteFile(templatePath, []byte("template"), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	cfg, err := buildGeneratorConfig("1.2.3", "v", templatePath)
	if err != nil {
		t.Fatalf("buildGeneratorConfig: %v", err)
	}

	if cfg.Template != templatePath {
		t.Fatalf("template = %q, want %q", cfg.Template, templatePath)
	}
	if cfg.Options.NextTag != "v1.2.3" {
		t.Fatalf("next tag = %q, want %q", cfg.Options.NextTag, "v1.2.3")
	}
	if cfg.Options.TagFilterPattern != "^v" {
		t.Fatalf("tag filter pattern = %q, want %q", cfg.Options.TagFilterPattern, "^v")
	}
	if cfg.Options.HeaderPattern == "" {
		t.Fatal("expected default header pattern")
	}
}

func TestBuildGeneratorConfigLoadsProjectConfig(t *testing.T) {
	tmpDir := t.TempDir()
	chglogDir := filepath.Join(tmpDir, ".chglog")
	if err := os.MkdirAll(chglogDir, 0755); err != nil {
		t.Fatalf("mkdir .chglog: %v", err)
	}

	templatePath := filepath.Join(chglogDir, "custom.tpl.md")
	if err := os.WriteFile(templatePath, []byte("template"), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	configPath := filepath.Join(chglogDir, "config.yml")
	configYAML := []byte("template: custom.tpl.md\ninfo:\n  repository_url: https://example.com\noptions:\n  header:\n    pattern: '^custom$'\n    pattern_maps:\n      - Type\n      - Scope\n      - Subject\n")
	if err := os.WriteFile(configPath, configYAML, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg, err := buildGeneratorConfig("2.0.0", "v", "")
	if err != nil {
		t.Fatalf("buildGeneratorConfig: %v", err)
	}

	if cfg.Template != filepath.Join(".chglog", "custom.tpl.md") {
		t.Fatalf("template = %q, want %q", cfg.Template, filepath.Join(".chglog", "custom.tpl.md"))
	}
	if cfg.Info.RepositoryURL != "https://example.com" {
		t.Fatalf("repo url = %q, want %q", cfg.Info.RepositoryURL, "https://example.com")
	}
	if cfg.Options.HeaderPattern != "^custom$" {
		t.Fatalf("header pattern = %q, want %q", cfg.Options.HeaderPattern, "^custom$")
	}
}

func TestBuildGeneratorConfigSupportsEmptyPrefix(t *testing.T) {
	cfg, err := buildGeneratorConfig("1.2.3", "", "")
	if err != nil {
		t.Fatalf("buildGeneratorConfig: %v", err)
	}

	if cfg.Options.NextTag != "1.2.3" {
		t.Fatalf("next tag = %q, want %q", cfg.Options.NextTag, "1.2.3")
	}
	if cfg.Options.TagFilterPattern != "^[0-9]" {
		t.Fatalf("tag filter pattern = %q, want %q", cfg.Options.TagFilterPattern, "^[0-9]")
	}
}

func TestNormalizeTagName(t *testing.T) {
	tests := []struct {
		tag    string
		prefix string
		want   string
	}{
		{tag: "1.2.3", prefix: "v", want: "v1.2.3"},
		{tag: "v1.2.3", prefix: "v", want: "v1.2.3"},
		{tag: "1.2.3", prefix: "", want: "1.2.3"},
	}

	for _, tt := range tests {
		if got := normalizeTagName(tt.tag, tt.prefix); got != tt.want {
			t.Fatalf("normalizeTagName(%q, %q) = %q, want %q", tt.tag, tt.prefix, got, tt.want)
		}
	}
}
