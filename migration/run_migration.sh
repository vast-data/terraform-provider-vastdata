#!/bin/bash
# Copyright (c) HashiCorp, Inc.


# Terraform Provider Migration Script Runner
# Migrates VastData Terraform configurations from provider 1.x to 2.0
#
# USAGE EXAMPLES:
#   # Run migration
#   ./run_migration.sh /path/to/old/configs /path/to/migrated/configs
#
#   # Show help
#   ./run_migration.sh --help
#
#   # Show version
#   ./run_migration.sh --version
#
#   # Run tests
#   ./run_migration.sh --test
#
#   # Clean up environment
#   ./run_migration.sh --clean
#
# REQUIREMENTS:
#   - Python 3.9 or higher
#   - terraform command (for validation)
#
# WHAT IT DOES:
#   - Renames resources: vastdata_administators_managers → vastdata_administrator_manager
#   - Updates attributes: type_ → type, permissions_list → permissions  
#   - Transforms schemas: Block lists to attributes, IP ranges, etc.
#   - Preserves dynamic blocks (requires manual review)
#   - Validates results with terraform init/apply

set -e  # Exit on any error

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATION_SCRIPT="$SCRIPT_DIR/migration_script.py"
VENV_DIR="$SCRIPT_DIR/venv"
MIN_PYTHON_VERSION="3.9"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to compare version numbers
version_compare() {
    local version1="$1"
    local version2="$2"
    
    if [[ "$version1" == "$version2" ]]; then
        return 0
    fi
    
    local IFS=.
    local ver1=($version1)
    local ver2=($version2)
    
    # Fill empty fields with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++)); do
        ver1[i]=0
    done
    for ((i=${#ver2[@]}; i<${#ver1[@]}; i++)); do
        ver2[i]=0
    done
    
    for ((i=0; i<${#ver1[@]}; i++)); do
        if [[ -z ${ver2[i]} ]]; then
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]})); then
            return 1  # version1 > version2
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]})); then
            return 2  # version1 < version2
        fi
    done
    return 0  # versions are equal
}

# Function to check Python version
check_python_version() {
    local python_cmd="$1"
    
    if ! command -v "$python_cmd" &> /dev/null; then
        return 1
    fi
    
    local python_version
    python_version=$("$python_cmd" -c "import sys; print('.'.join(map(str, sys.version_info[:2])))" 2>/dev/null)
    
    if [[ -z "$python_version" ]]; then
        return 1
    fi
    
    version_compare "$python_version" "$MIN_PYTHON_VERSION"
    local result=$?
    
    if [[ $result -eq 0 ]] || [[ $result -eq 1 ]]; then
        echo "$python_version"
        return 0
    else
        return 1
    fi
}

# Function to find suitable Python executable
find_python() {
    local python_candidates=("python3.12" "python3.11" "python3.10" "python3.9" "python3" "python")
    
    for python_cmd in "${python_candidates[@]}"; do
        local version
        if version=$(check_python_version "$python_cmd"); then
            log_success "Found suitable Python: $python_cmd (version $version)" >&2
            echo "$python_cmd"
            return 0
        fi
    done
    
    return 1
}

# Function to create virtual environment
create_venv() {
    local python_cmd="$1"
    
    if [[ -d "$VENV_DIR" ]]; then
        log_info "Virtual environment already exists at $VENV_DIR"
        return 0
    fi
    
    log_info "Creating virtual environment at $VENV_DIR"
    
    if ! "$python_cmd" -m venv "$VENV_DIR"; then
        log_error "Failed to create virtual environment"
        return 1
    fi
    
    log_success "Virtual environment created successfully"
    return 0
}

# Function to activate virtual environment
activate_venv() {
    local venv_activate="$VENV_DIR/bin/activate"
    
    if [[ ! -f "$venv_activate" ]]; then
        log_error "Virtual environment activation script not found at $venv_activate"
        return 1
    fi
    
    log_info "Activating virtual environment"
    # shellcheck source=/dev/null
    source "$venv_activate"
    
    # Verify activation
    if [[ "$VIRTUAL_ENV" != "$VENV_DIR" ]]; then
        log_error "Failed to activate virtual environment"
        return 1
    fi
    
    log_success "Virtual environment activated"
    return 0
}

# Function to check if migration script exists
check_migration_script() {
    if [[ ! -f "$MIGRATION_SCRIPT" ]]; then
        log_error "Migration script not found at $MIGRATION_SCRIPT"
        return 1
    fi
    
    if [[ ! -x "$MIGRATION_SCRIPT" ]]; then
        log_info "Making migration script executable"
        chmod +x "$MIGRATION_SCRIPT"
    fi
    
    return 0
}

