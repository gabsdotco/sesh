package config

import (
	"testing"

	"sesh/pkg/models"
)

func TestGetSessionOrphan(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "my-session",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	found, workspace, err := m.GetSession("my-session")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if found.Name != "my-session" {
		t.Errorf("session name = %v, want %v", found.Name, "my-session")
	}
	if workspace != "" {
		t.Errorf("workspace = %v, want empty string for orphan", workspace)
	}
}

func TestGetSessionInWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "work", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{
		Name:    "project",
		Windows: []models.Window{{Name: "dev", Panels: []models.Panel{{}}}},
	}
	if err := m.AddSession("work", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	found, workspaceName, err := m.GetSession("project")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if found.Name != "project" {
		t.Errorf("session name = %v, want %v", found.Name, "project")
	}
	if workspaceName != "work" {
		t.Errorf("workspace = %v, want %v", workspaceName, "work")
	}
}

func TestGetSessionNotFound(t *testing.T) {
	m := newTestManager(t)

	_, _, err := m.GetSession("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestAddSessionDuplicate(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "dup",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.AddOrphanSession(session)
	if err == nil {
		t.Error("expected error for duplicate session name")
	}
}

func TestAddSessionToNonexistentWorkspace(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	err := m.AddSession("nonexistent", session)
	if err == nil {
		t.Error("expected error adding to nonexistent workspace")
	}
}

func TestUpdateSession(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "myapp",
		Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	updated := &models.Session{
		Name: "myapp",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{}}},
			{Name: "terminal", Panels: []models.Panel{{}}},
		},
	}
	if err := m.UpdateSession(updated); err != nil {
		t.Fatalf("UpdateSession() error = %v", err)
	}

	found, _, err := m.GetSession("myapp")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(found.Windows) != 2 {
		t.Errorf("expected 2 windows, got %d", len(found.Windows))
	}
}

func TestUpdateSessionNotFound(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "nonexistent", Windows: []models.Window{}}
	err := m.UpdateSession(session)
	if err == nil {
		t.Error("expected error when updating nonexistent session")
	}
}

func TestUpdateSessionInWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "labs", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{
		Name:    "exp",
		Windows: []models.Window{{Name: "shell", Panels: []models.Panel{{}}}},
	}
	if err := m.AddSession("labs", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	updated := &models.Session{
		Name: "exp",
		Windows: []models.Window{
			{Name: "shell", Panels: []models.Panel{{}}},
			{Name: "build", Panels: []models.Panel{{}}},
		},
	}
	if err := m.UpdateSession(updated); err != nil {
		t.Fatalf("UpdateSession() error = %v", err)
	}

	found, wsName, err := m.GetSession("exp")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "labs" {
		t.Errorf("workspace = %v, want %v", wsName, "labs")
	}
	if len(found.Windows) != 2 {
		t.Errorf("expected 2 windows after update, got %d", len(found.Windows))
	}
}

func TestRemoveSessionOrphan(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "remove-me",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if err := m.RemoveSession("remove-me"); err != nil {
		t.Fatalf("RemoveSession() error = %v", err)
	}

	if m.SessionExists("remove-me") {
		t.Error("session should not exist after removal")
	}
}

func TestRemoveSessionFromWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{
		Name:    "proj",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	if err := m.RemoveSession("proj"); err != nil {
		t.Fatalf("RemoveSession() error = %v", err)
	}

	if m.SessionExists("proj") {
		t.Error("session should not exist after removal")
	}
}

func TestRemoveSessionNotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.RemoveSession("nonexistent")
	if err == nil {
		t.Error("expected error when removing nonexistent session")
	}
}

func TestSessionExists(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "existing",
		Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if !m.SessionExists("existing") {
		t.Error("SessionExists() = false, want true")
	}
	if m.SessionExists("nonexistent") {
		t.Error("SessionExists() = true, want false")
	}
}

