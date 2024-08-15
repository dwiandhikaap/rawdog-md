#!/usr/bin/env bash
set -euo pipefail

# Determine the platform
platform=$(uname -ms)

# Function to display error messages
error() {
    echo -e "\033[0;31merror\033[0m:" "$@" >&2
    exit 1
}

# Function to display information messages
info() {
    echo -e "\033[0;2m$@\033[0m"
}

# Check for required tools
command -v curl >/dev/null || error 'curl is required to install rawdog-md'
command -v tar >/dev/null || error 'tar is required to install rawdog-md'

# Set installation directory
INSTALL_DIR="$HOME/.local/share/rawdog-md"  # Change as needed
mkdir -p "$INSTALL_DIR" || error "Failed to create install directory \"$INSTALL_DIR\""

# Determine the target based on OS
version="v0.1.7"

case $platform in
    'Darwin x86_64')
        target="darwin-amd64"
        asset_name="rawd-$version-darwin-amd64.tar.gz"
        ;;
    'Darwin arm64')
        target="darwin-arm64"
        asset_name="rawd-$version-darwin-arm64.tar.gz"
        ;;
    'Linux x86_64')
        target="linux-amd64"
        asset_name="rawd-$version-linux-amd64.tar.gz"
        ;;
    'Linux aarch64')
        target="linux-arm64"
        asset_name="rawd-$version-linux-arm64.tar.gz"
        ;;
    *)
        error "Unsupported operating system: $platform"
        ;;
esac

# Set the GitHub repository
GITHUB_USER="dwiandhikaap"
GITHUB_REPO="rawdog-md"
DOWNLOAD_URL="https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/latest/download/$asset_name"

# Download the application to /tmp
info "Downloading the latest release..."
curl -L -o "/tmp/rawdog-md.tar.gz" "$DOWNLOAD_URL" || error "Failed to download from \"$DOWNLOAD_URL\""

# Extract the application
info "Extracting the files..."
tar -xzf "/tmp/rawdog-md.tar.gz" -C "$INSTALL_DIR" || error "Failed to extract the tar.gz file"

# Set executable permissions
chmod +x "$INSTALL_DIR/rawd" || error "Failed to set permissions on $INSTALL_DIR/rawd"

# Clean up the downloaded tar.gz file
rm "/tmp/rawdog-md.tar.gz"

# Add to PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc"
    info "Added $INSTALL_DIR to PATH. Please restart your terminal or run 'source ~/.bashrc' or 'source ~/.zshrc' for changes to take effect."
else
    info "$INSTALL_DIR is already in PATH."
fi

info "Installation complete! You can now use the 'rawd' command in your terminal."