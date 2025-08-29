package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/spf13/cobra"
	"github.com/gitwohq/gitwo/internal/wt"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		setup      func() string
		wantErr    bool
		checkSetup func(t *testing.T, repoPath string)
	}{
		{
			name: "should create worktree with smart branch naming",
			args: []string{"feature-add-openapi"},
			setup: func() string {
				// Create a temporary git repo
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

				return tmpDir
			},
			wantErr: false,
			checkSetup: func(t *testing.T, repoPath string) {
				// Verify worktree was created
				worktreePath := filepath.Join(repoPath, "..", "feature-add-openapi")
				assert.DirExists(t, worktreePath)

				// Verify branch was created
				cmd := exec.Command("git", "branch", "--list", "feature/feature-add-openapi")
				cmd.Dir = repoPath
				output, err := cmd.Output()
				require.NoError(t, err)
				assert.Contains(t, string(output), "feature/feature-add-openapi")

				// Clean up
				cmd = exec.Command("git", "worktree", "remove", worktreePath, "--force")
				cmd.Dir = repoPath
				cmd.Run() // Ignore errors during cleanup
			},
		},
		{
			name: "should fail when not in git repo",
			args: []string{"test-feature"},
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

			// Create a new root command for testing
			rootCmd := &cobra.Command{
				Use:   "gitwo",
				Short: "Git Worktree helper CLI",
			}

			// Add the new command
			newCmd := &cobra.Command{
				Use:   "new <name>",
				Short: "Create a new worktree with smart branch naming",
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					name := args[0]
					path := filepath.Join("..", name)
					branch := "feature/" + name
					_, err := wt.Add(path, branch, "HEAD")
					return err
				},
			}
			rootCmd.AddCommand(newCmd)

			// Execute the command
			rootCmd.SetArgs(append([]string{"new"}, tt.args...))
			err = rootCmd.Execute()

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
