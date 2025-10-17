# Dotman - Dotfiles Manager Planning & Status

## Project Status: ✅ COMPLETE & PRODUCTION READY

Dotman is a fully implemented, secure dotfiles manager written in Go that centralizes dotfiles in a git repository while maintaining their original locations through symlinks.

## Core Concept ✅ IMPLEMENTED
- Create `.dotman` folder in user's home directory ✅
- Move actual files/folders to `.dotman/<original-path>` ✅
- Create symlinks from original locations to `.dotman` locations ✅  
- Track managed files in JSON index file ✅
- **Security First**: HOME directory only restriction ✅

## Implemented Features ✅

### 1. Repository Management
- ✅ `dotman init` - Initialize new dotman repository in ~/.dotman
- ✅ `dotman clone <url>` - Clone existing dotfiles repository to ~/.dotman

### 2. File/Folder Management  
- ✅ `dotman add <path>...` - Add multiple files/folders (HOME only restriction)
- ✅ `dotman remove <path>...` - Remove multiple items from management
- ✅ `dotman deploy` - Deploy from index file (for new systems)
- ✅ `dotman status` - Show managed files status with smart directory handling

### 3. Advanced Status Operations
- ✅ `dotman status --sync` - Auto-discover and add unmanaged files from repo
- ✅ `dotman status --fix` - Fix broken or missing symlinks  
- ✅ `dotman status --cleanup` - Remove redundant file entries covered by directories
- ✅ `dotman status --dry-run` - Preview changes without applying

### 4. Security & Safety
- ✅ **HOME Directory Only**: Strict enforcement - files outside $HOME cannot be managed
- ✅ **Path Validation**: Prevents directory traversal and external access
- ✅ **Conflict Detection**: Checks for existing files before operations
- ✅ **Atomic Operations**: Rollback on failure during operations
- ✅ **Multi-operation Error Handling**: Continue processing on partial failures

### 5. Smart Directory Management ✅ 
- ✅ **Hierarchy Detection**: Recognizes when files are within managed directories
- ✅ **Cleanup Intelligence**: Removes redundant individual file entries
- ✅ **Status Filtering**: Only shows top-level managed items in status

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

## Implementation Status: ALL PHASES COMPLETE ✅

### Phase 1: Core Functionality ✅ COMPLETE
- ✅ Project structure setup
- ✅ Basic CLI with cobra  
- ✅ File operations (move, symlink)
- ✅ Index file management
- ✅ `dotman add` command with multi-file support
- ✅ `dotman status` command with advanced features
- ✅ `dotman remove` command with multi-file support
- ✅ `dotman deploy` command with sync option
- ✅ `dotman init` command
- ✅ `dotman clone` command

### Phase 2: Enhanced Security & Validation ✅ COMPLETE
- ✅ HOME directory restriction enforcement (strict security model)
- ✅ Path validation improvements with clear error messages
- ✅ Better error messages for invalid paths with actionable suggestions

### Phase 3: Git Integration ✅ COMPLETE  
- ✅ Git repository initialization with .gitignore
- ✅ Auto-commit on add/remove with descriptive $HOME/ path messages
- ✅ Git status integration in status command
- ✅ Enhanced commit messages for all operations

### Phase 4: Advanced Features ✅ COMPLETE
- ✅ Smart directory hierarchy management
- ✅ Cleanup functionality for redundant entries
- ✅ Symlink repair functionality  
- ✅ Auto-discovery of unmanaged files
- ✅ Dry-run mode for all operations
- ✅ Multi-operation error handling
- ✅ Comprehensive status reporting

### Phase 5: Enhanced User Experience ✅ COMPLETE
- ✅ Multi-file operations (add/remove multiple files at once)
- ✅ Smart commit messages with actual paths instead of generic counts
- ✅ Directory hierarchy intelligence (files within managed dirs aren't shown as separate)
- ✅ Comprehensive help system and documentation
- ✅ Error handling with actionable feedback

### Future Considerations (Not Currently Planned)
- [ ] Deploy scripts for dynamic paths (Firefox profiles, browser extensions)
- [ ] Multiple profile support (mentioned in original requirements but not implemented)
- [ ] Cross-platform support (currently Unix-like systems)
- [ ] Configuration file support
- [ ] Comprehensive testing suite

## Error Handling Considerations
- Existing symlinks at target locations
- Permission issues
- Git repository conflicts
- Partial failures during operations
- Missing original files during deploy

## Implemented Architecture ✅

### Dependencies Used
- ✅ `github.com/spf13/cobra` - CLI framework for command structure
- ✅ Shell `git` commands - Git operations via exec (simpler than go-git)
- ✅ Standard library for all file operations
- ✅ JSON encoding for index management

### Security Implementation ✅
- ✅ Strict path validation to prevent directory traversal
- ✅ HOME directory only restriction enforced at all entry points  
- ✅ Safe absolute vs relative paths handling with config.ExpandPath()
- ✅ File permissions preserved during operations
- ✅ Atomic operations with rollback on failure
- ✅ Symlink validation and repair functionality

### Key Technical Decisions ✅
- ✅ **Function-based API**: Used functions instead of methods (Go package restrictions)
- ✅ **Security First**: All path operations validated through config package
- ✅ **Simple Git Integration**: Shell commands instead of go-git library
- ✅ **JSON Index**: Human-readable tracking with timestamps
- ✅ **Comprehensive Error Handling**: Clear messages with actionable suggestions

## Current Feature Set Summary

### Core Commands
- `dotman init` - Initialize new repository
- `dotman clone <url>` - Clone existing repository  
- `dotman add <paths>...` - Add multiple files/directories
- `dotman remove <paths>...` - Remove multiple files/directories
- `dotman status [flags]` - Status with sync/fix/cleanup options
- `dotman deploy [flags]` - Deploy with optional sync

### Status Command Features
- `--sync, -s` - Auto-discover unmanaged files
- `--fix, -f` - Repair broken symlinks
- `--cleanup, -c` - Remove redundant entries
- `--dry-run, -n` - Preview mode for all operations

### Smart Features  
- Directory hierarchy detection
- Redundant entry cleanup
- Batch operation support
- Enhanced git commit messages
- Multi-operation error handling
- Dry-run preview mode

## Future Enhancement Ideas (Beyond Current Scope)

### Multi-Profile Support (Original Request)
The user mentioned wanting multiple repo support for different machines:
```
.dotman/profiles/github.com/Merith-TK/dotfiles.chromebook
```

This would require:
- Profile management system
- Profile switching commands
- Profile-specific index files
- Enhanced repository structure

### Deploy Scripts (Documented but Not Implemented)
For handling dynamic paths like Firefox profiles:
- Script-based deployment system
- Pattern matching for generated directory names  
- Cross-platform deployment handling

These features could be added in future versions but are not currently implemented.
