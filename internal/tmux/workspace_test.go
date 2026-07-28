package tmux

import (
	"fmt"
	"strings"
	"testing"
)

func TestKillWorkspace(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	err := client.KillWorkspace([]string{"api", "web"})
	if err != nil {
		t.Fatalf("KillWorkspace() error = %v", err)
	}

	foundKillApi := false
	foundKillWeb := false
	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "kill-session -t api") {
			foundKillApi = true
		}
		if strings.Contains(cmd, "kill-session -t web") {
			foundKillWeb = true
		}
	}
	if !foundKillApi || !foundKillWeb {
		t.Error("expected kill-session commands for both api and web")
	}
}

func TestKillWorkspaceSkipsNonexistent(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session -t =missing", Error: fmt.Errorf("not found")},
		},
	}
	client := NewClientWithRunner(runner)

	err := client.KillWorkspace([]string{"missing"})
	if err != nil {
		t.Fatalf("KillWorkspace() error = %v", err)
	}

	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "kill-session") && strings.Contains(cmd, "missing") {
			t.Error("should not call kill-session for nonexistent session")
		}
	}
}

func TestKillWorkspaceSwitchesFromAttachedSession(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "display-message", Output: "api"},
			{Pattern: "ls", Output: "api\nother"},
			{Pattern: "has-session", Error: fmt.Errorf("not found")},
		},
	}
	client := NewClientWithRunner(runner)

	t.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")

	err := client.KillWorkspace([]string{"api"})
	if err != nil {
		t.Fatalf("KillWorkspace() error = %v", err)
	}

	foundSwitch := false
	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "switch-client -t other") {
			foundSwitch = true
		}
	}
	if !foundSwitch {
		t.Errorf("expected switch-client command to fallback session\nCommands: %v", runner.Commands)
	}
}

func TestKillWorkspaceDetachesWhenNoFallback(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "display-message", Output: "api"},
			{Pattern: "ls", Output: "api"},
		},
	}
	client := NewClientWithRunner(runner)

	t.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")

	err := client.KillWorkspace([]string{"api"})
	if err != nil {
		t.Fatalf("KillWorkspace() error = %v", err)
	}

	foundDetach := false
	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "detach-client") {
			foundDetach = true
		}
	}
	if !foundDetach {
		t.Error("expected detach-client when no fallback session available")
	}
}

func TestKillWorkspaceNoSwitchWhenNotAttached(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "display-message", Output: "other-session"},
			{Pattern: "ls", Output: "other-session\napi"},
		},
	}
	client := NewClientWithRunner(runner)

	t.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")

	err := client.KillWorkspace([]string{"api"})
	if err != nil {
		t.Fatalf("KillWorkspace() error = %v", err)
	}

	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "switch-client") || strings.Contains(cmd, "detach-client") {
			t.Errorf("should not switch or detach when not attached to workspace session, got: %s", cmd)
		}
	}
}
