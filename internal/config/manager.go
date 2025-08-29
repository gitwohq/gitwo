package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the gitwo configuration
type Config struct {
	// Core settings
	WorktreesDir        string `yaml:"worktrees_dir"`
	NameTemplate        string `yaml:"name_template"`
	MainBranch          string `yaml:"main_branch"`
	EditorCmd           string `yaml:"editor_cmd"`
	PostAddOpenEditor   bool   `yaml:"post_add_open_editor"`
	AutoSwitch          bool   `yaml:"auto_switch"`
	DefaultBranchPrefix string `yaml:"default_branch_prefix"`

	// Language/Framework detection
	Language  string `yaml:"language"`
	Framework string `yaml:"framework"`

	// Shell configuration
	Shell ShellConfig `yaml:"shell"`

	// Hook configuration
	Hooks HooksConfig `yaml:"hooks"`
}

// ShellConfig represents shell-specific configuration
type ShellConfig struct {
	Type             string `yaml:"type"`
	AutoCD           bool   `yaml:"auto_cd"`
	WrapperInstalled bool   `yaml:"wrapper_installed"`
}

// HooksConfig represents hook configuration
type HooksConfig struct {
	Enabled bool   `yaml:"enabled"`
	PreAdd  []Hook `yaml:"pre_add"`
	PostAdd []Hook `yaml:"post_add"`
}

// Hook represents a single hook
type Hook struct {
	Type        string `yaml:"type"`
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
	Language    string `yaml:"language,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WorktreesDir:        "..",
		NameTemplate:        "${REPO}-${BRANCH}",
		MainBranch:          "origin/main",
		EditorCmd:           "code -g",
		PostAddOpenEditor:   true,
		AutoSwitch:          true,
		DefaultBranchPrefix: "feature/",
		Language:            "unknown",
		Framework:           "unknown",
		Shell: ShellConfig{
			Type:             "unknown",
			AutoCD:           true,
			WrapperInstalled: false,
		},
		Hooks: HooksConfig{
			Enabled: true,
			PreAdd:  []Hook{},
			PostAdd: []Hook{},
		},
	}
}

// LoadConfig loads configuration from .gitwo/config.yml
func LoadConfig(repoPath string) (*Config, error) {
	configPath := filepath.Join(repoPath, ".gitwo", "config.yml")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config: %w", err)
	}

	// Merge with defaults for missing fields
	config = mergeWithDefaults(config)

	return &config, nil
}

// SaveConfig saves configuration to .gitwo/config.yml
func SaveConfig(repoPath string, config *Config) error {
	gitwoDir := filepath.Join(repoPath, ".gitwo")

	// Create .gitwo directory if it doesn't exist
	if err := os.MkdirAll(gitwoDir, 0o755); err != nil {
		return fmt.Errorf("failed to create .gitwo directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configPath := filepath.Join(gitwoDir, "config.yml")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// mergeWithDefaults merges config with default values for missing fields
func mergeWithDefaults(config Config) Config {
	defaults := DefaultConfig()

	if config.WorktreesDir == "" {
		config.WorktreesDir = defaults.WorktreesDir
	}
	if config.NameTemplate == "" {
		config.NameTemplate = defaults.NameTemplate
	}
	if config.MainBranch == "" {
		config.MainBranch = defaults.MainBranch
	}
	if config.EditorCmd == "" {
		config.EditorCmd = defaults.EditorCmd
	}
	if config.DefaultBranchPrefix == "" {
		config.DefaultBranchPrefix = defaults.DefaultBranchPrefix
	}
	if config.Language == "" {
		config.Language = defaults.Language
	}
	if config.Framework == "" {
		config.Framework = defaults.Framework
	}

	// Merge shell config
	if config.Shell.Type == "" {
		config.Shell.Type = defaults.Shell.Type
	}

	// Merge hooks config
	if !config.Hooks.Enabled {
		config.Hooks.Enabled = defaults.Hooks.Enabled
	}

	return config
}

// GetEnvConfig returns configuration from environment variables
func GetEnvConfig() map[string]string {
	envConfig := make(map[string]string)

	// Core environment variables
	if editor := os.Getenv("GITWO_EDITOR"); editor != "" {
		envConfig["editor_cmd"] = editor
	}
	if open := os.Getenv("GITWO_OPEN"); open != "" {
		envConfig["post_add_open_editor"] = open
	}
	if sync := os.Getenv("GITWO_SYNC_ON_NEW"); sync != "" {
		envConfig["sync_on_new"] = sync
	}
	if docker := os.Getenv("GITWO_DOCKER_UP"); docker != "" {
		envConfig["docker_up"] = docker
	}
	if verbose := os.Getenv("GITWO_VERBOSE"); verbose != "" {
		envConfig["verbose"] = verbose
	}
	if color := os.Getenv("GITWO_COLOR"); color != "" {
		envConfig["color"] = color
	}
	if timeout := os.Getenv("GITWO_TIMEOUT"); timeout != "" {
		envConfig["timeout"] = timeout
	}

	return envConfig
}

// SaveHooks saves hooks configuration to .gitwo/hooks/<hookType>.yml
func SaveHooks(repoPath, hookType string, hooks []Hook) error {
	hooksDir := filepath.Join(repoPath, ".gitwo", "hooks")

	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(hooks)
	if err != nil {
		return fmt.Errorf("failed to marshal hooks: %w", err)
	}

	// Write to file
	hookPath := filepath.Join(hooksDir, hookType+".yml")
	if err := os.WriteFile(hookPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write hooks: %w", err)
	}

	return nil
}