# Function to run migration script
run_migration() {
    log_info "Running migration script: $MIGRATION_SCRIPT"
    
    # Check if arguments were provided
    if [[ $# -eq 0 ]]; then
        log_info "No arguments provided, showing migration script help:"
        python "$MIGRATION_SCRIPT" --help
        return 0
    fi
    
    # Run the migration script with provided arguments
    python "$MIGRATION_SCRIPT" "$@"
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        log_success "Migration script completed successfully"
    else
        log_error "Migration script failed with exit code $exit_code"
    fi
    
    return $exit_code
}

# Function to run tests
run_tests() {
    log_info "Running migration tool tests"
    
    # Find suitable Python executable
    local python_cmd
    if ! python_cmd=$(find_python); then
        log_error "Python $MIN_PYTHON_VERSION or higher is required"
        return 1
    fi
    
    # Create virtual environment if needed
    if ! create_venv "$python_cmd"; then
        return 1
    fi
    
    # Activate virtual environment
    if ! activate_venv; then
        return 1
    fi
    
    # Check if test dependencies are available
    if ! python -c "import pytest" 2>/dev/null; then
        log_warning "pytest not found, installing test dependencies..."
        if [[ -f "requirements-test.txt" ]]; then
            pip install -r requirements-test.txt
        else
            pip install pytest
        fi
    fi
    
    # Run pytest
    log_info "Executing test suite..."
    python -m pytest tests/ -v
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        log_success "All tests passed successfully"
    else
        log_error "Some tests failed (exit code: $exit_code)"
    fi
    
    return $exit_code
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS] [src_folder] [dst_folder]

Migrates VastData Terraform configurations from provider 1.x to 2.0.

Arguments:
  src_folder    Source folder containing .tf files to migrate
  dst_folder    Destination folder for converted files

Options:
  -h, --help    Show this help message
  --version     Show migration script version
  --test        Run migration tool test suite
  --clean       Remove virtual environment and exit

Examples:
  $0 ./old_configs ./new_configs    # Run migration
  $0 --test                         # Run tests
  $0 --version                      # Show version
  $0 --clean                        # Clean up

Requirements:
  - Python $MIN_PYTHON_VERSION or higher
  - terraform command (for validation)

EOF
}

# Function to clean up virtual environment
cleanup_venv() {
    if [[ -d "$VENV_DIR" ]]; then
        log_info "Removing virtual environment at $VENV_DIR"
        rm -rf "$VENV_DIR"
        log_success "Virtual environment removed"
    else
        log_info "No virtual environment found to clean up"
    fi
}

# Main function
main() {
    log_info "Starting Terraform Provider Migration Script Runner"
    
    # Parse command line arguments
    case "${1:-}" in
        -h|--help)
            show_usage
            exit 0
            ;;
        --clean)
            cleanup_venv
            exit 0
            ;;
        --version)
            if check_migration_script; then
                # Find suitable Python executable for version check
                local python_cmd
                if python_cmd=$(find_python); then
                    "$python_cmd" "$MIGRATION_SCRIPT" --version
                else
                    log_error "Python $MIN_PYTHON_VERSION or higher is required"
                    exit 1
                fi
            fi
            exit 0
            ;;
        --test)
            run_tests
            exit 0
            ;;
    esac
    
    # Check if migration script exists
    if ! check_migration_script; then
        exit 1
    fi
    
    # Find suitable Python executable
    local python_cmd
    if ! python_cmd=$(find_python); then
        log_error "Python $MIN_PYTHON_VERSION or higher is required but not found"
        log_error "Please install Python $MIN_PYTHON_VERSION+ and ensure it's in your PATH"
        exit 1
    fi
    
    # Create virtual environment
    if ! create_venv "$python_cmd"; then
        exit 1
    fi
    
    # Activate virtual environment
    if ! activate_venv; then
        exit 1
    fi
    
    # Check if terraform is available (required by migration script)
    if ! command -v terraform &> /dev/null; then
        log_warning "terraform command not found in PATH"
        log_warning "The migration script requires terraform for validation"
        log_warning "Please ensure terraform is installed and available in PATH"
    fi
    
    # Run migration script with provided arguments
    run_migration "$@"
    local exit_code=$?
    
    log_info "Migration script runner completed"
    exit $exit_code
}

# Run main function with all arguments
main "$@"
