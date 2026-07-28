package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/parser"
	"sesh/pkg/models"
)

var addWindowCmd = &cobra.Command{
	Use:   "add-window [session-name] [window-definition]",
	Short: "Add a window to an existing session",
	Long:  `Add a new window to a session. Window definition format: name[:panel_count]`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		sessionName := args[0]
		windowDef := args[1]

		window, err := parser.ParseWindowDefinition(windowDef)
		if err != nil {
			return fmt.Errorf("invalid window definition '%s': %w", windowDef, err)
		}

		workDir, _ := cmd.Flags().GetString("workdir")
		layout, _ := cmd.Flags().GetString("layout")
		window.WorkDir = workDir
		window.Layout = layout

		if err := app.cfgManager.AddWindow(sessionName, window); err != nil {
			return err
		}

		fmt.Printf("Added window '%s' to session '%s' (%d panels)\n", window.Name, sessionName, len(window.Panels))
		return nil
	},
}

var removeWindowCmd = &cobra.Command{
	Use:   "remove-window [session-name] [window-name]",
	Short: "Remove a window from a session",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		sessionName := args[0]
		windowName := args[1]

		if err := app.cfgManager.RemoveWindow(sessionName, windowName); err != nil {
			return err
		}

		fmt.Printf("Removed window '%s' from session '%s'\n", windowName, sessionName)
		return nil
	},
}

var addPanelCmd = &cobra.Command{
	Use:   "add-panel [session-name] [window-name]",
	Short: "Add a panel to a window",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		sessionName := args[0]
		windowName := args[1]

		command, _ := cmd.Flags().GetString("command")
		workDir, _ := cmd.Flags().GetString("workdir")

		panel := models.Panel{
			Command: command,
			WorkDir: workDir,
		}

		if err := app.cfgManager.AddPanel(sessionName, windowName, panel); err != nil {
			return err
		}

		fmt.Printf("Added panel to window '%s' in session '%s'\n", windowName, sessionName)
		return nil
	},
}

var removePanelCmd = &cobra.Command{
	Use:   "remove-panel [session-name] [window-name] [panel-index]",
	Short: "Remove a panel from a window by index",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		sessionName := args[0]
		windowName := args[1]

		var panelIndex int
		if _, err := fmt.Sscanf(args[2], "%d", &panelIndex); err != nil {
			return fmt.Errorf("invalid panel index '%s'", args[2])
		}

		if err := app.cfgManager.RemovePanel(sessionName, windowName, panelIndex); err != nil {
			return err
		}

		fmt.Printf("Removed panel %d from window '%s' in session '%s'\n", panelIndex, windowName, sessionName)
		return nil
	},
}

func init() {
	addWindowCmd.Flags().String("workdir", "", "Working directory for the window")
	addWindowCmd.Flags().String("layout", "", "Layout for the window")
	addPanelCmd.Flags().StringP("command", "c", "", "Command to run in the panel")
	addPanelCmd.Flags().String("workdir", "", "Working directory for the panel")

	rootCmd.AddCommand(addWindowCmd)
	rootCmd.AddCommand(removeWindowCmd)
	rootCmd.AddCommand(addPanelCmd)
	rootCmd.AddCommand(removePanelCmd)
}
