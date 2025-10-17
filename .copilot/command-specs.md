# Dotman Command Specifications

## Command Overview

Dotman is a HOME-directory-only dotfiles manager that uses symlinks to maintain files in their original locations while storing the actual files in a centralized git repository.

## Security Model

**HOME Directory Only**: All operations are restricted to files within the user's home directory (`$HOME`). Files outside of `$HOME` cannot be managed by dotman for security reasons.

## Commands

### `dotman init`

**Purpose**: Initialize a new dotman repository in `~/.dotman`

**Behavior**:
- Creates `~/.dotman` directory if it doesn't exist
- Initializes git repository on `main` branch
- Creates initial `index.json` file
- Creates `.gitignore` file
- Makes initial commit

**Error Conditions**:
- If `~/.dotman` already exists and is a valid dotman repository, do nothing (success)
- If `~/.dotman` exists but is not a git repository, initialize git repo
- If `~/.dotman` exists but is not empty and not a git repo, fail with error

**Example**:
```bash
dotman init
# Output: Initialized dotman repository in ~/.dotman
```

### `dotman clone <url>`

**Purpose**: Clone an existing dotfiles repository to `~/.dotman`

**Behavior**:
- Clones git repository from URL to `~/.dotman`
- Validates that the cloned repository has a valid `index.json`
- Does NOT automatically deploy files (use `dotman deploy` for that)

**Error Conditions**:
- If `~/.dotman` already exists, fail with error
- If URL is invalid or inaccessible, fail with error
- If cloned repository doesn't contain valid `index.json`, fail with error

**Example**:
```bash
dotman clone https://github.com/user/dotfiles.git
# Output: Cloned dotfiles repository to ~/.dotman
```

### `dotman add <path>`

**Purpose**: Add a file or directory to dotman management

**Flow**:
1. Validate `<path>` is within `$HOME`
2. Expand path to absolute path
3. Check if already managed (error if so)
4. Move file: `mv <path> ~/.dotman/<relative-path>`
5. Create symlink: `ln -s ~/.dotman/<relative-path> <path>`
6. Update `index.json`
7. Git add and commit

**Error Conditions**:
- Path is outside `$HOME`
- Path doesn't exist
- Path is already managed
- Path is a broken symlink

**Example**:
```bash
dotman add ~/.config/sway
# Output: Successfully added ~/.config/sway to dotman management

dotman add /etc/hosts
# Error: path must be inside home directory: /etc/hosts
```

### `dotman remove <path>`

**Purpose**: Remove a file or directory from dotman management

**Flow**:
1. Validate `<path>` is within `$HOME`
2. Expand path to absolute path
3. Check if managed (error if not)
4. Remove symlink: `rm <path>`
5. Restore original: `mv ~/.dotman/<relative-path> <path>`
6. Update `index.json`
7. Git add and commit

**Error Conditions**:
- Path is outside `$HOME`
- Path is not managed by dotman
- Original file missing from repository

**Example**:
```bash
dotman remove ~/.config/sway
# Output: Successfully removed ~/.config/sway from dotman management

dotman remove ~/.unmanaged-file
# Error: path is not managed by dotman: ~/.unmanaged-file
```

### `dotman deploy`

**Purpose**: Deploy all managed files from the repository

**Flow**:
1. Read `index.json`
2. For each managed file:
   - Check if repository file exists
   - Check if target location exists
   - If target is symlink, skip
   - If target is regular file, warn and skip
   - Create symlink from target to repository file

**Error Conditions**:
- Repository not initialized
- Repository files missing
- Permission errors

**Example**:
```bash
dotman deploy
# Output: 
# Deploying 3 file(s)...
# Deployed ~/.bashrc
# Deployed ~/.config/sway
# Skipping ~/.vimrc (symlink already exists)
# Deployment complete.
```

### `dotman fix [path]`

**Purpose**: Fix broken or missing symlinks for managed files

**Flow**:
1. If no path specified, check all managed files
2. For each file:
   - Verify repository file exists
   - Check symlink status (missing, broken, or wrong target)
   - Remove broken symlinks
   - Recreate correct symlinks
3. Report summary of fixes and issues

**Options**:
- `--dry-run`: Show what would be fixed without doing it

**Error Conditions**:
- Repository file missing (needs manual intervention)
- Original location has regular file (needs manual intervention)
- Permission errors

**Example**:
```bash
dotman fix
# Output:
# Checking 3 managed file(s) for issues...
# ✓ ~/.bashrc - OK
# 🔧 ~/.config/sway - Missing symlink - Fixed!
# ⚠️  ~/.vimrc - Exists but is not a symlink (manual intervention required)
# 
# Summary:
# - 1 file(s) fixed
# - 1 file(s) need manual intervention

dotman fix ~/.bashrc
# Output: ✓ ~/.bashrc is already correctly linked
```

### `dotman status`

**Purpose**: Show status of all managed files

**Output**:
- List of managed files with status indicators
- Git repository status (if uncommitted changes)

**Example**:
```bash
dotman status
# Output:
# Dotman is managing 3 file(s):
# 
# ✓ ~/.bashrc (file) - OK
# ✓ ~/.config/sway (directory) - OK
# ✗ ~/.vimrc (file) - Missing
```

## Path Handling Rules

1. **HOME Only**: All paths must resolve to locations within `$HOME`
2. **Absolute Resolution**: All paths are converted to absolute paths
3. **Tilde Expansion**: `~` and `~/` are expanded to `$HOME`
4. **Relative to HOME**: Paths are stored in index relative to `$HOME`
5. **Symlink Validation**: Original locations must be symlinks pointing to repository

## Repository Structure

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

## Error Handling

All commands should provide clear, actionable error messages:
- Specify what went wrong
- Suggest how to fix it
- Maintain system safety (no partial states)
