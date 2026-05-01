package interactive

import (
	"testing"

	"github.com/nicolaiort/shikai/internal/commits"
)

func TestVersionChoiceDefault(t *testing.T) {
	tests := []struct {
		name        string
		recommended commits.BumpType
		want        string
	}{
		{name: "major", recommended: commits.BumpMajor, want: "major"},
		{name: "minor", recommended: commits.BumpMinor, want: "minor"},
		{name: "patch", recommended: commits.BumpPatch, want: "patch"},
		{name: "empty falls back to patch", recommended: "", want: "patch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionChoiceDefault(tt.recommended); got != tt.want {
				t.Fatalf("versionChoiceDefault(%q) = %q, want %q", tt.recommended, got, tt.want)
			}
		})
	}
}

func TestVersionChoiceOptions(t *testing.T) {
	got := versionChoiceOptions()
	want := []string{"major", "minor", "patch"}
	if len(got) != len(want) {
		t.Fatalf("len(versionChoiceOptions()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("versionChoiceOptions()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
