package tmux

import (
	"fmt"

	"sesh/internal/expand"
	"sesh/internal/output"
	"sesh/pkg/models"
)

// CreateSession creates a new tmux session from a session definition
func (c *Client) CreateSession(session *models.Session) error {
	fullName := session.TmuxName()

	if c.SessionExists(fullName) {
		return fmt.Errorf("tmux session '%s' already exists", fullName)
	}

	if len(session.Windows) == 0 {
		return fmt.Errorf("session must have at least one window")
	}

	firstWindow := session.Windows[0]
	if err := c.run("new-session", "-d", "-s", fullName, "-n", firstWindow.Name); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Change to window working directory if specified
	if firstWindow.WorkDir != "" {
		target := getWindowTarget(fullName, firstWindow.Name)
		if err := c.changeDirectory(target, expand.Path(firstWindow.WorkDir)); err != nil {
			output.Warn("failed to change directory for window '%s': %v", firstWindow.Name, err)
		}
	}

	if err := c.setupPanels(fullName, firstWindow.Name, firstWindow.Panels, firstWindow.Layout); err != nil {
		return fmt.Errorf("failed to setup first window panels: %w", err)
	}

	for i := 1; i < len(session.Windows); i++ {
		window := session.Windows[i]
		if err := c.createWindow(fullName, window); err != nil {
			return fmt.Errorf("failed to create window '%s': %w", window.Name, err)
		}
	}

	return nil
}

// createWindow creates a new window in a session
func (c *Client) createWindow(sessionName string, window models.Window) error {
	if err := c.run("new-window", "-t", sessionName, "-n", window.Name); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	// Change to window working directory if specified
	if window.WorkDir != "" {
		target := getWindowTarget(sessionName, window.Name)
		if err := c.changeDirectory(target, expand.Path(window.WorkDir)); err != nil {
			output.Warn("failed to change directory for window '%s': %v", window.Name, err)
		}
	}

	return c.setupPanels(sessionName, window.Name, window.Panels, window.Layout)
}

// setupPanels configures panels in a window
func (c *Client) setupPanels(sessionName, windowName string, panels []models.Panel, layout string) error {
	if len(panels) == 0 {
		return nil
	}

	target := getWindowTarget(sessionName, windowName)

	// Create additional panels by splitting
	for i := 1; i < len(panels); i++ {
		if err := c.run("split-window", "-t", target); err != nil {
			output.Warn("failed to split window for panel %d: %v", i, err)
		}
	}

	// Apply layout if specified
	if layout != "" {
		if err := c.run("select-layout", "-t", target, layout); err != nil {
			output.Warn("failed to set layout: %v", err)
		}
	}

	// Send commands to panels
	for i, panel := range panels {
		panelTarget := getPaneTarget(sessionName, windowName, i)

		// Change to panel working directory if specified
		if panel.WorkDir != "" {
			if err := c.changeDirectory(panelTarget, expand.Path(panel.WorkDir)); err != nil {
				output.Warn("failed to change directory for panel %d: %v", i, err)
			}
		}

		if panel.Command != "" {
			if err := c.sendKeys(panelTarget, panel.Command); err != nil {
				output.Warn("failed to send command '%s' to panel: %v", panel.Command, err)
			}
		}
	}

	return nil
}

// KillSession kills a tmux session with automatic fallback
func (c *Client) KillSession(name string) error {
	sessions, err := c.GetSessions()
	if err != nil {
		return err
	}

	var fallbackSession string
	currentSession, _ := c.GetCurrentSession()

	for _, s := range sessions {
		if s != name {
			fallbackSession = s
			break
		}
	}

	if currentSession == name && fallbackSession != "" {
		if err := c.AttachSession(fallbackSession); err != nil {
			return fmt.Errorf("failed to switch to fallback session: %w", err)
		}
	}

	if err := c.run("kill-session", "-t", name); err != nil {
		return fmt.Errorf("failed to kill session: %w", err)
	}

	return nil
}

// SpawnSession creates and attaches to a session, or attaches if it exists
func (c *Client) SpawnSession(session *models.Session) error {
	fullName := session.TmuxName()

	if !c.SessionExists(fullName) {
		if err := c.CreateSession(session); err != nil {
			return err
		}
	}

	return c.AttachSession(fullName)
}

// RenameSession renames a running tmux session
func (c *Client) RenameSession(oldName, newName string) error {
	if !c.SessionExists(oldName) {
		return fmt.Errorf("tmux session '%s' is not running", oldName)
	}

	if c.SessionExists(newName) {
		return fmt.Errorf("tmux session '%s' already exists", newName)
	}

	if err := c.run("rename-session", "-t", oldName, newName); err != nil {
		return fmt.Errorf("failed to rename session: %w", err)
	}

	return nil
}

// Helper methods

func (c *Client) changeDirectory(target, dir string) error {
	return c.run("send-keys", "-t", target, fmt.Sprintf("cd %s", dir), "C-m")
}

func (c *Client) sendKeys(target, keys string) error {
	return c.run("send-keys", "-t", target, keys, "C-m")
}
