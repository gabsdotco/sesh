package config

import (
	"testing"

	"sesh/pkg/models"
)

func TestAddWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Description: "Development", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	found, err := m.GetWorkspace("dev")
	if err != nil {
		t.Fatalf("GetWorkspace() error = %v", err)
	}
	if found.Name != "dev" {
		t.Errorf("workspace name = %v, want %v", found.Name, "dev")
	}
	if found.Description != "Development" {
		t.Errorf("workspace description = %v, want %v", found.Description, "Development")
	}
}

func TestAddWorkspaceDuplicate(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	err := m.AddWorkspace(workspace)
	if err == nil {
		t.Error("expected error adding duplicate workspace")
	}
}

func TestGetWorkspaceNotFound(t *testing.T) {
	m := newTestManager(t)

	_, err := m.GetWorkspace("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent workspace")
	}
}

func TestListWorkspaces(t *testing.T) {
	m := newTestManager(t)

	w1 := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	w2 := &models.Workspace{Name: "staging", Sessions: []models.Session{}}
	if err := m.AddWorkspace(w1); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}
	if err := m.AddWorkspace(w2); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	workspaces, err := m.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces() error = %v", err)
	}
	if len(workspaces) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(workspaces))
	}
}

func TestWorkspaceExists(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	if !m.WorkspaceExists("dev") {
		t.Error("WorkspaceExists() = false, want true")
	}
	if m.WorkspaceExists("nonexistent") {
		t.Error("WorkspaceExists() = true, want false")
	}
}

func TestRemoveWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	if err := m.RemoveWorkspace("dev"); err != nil {
		t.Fatalf("RemoveWorkspace() error = %v", err)
	}

	if m.WorkspaceExists("dev") {
		t.Error("workspace should not exist after removal")
	}
}

func TestRemoveWorkspaceNotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.RemoveWorkspace("nonexistent")
	if err == nil {
		t.Error("expected error removing nonexistent workspace")
	}
}

func TestRemoveWorkspaceWithSessions(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	session := &models.Session{Name: "proj", Windows: []models.Window{{Name: "main", Panels: []models.Panel{{}}}}}
	if err := m.AddSession("dev", session); err != nil {
		t.Fatalf("AddSession() error = %v", err)
	}

	if err := m.RemoveWorkspace("dev"); err != nil {
		t.Fatalf("RemoveWorkspace() error = %v", err)
	}

	if m.SessionExists("proj") {
		t.Error("sessions in removed workspace should also be removed")
	}
}

func TestListWorkspacesEmpty(t *testing.T) {
	m := newTestManager(t)

	workspaces, err := m.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces() error = %v", err)
	}
	if len(workspaces) != 0 {
		t.Errorf("expected 0 workspaces, got %d", len(workspaces))
	}
}

func TestRenameWorkspace(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Description: "Development", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	if err := m.RenameWorkspace("dev", "development"); err != nil {
		t.Fatalf("RenameWorkspace() error = %v", err)
	}

	if m.WorkspaceExists("dev") {
		t.Error("old workspace name should not exist")
	}
	if !m.WorkspaceExists("development") {
		t.Error("new workspace name should exist")
	}

	found, err := m.GetWorkspace("development")
	if err != nil {
		t.Fatalf("GetWorkspace() error = %v", err)
	}
	if found.Description != "Development" {
		t.Errorf("description should be preserved, got %q", found.Description)
	}
}

func TestRenameWorkspaceNotFound(t *testing.T) {
	m := newTestManager(t)

	err := m.RenameWorkspace("nonexistent", "new-name")
	if err == nil {
		t.Error("expected error renaming nonexistent workspace")
	}
}

func TestRenameWorkspaceDuplicate(t *testing.T) {
	m := newTestManager(t)

	w1 := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	w2 := &models.Workspace{Name: "staging", Sessions: []models.Session{}}
	if err := m.AddWorkspace(w1); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}
	if err := m.AddWorkspace(w2); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	err := m.RenameWorkspace("dev", "staging")
	if err == nil {
		t.Error("expected error renaming to existing workspace name")
	}
}

func TestRenameWorkspaceSameName(t *testing.T) {
	m := newTestManager(t)

	workspace := &models.Workspace{Name: "dev", Sessions: []models.Session{}}
	if err := m.AddWorkspace(workspace); err != nil {
		t.Fatalf("AddWorkspace() error = %v", err)
	}

	err := m.RenameWorkspace("dev", "dev")
	if err == nil {
		t.Error("expected error when old and new names are the same")
	}
}
