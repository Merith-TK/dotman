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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of managed files",
	Long: `Show information about all files currently managed by dotman.
	
With --fix flag, repairs broken or missing symlinks.
With --cleanup flag, removes redundant individual file entries that are covered by managed directories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fix, _ := cmd.Flags().GetBool("fix")
		cleanup, _ := cmd.Flags().GetBool("cleanup")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return runStatus(fix, cleanup, dryRun)
	},
}

func init() {
	statusCmd.Flags().BoolP("fix", "f", false, "Fix broken or missing symlinks")
	statusCmd.Flags().BoolP("cleanup", "c", false, "Remove redundant file entries covered by managed directories")
	statusCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without doing it")
}

func runStatus(fix bool, cleanup bool, dryRun bool) error {
	if !config.DotmanDirExists(cfg) {
		fmt.Println("Dotman not initialized. Use 'dotman add' to start managing files.")
		return nil
	}

	// Run cleanup first if requested
	if cleanup {
		fmt.Println("Cleaning up redundant file entries...")
		if err := runCleanup(dryRun); err != nil {
			fmt.Printf("Warning: cleanup failed: %v\n", err)
		}
		fmt.Println()
	}

	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	if index.Count(idx) == 0 {
		fmt.Println("No files are currently managed by dotman.")
		return nil
	}

	fmt.Printf("Dotman is managing %d file(s):\n\n", index.Count(idx))

	// Get all managed directories first
	managedDirs := getManagedDirectories(idx)

	brokenCount := 0
	for _, file := range index.GetAllFiles(idx) {
		status := "✓"
		statusMsg := "OK"

		// Skip individual files that are within managed directories
		if file.Type == types.FileTypeFile && isWithinManagedDirectory(file.OriginalPath, managedDirs) {
			continue
		}

		// Check if symlink exists and is valid
		if fileops.PathExists(file.OriginalPath) {
			if !fileops.IsSymlink(file.OriginalPath) {
				status = "✗"
				statusMsg = "Not a symlink"
				brokenCount++
			}
		} else {
			status = "✗"
			statusMsg = "Missing"
			brokenCount++
		}

		fmt.Printf("%s %s (%s) - %s\n", status, file.OriginalPath, file.Type, statusMsg)
	}

	// Run fix if requested and there are broken symlinks
	if fix && brokenCount > 0 {
		fmt.Printf("\nFound %d broken symlink(s). ", brokenCount)
		if dryRun {
			fmt.Println("Would fix them (dry-run mode).")
		} else {
			fmt.Println("Fixing them...")
			if err := runFix(dryRun); err != nil {
				fmt.Printf("Fix failed: %v\n", err)
			}
		}
	} else if fix && brokenCount == 0 {
		fmt.Println("\nAll symlinks are working correctly.")
	}

	// Show git repository information if repository exists
	if git.IsGitRepo(cfg.DotmanDir) {
		fmt.Println("\nGit Repository Information:")

		// Show current branch
		if branch, err := git.GetCurrentBranch(cfg.DotmanDir); err == nil {
			fmt.Printf("Branch: %s\n", branch)
		} else {
			fmt.Printf("Branch: <unknown> (%v)\n", err)
		}

		// Show remote URL
		if remoteURL, err := git.GetRemoteURL(cfg.DotmanDir); err == nil {
			fmt.Printf("Remote: %s\n", remoteURL)
		} else {
			fmt.Printf("Remote: <not configured>\n")
		}

		// Show commit count
		if commitCount, err := git.GetCommitCount(cfg.DotmanDir); err == nil {
			fmt.Printf("Commits: %s\n", commitCount)
		}

		// Show uncommitted changes
		hasChanges, err := git.HasChanges(cfg.DotmanDir)
		if err == nil {
			if hasChanges {
				fmt.Println("Status: Uncommitted changes")
				if gitStatus, err := git.Status(cfg.DotmanDir); err == nil {
					// Parse and format git status output
					lines := strings.Split(strings.TrimSpace(gitStatus), "\n")
					for _, line := range lines {
						if len(line) >= 3 {
							status := line[:2]
							file := line[3:]
							switch status {
							case "??":
								fmt.Printf("  + %s (untracked)\n", file)
							case " M":
								fmt.Printf("  ~ %s (modified)\n", file)
							case "M ":
								fmt.Printf("  ~ %s (staged)\n", file)
							case " D":
								fmt.Printf("  - %s (deleted)\n", file)
							case "D ":
								fmt.Printf("  - %s (staged for deletion)\n", file)
							case "A ":
								fmt.Printf("  + %s (added)\n", file)
							default:
								fmt.Printf("  %s %s\n", status, file)
							}
						}
					}
				}
			} else {
				fmt.Println("Status: Clean (no changes)")
			}
		}
	} else {
		fmt.Println("\nGit Repository: Not initialized")
	}

	return nil
}

// runFix fixes broken or missing symlinks for managed files
func runFix(dryRun bool) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	if index.Count(idx) == 0 {
		fmt.Println("No files are managed by dotman.")
		return nil
	}

	fixed := 0
	problems := 0

	for _, file := range index.GetAllFiles(idx) {
		repoPath := filepath.Join(cfg.DotmanDir, file.RepoPath)

		// Check if repo file exists
		if !fileops.PathExists(repoPath) {
			fmt.Printf("⚠️  Repository file missing: %s\n", file.OriginalPath)
			problems++
			continue
		}

		// Check original location status
		if fileops.PathExists(file.OriginalPath) {
			if fileops.IsSymlink(file.OriginalPath) {
				// Check if symlink points to correct location
				if target, err := os.Readlink(file.OriginalPath); err == nil {
					if target == repoPath {
						continue // Already correct
					} else {
						fmt.Printf("⚠️  %s - Symlink points to wrong location: %s\n", file.OriginalPath, target)
						problems++
						continue
					}
				}
			} else {
				fmt.Printf("⚠️  %s - Exists but is not a symlink (manual intervention required)\n", file.OriginalPath)
				problems++
				continue
			}
		}

		// File is missing or broken symlink - can be fixed
		fmt.Printf("🔧 %s - Missing symlink", file.OriginalPath)
		if dryRun {
			fmt.Printf(" (would fix)\n")
		} else {
			// Remove broken symlink if it exists
			if fileops.PathExists(file.OriginalPath) {
				os.Remove(file.OriginalPath)
			}

			// Create new symlink
			if err := fileops.CreateSymlink(file.OriginalPath, repoPath); err != nil {
				fmt.Printf(" - Failed to fix: %v\n", err)
				problems++
				continue
			}
			fmt.Printf(" - Fixed!\n")
		}
		fixed++
	}

	if fixed > 0 {
		if dryRun {
			fmt.Printf("Would fix %d file(s)\n", fixed)
		} else {
			fmt.Printf("Fixed %d file(s)\n", fixed)
		}
	}

	return nil
}

// getManagedDirectories returns a list of all managed directory paths
func getManagedDirectories(idx *types.Index) []string {
	var dirs []string
	for _, file := range index.GetAllFiles(idx) {
		if file.Type == types.FileTypeDirectory {
			dirs = append(dirs, file.OriginalPath)
		}
	}
	return dirs
}

// isWithinManagedDirectory checks if a file path is within any of the managed directories
func isWithinManagedDirectory(filePath string, managedDirs []string) bool {
	for _, dir := range managedDirs {
		// Check if the file path starts with the directory path followed by a separator
		if strings.HasPrefix(filePath, dir+"/") || strings.HasPrefix(filePath, dir+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

// runCleanup removes redundant file entries that are covered by managed directories
func runCleanup(dryRun bool) error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	if index.Count(idx) == 0 {
		fmt.Println("No files are managed by dotman.")
		return nil
	}

	// Get all managed directories
	managedDirs := getManagedDirectories(idx)
	if len(managedDirs) == 0 {
		fmt.Println("No managed directories found - nothing to clean up.")
		return nil
	}

	// Find redundant file entries
	var redundantFiles []types.ManagedFile
	for _, file := range index.GetAllFiles(idx) {
		if file.Type == types.FileTypeFile && isWithinManagedDirectory(file.OriginalPath, managedDirs) {
			redundantFiles = append(redundantFiles, file)
		}
	}

	if len(redundantFiles) == 0 {
		fmt.Println("No redundant file entries found.")
		return nil
	}

	fmt.Printf("Found %d redundant file entries covered by managed directories:\n", len(redundantFiles))
	for _, file := range redundantFiles {
		fmt.Printf("  %s (covered by parent directory)\n", file.OriginalPath)
	}

	if dryRun {
		fmt.Println("\nDry-run mode: would remove these entries from the index")
		return nil
	}

	// Remove redundant entries from the index
	removed := 0
	var removedPaths []string
	for _, file := range redundantFiles {
		if index.RemoveFile(idx, file.OriginalPath) {
			removed++
			// Convert to $HOME relative path
			if homeRelPath, err := config.RelativeToHome(cfg, file.OriginalPath); err == nil {
				removedPaths = append(removedPaths, "$HOME/"+homeRelPath)
			} else {
				removedPaths = append(removedPaths, file.OriginalPath)
			}
			fmt.Printf("Removed %s from index\n", file.OriginalPath)
		}
	}

	if removed == 0 {
		fmt.Println("No entries were removed from the index.")
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

	// Create commit message with directory count
	dirCount := len(managedDirs)
	var commitMsg string
	if dirCount == 1 {
		commitMsg = fmt.Sprintf("Cleanup: remove %d redundant entries covered by 1 directory", removed)
	} else {
		commitMsg = fmt.Sprintf("Cleanup: remove %d redundant entries covered by %d directories", removed, dirCount)
	}

	if err := git.Commit(cfg.DotmanDir, commitMsg); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Printf("Successfully cleaned up %d redundant file entries\n", removed)
	return nil
}
