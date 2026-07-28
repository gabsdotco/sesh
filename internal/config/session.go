package config

import (
	"fmt"

	"sesh/pkg/models"
)

// Session operations (supports both workspace and orphan sessions)

// GetSession finds a session by name, returns session, workspace name (empty if orphan), and error
func (m *Manager) GetSession(name string) (*models.Session, string, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, "", err
	}

	// First check orphan sessions
	for i := range config.Sessions {
		if config.Sessions[i].Name == name {
			return &config.Sessions[i], "", nil
		}
	}

	// Then check workspace sessions
	for _, workspace := range config.Workspaces {
		for i := range workspace.Sessions {
			if workspace.Sessions[i].Name == name {
				return &workspace.Sessions[i], workspace.Name, nil
			}
		}
	}

	return nil, "", fmt.Errorf("session '%s' not found", name)
}

// GetSessionsByWorkspace returns sessions in a specific workspace
func (m *Manager) GetSessionsByWorkspace(workspaceName string) ([]models.Session, error) {
	workspace, err := m.GetWorkspace(workspaceName)
	if err != nil {
		return nil, err
	}

	return workspace.Sessions, nil
}

// GetOrphanSessions returns all sessions not in a workspace
func (m *Manager) GetOrphanSessions() ([]models.Session, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.Sessions, nil
}

// AddSession adds a session to a workspace (if workspaceName is empty, adds as orphan)
func (m *Manager) AddSession(workspaceName string, session *models.Session) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	// Check if session already exists anywhere
	if m.sessionExistsInConfig(config, session.Name) {
		return fmt.Errorf("session '%s' already exists", session.Name)
	}

	if workspaceName == "" {
		// Add as orphan session
		config.Sessions = append(config.Sessions, *session)
	} else {
		// Find workspace
		workspaceIdx := -1
		for i, w := range config.Workspaces {
			if w.Name == workspaceName {
				workspaceIdx = i
				break
			}
		}

		if workspaceIdx == -1 {
			return fmt.Errorf("workspace '%s' not found", workspaceName)
		}

		config.Workspaces[workspaceIdx].Sessions = append(config.Workspaces[workspaceIdx].Sessions, *session)
	}

	return m.SaveConfig(config)
}

// AddOrphanSession adds a session as an orphan (not in any workspace)
func (m *Manager) AddOrphanSession(session *models.Session) error {
	return m.AddSession("", session)
}

// UpdateSession updates a session (works for both workspace and orphan sessions)
func (m *Manager) UpdateSession(session *models.Session) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	// Check orphan sessions first
	for i := range config.Sessions {
		if config.Sessions[i].Name == session.Name {
			config.Sessions[i] = *session
			return m.SaveConfig(config)
		}
	}

	// Check workspace sessions
	for wi := range config.Workspaces {
		for si := range config.Workspaces[wi].Sessions {
			if config.Workspaces[wi].Sessions[si].Name == session.Name {
				config.Workspaces[wi].Sessions[si] = *session
				return m.SaveConfig(config)
			}
		}
	}

	return fmt.Errorf("session '%s' not found", session.Name)
}

// RemoveSession removes a session (works for both workspace and orphan sessions)
func (m *Manager) RemoveSession(name string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	// Check orphan sessions
	newOrphans := make([]models.Session, 0, len(config.Sessions))
	for _, s := range config.Sessions {
		if s.Name != name {
			newOrphans = append(newOrphans, s)
		} else {
			found = true
		}
	}
	config.Sessions = newOrphans

	// Check workspace sessions
	for wi := range config.Workspaces {
		newSessions := make([]models.Session, 0, len(config.Workspaces[wi].Sessions))
		for _, s := range config.Workspaces[wi].Sessions {
			if s.Name != name {
				newSessions = append(newSessions, s)
			} else {
				found = true
			}
		}
		config.Workspaces[wi].Sessions = newSessions
	}

	if !found {
		return fmt.Errorf("session '%s' not found", name)
	}

	return m.SaveConfig(config)
}

// SessionExists checks if a session exists anywhere
func (m *Manager) SessionExists(name string) bool {
	config, err := m.LoadConfig()
	if err != nil {
		return false
	}

	return m.sessionExistsInConfig(config, name)
}

