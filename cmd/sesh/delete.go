package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/output"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [project-name]",
	Short:   "Delete a session definition",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm", "remove"},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		killFlag, _ := cmd.Flags().GetBool("kill")

		if dryRun {
			fmt.Printf("[dry-run] Would delete session '%s' from config\n", name)
			if killFlag {
				fmt.Printf("[dry-run] Would kill tmux session '%s' if running\n", name)
			}
			return nil
		}

		if killFlag {
			session, _, err := app.cfgManager.GetSession(name)
			if err != nil {
				return err
			}

			if app.tmuxClient.SessionExists(session.TmuxName()) {
				if err := app.tmuxClient.KillSession(session.TmuxName()); err != nil {
					output.Warn("failed to kill tmux session: %v", err)
				}
			}
		}

		if err := app.cfgManager.RemoveSession(name); err != nil {
			return err
		}

		fmt.Printf("Deleted session '%s'\n", name)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolP("kill", "k", false, "Also kill the tmux session if running")
	deleteCmd.Flags().Bool("dry-run", false, "Show what would be deleted without making changes")

	rootCmd.AddCommand(deleteCmd)
}
