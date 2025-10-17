package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/git"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage git remote for dotman repository",
	Long: `Manage the git remote repository for your dotman dotfiles.

Use 'remote set <url>' to set the remote repository URL.
Use 'remote get' to show the current remote repository URL.`,
}

var remoteSetCmd = &cobra.Command{
	Use:   "set <url>",
	Short: "Set the git remote URL",
	Long: `Set the git remote URL for the dotman repository.

This is equivalent to running 'git remote set-url origin <url>' 
or 'git remote add origin <url>' if no remote exists.

Example:
  dotman remote set https://github.com/user/dotfiles.git`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoteSet(args[0])
	},
}

var remoteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current git remote URL",
	Long: `Show the current git remote URL for the dotman repository.

This is equivalent to running 'git remote get-url origin'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoteGet()
	},
}

func init() {
	remoteCmd.AddCommand(remoteSetCmd)
	remoteCmd.AddCommand(remoteGetCmd)
}

func runRemoteSet(url string) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	if !git.IsGitRepo(cfg.DotmanDir) {
		return fmt.Errorf("dotman directory is not a git repository")
	}

	// Check if origin remote already exists
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = cfg.DotmanDir
	
	if err := cmd.Run(); err != nil {
		// Remote doesn't exist, add it
		fmt.Printf("Adding remote origin: %s\n", url)
		addCmd := exec.Command("git", "remote", "add", "origin", url)
		addCmd.Dir = cfg.DotmanDir
		
		if output, err := addCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add remote: %s, %w", string(output), err)
		}
	} else {
		// Remote exists, update it
		fmt.Printf("Updating remote origin: %s\n", url)
		setCmd := exec.Command("git", "remote", "set-url", "origin", url)
		setCmd.Dir = cfg.DotmanDir
		
		if output, err := setCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set remote URL: %s, %w", string(output), err)
		}
	}

	fmt.Println("Remote origin set successfully")
	return nil
}

func runRemoteGet() error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	if !git.IsGitRepo(cfg.DotmanDir) {
		return fmt.Errorf("dotman directory is not a git repository")
	}

	// Get the remote URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = cfg.DotmanDir
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no remote origin configured")
	}

	fmt.Printf("Remote origin: %s", string(output))
	return nil
}