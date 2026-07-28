package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"sesh/pkg/models"
)

var workspaceSpawnCmd = &cobra.Command{
	Use:   "spawn [workspace-name]",
	Short: "Spawn sessions in a workspace (interactive)",
	Long: `Select which sessions to spawn from a workspace.
Shows an interactive prompt to choose sessions.
Use --all to spawn all sessions without prompting.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		workspaceName := args[0]
		allFlag, _ := cmd.Flags().GetBool("all")

		sessions, err := app.cfgManager.GetSessionsByWorkspace(workspaceName)
		if err != nil {
			return err
		}

		if len(sessions) == 0 {
			return fmt.Errorf("workspace '%s' has no sessions", workspaceName)
		}

		var selected []models.Session

		if allFlag {
			selected = sessions
		} else {
			options := make([]string, len(sessions))
			optionToSession := make(map[string]models.Session)
			for i, session := range sessions {
				option := fmt.Sprintf("%s (%d windows)", session.Name, len(session.Windows))
				options[i] = option
				optionToSession[option] = session
			}

			selectedOptions := []string{}
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Select sessions to spawn from workspace '%s':", workspaceName),
				Options: options,
				Default: options,
			}
			err = survey.AskOne(prompt, &selectedOptions)
			if err != nil {
				return fmt.Errorf("failed to get selection: %w", err)
			}

			if len(selectedOptions) == 0 {
				fmt.Println("No sessions selected. Exiting.")
				return nil
			}

			for _, option := range selectedOptions {
				selected = append(selected, optionToSession[option])
			}
		}

		for i, session := range selected {
			if i == 0 {
				fmt.Printf("Spawning session '%s' and attaching...\n", session.Name)
				if err := app.tmuxClient.SpawnSession(&session); err != nil {
					return fmt.Errorf("failed to spawn session '%s': %w", session.Name, err)
				}
			} else {
				if !app.tmuxClient.SessionExists(session.TmuxName()) {
					fmt.Printf("Creating session '%s'...\n", session.Name)
					if err := app.tmuxClient.CreateSession(&session); err != nil {
						fmt.Printf("Warning: failed to create session '%s': %v\n", session.Name, err)
					}
				}
			}
		}

		fmt.Printf("Spawned %d sessions from workspace '%s'\n", len(selected), workspaceName)
		return nil
	},
}

func init() {
	workspaceSpawnCmd.Flags().BoolP("all", "a", false, "Spawn all sessions without prompting")
	workspaceCmd.AddCommand(workspaceSpawnCmd)
}
