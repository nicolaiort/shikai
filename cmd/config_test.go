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
