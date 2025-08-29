package wt

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		config := DefaultConfig()
		assert.NotNil(t, config)
		assert.True(t, config.AutoSwitch, "AutoSwitch should be true by default")
	})

	t.Run("LoadConfig_WhenNoConfigFile", func(t *testing.T) {
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

		config, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.True(t, config.AutoSwitch, "Should return default config when no config file exists")
	})

	t.Run("LoadConfig_WithValidConfig", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create .gitwo directory and config file
		gitwoDir := filepath.Join(tmpDir, ".gitwo")
		err = os.MkdirAll(gitwoDir, 0o755)
		require.NoError(t, err)

		configContent := `auto_switch: true`
		configPath := filepath.Join(gitwoDir, "config.yml")
		err = os.WriteFile(configPath, []byte(configContent), 0o644)
		require.NoError(t, err)

		config, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.True(t, config.AutoSwitch, "Should load auto_switch: true from config file")
	})

	t.Run("LoadConfig_WithInvalidYAML", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create .gitwo directory and invalid config file
		gitwoDir := filepath.Join(tmpDir, ".gitwo")
		err = os.MkdirAll(gitwoDir, 0o755)
		require.NoError(t, err)

		invalidConfig := `auto_switch: invalid_value:`
		configPath := filepath.Join(gitwoDir, "config.yml")
		err = os.WriteFile(configPath, []byte(invalidConfig), 0o644)
		require.NoError(t, err)

		config, err := LoadConfig()
		assert.Error(t, err) // Should error with invalid YAML
		assert.NotNil(t, config) // But should still return default config
		assert.True(t, config.AutoSwitch, "Should return default config when YAML is invalid")
	})

	t.Run("SaveConfig", func(t *testing.T) {
		// Create a temporary git repo
		tmpDir, err := os.MkdirTemp("", "gitwo-config-test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Initialize git repo
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cmd := exec.Command("git", "init")
		require.NoError(t, cmd.Run())

		// Create config with auto_switch enabled
		config := &Config{AutoSwitch: true}

		// Save config
		err = SaveConfig(config)
		assert.NoError(t, err)

		// Verify config file was created
		configPath := filepath.Join(tmpDir, ".gitwo", "config.yml")
		assert.FileExists(t, configPath)

		// Load config and verify
		loadedConfig, err := LoadConfig()
		assert.NoError(t, err)
		assert.True(t, loadedConfig.AutoSwitch, "Saved config should have auto_switch: true")
	})

	t.Run("SaveConfig_WhenNotInGitRepo", func(t *testing.T) {
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

		config := &Config{AutoSwitch: true}
		err = SaveConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in a git repository")
	})
}
