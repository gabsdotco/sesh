package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		fmt.Println(app.cfgManager.GetConfigPath())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
