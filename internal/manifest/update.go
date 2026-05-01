package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

var manifestTypes = []string{
	"package.json",
	"Cargo.toml",
	"pyproject.toml",
	"go.mod",
	"pom.xml",
	"build.gradle",
	"pubspec.yaml",
}

var manifestFiles []string

// GetManifestFiles returns the list of detected manifest files.
func GetManifestFiles() []string {
	return manifestFiles
}

// UpdateAll detects and updates all manifest files in the repository.
func UpdateAll(version string) error {
	detected := detectManifests()
	updatable := make([]string, 0, len(detected))
	for _, path := range detected {
		if isUpdatableManifest(path) {
			updatable = append(updatable, path)
			continue
		}
		fmt.Printf("⚠️  Skipping unsupported manifest: %s\n", path)
	}
	manifestFiles = updatable

	if len(updatable) > 1 {
		return fmt.Errorf("multiple manifest files detected: %v\nPlease configure which one to update", updatable)
	}
	if len(updatable) == 0 {
		return nil // No manifest to update
	}

	path := updatable[0]
	switch path {
	case "package.json":
		return updatePackageJSON(path, version)
	case "Cargo.toml":
		return updateCargoToml(path, version)
	default:
		return fmt.Errorf("unsupported manifest type: %s", path)
	}
}

func detectManifests() []string {
	var found []string
	for _, name := range manifestTypes {
		if _, err := os.Stat(name); err == nil {
			found = append(found, name)
		}
	}
	return found
}

func updatePackageJSON(path string, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	pkg["version"] = version
	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

func isUpdatableManifest(path string) bool {
	switch path {
	case "package.json", "Cargo.toml":
		return true
	default:
		return false
	}
}

func updateCargoToml(path string, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	lines := ""
	inPackage := false
	for _, line := range splitLines(content) {
		if line == "[package]" {
			inPackage = true
		} else if inPackage && len(line) > 0 && line[0] == '[' {
			inPackage = false
		}
		if inPackage && len(line) > 8 && line[:8] == "version" {
			line = fmt.Sprintf("version = \"%s\"", version)
		}
		lines += line + "\n"
	}
	return os.WriteFile(path, []byte(lines), 0644)
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	return lines
}
