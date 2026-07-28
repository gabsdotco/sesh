package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for sesh.

To load completions:

Bash:
  $ source <(sesh completion bash)

  To load completions for each session, add to ~/.bashrc:
  eval "$(sesh completion bash)"

Zsh:
  If shell completion is not already enabled, enable it:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  To load completions for each session, add to ~/.zshrc:
  eval "$(sesh completion zsh)"

Fish:
  $ sesh completion fish | source

  To load completions for each session, add to ~/.config/fish/config.fish:
  sesh completion fish | source

PowerShell:
  PS> sesh completion powershell | Out-File | Invoke-Expression

  To load completions for each session, add to your PowerShell profile:
  Invoke-Expression (sesh completion powershell | Out-String)
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh, fish, powershell)", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
