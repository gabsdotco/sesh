package tmux

import (
	"fmt"
	"strings"

	"sesh/internal/output"
	"sesh/pkg/models"
)

type Warning struct {
	Message string
}

func (w Warning) Error() string {
	return w.Message
}

type Warnings []Warning

func (w Warnings) Error() string {
	if len(w) == 0 {
		return ""
	}
	messages := make([]string, len(w))
	for i, warning := range w {
		messages[i] = warning.Message
	}
	return strings.Join(messages, "; ")
}

// GetSessionStructure retrieves the complete structure of a tmux session
func (c *Client) GetSessionStructure(sessionName string) (*models.Session, error) {
	return c.GetSessionStructureWithWarnings(sessionName)
}

// GetSessionStructureWithWarnings retrieves the complete structure and collects non-fatal warnings
func (c *Client) GetSessionStructureWithWarnings(sessionName string) (*models.Session, error) {
	if !c.SessionExists(sessionName) {
		return nil, fmt.Errorf("session '%s' does not exist", sessionName)
	}

	windows, err := c.getWindows(sessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get windows: %w", err)
	}

	session := &models.Session{
		Name:    sessionName,
		Windows: make([]models.Window, 0, len(windows)),
	}

	var allWarnings Warnings
	for _, windowInfo := range windows {
		window, warnings, err := c.getWindowStructureWithWarnings(windowInfo.ID, windowInfo.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get window structure for '%s': %w", windowInfo.Name, err)
		}
		allWarnings = append(allWarnings, warnings...)
		session.Windows = append(session.Windows, *window)
	}

	for _, w := range allWarnings {
		output.Warn("%s", w.Message)
	}

	return session, nil
}

// windowInfo holds window identification
type windowInfo struct {
	ID   string
	Name string
}

// getWindows retrieves all windows in a session
func (c *Client) getWindows(sessionName string) ([]windowInfo, error) {
	output, err := c.runOutput("list-windows", "-t", sessionName, "-F", "#{window_id}|#{window_name}")
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}

	lines := strings.Split(output, "\n")
	windows := make([]windowInfo, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			parts := strings.SplitN(line, "|", 2)
			if len(parts) == 2 {
				windows = append(windows, windowInfo{ID: parts[0], Name: parts[1]})
			}
		}
	}

	return windows, nil
}

// getWindowStructureWithWarnings retrieves the structure of a single window, collecting warnings
func (c *Client) getWindowStructureWithWarnings(windowID, windowName string) (*models.Window, Warnings, error) {
	panelCount, err := c.getPanelCount(windowID)
	if err != nil {
		return nil, nil, err
	}

	layout, err := c.getLayout(windowID)
	if err != nil {
		return nil, nil, err
	}

	var warnings Warnings
	panels := make([]models.Panel, panelCount)
	var windowWorkDir string
	for i := 0; i < panelCount; i++ {
		command, cmdErr := c.getPanelCommand(windowID, i)
		if cmdErr != nil {
			warnings = append(warnings, Warning{Message: fmt.Sprintf("panel %d command: %v", i, cmdErr)})
		}
		workDir, dirErr := c.getPanelWorkingDir(windowID, i)
		if dirErr != nil {
			warnings = append(warnings, Warning{Message: fmt.Sprintf("panel %d working directory: %v", i, dirErr)})
		}
		panels[i] = models.Panel{
			Command: command,
			WorkDir: workDir,
		}
		if i == 0 && workDir != "" {
			windowWorkDir = workDir
		}
	}

	return &models.Window{
		Name:    windowName,
		Layout:  layout,
		WorkDir: windowWorkDir,
		Panels:  panels,
	}, warnings, nil
}

// getPanelCount returns the number of panels in a window
func (c *Client) getPanelCount(target string) (int, error) {
	output, err := c.runOutput("list-panes", "-t", target, "-F", "#{pane_id}")
	if err != nil {
		return 0, fmt.Errorf("failed to list panes: %w", err)
	}

	lines := strings.Split(output, "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}

	return count, nil
}

// getLayout returns the layout name for a window
func (c *Client) getLayout(target string) (string, error) {
	output, err := c.runOutput("display-message", "-p", "-t", target, "#{window_layout}")
	if err != nil {
		return "", fmt.Errorf("failed to get layout: %w", err)
	}

	// Map layout names to simpler ones
	switch {
	case strings.Contains(output, "even-horizontal"):
		return "even-horizontal", nil
	case strings.Contains(output, "even-vertical"):
		return "even-vertical", nil
	case strings.Contains(output, "main-horizontal"):
		return "main-horizontal", nil
	case strings.Contains(output, "main-vertical"):
		return "main-vertical", nil
	case strings.Contains(output, "tiled"):
		return "tiled", nil
	default:
		return "", nil
	}
}

// getPanelCommand returns the current command running in a panel
func (c *Client) getPanelCommand(target string, panelIndex int) (string, error) {
	panelTarget := fmt.Sprintf("%s.%d", target, panelIndex)
	output, err := c.runOutput("display-message", "-p", "-t", panelTarget, "#{pane_current_command}")
	if err != nil {
		return "", fmt.Errorf("failed to get panel command: %w", err)
	}

	// Filter out shell processes
	if output == "zsh" || output == "bash" || output == "sh" || output == "fish" {
		return "", nil
	}

	return output, nil
}

// getPanelWorkingDir returns the current working directory of a panel
func (c *Client) getPanelWorkingDir(target string, panelIndex int) (string, error) {
	panelTarget := fmt.Sprintf("%s.%d", target, panelIndex)
	output, err := c.runOutput("display-message", "-p", "-t", panelTarget, "#{pane_current_path}")
	if err != nil {
		return "", fmt.Errorf("failed to get panel working directory: %w", err)
	}

	return output, nil
}
