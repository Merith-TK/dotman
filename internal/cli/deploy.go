package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Merith-TK/dotman/internal/config"
	"github.com/Merith-TK/dotman/internal/fileops"
	"github.com/Merith-TK/dotman/internal/index"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy managed files",
	Long: `Deploy creates symlinks for all managed files.
Useful when setting up dotfiles on a new system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeploy()
	},
}

func runDeploy() error {
	if !config.DotmanDirExists(cfg) {
		return fmt.Errorf("dotman directory does not exist: %s", cfg.DotmanDir)
	}

	idx, err := index.Load(cfg.IndexFile)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	if index.Count(idx) == 0 {
		fmt.Println("No files to deploy.")
		return nil
	}

	fmt.Printf("Deploying %d file(s)...\n", index.Count(idx))

	for _, file := range index.GetAllFiles(idx) {
		repoPath := filepath.Join(cfg.DotmanDir, file.RepoPath)

		// Skip repository metadata
		if config.ShouldIgnoreRepoPath(cfg, file.RepoPath) {
			fmt.Printf("Skipping repository metadata: %s\n", file.RepoPath)
			continue
		}

		// Check if repo file exists
		if !fileops.PathExists(repoPath) {
			fmt.Printf("Warning: repo file missing for %s\n", file.OriginalPath)
			continue
		}

		// Check if original location already exists
		if fileops.PathExists(file.OriginalPath) {
			if fileops.IsSymlink(file.OriginalPath) {
				fmt.Printf("Skipping %s (symlink already exists)\n", file.OriginalPath)
				continue
			} else {
				fmt.Printf("Warning: %s exists and is not a symlink, skipping\n", file.OriginalPath)
				continue
			}
		}

		// Create symlink
		if err := fileops.CreateSymlink(file.OriginalPath, repoPath); err != nil {
			fmt.Printf("Error creating symlink for %s: %v\n", file.OriginalPath, err)
			continue
		}

		fmt.Printf("Deployed %s\n", file.OriginalPath)
	}

	fmt.Println("Deployment complete.")
	return nil
}
