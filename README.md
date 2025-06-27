# Click Guardian

A Windows application that prevents accidental double-clicks by introducing a customizable delay between mouse clicks.

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
2. **Subsequent Clicks**: Only allowed if they occur after the specified delay period
3. **Different Buttons**: Switching between left/right buttons resets the timer

## Usage

1. **Set Delay**: Enter the desired delay in milliseconds (1-5000ms)
2. **Start Protection**: Click "Start Protection" to begin monitoring
3. **Monitor Activity**: Watch the log to see which clicks are blocked/allowed
4. **Stop Protection**: Click "Stop Protection" when done

## Installation

### Prerequisites

- Windows operating system
- Go 1.24.1 or later (for building from source)

### Building from Source

```bash
git clone <repository-url>
cd click-guardian
go mod tidy

# Build using the build script (Windows)
scripts\build.bat

# Or build manually
go build -o dist\click-guardian.exe .\cmd\click-guardian
```

### Running

Run the executable from the dist folder:

```bash
dist\click-guardian.exe
```

For development, you can use:

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
- **Web**: 🚧 Possible via Fyne web compilation

## Troubleshooting

### "Failed to install mouse hook"

- Try running as administrator
- Check if antivirus is blocking the application
- Ensure no other mouse hook applications are running

### High CPU Usage

- Reduce log verbosity by clearing logs frequently
- Use reasonable delay values (avoid very small delays like 1ms)

### VSCode Build Constraint Errors

If you see errors like "build constraints exclude all Go files" in VSCode, this is a known issue with the Go language server and cross-platform dependencies. The errors are cosmetic - builds work fine.

**Quick fix:**

1. Open the workspace file: `click-guardian.code-workspace`
2. Or reload VSCode: `Ctrl+Shift+P` → "Developer: Reload Window"

See `docs/VSCODE_SETUP.md` for detailed solutions.

### Application Not Responding

- Stop protection and restart the application
- Check Windows Event Viewer for system errors

## Contributing

Feel free to submit issues and pull requests to improve this application.

## License

This project is open source and available under the MIT License.
