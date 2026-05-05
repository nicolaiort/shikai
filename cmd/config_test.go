package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigIgnoresMissingFile(t *testing.T) {
	tmpDir := t.TempDir()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	v := viper.New()
	if err := loadConfig(v); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
}

func TestLoadConfigReadsPushSetting(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".shikai.yml")
	if err := os.WriteFile(configPath, []byte("push: true\n"), 0644); err != nil {
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

	v := viper.New()
	if err := loadConfig(v); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if !v.GetBool("push") {
		t.Fatal("expected push to be enabled from config")
	}
}

func TestLoadConfigReadsHooks(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".shikai.yml")
	yaml := `hooks:
  before:
    - echo before
  after-changelog:
    - echo changelog
  after-tag:
    - echo tag
  after-done:
    - echo done
`
	if err := os.WriteFile(configPath, []byte(yaml), 0644); err != nil {
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

	v := viper.New()
	if err := loadConfig(v); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if got := v.GetStringSlice("hooks.before"); len(got) != 1 || got[0] != "echo before" {
		t.Fatalf("hooks.before = %#v", got)
	}
	if got := v.GetStringSlice("hooks.after-changelog"); len(got) != 1 || got[0] != "echo changelog" {
		t.Fatalf("hooks.after-changelog = %#v", got)
	}
	if got := v.GetStringSlice("hooks.after-tag"); len(got) != 1 || got[0] != "echo tag" {
		t.Fatalf("hooks.after-tag = %#v", got)
	}
	if got := v.GetStringSlice("hooks.after-done"); len(got) != 1 || got[0] != "echo done" {
		t.Fatalf("hooks.after-done = %#v", got)
	}
}

func TestLoadConfigReadsTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".shikai.yml")
	if err := os.WriteFile(configPath, []byte("template: templates/release-changelog.tpl.md\n"), 0644); err != nil {
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

	v := viper.New()
	if err := loadConfig(v); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if got := v.GetString("template"); got != "templates/release-changelog.tpl.md" {
		t.Fatalf("template = %q, want %q", got, "templates/release-changelog.tpl.md")
	}
}

func TestLoadConfigReadsTagPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".shikai.yml")
	if err := os.WriteFile(configPath, []byte("tag-prefix: \"\"\n"), 0644); err != nil {
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

	v := viper.New()
	v.SetDefault("tag-prefix", "v")
	if err := loadConfig(v); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if got := v.GetString("tag-prefix"); got != "" {
		t.Fatalf("tag-prefix = %q, want %q", got, "")
	}
}
