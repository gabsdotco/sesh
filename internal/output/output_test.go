package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("created session %s", "test")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "created session test") {
		t.Errorf("expected output to contain 'created session test', got %q", buf.String())
	}
}

func TestWarn(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Warn("failed to kill session %s", "test")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "Warning: failed to kill session test") {
		t.Errorf("expected output to contain 'Warning: failed to kill session test', got %q", buf.String())
	}
}

func TestInfoWritesToStdout(t *testing.T) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()
	os.Stdout = stdoutW
	os.Stderr = stderrW

	Info("hello %s", "world")

	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutBuf.ReadFrom(stdoutR)
	stderrBuf.ReadFrom(stderrR)

	if !strings.Contains(stdoutBuf.String(), "hello world") {
		t.Errorf("Info should write to stdout, got %q", stdoutBuf.String())
	}
	if stderrBuf.String() != "" {
		t.Errorf("Info should not write to stderr, got %q", stderrBuf.String())
	}
}

func TestWarnWritesToStderr(t *testing.T) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()
	os.Stdout = stdoutW
	os.Stderr = stderrW

	Warn("problem %s", "here")

	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutBuf.ReadFrom(stdoutR)
	stderrBuf.ReadFrom(stderrR)

	if stdoutBuf.String() != "" {
		t.Errorf("Warn should not write to stdout, got %q", stdoutBuf.String())
	}
	if !strings.Contains(stderrBuf.String(), "Warning: problem here") {
		t.Errorf("Warn should write to stderr with Warning prefix, got %q", stderrBuf.String())
	}
}

func TestErrorWritesToStderr(t *testing.T) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()
	os.Stdout = stdoutW
	os.Stderr = stderrW

	Error("bad %s", "thing")

	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutBuf.ReadFrom(stdoutR)
	stderrBuf.ReadFrom(stderrR)

	if stdoutBuf.String() != "" {
		t.Errorf("Error should not write to stdout, got %q", stdoutBuf.String())
	}
	if !strings.Contains(stderrBuf.String(), "bad thing") {
		t.Errorf("Error should write to stderr, got %q", stderrBuf.String())
	}
}

func TestErrorDoesNotHavePrefix(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Error("no prefix")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if strings.Contains(buf.String(), "Warning:") {
		t.Errorf("Error should not have 'Warning:' prefix, got %q", buf.String())
	}
}
