package editor

import (
	"os"
	"testing"
)

func TestNewEditor(t *testing.T) {
	origEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", origEditor)

	os.Unsetenv("EDITOR")
	e := NewEditor()
	if e.GetEditor() != "vim" {
		t.Errorf("expected default editor 'vim', got %v", e.GetEditor())
	}

	os.Setenv("EDITOR", "nano")
	e = NewEditor()
	if e.GetEditor() != "nano" {
		t.Errorf("expected editor from $EDITOR 'nano', got %v", e.GetEditor())
	}
}

func TestSetEditor(t *testing.T) {
	e := NewEditor()
	e.SetEditor("code")
	if e.GetEditor() != "code" {
		t.Errorf("expected editor 'code', got %v", e.GetEditor())
	}
}

func TestGetEditor(t *testing.T) {
	e := NewEditor()
	editor := e.GetEditor()
	if editor == "" {
		t.Error("expected non-empty editor name")
	}
}
