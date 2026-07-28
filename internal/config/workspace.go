package config

import (
	"fmt"

	"sesh/pkg/models"
)

// Workspace operations

func (m *Manager) GetWorkspace(name string) (*models.Workspace, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}

	for _, workspace := range config.Workspaces {
		if workspace.Name == name {
			return &workspace, nil
		}
	}

	return nil, fmt.Errorf("workspace '%s' not found", name)
}

func (m *Manager) ListWorkspaces() ([]models.Workspace, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.Workspaces, nil
}

func (m *Manager) WorkspaceExists(name string) bool {
	config, err := m.LoadConfig()
	if err != nil {
		return false
	}

	for _, w := range config.Workspaces {
		if w.Name == name {
			return true
		}
	}

	return false
}

func (m *Manager) AddWorkspace(workspace *models.Workspace) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	for _, w := range config.Workspaces {
		if w.Name == workspace.Name {
			return fmt.Errorf("workspace '%s' already exists", workspace.Name)
		}
	}

	config.Workspaces = append(config.Workspaces, *workspace)
	return m.SaveConfig(config)
}

func (m *Manager) RemoveWorkspace(name string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false
	newWorkspaces := make([]models.Workspace, 0, len(config.Workspaces))
	for _, w := range config.Workspaces {
		if w.Name != name {
			newWorkspaces = append(newWorkspaces, w)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("workspace '%s' not found", name)
	}

	config.Workspaces = newWorkspaces
	return m.SaveConfig(config)
}

func (m *Manager) RenameWorkspace(oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf("old and new names are the same")
	}

	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	if m.workspaceExistsInConfig(config, newName) {
		return fmt.Errorf("workspace '%s' already exists", newName)
	}

	found := false
	for i := range config.Workspaces {
		if config.Workspaces[i].Name == oldName {
			config.Workspaces[i].Name = newName
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("workspace '%s' not found", oldName)
	}

	return m.SaveConfig(config)
}

func (m *Manager) workspaceExistsInConfig(config *models.Config, name string) bool {
	for _, w := range config.Workspaces {
		if w.Name == name {
			return true
		}
	}
	return false
}
