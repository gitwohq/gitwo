package cmd

import (
	"fmt"

	"github.com/gitwohq/gitwo/internal/shell"
	"github.com/spf13/cobra"
)

var shellInitShell string

var shellInitCmd = &cobra.Command{
	Use:   "shell-init",
	Short: "Print a shell wrapper to auto-cd after `gitwo new`",
	Long: `Detects your shell (or use --shell) and prints a wrapper so 'gitwo new' auto-cd's.

Examples:
  gitwo shell-init                    # Auto-detect shell and print wrapper
  gitwo shell-init --shell bash       # Print bash wrapper
  gitwo shell-init --shell zsh        # Print zsh wrapper
  gitwo shell-init --shell fish       # Print fish wrapper
  gitwo shell-init --shell powershell # Print PowerShell wrapper

To use the wrapper:
  eval "$(gitwo shell-init)"          # bash/zsh
  gitwo shell-init --shell fish | source  # fish
  gitwo shell-init --shell pwsh | Out-String | Invoke-Expression  # PowerShell`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sh := shellInitShell
		if sh == "" {
			sh = shell.DetectShell()
		}

		wrapper := shell.GenerateWrapper(sh)
		fmt.Print(wrapper)
		return nil
	},
}

func init() {
	shellInitCmd.Flags().StringVar(&shellInitShell, "shell", "", "bash|zsh|fish|powershell")
	rootCmd.AddCommand(shellInitCmd)
}
