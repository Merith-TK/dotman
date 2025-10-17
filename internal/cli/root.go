package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/pkg/types"
)

var (
	cfg *types.Config
)

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "dotman",
	Short: "A dotfiles manager that centralizes configuration files",
	Long: `Dotman is a dotfiles manager that moves your configuration files to a 
centralized git repository and creates symlinks in their original locations.

This allows you to easily track, manage, and sync your dotfiles across systems
while keeping them in their expected locations for applications to find them.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.New()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		return nil
	},
}

func init() {
	// Disable auto-generated commands
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(syncCmd)

	// Add flags
	addCmd.Flags().BoolP("force", "f", false, "Force operation even if conflicts exist")
	addCmd.Flags().BoolP("dry-run", "n", false, "Show what would happen without doing it")
	addCmd.Flags().BoolP("backup", "b", false, "Create backup before operation")

	deployCmd.Flags().BoolP("force", "f", false, "Force deployment even if conflicts exist")
	deployCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without doing it")
	deployCmd.Flags().BoolP("backup", "b", false, "Create backup before operation")
}
