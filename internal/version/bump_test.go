package version

import (
	"testing"
)

func TestBump(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		bump     string
		expected string
	}{
		{"patch bump", "1.2.3", "patch", "1.2.4"},
		{"minor bump", "1.2.3", "minor", "1.3.0"},
		{"major bump", "1.2.3", "major", "2.0.0"},
		{"zero version", "0.0.0", "patch", "0.0.1"},
		{"starting version", "0.0.0", "major", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bump(tt.current, tt.bump)
			if got != tt.expected {
				t.Errorf("Bump(%q, %q) = %q, want %q", tt.current, tt.bump, got, tt.expected)
			}
		})
	}
}

func TestAddPrerelease(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		id       string
		expected string
	}{
		{"alpha", "1.2.3", "alpha", "1.2.3-alpha.0"},
		{"beta", "2.0.0", "beta", "2.0.0-beta.0"},
		{"rc", "1.0.0", "rc", "1.0.0-rc.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddPrerelease(tt.version, tt.id)
			if got != tt.expected {
				t.Errorf("AddPrerelease(%q, %q) = %q, want %q", tt.version, tt.id, got, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	major, minor, patch := parseVersion("1.2.3")
	if major != 1 || minor != 2 || patch != 3 {
		t.Errorf("parseVersion(1.2.3) = %d.%d.%d, want 1.2.3", major, minor, patch)
	}
}

func TestBumpEdgeCases(t *testing.T) {
	tests := []struct {
		current  string
		bump     string
		expected string
	}{
		{"9.9.9", "patch", "9.9.10"},
		{"9.9.9", "major", "10.0.0"},
	}

	for _, tt := range tests {
		got := Bump(tt.current, tt.bump)
		if got != tt.expected {
			t.Errorf("Bump(%q, %q) = %q, want %q", tt.current, tt.bump, got, tt.expected)
		}
	}
}
