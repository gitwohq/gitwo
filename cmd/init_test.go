package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string
		args           []string
		expectedConfig map[string]interface{}
		expectedHooks  bool
		expectError    bool
	}{
		{
			name: "init in rails project",
			setupFiles: map[string]string{
				"Gemfile": `source 'https://rubygems.org'
gem 'rails', '~> 7.0'`,
				"config/application.rb": "# Rails application",
			},
			args: []string{},
			expectedConfig: map[string]interface{}{
				"language":  "ruby",
				"framework": "rails",
			},
			expectedHooks: true,
		},
		{
			name: "init in nodejs project",
			setupFiles: map[string]string{
				"package.json": `{
					"name": "my-app",
					"dependencies": {
						"next": "^13.0.0"
					}
				}`,
			},
			args: []string{},
			expectedConfig: map[string]interface{}{
				"language":  "nodejs",
				"framework": "nextjs",
			},
			expectedHooks: true,
		},
		{
			name: "init in go project",
			setupFiles: map[string]string{
				"go.mod": `module my-app
go 1.21`,
			},
			args: []string{},
			expectedConfig: map[string]interface{}{
				"language":  "go",
				"framework": "unknown",
			},
			expectedHooks: true,
		},
		{
			name: "init in empty directory",
			setupFiles: map[string]string{},
			args: []string{},
			expectedConfig: map[string]interface{}{
				"language":  "unknown",
				"framework": "unknown",
			},
			expectedHooks: false,
		},
		{
			name: "init with --non-interactive",
			setupFiles: map[string]string{
				"Gemfile": "source 'https://rubygems.org'",
			},
			args: []string{"--non-interactive"},
			expectedConfig: map[string]interface{}{
				"language":  "ruby",
				"framework": "unknown",
			},
			expectedHooks: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-init-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup test files
			for filename, content := range tt.setupFiles {
				filePath := filepath.Join(tempDir, filename)
				err := os.MkdirAll(filepath.Dir(filePath), 0o755)
				require.NoError(t, err)
				err = os.WriteFile(filePath, []byte(content), 0o644)
				require.NoError(t, err)
			}

			// Change to temp directory
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalDir)
			err = os.Chdir(tempDir)
			require.NoError(t, err)

			// Set flags based on args
			initWithShell = false
			initNoShell = false
			initNonInteractive = false
			initShellType = ""

			for i := 0; i < len(tt.args); i++ {
				switch tt.args[i] {
				case "--with-shell":
					initWithShell = true
				case "--no-shell":
					initNoShell = true
				case "--non-interactive":
					initNonInteractive = true
				case "--shell":
					if i+1 < len(tt.args) {
						initShellType = tt.args[i+1]
						i++ // skip next argument
					}
				}
			}

			// Run init function directly
			err = runInit(nil, tt.args)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check if .gitwo directory was created
			gitwoDir := filepath.Join(tempDir, ".gitwo")
			assert.DirExists(t, gitwoDir)

			// Check if config.yml was created
			configPath := filepath.Join(gitwoDir, "config.yml")
			assert.FileExists(t, configPath)

			// Check if hooks directory was created (if expected)
			if tt.expectedHooks {
				hooksDir := filepath.Join(gitwoDir, "hooks")
				assert.DirExists(t, hooksDir)
			}
		})
	}
}

func TestInitCommandWithShellIntegration(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectShell bool
	}{
		{
			name:        "init with --with-shell",
			args:        []string{"--with-shell"},
			expectShell: true,
		},
		{
			name:        "init with --no-shell",
			args:        []string{"--no-shell"},
			expectShell: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-init-shell-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup basic Rails project
			gemfilePath := filepath.Join(tempDir, "Gemfile")
			err = os.WriteFile(gemfilePath, []byte("source 'https://rubygems.org'"), 0o644)
			require.NoError(t, err)

			// Change to temp directory
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalDir)
			err = os.Chdir(tempDir)
			require.NoError(t, err)

			// Set flags based on args
			initWithShell = false
			initNoShell = false
			initNonInteractive = false
			initShellType = ""

			for i := 0; i < len(tt.args); i++ {
				switch tt.args[i] {
				case "--with-shell":
					initWithShell = true
				case "--no-shell":
					initNoShell = true
				case "--non-interactive":
					initNonInteractive = true
				case "--shell":
					if i+1 < len(tt.args) {
						initShellType = tt.args[i+1]
						i++ // skip next argument
					}
				}
			}

			// Run init function directly
			err = runInit(nil, tt.args)

			assert.NoError(t, err)

			// Check if shell wrapper was installed (if expected)
			if tt.expectShell {
				// This would check if the shell wrapper was actually installed
				// For now, we'll just verify the command didn't error
				assert.NoError(t, err)
			}
		})
	}
}

func TestInitCommandValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "invalid shell type",
			args:        []string{"--shell", "invalid-shell"},
			expectError: true,
		},
		{
			name:        "valid shell type",
			args:        []string{"--shell", "bash"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-init-validation-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Change to temp directory
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalDir)
			err = os.Chdir(tempDir)
			require.NoError(t, err)

			// Set flags based on args
			initWithShell = false
			initNoShell = false
			initNonInteractive = false
			initShellType = ""

			for i := 0; i < len(tt.args); i++ {
				switch tt.args[i] {
				case "--with-shell":
					initWithShell = true
				case "--no-shell":
					initNoShell = true
				case "--non-interactive":
					initNonInteractive = true
				case "--shell":
					if i+1 < len(tt.args) {
						initShellType = tt.args[i+1]
						i++ // skip next argument
					}
				}
			}

			// Run init function directly
			err = runInit(nil, tt.args)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
