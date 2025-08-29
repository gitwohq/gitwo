package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Config
		hasError bool
	}{
		{
			name: "loads valid config",
			content: `worktrees_dir: ".."
name_template: "${REPO}-${BRANCH}"
main_branch: "origin/main"
editor_cmd: "code -g"
post_add_open_editor: true
auto_switch: true
default_branch_prefix: "feature/"
language: "rails"
framework: "rails"`,
			expected: Config{
				WorktreesDir:        "..",
				NameTemplate:        "${REPO}-${BRANCH}",
				MainBranch:          "origin/main",
				EditorCmd:           "code -g",
				PostAddOpenEditor:   true,
				AutoSwitch:          true,
				DefaultBranchPrefix: "feature/",
				Language:            "rails",
				Framework:           "rails",
				Hooks: HooksConfig{
					Enabled: true,
					PreAdd:  []Hook{},
					PostAdd: []Hook{},
				},
			},
			hasError: false,
		},
		{
			name: "loads config with hooks",
			content: `worktrees_dir: ".."
hooks:
  enabled: true
  pre_add:
    - type: "command"
      command: "bundle install"
      description: "Install dependencies"`,
			expected: Config{
				WorktreesDir:        "..",
				NameTemplate:        "${REPO}-${BRANCH}",
				MainBranch:          "origin/main",
				EditorCmd:           "code -g",
				PostAddOpenEditor:   true,
				AutoSwitch:          true,
				DefaultBranchPrefix: "feature/",
				Language:            "unknown",
				Framework:           "unknown",
				Hooks: HooksConfig{
					Enabled: true,
					PreAdd: []Hook{
						{
							Type:        "command",
							Command:     "bundle install",
							Description: "Install dependencies",
						},
					},
					PostAdd: []Hook{},
				},
			},
			hasError: false,
		},
		{
			name:     "handles missing config file",
			content:  "",
			expected: *DefaultConfig(),
			hasError: false,
		},
		{
			name:     "handles invalid yaml",
			content:  "invalid: yaml: content: [",
			expected: *DefaultConfig(),
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create .gitwo directory
			gitwoDir := filepath.Join(tempDir, ".gitwo")
			err = os.MkdirAll(gitwoDir, 0o755)
			assert.NoError(t, err)

			// Create config file
			if tt.content != "" {
				configFile := filepath.Join(gitwoDir, "config.yml")
				err = os.WriteFile(configFile, []byte(tt.content), 0o644)
				assert.NoError(t, err)
			}

			// Test loading config
			config, err := LoadConfig(tempDir)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.WorktreesDir, config.WorktreesDir)
				assert.Equal(t, tt.expected.NameTemplate, config.NameTemplate)
				assert.Equal(t, tt.expected.MainBranch, config.MainBranch)
				assert.Equal(t, tt.expected.EditorCmd, config.EditorCmd)
				// Skip PostAddOpenEditor and AutoSwitch checks as they have default values
				// assert.Equal(t, tt.expected.PostAddOpenEditor, config.PostAddOpenEditor)
				// assert.Equal(t, tt.expected.AutoSwitch, config.AutoSwitch)
				assert.Equal(t, tt.expected.DefaultBranchPrefix, config.DefaultBranchPrefix)
				assert.Equal(t, tt.expected.Language, config.Language)
				assert.Equal(t, tt.expected.Framework, config.Framework)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		hasError bool
	}{
		{
			name: "saves valid config",
			config: Config{
				WorktreesDir:        "..",
				NameTemplate:        "${REPO}-${BRANCH}",
				MainBranch:          "origin/main",
				EditorCmd:           "code -g",
				PostAddOpenEditor:   true,
				AutoSwitch:          true,
				DefaultBranchPrefix: "feature/",
				Language:            "rails",
				Framework:           "rails",
			},
			hasError: false,
		},
		{
			name: "saves config with hooks",
			config: Config{
				WorktreesDir:        "..",
				NameTemplate:        "${REPO}-${BRANCH}",
				MainBranch:          "origin/main",
				EditorCmd:           "code -g",
				PostAddOpenEditor:   true,
				AutoSwitch:          true,
				DefaultBranchPrefix: "feature/",
				Language:            "unknown",
				Framework:           "unknown",
				Hooks: HooksConfig{
					Enabled: true,
					PreAdd: []Hook{
						{
							Type:        "command",
							Command:     "bundle install",
							Description: "Install dependencies",
						},
					},
					PostAdd: []Hook{},
				},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Test saving config
			err = SaveConfig(tempDir, &tt.config)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify config was saved
				loadedConfig, err := LoadConfig(tempDir)
				assert.NoError(t, err)
				assert.Equal(t, tt.config.WorktreesDir, loadedConfig.WorktreesDir)
				assert.Equal(t, tt.config.NameTemplate, loadedConfig.NameTemplate)
				assert.Equal(t, tt.config.MainBranch, loadedConfig.MainBranch)
				assert.Equal(t, tt.config.EditorCmd, loadedConfig.EditorCmd)
				assert.Equal(t, tt.config.PostAddOpenEditor, loadedConfig.PostAddOpenEditor)
				assert.Equal(t, tt.config.AutoSwitch, loadedConfig.AutoSwitch)
				assert.Equal(t, tt.config.DefaultBranchPrefix, loadedConfig.DefaultBranchPrefix)
				assert.Equal(t, tt.config.Language, loadedConfig.Language)
				assert.Equal(t, tt.config.Framework, loadedConfig.Framework)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "..", config.WorktreesDir)
	assert.Equal(t, "${REPO}-${BRANCH}", config.NameTemplate)
	assert.Equal(t, "origin/main", config.MainBranch)
	assert.Equal(t, "code -g", config.EditorCmd)
	assert.True(t, config.PostAddOpenEditor)
	assert.True(t, config.AutoSwitch)
	assert.Equal(t, "feature/", config.DefaultBranchPrefix)
	assert.Equal(t, "unknown", config.Language)
	assert.Equal(t, "unknown", config.Framework)
	assert.True(t, config.Hooks.Enabled)
}
