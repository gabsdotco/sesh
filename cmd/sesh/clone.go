package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [name] [new-name]",
	Short: "Clone a session definition",
	Long:  `Create a copy of an existing session definition with a new name.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]
		newName := args[1]

		if err := app.cfgManager.CloneSession(name, newName); err != nil {
			return err
		}

		fmt.Printf("Cloned session '%s' -> '%s'\n", name, newName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
