# Dotman - Dotfiles Manager

A simple, efficient dotfiles manager written in Go that centralizes your configuration files in a git repository while maintaining their original locations through symlinks.

## Features

- **Centralized Management**: All dotfiles stored in `~/.dotman` git repository
- **Symlink Strategy**: Files stay in original locations, linked to central repo
- **Git Integration**: Automatic commits for all changes
- **File & Directory Support**: Manage individual files or entire directories
- **Security**: Only manages files within your home directory
- **Cleanup Tools**: Remove redundant entries and fix broken symlinks
- **Batch Operations**: Add or remove multiple files at once
- **Easy Deployment**: Set up dotfiles on new systems quickly

## Installation

```bash
# Build from source
git clone https://github.com/merith-tk/dotman
cd dotman
go build -o dotman cmd/dotman/main.go

# Install to system
sudo mv dotman /usr/local/bin/
```

## Quick Start

```bash
# Initialize a new dotman repository
dotman init

# Add your first dotfiles
dotman add ~/.bashrc ~/.bash_aliases

# Add a configuration directory
dotman add ~/.config/sway

# Check status of managed files
dotman status

# Clean up redundant entries (files covered by managed directories)
dotman status --cleanup

# On a new system: clone your dotfiles and deploy
dotman clone https://github.com/user/dotfiles.git
dotman deploy

# Auto-discover and add unmanaged files from your repo
dotman status --sync

# Fix broken or missing symlinks
dotman status --fix
```

## Commands

### `dotman init`
Initialize a new dotman repository in `~/.dotman`.

- Creates `~/.dotman` directory and initializes git repository
- Creates initial `index.json` and `.gitignore` files
- Makes initial commit

```bash
dotman init
```

### `dotman clone <repository-url>`
Clone an existing dotfiles repository to `~/.dotman`.

- Downloads remote dotfiles repository
- Validates repository structure
- Use `dotman deploy` after cloning to create symlinks

```bash
dotman clone https://github.com/user/dotfiles.git
```

### `dotman add <path>...`
Add files or directories to dotman management.

- Supports multiple paths in one command
- Moves files/directories to `~/.dotman/`
- Creates symlinks in original locations
- Updates index and commits with descriptive `$HOME/` paths

```bash
dotman add ~/.config/nvim ~/.bashrc ~/.ssh/config
dotman add ~/.config/sway    # Manages entire directory
```

### `dotman status [flags]`
Show status of all managed files with enhanced options.

**Flags:**
- `--sync, -s`: Auto-discover and add unmanaged files from repo
- `--fix, -f`: Fix broken or missing symlinks
- `--cleanup, -c`: Remove redundant file entries covered by directories
- `--dry-run, -n`: Show what would be done without doing it

```bash
dotman status                    # Basic status
dotman status --sync            # Discover unmanaged files
dotman status --cleanup         # Remove redundant entries
dotman status --fix --dry-run   # Preview symlink repairs
```

### `dotman remove <path>...`
Remove files or directories from dotman management.

- Supports multiple paths in one command
- Removes symlinks and restores original files
- Updates index and commits changes

```bash
dotman remove ~/.config/nvim ~/.old-config
```

### `dotman deploy [flags]`
Deploy all managed files by creating symlinks.

**Flags:**
- `--sync, -s`: Auto-discover and add unmanaged files before deploying

Perfect for setting up dotfiles on new systems.

```bash
dotman deploy            # Deploy tracked files only
dotman deploy --sync     # Discover and deploy all repo files
```

## How It Works

1. **Security First**: All operations are restricted to your `$HOME` directory - files outside home cannot be managed

2. **Repository Structure**: Dotman creates `~/.dotman/` and mirrors your home directory structure
   ```
   ~/.bashrc (original) → ~/.dotman/.bashrc (moved) ← ~/.bashrc (symlink)
   ~/.config/sway/ (original) → ~/.dotman/.config/sway/ (moved) ← ~/.config/sway/ (symlink)
   ```

3. **Smart Directory Management**: When you manage `~/.config/sway/`, individual files within it are automatically covered
   - Individual files like `~/.config/sway/config` don't need separate tracking
   - The `--cleanup` flag removes redundant individual file entries

