package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for wpsh.

To load completions:

Bash:
  $ source <(wpsh completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ wpsh completion bash > /etc/bash_completion.d/wpsh
  # macOS:
  $ wpsh completion bash > $(brew --prefix)/etc/bash_completion.d/wpsh

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ wpsh completion zsh > "${fpath[1]}/_wpsh"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ wpsh completion fish | source

  # To load completions for each session, execute once:
  $ wpsh completion fish > ~/.config/fish/completions/wpsh.fish

PowerShell:
  PS> wpsh completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> wpsh completion powershell > wpsh.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