// sessionExistsInConfig helper to check without reloading
func (m *Manager) sessionExistsInConfig(config *models.Config, name string) bool {
	// Check orphan sessions
	for _, s := range config.Sessions {
		if s.Name == name {
			return true
		}
	}

	// Check workspace sessions
	for _, workspace := range config.Workspaces {
		for _, s := range workspace.Sessions {
			if s.Name == name {
				return true
			}
		}
	}

	return false
}

// ListAllSessions returns all sessions (both workspace and orphan)
func (m *Manager) ListAllSessions() ([]models.Session, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}

	var allSessions []models.Session

	// Add orphan sessions
	allSessions = append(allSessions, config.Sessions...)

	// Add workspace sessions
	for _, workspace := range config.Workspaces {
		allSessions = append(allSessions, workspace.Sessions...)
	}

	return allSessions, nil
}

// GetSessionWorkspace returns the workspace name for a session (empty string if orphan)
func (m *Manager) GetSessionWorkspace(name string) (string, error) {
	_, workspaceName, err := m.GetSession(name)
	return workspaceName, err
}

// RenameSession renames a session in config (updates the Name field)
func (m *Manager) RenameSession(oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf("old and new names are the same")
	}

	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	if m.sessionExistsInConfig(config, newName) {
		return fmt.Errorf("session '%s' already exists", newName)
	}

	found := false

	// Check orphan sessions
	for i := range config.Sessions {
		if config.Sessions[i].Name == oldName {
			config.Sessions[i].Name = newName
			found = true
			break
		}
	}

	// Check workspace sessions
	if !found {
		for wi := range config.Workspaces {
			for si := range config.Workspaces[wi].Sessions {
				if config.Workspaces[wi].Sessions[si].Name == oldName {
					config.Workspaces[wi].Sessions[si].Name = newName
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("session '%s' not found", oldName)
	}

	return m.SaveConfig(config)
}

// MoveSessionToWorkspace moves a session to a workspace (from orphan or another workspace)
func (m *Manager) MoveSessionToWorkspace(sessionName, workspaceName string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	var session *models.Session
	found := false
	sourceWorkspace := ""

	// Search in orphans first
	newOrphans := make([]models.Session, 0, len(config.Sessions))
	for i := range config.Sessions {
		if config.Sessions[i].Name == sessionName {
			session = &config.Sessions[i]
			found = true
		} else {
			newOrphans = append(newOrphans, config.Sessions[i])
		}
	}

	// If not found in orphans, search in workspaces
	if !found {
		for wi := range config.Workspaces {
			newSessions := make([]models.Session, 0, len(config.Workspaces[wi].Sessions))
			for si := range config.Workspaces[wi].Sessions {
				if config.Workspaces[wi].Sessions[si].Name == sessionName {
					session = &config.Workspaces[wi].Sessions[si]
					found = true
					sourceWorkspace = config.Workspaces[wi].Name
				} else {
					newSessions = append(newSessions, config.Workspaces[wi].Sessions[si])
				}
			}
			if found {
				config.Workspaces[wi].Sessions = newSessions
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("session '%s' not found", sessionName)
	}

	if sourceWorkspace == workspaceName {
		return fmt.Errorf("session '%s' is already in workspace '%s'", sessionName, workspaceName)
	}

	// Find destination workspace
	workspaceIdx := -1
	for i, w := range config.Workspaces {
		if w.Name == workspaceName {
			workspaceIdx = i
			break
		}
	}

	if workspaceIdx == -1 {
		return fmt.Errorf("workspace '%s' not found", workspaceName)
	}

	// Check if session already exists in destination workspace
	for _, s := range config.Workspaces[workspaceIdx].Sessions {
		if s.Name == sessionName {
			return fmt.Errorf("session '%s' already exists in workspace '%s'", sessionName, workspaceName)
		}
	}

	// Move session
	config.Sessions = newOrphans
	config.Workspaces[workspaceIdx].Sessions = append(config.Workspaces[workspaceIdx].Sessions, *session)

	return m.SaveConfig(config)
}

// CloneSession creates a copy of a session with a new name
func (m *Manager) CloneSession(name, newName string) error {
	if name == newName {
		return fmt.Errorf("old and new names are the same")
	}

	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	if m.sessionExistsInConfig(config, newName) {
		return fmt.Errorf("session '%s' already exists", newName)
	}

	session, workspaceName, err := m.GetSession(name)
	if err != nil {
		return err
	}

	// Deep copy the session
	cloned := models.Session{
		Name:    newName,
		Windows: make([]models.Window, len(session.Windows)),
	}
	for i, w := range session.Windows {
		cloned.Windows[i] = models.Window{
			Name:    w.Name,
			Layout:  w.Layout,
			WorkDir: w.WorkDir,
			Panels:  make([]models.Panel, len(w.Panels)),
		}
		for j, p := range w.Panels {
			cloned.Windows[i].Panels[j] = models.Panel{
				Command: p.Command,
				WorkDir: p.WorkDir,
			}
		}
	}

	if workspaceName == "" {
		config.Sessions = append(config.Sessions, cloned)
	} else {
		for wi := range config.Workspaces {
			if config.Workspaces[wi].Name == workspaceName {
				config.Workspaces[wi].Sessions = append(config.Workspaces[wi].Sessions, cloned)
				break
			}
		}
	}

	return m.SaveConfig(config)
}

// AddWindow adds a window to a session
func (m *Manager) AddWindow(sessionName string, window models.Window) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	for i := range config.Sessions {
		if config.Sessions[i].Name == sessionName {
			config.Sessions[i].Windows = append(config.Sessions[i].Windows, window)
			found = true
			break
		}
	}

	if !found {
		for wi := range config.Workspaces {
			for si := range config.Workspaces[wi].Sessions {
				if config.Workspaces[wi].Sessions[si].Name == sessionName {
					config.Workspaces[wi].Sessions[si].Windows = append(config.Workspaces[wi].Sessions[si].Windows, window)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("session '%s' not found", sessionName)
	}

	return m.SaveConfig(config)
}

// RemoveWindow removes a window from a session by name
func (m *Manager) RemoveWindow(sessionName, windowName string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	for i := range config.Sessions {
		if config.Sessions[i].Name == sessionName {
			newWindows := make([]models.Window, 0, len(config.Sessions[i].Windows))
			for _, w := range config.Sessions[i].Windows {
				if w.Name != windowName {
					newWindows = append(newWindows, w)
				} else {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
			}
			config.Sessions[i].Windows = newWindows
			return m.SaveConfig(config)
		}
	}

	for wi := range config.Workspaces {
		for si := range config.Workspaces[wi].Sessions {
			if config.Workspaces[wi].Sessions[si].Name == sessionName {
				newWindows := make([]models.Window, 0, len(config.Workspaces[wi].Sessions[si].Windows))
				for _, w := range config.Workspaces[wi].Sessions[si].Windows {
					if w.Name != windowName {
						newWindows = append(newWindows, w)
					} else {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
				}
				config.Workspaces[wi].Sessions[si].Windows = newWindows
				return m.SaveConfig(config)
			}
		}
	}

	return fmt.Errorf("session '%s' not found", sessionName)
}

// AddPanel adds a panel to a window in a session
func (m *Manager) AddPanel(sessionName, windowName string, panel models.Panel) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	for i := range config.Sessions {
		if config.Sessions[i].Name == sessionName {
			for wi := range config.Sessions[i].Windows {
				if config.Sessions[i].Windows[wi].Name == windowName {
					config.Sessions[i].Windows[wi].Panels = append(config.Sessions[i].Windows[wi].Panels, panel)
					found = true
					break
				}
			}
			if found {
				return m.SaveConfig(config)
			}
			return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
		}
	}

	for wi := range config.Workspaces {
		for si := range config.Workspaces[wi].Sessions {
			if config.Workspaces[wi].Sessions[si].Name == sessionName {
				for wi2 := range config.Workspaces[wi].Sessions[si].Windows {
					if config.Workspaces[wi].Sessions[si].Windows[wi2].Name == windowName {
						config.Workspaces[wi].Sessions[si].Windows[wi2].Panels = append(config.Workspaces[wi].Sessions[si].Windows[wi2].Panels, panel)
						found = true
						break
					}
				}
				if found {
					return m.SaveConfig(config)
				}
				return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
			}
		}
	}

	return fmt.Errorf("session '%s' not found", sessionName)
}

// RemovePanel removes a panel from a window by index
func (m *Manager) RemovePanel(sessionName, windowName string, panelIndex int) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	for i := range config.Sessions {
		if config.Sessions[i].Name == sessionName {
			for wi := range config.Sessions[i].Windows {
				if config.Sessions[i].Windows[wi].Name == windowName {
					if panelIndex < 0 || panelIndex >= len(config.Sessions[i].Windows[wi].Panels) {
						return fmt.Errorf("panel index %d out of range (window has %d panels)", panelIndex, len(config.Sessions[i].Windows[wi].Panels))
					}
					config.Sessions[i].Windows[wi].Panels = append(
						config.Sessions[i].Windows[wi].Panels[:panelIndex],
						config.Sessions[i].Windows[wi].Panels[panelIndex+1:]...,
					)
					found = true
					break
				}
			}
			if found {
				return m.SaveConfig(config)
			}
			return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
		}
	}

	for wi := range config.Workspaces {
		for si := range config.Workspaces[wi].Sessions {
			if config.Workspaces[wi].Sessions[si].Name == sessionName {
				for wi2 := range config.Workspaces[wi].Sessions[si].Windows {
					if config.Workspaces[wi].Sessions[si].Windows[wi2].Name == windowName {
						if panelIndex < 0 || panelIndex >= len(config.Workspaces[wi].Sessions[si].Windows[wi2].Panels) {
							return fmt.Errorf("panel index %d out of range (window has %d panels)", panelIndex, len(config.Workspaces[wi].Sessions[si].Windows[wi2].Panels))
						}
						config.Workspaces[wi].Sessions[si].Windows[wi2].Panels = append(
							config.Workspaces[wi].Sessions[si].Windows[wi2].Panels[:panelIndex],
							config.Workspaces[wi].Sessions[si].Windows[wi2].Panels[panelIndex+1:]...,
						)
						found = true
						break
					}
				}
				if found {
					return m.SaveConfig(config)
				}
				return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
			}
		}
	}

	return fmt.Errorf("session '%s' not found", sessionName)
}

// SyncSession updates a session's definition in config with a new structure.
// The session must already exist in config. Only running sessions should be synced.
func (m *Manager) SyncSession(session *models.Session) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	found := false

	// Check orphan sessions
	for i := range config.Sessions {
		if config.Sessions[i].Name == session.Name {
			config.Sessions[i].Windows = session.Windows
			found = true
			break
		}
	}

	// Check workspace sessions
	if !found {
		for wi := range config.Workspaces {
			for si := range config.Workspaces[wi].Sessions {
				if config.Workspaces[wi].Sessions[si].Name == session.Name {
					config.Workspaces[wi].Sessions[si].Windows = session.Windows
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("session '%s' not found in config", session.Name)
	}

	return m.SaveConfig(config)
}

// MoveSessionToOrphan moves a session from a workspace to standalone (orphan)
func (m *Manager) MoveSessionToOrphan(sessionName string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return err
	}

	var session *models.Session
	found := false

	for wi := range config.Workspaces {
		newSessions := make([]models.Session, 0, len(config.Workspaces[wi].Sessions))
		for si := range config.Workspaces[wi].Sessions {
			if config.Workspaces[wi].Sessions[si].Name == sessionName {
				session = &config.Workspaces[wi].Sessions[si]
				found = true
			} else {
				newSessions = append(newSessions, config.Workspaces[wi].Sessions[si])
			}
		}
		if found {
			config.Workspaces[wi].Sessions = newSessions
			break
		}
	}

	if !found {
		return fmt.Errorf("session '%s' not found in any workspace", sessionName)
	}

	// Check if session already exists as orphan
	for _, s := range config.Sessions {
		if s.Name == sessionName {
			return fmt.Errorf("session '%s' already exists as a standalone session", sessionName)
		}
	}

	config.Sessions = append(config.Sessions, *session)

	return m.SaveConfig(config)
}
