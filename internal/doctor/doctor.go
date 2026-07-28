package doctor

import "strings"

type RenameCandidate struct {
	OldName string
	NewName string
}

type Issue struct {
	Type    string
	Name    string
	Message string
}

func Diagnose(configSessionNames []string, runningSessionNames []string) ([]Issue, []RenameCandidate) {
	configSet := make(map[string]bool)
	for _, name := range configSessionNames {
		configSet[name] = true
	}

	runningSet := make(map[string]bool)
	for _, name := range runningSessionNames {
		runningSet[name] = true
	}

	var issues []Issue
	var notRunning []string
	var untracked []string

	for _, name := range configSessionNames {
		if !runningSet[name] {
			notRunning = append(notRunning, name)
			issues = append(issues, Issue{
				Type:    "not-running",
				Name:    name,
				Message: "config session not running",
			})
		}
	}

	for _, name := range runningSessionNames {
		if !configSet[name] {
			untracked = append(untracked, name)
			issues = append(issues, Issue{
				Type:    "untracked",
				Name:    name,
				Message: "running session not in config",
			})
		}
	}

	candidates := DetectPossibleRenames(notRunning, untracked)

	return issues, candidates
}

func DetectPossibleRenames(notRunning []string, untracked []string) []RenameCandidate {
	var candidates []RenameCandidate
	for _, old := range notRunning {
		for _, new_ := range untracked {
			if IsSimilarName(old, new_) {
				candidates = append(candidates, RenameCandidate{OldName: old, NewName: new_})
			}
		}
	}
	return candidates
}

func IsSimilarName(a, b string) bool {
	prefixA := NamePrefix(a)
	prefixB := NamePrefix(b)
	if prefixA != "" && prefixB != "" && prefixA == prefixB && a != b {
		return true
	}
	if a != b && len(a) >= 3 && len(b) >= 3 && strings.HasPrefix(b, a[:3]) {
		return true
	}
	return false
}

func FilterByType(issues []Issue, issueType string) []Issue {
	var result []Issue
	for _, issue := range issues {
		if issue.Type == issueType {
			result = append(result, issue)
		}
	}
	return result
}

func NamePrefix(name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}
