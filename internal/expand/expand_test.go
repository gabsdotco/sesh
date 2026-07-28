package expand

import (
	"os"
	"testing"
)

func TestPath_Tilde(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	result := Path("~/work/project")
	expected := homeDir + "/work/project"
	if result != expected {
		t.Errorf("Path(\"~/work/project\") = %q, want %q", result, expected)
	}
}

func TestPath_DollarHome(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	result := Path("$HOME/work/project")
	expected := homeDir + "/work/project"
	if result != expected {
		t.Errorf("Path(\"$HOME/work/project\") = %q, want %q", result, expected)
	}
}

func TestPath_BracedHome(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	result := Path("${HOME}/work/project")
	expected := homeDir + "/work/project"
	if result != expected {
		t.Errorf("Path(\"${HOME}/work/project\") = %q, want %q", result, expected)
	}
}

func TestPath_EnvVar(t *testing.T) {
	os.Setenv("SESH_TEST_VAR", "/custom/path")
	defer os.Unsetenv("SESH_TEST_VAR")

	result := Path("$SESH_TEST_VAR/sub")
	expected := "/custom/path/sub"
	if result != expected {
		t.Errorf("Path(\"$SESH_TEST_VAR/sub\") = %q, want %q", result, expected)
	}
}

func TestPath_PlainPath(t *testing.T) {
	result := Path("/absolute/path/to/dir")
	if result != "/absolute/path/to/dir" {
		t.Errorf("Path(\"/absolute/path/to/dir\") = %q, want \"/absolute/path/to/dir\"", result)
	}
}

func TestPath_Empty(t *testing.T) {
	result := Path("")
	if result != "" {
		t.Errorf("Path(\"\") = %q, want \"\"", result)
	}
}

func TestPath_JustTilde(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	result := Path("~")
	if result != homeDir {
		t.Errorf("Path(\"~\") = %q, want %q", result, homeDir)
	}
}
