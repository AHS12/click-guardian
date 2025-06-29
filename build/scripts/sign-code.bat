@echo off
setlocal

echo =====================================
echo  Code Signing Script
echo =====================================
echo.

if "%1"=="" (
    echo Usage: sign-code.bat ^<executable^> [certificate_file] [password]
    echo.
    echo Examples:
    echo   sign-code.bat dist\click-guardian-gui.exe
    echo   sign-code.bat dist\click-guardian-gui.exe mycert.pfx mypassword
    echo.
    exit /b 1
)

set EXECUTABLE=%1
set CERT_FILE=%2
set CERT_PASSWORD=%3

REM Default values if not provided
if "%CERT_FILE%"=="" set CERT_FILE=signing-certificate.pfx
if "%CERT_PASSWORD%"=="" set /P CERT_PASSWORD=Enter certificate password: 

echo File to sign: %EXECUTABLE%
echo Certificate: %CERT_FILE%
echo.

REM Check if file exists
if not exist "%EXECUTABLE%" (
    echo ❌ File not found: %EXECUTABLE%
    exit /b 1
)

REM Check if certificate exists
if not exist "%CERT_FILE%" (
    echo ❌ Certificate not found: %CERT_FILE%
    echo.
    echo To obtain a code signing certificate:
    echo 1. Purchase from a trusted CA (DigiCert, Sectigo, etc.)
    echo 2. Or create a self-signed certificate for testing:
    echo    makecert -sv mykey.pvk -n "CN=YourName" mycert.cer
    echo    pvk2pfx -pvk mykey.pvk -spc mycert.cer -pfx mycert.pfx
    echo.
    exit /b 1
)

REM Check if signtool is available
signtool sign /? >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ signtool not found!
    echo Install Windows SDK from:
    echo https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/
    exit /b 1
)

echo Signing executable...
signtool sign ^
    /f "%CERT_FILE%" ^
    /p "%CERT_PASSWORD%" ^
    /t http://timestamp.digicert.com ^
    /d "Click Guardian" ^
    /du "https://github.com/your-repo" ^
    "%EXECUTABLE%"

if %ERRORLEVEL% EQU 0 (
    echo ✅ Successfully signed: %EXECUTABLE%
    
    REM Verify the signature
    echo Verifying signature...
    signtool verify /pa "%EXECUTABLE%"
    
    if %ERRORLEVEL% EQU 0 (
        echo ✅ Signature verified successfully!
    ) else (
        echo ⚠️  Signature verification failed
    )
) else (
    echo ❌ Signing failed!
    echo.
    echo Common issues:
    echo - Wrong certificate password
    echo - Certificate expired
    echo - Certificate not trusted for code signing
    echo - Network issues with timestamp server
)

pause