4. **Index Tracking**: A `index.json` file tracks all managed items:
   ```json
   {
     "version": "1.0",
     "managed_files": [
       {
         "original_path": "/home/user/.config/sway",
         "repo_path": ".config/sway",
         "type": "directory",
         "added_date": "2025-10-16T10:00:00Z"
       },
       {
         "original_path": "/home/user/.bashrc",
         "repo_path": ".bashrc",
         "type": "file",
         "added_date": "2025-10-16T10:05:00Z"
       }
     ]
   }
   ```

5. **Git Integration**: All changes are automatically committed with descriptive messages using `$HOME/` paths:
   - `Add $HOME/.config/sway to dotman management`
   - `Remove $HOME/.bashrc from dotman management`
   - `Cleanup: remove 21 redundant entries covered by 3 directories`

## Directory Structure

```
~/.dotman/
├── .git/                    # Git repository
├── index.json              # Managed files index
├── .gitignore              # Generated gitignore
├── .config/                # Mirrored home structure
│   ├── sway/
│   └── nvim/
├── .bashrc
└── .vimrc
```

## Advanced Features

### Smart Directory Hierarchy
- **Intelligent Tracking**: Manages directories as units while supporting individual files
- **Hierarchy Detection**: Recognizes when files are covered by managed directories
- **Cleanup Tools**: `--cleanup` flag removes redundant entries automatically

### Multi-Operation Support
- **Batch Operations**: Add or remove multiple files in single commands
- **Error Handling**: Continues processing remaining files if some operations fail
- **Progress Reporting**: Clear feedback on successes and failures

### Deployment Tools
- **Auto-Discovery**: `--sync` flag finds unmanaged files in your repository
- **Repair Tools**: `--fix` flag repairs broken or missing symlinks
- **Dry-Run Mode**: Preview changes with `--dry-run` before applying

### Enhanced Git Integration
- **Descriptive Commits**: Uses `$HOME/` relative paths in commit messages
- **Smart Commit Messages**: 
  - Individual files: `Add $HOME/.bashrc to dotman management`
  - Multiple files: `Sync: add $HOME/.config/foo, $HOME/.config/bar to index`
  - Cleanup operations: `Cleanup: remove 21 redundant entries covered by 3 directories`

## Safety Features

- **🔒 HOME Directory Only**: Strict path validation prevents managing files outside `$HOME`
- **🛡️ Conflict Detection**: Checks for existing files and symlinks before operations
- **🔗 Symlink Verification**: Validates symlinks during status checks and repairs
- **⚛️ Atomic Operations**: Rolls back changes on failure during add operations
- **🧪 Dry-Run Support**: Preview changes without applying them
- **📊 Error Reporting**: Clear error messages with actionable suggestions

## Use Cases

### Personal Configuration Management
- **Development Environment**: Manage editor configs, shell settings, and development tools
- **Desktop Environment**: Sync window manager, desktop, and application configurations
- **Command Line Tools**: Centralize shell aliases, functions, and tool configurations

### Multi-System Workflows  
- **New System Setup**: Clone and deploy your entire configuration with two commands
- **System Migration**: Easily move your dotfiles between computers
- **Shared Configurations**: Maintain consistent environments across work and personal systems

### Team and Project Collaboration
- **Onboarding**: Share standardized development environment configurations
- **Project Templates**: Distribute project-specific configuration templates
- **Documentation**: Version-controlled configuration with commit history

### Maintenance and Organization
- **Configuration Cleanup**: Remove redundant entries and fix broken configurations
- **Backup Strategy**: Git-based versioning with descriptive commit messages
- **Selective Management**: Choose exactly which files and directories to track

### Example Workflows

**Daily Use:**
```bash
# Add new configuration
dotman add ~/.config/new-app

# Check what's managed and clean up
dotman status --cleanup

# Commit any manual changes in ~/.dotman
cd ~/.dotman && git add . && git commit -m "Update configurations"
```

**New System Setup:**
```bash
# One-time setup
dotman clone https://github.com/user/dotfiles.git
dotman deploy --sync  # Deploy tracked files and discover any untracked ones

# System is ready with all your configurations!
```

**Maintenance:**
```bash
# Fix any broken symlinks
dotman status --fix

# Find files that were added to repo but not tracked
dotman status --sync

# Clean up redundant entries
dotman status --cleanup
```

## Requirements

- Go 1.19 or later
- Git installed and available in PATH
- Unix-like system (Linux, macOS)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
