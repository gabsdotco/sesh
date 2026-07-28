package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"sesh/internal/config"
	"sesh/internal/editor"
	"sesh/internal/tmux"
)

type appContext struct {
	cfgManager *config.Manager
	tmuxClient *tmux.Client
	edit       *editor.Editor
}

func appFromContext(cmd *cobra.Command) *appContext {
	val := cmd.Context().Value(appKey)
	if val == nil {
		return nil
	}
	return val.(*appContext)
}

var appKey = struct{}{}

const version = "0.1.0"

const logo = "\nセッション"

var rootCmd = &cobra.Command{
	Use:     "sesh",
	Short:   "Sesh - Manage predefined TMUX sessions",
	Version: version,
	Long: logo + `

Sesh is a CLI tool that helps you manage predefined 
TMUX sessions with windows and panels. Define your sessions in a YAML 
configuration file and spawn them on demand.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfgManager, err := config.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize config manager: %w", err)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		tmuxClient := tmux.NewClient()
		tmuxClient.Verbose = verbose

		app := &appContext{
			cfgManager: cfgManager,
			tmuxClient: tmuxClient,
			edit:       editor.NewEditor(),
		}

		ctx := context.WithValue(cmd.Context(), appKey, app)
		cmd.SetContext(ctx)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output (log tmux commands)")
}
