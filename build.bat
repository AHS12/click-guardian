@echo off
echo Building Double-Click Fix (GUI version)...
go build -ldflags "-s -w -H=windowsgui" -o go-double-click-fix-gui.exe .
if %ERRORLEVEL% EQU 0 (
    echo Build successful! Created go-double-click-fix-gui.exe
) else (
    echo Build failed, creating console version instead...
    go build -o go-double-click-fix.exe .
)
pause
