package validate

import (
	"fmt"

	"sesh/pkg/models"
)

type Issue struct {
	Severity string
	Message  string
}

const (
	Error   = "error"
	Warning = "warning"
)

func Validate(config *models.Config) []Issue {
	var issues []Issue

	sessionNames := make(map[string]int)

	if len(config.Workspaces) == 0 && len(config.Sessions) == 0 {
		issues = append(issues, Issue{Warning, "config has no workspaces or sessions defined"})
	}

	for i, w := range config.Workspaces {
		if w.Name == "" {
			issues = append(issues, Issue{Error, fmt.Sprintf("workspace at index %d has no name", i)})
		}

		for j, dup := range config.Workspaces {
			if i != j && w.Name == dup.Name && w.Name != "" {
				issues = append(issues, Issue{Error, fmt.Sprintf("duplicate workspace name '%s'", w.Name)})
				break
			}
		}

		for j, s := range w.Sessions {
			if s.Name == "" {
				issues = append(issues, Issue{Error, fmt.Sprintf("workspace '%s' has session at index %d with no name", w.Name, j)})
			}

			sessionNames[s.Name]++

			issues = append(issues, validateSession(s, fmt.Sprintf("workspace '%s'", w.Name))...)
		}
	}

	for i, s := range config.Sessions {
		if s.Name == "" {
			issues = append(issues, Issue{Error, fmt.Sprintf("standalone session at index %d has no name", i)})
		}

		sessionNames[s.Name]++

		issues = append(issues, validateSession(s, "standalone")...)
	}

	for name, count := range sessionNames {
		if count > 1 {
			issues = append(issues, Issue{Error, fmt.Sprintf("duplicate session name '%s' (appears %d times)", name, count)})
		}
	}

	return issues
}

func FilterBySeverity(issues []Issue, severity string) []Issue {
	var result []Issue
	for _, issue := range issues {
		if issue.Severity == severity {
			result = append(result, issue)
		}
	}
	return result
}

func validateSession(s models.Session, context string) []Issue {
	var issues []Issue

	if len(s.Windows) == 0 {
		issues = append(issues, Issue{Warning, fmt.Sprintf("session '%s' (%s) has no windows defined", s.Name, context)})
	}

	windowNames := make(map[string]int)
	for _, w := range s.Windows {
		if w.Name == "" {
			issues = append(issues, Issue{Error, fmt.Sprintf("session '%s' has a window with no name", s.Name)})
		}

		windowNames[w.Name]++

		if w.Layout != "" {
			validLayouts := map[string]bool{
				"even-horizontal": true,
				"even-vertical":   true,
				"main-horizontal": true,
				"main-vertical":   true,
				"tiled":           true,
			}
			if !validLayouts[w.Layout] {
				issues = append(issues, Issue{Warning, fmt.Sprintf("session '%s', window '%s': unknown layout '%s'", s.Name, w.Name, w.Layout)})
			}
		}
	}

	for name, count := range windowNames {
		if count > 1 {
			issues = append(issues, Issue{Warning, fmt.Sprintf("session '%s' has duplicate window name '%s'", s.Name, name)})
		}
	}

	return issues
}
