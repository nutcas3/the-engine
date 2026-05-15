#!/bin/bash

set -e

# Sovereign Engine Installation Script
# This script installs the Sovereign Engine CLI and sets up the environment

VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.engine"
REPO_URL="https://github.com/yourusername/the-engine"
BINARY_BASE_URL="$REPO_URL/releases/download/v${VERSION}"

echo "🚀 Installing Sovereign Engine v${VERSION}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "📥 Detected platform: ${OS}-${ARCH}"

# Create config directory
echo "📁 Creating config directory..."
mkdir -p "$CONFIG_DIR"

# Download pre-built binaries
echo "� Downloading pre-built binaries..."
BINARY_NAME="engine-${OS}-${ARCH}"
DOWNLOAD_URL="${BINARY_BASE_URL}/${BINARY_NAME}.tar.gz"

echo "   Downloading from: $DOWNLOAD_URL"
if command -v curl &> /dev/null; then
    curl -L -o /tmp/engine.tar.gz "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O /tmp/engine.tar.gz "$DOWNLOAD_URL"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Extract binaries
echo "📦 Extracting binaries..."
tar -xzf /tmp/engine.tar.gz -C /tmp
rm /tmp/engine.tar.gz

# Install CLI binary
echo "📦 Installing CLI binary..."
sudo cp /tmp/engine "$INSTALL_DIR/engine"
sudo chmod +x "$INSTALL_DIR/engine"

# Install encryption tool
echo "📦 Installing encryption tool..."
sudo cp /tmp/engine-encrypt "$INSTALL_DIR/engine-encrypt"
sudo chmod +x "$INSTALL_DIR/engine-encrypt"

# Install web server binary
echo "📦 Installing web server binary..."
sudo cp /tmp/engine-web "$INSTALL_DIR/engine-web"
sudo chmod +x "$INSTALL_DIR/engine-web"

# Cleanup
rm -rf /tmp/engine /tmp/engine-encrypt /tmp/engine-web

# Initialize secure environment
echo "🔒 Setting up secure environment management..."
MASTER_KEY=$(openssl rand -base64 32)
echo "⚠️  IMPORTANT: Save this master key securely!"
echo "   MASTER_KEY=$MASTER_KEY"
echo "   Add this to your environment: export ENGINE_MASTER_KEY=$MASTER_KEY"

# Initialize configuration
echo "⚙️  Initializing configuration..."
cat > "$CONFIG_DIR/config.yaml" << EOF
version: ${VERSION}
compositions_dir: ./compositions
providers:
  - aws
  - azure
  - gcp
  - hetzner
  - ovh
  - digitalocean
secure_env_file: $CONFIG_DIR/secure.env
web_port: 8080
EOF

# Create secure environment file template
cat > "$CONFIG_DIR/secure.env" << EOF
{
  "variables": {
    "AWS_ACCESS_KEY_ID": "",
    "AWS_SECRET_ACCESS_KEY": "",
    "AZURE_CLIENT_ID": "",
    "AZURE_CLIENT_SECRET": "",
    "GCP_CREDENTIALS": "",
    "HETZNER_API_TOKEN": "",
    "OVH_APPLICATION_KEY": "",
    "OVH_APPLICATION_SECRET": "",
    "DIGITALOCEAN_TOKEN": ""
  }
}
EOF

chmod 600 "$CONFIG_DIR/secure.env"

# Install shell completions
echo "🔧 Installing shell completions..."
# Bash completion
mkdir -p "$HOME/.local/share/bash-completion/completions"
cat > "$HOME/.local/share/bash-completion/completions/engine" << 'EOF'
_engine_completion() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    commands="deploy list status cost help"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
    fi
}
complete -F _engine_completion engine
EOF

# Zsh completion
mkdir -p "$HOME/.zfunc"
cat > "$HOME/.zfunc/_engine" << 'EOF'
#compdef engine
_engine() {
    local -a commands
    commands=(
        'deploy:Deploy a composition'
        'list:List deployments'
        'status:Check deployment status'
        'cost:Get cost information'
        'help:Show help'
    )
    _describe 'command' commands
}
EOF

# Add to zsh path if not already there
if ! grep -q "fpath=($HOME/.zfunc)" "$HOME/.zshrc" 2>/dev/null; then
    echo "fpath=($HOME/.zfunc)" >> "$HOME/.zshrc"
fi

# Load completion in zsh
if ! grep -q "autoload -U compinit && compinit" "$HOME/.zshrc" 2>/dev/null; then
    echo "autoload -U compinit && compinit" >> "$HOME/.zshrc"
fi

echo "✅ Installation complete!"
echo ""
echo "🎉 Sovereign Engine v${VERSION} has been installed successfully."
echo ""
echo "📝 Configuration file: $CONFIG_DIR/config.yaml"
echo "🔒 Secure environment file: $CONFIG_DIR/secure.env"
echo "🔐 Master key: $MASTER_KEY"
echo ""
echo "🚀 Quick start:"
echo "   export ENGINE_MASTER_KEY=$MASTER_KEY"
echo "   engine-encrypt -key AWS_ACCESS_KEY_ID -value 'your_key' -file $CONFIG_DIR/secure.env"
echo "   engine deploy --provider aws --composition compute"
echo "   engine list"
echo "   engine status"
echo ""
echo "🌐 Web interface:"
echo "   engine-web  # Start web server on http://localhost:8080"
echo ""
echo "📚 For more information, run: engine help"
echo "🔐 For secure env management, run: engine-encrypt -h"
