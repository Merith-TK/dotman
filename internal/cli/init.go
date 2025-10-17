package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/git"
	"github.com/Merith-TK/dotman/internal/index"
	"github.com/Merith-TK/dotman/pkg/types"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new dotman repo",
	Long: `Initialize creates a new dotman repo in ~/.dotman.

If ~/.dotman already exists and is a valid dotman repo, this command does nothing.
If ~/.dotman exists but is not a git repository, it will be initialized as one.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone <url>",
	Short: "Clone an existing dotfiles repo",
	Long: `Clone downloads an existing dotfiles repo to ~/.dotman.

This command will fail if ~/.dotman already exists.
After cloning, use 'dotman deploy' to create symlinks.

Example:
  dotman clone https://github.com/user/dotfiles.git`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClone(args[0])
	},
}

func runInit() error {
	// Check if dotman directory already exists
	if config.DotmanDirExists(cfg) {
		// Check if it's already a git repository
		if git.IsGitRepo(cfg.DotmanDir) {
			fmt.Println("Dotman repo already initialized at", cfg.DotmanDir)
			return nil
		} else {
			// Directory exists but is not a git repo - check if it's empty or has dotman files
			if config.IndexFileExists(cfg) {
				// Has index file, so initialize git
				fmt.Println("Initializing git repository in existing dotman directory...")
				return git.EnsureRepo(cfg.DotmanDir)
			} else {
				// Directory exists but doesn't look like dotman - error
				return fmt.Errorf("directory %s exists but is not a dotman repo", cfg.DotmanDir)
			}
		}
	}

	// Create dotman directory
	if err := config.EnsureDotmanDir(cfg); err != nil {
		return fmt.Errorf("failed to create dotman directory: %w", err)
	}

	// Initialize git repository with initial files
	if err := git.EnsureRepo(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create empty index
	idx := &types.Index{
		Version:      config.DefaultVersion,
		ManagedFiles: make([]types.ManagedFile, 0),
	}

	if err := index.Save(idx, cfg.IndexFile); err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}

	// Commit the initial index
	if err := git.Add(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to stage initial files: %w", err)
	}

	if err := git.Commit(cfg.DotmanDir, "Initialize dotman repository with empty index"); err != nil {
		return fmt.Errorf("failed to commit initial files: %w", err)
	}

	fmt.Printf("Initialized dotman repo in %s\n", cfg.DotmanDir)
	return nil
}

func runClone(url string) error {
	// Check if dotman directory already exists
	if config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory already exists: %s", cfg.DotmanDir)
	}

	// Clone the repository
	fmt.Printf("Cloning dotfiles repo from %s...\n", url)

	cmd := exec.Command("git", "clone", url, cfg.DotmanDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %s, %w", string(output), err)
	}

	// Validate that the cloned repository has a valid index file
	if !config.IndexFileExists(cfg) {
		// Clean up the failed clone
		os.RemoveAll(cfg.DotmanDir)
		return fmt.Errorf("cloned repository does not contain a valid dotman index.json file")
	}

	// Validate the index file can be loaded
	_, err := index.Load(cfg.IndexFile)
	if err != nil {
		// Clean up the failed clone
		os.RemoveAll(cfg.DotmanDir)
		return fmt.Errorf("cloned repository has invalid index.json: %w", err)
	}

	fmt.Printf("Successfully cloned dotfiles repo to %s\n", cfg.DotmanDir)
	fmt.Println("Use 'dotman sync' to discover and deploy all files in the repo,")
	fmt.Println("or 'dotman deploy' to deploy only files already in the index.")
	return nil
}
