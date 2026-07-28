package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"sesh/internal/validate"
	"sesh/pkg/models"
)

type Manager struct {
	configPath string
}

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "sesh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "sessions.yaml")
	return NewManagerWithPath(configPath)
}

func NewManagerWithPath(configPath string) (*Manager, error) {
	manager := &Manager{configPath: configPath}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := manager.SaveConfig(&models.Config{
			Workspaces: []models.Workspace{},
			Sessions:   []models.Session{},
		}); err != nil {
			return nil, fmt.Errorf("failed to create initial config: %w", err)
		}
	}

	return manager, nil
}

func (m *Manager) GetConfigPath() string {
	return m.configPath
}

func (m *Manager) LoadConfig() (*models.Config, error) {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (m *Manager) SaveConfig(config *models.Config) error {
	// Backup existing config before overwriting (remove old backup first)
	backupPath := m.configPath + ".bak"
	if existing, err := os.ReadFile(m.configPath); err == nil && len(existing) > 0 {
		_ = os.Remove(backupPath)
		if err := os.WriteFile(backupPath, existing, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (m *Manager) ValidateConfig() ([]validate.Issue, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}
	return validate.Validate(config), nil
}
