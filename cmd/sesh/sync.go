package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [session-name]",
	Short: "Update existing config session to match current tmux state",
	Long: `Update a session already in config to reflect its current tmux structure.

Use this when you've modified a running session (added/removed windows or panels)
and want the config file updated to match.

Only sessions already in config are synced. Untracked sessions are ignored.
Use 'save' instead to add a new untracked session to config.

If no session name is provided and --all is set, syncs all running sessions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)

		allFlag, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if allFlag {
			sessions, err := app.tmuxClient.GetSessions()
			if err != nil {
				return fmt.Errorf("failed to list tmux sessions: %w", err)
			}

			var syncedCount int
			for _, name := range sessions {
				if !app.cfgManager.SessionExists(name) {
					continue
				}

				structure, err := app.tmuxClient.GetSessionStructure(name)
				if err != nil {
					fmt.Printf("Warning: failed to introspect session '%s': %v\n", name, err)
					continue
				}

				if dryRun {
					fmt.Printf("[dry-run] Would sync session '%s' (%d windows)\n", name, len(structure.Windows))
					syncedCount++
					continue
				}

				if err := app.cfgManager.SyncSession(structure); err != nil {
					fmt.Printf("Warning: failed to sync session '%s': %v\n", name, err)
					continue
				}

				fmt.Printf("Synced session '%s' (%d windows)\n", name, len(structure.Windows))
				syncedCount++
			}

			if dryRun {
				fmt.Printf("[dry-run] Would sync %d session(s)\n", syncedCount)
			} else {
				fmt.Printf("Synced %d session(s)\n", syncedCount)
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("session name required (or use --all)")
		}

		name := args[0]

		if !app.tmuxClient.SessionExists(name) {
			return fmt.Errorf("session '%s' is not running", name)
		}

		if !app.cfgManager.SessionExists(name) {
			return fmt.Errorf("session '%s' not found in config", name)
		}

		structure, err := app.tmuxClient.GetSessionStructure(name)
		if err != nil {
			return fmt.Errorf("failed to introspect session: %w", err)
		}

		if dryRun {
			fmt.Printf("[dry-run] Would sync session '%s' (%d windows)\n", name, len(structure.Windows))
			return nil
		}

		if err := app.cfgManager.SyncSession(structure); err != nil {
			return fmt.Errorf("failed to sync session: %w", err)
		}

		fmt.Printf("Synced session '%s' (%d windows)\n", name, len(structure.Windows))
		return nil
	},
}

func init() {
	syncCmd.Flags().BoolP("all", "a", false, "Sync all running sessions that exist in config")
	syncCmd.Flags().Bool("dry-run", false, "Show what would be synced without making changes")
	rootCmd.AddCommand(syncCmd)
}
