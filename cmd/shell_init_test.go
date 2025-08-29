package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/gitwohq/gitwo/internal/shell"
	"github.com/stretchr/testify/assert"
)

func TestShellInitCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		shellEnv string
		expected string
	}{
		{
			name:     "generates bash wrapper",
			args:     []string{"--shell", "bash"},
			expected: "gitwo()",
		},
		{
			name:     "generates zsh wrapper",
			args:     []string{"--shell", "zsh"},
			expected: "gitwo()",
		},
		{
			name:     "generates fish wrapper",
			args:     []string{"--shell", "fish"},
			expected: "function gitwo",
		},
		{
			name:     "generates powershell wrapper",
			args:     []string{"--shell", "powershell"},
			expected: "function gitwo",
		},
		{
			name:     "auto-detects zsh",
			args:     []string{},
			shellEnv: "/bin/zsh",
			expected: "gitwo()",
		},
		{
			name:     "auto-detects bash",
			args:     []string{},
			shellEnv: "/bin/bash",
			expected: "gitwo()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.shellEnv != "" {
				os.Setenv("SHELL", tt.shellEnv)
				defer os.Unsetenv("SHELL")
			}

			// Test the shell wrapper generation directly
			var shellType string
			if len(tt.args) > 1 && tt.args[0] == "--shell" {
				shellType = tt.args[1]
			} else {
				shellType = shell.DetectShell()
			}

			wrapper := shell.GenerateWrapper(shellType)

			// Assertions
			assert.Contains(t, wrapper, tt.expected)
			// Check for appropriate GITWO_WRAPPER syntax based on shell
			if strings.Contains(shellType, "fish") {
				assert.Contains(t, wrapper, "GITWO_WRAPPER 1")
			} else if strings.Contains(shellType, "powershell") {
				assert.Contains(t, wrapper, "GITWO_WRAPPER = \"1\"")
			} else {
				assert.Contains(t, wrapper, "GITWO_WRAPPER=1")
			}
		})
	}
}

func TestShellInitInvalidShell(t *testing.T) {
	// Test with invalid shell - should default to bash
	wrapper := shell.GenerateWrapper("invalid")
	assert.Contains(t, wrapper, "gitwo()")
}

func TestShellInitWrapperContent(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		required []string
	}{
		{
			name:  "bash wrapper content",
			shell: "bash",
			required: []string{
				"gitwo()",
				"builtin cd",
				"GITWO_WRAPPER=1",
				"command gitwo",
			},
		},
		{
			name:  "fish wrapper content",
			shell: "fish",
			required: []string{
				"function gitwo",
				"cd",
				"GITWO_WRAPPER 1",
				"command gitwo",
			},
		},
		{
			name:  "powershell wrapper content",
			shell: "powershell",
			required: []string{
				"function gitwo",
				"Set-Location",
				"GITWO_WRAPPER = \"1\"",
				"& gitwo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := shell.GenerateWrapper(tt.shell)

			for _, required := range tt.required {
				assert.Contains(t, wrapper, required, "Wrapper for %s missing: %s", tt.shell, required)
			}
		})
	}
}
