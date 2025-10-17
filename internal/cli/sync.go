package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/fileops"
	"github.com/Merith-TK/dotman/internal/git"
	"github.com/Merith-TK/dotman/internal/index"
	"github.com/Merith-TK/dotman/pkg/types"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotman repository with git remote",
	Long: `Sync handles git operations for the dotman repository.

With --pull flag, pulls changes from git remote.
With --push flag, pushes local changes to git remote.
Without flags, discovers and adds unmanaged files in the repo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pull, _ := cmd.Flags().GetBool("pull")
		push, _ := cmd.Flags().GetBool("push")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if pull {
			return runSyncPull(dryRun)
		}
		if push {
			return runSyncPush(dryRun)
		}

		// Default behavior: discover unmanaged files
		return runSyncDiscover(dryRun, false)
	},
}

func init() {
	syncCmd.Flags().BoolP("pull", "", false, "Pull changes from git remote")
	syncCmd.Flags().BoolP("push", "", false, "Push local changes to git remote")
	syncCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without doing it")
}

func runSyncPull(dryRun bool) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	if !git.IsGitRepo(cfg.DotmanDir) {
		return fmt.Errorf("dotman directory is not a git repository")
	}

	fmt.Println("Pulling changes from git remote...")

	if dryRun {
		fmt.Println("Dry-run mode: would pull changes from remote")
		return nil
	}

	if err := git.Pull(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to pull from remote: %w", err)
	}

	fmt.Println("Successfully pulled changes from remote")
	return nil
}

func runSyncPush(dryRun bool) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	if !git.IsGitRepo(cfg.DotmanDir) {
		return fmt.Errorf("dotman directory is not a git repository")
	}

	// Check if there are any changes to push
	hasChanges, err := git.HasChanges(cfg.DotmanDir)
	if err != nil {
		return fmt.Errorf("failed to check for changes: %w", err)
	}

	if hasChanges {
		fmt.Println("Warning: You have uncommitted changes. Commit them first or they won't be pushed.")
	}

	fmt.Println("Pushing changes to git remote...")

	if dryRun {
		fmt.Println("Dry-run mode: would push changes to remote")
		return nil
	}

	if err := git.Push(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to push to remote: %w", err)
	}

	fmt.Println("Successfully pushed changes to remote")
	return nil
}

// runSyncDiscover scans the .dotman directory for unmanaged files and adds them to the index
func runSyncDiscover(dryRun bool, autoAdd bool) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	// Load current index
	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// Find unmanaged files in the repo
	unmanaged, err := findUnmanagedFiles(cfg.DotmanDir, idx)
	if err != nil {
		return fmt.Errorf("failed to scan repo: %w", err)
	}

	if len(unmanaged) == 0 {
		fmt.Println("All repo files are already managed in the index.")
		return nil
	}

	fmt.Printf("Found %d unmanaged file(s) in repo:\n", len(unmanaged))
	for _, file := range unmanaged {
		fmt.Printf("  %s\n", file)
	}

	if dryRun {
		fmt.Println("\nDry-run mode: would add these files to the index")
		return nil
	}

	if !autoAdd {
		fmt.Print("\nAdd these files to the index? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Sync cancelled.")
			return nil
		}
	}

	// Add unmanaged files to the index
	added := 0
	var addedPaths []string
	for _, repoPath := range unmanaged {
		if err := addUnmanagedFile(idx, repoPath); err != nil {
			fmt.Printf("Warning: failed to add %s: %v\n", repoPath, err)
			continue
		}
		added++
		addedPaths = append(addedPaths, "$HOME/"+repoPath)
		if !autoAdd {
			fmt.Printf("Added %s to index\n", repoPath)
		}
	}

	if added == 0 {
		fmt.Println("No files were added to the index.")
		return nil
	}

	// Save updated index
	if err := index.Save(idx, cfg.IndexFile); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	// Commit the changes
	if err := git.Add(cfg.DotmanDir); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create commit message with actual paths
	var commitMsg string
	if len(addedPaths) == 1 {
		commitMsg = fmt.Sprintf("Sync: add %s to index", addedPaths[0])
	} else if len(addedPaths) <= 3 {
		commitMsg = fmt.Sprintf("Sync: add %s to index", strings.Join(addedPaths, ", "))
	} else {
		commitMsg = fmt.Sprintf("Sync: add %d files to index (%s, ...)", len(addedPaths), strings.Join(addedPaths[:2], ", "))
	}

	if err := git.Commit(cfg.DotmanDir, commitMsg); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	if autoAdd {
		fmt.Printf("Auto-synced %d file(s)\n", added)
	} else {
		fmt.Printf("Successfully synced %d file(s)\n", added)
	}

	return nil
}

// findUnmanagedFiles scans the repo directory and returns files not in the index
func findUnmanagedFiles(repoDir string, idx *types.Index) ([]string, error) {
	var unmanaged []string

	// Get all managed directories first
	managedDirs := getManagedDirectories(idx)

	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory and index.json
		if strings.Contains(path, ".git") || strings.HasSuffix(path, "index.json") {
			if info.IsDir() && strings.HasSuffix(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories - we only track files
		if info.IsDir() {
			return nil
		}

		// Get relative path from repo root
		relPath, err := filepath.Rel(repoDir, path)
		if err != nil {
			return err
		}

		// Check if this file is already managed in the index
		originalPath := filepath.Join(cfg.HomeDir, relPath)
		if !index.IsManaged(idx, originalPath) {
			// Also check if this file is covered by a managed directory
			if !isWithinManagedDirectory(originalPath, managedDirs) {
				unmanaged = append(unmanaged, relPath)
			}
		}

		return nil
	})

	return unmanaged, err
}

// addUnmanagedFile adds a single unmanaged file to the index
func addUnmanagedFile(idx *types.Index, repoPath string) error {
	// Calculate the original path (where the symlink should be)
	originalPath := filepath.Join(cfg.HomeDir, repoPath)

	// Get the full repository path
	fullRepoPath := filepath.Join(cfg.DotmanDir, repoPath)

	// Get file type
	fileType := fileops.GetFileType(fullRepoPath)

	// Add to index
	index.AddFile(idx, originalPath, repoPath, fileType)

	return nil
}
