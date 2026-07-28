package tmux

import (
	"os"
	"os/exec"
)

// AttachSession attaches to a tmux session (or switches if already inside tmux)
func (c *Client) AttachSession(name string) error {
	if os.Getenv("TMUX") != "" {
		return c.run("switch-client", "-t", name)
	}

	cmd := exec.Command("tmux", "attach-session", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
