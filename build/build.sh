#!/bin/bash

# Error management
set -e
set -o pipefail

on_error() { echo "Script failed. Error on line: $1. Terminating..."; }
trap 'on_error $LINENO' ERR

echo
echo "========================================================================="
echo "                         StepKeys Release Builder                        "
echo "========================================================================="

# Detect platform and arch
OS="$(go env GOOS)"
ARCH="$(go env GOARCH)"

echo -e "\nDetected platform:"
echo "  OS:   $OS"
echo "  ARCH: $ARCH"

echo
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
GUI_DIR="$ROOT_DIR/gui"
GUI_DIST_DIR="$GUI_DIR/dist"
SERVER_DIR="$ROOT_DIR/server"
SERVER_GUI_DIR="$SERVER_DIR/web/gui"
BIN_DIR="$ROOT_DIR/build/bin"

echo -e "\nRoot: $ROOT_DIR"
echo "GUI: $GUI_DIR"
echo "Server: $SERVER_DIR"
echo "Build Output: $BIN_DIR"

echo
read -rp "Release version (eg. 1.0.0): " VERSION

# Cleanup
echo -e "\nCleaning up..."

rm -rf "$BIN_DIR"
rm -rf "$GUI_DIST_DIR"
rm -rf "$SERVER_GUI_DIR"

mkdir -p "$BIN_DIR"
mkdir -p "$GUI_DIST_DIR"
mkdir -p "$SERVER_GUI_DIR"

# Build GUI
echo -e "\nBuilding GUI..."

cd "$GUI_DIR"
npm ci --silent
npm run lint
npm run build

# Copy built GUI to server
echo -e "\nCopying built GUI to server..."
cp -r "$GUI_DIST_DIR/"* "$SERVER_GUI_DIR/"

# Compiling
echo -e "\nCompiling native binary..."

cd "$SERVER_DIR"
OUTPUT_NAME="stepkeys-${OS}-${ARCH}-${VERSION}"
OUTPUT_PATH="$BIN_DIR/$OUTPUT_NAME"

GOOS="$OS" GOARCH="$ARCH" go build -o "$OUTPUT_PATH" .

echo -e "Release build completed successfully!"

# Start the built binary
echo
read -rp "Start StepKeys (y/N): " RUN_NOW

if [[ "$RUN_NOW" =~ ^[Yy]$ ]]; then
  echo -e "\nStarting StepKeys..."
  chmod +x "$OUTPUT_PATH"
  "$OUTPUT_PATH"
fi

echo