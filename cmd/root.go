package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// these variables are filled by main through ldflags
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	rootCmd = &cobra.Command{
		Use:   "gitwo",
		Short: "Git Worktree helper CLI",
		Long:  "gitwo — a small CLI helper around git worktree.",
	}
)

func init() {
	// version via `gitwo version` or `gitwo --version`
	rootCmd.Version = fmt.Sprintf("%s (%s) %s", Version, Commit, Date)
	rootCmd.SetVersionTemplate("gitwo {{.Version}}\n")
	rootCmd.Flags().BoolP("version", "v", false, "print version and exit")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Print(rootCmd.VersionTemplate())
			return
		}
		_ = cmd.Help()
	}
}

// Execute is the entrypoint called from main
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
