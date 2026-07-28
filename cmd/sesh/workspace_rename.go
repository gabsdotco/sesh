package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var workspaceRenameCmd = &cobra.Command{
	Use:   "rename [old-name] [new-name]",
	Short: "Rename a workspace",
	Long:  `Rename a workspace in the config file.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		oldName := args[0]
		newName := args[1]

		if err := app.cfgManager.RenameWorkspace(oldName, newName); err != nil {
			return err
		}

		fmt.Printf("Renamed workspace '%s' -> '%s'\n", oldName, newName)
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceRenameCmd)
}
