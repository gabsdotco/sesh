package editor

import (
	"fmt"
	"os"
	"os/exec"
)

type Editor struct {
	editor string
}

func NewEditor() *Editor {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return &Editor{editor: editor}
}

func (e *Editor) OpenFile(path string) error {
	cmd := exec.Command(e.editor, path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	return nil
}

func (e *Editor) SetEditor(editor string) {
	e.editor = editor
}

func (e *Editor) GetEditor() string {
	return e.editor
}
