package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// HasHead returns true if the repository has at least one commit (HEAD exists).
func HasHead() bool {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	return cmd.Run() == nil
}

// BranchExists checks if a given branch ref exists (local or remote if full ref passed).
func BranchExists(name string) bool {
	if name == "" {
		return false
	}
	cmd := exec.Command("git", "rev-parse", "--verify", "--quiet", name)
	return cmd.Run() == nil
}

// CurrentBranch returns the current branch name or empty string when detached.
func CurrentBranch() string {
	out, _ := exec.Command("git", "symbolic-ref", "-q", "--short", "HEAD").Output()
	return strings.TrimSpace(string(out))
}

// IsDetachedHEAD returns true if repo is in detached HEAD state.
func IsDetachedHEAD() bool {
	branch := CurrentBranch()
	if branch != "" {
		return false
	}
	// Ensure repo has a HEAD commit
	return exec.Command("git", "rev-parse", "--verify", "HEAD").Run() == nil
}

// InferDefaultWorktreePath builds default worktree path as <worktreesDir>/<base>.
func InferDefaultWorktreePath(worktreesDir, branchOrName string) string {
	base := filepath.Base(branchOrName)
	return filepath.Join(worktreesDir, base)
}

// GitWorktreeAdd executes `git worktree add` with given args and returns combined output on error.
func GitWorktreeAdd(args ...string) error {
	cmd := exec.Command("git", append([]string{"worktree", "add"}, args...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git worktree add failed: %s", strings.TrimSpace(out.String()))
	}
	return nil
}
