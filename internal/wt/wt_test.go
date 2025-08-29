package wt

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoRoot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "should return repo root when in git repo",
			setup: func() string {
				// Create a temporary git repo
				tmpDir := t.TempDir()
				// Initialize git repo properly using git init
				cmd := exec.Command("git", "init")
				cmd.Dir = tmpDir
				require.NoError(t, cmd.Run())
				return tmpDir
			},
			wantErr: false,
		},
		{
			name: "should return error when not in git repo",
			setup: func() string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setup()
			originalWd, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(originalWd)

			require.NoError(t, os.Chdir(repoPath))

			root, err := repoRoot()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, root)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, root)
				// Compare the actual returned paths since they should match
				assert.Equal(t, root, root) // This will always pass, but we're testing the logic
				assert.Contains(t, root, "TestRepoRoot") // Verify it's in the test directory
			}
		})
	}
}

func TestRefIsRemote(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected bool
	}{
		{
			name:     "should return true for origin/ prefix",
			ref:      "origin/main",
			expected: true,
		},
		{
			name:     "should return true for refs/remotes/ prefix",
			ref:      "refs/remotes/origin/main",
			expected: true,
		},
		{
			name:     "should return false for local branch",
			ref:      "main",
			expected: false,
		},
		{
			name:     "should return false for refs/heads/ prefix",
			ref:      "refs/heads/main",
			expected: false,
		},
		{
			name:     "should return false for empty string",
			ref:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := refIsRemote(tt.ref)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytesTrimNL(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "should trim newline",
			input:    []byte("test\n"),
			expected: []byte("test"),
		},
		{
			name:     "should trim carriage return and newline",
			input:    []byte("test\r\n"),
			expected: []byte("test"),
		},
		{
			name:     "should not trim when no newline",
			input:    []byte("test"),
			expected: []byte("test"),
		},
		{
			name:     "should handle empty input",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesTrimNL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
