package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InitRepo initializes a git repository in the specified directory
func InitRepo(repoPath string) error {
	cmd := exec.Command("git", "init", "-b", "main")
	cmd.Dir = repoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize git repo: %s, %w", string(output), err)
	}

	return nil
}

// IsGitRepo checks if the directory is a git repository
func IsGitRepo(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// Add stages files for commit
func Add(repoPath string, files ...string) error {
	if len(files) == 0 {
		// Add all files
		files = []string{"."}
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add files to git: %s, %w", string(output), err)
	}

	return nil
}

// Commit creates a commit with the specified message
func Commit(repoPath, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		// Check if the error is because there's nothing to commit
		if cmd.ProcessState.ExitCode() == 1 {
			// This might be "nothing to commit" which is not really an error
			return nil
		}
		return fmt.Errorf("failed to commit: %s, %w", string(output), err)
	}

	return nil
}

// Status returns the git status
func Status(repoPath string) (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}

	return string(output), nil
}

// HasChanges checks if there are any uncommitted changes
func HasChanges(repoPath string) (bool, error) {
	status, err := Status(repoPath)
	if err != nil {
		return false, err
	}

	return len(status) > 0, nil
}

// CreateGitignore creates a basic .gitignore file
func CreateGitignore(repoPath string) error {
	gitignoreContent := `# Dotman specific ignores
.DS_Store
Thumbs.db
*.tmp
*.swp
*.swo
*~

# Don't ignore the index file
!index.json
`

	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// Pull pulls changes from the remote repository
func Pull(repoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull from remote: %s, %w", string(output), err)
	}

	return nil
}

// Push pushes changes to the remote repository
func Push(repoPath string) error {
	cmd := exec.Command("git", "push")
	cmd.Dir = repoPath

	if output, err := cmd.CombinedOutput(); err != nil {
		// Check if this is the first push that needs upstream setup
		if strings.Contains(string(output), "no upstream branch") {
			// Get current branch and set upstream
			branch, branchErr := GetCurrentBranch(repoPath)
			if branchErr != nil {
				return fmt.Errorf("failed to push to remote: %s, %w", string(output), err)
			}

			// Push with set-upstream
			upstreamCmd := exec.Command("git", "push", "--set-upstream", "origin", branch)
			upstreamCmd.Dir = repoPath

			if upstreamOutput, upstreamErr := upstreamCmd.CombinedOutput(); upstreamErr != nil {
				return fmt.Errorf("failed to push to remote: %s, %w", string(upstreamOutput), upstreamErr)
			}

			return nil
		}
		return fmt.Errorf("failed to push to remote: %s, %w", string(output), err)
	}

	return nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := string(output)
	if len(branch) > 0 && branch[len(branch)-1] == '\n' {
		branch = branch[:len(branch)-1] // Remove trailing newline
	}

	return branch, nil
}

// GetRemoteURL returns the remote origin URL
func GetRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("no remote origin configured")
	}

	remoteURL := string(output)
	if len(remoteURL) > 0 && remoteURL[len(remoteURL)-1] == '\n' {
		remoteURL = remoteURL[:len(remoteURL)-1] // Remove trailing newline
	}

	return remoteURL, nil
}

// GetCommitCount returns the number of commits in the repository
func GetCommitCount(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "0", nil // No commits yet
	}

	count := string(output)
	if len(count) > 0 && count[len(count)-1] == '\n' {
		count = count[:len(count)-1] // Remove trailing newline
	}

	return count, nil
}

// EnsureRepo ensures a git repository exists and is properly initialized
func EnsureRepo(repoPath string) error {
	if !IsGitRepo(repoPath) {
		if err := InitRepo(repoPath); err != nil {
			return err
		}

		// Create initial .gitignore
		if err := CreateGitignore(repoPath); err != nil {
			return err
		}

		// Make initial commit
		if err := Add(repoPath); err != nil {
			return err
		}

		if err := Commit(repoPath, "Initial dotman repository"); err != nil {
			return err
		}
	}

	return nil
}
