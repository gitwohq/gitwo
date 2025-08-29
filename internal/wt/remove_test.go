package wt

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		name       string
		worktree   string
		setup      func() string
		wantErr    bool
		checkSetup func(t *testing.T, repoPath string)
	}{
		{
			name:     "should remove worktree by full path",
			worktree: "../test-worktree",
			setup: func() string {
				// Create a temporary git repo with worktree
				tmpDir := t.TempDir()
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				// Create initial commit
				require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644))
				cmd = exec.Command("git", "add", "test.txt")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				cmd = exec.Command("git", "commit", "-m", "Initial commit")
				cmd.Dir = tmpDir
				cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com")
				require.NoError(t, cmd.Run())

				// Create a worktree
				worktreePath := filepath.Join(tmpDir, "..", "test-worktree")
				cmd = exec.Command("git", "worktree", "add", "-b", "feature/test", worktreePath, "HEAD")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				return tmpDir
			},
			wantErr: false,
			checkSetup: func(t *testing.T, repoPath string) {
				// Verify worktree was removed
				worktreePath := filepath.Join(repoPath, "..", "test-worktree")
				assert.NoDirExists(t, worktreePath)

				// Verify branch still exists (worktree removal doesn't delete branch)
				cmd := exec.Command("git", "branch", "--list", "feature/test")
				cmd.Dir = repoPath
				output, err := cmd.Output()
				require.NoError(t, err)
				assert.Contains(t, string(output), "feature/test")
			},
		},
		{
			name:     "should remove worktree by name only",
			worktree: "test-worktree",
			setup: func() string {
				// Create a temporary git repo with worktree
				tmpDir := t.TempDir()
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				// Create initial commit
				require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644))
				cmd = exec.Command("git", "add", "test.txt")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				cmd = exec.Command("git", "commit", "-m", "Initial commit")
				cmd.Dir = tmpDir
				cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com")
				require.NoError(t, cmd.Run())

				// Create a worktree
				worktreePath := filepath.Join(tmpDir, "..", "test-worktree")
				cmd = exec.Command("git", "worktree", "add", "-b", "feature/test", worktreePath, "HEAD")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				return tmpDir
			},
			wantErr: false,
			checkSetup: func(t *testing.T, repoPath string) {
				// Verify worktree was removed
				worktreePath := filepath.Join(repoPath, "..", "test-worktree")
				assert.NoDirExists(t, worktreePath)
			},
		},
		{
			name:     "should fail when worktree doesn't exist",
			worktree: "non-existent-worktree",
			setup: func() string {
				tmpDir := t.TempDir()
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())
				return tmpDir
			},
			wantErr: true,
		},
		{
			name:     "should fail when not in git repo",
			worktree: "test-worktree",
			setup: func() string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setup()
			originalWd, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalWd)

			require.NoError(t, os.Chdir(repoPath))

			err = Remove(tt.worktree)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkSetup != nil {
					tt.checkSetup(t, repoPath)
				}
			}
		})
	}
}

func TestRemove_Validation(t *testing.T) {
	tests := []struct {
		name     string
		worktree string
		expected string
	}{
		{
			name:     "should validate worktree is not empty",
			worktree: "",
			expected: "worktree cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Remove(tt.worktree)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}
