package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DetectShell detects the current shell type
func DetectShell() string {
	// Windows: prefer PowerShell
	if runtime.GOOS == "windows" || os.Getenv("PSModulePath") != "" {
		return "powershell"
	}

	// Check SHELL environment variable
	if shell := os.Getenv("SHELL"); shell != "" {
		switch {
		case strings.Contains(shell, "zsh"):
			return "zsh"
		case strings.Contains(shell, "fish"):
			return "fish"
		case strings.Contains(shell, "bash"):
			return "bash"
		}
	}

	// Default to bash
	return "bash"
}

// GetShellProfilePath returns the path to the shell profile file
func GetShellProfilePath(shell string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	switch shell {
	case "zsh":
		return filepath.Join(homeDir, ".zshrc")
	case "bash":
		return filepath.Join(homeDir, ".bashrc")
	case "fish":
		return filepath.Join(homeDir, ".config", "fish", "config.fish")
	case "powershell":
		// Windows PowerShell profile
		if runtime.GOOS == "windows" {
			if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
				return filepath.Join(userProfile, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
			}
		}
		return filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	default:
		return filepath.Join(homeDir, ".bashrc")
	}
}

// IsWrapperInstalled checks if the gitwo wrapper is installed in the shell profile
func IsWrapperInstalled(shell string) bool {
	profilePath := GetShellProfilePath(shell)

	// Read profile file
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return false
	}

	content := string(data)
	return strings.Contains(content, "# >>> gitwo shell wrapper >>>")
}

// GetHomeDir returns the user's home directory
func GetHomeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return "."
}

// GetWindowsHome returns the Windows user profile directory
func GetWindowsHome() string {
	if runtime.GOOS == "windows" {
		if up := os.Getenv("USERPROFILE"); up != "" {
			return up
		}
	}
	return GetHomeDir()
}
