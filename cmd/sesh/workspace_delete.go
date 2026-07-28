package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"sesh/internal/output"
)

var workspaceDeleteCmd = &cobra.Command{
	Use:     "delete [workspace-name]",
	Short:   "Delete a workspace",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceName := args[0]

		workspace, err := app.cfgManager.GetWorkspace(workspaceName)
		if err != nil {
			return err
		}

		if len(workspace.Sessions) == 0 {
			if err := app.cfgManager.RemoveWorkspace(workspaceName); err != nil {
				return err
			}
			fmt.Printf("Deleted empty workspace '%s'\n", workspaceName)
			return nil
		}

		var runningSessions []string
		for _, s := range workspace.Sessions {
			if app.tmuxClient.SessionExists(s.TmuxName()) {
				runningSessions = append(runningSessions, s.Name)
			}
		}

		if len(runningSessions) > 0 {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Workspace '%s' has %d running sessions. Kill them and delete?", workspaceName, len(runningSessions)),
				Default: true,
			}
			if err := survey.AskOne(prompt, &confirm); err != nil {
				return fmt.Errorf("failed to get confirmation: %w", err)
			}

			if confirm {
				sessionNames := make([]string, len(workspace.Sessions))
				for i, s := range workspace.Sessions {
					sessionNames[i] = s.TmuxName()
				}
				if err := app.tmuxClient.KillWorkspace(sessionNames); err != nil {
					output.Warn("failed to kill some sessions: %v", err)
				}
				fmt.Printf("Killed %d running sessions\n", len(runningSessions))
			} else {
				fmt.Println("Keeping running sessions")
			}
		}

		if err := app.cfgManager.RemoveWorkspace(workspaceName); err != nil {
			return err
		}

		fmt.Printf("Deleted workspace '%s' with %d sessions\n", workspaceName, len(workspace.Sessions))
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceDeleteCmd)
}
