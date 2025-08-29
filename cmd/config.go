package cmd

import (
	"fmt"

	"github.com/gitwohq/gitwo/internal/wt"
	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage gitwo configuration",
		Long: `Manage gitwo configuration settings.

Configuration is stored in .gitwo/config.yml in your repository root.`,
	}

	// Show current config
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := wt.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Printf("Current gitwo configuration:\n")
			fmt.Printf("  auto_switch: %t\n", config.AutoSwitch)
			return nil
		},
	}

	// Set auto-switch
	setAutoSwitchCmd := &cobra.Command{
		Use:   "auto-switch [true|false]",
		Short: "Set auto-switch behavior",
		Long: `Set whether to automatically switch to new worktrees after creation.

Examples:
  gitwo config auto-switch true   # Enable auto-switch
  gitwo config auto-switch false  # Disable auto-switch`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := args[0]

			// Load current config
			config, err := wt.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Parse boolean value
			switch value {
			case "true", "1", "yes", "on":
				config.AutoSwitch = true
			case "false", "0", "no", "off":
				config.AutoSwitch = false
			default:
				return fmt.Errorf("invalid value: %s (use true/false, 1/0, yes/no, on/off)", value)
			}

			// Save config
			if err := wt.SaveConfig(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("auto_switch set to: %t\n", config.AutoSwitch)
			return nil
		},
	}

	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(setAutoSwitchCmd)
	rootCmd.AddCommand(configCmd)
}
