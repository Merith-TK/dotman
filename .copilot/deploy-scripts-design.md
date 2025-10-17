# Deploy Scripts Feature Design

## Overview
Deploy scripts are custom scripts that handle deployment of dotfiles to dynamic or complex locations that can't be handled by simple symlinks. This is particularly useful for applications that use generated directory names or have complex configuration structures.

## Use Cases

### Firefox Profile Management
- **Problem**: Firefox profiles have random names like `4yi8ybqt.default-release`
- **Solution**: Deploy script that finds the active profile and deploys `userChrome.css`
- **Path**: `~/.mozilla/firefox/*/chrome/userChrome.css`

### Browser Extensions
- **Problem**: Chrome/Edge extensions have random IDs
- **Solution**: Deploy script that finds extension directories and deploys configs
- **Path**: `~/.config/google-chrome/Default/Extensions/*/`

### VSCode Extensions
- **Problem**: VSCode extensions with version-specific paths
- **Solution**: Deploy script that handles extension configuration deployment
- **Path**: `~/.vscode/extensions/*/`

## Design Specification

### Deploy Script Format
Deploy scripts would be stored in `.dotman/scripts/` directory as executable files or configuration definitions.

#### Option 1: Shell Scripts
```bash
#!/bin/bash
# .dotman/scripts/firefox-userchrome.sh

# Find active Firefox profile
PROFILE_DIR=$(find ~/.mozilla/firefox -maxdepth 1 -name "*.default*" -type d | head -1)

if [ -z "$PROFILE_DIR" ]; then
    echo "No Firefox profile found"
    exit 1
fi

# Create chrome directory if it doesn't exist
mkdir -p "$PROFILE_DIR/chrome"

# Deploy userChrome.css
ln -sf "$DOTMAN_REPO/firefox/userChrome.css" "$PROFILE_DIR/chrome/userChrome.css"

echo "Deployed userChrome.css to $PROFILE_DIR/chrome/"
```

#### Option 2: YAML Configuration
```yaml
# .dotman/scripts/firefox.yaml
name: "Firefox userChrome.css"
description: "Deploy userChrome.css to active Firefox profile"

target_discovery:
  method: "glob"
  pattern: "~/.mozilla/firefox/*.default*/chrome/"
  create_path: true

deployments:
  - source: "firefox/userChrome.css"
    target: "userChrome.css"
    method: "symlink"

validation:
  - check_file_exists: "~/.mozilla/firefox/"
  - warn_if_missing: "Firefox profile directory"
```

#### Option 3: Go-based Plugins
```go
// .dotman/scripts/firefox.go
package main

import (
    "path/filepath"
    "os"
    "fmt"
)

func Deploy(dotmanRepo string, homeDir string) error {
    // Find Firefox profile
    profiles, err := filepath.Glob(filepath.Join(homeDir, ".mozilla/firefox/*.default*"))
    if err != nil || len(profiles) == 0 {
        return fmt.Errorf("no Firefox profile found")
    }
    
    profile := profiles[0]
    chromeDir := filepath.Join(profile, "chrome")
    
    // Create chrome directory
    os.MkdirAll(chromeDir, 0755)
    
    // Create symlink
    source := filepath.Join(dotmanRepo, "firefox/userChrome.css")
    target := filepath.Join(chromeDir, "userChrome.css")
    
    return os.Symlink(source, target)
}
```

## Implementation Plan

### Phase 1: Basic Script Support
1. **Script Directory**: Create `.dotman/scripts/` directory
2. **Script Registry**: Track scripts in `index.json` or separate `scripts.json`
3. **Script Execution**: `dotman deploy --scripts` command
4. **Environment Variables**: Provide `$DOTMAN_REPO`, `$HOME` to scripts

### Phase 2: Pattern-Based Discovery
1. **Glob Patterns**: Support for `~/.mozilla/firefox/*/chrome/`
2. **Multiple Matches**: Handle multiple profile scenarios
3. **Path Creation**: Automatically create intermediate directories
4. **Conflict Resolution**: Handle existing files/symlinks

### Phase 3: Validation & Safety
1. **Pre-flight Checks**: Validate target applications exist
2. **Dry Run Mode**: Show what would be deployed without doing it
3. **Rollback Support**: Ability to undo script deployments
4. **Error Handling**: Graceful failure and reporting

## Command Interface

### New Commands
```bash
# Run all deploy scripts
dotman deploy --scripts

# Run specific script
dotman deploy --script firefox

# List available scripts
dotman scripts list

# Validate scripts without running
dotman scripts validate

# Add a new script
dotman scripts add firefox.sh

# Remove a script
dotman scripts remove firefox
```

### Enhanced Commands
```bash
# Deploy both symlinks and scripts
dotman deploy --all

# Status including script deployments
dotman status --scripts
```

## Storage Structure

```
~/.dotman/
├── scripts/
│   ├── firefox.sh           # Firefox userChrome deployment
│   ├── vscode-settings.yaml # VSCode settings deployment
│   └── chrome-extensions.go # Chrome extension deployment
├── firefox/
│   └── userChrome.css       # Source file
├── vscode/
│   └── settings.json        # VSCode settings
└── index.json               # Main index (enhanced)
```

## Enhanced Index Format

```json
{
  "version": "1.0",
  "managed_files": [...],
  "deploy_scripts": [
    {
      "name": "firefox",
      "script_path": "scripts/firefox.sh",
      "description": "Deploy Firefox userChrome.css",
      "targets": [
        "~/.mozilla/firefox/*/chrome/userChrome.css"
      ],
      "dependencies": ["firefox/userChrome.css"],
      "enabled": true,
      "last_run": "2025-08-01T10:00:00Z"
    }
  ]
}
```

## Benefits

1. **Handles Complex Paths**: Solves the Firefox profile problem elegantly
2. **Extensible**: Can handle any application with dynamic paths
3. **Version Controlled**: Scripts are part of the dotfiles repository
4. **Cross-Platform**: Can have platform-specific scripts
5. **Maintainable**: Clear separation between simple symlinks and complex deployments

## Security Considerations

1. **Script Validation**: Verify scripts before execution
2. **Sandboxing**: Limit script capabilities to dotfiles deployment
3. **User Confirmation**: Prompt before running unknown scripts
4. **Audit Trail**: Log script executions and changes

This feature would make dotman much more powerful for handling real-world dotfiles scenarios while maintaining the simplicity of the core symlink-based approach.
