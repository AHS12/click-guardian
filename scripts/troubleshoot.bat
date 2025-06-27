@echo off
echo Troubleshooting Go/CGO setup for VSCode...
echo.

REM Save current directory and navigate to project root
set SCRIPT_DIR=%~dp0
cd /d "%SCRIPT_DIR%.."
echo Current directory: %CD%
echo.

echo Checking Go version:
go version
echo.

echo Checking CGO status:
go env CGO_ENABLED
echo.

echo Checking GOOS/GOARCH:
go env GOOS
go env GOARCH
echo.

echo Testing CGO compilation:
echo package main > test_cgo.go
echo. >> test_cgo.go
echo /* >> test_cgo.go
echo #include ^<stdio.h^> >> test_cgo.go
echo */ >> test_cgo.go
echo import "C" >> test_cgo.go
echo. >> test_cgo.go
echo func main() { >> test_cgo.go
echo     println("CGO test") >> test_cgo.go
echo } >> test_cgo.go

go build -o test_cgo.exe test_cgo.go
if %ERRORLEVEL% EQU 0 (
    echo ✅ CGO compilation works!
    del test_cgo.exe
) else (
    echo ❌ CGO compilation failed!
    echo.
    echo Possible solutions:
    echo 1. Install TDM-GCC or MinGW-w64
    echo 2. Ensure gcc is in your PATH
    echo 3. Install Microsoft C++ Build Tools
)

del test_cgo.go 2>nul

echo.
echo Cleaning Go module cache:
go clean -modcache

echo.
echo Downloading dependencies:
go mod download

echo.
echo Testing project build:
go build -o dist\test-troubleshoot.exe .\cmd\doubleclick-fix
if %ERRORLEVEL% EQU 0 (
    echo ✅ Project builds successfully!
    del dist\test-troubleshoot.exe 2>nul
) else (
    echo ❌ Project build failed!
)

echo.
echo Troubleshooting complete!
REM Return to original directory
cd /d "%SCRIPT_DIR%"
pause
