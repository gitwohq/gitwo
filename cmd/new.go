package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gitwohq/gitwo/internal/wt"
	"github.com/spf13/cobra"
)

var (
	newStartRef       string
	newAutoSwitch     bool
	newOutputShell    bool
	newCreateScript   bool
	newSourceFunction bool
	newAutoSource     bool
)

func init() {
	newCmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new worktree with smart branch naming and auto-switch",
		Long: `Create a new worktree with automatic branch naming.

The branch name will be automatically generated from the worktree name:
- Converts kebab-case to feature/ prefix
- Example: 'feature-add-openapi' becomes 'feature/add-openapi'

Usage:
  gitwo new feature-add-openapi-related-workflows
  # Creates: ../feature-add-openapi-related-workflows with branch feature/add-openapi-related-workflows

Auto-switching options:
  --switch    Show auto-switch instructions (configurable via gitwo config)
  --shell     Output shell command for instant switching

Examples:
  gitwo new my-feature                    # Create worktree with auto-switch function
  gitwo new my-feature --switch          # Show auto-switch instructions
  gitwo new my-feature --shell           # Output shell command
  gitwo new my-feature --script          # Create switch script
  gitwo new my-feature --source          # Output shell function
  gitwo new my-feature --auto-source     # Provide auto-source instructions
  eval "$(gitwo new my-feature --shell | tail -1)"  # Auto-switch instantly`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Generate path: ../<name>
			path := filepath.Join("..", name)

			// Generate branch name: feature/<name>
			branch := fmt.Sprintf("feature/%s", name)

			// Load configuration
			config, err := wt.LoadConfig()
			if err != nil {
				// Use default config if loading fails
				config = wt.DefaultConfig()
			}

			// Override config with flag if explicitly provided
			if cmd.Flags().Changed("switch") {
				config.AutoSwitch = newAutoSwitch
			}

			// Use silent mode if --shell flag is used
			var result *wt.AddResult
			if newOutputShell {
				result, err = wt.AddWithConfigSilent(path, branch, newStartRef, config, true)
			} else {
				result, err = wt.AddWithConfig(path, branch, newStartRef, config)
			}
			if err != nil {
				return err
			}

			// If --shell flag is used, only output the shell command
			if newOutputShell {
				fmt.Printf("cd %s\n", path)
				return nil
			}

			// Print footer with next steps
			if config != nil && config.AutoSwitch {
				wt.PrintFooterWithSwitch(path, result.Branch)
			} else {
				wt.PrintFooter(path, result.Branch)
			}

			// Output shell command if auto-switch is enabled by default
			if config != nil && config.AutoSwitch {
				fmt.Printf("\n# Shell command to switch to worktree:\n")
				fmt.Printf("cd %s\n", path)
			}

			// For true auto-switch, output a shell function that can be sourced
			if config != nil && config.AutoSwitch {
				fmt.Printf("\n🚀 Auto-switch function (copy & paste):\n")
				fmt.Printf("gitwo-switch-%s() {\n", name)
				fmt.Printf("    cd %s\n", path)
				fmt.Printf("    echo \"✅ Switched to worktree: %s\"\n", path)
				fmt.Printf("    echo \"   Branch: %s\"\n", result.Branch)
				fmt.Printf("    git status\n")
				fmt.Printf("}\n")
				fmt.Printf("\n💡 Then run: gitwo-switch-%s\n", name)
				fmt.Printf("\n🔧 Or use eval: eval \"$(gitwo new %s --shell | tail -1)\"\n", name)
				fmt.Printf("\n⚡ Or auto-source: source <(gitwo new %s --source | grep -A 5 'gitwo-switch-%s()')\n", name, name)
				fmt.Printf("\n🔄 Or one-liner: cd %s && echo \"✅ Switched to worktree: %s\" && git status\n", path, path)
			}

			// Output shell function for sourcing
			if newSourceFunction {
				fmt.Printf("\n# Add this to your shell profile (.bashrc, .zshrc, etc.):\n")
				fmt.Printf("gitwo-switch-%s() {\n", name)
				fmt.Printf("    cd %s\n", path)
				fmt.Printf("    echo \"✅ Switched to worktree: %s\"\n", path)
				fmt.Printf("    echo \"   Branch: %s\"\n", result.Branch)
				fmt.Printf("    git status\n")
				fmt.Printf("}\n")
				fmt.Printf("\n# Then run: gitwo-switch-%s\n", name)
			}

			// Auto-source the function if requested
			if newAutoSource {
				fmt.Printf("\n🔄 Auto-sourcing switch function...\n")
				fmt.Printf("   Run: source <(gitwo new %s --source | grep -A 5 'gitwo-switch-%s()')\n", name, name)
				fmt.Printf("\n💡 Alternative for Oh My Zsh:\n")
				fmt.Printf("   alias gitwo-switch-%s='cd %s && echo \"✅ Switched to worktree: %s\" && git status'\n", name, path, path)
			}

			// Create a shell script for easy switching
			if newCreateScript {
				scriptName := fmt.Sprintf("switch-to-%s.sh", name)
				scriptContent := fmt.Sprintf(`#!/bin/bash
# Auto-generated script to switch to worktree: %s
cd %s
echo "✅ Switched to worktree: %s"
echo "   Branch: %s"
echo "   Path: %s"
`, name, path, path, result.Branch, path)

				if err := os.WriteFile(scriptName, []byte(scriptContent), 0o755); err != nil {
					fmt.Printf("\n⚠️  Warning: Could not create switch script: %v\n", err)
				} else {
					fmt.Printf("\n📜 Created switch script: %s\n", scriptName)
					fmt.Printf("   Run: ./%s\n", scriptName)
				}
			}

			return nil
		},
	}

	newCmd.Flags().StringVar(&newStartRef, "start-point", "HEAD", "start point ref (default HEAD)")
	newCmd.Flags().BoolVar(&newAutoSwitch, "switch", false, "automatically switch to the new worktree after creation")
	newCmd.Flags().BoolVar(&newOutputShell, "shell", false, "output shell command for switching to worktree")
	newCmd.Flags().BoolVar(&newCreateScript, "script", false, "create a shell script for easy switching")
	newCmd.Flags().BoolVar(&newSourceFunction, "source", false, "output shell function for sourcing")
	newCmd.Flags().BoolVar(&newAutoSource, "auto-source", false, "provide auto-source instructions")

	rootCmd.AddCommand(newCmd)
}
