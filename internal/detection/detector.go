package detection

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectLanguage detects the programming language of the project
func DetectLanguage(repoPath string) string {
	// Check for Ruby/Rails
	if hasFile(repoPath, "Gemfile") {
		return "ruby"
	}

	// Check for Node.js
	if hasFile(repoPath, "package.json") {
		return "nodejs"
	}

	// Check for Go
	if hasFile(repoPath, "go.mod") {
		return "go"
	}

	// Check for Java
	if hasFile(repoPath, "pom.xml") || hasFile(repoPath, "build.gradle") {
		return "java"
	}

	// Check for Python
	if hasFile(repoPath, "requirements.txt") || hasFile(repoPath, "setup.py") || hasFile(repoPath, "pyproject.toml") {
		return "python"
	}

	// Check for PHP
	if hasFile(repoPath, "composer.json") {
		return "php"
	}

	// Check for Rust
	if hasFile(repoPath, "Cargo.toml") {
		return "rust"
	}

	// Check for C/C++
	if hasFile(repoPath, "CMakeLists.txt") || hasFile(repoPath, "Makefile") {
		return "c"
	}

	return "unknown"
}

// DetectFramework detects the framework for a given language
func DetectFramework(repoPath, language string) string {
	switch language {
	case "ruby":
		return detectRubyFramework(repoPath)
	case "nodejs":
		return detectNodeJSFramework(repoPath)
	case "go":
		return detectGoFramework(repoPath)
	case "java":
		return detectJavaFramework(repoPath)
	case "python":
		return detectPythonFramework(repoPath)
	case "php":
		return detectPHPFramework(repoPath)
	case "rust":
		return detectRustFramework(repoPath)
	default:
		return "unknown"
	}
}

// detectRubyFramework detects Ruby frameworks
func detectRubyFramework(repoPath string) string {
	// Check for Rails
	if hasFile(repoPath, "config/application.rb") && hasFile(repoPath, "bin/rails") {
		return "rails"
	}

	// Check for Sinatra
	if hasFile(repoPath, "config.ru") && hasFile(repoPath, "app.rb") {
		return "sinatra"
	}

	// Check for Rack
	if hasFile(repoPath, "config.ru") {
		return "rack"
	}

	return "unknown"
}

// detectNodeJSFramework detects Node.js frameworks
func detectNodeJSFramework(repoPath string) string {
	// Check for Next.js
	if hasFile(repoPath, "next.config.js") || hasFile(repoPath, "next.config.mjs") {
		return "nextjs"
	}

	// Check for React
	if hasFile(repoPath, "package.json") && containsInFile(repoPath, "package.json", "react") {
		return "react"
	}

	// Check for Vue
	if hasFile(repoPath, "package.json") && containsInFile(repoPath, "package.json", "vue") {
		return "vue"
	}

	// Check for Express
	if hasFile(repoPath, "package.json") && containsInFile(repoPath, "package.json", "express") {
		return "express"
	}

	// Check for Nuxt
	if hasFile(repoPath, "nuxt.config.js") || hasFile(repoPath, "nuxt.config.ts") {
		return "nuxt"
	}

	return "unknown"
}

// detectGoFramework detects Go frameworks
func detectGoFramework(repoPath string) string {
	// Check for Gin
	if hasFile(repoPath, "go.mod") && containsInFile(repoPath, "go.mod", "gin-gonic/gin") {
		return "gin"
	}

	// Check for Echo
	if hasFile(repoPath, "go.mod") && containsInFile(repoPath, "go.mod", "labstack/echo") {
		return "echo"
	}

	// Check for Fiber
	if hasFile(repoPath, "go.mod") && containsInFile(repoPath, "go.mod", "gofiber/fiber") {
		return "fiber"
	}

	// Check for Chi
	if hasFile(repoPath, "go.mod") && containsInFile(repoPath, "go.mod", "go-chi/chi") {
		return "chi"
	}

	return "unknown"
}

// detectJavaFramework detects Java frameworks
func detectJavaFramework(repoPath string) string {
	// Check for Spring Boot
	if hasFile(repoPath, "pom.xml") && containsInFile(repoPath, "pom.xml", "spring-boot") {
		return "spring"
	}

	// Check for Spring Boot (Gradle)
	if hasFile(repoPath, "build.gradle") && containsInFile(repoPath, "build.gradle", "spring-boot") {
		return "spring"
	}

	// Check for Quarkus
	if hasFile(repoPath, "pom.xml") && containsInFile(repoPath, "pom.xml", "quarkus") {
		return "quarkus"
	}

	// Check for Micronaut
	if hasFile(repoPath, "build.gradle") && containsInFile(repoPath, "build.gradle", "micronaut") {
		return "micronaut"
	}

	return "unknown"
}

// detectPythonFramework detects Python frameworks
func detectPythonFramework(repoPath string) string {
	// Check for Django
	if hasFile(repoPath, "manage.py") || hasFile(repoPath, "requirements.txt") && containsInFile(repoPath, "requirements.txt", "django") {
		return "django"
	}

	// Check for Flask
	if hasFile(repoPath, "requirements.txt") && containsInFile(repoPath, "requirements.txt", "flask") {
		return "flask"
	}

	// Check for FastAPI
	if hasFile(repoPath, "requirements.txt") && containsInFile(repoPath, "requirements.txt", "fastapi") {
		return "fastapi"
	}

	return "unknown"
}

// detectPHPFramework detects PHP frameworks
func detectPHPFramework(repoPath string) string {
	// Check for Laravel
	if hasFile(repoPath, "artisan") {
		return "laravel"
	}

	// Check for Symfony
	if hasFile(repoPath, "bin/console") {
		return "symfony"
	}

	// Check for CodeIgniter
	if hasFile(repoPath, "index.php") && hasFile(repoPath, "application") {
		return "codeigniter"
	}

	return "unknown"
}

// detectRustFramework detects Rust frameworks
func detectRustFramework(repoPath string) string {
	// Check for Actix
	if hasFile(repoPath, "Cargo.toml") && containsInFile(repoPath, "Cargo.toml", "actix-web") {
		return "actix"
	}

	// Check for Rocket
	if hasFile(repoPath, "Cargo.toml") && containsInFile(repoPath, "Cargo.toml", "rocket") {
		return "rocket"
	}

	// Check for Warp
	if hasFile(repoPath, "Cargo.toml") && containsInFile(repoPath, "Cargo.toml", "warp") {
		return "warp"
	}

	return "unknown"
}

// hasFile checks if a file exists in the repository
func hasFile(repoPath, filename string) bool {
	filePath := filepath.Join(repoPath, filename)
	_, err := os.Stat(filePath)
	return err == nil
}

// containsInFile checks if a file contains a specific string
func containsInFile(repoPath, filename, content string) bool {
	filePath := filepath.Join(repoPath, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), content)
}
