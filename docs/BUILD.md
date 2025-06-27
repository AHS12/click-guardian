# Click Guardian - Cross-Platform Build Instructions

## Quick Start

### Windows

```cmd
scripts\build.bat
```

### Linux/macOS

```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

## Manual Build Commands

### Single Platform Build

```bash
# Windows GUI version
go build -ldflags "-s -w -H=windowsgui" -o dist/click-guardian-gui.exe ./cmd/click-guardian

# Windows Console version
go build -ldflags "-s -w" -o dist/click-guardian.exe ./cmd/click-guardian

# Linux/macOS
go build -ldflags "-s -w" -o dist/click-guardian ./cmd/click-guardian
```

### Cross-Platform Build

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-windows-amd64.exe ./cmd/click-guardian

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-linux-amd64 ./cmd/click-guardian

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/click-guardian-darwin-amd64 ./cmd/click-guardian
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/click-guardian-darwin-arm64 ./cmd/click-guardian
```

## Development

### Run in Development Mode

```bash
# Windows
scripts\dev.bat

# Linux/macOS
go run ./cmd/click-guardian
```

### Clean Build Artifacts

```bash
# Windows
rmdir /s /q dist

# Linux/macOS
rm -rf dist
```

## Dependencies

### Prerequisites

- Go 1.24.1 or later
- C compiler (for CGO on Windows)

### Install Dependencies

```bash
go mod tidy
```
