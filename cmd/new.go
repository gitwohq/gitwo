package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitwohq/gitwo/internal/gitutil"
	"github.com/spf13/cobra"
)

var (
	newStartRef       string
	newAutoSwitch     bool
	newOutputShell    bool
	newCreateScript   bool
	newSourceFunction bool
	newAutoSource     bool
	newPrefix         string
	newWorktreesDir   string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new branch (by default 'feature/<name>') and attach a worktree",
	Long: strings.TrimSpace(`
Create a new branch (default prefix 'feature/') from --start-point (default HEAD) and attach a worktree.

Examples:
  gitwo new auth-refactor
  gitwo new hotfix-123 --prefix ''
  gitwo new rfc-xx --start-point origin/main
`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Avoid := for err anywhere in this function
		var err error

		name := args[0]
		if name == "" {
			return fmt.Errorf("name is required, e.g. 'gitwo new checkout-timeouts'")
		}

		// Guard: repo must have at least one commit
		if !gitutil.HasHead() {
			return fmt.Errorf("this repository has no commits yet.\nMake an initial commit before using 'gitwo new'")
		}

		// Compute branch and path
		branch := name
		if newPrefix != "" {
			branch = fmt.Sprintf("%s%s", newPrefix, name)
		}

		worktreesDir := newWorktreesDir
		if worktreesDir == "" {
			// default to ".gitwo" in repo root unless already configured elsewhere
			worktreesDir = ".gitwo"
		}
		err = os.MkdirAll(worktreesDir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create worktrees dir %q: %w", worktreesDir, err)
		}

		path := filepath.Join(worktreesDir, name)

		// Build git worktree add args
		var wtArgs []string
		if newStartRef == "" {
			wtArgs = []string{"-b", branch, path, "HEAD"}
		} else {
			wtArgs = []string{"-b", branch, path, newStartRef}
		}

		// Run git worktree add
		err = gitutil.GitWorktreeAdd(wtArgs...)
		if err != nil {
			lower := strings.ToLower(err.Error())
			if strings.Contains(lower, "invalid reference: head") {
				return fmt.Errorf("cannot create a branch from HEAD: repository appears to have no commits.\nMake an initial commit, or specify --start-point <ref>")
			}
			return err
		}

		// Print guidance
		fmt.Fprintf(cmd.OutOrStdout(), "Preparing worktree (new branch %q) at %s\n", branch, path)

		// TODO: hook up shell helpers if needed
		_ = newAutoSwitch
		_ = newOutputShell
		_ = newCreateScript
		_ = newSourceFunction
		_ = newAutoSource

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&newStartRef, "start-point", "HEAD", "start point ref (default HEAD)")
	newCmd.Flags().StringVar(&newPrefix, "prefix", "feature/", "branch prefix to use (empty to disable)")
	newCmd.Flags().StringVar(&newWorktreesDir, "worktrees-dir", "", "directory to place worktrees (default ./.gitwo)")

	// keep existing flags (if used by your shell helpers)
	newCmd.Flags().BoolVar(&newAutoSwitch, "switch", false, "automatically switch to the new worktree after creation")
	newCmd.Flags().BoolVar(&newOutputShell, "shell", false, "output shell command for switching to worktree")
	newCmd.Flags().BoolVar(&newCreateScript, "script", false, "create a shell script for easy switching")
	newCmd.Flags().BoolVar(&newSourceFunction, "source", false, "output shell function for sourcing")
	newCmd.Flags().BoolVar(&newAutoSource, "auto-source", false, "provide auto-source instructions")
}

