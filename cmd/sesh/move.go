package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move [session-name]",
	Short: "Move a session to a different workspace or to standalone",
	Long: `Move a session between workspaces or convert between standalone and workspace sessions.

Examples:
  sesh move myproject --workspace work          # Move to a workspace
  sesh move myproject --standalone              # Move from workspace to standalone
  sesh move myproject -w work                   # Shorthand for --workspace`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)
		name := args[0]

		workspaceFlag, _ := cmd.Flags().GetString("workspace")
		standaloneFlag, _ := cmd.Flags().GetBool("standalone")

		if workspaceFlag == "" && !standaloneFlag {
			return fmt.Errorf("specify either --workspace or --standalone")
		}

		if workspaceFlag != "" && standaloneFlag {
			return fmt.Errorf("cannot specify both --workspace and --standalone")
		}

		_, currentWorkspace, err := app.cfgManager.GetSession(name)
		if err != nil {
			return fmt.Errorf("session '%s' not found", name)
		}

		if standaloneFlag {
			if currentWorkspace == "" {
				return fmt.Errorf("session '%s' is already a standalone session", name)
			}

			if err := app.cfgManager.MoveSessionToOrphan(name); err != nil {
				return fmt.Errorf("failed to move session: %w", err)
			}

			fmt.Printf("Moved session '%s' from workspace '%s' to standalone\n", name, currentWorkspace)
			return nil
		}

		if !app.cfgManager.WorkspaceExists(workspaceFlag) {
			return fmt.Errorf("workspace '%s' not found", workspaceFlag)
		}

		if currentWorkspace == workspaceFlag {
			return fmt.Errorf("session '%s' is already in workspace '%s'", name, workspaceFlag)
		}

		if err := app.cfgManager.MoveSessionToWorkspace(name, workspaceFlag); err != nil {
			return fmt.Errorf("failed to move session: %w", err)
		}

		if currentWorkspace == "" {
			fmt.Printf("Moved standalone session '%s' to workspace '%s'\n", name, workspaceFlag)
		} else {
			fmt.Printf("Moved session '%s' from workspace '%s' to workspace '%s'\n", name, currentWorkspace, workspaceFlag)
		}

		return nil
	},
}

func init() {
	moveCmd.Flags().StringP("workspace", "w", "", "Destination workspace")
	moveCmd.Flags().Bool("standalone", false, "Move session to standalone (orphan)")

	rootCmd.AddCommand(moveCmd)
}
