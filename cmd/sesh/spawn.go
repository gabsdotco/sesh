package main

import (
	"github.com/spf13/cobra"
)

var spawnCmd = &cobra.Command{
	Use:   "spawn [project-name]",
	Short: "Spawn a TMUX session",
	Long:  `Create and attach to a TMUX session if it doesn't exist, or attach to an existing one.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]

		session, _, err := app.cfgManager.GetSession(name)
		if err != nil {
			return err
		}

		return app.tmuxClient.SpawnSession(session)
	},
}

func init() {
	rootCmd.AddCommand(spawnCmd)
}
