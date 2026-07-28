package tmux

import (
	"fmt"
	"testing"
)

func TestSessionExists_True(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	result := client.SessionExists("my-session")
	if !result {
		t.Error("SessionExists() = false, want true")
	}

	hasSession := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "has-session", "-t", "=my-session") {
			hasSession = true
		}
	}
	if !hasSession {
		t.Error("expected has-session command to be executed")
	}
}

func TestSessionExists_False(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	result := client.SessionExists("nonexistent")
	if result {
		t.Error("SessionExists() = true, want false")
	}
}

func TestGetSessions(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Output: "session1\nsession2\nsession3"},
		},
	}
	client := NewClientWithRunner(runner)

	sessions, err := client.GetSessions()
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 3 {
		t.Errorf("expected 3 sessions, got %d", len(sessions))
	}
	if sessions[0] != "session1" {
		t.Errorf("session[0] = %v, want %v", sessions[0], "session1")
	}
	if sessions[1] != "session2" {
		t.Errorf("session[1] = %v, want %v", sessions[1], "session2")
	}
	if sessions[2] != "session3" {
		t.Errorf("session[2] = %v, want %v", sessions[2], "session3")
	}
}

func TestGetSessionsSingle(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Output: "only-session"},
		},
	}
	client := NewClientWithRunner(runner)

	sessions, err := client.GetSessions()
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0] != "only-session" {
		t.Errorf("session[0] = %v, want %v", sessions[0], "only-session")
	}
}

func TestGetSessionsEmpty(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Output: ""},
		},
	}
	client := NewClientWithRunner(runner)

	sessions, err := client.GetSessions()
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sessions))
	}
}

func TestGetSessionsNoSessionsRunning(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Error: fmt.Errorf("no sessions found on server")},
		},
	}
	client := NewClientWithRunner(runner)

	sessions, err := client.GetSessions()
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions for 'no sessions' error, got %d", len(sessions))
	}
}

func TestGetSessionsOtherError(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Error: fmt.Errorf("connection refused")},
		},
	}
	client := NewClientWithRunner(runner)

	_, err := client.GetSessions()
	if err == nil {
		t.Error("expected error for non-'no sessions' error")
	}
}

func TestIsNoSessionsError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"no sessions error", fmt.Errorf("no sessions found on server"), true},
		{"no session error lowercase", fmt.Errorf("tmux: no session"), true},
		{"No Session capitalized", fmt.Errorf("No Session found"), true},
		{"other error", fmt.Errorf("connection refused"), false},
		{"unrelated error", fmt.Errorf("permission denied"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNoSessionsError(tt.err)
			if result != tt.expected {
				t.Errorf("isNoSessionsError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsTmuxRunning_True(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	result := client.IsTmuxRunning()
	if !result {
		t.Error("IsTmuxRunning() = false, want true")
	}
}

func TestIsTmuxRunning_False(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls", Error: fmt.Errorf("server not running")},
		},
	}
	client := NewClientWithRunner(runner)

	result := client.IsTmuxRunning()
	if result {
		t.Error("IsTmuxRunning() = true, want false")
	}
}
