package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitwohq/gitwo/internal/shell"
	"github.com/spf13/cobra"
)

var shellInstallShell string

const (
	beginMarker = "# >>> gitwo shell wrapper >>>"
	endMarker   = "# <<< gitwo shell wrapper <<<"
)

var shellInstallCmd = &cobra.Command{
	Use:   "shell-install",
	Short: "Persist the gitwo auto-cd wrapper into your shell profile",
	Long: `Installs the gitwo wrapper into your shell profile for persistent auto-cd functionality.

Examples:
  gitwo shell-install                    # Auto-detect shell and install wrapper
  gitwo shell-install --shell bash       # Install bash wrapper
  gitwo shell-install --shell zsh        # Install zsh wrapper
  gitwo shell-install --shell fish       # Install fish wrapper
  gitwo shell-install --shell powershell # Install PowerShell wrapper

After installation, restart your shell or reload the profile:
  source ~/.bashrc                       # bash
  source ~/.zshrc                        # zsh
  source ~/.config/fish/config.fish      # fish
  . ~/Documents/PowerShell/Microsoft.PowerShell_profile.ps1  # PowerShell`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sh := shellInstallShell
		if sh == "" {
			sh = shell.DetectShell()
		}

		profilePath := shell.GetShellProfilePath(sh)
		wrapperContent := shell.GenerateWrapper(sh)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(profilePath), 0o755); err != nil {
			return fmt.Errorf("failed to create profile directory: %w", err)
		}

		// Read existing content
		existingContent := readFileOrEmpty(profilePath)

		// Remove existing wrapper block
		cleanContent := removeBlock(existingContent, beginMarker, endMarker)

		// Add new wrapper block
		newContent := cleanContent + "\n" + beginMarker + "\n" + wrapperContent + "\n" + endMarker + "\n"

		// Write to file
		if err := os.WriteFile(profilePath, []byte(newContent), 0o644); err != nil {
			return fmt.Errorf("failed to write profile: %w", err)
		}

		fmt.Printf("✅ Installed wrapper into %s\n", profilePath)

		// Print reload instructions
		switch sh {
		case "bash", "zsh":
			fmt.Printf("   Reload your shell or run: source %s\n", profilePath)
		case "fish":
			fmt.Printf("   Reload your shell or run: source %s\n", profilePath)
		case "powershell":
			fmt.Printf("   Restart PowerShell or run: . '%s'\n", profilePath)
		}

		return nil
	},
}

func init() {
	shellInstallCmd.Flags().StringVar(&shellInstallShell, "shell", "", "bash|zsh|fish|powershell")
	rootCmd.AddCommand(shellInstallCmd)
}

// readFileOrEmpty reads a file or returns empty string if it doesn't exist
func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// removeBlock removes a block of text between markers
func removeBlock(content, begin, end string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	inBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == begin {
			inBlock = true
			continue
		}
		if strings.TrimSpace(line) == end {
			inBlock = false
			continue
		}
		if !inBlock {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
