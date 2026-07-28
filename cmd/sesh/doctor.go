package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sesh/internal/doctor"
	"sesh/internal/validate"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for config and tmux inconsistencies",
	Long: `Validate your session config and compare against running TMUX sessions:
  - Config validation errors and warnings
  - Config sessions not running (stale config)
  - Running TMUX sessions not in config (untracked)
  - Possible renames`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		hasIssues := false

		validationIssues, err := app.cfgManager.ValidateConfig()
		if err != nil {
			return fmt.Errorf("failed to validate config: %w", err)
		}

		if len(validationIssues) > 0 {
			hasIssues = true
			errors := validate.FilterBySeverity(validationIssues, validate.Error)
			warnings := validate.FilterBySeverity(validationIssues, validate.Warning)

			if len(errors) > 0 {
				fmt.Fprintf(os.Stderr, "\033[31mConfig errors:\033[0m\n")
				for _, issue := range errors {
					fmt.Fprintf(os.Stderr, "  - %s\n", issue.Message)
				}
				fmt.Fprintln(os.Stderr)
			}

			if len(warnings) > 0 {
				fmt.Fprintf(os.Stderr, "\033[33mConfig warnings:\033[0m\n")
				for _, issue := range warnings {
					fmt.Fprintf(os.Stderr, "  - %s\n", issue.Message)
				}
				fmt.Fprintln(os.Stderr)
			}
		}

		if !app.tmuxClient.IsTmuxRunning() {
			if !hasIssues {
				fmt.Println("Config is valid. tmux is not running, skipping session sync check.")
			}
			if hasIssues {
				os.Exit(1)
			}
			return nil
		}

		config, err := app.cfgManager.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		runningSessions, err := app.tmuxClient.GetSessions()
		if err != nil {
			return fmt.Errorf("failed to list tmux sessions: %w", err)
		}

		var configSessionNames []string
		for _, s := range config.Sessions {
			configSessionNames = append(configSessionNames, s.Name)
		}
		for _, w := range config.Workspaces {
			for _, s := range w.Sessions {
				configSessionNames = append(configSessionNames, s.Name)
			}
		}

		issues, candidates := doctor.Diagnose(configSessionNames, runningSessions)

		notRunning := doctor.FilterByType(issues, "not-running")
		untracked := doctor.FilterByType(issues, "untracked")

		if len(notRunning) > 0 {
			hasIssues = true
			fmt.Fprintf(os.Stderr, "\033[33mConfig sessions not running:\033[0m\n")
			for _, issue := range notRunning {
				fmt.Fprintf(os.Stderr, "  - %s\n", issue.Name)
			}
			fmt.Fprintln(os.Stderr)
		}

		if len(untracked) > 0 {
			hasIssues = true
			fmt.Fprintf(os.Stderr, "\033[33mRunning TMUX sessions not in config:\033[0m\n")
			for _, issue := range untracked {
				fmt.Fprintf(os.Stderr, "  - %s\n", issue.Name)
			}
			fmt.Fprintln(os.Stderr)

			if len(candidates) > 0 {
				fmt.Fprintf(os.Stderr, "\033[36mPossible renames (use 'sesh rename' to fix):\033[0m\n")
				for _, c := range candidates {
					fmt.Fprintf(os.Stderr, "  - %s -> %s\n", c.OldName, c.NewName)
				}
				fmt.Fprintln(os.Stderr)
			}
		}

		if !hasIssues {
			fmt.Println("\033[32mAll good!\033[0m Config and tmux sessions are in sync.")
		} else {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
