package wt

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddWithConfig(t *testing.T) {
	t.Run("should_create_worktree_with_auto_switch_enabled", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create initial commit
		err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0o644)
		require.NoError(t, err)

		cmd = exec.Command("git", "add", "README.md")
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "commit", "-m", "Initial commit")
		require.NoError(t, cmd.Run())

		// Test with auto-switch enabled
		config := &Config{AutoSwitch: true}
		result, err := AddWithConfig("../test-worktree", "feature/test", "HEAD", config)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "../test-worktree", result.Path)
		assert.Equal(t, "feature/test", result.Branch)
		assert.NotEmpty(t, result.Head)

		// Clean up
		cmd = exec.Command("git", "worktree", "remove", "../test-worktree")
		cmd.Run() // Ignore errors for cleanup
	})

	t.Run("should_create_worktree_with_auto_switch_disabled", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create initial commit
		err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0o644)
		require.NoError(t, err)

		cmd = exec.Command("git", "add", "README.md")
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "commit", "-m", "Initial commit")
		require.NoError(t, cmd.Run())

		// Test with auto-switch disabled
		config := &Config{AutoSwitch: false}
		result, err := AddWithConfig("../test-worktree", "feature/test", "HEAD", config)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "../test-worktree", result.Path)
		assert.Equal(t, "feature/test", result.Branch)
		assert.NotEmpty(t, result.Head)

		// Clean up
		cmd = exec.Command("git", "worktree", "remove", "../test-worktree")
		cmd.Run() // Ignore errors for cleanup
	})

	t.Run("should_create_worktree_with_nil_config", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create initial commit
		err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0o644)
		require.NoError(t, err)

		cmd = exec.Command("git", "add", "README.md")
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "commit", "-m", "Initial commit")
		require.NoError(t, cmd.Run())

		// Test with nil config (should behave like auto-switch disabled)
		result, err := AddWithConfig("../test-worktree", "feature/test", "HEAD", nil)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "../test-worktree", result.Path)
		assert.Equal(t, "feature/test", result.Branch)
		assert.NotEmpty(t, result.Head)

		// Clean up
		cmd = exec.Command("git", "worktree", "remove", "../test-worktree")
		cmd.Run() // Ignore errors for cleanup
	})
}
