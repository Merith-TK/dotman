package index

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Merith-TK/dotman/pkg/types"
)

// Load reads and parses the index.json file
func Load(indexPath string) (*types.Index, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty index if file doesn't exist
			return &types.Index{
				Version:      "1.0",
				ManagedFiles: make([]types.ManagedFile, 0),
			}, nil
		}
		return nil, fmt.Errorf("failed to read index file: %w", err)
	}

	var index types.Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index file: %w", err)
	}

	return &index, nil
}

// Save writes the index to the index.json file
func Save(index *types.Index, indexPath string) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// AddFile adds a managed file to the index
func AddFile(idx *types.Index, originalPath, repoPath string, fileType types.FileType) {
	managedFile := types.ManagedFile{
		OriginalPath: originalPath,
		RepoPath:     repoPath,
		Type:         fileType,
		AddedDate:    time.Now(),
	}
	
	idx.ManagedFiles = append(idx.ManagedFiles, managedFile)
}

// RemoveFile removes a managed file from the index by original path
func RemoveFile(idx *types.Index, originalPath string) bool {
	for i, file := range idx.ManagedFiles {
		if file.OriginalPath == originalPath {
			// Remove element at index i
			idx.ManagedFiles = append(idx.ManagedFiles[:i], idx.ManagedFiles[i+1:]...)
			return true
		}
	}
	return false
}

// FindFile finds a managed file by original path
func FindFile(idx *types.Index, originalPath string) (*types.ManagedFile, bool) {
	for _, file := range idx.ManagedFiles {
		if file.OriginalPath == originalPath {
			return &file, true
		}
	}
	return nil, false
}

// IsManaged checks if a path is already managed
func IsManaged(idx *types.Index, originalPath string) bool {
	_, found := FindFile(idx, originalPath)
	return found
}

// GetAllFiles returns all managed files
func GetAllFiles(idx *types.Index) []types.ManagedFile {
	return idx.ManagedFiles
}

// Count returns the number of managed files
func Count(idx *types.Index) int {
	return len(idx.ManagedFiles)
}
