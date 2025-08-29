package shell

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		psModule string
		expected string
	}{
		{
			name:     "detects zsh",
			shellEnv: "/bin/zsh",
			expected: "zsh",
		},
		{
			name:     "detects bash",
			shellEnv: "/bin/bash",
			expected: "bash",
		},
		{
			name:     "detects fish",
			shellEnv: "/usr/bin/fish",
			expected: "fish",
		},
		{
			name:     "detects powershell on windows",
			shellEnv: "",
			psModule: "C:\\Program Files\\PowerShell\\Modules",
			expected: "powershell",
		},
		{
			name:     "defaults to bash when unknown",
			shellEnv: "/bin/unknown",
			expected: "bash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.shellEnv != "" {
				os.Setenv("SHELL", tt.shellEnv)
			} else {
				os.Unsetenv("SHELL")
			}
			if tt.psModule != "" {
				os.Setenv("PSModulePath", tt.psModule)
			} else {
				os.Unsetenv("PSModulePath")
			}

			result := DetectShell()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetShellProfilePath(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		expected string
	}{
		{
			name:     "zsh profile path",
			shell:    "zsh",
			expected: ".zshrc",
		},
		{
			name:     "bash profile path",
			shell:    "bash",
			expected: ".bashrc",
		},
		{
			name:     "fish profile path",
			shell:    "fish",
			expected: ".config/fish/config.fish",
		},
		{
			name:     "powershell profile path",
			shell:    "powershell",
			expected: "Documents/PowerShell/Microsoft.PowerShell_profile.ps1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetShellProfilePath(tt.shell)
			assert.Contains(t, result, tt.expected)
		})
	}
}
