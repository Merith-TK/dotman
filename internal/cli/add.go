package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/fileops"
	"github.com/Merith-TK/dotman/internal/git"
	"github.com/Merith-TK/dotman/internal/index"
)

var addCmd = &cobra.Command{
	Use:   "add <path>...",
	Short: "Add files to dotman",
	Long: `Add files to dotman management. Files are moved to the dotman repo
and symlinks are created in their original locations.

Examples:
  dotman add ~/.config/sway
  dotman add ~/.bashrc ~/.bash_aliases
  dotman add ~/.bash*`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAddMultiple(args)
	},
}

func runAdd(path string) error {
	// Expand the path
	expandedPath, err := config.ExpandPath(cfg, path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Check if path exists
	if !fileops.PathExists(expandedPath) {
		return fmt.Errorf("path does not exist: %s", expandedPath)
	}

	// Check if path is inside home directory
	if !config.IsInsideHome(cfg, expandedPath) {
		return fmt.Errorf("path must be inside home directory: %s", expandedPath)
	}

	// Ensure dotman directory exists
	if err := config.EnsureDotmanDir(cfg); err != nil {
		return fmt.Errorf("failed to create dotman directory: %w", err)
	}

	// Ensure git repository is initialized
	if err := git.EnsureRepo(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Load index
	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// Check if already managed
	if index.IsManaged(idx, expandedPath) {
		return fmt.Errorf("path is already managed: %s", expandedPath)
	}

	// Calculate repo path
	relativePath, err := config.RelativeToHome(cfg, expandedPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	repoPath := filepath.Join(cfg.DotmanDir, relativePath)

	// Get file type
	fileType := fileops.GetFileType(expandedPath)

	fmt.Printf("Adding %s to dotman management...\n", expandedPath)

	// Move file to repo
	if err := fileops.MoveToRepo(expandedPath, repoPath); err != nil {
		return fmt.Errorf("failed to move file to repo: %w", err)
	}

	// Create symlink
	if err := fileops.CreateSymlink(expandedPath, repoPath); err != nil {
		// Try to restore the file if symlink creation fails
		os.Rename(repoPath, expandedPath)
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Add to index
	index.AddFile(idx, expandedPath, relativePath, fileType)

	// Save index
	if err := index.Save(idx, cfg.IndexFile); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	// Commit changes
	if err := git.Add(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Convert to $HOME relative path for commit message
	homeRelPath, err := config.RelativeToHome(cfg, expandedPath)
	if err != nil {
		// Fallback to repo path if conversion fails
		homeRelPath = relativePath
	}
	homePath := "$HOME/" + homeRelPath

	commitMsg := fmt.Sprintf("Add %s to dotman management", homePath)
	if err := git.Commit(cfg.DotmanDir, commitMsg); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Printf("Successfully added %s to dotman management\n", path)
	return nil
}

func runAddMultiple(paths []string) error {
	var successCount int
	var failures []string

	for _, path := range paths {
		err := runAdd(path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", path, err))
		} else {
			successCount++
		}
	}

	if len(failures) > 0 {
		fmt.Printf("\nCompleted with %d successes and %d failures:\n", successCount, len(failures))
		for _, failure := range failures {
			fmt.Printf("  Error: %s\n", failure)
		}
		if successCount == 0 {
			return fmt.Errorf("all operations failed")
		}
	} else if successCount > 1 {
		fmt.Printf("\nSuccessfully added %d files to dotman management\n", successCount)
	}

	return nil
}
