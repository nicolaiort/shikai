package cmd

import "testing"

func TestNormalizeReleaseRefs(t *testing.T) {
	tests := []struct {
		name        string
		latestTag   string
		wantVersion string
		wantRef     string
	}{
		{name: "no tags", latestTag: "", wantVersion: "0.0.0", wantRef: ""},
		{name: "tag present", latestTag: "v1.2.3", wantVersion: "1.2.3", wantRef: "v1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotRef := normalizeReleaseRefs(tt.latestTag)
			if gotVersion != tt.wantVersion || gotRef != tt.wantRef {
				t.Fatalf("normalizeReleaseRefs(%q) = (%q, %q), want (%q, %q)", tt.latestTag, gotVersion, gotRef, tt.wantVersion, tt.wantRef)
			}
		})
	}
}

func TestEffectiveChangelogPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "default path", path: "", want: "CHANGELOG.md"},
		{name: "custom path", path: "docs/CHANGELOG.md", want: "docs/CHANGELOG.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := effectiveChangelogPath(tt.path); got != tt.want {
				t.Fatalf("effectiveChangelogPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
