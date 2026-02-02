#!/bin/bash

set -e
set -o pipefail

on_error() { echo "Script failed. Error on line: $1. Terminating..."; }
trap 'on_error $LINENO' ERR

echo
echo "========================================================================="
echo "                            StepKeys Installer                           "
echo "========================================================================="
echo

# Detect platform
OS="$(uname -s)"
ARCH="$(uname -m)"

if [[ "$OS" == "Darwin" ]]; then
    PLATFORM="macos-arm64"
elif [[ "$OS" == "Linux" ]]; then
    PLATFORM="linux-amd64"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Fail on macOS x86_64 / amd64 and Linux arm64
if [[ "$OS" == "Darwin" && "$ARCH" != "arm64" ]]; then
    echo "Unsupported architecture on macOS: $ARCH. Only an arm64 build is available. You can build StepKeys from source."
    exit 1
elif [[ "$OS" == "Linux" && "$ARCH" != "x86_64" ]]; then
    echo "Unsupported architecture on Linux: $ARCH. Only an amd64 build is available. You can build StepKeys from source."
    exit 1
fi

# Handle update argument
UPDATE=false
if [[ "$1" == "update" ]]; then
    UPDATE=true

    echo "Updating StepKeys..."

    echo -e "\nRequesting StepKeys shutdown via API..."
    curl -s -X POST http://localhost:18000/api/quit || true

    echo -e "Waiting for StepKeys to stop...\n"
    sleep 2
fi

# Latest release URL
REPO="https://github.com/BrNi05/StepKeys/releases/latest"

# Installation directory
TMP_DIR="$(mktemp -d)"
INSTALL_DIR="$HOME/.stepkeys"
mkdir -p "$INSTALL_DIR"

echo "Detected platform: $PLATFORM"
echo "Downloading latest StepKeys release..."

# Get latest version by following GitHub redirect
VERSION=$(curl -fsSL -o /dev/null -w '%{url_effective}' "$REPO" | sed -E 's|.*/tag/v?||')
BINARY_NAME="stepkeys-${PLATFORM}-${VERSION}"
DOWNLOAD_URL="https://github.com/BrNi05/StepKeys/releases/download/v${VERSION}/${BINARY_NAME}"

# Download the binary
echo -e "\nDownloading $DOWNLOAD_URL..."
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME"
chmod +x "$TMP_DIR/$BINARY_NAME"

# macOS quarantine removal
if [[ "$OS" == "Darwin" ]]; then
    xattr -d com.apple.quarantine "$TMP_DIR/$BINARY_NAME" || true
fi

# Move binary to install dir
mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/stepkeys"
echo -e "\nStepKeys installed to $INSTALL_DIR"

rm -rf "$TMP_DIR"

# dotenv setup
if [[ "$UPDATE" == false ]]; then
    echo -e "\nSetting up configuration (.env)..."

    echo -e "\nDetecting available serial devices..."
    if [[ "$OS" == "Linux" ]]; then
        DEVICES=$(ls /dev/ttyACM* 2>/dev/null || echo "none")
    elif [[ "$OS" == "Darwin" ]]; then
        DEVICES=$(ls /dev/cu.* 2>/dev/null || echo "none")
    fi

    echo "Available serial devices: $DEVICES"

    read -rp "Enter the serial port to use (leave empty to skip): " SERIAL_PORT
    read -rp "Enter baud rate [default 115200]: " BAUD_RATE

    # Write .env
    ENV_FILE="$INSTALL_DIR/.env"
    echo "SERIAL_PORT=$SERIAL_PORT" > "$ENV_FILE"
    echo "BAUD_RATE=$BAUD_RATE" >> "$ENV_FILE"
    echo "VERSION=$VERSION" >> "$ENV_FILE" # auto-generated, used for update checks
else
    echo -e "\nUpdate mode: keeping existing .env configuration"
fi

# Linux permission warning
if [ "$OS" == "Linux" ] && [ "$UPDATE" == false ]; then
    echo -e "\nNote: On Linux, ensure your user has permission to access the serial port."
    echo "You may need to add your user to the 'dialout' or 'uucp' group. See the StepKeys docs for details."
fi

# Start StepKeys
echo -e "\nStarting StepKeys in the background..."
nohup "$INSTALL_DIR/stepkeys" > "$INSTALL_DIR/stepkeys.log" 2>&1 &

echo -e "\nStepKeys started successfully!"
echo -e "Logs: $INSTALL_DIR/stepkeys.log\n"
