@echo off
setlocal

echo Building Click Guardian...

REM Navigate to project root
cd /d "%~dp0.."
echo Building from: %CD%
echo.

REM Verify go.mod exists
if not exist "go.mod" (
    echo ❌ go.mod not found in current directory!
    echo Current directory contents:
    dir /b
    exit /b 1
)

REM Create dist directory if it doesn't exist
if not exist "dist" mkdir dist

REM Clean previous builds
del /Q dist\*.exe 2>nul

REM Verify the source path exists
if not exist "cmd\click-guardian" (
    echo ❌ Source directory cmd\click-guardian not found!
    echo Available directories:
    dir /b /ad
    exit /b 1
)

REM Build GUI version (no console window)
echo Building GUI version...
go build -ldflags "-s -w -H=windowsgui" -o dist\click-guardian-gui.exe .\cmd\click-guardian

if %ERRORLEVEL% EQU 0 (
    echo ✅ GUI build successful! Created dist\click-guardian-gui.exe
) else (
    echo ❌ GUI build failed with error code %ERRORLEVEL%!
    goto console_build
)

REM Build console version (with console window for debugging)
:console_build
echo Building console version...
go build -ldflags "-s -w" -o dist\click-guardian.exe .\cmd\click-guardian

if %ERRORLEVEL% EQU 0 (
    echo ✅ Console build successful! Created dist\click-guardian.exe
) else (
    echo ❌ Console build failed with error code %ERRORLEVEL%!
    exit /b 1
)

echo.
echo Build complete! Executables are in the dist folder.
echo - dist\click-guardian-gui.exe (recommended for normal use)
echo - dist\click-guardian.exe (for debugging/console output)

REM Don't pause in CI environment
if not defined GITHUB_ACTIONS pause