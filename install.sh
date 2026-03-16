#!/bin/bash
set -e

# WPSH Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/shariffff/wpsh/main/install.sh | bash

REPO="shariffff/wpsh"
BINARY_NAME="wpsh"

# Install directory (like Bun's ~/.bun)
install_dir="${WORDMON_INSTALL:-$HOME/.wpsh}"
bin_dir="$install_dir/bin"
ansible_dir="$install_dir/ansible"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

info() { echo -e "${GREEN}$1${NC}"; }
warn() { echo -e "${YELLOW}$1${NC}"; }
error() { echo -e "${RED}error${NC}: $1" >&2; exit 1; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Add to shell config and return the config file path
setup_shell() {
    local shell_name=$(basename "$SHELL")
    local config_file=""
    local shell_export=""

    case $shell_name in
        fish)
            config_file="$HOME/.config/fish/config.fish"
            shell_export="set --export WORDMON_INSTALL \"$install_dir\"\nset --export PATH \$WORDMON_INSTALL/bin \$PATH"
            ;;
        zsh)
            config_file="$HOME/.zshrc"
            shell_export="export WORDMON_INSTALL=\"$install_dir\"\nexport PATH=\"\$WORDMON_INSTALL/bin:\$PATH\""
            ;;
        bash)
            if [[ -f "$HOME/.bashrc" ]]; then
                config_file="$HOME/.bashrc"
            else
                config_file="$HOME/.bash_profile"
            fi
            shell_export="export WORDMON_INSTALL=\"$install_dir\"\nexport PATH=\"\$WORDMON_INSTALL/bin:\$PATH\""
            ;;
        *)
            echo ""
            warn "Could not detect shell. Manually add to your shell config:"
            echo ""
            echo "  export WORDMON_INSTALL=\"$install_dir\""
            echo "  export PATH=\"\$WORDMON_INSTALL/bin:\$PATH\""
            echo ""
            return
            ;;
    esac

    # Check if already configured
    if [[ -f "$config_file" ]] && grep -q "WORDMON_INSTALL" "$config_file" 2>/dev/null; then
        # Already configured, just export for current session
        export WORDMON_INSTALL="$install_dir"
        export PATH="$bin_dir:$PATH"
        return
    fi

    # Create config file if it doesn't exist
    if [[ ! -f "$config_file" ]]; then
        mkdir -p "$(dirname "$config_file")"
        touch "$config_file"
    fi

    # Check if writable
    if [[ ! -w "$config_file" ]]; then
        warn "Could not write to $config_file. Manually add:"
        echo ""
        echo -e "  $shell_export"
        echo ""
        return
    fi

    # Append to config
    {
        echo ""
        echo "# WPSH"
        echo -e "$shell_export"
    } >> "$config_file"

    echo -e "${DIM}Added to $config_file${NC}"

    # Export for current session
    export WORDMON_INSTALL="$install_dir"
    export PATH="$bin_dir:$PATH"
}

# Download and install
install_wpsh() {
    local os=$(detect_os)
    local arch=$(detect_arch)
    local version="${WORDMON_VERSION:-$(get_latest_version)}"

    if [[ -z "$version" ]]; then
        error "Could not determine latest version. Set WORDMON_VERSION manually."
    fi

    # Remove 'v' prefix if present for filename
    local version_num="${version#v}"

    echo -e "${DIM}Installing WPSH ${version} (${os}/${arch})${NC}"

    # Construct download URL
    local archive_name="wpsh_${version_num}_${os}_${arch}"
    local filename="${archive_name}"
    if [[ "$os" = "windows" ]]; then
        filename="${filename}.zip"
    else
        filename="${filename}.tar.gz"
    fi

    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"

    # Create temp directory
    local tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT

    # Download
    echo -e "${DIM}Downloading...${NC}"
    if ! curl -fsSL "$url" -o "${tmp_dir}/${filename}" 2>/dev/null; then
        error "Failed to download from ${url}"
    fi

    # Extract
    cd "${tmp_dir}"
    if [[ "$os" = "windows" ]]; then
        unzip -q "${filename}"
    else
        tar -xzf "${filename}"
    fi

    # Create install directories
    mkdir -p "${bin_dir}"

    # Find and install binary
    local binary="${BINARY_NAME}"
    if [[ "$os" = "windows" ]]; then
        binary="${BINARY_NAME}.exe"
    fi

    # Binary is inside the archive_name folder
    if [[ -f "${archive_name}/${binary}" ]]; then
        mv "${archive_name}/${binary}" "${bin_dir}/"
    else
        error "Binary not found in archive"
    fi

    chmod +x "${bin_dir}/${binary}"

    # Install ansible directory
    if [[ -d "${archive_name}/ansible" ]]; then
        # Remove old ansible dir if exists
        rm -rf "${ansible_dir}"
        mv "${archive_name}/ansible" "${ansible_dir}"
        echo -e "${DIM}Installed ansible playbooks to ${ansible_dir}${NC}"
    else
        warn "Ansible directory not found in archive"
    fi
}

# Main
main() {
    echo ""
    echo -e "${BOLD}WPSH${NC} Installer"
    echo ""

    # Check for required tools
    command -v curl >/dev/null 2>&1 || error "curl is required but not installed"
    command -v tar >/dev/null 2>&1 || error "tar is required but not installed"

    install_wpsh
    setup_shell

    echo ""
    echo -e "${GREEN}WPSH was installed successfully!${NC}"
    echo ""
    echo "Run the following to get started:"
    echo ""
    echo -e "  ${BOLD}source ~/.$(basename $SHELL)rc && wpsh init${NC}"
    echo ""
    echo -e "${DIM}Or restart your terminal and run: wpsh init${NC}"
    echo ""
}

main "$@"
