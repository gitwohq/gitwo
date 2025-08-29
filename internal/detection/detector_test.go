package detection

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected string
	}{
		{
			name: "detects ruby",
			files: map[string]string{
				"Gemfile": "source 'https://rubygems.org'",
			},
			expected: "ruby",
		},
		{
			name: "detects rails",
			files: map[string]string{
				"Gemfile":               "source 'https://rubygems.org'\ngem 'rails'",
				"config/application.rb": "class Application < Rails::Application",
			},
			expected: "ruby",
		},
		{
			name: "detects nodejs",
			files: map[string]string{
				"package.json": `{"name": "test", "version": "1.0.0"}`,
			},
			expected: "nodejs",
		},
		{
			name: "detects nextjs",
			files: map[string]string{
				"package.json": `{"name": "test", "dependencies": {"next": "^13.0.0"}}`,
			},
			expected: "nodejs",
		},
		{
			name: "detects go",
			files: map[string]string{
				"go.mod": "module test",
			},
			expected: "go",
		},
		{
			name: "detects java",
			files: map[string]string{
				"pom.xml": "<project><groupId>test</groupId></project>",
			},
			expected: "java",
		},
		{
			name:     "defaults to unknown",
			files:    map[string]string{},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create test files
			for filename, content := range tt.files {
				filePath := filepath.Join(tempDir, filename)
				err := os.MkdirAll(filepath.Dir(filePath), 0o755)
				assert.NoError(t, err)
				err = os.WriteFile(filePath, []byte(content), 0o644)
				assert.NoError(t, err)
			}

			// Test detection
			result := DetectLanguage(tempDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectFramework(t *testing.T) {
	tests := []struct {
		name     string
		language string
		files    map[string]string
		expected string
	}{
		{
			name:     "detects rails framework",
			language: "ruby",
			files: map[string]string{
				"config/application.rb": "class Application < Rails::Application",
				"bin/rails":             "#!/usr/bin/env ruby",
			},
			expected: "rails",
		},
		{
			name:     "detects nextjs framework",
			language: "nodejs",
			files: map[string]string{
				"package.json":   `{"dependencies": {"next": "^13.0.0"}}`,
				"next.config.js": "module.exports = {}",
			},
			expected: "nextjs",
		},
		{
			name:     "detects react framework",
			language: "nodejs",
			files: map[string]string{
				"package.json": `{"dependencies": {"react": "^18.0.0"}}`,
			},
			expected: "react",
		},
		{
			name:     "detects gin framework",
			language: "go",
			files: map[string]string{
				"go.mod": "module test\ngo 1.19\nrequire github.com/gin-gonic/gin v1.9.0",
			},
			expected: "gin",
		},
		{
			name:     "detects spring framework",
			language: "java",
			files: map[string]string{
				"pom.xml": "<project><dependencies><dependency><groupId>org.springframework.boot</groupId><artifactId>spring-boot-starter-web</artifactId></dependency></dependencies></project>",
			},
			expected: "spring",
		},
		{
			name:     "defaults to unknown framework",
			language: "ruby",
			files:    map[string]string{},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "gitwo-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create test files
			for filename, content := range tt.files {
				filePath := filepath.Join(tempDir, filename)
				err := os.MkdirAll(filepath.Dir(filePath), 0o755)
				assert.NoError(t, err)
				err = os.WriteFile(filePath, []byte(content), 0o644)
				assert.NoError(t, err)
			}

			// Test detection
			result := DetectFramework(tempDir, tt.language)
			assert.Equal(t, tt.expected, result)
		})
	}
}
