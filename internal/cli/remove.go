package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/fileops"
	"github.com/Merith-TK/dotman/internal/git"
	"github.com/Merith-TK/dotman/internal/index"
)

var removeCmd = &cobra.Command{
	Use:   "remove <path>...",
	Short: "Remove files from dotman",
	Long: `Remove files from dotman management. Files are restored from the repo
back to their original locations and removed from management.

Examples:
  dotman remove ~/.config/sway
  dotman remove ~/.bashrc ~/.bash_aliases
  dotman remove ~/.bash*`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemoveMultiple(args)
	},
}

func runRemove(path string) error {
	// Expand the path
	expandedPath, err := config.ExpandPath(cfg, path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Load index
	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// Check if managed
	managedFile, found := index.FindFile(idx, expandedPath)
	if !found {
		return fmt.Errorf("path is not managed by dotman: %s", expandedPath)
	}

	repoPath := filepath.Join(cfg.DotmanDir, managedFile.RepoPath)

	fmt.Printf("Removing %s from dotman management...\n", expandedPath)

	// Remove symlink and restore original
	if err := fileops.RemoveSymlink(expandedPath, repoPath); err != nil {
		return fmt.Errorf("failed to remove symlink and restore file: %w", err)
	}

	// Remove from index
	index.RemoveFile(idx, expandedPath)

	// Save index
	if err := index.Save(idx, cfg.IndexFile); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	// Commit changes
	if err := git.Add(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Convert to $HOME relative path for commit message
	homeRelPath, err := config.RelativeToHome(cfg, managedFile.OriginalPath)
	if err != nil {
		// Fallback to repo path if conversion fails
		homeRelPath = managedFile.RepoPath
	}
	homePath := "$HOME/" + homeRelPath

	commitMsg := fmt.Sprintf("Remove %s from dotman management", homePath)
	if err := git.Commit(cfg.DotmanDir, commitMsg); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Printf("Successfully removed %s from dotman management\n", path)
	return nil
}

func runRemoveMultiple(paths []string) error {
	var successCount int
	var failures []string

	for _, path := range paths {
		err := runRemove(path)
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
		fmt.Printf("\nSuccessfully removed %d files from dotman management\n", successCount)
	}

	return nil
}
