# Click Guardian - Build Instructions

## Quick Start

### Windows

```cmd
# Quick development build and run
scripts\dev.bat

# Production build
scripts\build.bat
```

### Linux/macOS

```bash
# Make executable (first time only)
chmod +x scripts/build.sh

# Build for all platforms
./scripts/build.sh
```

## Build Scripts

This project includes several build and development scripts in the `scripts/` directory:

### `build.bat` (Windows)

Main Windows build script that creates both GUI and console versions.

- Creates executables in `dist/` directory
- Builds GUI version (no console window) and console version (for debugging)
- Handles build failures gracefully
- Can be run from anywhere in the project

### `build.sh` (Linux/macOS)

Cross-platform build script for multiple platforms.

- Requires bash shell
- Builds for Windows, Linux, and macOS
- Creates both GUI and console versions for Windows
- Outputs with platform-specific naming convention

### `dev.bat` (Windows)

Development script for quick testing.

- Runs the application directly without building
- Useful for rapid development iteration
- Automatically navigates to project root

### `troubleshoot.bat` (Windows)

Troubleshooting script for Go/CGO setup issues.

- Tests CGO compilation
- Verifies project builds
- Cleans and refreshes dependencies
- Helpful for VSCode setup problems

### Script Usage

All scripts automatically navigate to the project root directory, so they can be run from anywhere within the project:

```cmd
# From any directory in the project
scripts\dev.bat
scripts\build.bat
scripts\troubleshoot.bat
```

## Manual Build Commands

### Single Platform Build

```bash
# Windows GUI version (recommended for end users)
go build -ldflags "-s -w -H=windowsgui" -o dist/click-guardian-gui.exe ./cmd/click-guardian

# Windows Console version (for debugging)
go build -ldflags "-s -w" -o dist/click-guardian.exe ./cmd/click-guardian

# Linux/macOS
go build -ldflags "-s -w" -o dist/click-guardian ./cmd/click-guardian
```

### Cross-Platform Build

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-windows-amd64.exe ./cmd/click-guardian

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-linux-amd64 ./cmd/click-guardian

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-darwin-amd64 ./cmd/click-guardian

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/click-guardian-darwin-arm64 ./cmd/click-guardian
```

## Build Output

All builds output to the `dist/` directory (git-ignored).

### Naming Convention

- `click-guardian-gui.exe` - Windows GUI version (no console, recommended)
- `click-guardian.exe` - Windows console version (shows debug output)
- `click-guardian-{os}-{arch}` - Cross-platform builds

## Development

### Prerequisites

- **Go 1.24.1 or later**
- **C compiler** (for CGO on Windows - typically MinGW or Visual Studio Build Tools)
- **Git** (for cloning the repository)

### Setup

```bash
# Clone and setup
git clone [repository-url]
cd click-guardian
go mod tidy
```

### Development Workflow

```bash
# Quick development testing (Windows)
scripts\dev.bat

# Or run manually (any platform)
go run ./cmd/click-guardian
```

### Clean Build Artifacts

```bash
# Windows
rmdir /s /q dist

# Linux/macOS
rm -rf dist
```

## Troubleshooting

### Build Issues

**CGO compilation errors:**

```bash
# Run troubleshooting script (Windows)
scripts\troubleshoot.bat
```

**Module issues:**

```bash
go clean -modcache
go mod download
go mod tidy
```

**VSCode setup issues:**
See [VSCode Setup Guide](VSCODE_SETUP.md) for workspace configuration.

### Platform-Specific Notes

**Windows:**

- CGO is required for Windows API integration
- GUI version (`-gui.exe`) has no console window
- Console version (`.exe`) shows debug output

**Linux/macOS:**

- Currently builds but mouse hook functionality is not implemented
- Future releases will include X11/Wayland support
