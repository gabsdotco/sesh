package tmux

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetSessionStructure(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Output: ""},
			{Pattern: "list-windows -t", Output: "@0|editor\n@1|terminal"},
			{Pattern: "list-panes", Output: "0"},
			{Pattern: "display-message", Output: "even-horizontal"},
		},
	}
	client := NewClientWithRunner(runner)

	session, err := client.GetSessionStructure("my-session")
	if err != nil {
		t.Fatalf("GetSessionStructure() error = %v", err)
	}
	if session.Name != "my-session" {
		t.Errorf("session name = %v, want %v", session.Name, "my-session")
	}
	if len(session.Windows) != 2 {
		t.Errorf("expected 2 windows, got %d", len(session.Windows))
	}
}

func TestGetSessionStructureNotFound(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	_, err := client.GetSessionStructure("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestGetWindowStructure(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-panes", Output: "0\n1\n2"},
			{Pattern: "window_layout", Output: "main-horizontal"},
			{Pattern: "pane_current_command", Output: "vim"},
			{Pattern: "pane_current_path", Output: "/home/user/project"},
		},
	}
	client := NewClientWithRunner(runner)

	window, warnings, err := client.getWindowStructureWithWarnings("@0", "dev")
	if err != nil {
		t.Fatalf("getWindowStructureWithWarnings() error = %v", err)
	}
	if window.Name != "dev" {
		t.Errorf("window name = %v, want %v", window.Name, "dev")
	}
	if window.Layout != "main-horizontal" {
		t.Errorf("window layout = %v, want %v", window.Layout, "main-horizontal")
	}
	if len(window.Panels) != 3 {
		t.Errorf("expected 3 panels, got %d", len(window.Panels))
	}
	if len(warnings) > 0 {
		t.Logf("Warnings collected: %v", warnings)
	}
}

func TestGetPanelCount(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-panes", Output: "0\n1\n2"},
		},
	}
	client := NewClientWithRunner(runner)

	count, err := client.getPanelCount("@0")
	if err != nil {
		t.Fatalf("getPanelCount() error = %v", err)
	}
	if count != 3 {
		t.Errorf("panel count = %d, want 3", count)
	}
}

func TestGetPanelCountSingle(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-panes", Output: "0"},
		},
	}
	client := NewClientWithRunner(runner)

	count, err := client.getPanelCount("@0")
	if err != nil {
		t.Fatalf("getPanelCount() error = %v", err)
	}
	if count != 1 {
		t.Errorf("panel count = %d, want 1", count)
	}
}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{"even-horizontal", "even-horizontal", "even-horizontal"},
		{"even-vertical", "even-vertical", "even-vertical"},
		{"main-horizontal", "main-horizontal", "main-horizontal"},
		{"main-vertical", "main-vertical", "main-vertical"},
		{"tiled", "tiled", "tiled"},
		{"custom layout", "5f83,191x45,0,0{95x45,0,0,1,95x45,0,1,2}", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{
				Responses: []MockResponse{
					{Pattern: "display-message", Output: tt.output},
				},
			}
			client := NewClientWithRunner(runner)

			layout, err := client.getLayout("@0")
			if err != nil {
				t.Fatalf("getLayout() error = %v", err)
			}
			if layout != tt.expected {
				t.Errorf("getLayout() = %v, want %v", layout, tt.expected)
			}
		})
	}
}

func TestGetPanelCommand(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{"vim", "vim", "vim"},
		{"python", "python", "python"},
		{"zsh filtered", "zsh", ""},
		{"bash filtered", "bash", ""},
		{"sh filtered", "sh", ""},
		{"fish filtered", "fish", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &MockRunner{
				Responses: []MockResponse{
					{Pattern: "display-message", Output: tt.output},
				},
			}
			client := NewClientWithRunner(runner)

			cmd, err := client.getPanelCommand("@0", 0)
			if err != nil {
				t.Fatalf("getPanelCommand() error = %v", err)
			}
			if cmd != tt.expected {
				t.Errorf("getPanelCommand() = %v, want %v", cmd, tt.expected)
			}
		})
	}
}

