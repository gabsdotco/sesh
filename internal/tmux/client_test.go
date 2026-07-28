package tmux

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestNewClientWithRunner(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)
	if client == nil {
		t.Fatal("NewClientWithRunner() returned nil")
	}
	if client.runner != runner {
		t.Error("runner not set correctly")
	}
}

func TestClientRun(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	err := client.run("new-session", "-d", "-s", "test")
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(runner.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(runner.Commands))
	}
	if runner.Commands[0] != "tmux new-session -d -s test" {
		t.Errorf("command = %v, want %v", runner.Commands[0], "tmux new-session -d -s test")
	}
}

func TestClientRunOutput(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "display-message", Output: "my-session"},
		},
	}
	client := NewClientWithRunner(runner)

	output, err := client.runOutput("display-message", "-p", "#S")
	if err != nil {
		t.Fatalf("runOutput() error = %v", err)
	}
	if output != "my-session" {
		t.Errorf("output = %v, want %v", output, "my-session")
	}
}

func TestClientRunOutputError(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	_, err := client.runOutput("has-session", "-t", "nonexistent")
	if err == nil {
		t.Error("expected error from runOutput")
	}
}

func TestGetPaneTarget(t *testing.T) {
	result := getPaneTarget("my-session", "editor", 0)
	expected := "my-session:editor.0"
	if result != expected {
		t.Errorf("getPaneTarget() = %v, want %v", result, expected)
	}
}

func TestGetWindowTarget(t *testing.T) {
	result := getWindowTarget("my-session", "editor")
	expected := "my-session:editor"
	if result != expected {
		t.Errorf("getWindowTarget() = %v, want %v", result, expected)
	}
}

func TestRealRunner(t *testing.T) {
	runner := &RealRunner{}
	if runner == nil {
		t.Fatal("RealRunner should not be nil")
	}
}
