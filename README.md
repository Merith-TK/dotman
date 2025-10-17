# Dotman - Dotfiles Manager

A simple, efficient dotfiles manager written in Go that centralizes your configuration files in a git repository while maintaining their original locations through symlinks.

## Features

- **Centralized Management**: All your dotfiles are stored in `~/.dotman` git repository
- **Symlink Strategy**: Original file locations are preserved via symlinks
- **Git Integration**: Automatic git commits for tracking changes
- **File & Directory Support**: Manage both individual files and entire directories
- **Simple CLI**: Easy-to-use command-line interface
- **Safe Operations**: Built-in safeguards and validation

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
# Add your first dotfile
dotman add ~/.bashrc

# Add a configuration directory
dotman add ~/.config/sway

# Check status of managed files
dotman status

# Deploy on a new system (after cloning your dotfiles repo)
dotman deploy
```

## Commands

### `dotman add <path>`
Adds a file or directory to dotman management.

- Moves the file/directory to `~/.dotman/`
- Creates a symlink in the original location
- Updates the index and commits to git

```bash
dotman add ~/.config/nvim
dotman add ~/.bashrc
dotman add ~/.ssh/config
```

### `dotman status`
Shows the status of all managed files.

```bash
dotman status
```

### `dotman remove <path>`
Removes a file or directory from dotman management.

- Removes the symlink
- Restores the original file from the repository
- Updates the index and commits to git

```bash
dotman remove ~/.config/nvim
```

### `dotman deploy`
Deploys all managed files by creating symlinks from the repository.

Useful when setting up dotfiles on a new system after cloning your `.dotman` repository.

```bash
dotman deploy
```

## How It Works

1. **Initialization**: On first use, dotman creates `~/.dotman/` and initializes it as a git repository

2. **Adding Files**: When you add a file:
   ```
   ~/.bashrc (original) → ~/.dotman/.bashrc (moved)
   ~/.bashrc (symlink) → ~/.dotman/.bashrc (points to)
   ```

3. **Index Tracking**: An `index.json` file tracks all managed files:
   ```json
   {
     "version": "1.0",
     "managed_files": [
       {
         "original_path": "/home/user/.bashrc",
         "repo_path": ".bashrc",
         "type": "file",
         "added_date": "2025-08-01T10:00:00Z"
       }
     ]
   }
   ```

4. **Git Integration**: All changes are automatically committed with descriptive messages

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

## Safety Features

- **Path Validation**: Only allows files within the home directory
- **Conflict Detection**: Checks for existing files before operations
- **Symlink Verification**: Validates symlinks during status checks
- **Atomic Operations**: Rolls back on failure during add operations

## Use Cases

- **Personal Dotfiles**: Manage your personal configuration files
- **Development Environment**: Sync development tools and editor configs
- **System Setup**: Quickly deploy your environment on new machines
- **Backup & Versioning**: Keep a versioned backup of all your configurations

## Requirements

- Go 1.19 or later
- Git installed and available in PATH
- Unix-like system (Linux, macOS)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
