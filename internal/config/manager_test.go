package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"sesh/pkg/models"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "sessions.yaml")
	m := &Manager{configPath: configPath}
	err := m.SaveConfig(&models.Config{
		Workspaces: []models.Workspace{},
		Sessions:   []models.Session{},
	})
	if err != nil {
		t.Fatalf("failed to create initial config: %v", err)
	}
	return m
}

func TestNewManager(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}
	m, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	expectedPath := filepath.Join(homeDir, ".config", "sesh", "sessions.yaml")
	if m.GetConfigPath() != expectedPath {
		t.Errorf("GetConfigPath() = %v, want %v", m.GetConfigPath(), expectedPath)
	}
}

func TestLoadConfig(t *testing.T) {
	m := newTestManager(t)

	cfg, err := m.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if len(cfg.Workspaces) != 0 {
		t.Errorf("expected empty workspaces, got %d", len(cfg.Workspaces))
	}
	if len(cfg.Sessions) != 0 {
		t.Errorf("expected empty sessions, got %d", len(cfg.Sessions))
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	m := newTestManager(t)

	cfg := &models.Config{
		Workspaces: []models.Workspace{
			{
				Name:        "labs",
				Description: "Lab projects",
				Sessions: []models.Session{
					{
						Name: "test-session",
						Windows: []models.Window{
							{
								Name: "editor",
								Panels: []models.Panel{
									{Command: "vim", WorkDir: "/tmp"},
								},
							},
						},
					},
				},
			},
		},
		Sessions: []models.Session{
			{
				Name: "standalone",
				Windows: []models.Window{
					{Name: "main", Panels: []models.Panel{{}}},
				},
			},
		},
	}

	if err := m.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	loaded, err := m.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if len(loaded.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(loaded.Workspaces))
	}
	if loaded.Workspaces[0].Name != "labs" {
		t.Errorf("workspace name = %v, want %v", loaded.Workspaces[0].Name, "labs")
	}
	if len(loaded.Sessions) != 1 {
		t.Errorf("expected 1 standalone session, got %d", len(loaded.Sessions))
	}
	if loaded.Sessions[0].Name != "standalone" {
		t.Errorf("session name = %v, want %v", loaded.Sessions[0].Name, "standalone")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	m := &Manager{configPath: "/nonexistent/path/config.yaml"}
	_, err := m.LoadConfig()
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestSaveConfigInvalidPath(t *testing.T) {
	m := &Manager{configPath: "/nonexistent/dir/config.yaml"}
	cfg := &models.Config{}
	err := m.SaveConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSaveConfigCreatesBackup(t *testing.T) {
	m := newTestManager(t)

	// Create initial config
	initial := &models.Config{
		Sessions: []models.Session{
			{Name: "test", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
		},
	}
	if err := m.SaveConfig(initial); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Save a new config
	updated := &models.Config{
		Sessions: []models.Session{
			{Name: "test", Windows: []models.Window{
				{Name: "editor", Panels: []models.Panel{{}}},
				{Name: "terminal", Panels: []models.Panel{{}}},
			}},
		},
	}
	if err := m.SaveConfig(updated); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Check backup exists and has initial content
	backupPath := m.GetConfigPath() + ".bak"
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("failed to read backup: %v", err)
	}

	var backupConfig models.Config
	if err := yaml.Unmarshal(backupData, &backupConfig); err != nil {
		t.Fatalf("failed to parse backup: %v", err)
	}

	if len(backupConfig.Sessions) != 1 {
		t.Errorf("expected 1 session in backup, got %d", len(backupConfig.Sessions))
	}
	if len(backupConfig.Sessions[0].Windows) != 1 {
		t.Errorf("expected 1 window in backup, got %d", len(backupConfig.Sessions[0].Windows))
	}

	// Check current config has updated content
	current, err := m.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if len(current.Sessions[0].Windows) != 2 {
		t.Errorf("expected 2 windows in current config, got %d", len(current.Sessions[0].Windows))
	}
}

func TestValidateConfig(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}},
	}
	if err := m.AddSession("", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	issues, err := m.ValidateConfig()
	if err != nil {
		t.Fatalf("ValidateConfig() error = %v", err)
	}

	errorCount := 0
	for _, issue := range issues {
		if issue.Severity == "error" {
			errorCount++
		}
	}
	if errorCount > 0 {
		t.Errorf("valid config should have no errors, got %d", errorCount)
	}
}

func TestValidateConfigEmpty(t *testing.T) {
	m := newTestManager(t)

	issues, err := m.ValidateConfig()
	if err != nil {
		t.Fatalf("ValidateConfig() error = %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Message == "config has no workspaces or sessions defined" {
			found = true
		}
	}
	if !found {
		t.Error("empty config should have warning about no sessions")
	}
}

func TestValidateConfigDuplicateSession(t *testing.T) {
	m := newTestManager(t)

	config, err := m.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	config.Workspaces = []models.Workspace{
		{
			Name: "work",
			Sessions: []models.Session{
				{Name: "duplicate", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
			},
		},
	}
	config.Sessions = []models.Session{
		{Name: "duplicate", Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}}},
	}
	if err := m.SaveConfig(config); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	issues, err := m.ValidateConfig()
	if err != nil {
		t.Fatalf("ValidateConfig() error = %v", err)
	}

	found := false
	for _, issue := range issues {
		if issue.Severity == "error" && issue.Message == "duplicate session name 'duplicate' (appears 2 times)" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate session name error, got: %v", issues)
	}
}
