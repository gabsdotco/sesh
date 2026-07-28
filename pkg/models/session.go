package models

type Panel struct {
	Command string `yaml:"command,omitempty"`
	WorkDir string `yaml:"workdir,omitempty"`
}

type Window struct {
	Name    string  `yaml:"name"`
	Layout  string  `yaml:"layout,omitempty"`
	WorkDir string  `yaml:"workdir,omitempty"`
	Panels  []Panel `yaml:"panels"`
}

type Session struct {
	Name    string   `yaml:"name"`
	Windows []Window `yaml:"windows"`
}

type Workspace struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Sessions    []Session `yaml:"sessions"`
}

type Config struct {
	Workspaces []Workspace `yaml:"workspaces"`
	Sessions   []Session   `yaml:"sessions,omitempty"`
}

// TmuxName returns the name used to identify this session in tmux.
// Currently this is just the session Name — session names must be globally unique
// across all workspaces.
func (s *Session) TmuxName() string {
	return s.Name
}
