package types

import "time"

// ManagedFile represents a file or directory managed by dotman
type ManagedFile struct {
	OriginalPath string    `json:"original_path"` // Original location (e.g., ~/.config/sway)
	RepoPath     string    `json:"repo_path"`     // Path within .dotman repo (e.g., .config/sway)
	Type         FileType  `json:"type"`          // file or directory
	AddedDate    time.Time `json:"added_date"`    // When it was added to management
}

// FileType represents whether the managed item is a file or directory
type FileType string

const (
	FileTypeFile      FileType = "file"
	FileTypeDirectory FileType = "directory"
)

// Index represents the dotman index file structure
type Index struct {
	Version      string        `json:"version"`
	ManagedFiles []ManagedFile `json:"managed_files"`
}

// Config represents dotman configuration
type Config struct {
	DotmanDir string // Path to .dotman directory (usually ~/.dotman)
	HomeDir   string // User's home directory
	IndexFile string // Path to index.json file
}

// Operation represents a file operation result
type Operation struct {
	Success bool
	Path    string
	Error   error
	Message string
}

// AddOptions represents options for the add command
type AddOptions struct {
	Force   bool // Force operation even if conflicts exist
	DryRun  bool // Show what would happen without doing it
	Backup  bool // Create backup before operation
	Message string // Custom commit message
}

// DeployOptions represents options for the deploy command
type DeployOptions struct {
	Force   bool // Force deployment even if conflicts exist
	DryRun  bool // Show what would happen without doing it
	Backup  bool // Create backup before operation
}
