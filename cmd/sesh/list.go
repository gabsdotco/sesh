package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"sesh/pkg/models"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all session definitions",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := appFromContext(cmd)

		runningOnly, _ := cmd.Flags().GetBool("running")
		workspaceFilter, _ := cmd.Flags().GetString("workspace")
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")

		configSessionNames := make(map[string]bool)

		orphans, err := app.cfgManager.GetOrphanSessions()
		if err != nil {
			return err
		}

		workspaces, err := app.cfgManager.ListWorkspaces()
		if err != nil {
			return err
		}

		showStandalone := workspaceFilter == ""
		showWorkspaces := workspaceFilter != "" || !standaloneOnly

		if showStandalone {
			var visibleOrphans []models.Session
			for _, session := range orphans {
				if runningOnly && !app.tmuxClient.SessionExists(session.TmuxName()) {
					continue
				}
				visibleOrphans = append(visibleOrphans, session)
			}

			if runningOnly && len(visibleOrphans) == 0 {
				// skip standalone section when --running and nothing is running
			} else {
				fmt.Println("[standalone]")
				if len(visibleOrphans) == 0 {
					fmt.Println("  (no sessions)")
				} else {
					for _, session := range visibleOrphans {
						symbol := "○"
						if app.tmuxClient.SessionExists(session.TmuxName()) {
							symbol = "●"
						}
						configSessionNames[session.Name] = true
						fmt.Printf("  %s %s (%d windows)\n", symbol, session.Name, len(session.Windows))
					}
				}
				fmt.Println()
			}
		}

		if showWorkspaces {
			var visibleWorkspaces []struct {
				workspace models.Workspace
				sessions  []models.Session
			}

			for _, workspace := range workspaces {
				if workspaceFilter != "" && workspace.Name != workspaceFilter {
					continue
				}

				var visibleSessions []models.Session
				for _, session := range workspace.Sessions {
					if runningOnly && !app.tmuxClient.SessionExists(session.TmuxName()) {
						continue
					}
					visibleSessions = append(visibleSessions, session)
				}

				if runningOnly && len(visibleSessions) == 0 {
					continue
				}

				visibleWorkspaces = append(visibleWorkspaces, struct {
					workspace models.Workspace
					sessions  []models.Session
				}{workspace, visibleSessions})
			}

			if len(visibleWorkspaces) > 0 {
				fmt.Println("[workspaces]")
			}

			for _, vw := range visibleWorkspaces {
				workspace := vw.workspace
				sessions := vw.sessions

				if workspace.Description != "" {
					fmt.Printf("  %s: %s\n", workspace.Name, workspace.Description)
				} else {
					fmt.Printf("  %s\n", workspace.Name)
				}

				if len(sessions) == 0 {
					fmt.Println("    (no sessions)")
					continue
				}

				for _, session := range sessions {
					symbol := "○"
					if app.tmuxClient.SessionExists(session.TmuxName()) {
						symbol = "●"
					}
					configSessionNames[session.Name] = true
					fmt.Printf("    %s %s (%d windows)\n", symbol, session.Name, len(session.Windows))
				}
			}
		}

		if !runningOnly && !standaloneOnly && workspaceFilter == "" && app.tmuxClient.IsTmuxRunning() {
			runningSessions, err := app.tmuxClient.GetSessions()
			if err == nil {
				var untracked []string
				for _, name := range runningSessions {
					if !configSessionNames[name] {
						untracked = append(untracked, name)
					}
				}
				if len(untracked) > 0 {
					fmt.Println()
					fmt.Println("[untracked]")
					for _, name := range untracked {
						fmt.Printf("  ◌ %s\n", name)
					}
				}
			}
		}

		if len(workspaces) == 0 && len(orphans) == 0 {
			fmt.Println("No sessions defined. Use 'sesh create' to add one.")
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolP("running", "r", false, "Show only running sessions")
	listCmd.Flags().StringP("workspace", "w", "", "Show only sessions in a workspace")
	listCmd.Flags().BoolP("standalone", "s", false, "Show only standalone sessions")
	rootCmd.AddCommand(listCmd)
}
