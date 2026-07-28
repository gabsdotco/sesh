package tmux

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// CommandRunner abstracts command execution for testability
type CommandRunner interface {
	Run(name string, args ...string) error
	Output(name string, args ...string) (string, error)
}

// RealRunner executes commands via os/exec
type RealRunner struct{}

func (r *RealRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *RealRunner) Output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Client provides an interface to interact with tmux
type Client struct {
	runner  CommandRunner
	Verbose bool
}

// NewClient creates a new tmux client
func NewClient() *Client {
	return &Client{runner: &RealRunner{}}
}

// NewClientWithRunner creates a new tmux client with a custom command runner (for testing)
func NewClientWithRunner(runner CommandRunner) *Client {
	return &Client{runner: runner}
}

// run executes a tmux command with the given arguments
func (c *Client) run(args ...string) error {
	if c.Verbose {
		log.Printf("[tmux] %s %s", "tmux", strings.Join(args, " "))
	}
	return c.runner.Run("tmux", args...)
}

// runOutput executes a tmux command and returns the output
func (c *Client) runOutput(args ...string) (string, error) {
	if c.Verbose {
		log.Printf("[tmux] %s %s", "tmux", strings.Join(args, " "))
	}
	return c.runner.Output("tmux", args...)
}

// runOutputSafe executes a tmux command and returns the output,
// ignoring exit errors (useful for commands that might fail)
func (c *Client) runOutputSafe(args ...string) string {
	output, _ := c.runOutput(args...)
	return output
}

// getPaneTarget returns a formatted target string for a pane
func getPaneTarget(session, window string, paneIndex int) string {
	return fmt.Sprintf("%s:%s.%d", session, window, paneIndex)
}

// getWindowTarget returns a formatted target string for a window
func getWindowTarget(session, window string) string {
	return fmt.Sprintf("%s:%s", session, window)
}
