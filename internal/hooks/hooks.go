package hooks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Hook represents a single hook
type Hook struct {
	Type        string `yaml:"type"`
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
	Language    string `yaml:"language,omitempty"`
}

// LoadHooks loads hooks from .gitwo/hooks/<hookType>.yml
func LoadHooks(repoPath, hookType string) ([]Hook, error) {
	hookPath := filepath.Join(repoPath, ".gitwo", "hooks", hookType+".yml")

	// If hook file doesn't exist, return empty slice
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return []Hook{}, nil
	}

	// Read hook file
	data, err := os.ReadFile(hookPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hooks: %w", err)
	}

	// Parse YAML
	var hookFile struct {
		Hooks []Hook `yaml:"hooks"`
	}
	if err := yaml.Unmarshal(data, &hookFile); err != nil {
		return nil, fmt.Errorf("failed to parse hooks: %w", err)
	}

	// Validate hooks
	for _, hook := range hookFile.Hooks {
		if err := ValidateHook(hook); err != nil {
			return nil, fmt.Errorf("invalid hook: %w", err)
		}
	}

	return hookFile.Hooks, nil
}

// ExecuteHooks executes a list of hooks with the given environment variables
func ExecuteHooks(hooks []Hook, env map[string]string) error {
	for _, hook := range hooks {
		if err := ExecuteHook(hook, env); err != nil {
			return fmt.Errorf("hook execution failed: %w", err)
		}
	}
	return nil
}

// ExecuteHook executes a single hook
func ExecuteHook(hook Hook, env map[string]string) error {
	switch hook.Type {
	case "command":
		return executeCommandHook(hook, env)
	default:
		return fmt.Errorf("unknown hook type: %s", hook.Type)
	}
}

// executeCommandHook executes a command hook
func executeCommandHook(hook Hook, env map[string]string) error {
	// Get timeout from environment
	timeout := 30 * time.Second // default timeout
	if timeoutStr := env["GITWO_TIMEOUT"]; timeoutStr != "" {
		if t, err := time.ParseDuration(timeoutStr + "s"); err == nil {
			timeout = t
		}
	}

	// Create command
	cmd := exec.Command("sh", "-c", hook.Command)

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd = exec.CommandContext(ctx, "sh", "-c", hook.Command)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %s, output: %s", err, string(output))
	}

	return nil
}

// ValidateHook validates a hook configuration
func ValidateHook(hook Hook) error {
	if hook.Type == "" {
		return fmt.Errorf("hook type is required")
	}

	switch hook.Type {
	case "command":
		if hook.Command == "" {
			return fmt.Errorf("command is required for command hooks")
		}
	default:
		return fmt.Errorf("unknown hook type: %s", hook.Type)
	}

	return nil
}

// CreateHookEnvironment creates environment variables for hook execution
func CreateHookEnvironment(repoPath, branch, worktreePath string, config map[string]string) map[string]string {
	env := make(map[string]string)

	// Core gitwo variables
	env["GITWO_ACTION"] = "add"
	env["GITWO_REPO"] = filepath.Base(repoPath)
	env["GITWO_BRANCH"] = branch
	env["GITWO_PATH"] = worktreePath

	// Configuration variables
	if worktreesDir := config["worktrees_dir"]; worktreesDir != "" {
		env["GITWO_WORKTREES_DIR"] = worktreesDir
	}
	if mainBranch := config["main_branch"]; mainBranch != "" {
		env["GITWO_MAIN_BRANCH"] = mainBranch
	}
	if nameTemplate := config["name_template"]; nameTemplate != "" {
		env["GITWO_NAME_TEMPLATE"] = nameTemplate
	}
	if editor := config["editor_cmd"]; editor != "" {
		env["GITWO_EDITOR"] = editor
	}

	// Environment overrides
	if editor := os.Getenv("GITWO_EDITOR"); editor != "" {
		env["GITWO_EDITOR"] = editor
	}
	if open := os.Getenv("GITWO_OPEN"); open != "" {
		env["GITWO_OPEN"] = open
	}
	if sync := os.Getenv("GITWO_SYNC_ON_NEW"); sync != "" {
		env["GITWO_SYNC_ON_NEW"] = sync
	}
	if docker := os.Getenv("GITWO_DOCKER_UP"); docker != "" {
		env["GITWO_DOCKER_UP"] = docker
	}
	if verbose := os.Getenv("GITWO_VERBOSE"); verbose != "" {
		env["GITWO_VERBOSE"] = verbose
	}
	if color := os.Getenv("GITWO_COLOR"); color != "" {
		env["GITWO_COLOR"] = color
	}
	if timeout := os.Getenv("GITWO_TIMEOUT"); timeout != "" {
		env["GITWO_TIMEOUT"] = timeout
	}

	return env
}
