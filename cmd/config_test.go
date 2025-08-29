package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCommand(t *testing.T) {
	t.Run("should_show_default_config", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Run config show command
		rootCmd.SetArgs([]string{"config", "show"})
		err = rootCmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("should_set_auto_switch_true", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Set auto-switch to true
		rootCmd.SetArgs([]string{"config", "auto-switch", "true"})
		err = rootCmd.Execute()
		assert.NoError(t, err)

		// Verify config file was created
		configPath := filepath.Join(tmpDir, ".gitwo", "config.yml")
		assert.FileExists(t, configPath)

		// Check config content
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "auto_switch: true")
	})

	t.Run("should_set_auto_switch_false", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Set auto-switch to false
		rootCmd.SetArgs([]string{"config", "auto-switch", "false"})
		err = rootCmd.Execute()
		assert.NoError(t, err)

		// Verify config file was created
		configPath := filepath.Join(tmpDir, ".gitwo", "config.yml")
		assert.FileExists(t, configPath)

		// Check config content
		content, err := os.ReadFile(configPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "auto_switch: false")
	})

	t.Run("should_accept_various_boolean_values", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		testCases := []struct {
			input    string
			expected string
		}{
			{"1", "auto_switch: true"},
			{"0", "auto_switch: false"},
			{"yes", "auto_switch: true"},
			{"no", "auto_switch: false"},
			{"on", "auto_switch: true"},
			{"off", "auto_switch: false"},
		}

		for _, tc := range testCases {
			t.Run("should_accept_"+tc.input, func(t *testing.T) {
				// Set auto-switch
				rootCmd.SetArgs([]string{"config", "auto-switch", tc.input})
				err = rootCmd.Execute()
				assert.NoError(t, err)

				// Check config content
				configPath := filepath.Join(tmpDir, ".gitwo", "config.yml")
				content, err := os.ReadFile(configPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), tc.expected)
			})
		}
	})

	t.Run("should_reject_invalid_boolean_value", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Try to set invalid value
		rootCmd.SetArgs([]string{"config", "auto-switch", "invalid"})
		err = rootCmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value")
	})

	t.Run("should_fail_when_not_in_git_repo", func(t *testing.T) {
		// Create a temporary directory outside of git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Change to temp directory (not a git repo)
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalDir)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		// Try to set config
		rootCmd.SetArgs([]string{"config", "auto-switch", "true"})
		err = rootCmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in a git repository")
	})
}
