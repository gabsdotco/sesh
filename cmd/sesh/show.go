package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/expand"
)

var showCmd = &cobra.Command{
	Use:   "show [session-name]",
	Short: "Show session details",
	Long:  `Display detailed information about a session including windows, panels, layout, and working directories.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]

		session, workspaceName, err := app.cfgManager.GetSession(name)
		if err != nil {
			return err
		}

		if workspaceName != "" {
			fmt.Printf("Session: %s (workspace: %s)\n", session.Name, workspaceName)
		} else {
			fmt.Printf("Session: %s (standalone)\n", session.Name)
		}

		running := app.tmuxClient.SessionExists(session.TmuxName())
		status := "not running"
		if running {
			status = "running"
		}
		fmt.Printf("Status:  %s\n", status)
		fmt.Printf("Windows: %d\n\n", len(session.Windows))

		for i, window := range session.Windows {
			fmt.Printf("  Window: %s\n", window.Name)
			if window.WorkDir != "" {
				fmt.Printf("    WorkDir: %s (resolves to %s)\n", window.WorkDir, expand.Path(window.WorkDir))
			}
			if window.Layout != "" {
				fmt.Printf("    Layout:  %s\n", window.Layout)
			}
			for j, panel := range window.Panels {
				label := fmt.Sprintf("Panel %d", j+1)
				if panel.Command != "" || panel.WorkDir != "" {
					fmt.Printf("    %s:\n", label)
					if panel.Command != "" {
						fmt.Printf("      Command: %s\n", panel.Command)
					}
					if panel.WorkDir != "" {
						fmt.Printf("      WorkDir: %s (resolves to %s)\n", panel.WorkDir, expand.Path(panel.WorkDir))
					}
				} else {
					fmt.Printf("    %s: (default shell)\n", label)
				}
			}
			if i < len(session.Windows)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
