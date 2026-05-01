package commits

import (
	"regexp"
	"testing"
)

func TestCommitPatternMatchesConventionalCommits(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  CommitType
		wantScope string
		wantBody  string
		wantBreak bool
		wantMatch bool
	}{
		{name: "feat with scope", input: "feat(auth): add login form", wantType: Type_feat, wantScope: "auth", wantBody: "add login form", wantMatch: true},
		{name: "fix without scope", input: "fix: handle null response", wantType: Type_fix, wantBody: "handle null response", wantMatch: true},
		{name: "breaking feat", input: "feat(api)!: change response format", wantType: Type_feat, wantScope: "api", wantBody: "change response format", wantBreak: true, wantMatch: true},
		{name: "docs", input: "docs(readme): update installation", wantType: Type_docs, wantScope: "readme", wantBody: "update installation", wantMatch: true},
		{name: "non-conforming", input: "fixed stuff", wantMatch: false},
	}

	pattern := regexp.MustCompile(`^(feat|fix|docs|style|refactor|test|chore|revert)(?:\(([^)]*)\))?(!)?:\s*(.+)`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pattern.FindStringSubmatch(tt.input)
			if !tt.wantMatch {
				if got != nil {
					t.Fatalf("expected no match, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected match, got nil")
			}
			if CommitType(got[1]) != tt.wantType {
				t.Fatalf("type = %q, want %q", got[1], tt.wantType)
			}
			if got[2] != tt.wantScope {
				t.Fatalf("scope = %q, want %q", got[2], tt.wantScope)
			}
			if got[4] != tt.wantBody {
				t.Fatalf("body = %q, want %q", got[4], tt.wantBody)
			}
			if (got[3] == "!") != tt.wantBreak {
				t.Fatalf("breaking = %v, want %v", got[3] == "!", tt.wantBreak)
			}
		})
	}
}

func TestAnalyzeBumpType(t *testing.T) {
	tests := []struct {
		name       string
		commits    []Commit
		wantBump   BumpType
		wantBreaks int
	}{
		{name: "feat wins over fix", commits: []Commit{{Type: Type_feat}, {Type: Type_fix}}, wantBump: BumpMinor},
		{name: "fix only is patch", commits: []Commit{{Type: Type_fix}}, wantBump: BumpPatch},
		{name: "breaking is major", commits: []Commit{{Type: Type_refactor, IsBreaking: true}}, wantBump: BumpMajor, wantBreaks: 1},
		{name: "docs only is patch", commits: []Commit{{Type: Type_docs}}, wantBump: BumpPatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBump, gotBreaks := AnalyzeBumpType(tt.commits)
			if gotBump != tt.wantBump {
				t.Fatalf("bump = %q, want %q", gotBump, tt.wantBump)
			}
			if gotBreaks != tt.wantBreaks {
				t.Fatalf("breaks = %d, want %d", gotBreaks, tt.wantBreaks)
			}
		})
	}
}
