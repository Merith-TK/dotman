# Dotman - Dotfiles Manager Planning

## Overview
Dotman is a Go-based dotfiles manager that centralizes dotfiles in a git repository while maintaining their original locations through symlinks.

## Core Concept
- Create `.dotman` folder in user's home directory
- Move actual files/folders to `.dotman/<original-path>`
- Create symlinks from original locations to `.dotman` locations
- Track managed files in an index file

## Key Features

### 1. Repository Management
- `dotman init` - Initialize a new dotman repository in ~/.dotman
- `dotman clone <url>` - Clone existing dotfiles repository to ~/.dotman

### 2. File/Folder Management
- `dotman add <path>` - Add file/folder to management (HOME only)
- `dotman remove <path>` - Remove from management (restore original)
- `dotman deploy` - Deploy from index file (for new systems)
- `dotman status` - Show managed files status

### 3. Security & Safety
- **HOME Directory Only**: All managed files must be within $HOME
- **Path Validation**: Prevents directory traversal and external access
- **Conflict Detection**: Checks for existing files before operations

### 2. Repository Structure
```
~/.dotman/
├── .git/                    # Git repository
├── index.json              # Tracking file for managed dotfiles
├── .config/                # Mirrored structure
│   └── sway/
├── .bashrc
└── .vimrc
```

### 3. Index File Format (JSON)
```json
{
  "version": "1.0",
  "managed_files": [
    {
      "original_path": "~/.config/sway",
      "repo_path": ".config/sway",
      "type": "directory",
      "added_date": "2025-08-01T10:00:00Z"
    },
    {
      "original_path": "~/.bashrc",
      "repo_path": ".bashrc", 
      "type": "file",
      "added_date": "2025-08-01T10:05:00Z"
    }
  ]
}
```

## Technical Architecture

### Core Components
1. **Config Manager** - Handle `.dotman` directory setup
2. **File Operations** - Move files, create symlinks
3. **Index Manager** - Read/write index.json
4. **Git Integration** - Initialize repo, commit changes
5. **CLI Interface** - Command parsing and execution

### Go Packages Structure
```
dotman/
├── cmd/
│   └── dotman/
│       └── main.go          # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── fileops/
│   │   └── fileops.go       # File operations (move, symlink)
│   ├── index/
│   │   └── index.go         # Index file management
│   ├── git/
│   │   └── git.go           # Git operations
│   └── cli/
│       └── commands.go      # Command implementations
├── pkg/
│   └── types/
│       └── types.go         # Shared types and structs
├── go.mod
├── go.sum
└── README.md
```

## Implementation Phases

### Phase 1: Core Functionality
- [x] Project structure setup
- [x] Basic CLI with cobra
- [x] File operations (move, symlink)
- [x] Index file management
- [x] `dotman add` command
- [x] `dotman status` command
- [x] `dotman remove` command
- [x] `dotman deploy` command
- [ ] `dotman init` command
- [ ] `dotman clone` command

### Phase 2: Enhanced Security & Validation
- [x] HOME directory restriction enforcement
- [ ] Path validation improvements
- [ ] Better error messages for invalid paths

### Phase 2: Git Integration
- [ ] Git repository initialization
- [ ] Auto-commit on add/remove
- [ ] Git status integration

### Phase 3: Advanced Features
- [ ] `dotman deploy` command
- [ ] `dotman status` command
- [ ] `dotman remove` command
- [ ] Conflict resolution
- [ ] Backup before operations

### Phase 4: Advanced Features
- [ ] Deploy scripts for dynamic paths
- [ ] `dotman remove` command
- [ ] `dotman status` command
- [ ] Conflict resolution
- [ ] Backup before operations

### Phase 5: Deploy Scripts & Dynamic Paths
- [ ] Deploy script configuration system
- [ ] Support for pattern-based file discovery
- [ ] Firefox profile handling (userChrome.css, etc.)
- [ ] Browser extension deployment
- [ ] Application-specific deployment handlers
- [ ] Script validation and error handling

### Phase 6: Polish
- [ ] Error handling improvements
- [ ] Configuration file support
- [ ] Multiple profiles support
- [ ] Testing suite

## Error Handling Considerations
- Existing symlinks at target locations
- Permission issues
- Git repository conflicts
- Partial failures during operations
- Missing original files during deploy

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/go-git/go-git/v5` - Git operations (or shell git)
- Standard library for file operations

## Security Considerations
- Validate paths to prevent directory traversal
- Handle absolute vs relative paths safely
- Preserve file permissions
- Backup strategy for destructive operations
