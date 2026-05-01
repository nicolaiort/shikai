package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Supported manifest types
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
	manifestFiles = detected

	if len(detected) > 1 {
		return fmt.Errorf("multiple manifest files detected: %v\nPlease configure which one to update", detected)
	}
	if len(detected) == 0 {
		return nil // No manifest to update
	}

	path := detected[0]
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		return updatePackageJSON(path, version)
	case ".toml":
		return updateTOML(path, version)
	case ".yaml", ".yml":
		return updateYAML(path, version)
	case ".gradle":
		return updateGradle(path, version)
	case ".xml":
		return updateXML(path, version)
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

func updateTOML(path string, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// Simple line-based replacement for [package] section
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

func updateYAML(path string, version string) error {
	return errors.New("YAML manifest update not yet implemented")
}

func updateGradle(path string, version string) error {
	return errors.New("Gradle manifest update not yet implemented")
}

func updateXML(path string, version string) error {
	return errors.New("XML manifest update not yet implemented")
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