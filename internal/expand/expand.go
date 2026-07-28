package expand

import (
	"os"
	"path/filepath"
	"strings"
)

func Path(path string) string {
	if path == "" {
		return path
	}

	if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}

	if strings.HasPrefix(path, "$HOME/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[6:])
	}

	if strings.HasPrefix(path, "${HOME}/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[7:])
	}

	return os.ExpandEnv(path)
}
