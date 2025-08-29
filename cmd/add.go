package cmd

import (
	"fmt"

	"github.com/gitwohq/gitwo/internal/wt"
	"github.com/spf13/cobra"
)

var (
	addBranch   string
	addStartRef string
)

func init() {
	addCmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Create a new worktree",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if addBranch == "" {
				return fmt.Errorf("branch is required, use -b or --branch")
			}
			path := args[0]
			_, err := wt.Add(path, addBranch, addStartRef)
			return err
		},
	}
	addCmd.Flags().StringVarP(&addBranch, "branch", "b", "", "branch name (required)")
	addCmd.Flags().StringVar(&addStartRef, "start-point", "HEAD", "start point ref (default HEAD)")

	rootCmd.AddCommand(addCmd)
}
