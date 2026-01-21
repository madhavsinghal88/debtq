#!/bin/bash

set -e

BINARY_NAME="debtq"
INSTALL_DIR="$HOME/.local/bin"

echo "==================================="
echo "  DebtQ - Personal Money Tracker"
echo "==================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go from https://go.dev/dl/"
    exit 1
fi

echo "Building $BINARY_NAME..."
go build -ldflags="-s -w" -o "$BINARY_NAME" ./cmd/main.go

if [ ! -f "$BINARY_NAME" ]; then
    echo "Error: Build failed."
    exit 1
fi

echo "Build successful!"
echo ""

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Move binary to install directory
echo "Installing to $INSTALL_DIR..."
mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "Installation complete!"
echo ""

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "WARNING: $INSTALL_DIR is not in your PATH."
    echo ""
    echo "Add the following line to your ~/.bashrc or ~/.zshrc:"
    echo ""
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "Then run: source ~/.bashrc (or source ~/.zshrc)"
    echo ""
fi

echo "Run 'debtq' to start the application."
echo ""

# Create default config directory
CONFIG_DIR="$HOME/.config/debtq"
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Creating config directory at $CONFIG_DIR..."
    mkdir -p "$CONFIG_DIR"
fi

# Create default obsidian directory
OBSIDIAN_DIR="$HOME/Documents/obsidian-notes/debtq"
if [ ! -d "$OBSIDIAN_DIR" ]; then
    echo "Creating Obsidian vault directory at $OBSIDIAN_DIR..."
    mkdir -p "$OBSIDIAN_DIR"
fi

echo ""
echo "Setup complete! Enjoy using DebtQ!"
