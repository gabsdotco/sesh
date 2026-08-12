package tmux

import (
	"fmt"
	"os"
	"strings"
)

// IsTmuxRunning checks if tmux server is running
func (c *Client) IsTmuxRunning() bool {
	_, err := c.runOutput("ls")
	return err == nil
}

// SessionExists checks if a tmux session with the exact name exists
func (c *Client) SessionExists(name string) bool {
	_, err := c.runOutput("has-session", "-t", "="+name)
	return err == nil
}

// GetCurrentSession returns the name of the current tmux session
func (c *Client) GetCurrentSession() (string, error) {
	if os.Getenv("TMUX") == "" {
		return "", fmt.Errorf("not running inside tmux")
	}

	output, err := c.runOutput("display-message", "-p", "#S")
	if err != nil {
		return "", fmt.Errorf("failed to get current session: %w", err)
	}

	return output, nil
}

// GetSessions returns a list of all tmux session names
func (c *Client) GetSessions() ([]string, error) {
	output, err := c.runOutput("ls", "-F", "#{session_name}")
	if err != nil {
		if isNoSessionsError(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	sessions := make([]string, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			sessions = append(sessions, line)
		}
	}

	return sessions, nil
}

// isNoSessionsError checks if the error indicates no tmux sessions exist
func isNoSessionsError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "no session")
}
