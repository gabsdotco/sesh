package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill [session-name]",
	Short: "Kill a running TMUX session",
	Long: `Kill a running TMUX session and fallback to another open session if currently attached.
If the session name is not in config, attempts to kill it as an untracked tmux session.
Use --all to kill all running sessions, or --workspace to kill all sessions in a workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)

		allFlag, _ := cmd.Flags().GetBool("all")
		workspaceFlag, _ := cmd.Flags().GetString("workspace")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if allFlag {
			sessions, err := app.tmuxClient.GetSessions()
			if err != nil {
				return fmt.Errorf("failed to list sessions: %w", err)
			}

			if dryRun {
				for _, name := range sessions {
					fmt.Printf("[dry-run] Would kill session '%s'\n", name)
				}
				fmt.Printf("[dry-run] Would kill %d session(s)\n", len(sessions))
				return nil
			}

			for _, name := range sessions {
				if err := app.tmuxClient.KillSession(name); err != nil {
					fmt.Printf("Warning: failed to kill session '%s': %v\n", name, err)
				}
			}

			fmt.Printf("Killed %d session(s)\n", len(sessions))
			return nil
		}

		if workspaceFlag != "" {
			sessions, err := app.cfgManager.GetSessionsByWorkspace(workspaceFlag)
			if err != nil {
				return err
			}

			sessionNames := make([]string, len(sessions))
			for i, s := range sessions {
				sessionNames[i] = s.TmuxName()
			}

			if dryRun {
				for _, name := range sessionNames {
					fmt.Printf("[dry-run] Would kill session '%s' in workspace '%s'\n", name, workspaceFlag)
				}
				fmt.Printf("[dry-run] Would kill %d session(s) in workspace '%s'\n", len(sessions), workspaceFlag)
				return nil
			}

			if err := app.tmuxClient.KillWorkspace(sessionNames); err != nil {
				return fmt.Errorf("failed to kill workspace sessions: %w", err)
			}

			fmt.Printf("Killed %d session(s) in workspace '%s'\n", len(sessions), workspaceFlag)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("session name required (or use --all or --workspace)")
		}

		name := args[0]

		session, _, err := app.cfgManager.GetSession(name)
		if err != nil {
			if !app.tmuxClient.SessionExists(name) {
				return fmt.Errorf("session '%s' not found in config and is not running", name)
			}

			if dryRun {
				fmt.Printf("[dry-run] Would kill untracked session '%s'\n", name)
				return nil
			}

			fmt.Printf("Killing untracked session '%s'\n", name)
			return app.tmuxClient.KillSession(name)
		}

		fullName := session.TmuxName()
		if !app.tmuxClient.SessionExists(fullName) {
			return fmt.Errorf("session '%s' is not running", fullName)
		}

		if dryRun {
			fmt.Printf("[dry-run] Would kill session '%s'\n", fullName)
			return nil
		}

		return app.tmuxClient.KillSession(fullName)
	},
}

func init() {
	killCmd.Flags().BoolP("all", "a", false, "Kill all running sessions")
	killCmd.Flags().StringP("workspace", "w", "", "Kill all sessions in a workspace")
	killCmd.Flags().Bool("dry-run", false, "Show what would be killed without making changes")

	rootCmd.AddCommand(killCmd)
}
