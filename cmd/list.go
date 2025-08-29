package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/gitwohq/gitwo/internal/wt"
	"github.com/spf13/cobra"
)

var listVerbose bool

func init() {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List worktrees",
		Long: `List worktrees with optional detailed information.

Examples:
  gitwo list                    # Basic list
  gitwo list --verbose          # Detailed information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := wt.List()
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Println("No worktrees found.")
				return nil
			}

			// Get current working directory
			currentDir, err := os.Getwd()
			if err != nil {
				currentDir = ""
			}
			currentDir, _ = filepath.Abs(currentDir)

			tw := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)

			if listVerbose {
				fmt.Fprintln(tw, "PATH\tBRANCH\tHEAD\tSTATUS")
				for _, it := range items {
					status := getWorktreeStatus(it.Path)
					path := formatPath(it.Path, currentDir)
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", path, it.Branch, it.Head, status)
				}
			} else {
				fmt.Fprintln(tw, "PATH\tBRANCH\tHEAD")
				for _, it := range items {
					path := formatPath(it.Path, currentDir)
					fmt.Fprintf(tw, "%s\t%s\t%s\n", path, it.Branch, it.Head)
				}
			}

			return tw.Flush()
		},
	}

	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show detailed information")

	rootCmd.AddCommand(listCmd)
}

func getWorktreeStatus(worktreePath string) string {
	// Simple status check - check if there are uncommitted changes
	// For now, return a basic status
	// TODO: Implement actual status checking
	return "CLEAN"
}

func formatPath(path, currentDir string) string {
	absPath, _ := filepath.Abs(path)
	if absPath == currentDir {
		return fmt.Sprintf("→ %s", path)
	}
	return path
}
