@echo off
setlocal

echo Building Double-Click Fix...

REM Navigate to project root
cd /d "%~dp0.."
echo Building from: %CD%
echo.

REM Create dist directory if it doesn't exist
if not exist "dist" mkdir dist

REM Clean previous builds
del /Q dist\*.exe 2>nul

REM Build GUI version (no console window)
echo Building GUI version...
go build -ldflags "-s -w -H=windowsgui" -o dist\go-double-click-fix-gui.exe .\cmd\doubleclick-fix

if %ERRORLEVEL% EQU 0 (
    echo ✅ GUI build successful! Created dist\go-double-click-fix-gui.exe
) else (
    echo ❌ GUI build failed!
    goto console_build
)

REM Build console version (with console window for debugging)
:console_build
echo Building console version...
go build -ldflags "-s -w" -o dist\go-double-click-fix.exe .\cmd\doubleclick-fix

if %ERRORLEVEL% EQU 0 (
    echo ✅ Console build successful! Created dist\go-double-click-fix.exe
) else (
    echo ❌ Console build failed!
    goto end
)

echo.
echo Build complete! Executables are in the dist folder.
echo - dist\go-double-click-fix-gui.exe (recommended for normal use)
echo - dist\go-double-click-fix.exe (for debugging/console output)

:end
pause
