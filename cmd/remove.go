package cmd

import (
	"fmt"

	"github.com/gitwohq/gitwo/internal/wt"
	"github.com/spf13/cobra"
)

func init() {
	removeCmd := &cobra.Command{
		Use:     "remove <worktree>",
		Aliases: []string{"rm"},
		Short:   "Remove a worktree",
		Long: `Remove a worktree by path or name.

The worktree can be specified in multiple ways:
- Full path: gitwo remove ../feature-branch
- Name only: gitwo remove feature-branch
- Short alias: gitwo rm feature-branch

Examples:
  gitwo remove ../test-feature
  gitwo remove test-feature
  gitwo rm test-feature`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktree := args[0]
			
			fmt.Printf("Removing worktree: %s\n", worktree)
			
			if err := wt.Remove(worktree); err != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}
			
			fmt.Printf("Successfully removed worktree: %s\n", worktree)
			return nil
		},
	}

	rootCmd.AddCommand(removeCmd)
}
