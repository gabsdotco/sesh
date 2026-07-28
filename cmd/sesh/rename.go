package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [old-name] [new-name]",
	Short: "Rename a session",
	Long: `Rename a session in the config file and in the active TMUX session (if running).

Examples:
  sesh rename myapp api-server
  sesh rename backend backend-api --workspace work`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		oldName := args[0]
		newName := args[1]

		// Check that the session exists in config
		_, _, err := app.cfgManager.GetSession(oldName)
		if err != nil {
			return fmt.Errorf("session '%s' not found", oldName)
		}

		// Check that the new name doesn't already exist
		if app.cfgManager.SessionExists(newName) {
			return fmt.Errorf("session '%s' already exists", newName)
		}

		// If the TMUX session is running, rename it there too
		if app.tmuxClient.SessionExists(oldName) {
			if app.tmuxClient.SessionExists(newName) {
				return fmt.Errorf("a TMUX session '%s' is already running", newName)
			}
			if err := app.tmuxClient.RenameSession(oldName, newName); err != nil {
				return fmt.Errorf("failed to rename TMUX session: %w", err)
			}
			fmt.Printf("Renamed TMUX session '%s' -> '%s'\n", oldName, newName)
		}

		// Update the config
		if err := app.cfgManager.RenameSession(oldName, newName); err != nil {
			return fmt.Errorf("failed to rename session in config: %w", err)
		}

		fmt.Printf("Renamed session '%s' -> '%s' in config\n", oldName, newName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
