package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBashZshWrapper(t *testing.T) {
	wrapper := GenerateBashZshWrapper()

	assert.Contains(t, wrapper, "gitwo()")
	assert.Contains(t, wrapper, "builtin cd")
	assert.Contains(t, wrapper, "GITWO_WRAPPER=1")
	assert.Contains(t, wrapper, "command gitwo")
}

func TestGenerateFishWrapper(t *testing.T) {
	wrapper := GenerateFishWrapper()

	assert.Contains(t, wrapper, "function gitwo")
	assert.Contains(t, wrapper, "cd")
	assert.Contains(t, wrapper, "GITWO_WRAPPER 1")
	assert.Contains(t, wrapper, "command gitwo")
}

func TestGeneratePowerShellWrapper(t *testing.T) {
	wrapper := GeneratePowerShellWrapper()

	assert.Contains(t, wrapper, "function gitwo")
	assert.Contains(t, wrapper, "Set-Location")
	assert.Contains(t, wrapper, "GITWO_WRAPPER = \"1\"")
	assert.Contains(t, wrapper, "& gitwo")
}

func TestWrapperSanity(t *testing.T) {
	tests := []struct {
		name    string
		shell   string
		checker func(string) bool
	}{
		{
			name:  "bash wrapper contains required elements",
			shell: "bash",
			checker: func(wrapper string) bool {
				return strings.Contains(wrapper, "gitwo()") &&
					strings.Contains(wrapper, "builtin cd") &&
					strings.Contains(wrapper, "GITWO_WRAPPER=1")
			},
		},
		{
			name:  "zsh wrapper contains required elements",
			shell: "zsh",
			checker: func(wrapper string) bool {
				return strings.Contains(wrapper, "gitwo()") &&
					strings.Contains(wrapper, "builtin cd") &&
					strings.Contains(wrapper, "GITWO_WRAPPER=1")
			},
		},
		{
			name:  "fish wrapper contains required elements",
			shell: "fish",
			checker: func(wrapper string) bool {
				return strings.Contains(wrapper, "function gitwo") &&
					strings.Contains(wrapper, "cd") &&
					strings.Contains(wrapper, "GITWO_WRAPPER 1")
			},
		},
		{
			name:  "powershell wrapper contains required elements",
			shell: "powershell",
			checker: func(wrapper string) bool {
				return strings.Contains(wrapper, "function gitwo") &&
					strings.Contains(wrapper, "Set-Location") &&
					strings.Contains(wrapper, "GITWO_WRAPPER = \"1\"")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wrapper string
			switch tt.shell {
			case "bash", "zsh":
				wrapper = GenerateBashZshWrapper()
			case "fish":
				wrapper = GenerateFishWrapper()
			case "powershell":
				wrapper = GeneratePowerShellWrapper()
			}

			assert.True(t, tt.checker(wrapper), "Wrapper for %s missing required elements", tt.shell)
		})
	}
}
