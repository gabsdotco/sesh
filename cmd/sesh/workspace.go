package main

import (
	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Long:  `Workspace commands for managing groups of sessions.`,
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
