package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/parser"
	"sesh/pkg/models"
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new session definition",
	Long: `Create a new TMUX session definition with custom windows and panels.

Define windows using the --window flag with format: name[:panel_count]
Examples:
  sesh create myproject --window "dev" --window "logs:2" --window "shell"
  sesh create myproject -w "editor:2" -w "terminal" -w "build:3"

If no workspace is specified, the session will be created as an orphan (not in any workspace).
If panel_count is omitted, it defaults to 1.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]

		workspaceFlag, _ := cmd.Flags().GetString("workspace")

		windowFlags, _ := cmd.Flags().GetStringArray("window")
		if len(windowFlags) == 0 {
			return fmt.Errorf("at least one window must be defined using --window flag")
		}

		windows := make([]models.Window, 0, len(windowFlags))
		for _, windowDef := range windowFlags {
			window, err := parser.ParseWindowDefinition(windowDef)
			if err != nil {
				return fmt.Errorf("invalid window definition '%s': %w", windowDef, err)
			}
			windows = append(windows, window)
		}

		session := &models.Session{
			Name:    name,
			Windows: windows,
		}

		// Create workspace if specified and doesn't exist
		if workspaceFlag != "" && !app.cfgManager.WorkspaceExists(workspaceFlag) {
			workspace := &models.Workspace{
				Name:        workspaceFlag,
				Description: "",
				Sessions:    []models.Session{},
			}

			if err := app.cfgManager.AddWorkspace(workspace); err != nil {
				return fmt.Errorf("failed to create workspace: %w", err)
			}

			fmt.Printf("Created workspace '%s'\n", workspaceFlag)
		}

		if err := app.cfgManager.AddSession(workspaceFlag, session); err != nil {
			return err
		}

		totalPanels := 0
		for _, w := range windows {
			totalPanels += len(w.Panels)
		}

		if workspaceFlag == "" {
			fmt.Printf("Created orphan session '%s' with %d windows and %d total panels\n", name, len(windows), totalPanels)
		} else {
			fmt.Printf("Created session '%s' in workspace '%s' with %d windows and %d total panels\n", name, workspaceFlag, len(windows), totalPanels)
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringArrayP("window", "w", []string{}, "Window definition (name or name:panel_count). Can be used multiple times.")
	createCmd.MarkFlagRequired("window")
	createCmd.Flags().String("workspace", "", "Workspace to create the session in (if omitted, creates orphan session)")

	rootCmd.AddCommand(createCmd)
}
