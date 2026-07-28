package main

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the session configuration file",
	Long:  `Open the session configuration file in your default editor (vim/nvim or $EDITOR).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		editorFlag, _ := cmd.Flags().GetString("editor")
		if editorFlag != "" {
			app.edit.SetEditor(editorFlag)
		}

		return app.edit.OpenFile(app.cfgManager.GetConfigPath())
	},
}

func init() {
	editCmd.Flags().String("editor", "", "Editor to use (vim, nvim, etc.)")

	rootCmd.AddCommand(editCmd)
}
