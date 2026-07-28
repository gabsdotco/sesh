package parser

import (
	"fmt"
	"strconv"
	"strings"

	"sesh/pkg/models"
)

func ParseWindowDefinition(def string) (models.Window, error) {
	parts := strings.Split(def, ":")
	if len(parts) == 0 || len(parts) > 2 {
		return models.Window{}, fmt.Errorf("invalid format, expected 'name' or 'name:panel_count'")
	}

	windowName := strings.TrimSpace(parts[0])
	if windowName == "" {
		return models.Window{}, fmt.Errorf("window name cannot be empty")
	}

	panelCount := 1
	if len(parts) == 2 {
		count, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || count < 1 {
			return models.Window{}, fmt.Errorf("panel count must be a positive integer")
		}
		panelCount = count
	}

	panels := make([]models.Panel, panelCount)

	return models.Window{
		Name:   windowName,
		Panels: panels,
	}, nil
}
