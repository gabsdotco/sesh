package models

import "testing"

func TestSessionTmuxName(t *testing.T) {
	session := &Session{Name: "my-session"}
	result := session.TmuxName()
	if result != "my-session" {
		t.Errorf("TmuxName() = %v, want %v", result, "my-session")
	}
}

func TestSessionTmuxNameEmpty(t *testing.T) {
	session := &Session{Name: ""}
	result := session.TmuxName()
	if result != "" {
		t.Errorf("TmuxName() = %v, want empty string", result)
	}
}

func TestConfigYAML(t *testing.T) {
	config := &Config{
		Workspaces: []Workspace{
			{
				Name:        "dev",
				Description: "Development",
				Sessions: []Session{
					{
						Name: "myapp",
						Windows: []Window{
							{
								Name:   "editor",
								Layout: "main-horizontal",
								Panels: []Panel{
									{Command: "vim", WorkDir: "/home/user/project"},
									{WorkDir: "/home/user/logs"},
								},
							},
						},
					},
				},
			},
		},
		Sessions: []Session{
			{
				Name: "standalone",
				Windows: []Window{
					{Name: "shell", Panels: []Panel{{}}},
				},
			},
		},
	}

	if len(config.Workspaces) != 1 {
		t.Errorf("expected 1 workspace, got %d", len(config.Workspaces))
	}
	if config.Workspaces[0].Name != "dev" {
		t.Errorf("workspace name = %v, want %v", config.Workspaces[0].Name, "dev")
	}
	if config.Workspaces[0].Sessions[0].Name != "myapp" {
		t.Errorf("session name = %v, want %v", config.Workspaces[0].Sessions[0].Name, "myapp")
	}
	if len(config.Sessions) != 1 {
		t.Errorf("expected 1 standalone session, got %d", len(config.Sessions))
	}
	if config.Sessions[0].Name != "standalone" {
		t.Errorf("session name = %v, want %v", config.Sessions[0].Name, "standalone")
	}
}

func TestPanelOmitEmpty(t *testing.T) {
	panel := Panel{}
	if panel.Command != "" || panel.WorkDir != "" {
		t.Error("empty panel fields should be zero values")
	}

	panelWithValues := Panel{Command: "vim", WorkDir: "/tmp"}
	if panelWithValues.Command != "vim" {
		t.Errorf("command = %v, want %v", panelWithValues.Command, "vim")
	}
	if panelWithValues.WorkDir != "/tmp" {
		t.Errorf("workdir = %v, want %v", panelWithValues.WorkDir, "/tmp")
	}
}

func TestWindowOmitEmpty(t *testing.T) {
	window := Window{Name: "main", Panels: []Panel{{}}}
	if window.Layout != "" {
		t.Errorf("empty layout should be zero value, got %v", window.Layout)
	}
	if window.WorkDir != "" {
		t.Errorf("empty workdir should be zero value, got %v", window.WorkDir)
	}
}

func TestWorkspaceOmitEmpty(t *testing.T) {
	workspace := Workspace{Name: "dev", Sessions: []Session{}}
	if workspace.Description != "" {
		t.Errorf("empty description should be zero value, got %v", workspace.Description)
	}
}
