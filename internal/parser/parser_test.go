package parser

import (
	"testing"
)

func TestParseWindowDefinitionNameOnly(t *testing.T) {
	window, err := ParseWindowDefinition("editor")
	if err != nil {
		t.Fatalf("ParseWindowDefinition() error = %v", err)
	}
	if window.Name != "editor" {
		t.Errorf("window name = %v, want %v", window.Name, "editor")
	}
	if len(window.Panels) != 1 {
		t.Errorf("expected 1 panel, got %d", len(window.Panels))
	}
}

func TestParseWindowDefinitionWithCount(t *testing.T) {
	window, err := ParseWindowDefinition("editor:3")
	if err != nil {
		t.Fatalf("ParseWindowDefinition() error = %v", err)
	}
	if window.Name != "editor" {
		t.Errorf("window name = %v, want %v", window.Name, "editor")
	}
	if len(window.Panels) != 3 {
		t.Errorf("expected 3 panels, got %d", len(window.Panels))
	}
}

func TestParseWindowDefinitionEmptyName(t *testing.T) {
	_, err := ParseWindowDefinition("")
	if err == nil {
		t.Error("expected error for empty window name")
	}
}

func TestParseWindowDefinitionInvalidCount(t *testing.T) {
	_, err := ParseWindowDefinition("editor:0")
	if err == nil {
		t.Error("expected error for zero panel count")
	}

	_, err = ParseWindowDefinition("editor:abc")
	if err == nil {
		t.Error("expected error for non-numeric panel count")
	}
}

func TestParseWindowDefinitionSpaces(t *testing.T) {
	window, err := ParseWindowDefinition("  editor  :  2  ")
	if err != nil {
		t.Fatalf("ParseWindowDefinition() error = %v", err)
	}
	if window.Name != "editor" {
		t.Errorf("window name = %v, want %v", window.Name, "editor")
	}
	if len(window.Panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(window.Panels))
	}
}

func TestParseWindowDefinitionTooManyParts(t *testing.T) {
	_, err := ParseWindowDefinition("editor:2:extra")
	if err == nil {
		t.Error("expected error for too many colon-separated parts")
	}
}