func TestGetPanelWorkingDir(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "display-message", Output: "/home/user/project"},
		},
	}
	client := NewClientWithRunner(runner)

	dir, err := client.getPanelWorkingDir("@0", 0)
	if err != nil {
		t.Fatalf("getPanelWorkingDir() error = %v", err)
	}
	if dir != "/home/user/project" {
		t.Errorf("getPanelWorkingDir() = %v, want %v", dir, "/home/user/project")
	}
}

func TestGetWindows(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-windows", Output: "@0|editor\n@1|terminal\n@2|logs"},
		},
	}
	client := NewClientWithRunner(runner)

	windows, err := client.getWindows("my-session")
	if err != nil {
		t.Fatalf("getWindows() error = %v", err)
	}
	if len(windows) != 3 {
		t.Errorf("expected 3 windows, got %d", len(windows))
	}
	if windows[0].ID != "@0" || windows[0].Name != "editor" {
		t.Errorf("window[0] = %+v, want {ID:@0, Name:editor}", windows[0])
	}
	if windows[1].ID != "@1" || windows[1].Name != "terminal" {
		t.Errorf("window[1] = %+v, want {ID:@1, Name:terminal}", windows[1])
	}
}

func TestGetWindowsEmpty(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-windows", Output: ""},
		},
	}
	client := NewClientWithRunner(runner)

	windows, err := client.getWindows("my-session")
	if err != nil {
		t.Fatalf("getWindows() error = %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected 0 windows for empty output, got %d", len(windows))
	}
}

func TestWarningType(t *testing.T) {
	w := Warning{Message: "test warning"}
	if w.Error() != "test warning" {
		t.Errorf("Warning.Error() = %v, want %v", w.Error(), "test warning")
	}

	warnings := Warnings{
		{Message: "warning 1"},
		{Message: "warning 2"},
	}
	errMsg := warnings.Error()
	if !strings.Contains(errMsg, "warning 1") || !strings.Contains(errMsg, "warning 2") {
		t.Errorf("Warnings.Error() = %v, want combined message", errMsg)
	}
}

func TestGetSessionStructureWithWarnings(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Output: ""},
			{Pattern: "list-windows", Output: "@0|main"},
			{Pattern: "list-panes", Output: "0"},
			{Pattern: "window_layout", Output: "even-horizontal"},
			{Pattern: "pane_current_command", Output: "vim"},
			{Pattern: "pane_current_path", Output: "/home/user"},
		},
	}
	client := NewClientWithRunner(runner)

	session, err := client.GetSessionStructureWithWarnings("test-session")
	if err != nil {
		t.Fatalf("GetSessionStructureWithWarnings() error = %v", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session")
	}
	if session.Name != "test-session" {
		t.Errorf("session name = %v, want %v", session.Name, "test-session")
	}
}

func TestGetSessionStructureWithPanelWarnings(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Output: ""},
			{Pattern: "list-windows", Output: "@0|main"},
			{Pattern: "list-panes", Output: "0"},
			{Pattern: "window_layout", Output: "tiled"},
			{Pattern: "pane_current_command", Error: fmt.Errorf("failed to get command")},
			{Pattern: "pane_current_path", Error: fmt.Errorf("failed to get path")},
		},
	}
	client := NewClientWithRunner(runner)

	session, err := client.GetSessionStructureWithWarnings("error-session")
	if err != nil {
		t.Fatalf("GetSessionStructureWithWarnings() error = %v", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session even with panel warnings")
	}
	if len(session.Windows) != 1 {
		t.Errorf("expected 1 window, got %d", len(session.Windows))
	}
}

func TestGetWindowStructureWithWarnings(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "list-panes", Output: "0"},
			{Pattern: "window_layout", Output: "main-horizontal"},
			{Pattern: "pane_current_command", Output: "vim"},
			{Pattern: "pane_current_path", Output: "/home/user/project"},
		},
	}
	client := NewClientWithRunner(runner)

	window, warnings, err := client.getWindowStructureWithWarnings("@0", "dev")
	if err != nil {
		t.Fatalf("getWindowStructureWithWarnings() error = %v", err)
	}
	if window.Name != "dev" {
		t.Errorf("window name = %v, want %v", window.Name, "dev")
	}
	if window.Layout != "main-horizontal" {
		t.Errorf("window layout = %v, want %v", window.Layout, "main-horizontal")
	}
	if len(window.Panels) != 1 {
		t.Errorf("expected 1 panel, got %d", len(window.Panels))
	}
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings, got %d", len(warnings))
	}
}
