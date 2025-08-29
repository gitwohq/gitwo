package wt

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "should list worktrees successfully",
			setup: func() string {
				// Create a temporary git repo with worktrees
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
		},
		{
			name: "should return error when not in git repo",
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

			// Clean up worktrees after test
			if !tt.wantErr {
				defer func() {
					worktreePath := filepath.Join(repoPath, "..", "test-worktree")
					if _, err := os.Stat(worktreePath); err == nil {
						cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
						cmd.Dir = repoPath
						cmd.Run() // Ignore errors during cleanup
					}
				}()
			}

			items, err := List()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, items)
				assert.Len(t, items, 2) // Main worktree + created worktree
				
				// Verify main worktree is listed
				foundMain := false
				for _, item := range items {
					// Use contains check to handle path normalization differences
					if strings.Contains(item.Path, "TestList") && item.Branch == "main" {
						foundMain = true
						assert.NotEmpty(t, item.Head)
						break
					}
				}
				assert.True(t, foundMain, "Main worktree should be listed")
			}
		})
	}
}

func TestWorktreeItem_String(t *testing.T) {
	item := WorktreeItem{
		Path:   "/path/to/worktree",
		Head:   "abc123",
		Branch: "feature/test",
	}

	expected := "WorktreeItem{Path: /path/to/worktree, Branch: feature/test, Head: abc123}"
	assert.Equal(t, expected, item.String())
}

func TestList_EmptyRepo(t *testing.T) {
	// Test listing in a repo with no worktrees
	tmpDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	require.NoError(t, cmd.Run())

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	require.NoError(t, os.Chdir(tmpDir))

	items, err := List()
	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.Len(t, items, 1) // Only main worktree
	// Use contains check instead of exact match to handle path differences
	assert.Contains(t, items[0].Path, "TestList_EmptyRepo")
}
