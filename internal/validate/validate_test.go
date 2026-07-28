package validate

import (
	"testing"

	"sesh/pkg/models"
)

func TestValidate_EmptyConfig(t *testing.T) {
	config := &models.Config{}
	issues := Validate(config)
	if len(issues) == 0 {
		t.Error("expected warning for empty config")
	}
	if issues[0].Severity != Warning {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	config := &models.Config{
		Workspaces: []models.Workspace{
			{
				Name: "work",
				Sessions: []models.Session{
					{
						Name: "backend",
						Windows: []models.Window{
							{Name: "editor", Panels: []models.Panel{{}}},
						},
					},
				},
			},
		},
	}
	issues := Validate(config)
	errorCount := 0
	for _, issue := range issues {
		if issue.Severity == Error {
			errorCount++
		}
	}
	if errorCount > 0 {
		t.Errorf("expected no errors, got %d: %v", errorCount, issues)
	}
}

func TestValidate_DuplicateSessionNames(t *testing.T) {
	config := &models.Config{
		Workspaces: []models.Workspace{
			{
				Name: "work",
				Sessions: []models.Session{
					{Name: "api", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
				},
			},
			{
				Name: "personal",
				Sessions: []models.Session{
					{Name: "api", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
				},
			},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Error && issue.Message == "duplicate session name 'api' (appears 2 times)" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate session name error, got: %v", issues)
	}
}

func TestValidate_DuplicateWorkspaceNames(t *testing.T) {
	config := &models.Config{
		Workspaces: []models.Workspace{
			{Name: "work", Sessions: []models.Session{}},
			{Name: "work", Sessions: []models.Session{}},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Error && issue.Message == "duplicate workspace name 'work'" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate workspace name error, got: %v", issues)
	}
}

func TestValidate_EmptyWorkspaceName(t *testing.T) {
	config := &models.Config{
		Workspaces: []models.Workspace{
			{Name: "", Sessions: []models.Session{}},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Error && issue.Message == "workspace at index 0 has no name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected empty workspace name error, got: %v", issues)
	}
}

func TestValidate_EmptySessionName(t *testing.T) {
	config := &models.Config{
		Sessions: []models.Session{
			{Name: "", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Error && issue.Message == "standalone session at index 0 has no name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected empty session name error, got: %v", issues)
	}
}

func TestValidate_NoWindows(t *testing.T) {
	config := &models.Config{
		Sessions: []models.Session{
			{Name: "test", Windows: []models.Window{}},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Warning && issue.Message == "session 'test' (standalone) has no windows defined" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected no-windows warning, got: %v", issues)
	}
}

func TestValidate_InvalidLayout(t *testing.T) {
	config := &models.Config{
		Sessions: []models.Session{
			{
				Name: "test",
				Windows: []models.Window{
					{Name: "editor", Layout: "invalid-layout", Panels: []models.Panel{{}}},
				},
			},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Warning && issue.Message == "session 'test', window 'editor': unknown layout 'invalid-layout'" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected invalid layout warning, got: %v", issues)
	}
}

func TestValidate_EmptyWindowName(t *testing.T) {
	config := &models.Config{
		Sessions: []models.Session{
			{
				Name: "test",
				Windows: []models.Window{
					{Name: "", Panels: []models.Panel{{}}},
				},
			},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Error && issue.Message == "session 'test' has a window with no name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected empty window name error, got: %v", issues)
	}
}

func TestValidate_DuplicateWindowNames(t *testing.T) {
	config := &models.Config{
		Sessions: []models.Session{
			{
				Name: "test",
				Windows: []models.Window{
					{Name: "editor", Panels: []models.Panel{{}}},
					{Name: "editor", Panels: []models.Panel{{}}},
				},
			},
		},
	}
	issues := Validate(config)
	found := false
	for _, issue := range issues {
		if issue.Severity == Warning && issue.Message == "session 'test' has duplicate window name 'editor'" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate window name warning, got: %v", issues)
	}
}

func TestValidate_ValidLayouts(t *testing.T) {
	layouts := []string{"even-horizontal", "even-vertical", "main-horizontal", "main-vertical", "tiled"}
	for _, layout := range layouts {
		config := &models.Config{
			Sessions: []models.Session{
				{
					Name: "test",
					Windows: []models.Window{
						{Name: "editor", Layout: layout, Panels: []models.Panel{{}}},
					},
				},
			},
		}
		issues := Validate(config)
		for _, issue := range issues {
			if issue.Severity == Warning && issue.Message != "config has no workspaces or sessions defined" {
				t.Errorf("layout '%s' should be valid but got: %v", layout, issue)
			}
		}
	}
}
