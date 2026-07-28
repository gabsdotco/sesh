package tmux

import (
	"sesh/internal/output"
)

// KillWorkspace kills all sessions in a workspace, falling back to any non-workspace session
func (c *Client) KillWorkspace(workspaceSessionNames []string) error {
	workspaceSet := make(map[string]bool)
	for _, name := range workspaceSessionNames {
		workspaceSet[name] = true
	}

	currentSession, _ := c.GetCurrentSession()

	var fallbackSession string
	allSessions, err := c.GetSessions()
	if err == nil {
		for _, s := range allSessions {
			if !workspaceSet[s] {
				fallbackSession = s
				break
			}
		}
	}

	if currentSession != "" && workspaceSet[currentSession] {
		if fallbackSession != "" {
			if err := c.run("switch-client", "-t", fallbackSession); err != nil {
				output.Warn("failed to switch to fallback session: %v", err)
			}
		} else {
			if err := c.run("detach-client"); err != nil {
				output.Warn("failed to detach: %v", err)
			}
		}
	}

	for _, name := range workspaceSessionNames {
		if c.SessionExists(name) {
			if err := c.run("kill-session", "-t", name); err != nil {
				output.Warn("failed to kill session '%s': %v", name, err)
			}
		}
	}

	return nil
}
