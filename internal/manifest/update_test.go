package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateAllSkipsUnsupportedGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	writeFile(t, tmpDir, "go.mod", "module example.com/demo\n\ngo 1.26\n")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := UpdateAll("1.2.3"); err != nil {
		t.Fatalf("UpdateAll returned error: %v", err)
	}

	if got := GetManifestFiles(); len(got) != 0 {
		t.Fatalf("GetManifestFiles() = %v, want empty", got)
	}
}

func TestUpdateAllUpdatesPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()
	writeFile(t, tmpDir, "package.json", `{"name":"demo","version":"0.1.0"}`)

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := UpdateAll("1.2.3"); err != nil {
		t.Fatalf("UpdateAll returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "package.json"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "{\n  \"name\": \"demo\",\n  \"version\": \"1.2.3\"\n}" {
		t.Fatalf("package.json = %q", string(data))
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
