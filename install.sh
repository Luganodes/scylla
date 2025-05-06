#!/bin/bash

# Configuration variables
REPO="Luganodes/scylla"
DESTINATION_DIR="/usr/local/bin"
BINARY_NAME="scylla"
SERVICE_FILE="/etc/systemd/system/scylla.service"
TMP_DIR="/tmp/scylla-install"

# Function to check OS type and architecture
check_system() {
  # Detect OS
  case "$(uname -s)" in
    Linux*)     OS="linux" ;;
    Darwin*)    OS="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
    *)          echo "Unsupported OS"; exit 1 ;;
  esac

  # Detect architecture
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64)     ARCH="amd64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    *)          echo "Unsupported architecture: $ARCH"; exit 1 ;;
  esac

  echo "Detected system: $OS/$ARCH"
}

# Function to get the latest release tag
get_latest_tag() {
  echo "Fetching latest release tag..."
  LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  
  if [ -z "$LATEST_TAG" ]; then
    echo "Failed to fetch latest tag. Falling back to tags list..."
    LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/tags" | grep '"name":' | head -n 1 | sed -E 's/.*"([^"]+)".*/\1/')
  fi
  
  if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not determine the latest release tag"
    exit 1
  fi
  
  echo "Latest release: $LATEST_TAG"
}

# Function to download and extract the binary
download_binary() {
  echo "Creating temporary directory..."
  mkdir -p "$TMP_DIR"
  
  DOWNLOAD_FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$DOWNLOAD_FILENAME"
  
  echo "Downloading $DOWNLOAD_URL..."
  if ! curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/$DOWNLOAD_FILENAME"; then
    echo "Error: Failed to download the binary"
    rm -rf "$TMP_DIR"
    exit 1
  fi
}

# Function to install the binary
install_binary() {
  # Check if the user is root or has sudo privileges
  if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo or as root."
    exit 1
  fi
  
  # Check if the destination directory exists
  if [ ! -d "$DESTINATION_DIR" ]; then
    echo "Creating $DESTINATION_DIR..."
    mkdir -p "$DESTINATION_DIR"
  fi
  
  # Copy the binary to the destination directory
  echo "Installing $BINARY_NAME to $DESTINATION_DIR..."
  cp "$TMP_DIR/${BINARY_NAME}-${OS}-${ARCH}" "$DESTINATION_DIR/$BINARY_NAME"
  chmod +x "$DESTINATION_DIR/$BINARY_NAME"
  
  # Clean up temporary directory
  rm -rf "$TMP_DIR"
  
  # Verify installation
  if [ -f "$DESTINATION_DIR/$BINARY_NAME" ]; then
    echo "$BINARY_NAME successfully installed in $DESTINATION_DIR"
  else
    echo "Installation failed."
    exit 1
  fi
}

# Function to install the service
install_service() {
  echo "Installing service..."
  
  # Create service file
  cat <<EOF | sudo tee $SERVICE_FILE > /dev/null
[Unit]
Description=Scylla Slashing Observer Service
After=network.target

[Service]
User=$USER
Type=simple
ExecStart=$DESTINATION_DIR/$BINARY_NAME start --ethereum.rpc https://ethereum-rpc.publicnode.com --ethereum.ws wss://ethereum-rpc.publicnode.com
Restart=on-failure
RestartSec=10
Environment="HOME=$HOME"

[Install]
WantedBy=multi-user.target
EOF

  echo "Reloading systemd daemon..."
  sudo systemctl daemon-reload
}

echo "Scylla Installer"
echo "----------------"

check_system

get_latest_tag

download_binary

install_binary

# Ask to install service
if [ "$OS" = "linux" ]; then
  read -p "Do you want to install Scylla as a system service? (y/n): " install_svc
  if [[ "$install_svc" =~ ^[Yy]$ ]]; then
    install_service
    echo ""
    echo "Service installed. To start Scylla, run:"
    echo "sudo systemctl enable scylla"
    echo "sudo systemctl start scylla"
    echo ""
    echo "To check status:"
    echo "sudo systemctl status scylla"
  fi
fi

# Success message
echo ""
echo "Installation complete!"
echo "Run '$BINARY_NAME --help' for usage information"