package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all workspaces",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaces, err := app.cfgManager.ListWorkspaces()
		if err != nil {
			return err
		}

		if len(workspaces) == 0 {
			fmt.Println("No workspaces defined. Use 'sesh create' to create sessions in workspaces.")
			return nil
		}

		fmt.Println("Workspaces:")
		for _, workspace := range workspaces {
			desc := ""
			if workspace.Description != "" {
				desc = fmt.Sprintf(" (%s)", workspace.Description)
			}
			fmt.Printf("  - %s%s: %d sessions\n", workspace.Name, desc, len(workspace.Sessions))
		}

		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceListCmd)
}
