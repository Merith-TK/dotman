package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Merith-TK/dotman/pkg/types"
)

const (
	DotmanDirName  = ".dotman"
	IndexFileName  = "index.json"
	DefaultVersion = "1.0"
)

// New creates a new Config with default values
func New() (*types.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dotmanDir := filepath.Join(homeDir, DotmanDirName)
	indexFile := filepath.Join(dotmanDir, IndexFileName)

	return &types.Config{
		DotmanDir: dotmanDir,
		HomeDir:   homeDir,
		IndexFile: indexFile,
	}, nil
}

// EnsureDotmanDir creates the .dotman directory if it doesn't exist
func EnsureDotmanDir(cfg *types.Config) error {
	return os.MkdirAll(cfg.DotmanDir, 0755)
}

// DotmanDirExists checks if the .dotman directory exists
func DotmanDirExists(cfg *types.Config) bool {
	_, err := os.Stat(cfg.DotmanDir)
	return err == nil
}

// IndexFileExists checks if the index.json file exists
func IndexFileExists(cfg *types.Config) bool {
	_, err := os.Stat(cfg.IndexFile)
	return err == nil
}

// ExpandPath expands ~ to the user's home directory and returns absolute path
// Only allows paths within the home directory for security
func ExpandPath(cfg *types.Config, path string) (string, error) {
	var expandedPath string

	if path == "~" {
		expandedPath = cfg.HomeDir
	} else if filepath.HasPrefix(path, "~/") {
		expandedPath = filepath.Join(cfg.HomeDir, path[2:])
	} else {
		// Convert to absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		expandedPath = absPath
	}

	// Security check: ensure path is within home directory
	if !IsInsideHome(cfg, expandedPath) {
		return "", fmt.Errorf("path must be inside home directory: %s", expandedPath)
	}

	return expandedPath, nil
}

// RelativeToHome returns the path relative to the user's home directory
func RelativeToHome(cfg *types.Config, absolutePath string) (string, error) {
	// Ensure the path is inside home first
	if !IsInsideHome(cfg, absolutePath) {
		return "", fmt.Errorf("path is outside home directory: %s", absolutePath)
	}

	return filepath.Rel(cfg.HomeDir, absolutePath)
}

// IsInsideHome checks if the given path is inside the user's home directory
func IsInsideHome(cfg *types.Config, path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Ensure home directory is absolute
	absHome, err := filepath.Abs(cfg.HomeDir)
	if err != nil {
		return false
	}

	rel, err := filepath.Rel(absHome, absPath)
	if err != nil {
		return false
	}

	// If the relative path starts with "..", it's outside the home directory
	return !filepath.HasPrefix(rel, "..")
}

// ShouldIgnoreRepoPath returns true if the given repo-relative path refers to
// metadata that should never be tracked or deployed by dotman.
// We hardcode ignoring the .dotman directory and README.md (case-insensitive).
func ShouldIgnoreRepoPath(cfg *types.Config, repoRelPath string) bool {
	// Normalize path separators
	rel := filepath.Clean(repoRelPath)

	// Ignore dotman metadata directory at repo root
	if rel == ".dotman" || rel == ".dotman"+string(filepath.Separator) {
		return true
	}

	// Ignore README.md (case-insensitive) at repo root
	lower := strings.ToLower(rel)
	if lower == "readme.md" {
		return true
	}

	// If the path is inside a .dotman directory, ignore it
	if strings.HasPrefix(rel, ".dotman"+string(filepath.Separator)) {
		return true
	}

	return false
}
