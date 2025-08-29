package wt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Remove(worktree string) error {
	// Validate input
	if worktree == "" {
		return fmt.Errorf("worktree cannot be empty")
	}

	// Ensure inside a repo
	if _, err := repoRoot(); err != nil {
		return err
	}

	// Determine the actual worktree path
	worktreePath := worktree
	
	// If worktree doesn't start with / or ../, assume it's a relative path
	if !strings.HasPrefix(worktree, "/") && !strings.HasPrefix(worktree, "../") && !strings.HasPrefix(worktree, "./") {
		worktreePath = filepath.Join("..", worktree)
	}

	// Check if worktree exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree does not exist: %s", worktreePath)
	}

	// Remove the worktree using git worktree remove
	if err := git("worktree", "remove", worktreePath); err != nil {
		return fmt.Errorf("failed to remove worktree %s: %w", worktreePath, err)
	}

	return nil
}
