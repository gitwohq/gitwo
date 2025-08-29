package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadHooks(t *testing.T) {
	tests := []struct {
		name     string
		hookType string
		content  string
		expected int
	}{
		{
			name:     "loads pre_add hooks",
			hookType: "pre_add",
			content: `hooks:
  - type: "command"
    command: "bundle install"
    description: "Install dependencies"`,
			expected: 1,
		},
		{
			name:     "loads post_add hooks",
			hookType: "post_add",
			content: `hooks:
  - type: "command"
    command: "bin/setup"
    description: "Run setup script"`,
			expected: 1,
		},
		{
			name:     "loads multiple hooks",
			hookType: "post_add",
			content: `hooks:
  - type: "command"
    command: "bundle install"
    description: "Install dependencies"
  - type: "command"
    command: "bin/setup"
    description: "Run setup script"`,
			expected: 2,
		},
		{
			name:     "handles empty hooks file",
			hookType: "pre_add",
			content:  "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create .gitwo/hooks directory
			gitwoDir := filepath.Join(tempDir, ".gitwo")
			hooksDir := filepath.Join(gitwoDir, "hooks")
			err = os.MkdirAll(hooksDir, 0o755)
			assert.NoError(t, err)

			// Create hooks file
			if tt.content != "" {
				hookFile := filepath.Join(hooksDir, tt.hookType+".yml")
				err = os.WriteFile(hookFile, []byte(tt.content), 0o644)
				assert.NoError(t, err)
			}

			// Test loading hooks
			hooks, err := LoadHooks(tempDir, tt.hookType)
			assert.NoError(t, err)
			assert.Len(t, hooks, tt.expected)
		})
	}
}

func TestExecuteHooks(t *testing.T) {
	tests := []struct {
		name     string
		hooks    []Hook
		env      map[string]string
		expected bool
	}{
		{
			name: "executes command hook successfully",
			hooks: []Hook{
				{
					Type:        "command",
					Command:     "echo 'test'",
					Description: "Test command",
				},
			},
			expected: true,
		},
		{
			name: "handles multiple hooks",
			hooks: []Hook{
				{
					Type:        "command",
					Command:     "echo 'hook1'",
					Description: "First hook",
				},
				{
					Type:        "command",
					Command:     "echo 'hook2'",
					Description: "Second hook",
				},
			},
			expected: true,
		},
		{
			name:     "handles empty hooks",
			hooks:    []Hook{},
			expected: true,
		},
		{
			name: "sets environment variables",
			hooks: []Hook{
				{
					Type:        "command",
					Command:     "echo $GITWO_REPO",
					Description: "Test env var",
				},
			},
			env: map[string]string{
				"GITWO_REPO": "test-repo",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.env {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Test hook execution
			err := ExecuteHooks(tt.hooks, tt.env)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestHookValidation(t *testing.T) {
	tests := []struct {
		name    string
		hook    Hook
		isValid bool
	}{
		{
			name: "valid command hook",
			hook: Hook{
				Type:        "command",
				Command:     "bundle install",
				Description: "Install dependencies",
			},
			isValid: true,
		},
		{
			name: "invalid hook type",
			hook: Hook{
				Type:        "invalid",
				Command:     "bundle install",
				Description: "Install dependencies",
			},
			isValid: false,
		},
		{
			name: "missing command",
			hook: Hook{
				Type:        "command",
				Description: "Install dependencies",
			},
			isValid: false,
		},
		{
			name: "empty command",
			hook: Hook{
				Type:        "command",
				Command:     "",
				Description: "Install dependencies",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHook(tt.hook)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