func TestGetOrphanSessions(t *testing.T) {
	m := newTestManager(t)

	s1 := &models.Session{Name: "s1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	s2 := &models.Session{Name: "s2", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(s1); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}
	if err := m.AddOrphanSession(s2); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	orphans, err := m.GetOrphanSessions()
	if err != nil {
		t.Fatalf("GetOrphanSessions() error = %v", err)
	}
	if len(orphans) != 2 {
		t.Errorf("expected 2 orphan sessions, got %d", len(orphans))
	}
}

func TestGetSessionsByWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "team", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	s1 := &models.Session{Name: "proj1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	s2 := &models.Session{Name: "proj2", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("team", s1); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}
	if err := m.AddSession("team", s2); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	sessions, err := m.GetSessionsByWorkspace("team")
	if err != nil {
		t.Fatalf("GetSessionsByWorkspace() error = %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestGetSessionsByWorkspaceNotFound(t *testing.T) {
	m := newTestManager(t)

	_, err := m.GetSessionsByWorkspace("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent workspace")
	}
}

func TestListAllSessions(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "work", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	orphan := &models.Session{Name: "orphan1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	wsSession := &models.Session{Name: "ws1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}

	if err := m.AddOrphanSession(orphan); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}
	if err := m.AddSession("work", wsSession); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	all, err := m.ListAllSessions()
	if err != nil {
		t.Fatalf("ListAllSessions() error = %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 total sessions, got %d", len(all))
	}
}

func TestGetSessionWorkspace(t *testing.T) {
	m := newTestManager(t)

	orphan := &models.Session{Name: "orphan1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(orphan); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	ws, err := m.GetSessionWorkspace("orphan1")
	if err != nil {
		t.Fatalf("GetSessionWorkspace() error = %v", err)
	}
	if ws != "" {
		t.Errorf("workspace = %v, want empty string for orphan", ws)
	}

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}
	wsSession := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", wsSession); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	ws, err = m.GetSessionWorkspace("proj")
	if err != nil {
		t.Fatalf("GetSessionWorkspace() error = %v", err)
	}
	if ws != "dev" {
		t.Errorf("workspace = %v, want %v", ws, "dev")
	}
}

func TestMoveSessionToWorkspace(t *testing.T) {
	m := newTestManager(t)

	orphan := &models.Session{Name: "mover", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(orphan); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	workspace := &models.Workspace{Name: "target", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	if err := m.MoveSessionToWorkspace("mover", "target"); err != nil {
		t.Fatalf("MoveSessionToWorkspace() error = %v", err)
	}

	found, wsName, err := m.GetSession("mover")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "target" {
		t.Errorf("workspace = %v, want %v", wsName, "target")
	}
	if found.Name != "mover" {
		t.Errorf("session name = %v, want %v", found.Name, "mover")
	}

	orphans, _ := m.GetOrphanSessions()
	for _, s := range orphans {
		if s.Name == "mover" {
			t.Error("session should no longer be an orphan")
		}
	}
}

func TestMoveSessionBetweenWorkspaces(t *testing.T) {
	m := newTestManager(t)

	w1 := &models.Workspace{Name: "source", Sessions: []models.Session{}}
	w2 := &models.Workspace{Name: "dest", Sessions: []models.Session{}}
	if err := m.AddWorkspace(w1); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}
	if err := m.AddWorkspace(w2); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("source", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	if err := m.MoveSessionToWorkspace("proj", "dest"); err != nil {
		t.Fatalf("MoveSessionToWorkspace() error = %v", err)
	}

	_, wsName, err := m.GetSession("proj")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "dest" {
		t.Errorf("workspace = %v, want %v", wsName, "dest")
	}

	sourceWS, _ := m.GetWorkspace("source")
	for _, s := range sourceWS.Sessions {
		if s.Name == "proj" {
			t.Error("session should no longer be in source workspace")
		}
	}
}

func TestMoveSessionToWorkspaceSameWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	err := m.MoveSessionToWorkspace("proj", "dev")
	if err == nil {
		t.Error("expected error moving to same workspace")
	}
}

func TestMoveSessionToOrphan(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	if err := m.MoveSessionToOrphan("proj"); err != nil {
		t.Fatalf("MoveSessionToOrphan() error = %v", err)
	}

	_, wsName, err := m.GetSession("proj")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "" {
		t.Errorf("workspace = %v, want empty string for orphan", wsName)
	}

	orphans, _ := m.GetOrphanSessions()
	found := false
	for _, s := range orphans {
		if s.Name == "proj" {
			found = true
		}
	}
	if !found {
		t.Error("session should now be an orphan")
	}
}

func TestMoveSessionToOrphanAlreadyOrphan(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.MoveSessionToOrphan("proj")
	if err == nil {
		t.Error("expected error moving already-orphan session to standalone")
	}
}

func TestMoveSessionToOrphanNotFoundInWorkspace(t *testing.T) {
	m := newTestManager(t)

	err := m.MoveSessionToOrphan("nonexistent")
	if err == nil {
		t.Error("expected error moving nonexistent session to orphan")
	}
}

func TestMoveSessionToWorkspaceNoOrphan(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "target", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	err := m.MoveSessionToWorkspace("nonexistent", "target")
	if err == nil {
		t.Error("expected error moving nonexistent session")
	}
}

func TestMoveSessionToWorkspaceDuplicate(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	orphan := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(orphan); err == nil {
		t.Fatal("expected error: session name 'proj' already exists in workspace")
	}
}

func TestDuplicateSessionNameAcrossWorkspaces(t *testing.T) {
	m := newTestManager(t)

	w1 := &models.Workspace{Name: "team1", Sessions: []models.Session{}}
	w2 := &models.Workspace{Name: "team2", Sessions: []models.Session{}}
	if err := m.AddWorkspace(w1); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}
	if err := m.AddWorkspace(w2); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("team1", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	err := m.AddSession("team2", session)
	if err == nil {
		t.Error("expected error: session name must be unique across all workspaces")
	}
}

func TestRenameSessionOrphan(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "old-name", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if err := m.RenameSession("old-name", "new-name"); err != nil {
		t.Fatalf("RenameSession() error = %v", err)
	}

	if m.SessionExists("old-name") {
		t.Error("old name should not exist after rename")
	}
	if !m.SessionExists("new-name") {
		t.Error("new name should exist after rename")
	}

	session, wsName, err := m.GetSession("new-name")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if session.Name != "new-name" {
		t.Errorf("session name = %v, want %v", session.Name, "new-name")
	}
	if wsName != "" {
		t.Errorf("workspace = %v, want empty string for orphan", wsName)
	}
}

func TestRenameSessionInWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "old-name", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	if err := m.RenameSession("old-name", "new-name"); err != nil {
		t.Fatalf("RenameSession() error = %v", err)
	}

	if m.SessionExists("old-name") {
		t.Error("old name should not exist after rename")
	}
	if !m.SessionExists("new-name") {
		t.Error("new name should exist after rename")
	}

	_, wsName, err := m.GetSession("new-name")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "dev" {
		t.Errorf("workspace = %v, want %v", wsName, "dev")
	}
}

func TestRenameSessionSameName(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "myapp", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.RenameSession("myapp", "myapp")
	if err == nil {
		t.Error("expected error renaming to same name")
	}
}

func TestRenameSessionDuplicateName(t *testing.T) {
	m := newTestManager(t)

	session1 := &models.Session{Name: "app1", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	session2 := &models.Session{Name: "app2", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session1); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}
	if err := m.AddOrphanSession(session2); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.RenameSession("app1", "app2")
	if err == nil {
		t.Error("expected error renaming to existing name")
	}
}

func TestRenameSessionNotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.RenameSession("nonexistent", "new-name")
	if err == nil {
		t.Error("expected error renaming nonexistent session")
	}
}

