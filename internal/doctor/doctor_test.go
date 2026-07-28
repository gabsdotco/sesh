package doctor

import "testing"

func TestDiagnose_NoIssues(t *testing.T) {
	issues, candidates := Diagnose(
		[]string{"session1", "session2"},
		[]string{"session1", "session2"},
	)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
	if len(candidates) != 0 {
		t.Errorf("expected no rename candidates, got %d", len(candidates))
	}
}

func TestDiagnose_NotRunning(t *testing.T) {
	issues, _ := Diagnose(
		[]string{"session1", "session2"},
		[]string{"session1"},
	)
	notRunningCount := 0
	for _, issue := range issues {
		if issue.Type == "not-running" {
			notRunningCount++
		}
	}
	if notRunningCount != 1 {
		t.Errorf("expected 1 not-running issue, got %d", notRunningCount)
	}
}

func TestDiagnose_Untracked(t *testing.T) {
	issues, _ := Diagnose(
		[]string{"session1"},
		[]string{"session1", "unknown-session"},
	)
	untrackedCount := 0
	for _, issue := range issues {
		if issue.Type == "untracked" {
			untrackedCount++
		}
	}
	if untrackedCount != 1 {
		t.Errorf("expected 1 untracked issue, got %d", untrackedCount)
	}
}

func TestDiagnose_PossibleRename(t *testing.T) {
	_, candidates := Diagnose(
		[]string{"work/backend"},
		[]string{"work/backend-api"},
	)
	if len(candidates) != 1 {
		t.Fatalf("expected 1 rename candidate, got %d", len(candidates))
	}
	if candidates[0].OldName != "work/backend" || candidates[0].NewName != "work/backend-api" {
		t.Errorf("unexpected candidate: %v", candidates[0])
	}
}

func TestIsSimilarName_SamePrefix(t *testing.T) {
	if !IsSimilarName("work/api", "work/api-v2") {
		t.Error("expected similar names with same prefix")
	}
}

func TestIsSimilarName_CommonPrefix(t *testing.T) {
	if !IsSimilarName("backend", "backend-api") {
		t.Error("expected similar names with common prefix")
	}
}

func TestIsSimilarName_Different(t *testing.T) {
	if IsSimilarName("frontend", "backend") {
		t.Error("expected different names to not be similar")
	}
}

func TestIsSimilarName_SameName(t *testing.T) {
	if IsSimilarName("frontend", "frontend") {
		t.Error("same name should not be a rename candidate")
	}
}

func TestNamePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"work/api", "work"},
		{"personal/editor", "personal"},
		{"standalone", ""},
		{"a/b/c", "a"},
	}
	for _, tt := range tests {
		result := NamePrefix(tt.input)
		if result != tt.expected {
			t.Errorf("NamePrefix(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDiagnose_Empty(t *testing.T) {
	issues, candidates := Diagnose(nil, nil)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
	if len(candidates) != 0 {
		t.Errorf("expected no candidates, got %d", len(candidates))
	}
}

func TestDiagnose_AllUntracked(t *testing.T) {
	issues, _ := Diagnose(nil, []string{"foo", "bar"})
	untrackedCount := 0
	for _, issue := range issues {
		if issue.Type == "untracked" {
			untrackedCount++
		}
	}
	if untrackedCount != 2 {
		t.Errorf("expected 2 untracked issues, got %d", untrackedCount)
	}
}

func TestDiagnose_AllNotRunning(t *testing.T) {
	issues, _ := Diagnose([]string{"foo", "bar"}, nil)
	notRunningCount := 0
	for _, issue := range issues {
		if issue.Type == "not-running" {
			notRunningCount++
		}
	}
	if notRunningCount != 2 {
		t.Errorf("expected 2 not-running issues, got %d", notRunningCount)
	}
}
