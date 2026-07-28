package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sesh/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate session configuration",
	Long:    `Check your session config for errors and warnings without connecting to tmux.`,
	Aliases: []string{"check"},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		validationIssues, err := app.cfgManager.ValidateConfig()
		if err != nil {
			return fmt.Errorf("failed to validate config: %w", err)
		}

		if len(validationIssues) == 0 {
			fmt.Println("\033[32mConfig is valid.\033[0m")
			return nil
		}

		errors := validate.FilterBySeverity(validationIssues, validate.Error)
		warnings := validate.FilterBySeverity(validationIssues, validate.Warning)

		if len(errors) > 0 {
			fmt.Fprintf(os.Stderr, "\033[31mErrors:\033[0m\n")
			for _, issue := range errors {
				fmt.Fprintf(os.Stderr, "  - %s\n", issue.Message)
			}
			fmt.Fprintln(os.Stderr)
		}

		if len(warnings) > 0 {
			fmt.Fprintf(os.Stderr, "\033[33mWarnings:\033[0m\n")
			for _, issue := range warnings {
				fmt.Fprintf(os.Stderr, "  - %s\n", issue.Message)
			}
			fmt.Fprintln(os.Stderr)
		}

		if len(errors) > 0 {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(validateCmd)
}
