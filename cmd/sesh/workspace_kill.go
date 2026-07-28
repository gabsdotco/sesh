package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceKillCmd = &cobra.Command{
	Use:   "kill [workspace-name]",
	Short: "Kill all sessions in a workspace",
	Long: `Kill all running TMUX sessions in a workspace.
Will fallback to any other running session not in this workspace.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceName := args[0]

		sessions, err := app.cfgManager.GetSessionsByWorkspace(workspaceName)
		if err != nil {
			return err
		}

		if len(sessions) == 0 {
			return fmt.Errorf("workspace '%s' has no sessions", workspaceName)
		}

		workspaceSessionNames := make([]string, len(sessions))
		for i, s := range sessions {
			workspaceSessionNames[i] = s.TmuxName()
		}

		if err := app.tmuxClient.KillWorkspace(workspaceSessionNames); err != nil {
			return fmt.Errorf("failed to kill workspace: %w", err)
		}

		fmt.Printf("Killed workspace '%s' with %d sessions\n", workspaceName, len(sessions))
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceKillCmd)
}
