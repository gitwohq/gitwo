package wt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runOut(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func git(args ...string) error {
	return run("git", args...)
}

func gitSilent(args ...string) error {
	cmd := exec.Command("git", args...)
	// Don't output stdout/stderr for silent operations
	return cmd.Run()
}

func gitOut(args ...string) ([]byte, error) {
	return runOut("git", args...)
}

func repoRoot() (string, error) {
	out, err := gitOut("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not a git repository (run inside a repo)")
	}
	return filepath.Clean(string(bytesTrimNL(out))), nil
}

func bytesTrimNL(b []byte) []byte {
	return bytes.TrimRight(b, "\r\n")
}

func refIsRemote(ref string) bool {
	return strings.HasPrefix(ref, "origin/") || strings.HasPrefix(ref, "refs/remotes/")
}

// ChangeToWorktree changes the current directory to the worktree path
func ChangeToWorktree(worktreePath string) error {
	// Get absolute path
	absPath, err := filepath.Abs(worktreePath)
	if err != nil {
		return fmt.Errorf("failed to resolve worktree path: %w", err)
	}

	// Change directory
	if err := os.Chdir(absPath); err != nil {
		return fmt.Errorf("failed to change to worktree directory: %w", err)
	}

	return nil
}
