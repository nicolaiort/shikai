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
		{name: "tag present with default prefix", latestTag: "v1.2.3", wantVersion: "1.2.3", wantRef: "v1.2.3"},
		{name: "tag present without prefix", latestTag: "1.2.3", wantVersion: "1.2.3", wantRef: "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotRef := normalizeReleaseRefs(tt.latestTag, "v")
			if gotVersion != tt.wantVersion || gotRef != tt.wantRef {
				t.Fatalf("normalizeReleaseRefs(%q, %q) = (%q, %q), want (%q, %q)", tt.latestTag, "v", gotVersion, gotRef, tt.wantVersion, tt.wantRef)
			}
		})
	}
}

func TestNormalizeReleaseRefsWithoutPrefix(t *testing.T) {
	gotVersion, gotRef := normalizeReleaseRefs("1.2.3", "")
	if gotVersion != "1.2.3" || gotRef != "1.2.3" {
		t.Fatalf("normalizeReleaseRefs without prefix = (%q, %q), want (%q, %q)", gotVersion, gotRef, "1.2.3", "1.2.3")
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

func TestShouldSkipPush(t *testing.T) {
	tests := []struct {
		name                  string
		versionSelectedByFlag bool
		push                  bool
		want                  bool
	}{
		{name: "interactive and no push", versionSelectedByFlag: false, push: false, want: false},
		{name: "interactive and push", versionSelectedByFlag: false, push: true, want: false},
		{name: "flag-selected and no push", versionSelectedByFlag: true, push: false, want: true},
		{name: "flag-selected and push", versionSelectedByFlag: true, push: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkipPush(tt.versionSelectedByFlag, tt.push); got != tt.want {
				t.Fatalf("shouldSkipPush(%v, %v) = %v, want %v", tt.versionSelectedByFlag, tt.push, got, tt.want)
			}
		})
	}
}
