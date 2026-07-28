package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/pkg/models"
)

var workspaceCreateCmd = &cobra.Command{
	Use:   "create [workspace-name]",
	Short: "Create a new workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceName := args[0]
		workspaceDescription, _ := cmd.Flags().GetString("description")

		workspace := &models.Workspace{
			Name:        workspaceName,
			Description: workspaceDescription,
			Sessions:    []models.Session{},
		}

		if !app.cfgManager.WorkspaceExists(workspaceName) {
			if err := app.cfgManager.AddWorkspace(workspace); err != nil {
				return fmt.Errorf("failed to create workspace: %w", err)
			}

			fmt.Printf("Created workspace '%s'\n", workspaceName)

			return nil
		}

		return fmt.Errorf("workspace '%s' already exists", workspaceName)
	},
}

func init() {
	workspaceCreateCmd.Flags().StringP("description", "d", "", "Workspace description")
	workspaceCmd.AddCommand(workspaceCreateCmd)
}
