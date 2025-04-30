#!/bin/bash

# Full Login Sequence Test Script
# This script tests the complete login sequence using the Go implementation
# It targets the renewal server with the login credentials botijo0/Melon.77

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$BASE_DIR/full_login_sequence_test.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

# Create log file
> "$LOG_FILE"

# Function to log messages
log() {
    echo -e "${BLUE}[$(date +"%Y-%m-%d %H:%M:%S")]${NC} $1"
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] $1" >> "$LOG_FILE"
}

# Function to log success messages
log_success() {
    echo -e "${GREEN}[$(date +"%Y-%m-%d %H:%M:%S")]${NC} $1"
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] $1" >> "$LOG_FILE"
}

# Function to log error messages
log_error() {
    echo -e "${RED}[$(date +"%Y-%m-%d %H:%M:%S")]${NC} $1"
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] $1" >> "$LOG_FILE"
}

# Function to log warning messages
log_warning() {
    echo -e "${YELLOW}[$(date +"%Y-%m-%d %H:%M:%S")]${NC} $1"
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] $1" >> "$LOG_FILE"
}

# Create test data file for login
create_login_test_data() {
    local test_data_dir="$BASE_DIR/test_data/full_login_sequence"
    mkdir -p "$test_data_dir"
    
    cat > "$test_data_dir/renewal_login.json" << EOF
{
    "server_type": "0",
    "server_ip": "192.168.5.219",
    "server_port": 6900,
    "username": "botijo0",
    "password": "Melon.77",
    "version": 1,
    "gender": 0
}
EOF
    
    log "Created test data file: $test_data_dir/renewal_login.json"
}

# Function to run the full login sequence test
run_full_login_sequence_test() {
    log "===== Running Full Login Sequence Test ====="
    
    # Create results directory if it doesn't exist
    mkdir -p "$BASE_DIR/results"
    
    # Build the test binary
    log "Building full_login_sequence_test binary..."
    
    # Change to the verification directory
    pushd "$BASE_DIR" > /dev/null
    
    # Build the test binary
    go build -o "full_login_sequence_test_bin" "full_login_sequence_test.go"
    
    if [ $? -ne 0 ]; then
        log_error "Failed to build full_login_sequence_test binary"
        popd > /dev/null
        return 1
    fi
    
    log_success "Successfully built full_login_sequence_test binary"
    
    # Run the test
    log "Running full login sequence test..."
    
    ./full_login_sequence_test_bin > "$BASE_DIR/results/full_login_sequence_result.txt" 2> "$BASE_DIR/results/full_login_sequence_error.txt"
    
    # Check if the test was successful
    if grep -q "Full login sequence test passed!" "$BASE_DIR/results/full_login_sequence_result.txt"; then
        log_success "Full login sequence test successful"
        popd > /dev/null
        return 0
    else
        log_error "Full login sequence test failed"
        cat "$BASE_DIR/results/full_login_sequence_error.txt"
        popd > /dev/null
        return 1
    fi
}

# Main function
main() {
    log "Starting Full Login Sequence Test"
    log "This test uses the actual Go implementation for the complete login sequence"
    
    # Create test data
    create_login_test_data
    
    # Run the full login sequence test
    if ! run_full_login_sequence_test; then
        log_error "Full login sequence test failed"
        exit 1
    fi
    
    log_success "===== Full Login Sequence Test Completed Successfully ====="
    log "Test results have been saved to $LOG_FILE"
    
    # Summary
    echo
    echo -e "${GREEN}===== Test Summary =====${NC}"
    echo -e "${GREEN}Full Login Sequence:${NC} Success (ACTUAL IMPLEMENTATION)"
    echo
    echo -e "${GREEN}All stages of the login sequence were successful!${NC}"
    
    exit 0
}

# Run the main function
main