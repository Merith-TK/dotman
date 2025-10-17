package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Merith-TK/dotman/pkg/types"
)

// MoveToRepo moves a file or directory from its original location to the dotman repo
func MoveToRepo(originalPath, repoPath string) error {
	// Ensure the destination directory exists
	destDir := filepath.Dir(repoPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Move the file/directory
	if err := os.Rename(originalPath, repoPath); err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", originalPath, repoPath, err)
	}

	return nil
}

// CreateSymlink creates a symlink from the original location to the repo location
func CreateSymlink(originalPath, repoPath string) error {
	// Ensure the parent directory of the symlink exists
	parentDir := filepath.Dir(originalPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory for symlink: %w", err)
	}

	// Create the symlink
	if err := os.Symlink(repoPath, originalPath); err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", originalPath, repoPath, err)
	}

	return nil
}

// RemoveSymlink removes a symlink and restores the original file from repo
func RemoveSymlink(originalPath, repoPath string) error {
	// Check if the original path is actually a symlink
	linkInfo, err := os.Lstat(originalPath)
	if err != nil {
		return fmt.Errorf("failed to stat original path: %w", err)
	}

	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("original path is not a symlink")
	}

	// Remove the symlink
	if err := os.Remove(originalPath); err != nil {
		return fmt.Errorf("failed to remove symlink: %w", err)
	}

	// Move the file back from repo to original location
	if err := os.Rename(repoPath, originalPath); err != nil {
		return fmt.Errorf("failed to restore file from repo: %w", err)
	}

	return nil
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsSymlink checks if a path is a symlink
func IsSymlink(path string) bool {
	linkInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return linkInfo.Mode()&os.ModeSymlink != 0
}

// IsDirectory checks if a path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileType determines if a path is a file or directory
func GetFileType(path string) types.FileType {
	if IsDirectory(path) {
		return types.FileTypeDirectory
	}
	return types.FileTypeFile
}

// BackupPath creates a backup of a file or directory by copying it with a .backup suffix
func BackupPath(path string) error {
	backupPath := path + ".backup"
	
	// Remove existing backup if it exists
	if PathExists(backupPath) {
		if err := os.RemoveAll(backupPath); err != nil {
			return fmt.Errorf("failed to remove existing backup: %w", err)
		}
	}

	if IsDirectory(path) {
		return copyDir(path, backupPath)
	}
	return copyFile(path, backupPath)
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
