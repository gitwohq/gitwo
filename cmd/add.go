package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/gitwohq/gitwo/internal/gitutil"
)

var (
	addBranchFlag string // deprecated alias for positional branch
	addWorktrees  string
	addPath       string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Attach an existing branch as a worktree (strictly like 'git worktree add <path> <branch>')",
	Long: strings.TrimSpace(`
Attach an existing branch to a new worktree. This command will NOT create a branch.
If you need to create one, use: gitwo new <name> [--start-point <ref>]

Examples:
  gitwo add feature/foo
  gitwo add release/1.2 --path ./_wt/release-1-2
`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := strings.TrimSpace(args[0])
		if branch == "" {
			return fmt.Errorf("branch is required, e.g. 'gitwo add feature/foo'")
		}

		// Deprecate --branch in favor of positional
		if addBranchFlag != "" && addBranchFlag != branch {
			fmt.Fprintln(cmd.ErrOrStderr(), "warning: --branch is deprecated; pass the branch as a positional: 'gitwo add <branch>'")
			branch = addBranchFlag
		}

		// Guard: repo must have at least one commit
		if !gitutil.HasHead() {
			return fmt.Errorf("this repository has no commits yet.\nMake an initial commit before using 'gitwo add'")
		}

		// Ensure branch exists (strict with Git semantics)
		if !gitutil.BranchExists(branch) {
			base := filepath.Base(branch)
			return fmt.Errorf("branch %q does not exist.\nTo create it from HEAD: gitwo new %s\nOr from another ref: gitwo new %s --start-point origin/main",
				branch, base, base)
		}

		// Determine worktree path
		worktreesDir := addWorktrees
		if worktreesDir == "" {
			worktreesDir = ".gitwo"
		}
		if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
			return fmt.Errorf("failed to create worktrees dir %q: %w", worktreesDir, err)
		}

		path := addPath
		if path == "" {
			path = gitutil.InferDefaultWorktreePath(worktreesDir, branch)
		}

		// Guard: path must not be an existing non-empty directory
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			entries, _ := os.ReadDir(path)
			if len(entries) > 0 {
				return fmt.Errorf("target path %q already exists and is not empty; choose another path or remove it", path)
			}
		}

		// Execute: git worktree add <path> <branch>
		if err := gitutil.GitWorktreeAdd(path, branch); err != nil {
			// Friendlier message for common cases
			if strings.Contains(strings.ToLower(err.Error()), "is already checked out at") {
				return fmt.Errorf("branch %q is already attached to a worktree.\nUse 'git worktree list' to locate it.", branch)
			}
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Attached branch %q at %s\n", branch, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Deprecated flag: keep for compatibility but discourage
	addCmd.Flags().StringVarP(&addBranchFlag, "branch", "b", "", "[deprecated] specify branch (use positional arg instead)")
	_ = addCmd.Flags().MarkDeprecated("branch", "use positional: gitwo add <branch>")

	// Users may still pass --path
	addCmd.Flags().StringVar(&addPath, "path", "", "explicit worktree path (default: <worktrees-dir>/<basename(branch)>)")
	addCmd.Flags().StringVar(&addWorktrees, "worktrees-dir", "", "directory to place worktrees (default ./.gitwo)")
}
