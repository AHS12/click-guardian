<p align="center">
  <img src="assets/icon-modern-shield.svg" alt="Click Guardian Logo" width="128" height="128">
</p>

<h1 align="center">Click Guardian</h1>

<p align="center">
  <strong>An application that prevents accidental double-clicks by introducing a customizable delay between mouse clicks - currently available on Windows.</strong>
</p>

## Features

- 🎯 **Strict Double-Click Blocking**: Ensures no double-clicks are allowed under any circumstances
- ⚙️ **Customizable Delay**: Set delay from 5ms to 500ms (default: 50ms)
- � **Adaptive Protection**: Automatically increases delay when faulty mouse hardware is detected (never decreases below user setting)
- �📊 **Real-time Logging**: Detailed logs for allowed and blocked clicks, including reasons and timestamps
- 🖥️ **Modern GUI**: Clean and intuitive Fyne-based interface
- 🚀 **Lightweight**: Minimal resource usage
- 🛡️ **Safe**: Only monitors clicks, doesn't interfere with other mouse operations

## How It Works

The application installs a low-level mouse hook that monitors left and right mouse button clicks. When a click is detected:

1. **First Click**: Always allowed and logged
2. **Subsequent Clicks**: Strictly blocked if they occur within the specified delay period for that specific button
3. **Independent Timers**: Left and right mouse buttons have separate timers - switching between buttons doesn't reset the protection
4. **Adaptive Protection**: Automatically detects faulty mouse hardware patterns (like low-pressure clicks) and increases the protection delay accordingly - never reduces below your selected setting

The adaptive system ensures maximum protection against problematic mice while maintaining your chosen baseline delay for normal operation.

## Quick Start

1. **Set Delay**: Enter your desired delay in milliseconds (5-500ms)
2. **Start Protection**: Click "Start Protection" to begin monitoring clicks
3. **Monitor Activity**: Watch the real-time log showing allowed/blocked clicks
4. **Stop Protection**: Click "Stop Protection" when finished

_Tip: Start with the default 50ms delay - it works well for most users._

## Installation

### Download Release

_Coming soon - pre-built executables will be available from the releases page_

### Build from Source

For detailed build instructions, see [**Development Guide**](docs/DEVELOPMENT.md)

## Configuration

- **Default Delay**: 50ms (good for most users)
- **Recommended Range**: 30-100ms for most applications
- **Gaming**: 10-30ms for fast-paced games
- **Accessibility**: 100-500ms for users with motor difficulties

## Documentation

- 📖 [**Development Guide**](docs/DEVELOPMENT.md) - Building, project structure, and development setup
- ⚙️ [**VSCode Setup**](docs/VSCODE_SETUP.md) - IDE configuration and troubleshooting
- 🔧 [**Build Instructions**](docs/BUILD.md) - Detailed build documentation

## Cross-Platform Support

- **Windows**: ✅ Fully supported (current)
- **Linux**: 🚧 Planned (X11/Wayland support)
- **macOS**: 🚧 Planned

## Contributing

Feel free to submit issues and pull requests to improve this application. See the [Development Guide](docs/DEVELOPMENT.md) for getting started.

## License

This project is open source and available under the MIT License.