func TestSyncSession(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name: "myapp",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{Command: "vim"}}},
		},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	updated := &models.Session{
		Name: "myapp",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{Command: "nvim"}, {Command: "git status"}}},
			{Name: "terminal", Panels: []models.Panel{{}}},
		},
	}

	if err := m.SyncSession(updated); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	s, _, err := m.GetSession("myapp")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(s.Windows) != 2 {
		t.Errorf("expected 2 windows after sync, got %d", len(s.Windows))
	}
	if s.Windows[0].Panels[0].Command != "nvim" {
		t.Errorf("expected panel command 'nvim', got %q", s.Windows[0].Panels[0].Command)
	}
}

func TestSyncSessionInWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{
		Name:    "api",
		Windows: []models.Window{{Name: "editor", Panels: []models.Panel{{}}}},
	}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	updated := &models.Session{
		Name: "api",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{Command: "go run main.go"}}},
			{Name: "logs", Panels: []models.Panel{{Command: "tail -f app.log"}}},
		},
	}

	if err := m.SyncSession(updated); err != nil {
		t.Fatalf("SyncSession() error = %v", err)
	}

	s, wsName, err := m.GetSession("api")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "dev" {
		t.Errorf("expected workspace 'dev', got %q", wsName)
	}
	if len(s.Windows) != 2 {
		t.Errorf("expected 2 windows after sync, got %d", len(s.Windows))
	}
}

func TestSyncSessionNotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.SyncSession(&models.Session{Name: "nonexistent", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}})
	if err == nil {
		t.Error("expected error syncing nonexistent session")
	}
}

func TestCloneSession(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name: "original",
		Windows: []models.Window{
			{Name: "editor", Panels: []models.Panel{{Command: "vim"}, {Command: "git status"}}},
			{Name: "terminal", Panels: []models.Panel{{}}},
		},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if err := m.CloneSession("original", "copy"); err != nil {
		t.Fatalf("CloneSession() error = %v", err)
	}

	if !m.SessionExists("copy") {
		t.Error("cloned session should exist")
	}

	cloned, wsName, err := m.GetSession("copy")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if wsName != "" {
		t.Errorf("cloned session should be orphan, got workspace %q", wsName)
	}
	if len(cloned.Windows) != 2 {
		t.Errorf("expected 2 windows, got %d", len(cloned.Windows))
	}
	if len(cloned.Windows[0].Panels) != 2 {
		t.Errorf("expected 2 panels in first window, got %d", len(cloned.Windows[0].Panels))
	}
	if cloned.Windows[0].Panels[0].Command != "vim" {
		t.Errorf("expected panel command 'vim', got %q", cloned.Windows[0].Panels[0].Command)
	}
}

func TestCloneSessionSameName(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "test", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.CloneSession("test", "test")
	if err == nil {
		t.Error("expected error when cloning to same name")
	}
}

func TestCloneSessionDuplicateName(t *testing.T) {
	m := newTestManager(t)

	if err := m.AddOrphanSession(&models.Session{Name: "a", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}
	if err := m.AddOrphanSession(&models.Session{Name: "b", Windows: []models.Window{{Name: "w", Panels: []models.Panel{{}}}}}); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	err := m.CloneSession("a", "b")
	if err == nil {
		t.Error("expected error when cloning to existing name")
	}
}

func TestAddWindow(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "test", Windows: []models.Window{{Name: "w1", Panels: []models.Panel{{}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	window := models.Window{Name: "w2", Panels: []models.Panel{{}, {}}}
	if err := m.AddWindow("test", window); err != nil {
		t.Fatalf("AddWindow() error = %v", err)
	}

	s, _, err := m.GetSession("test")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(s.Windows) != 2 {
		t.Errorf("expected 2 windows, got %d", len(s.Windows))
	}
}

func TestRemoveWindow(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{
		Name:    "test",
		Windows: []models.Window{{Name: "w1", Panels: []models.Panel{{}}}, {Name: "w2", Panels: []models.Panel{{}}}},
	}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if err := m.RemoveWindow("test", "w1"); err != nil {
		t.Fatalf("RemoveWindow() error = %v", err)
	}

	s, _, err := m.GetSession("test")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(s.Windows) != 1 {
		t.Errorf("expected 1 window, got %d", len(s.Windows))
	}
	if s.Windows[0].Name != "w2" {
		t.Errorf("expected window 'w2', got %q", s.Windows[0].Name)
	}
}

func TestAddPanel(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "test", Windows: []models.Window{{Name: "w1", Panels: []models.Panel{{Command: "a"}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	panel := models.Panel{Command: "b"}
	if err := m.AddPanel("test", "w1", panel); err != nil {
		t.Fatalf("AddPanel() error = %v", err)
	}

	s, _, err := m.GetSession("test")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(s.Windows[0].Panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(s.Windows[0].Panels))
	}
}

func TestRemovePanel(t *testing.T) {
	m := newTestManager(t)

	session := &models.Session{Name: "test", Windows: []models.Window{{Name: "w1", Panels: []models.Panel{{Command: "a"}, {Command: "b"}, {Command: "c"}}}}}
	if err := m.AddOrphanSession(session); err != nil {
		t.Fatalf("AddOrphanSession() error = %v", err)
	}

	if err := m.RemovePanel("test", "w1", 1); err != nil {
		t.Fatalf("RemovePanel() error = %v", err)
	}

	s, _, err := m.GetSession("test")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if len(s.Windows[0].Panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(s.Windows[0].Panels))
	}
	if s.Windows[0].Panels[0].Command != "a" || s.Windows[0].Panels[1].Command != "c" {
		t.Errorf("expected panels 'a' and 'c', got %v", s.Windows[0].Panels)
	}
}
