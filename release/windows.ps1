Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Show-Error {
    param($message)
    Write-Host "ERROR: $message" -ForegroundColor Red
    if (Test-Path $TMP_FILE) { Remove-Item $TMP_FILE -Force }
    Exit 1
}

Write-Host ""
Write-Host "========================================================================="
Write-Host "                            StepKeys Installer                           "
Write-Host "========================================================================="
Write-Host ""

# Detect platform
$ARCH = (Get-CimInstance Win32_OperatingSystem).OSArchitecture
if ($ARCH -notmatch "64") {
    Show-Error "Unsupported architecture: $ARCH. Only x64 (64-bit) build is available. You can build StepKeys from source."
}

# Handle update argument
$UPDATE = $false
if ($args.Length -ge 1 -and $args[0] -eq "update") {
    $UPDATE = $true
    Write-Host "Updating StepKeys..."
    
    Write-Host "`nRequesting StepKeys shutdown via API..."
    try {
        Invoke-RestMethod -Uri http://localhost:18000/api/quit -Method Post -ErrorAction SilentlyContinue
    } catch {}
    
    Write-Host "Waiting for StepKeys to stop...`n"
    Start-Sleep -Seconds 2
}

# Latest release URL
$REPO = "https://github.com/BrNi05/StepKeys/releases/latest"

# Installation directory
$TMP_FILE = [System.IO.Path]::GetTempFileName() + ".exe"
$INSTALL_DIR = "$HOME\.stepkeys"
if (-not (Test-Path $INSTALL_DIR)) { New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null }

Write-Host "Detected platform: Windows-x64"
Write-Host "Downloading latest StepKeys release..."

# Get latest version by following GitHub redirect
try {
    $response = Invoke-WebRequest -Uri $REPO -MaximumRedirection 0 -ErrorAction SilentlyContinue
    $finalUrl = $response.Headers["Location"]
    $VERSION = ($finalUrl -split "/")[-1] -replace "^v", ""
} catch {
    Show-Error "Failed to detect latest version: $_"
}

$BINARY_NAME = "stepkeys-windows-amd64-$VERSION.exe"
$DOWNLOAD_URL = "https://github.com/BrNi05/StepKeys/releases/download/v$VERSION/$BINARY_NAME"

# Download the binary
Write-Host "`nDownloading $DOWNLOAD_URL..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TMP_FILE

# Move binary to install dir
Move-Item -Path $TMP_FILE -Destination "$INSTALL_DIR\stepkeys.exe" -Force
Write-Host "`nStepKeys installed to $INSTALL_DIR"

# dotenv setup
if (-not $UPDATE) {
    Write-Host "`nSetting up configuration (.env)..."

    Write-Host "`nDetecting available serial devices..."
    $ports = Get-CimInstance Win32_SerialPort | Select-Object -ExpandProperty DeviceID
    if (-not $ports) { $ports = "none found" }
    Write-Host "Available serial devices: $ports"

    Write-Host ""
    $SERIAL_PORT = Read-Host "Enter the serial port to use (leave empty to skip)"
    $BAUD_RATE = Read-Host "Enter baud rate [default 115200]"

    if ([string]::IsNullOrWhiteSpace($BAUD_RATE)) {
        $BAUD_RATE = "115200"
    }

    # Write .env
    # VERSION is auto-generated and is used for update checks
    $ENV_FILE = "$INSTALL_DIR\.env"
    @"
SERIAL_PORT=$SERIAL_PORT
BAUD_RATE=$BAUD_RATE
VERSION=$VERSION
"@ | Out-File -FilePath $ENV_FILE -Encoding UTF8
} else {
    Write-Host "`nUpdate mode: keeping existing .env configuration"

    $ENV_FILE = "$INSTALL_DIR\.env"
    if (Test-Path $ENV_FILE) {
        $content = Get-Content $ENV_FILE

        if ($content -match '^VERSION=') {
            $content = $content -replace '^VERSION=.*', "VERSION=$VERSION" # update existing VERSION
        } else {
            $content += "VERSION=$VERSION"
        }

        $content | Set-Content -Path $ENV_FILE -Encoding UTF8
    } else {
        "VERSION=$VERSION" | Set-Content -Path $ENV_FILE -Encoding UTF8
    }
}

# Start StepKeys
Write-Host "`nStarting StepKeys in the background..."
Start-Process -FilePath "$INSTALL_DIR\stepkeys.exe" -NoNewWindow

Write-Host "`nStepKeys started successfully!"
Write-Host "Logs: $INSTALL_DIR\stepkeys.log"
