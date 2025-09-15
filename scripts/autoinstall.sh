#!/usr/bin/env bash

set -e

GITHUB_REPO="andresantos-blip/gamedl"
BINARY_NAME="gamedl"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)
            os="Linux"
            ;;
        Darwin*)
            os="Darwin"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            log_error "Supported platforms: Linux, macOS (Darwin)"
            exit 1
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            log_error "Supported architectures: x86_64, arm64"
            exit 1
            ;;
    esac
    
    echo "${os}_${arch}"
}

# Get the latest release version
get_latest_version() {
    local api_url="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    local version
    
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null)
    else
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Check if we got a valid version
    if [[ -z "${version}" ]] || [[ "${version}" == *"Not Found"* ]]; then
        log_error "Failed to fetch latest release. The repository may be private or have no releases yet."
        log_error "You can specify a version manually by setting the VERSION environment variable:"
        log_error "  VERSION=v1.0.0 bash <(curl -sSfL <installer-url>)"
        return 1
    fi
    
    echo "${version}"
}

# Download and extract the binary
download_and_install() {
    local platform="$1"
    local version="$2"
    local archive_name="${BINARY_NAME}_${platform}.tar.gz"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${archive_name}"
    local temp_dir
    
    # Create temporary directory
    temp_dir=$(mktemp -d)
    trap "rm -rf ${temp_dir}" EXIT
    
    log_info "Downloading ${BINARY_NAME} ${version} for ${platform}..."
    log_info "Download URL: ${download_url}"
    
    # Download the archive
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "${temp_dir}/${archive_name}" "${download_url}"; then
            log_error "Failed to download ${archive_name}"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "${temp_dir}/${archive_name}" "${download_url}"; then
            log_error "Failed to download ${archive_name}"
            exit 1
        fi
    else
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Extract the archive
    log_info "Extracting ${archive_name}..."
    if ! tar -xzf "${temp_dir}/${archive_name}" -C "${temp_dir}"; then
        log_error "Failed to extract ${archive_name}"
        exit 1
    fi
    
    # Find install location
    local install_dir
    local binary_path="${temp_dir}/${BINARY_NAME}"
    
    if [[ ! -f "${binary_path}" ]]; then
        log_error "Binary ${BINARY_NAME} not found in archive"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "${binary_path}"
    
    # Determine installation directory
    if [[ -w "/usr/local/bin" ]]; then
        install_dir="/usr/local/bin"
    elif [[ -w "${HOME}/.local/bin" ]]; then
        install_dir="${HOME}/.local/bin"
        # Create directory if it doesn't exist
        mkdir -p "${install_dir}"
    elif [[ -w "${HOME}/bin" ]]; then
        install_dir="${HOME}/bin"
        # Create directory if it doesn't exist
        mkdir -p "${install_dir}"
    else
        log_warn "No writable directory found in PATH. Attempting to install to /usr/local/bin with sudo..."
        install_dir="/usr/local/bin"
        sudo=true
    fi
    
    # Install the binary
    log_info "Installing ${BINARY_NAME} to ${install_dir}..."
    if [[ "${sudo:-false}" == "true" ]]; then
        if ! sudo cp "${binary_path}" "${install_dir}/"; then
            log_error "Failed to install ${BINARY_NAME} to ${install_dir}"
            log_error "You may need to manually copy the binary from: ${binary_path}"
            exit 1
        fi
    else
        if ! cp "${binary_path}" "${install_dir}/"; then
            log_error "Failed to install ${BINARY_NAME} to ${install_dir}"
            log_error "You may need to manually copy the binary from: ${binary_path}"
            exit 1
        fi
    fi
    
    log_info "Successfully installed ${BINARY_NAME} to ${install_dir}/${BINARY_NAME}"
    
    # Check if the install directory is in PATH
    if [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        log_warn "Warning: ${install_dir} is not in your PATH"
        log_warn "You may need to add it to your PATH or run the binary with the full path:"
        log_warn "  ${install_dir}/${BINARY_NAME}"
        
        # Provide instructions for adding to PATH
        case "${install_dir}" in
            "${HOME}/.local/bin"|"${HOME}/bin")
                log_info "To add ${install_dir} to your PATH, add this line to your shell profile:"
                log_info "  export PATH=\"${install_dir}:\$PATH\""
                ;;
        esac
    fi
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        local installed_version
        installed_version=$("${BINARY_NAME}" version 2>/dev/null || "${BINARY_NAME}" --version 2>/dev/null || echo "unknown")
        log_info "✓ ${BINARY_NAME} is installed and available in PATH"
        log_info "  Version: ${installed_version}"
    else
        log_warn "✗ ${BINARY_NAME} is not in PATH, but may be installed in a specific directory"
        log_info "Try running: ${BINARY_NAME} version"
    fi
}

# Main installation flow
main() {
    log_info "Installing ${BINARY_NAME}..."
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: ${platform}"
    
    # Get version (from environment or latest release)
    local version
    if [[ -n "${VERSION:-}" ]]; then
        version="${VERSION}"
        log_info "Using specified version: ${version}"
    else
        version=$(get_latest_version)
        if [[ $? -ne 0 ]] || [[ -z "${version}" ]]; then
            log_error "Failed to get latest version from GitHub releases"
            exit 1
        fi
        log_info "Latest version: ${version}"
    fi
    
    # Download and install
    download_and_install "${platform}" "${version}"
    
    # Verify installation
    verify_installation
    
    log_info "Installation complete! 🎉"
    log_info "Run '${BINARY_NAME} --help' to get started."
}

# Check if running with bash
if [[ "${BASH_VERSION:-}" == "" ]]; then
    log_error "This script requires bash. Please run with: bash <(curl -sSfL <url>)"
    exit 1
fi

# Run main function
main "$@"