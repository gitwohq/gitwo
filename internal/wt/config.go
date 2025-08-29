package wt

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the gitwo configuration
type Config struct {
	AutoSwitch bool `yaml:"auto_switch"` // Automatically switch to new worktree after creation
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AutoSwitch: true, // Default to true for convenience
	}
}

// LoadConfig loads configuration from .gitwo/config.yml
func LoadConfig() (*Config, error) {
	repoRoot, err := repoRoot()
	if err != nil {
		return DefaultConfig(), nil // Return default if not in a repo
	}

	configPath := filepath.Join(repoRoot, ".gitwo", "config.yml")

	// Check if config file exists
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

	return &config, nil
}

// SaveConfig saves configuration to .gitwo/config.yml
func SaveConfig(config *Config) error {
	repoRoot, err := repoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// Create .gitwo directory if it doesn't exist
	gitwoDir := filepath.Join(repoRoot, ".gitwo")
	if err := os.MkdirAll(gitwoDir, 0o755); err != nil {
		return fmt.Errorf("failed to create .gitwo directory: %w", err)
	}

	configPath := filepath.Join(gitwoDir, "config.yml")

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
