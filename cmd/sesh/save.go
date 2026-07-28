package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/pkg/models"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Capture current tmux session into sesh config",
	Long: `Save the current TMUX session structure to the sesh config file.

Primary use case: you created a session manually in tmux and want to persist
it to sesh config so it can be recreated later.

If the session already exists in config, it updates the structure (same as 'sync').
If the session is not in config, it adds it as a new entry.

This command captures:
- Window names
- Number of panels in each window
- Current layout
- Running commands in each panel

Note: You must be inside a TMUX session to use this command.
For sessions already in config, consider using 'sync' or tmux hooks instead.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceFlag, _ := cmd.Flags().GetString("workspace")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Get current session
		sessionName, err := app.tmuxClient.GetCurrentSession()
		if err != nil {
			return fmt.Errorf("failed to get current session: %w\nMake sure you're running inside a TMUX session", err)
		}

		// Get current session structure from TMUX
		session, err := app.tmuxClient.GetSessionStructure(sessionName)
		if err != nil {
			return fmt.Errorf("failed to inspect session: %w", err)
		}

		// Check if session already exists in config
		exists := app.cfgManager.SessionExists(sessionName)

		if exists {
			if dryRun {
				fmt.Printf("[dry-run] Would update session '%s' with %d windows and %d total panels\n",
					sessionName, len(session.Windows), countPanels(session))
				return nil
			}

			// Update existing session
			if err := app.cfgManager.UpdateSession(session); err != nil {
				return fmt.Errorf("failed to update session: %w", err)
			}

			fmt.Printf("Updated session '%s' with %d windows and %d total panels\n",
				sessionName, len(session.Windows), countPanels(session))
		} else {
			if dryRun {
				if workspaceFlag == "" {
					fmt.Printf("[dry-run] Would save new orphan session '%s' with %d windows and %d total panels\n",
						sessionName, len(session.Windows), countPanels(session))
				} else {
					fmt.Printf("[dry-run] Would save new session '%s' to workspace '%s' with %d windows and %d total panels\n",
						sessionName, workspaceFlag, len(session.Windows), countPanels(session))
				}
				return nil
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

			// Add new session
			if err := app.cfgManager.AddSession(workspaceFlag, session); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			if workspaceFlag == "" {
				fmt.Printf("Saved new orphan session '%s' with %d windows and %d total panels\n",
					sessionName, len(session.Windows), countPanels(session))
			} else {
				fmt.Printf("Saved new session '%s' to workspace '%s' with %d windows and %d total panels\n",
					sessionName, workspaceFlag, len(session.Windows), countPanels(session))
			}
		}

		return nil
	},
}

func countPanels(session *models.Session) int {
	count := 0
	for _, w := range session.Windows {
		count += len(w.Panels)
	}
	return count
}

func init() {
	saveCmd.Flags().String("workspace", "", "Workspace to save the session to (if omitted, creates orphan session)")
	saveCmd.Flags().Bool("dry-run", false, "Show what would be saved without making changes")

	rootCmd.AddCommand(saveCmd)
}
