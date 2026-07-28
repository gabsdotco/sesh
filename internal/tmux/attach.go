package tmux

import (
	"os"
)

// AttachSession attaches to a tmux session (or switches if already inside tmux)
func (c *Client) AttachSession(name string) error {
	if os.Getenv("TMUX") != "" {
		return c.run("switch-client", "-t", name)
	}

	return c.run("attach-session", "-t", name)
}
