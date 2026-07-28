package tmux

import (
	"fmt"
	"strings"
	"testing"

	"sesh/pkg/models"
)

func TestCreateSession(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session -t =test", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name: "test",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{}}},
		},
	}

	err := client.CreateSession(session)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	hasNewSession := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "new-session", "-d", "-s", "test") {
			hasNewSession = true
		}
	}
	if !hasNewSession {
		t.Error("expected new-session command to be run")
	}
}

func TestCreateSessionAlreadyExists(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name: "test",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{}}},
		},
	}

	err := client.CreateSession(session)
	if err == nil {
		t.Error("expected error when session already exists")
	}
}

func TestCreateSessionNoWindows(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{},
	}

	err := client.CreateSession(session)
	if err == nil {
		t.Error("expected error for session with no windows")
	}
}

func TestCreateSessionWithWorkingDir(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name: "test",
		Windows: []models.Window{
			{
				Name:    "editor",
				WorkDir: "/home/user/project",
				Panels:  []models.Panel{{}},
			},
		},
	}

	err := client.CreateSession(session)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	hasCd := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "send-keys", "cd /home/user/project", "C-m") {
			hasCd = true
		}
	}
	if !hasCd {
		t.Error("expected cd command to be sent for working directory")
	}
}

func TestCreateSessionWithPanels(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name: "test",
		Windows: []models.Window{
			{
				Name:   "editor",
				Layout: "even-horizontal",
				Panels: []models.Panel{
					{Command: "vim", WorkDir: "/home/user/project"},
					{WorkDir: "/home/user/logs"},
				},
			},
		},
	}

	err := client.CreateSession(session)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	hasSplit := false
	hasLayout := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "split-window") {
			hasSplit = true
		}
		if containsAll(cmd, "select-layout", "even-horizontal") {
			hasLayout = true
		}
	}
	if !hasSplit {
		t.Error("expected split-window command for additional panel")
	}
	if !hasLayout {
		t.Error("expected select-layout command for layout")
	}
}

func TestCreateSessionMultipleWindows(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name: "test",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{}}},
			{Name: "terminal", Panels: []models.Panel{{}}},
		},
	}

	err := client.CreateSession(session)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	hasNewWindow := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "new-window", "-t", "test", "-n", "terminal") {
			hasNewWindow = true
		}
	}
	if !hasNewWindow {
		t.Error("expected new-window command for second window")
	}
}

func TestSpawnSessionNew(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}},
	}

	err := client.SpawnSession(session)
	if err != nil {
		t.Fatalf("SpawnSession() error = %v", err)
	}

	hasNewSession := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "new-session", "-d", "-s", "test") {
			hasNewSession = true
		}
	}
	if !hasNewSession {
		t.Error("expected new-session command when spawning new session")
	}
}

func TestSpawnSessionExisting(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}},
	}

	err := client.SpawnSession(session)
	if err != nil {
		t.Fatalf("SpawnSession() error = %v", err)
	}

	for _, cmd := range runner.Commands {
		if containsAll(cmd, "new-session") {
			t.Errorf("should not create new session when it already exists, got: %v", cmd)
		}
	}
}

func TestKillSession(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "ls -F", Output: "session1\nsession2"},
			{Pattern: "display-message -p #S", Output: "session1"},
		},
	}
	client := NewClientWithRunner(runner)

	err := client.KillSession("session2")
	if err != nil {
		t.Fatalf("KillSession() error = %v", err)
	}

	hasKill := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "kill-session", "-t", "session2") {
			hasKill = true
		}
	}
	if !hasKill {
		t.Error("expected kill-session command")
	}
}

func TestChangeDirectory(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	err := client.changeDirectory("test:editor.0", "/home/user")
	if err != nil {
		t.Fatalf("changeDirectory() error = %v", err)
	}

	if len(runner.Commands) == 0 {
		t.Fatal("expected command to be run")
	}

	cmd := runner.Commands[0]
	if !containsAll(cmd, "send-keys", "-t", "test:editor.0", "cd /home/user", "C-m") {
		t.Errorf("unexpected command: %v", cmd)
	}
}

func TestSendKeys(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	err := client.sendKeys("test:editor.0", "vim")
	if err != nil {
		t.Fatalf("sendKeys() error = %v", err)
	}

	if len(runner.Commands) == 0 {
		t.Fatal("expected command to be run")
	}

	cmd := runner.Commands[0]
	if !containsAll(cmd, "send-keys", "-t", "test:editor.0", "vim", "C-m") {
		t.Errorf("unexpected command: %v", cmd)
	}
}

func TestRenameSession(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session -t =old-name"},
			{Pattern: "has-session -t =new-name", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	err := client.RenameSession("old-name", "new-name")
	if err != nil {
		t.Fatalf("RenameSession() error = %v", err)
	}

	hasRename := false
	for _, cmd := range runner.Commands {
		if containsAll(cmd, "rename-session", "-t", "old-name", "new-name") {
			hasRename = true
		}
	}
	if !hasRename {
		t.Error("expected rename-session command")
	}
}

func TestRenameSessionNotRunning(t *testing.T) {
	runner := &MockRunner{
		Responses: []MockResponse{
			{Pattern: "has-session", Error: fmt.Errorf("session not found")},
		},
	}
	client := NewClientWithRunner(runner)

	err := client.RenameSession("old-name", "new-name")
	if err == nil {
		t.Error("expected error renaming session that is not running")
	}
}

func TestRenameSessionTargetExists(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	// With no mock errors, both SessionExists calls return true (no error = exists)
	err := client.RenameSession("old-name", "new-name")
	if err == nil {
		t.Error("expected error when target session name already exists")
	}
}

func TestAttachSessionSwitch(t *testing.T) {
	runner := &MockRunner{}
	client := NewClientWithRunner(runner)

	t.Setenv("TMUX", "1")

	err := client.AttachSession("test-session")
	if err != nil {
		t.Fatalf("AttachSession() error = %v", err)
	}

	found := false
	for _, cmd := range runner.Commands {
		if strings.Contains(cmd, "switch-client -t test-session") {
			found = true
		}
	}
	if !found {
		t.Error("expected switch-client command when TMUX is set")
	}
}

func containsAll(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if !containsSubstr(s, sub) {
			return false
		}
	}
	return true
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || strings.Contains(s, sub))
}
