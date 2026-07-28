package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/expand"
)

var workspaceShowCmd = &cobra.Command{
	Use:   "show [workspace-name]",
	Short: "Show workspace details",
	Long:  `Display detailed information about a workspace and its sessions.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceName := args[0]

		workspace, err := app.cfgManager.GetWorkspace(workspaceName)
		if err != nil {
			return err
		}

		fmt.Printf("Workspace: %s\n", workspace.Name)
		if workspace.Description != "" {
			fmt.Printf("Description: %s\n", workspace.Description)
		}
		fmt.Printf("Sessions: %d\n\n", len(workspace.Sessions))

		for i, session := range workspace.Sessions {
			running := app.tmuxClient.SessionExists(session.TmuxName())
			status := "not running"
			if running {
				status = "running"
			}
			fmt.Printf("  Session: %s (%s)\n", session.Name, status)
			fmt.Printf("  Windows: %d\n", len(session.Windows))

			for _, window := range session.Windows {
				fmt.Printf("    - %s", window.Name)
				if window.WorkDir != "" {
					fmt.Printf(" [workdir: %s -> %s]", window.WorkDir, expand.Path(window.WorkDir))
				}
				if window.Layout != "" {
					fmt.Printf(" [layout: %s]", window.Layout)
				}
				fmt.Printf(" (%d panels)\n", len(window.Panels))
			}

			if i < len(workspace.Sessions)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceShowCmd)
}
