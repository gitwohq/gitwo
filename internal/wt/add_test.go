package wt

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		branch     string
		startPoint string
		setup      func() string
		wantErr    bool
	}{
		{
			name:       "should create worktree successfully",
			path:       "../test-worktree",
			branch:     "feature/test",
			startPoint: "HEAD",
			setup: func() string {
				// Create a temporary git repo with a commit
				tmpDir := t.TempDir()
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				// Create a file and commit it
				require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644))
				cmd = exec.Command("git", "add", "test.txt")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())

				cmd = exec.Command("git", "commit", "-m", "Initial commit")
				cmd.Dir = tmpDir
				cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com")
				require.NoError(t, cmd.Run())

				return tmpDir
			},
			wantErr: false,
		},
		{
			name:       "should fail when not in git repo",
			path:       "../test-worktree",
			branch:     "feature/test",
			startPoint: "HEAD",
			setup: func() string {
				return t.TempDir()
			},
			wantErr: true,
		},
		{
			name:       "should fail with empty branch name",
			path:       "../test-worktree",
			branch:     "",
			startPoint: "HEAD",
			setup: func() string {
				tmpDir := t.TempDir()
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())
				return tmpDir
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

			// Clean up any existing worktree
			if !tt.wantErr {
				defer func() {
					// Try to remove the worktree if it was created
					worktreePath := filepath.Join(repoPath, "..", filepath.Base(tt.path))
					if _, err := os.Stat(worktreePath); err == nil {
						cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
						cmd.Dir = repoPath
						cmd.Run() // Ignore errors during cleanup
					}
				}()
			}

			_, err = Add(tt.path, tt.branch, tt.startPoint)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify worktree was created
				worktreePath := filepath.Join(repoPath, "..", filepath.Base(tt.path))
				assert.DirExists(t, worktreePath)
				
				// Verify branch was created
				cmd := exec.Command("git", "branch", "--list", tt.branch)
				cmd.Dir = repoPath
				output, err := cmd.Output()
				require.NoError(t, err)
				assert.Contains(t, string(output), tt.branch)
			}
		})
	}
}

func TestAdd_Validation(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		branch     string
		startPoint string
		expected   string
	}{
		{
			name:       "should validate path is not empty",
			path:       "",
			branch:     "test",
			startPoint: "HEAD",
			expected:   "path cannot be empty",
		},
		{
			name:       "should validate branch is not empty",
			path:       "../test",
			branch:     "",
			startPoint: "HEAD",
			expected:   "branch cannot be empty",
		},
		{
			name:       "should validate start point is not empty",
			path:       "../test",
			branch:     "test",
			startPoint: "",
			expected:   "start point cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary git repo for validation tests
			tmpDir, err := os.MkdirTemp("", "gitwo-validation-test")
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

			_, err = Add(tt.path, tt.branch, tt.startPoint)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}
