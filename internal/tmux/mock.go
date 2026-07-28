package tmux

import (
	"fmt"
	"strings"
)

type MockResponse struct {
	Pattern string
	Output  string
	Error   error
}

type MockRunner struct {
	Commands  []string
	Responses []MockResponse
}

func (m *MockRunner) Run(name string, args ...string) error {
	cmd := fmt.Sprintf("%s %s", name, strings.Join(args, " "))

	m.Commands = append(m.Commands, cmd)

	for _, resp := range m.Responses {
		if strings.Contains(cmd, resp.Pattern) {
			return resp.Error
		}
	}

	return nil
}

func (m *MockRunner) Output(name string, args ...string) (string, error) {
	cmd := fmt.Sprintf("%s %s", name, strings.Join(args, " "))

	m.Commands = append(m.Commands, cmd)

	for _, resp := range m.Responses {
		if strings.Contains(cmd, resp.Pattern) {
			return resp.Output, resp.Error
		}
	}

	return "", nil
}
