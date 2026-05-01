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
		{name: "tag present", latestTag: "1.2.3", wantVersion: "1.2.3", wantRef: "1.2.3"},
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
