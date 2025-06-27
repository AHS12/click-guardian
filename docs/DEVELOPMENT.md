# Development Guide

## Project Structure

```
go-double-click-fix/
├── cmd/
│   └── doubleclick-fix/          # Main application entry point
│       └── main.go
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   ├── gui/                      # GUI application logic
│   ├── hooks/                    # Platform-specific mouse hooks
│   └── logger/                   # Logging functionality
├── pkg/                          # Public packages (for future use)
│   └── platform/                 # Platform detection utilities
├── scripts/                      # Build and development scripts
├── dist/                         # Build outputs (git ignored)
├── docs/                         # Documentation
└── assets/                       # Static assets (icons, etc.)
```

## Building

### Windows Development

```bash
# Run in development mode
scripts\dev.bat

# Build for production
scripts\build.bat
```

### Cross-platform Build

```bash
# Make script executable (Linux/macOS)
chmod +x scripts/build.sh

# Build for all platforms
scripts/build.sh
```

## Architecture

### Core Components

1. **GUI Layer** (`internal/gui/`)

   - Fyne-based user interface
   - Application state management
   - User interaction handling

2. **Hook Layer** (`internal/hooks/`)

   - Platform-specific mouse hook implementations
   - Windows: Low-level mouse hook using Windows API
   - Other platforms: Placeholder implementation

3. **Configuration** (`internal/config/`)

   - Application settings
   - Validation logic
   - Default values

4. **Logging** (`internal/logger/`)
   - Real-time log display
   - Message formatting
   - Thread-safe logging

### Platform Support

Currently fully supported:

- Windows (x86, x64)

Planned support:

- Linux (X11/Wayland)
- macOS
- Web (via Fyne web compilation)

## Adding New Platforms

1. Create platform-specific hook implementation in `internal/hooks/`
2. Use build tags to conditionally compile: `//go:build yourplatform`
3. Implement the `MouseHook` interface
4. Update build scripts to include the new platform

## Dependencies

- **Fyne v2**: GUI framework
- **Windows API**: For mouse hooking on Windows (CGO)

## Development Tips

- Use `scripts/dev.bat` for quick testing
- Check `internal/hooks/hook_unsupported.go` for reference implementation
- All builds output to `dist/` directory
- Use the console version for debugging output
