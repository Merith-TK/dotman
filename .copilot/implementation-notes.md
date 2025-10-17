# Implementation Notes

## Current Progress
- [x] Initial planning completed
- [x] Project structure defined
- [x] Go module initialization
- [x] Basic CLI setup with cobra
- [x] Core types and interfaces implemented
- [x] File operations package
- [x] Index management
- [x] Git integration
- [x] CLI commands (add, remove, status, deploy, init, clone)
- [x] Configuration management with HOME-only security
- [x] README documentation
- [x] Command specifications document
- [x] Working implementation with all core features

## Testing Results
- [x] Build succeeds without errors
- [x] CLI help system works
- [x] `dotman init` creates repository correctly
- [x] `dotman status` shows proper status
- [x] Path validation rejects files outside HOME
- [x] `dotman add` moves files and creates symlinks
- [x] `dotman remove` restores files correctly
- [x] Symlinks work and preserve content

## Key Decisions Made

### 1. Index File Format
- Using JSON for simplicity and readability
- Storing both original and repo paths for flexibility
- Including metadata like type and date added

### 2. Directory Structure
- Following Go project layout standards
- Separating internal packages from public API
- Using cmd/ for executables, internal/ for private code

### 3. File Operations Strategy
- Move files to `.dotman` first, then create symlinks
- Preserve directory structure in repo
- Handle both files and directories uniformly

## Next Steps
1. Initialize Go module
2. Create basic project structure
3. Implement core types and interfaces
4. Build CLI framework with cobra
5. Implement file operations

## Technical Decisions

### Path Handling
- Use `filepath.Abs()` for absolute paths
- Expand `~` to user home directory
- Store relative paths in index for portability

### Symlink Strategy
- Always create absolute symlinks for reliability
- Check if target already exists before creating
- Handle broken symlinks gracefully

### Git Integration
- Initialize repo automatically on first use
- Auto-commit with descriptive messages
- Support for custom commit messages later

## Error Scenarios to Handle
1. Target file/directory already exists
2. Permission denied on file operations
3. Git repository in inconsistent state
4. Broken symlinks in original locations
5. Index file corruption
6. Partial failures during batch operations

## Next Steps - Future Features

### Deploy Scripts (Planned)
- Support for dynamic path deployment (Firefox profiles, browser extensions)
- Script-based deployment for complex scenarios
- Pattern matching for generated directory names
- Cross-platform deployment handling

### Potential Implementation Approaches
1. **Shell Scripts**: Simple executable scripts in `.dotman/scripts/`
2. **YAML Configuration**: Declarative deployment definitions
3. **Go Plugins**: Native Go-based deployment handlers

### Use Case: Firefox userChrome.css
- Problem: Firefox profiles have random names like `4yi8ybqt.default-release`
- Solution: Deploy script that discovers active profile and creates symlink
- Command: `dotman deploy --scripts` or `dotman script run firefox`
