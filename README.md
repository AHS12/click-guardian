<p align="center">
  <img src="assets/icon-modern-shield.svg" alt="Click Guardian Logo" width="128" height="128">
</p>

<h1 align="center">Click Guardian</h1>

<p align="center">
  <strong>An application that prevents accidental double-clicks by introducing a customizable delay between mouse clicks - currently available on Windows.</strong>
</p>

## Features

- 🎯 **Smart Click Protection**: Prevents double-clicks within a specified time window
- ⚙️ **Customizable Delay**: Set delay from 1ms to 5000ms (default: 50ms)
- 📊 **Real-time Logging**: See exactly which clicks are allowed or blocked
- 🖥️ **Modern GUI**: Clean and intuitive Fyne-based interface
- 🚀 **Lightweight**: Minimal resource usage
- 🛡️ **Safe**: Only monitors clicks, doesn't interfere with other mouse operations

## How It Works

The application installs a low-level mouse hook that monitors left and right mouse button clicks. When a click is detected:

1. **First Click**: Always allowed and logged
2. **Subsequent Clicks**: Only allowed if they occur after the specified delay period for that specific button

## Usage

1. **Set Delay**: Enter your desired delay in milliseconds (1-5000ms)
2. **Start Protection**: Click "Start Protection" to begin monitoring clicks
3. **Monitor Activity**: Watch the real-time log showing allowed/blocked clicks
4. **Stop Protection**: Click "Stop Protection" when finished

_Tip: Start with the default 50ms delay - it works well for most users._

## Installation

### Option 1: Download Release

_Coming soon - pre-built executables will be available from the releases page_

### Option 2: Build from Source

**Prerequisites:**

- Windows operating system
- Go 1.24.1 or later

**Steps:**

```bash
# Clone the repository then
cd click-guardian
go mod tidy

# Build using the build script (Windows)
scripts\build.bat

# Or build manually
go build -o dist\click-guardian.exe .\cmd\click-guardian
```

### Running the Application

After building, run the executable:

```bash
dist\click-guardian.exe
```

**For Development:**

```bash
scripts\dev.bat
```

## Project Structure

```
click-guardian/
├── cmd/click-guardian/          # Main application entry point
├── internal/                     # Private application packages
│   ├── config/                   # Configuration management
│   ├── gui/                      # GUI application logic
│   ├── hooks/                    # Platform-specific mouse hooks
│   └── logger/                   # Logging functionality
├── scripts/                      # Build and development scripts
├── dist/                         # Build outputs (executables)
├── docs/                         # Documentation
└── assets/                       # Static assets
```

## Configuration

- **Default Delay**: 50ms (good for most users)
- **Recommended Range**: 30-100ms for most applications
- **Gaming**: 10-30ms for fast-paced games
- **Accessibility**: 100-500ms for users with motor difficulties

## Technical Details

- **Platform**: Currently Windows only (uses Windows API)
- **Framework**: Fyne v2 for GUI
- **Hook Type**: Low-level mouse hook (WH_MOUSE_LL)
- **Permissions**: May require administrator privileges on some systems
- **Architecture**: Modular design prepared for cross-platform expansion

### Cross-Platform Roadmap

The project is structured to support multiple platforms in the future:

- **Windows**: ✅ Fully supported (current)
- **Linux**: 🚧 Planned (X11/Wayland support)
- **macOS**: 🚧 Planned

### VSCode Build Constraint Errors

If you see errors like "build constraints exclude all Go files" in VSCode, this is a known issue with the Go language server and cross-platform dependencies. The errors are cosmetic - builds work fine.

**Quick fix:**

1. Open the workspace file: `click-guardian.code-workspace`
2. Or reload VSCode: `Ctrl+Shift+P` → "Developer: Reload Window"

See `docs/VSCODE_SETUP.md` for detailed solutions.

## Contributing

Feel free to submit issues and pull requests to improve this application.

## License

This project is open source and available under the MIT License.
