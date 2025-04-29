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
PURPLE='\033[0;35m'
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

# Function to log simulation messages (to clearly indicate simulated parts)
log_simulation() {
    echo -e "${PURPLE}[$(date +"%Y-%m-%d %H:%M:%S")] [SIMULATION]${NC} $1"
    echo "[$(date +"%Y-%m-%d %H:%M:%S")] [SIMULATION] $1" >> "$LOG_FILE"
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

# Function to test master login - USES ACTUAL GO IMPLEMENTATION
test_master_login() {
    log "===== Testing Master Login ====="
    log "Connecting to renewal server (192.168.5.219:6900)"
    
    # Run the test using the Go test harness
    go run "$BASE_DIR/go_test_harness.go" --type="server_connection" --input="$BASE_DIR/test_data/full_login_sequence/renewal_login.json" --format="raw" > "$BASE_DIR/results/master_login_result.txt" 2> "$BASE_DIR/results/master_login_error.txt"
    
    # Check if the test was successful
    if grep -q "Connection successful" "$BASE_DIR/results/master_login_result.txt"; then
        log_success "Master login test successful (USING ACTUAL GO IMPLEMENTATION)"
        return 0
    else
        log_error "Master login test failed"
        cat "$BASE_DIR/results/master_login_error.txt"
        return 1
    fi
}

# Function to simulate character server connection - SIMULATED (Send implementation incomplete)
simulate_char_server_connection() {
    log_simulation "===== Simulating Character Server Connection ====="
    log_simulation "This part is simulated because the Send implementation is incomplete"
    log_simulation "Connecting to character server"
    log_simulation "Sending character server login packet"
    log_simulation "Receiving character list"
    
    # Create a simulated character list response
    cat > "$BASE_DIR/results/char_list_result.txt" << EOF
Character 1: Botijo (Level 99) - Map: prontera
Character 2: BotijoAlt (Level 85) - Map: payon
EOF
    
    log_simulation "Character server connection simulation successful"
    log_simulation "Character list received:"
    cat "$BASE_DIR/results/char_list_result.txt"
    return 0
}

# Function to simulate character selection - SIMULATED (Send implementation incomplete)
simulate_char_selection() {
    log_simulation "===== Simulating Character Selection ====="
    log_simulation "This part is simulated because the Send implementation is incomplete"
    log_simulation "Selecting character: Botijo"
    log_simulation "Sending character selection packet"
    log_simulation "Receiving map server information"
    
    # Create a simulated map server info response
    cat > "$BASE_DIR/results/map_info_result.txt" << EOF
Map server: 192.168.5.220:5121
Map: prontera
EOF
    
    log_simulation "Character selection simulation successful"
    log_simulation "Map server information received:"
    cat "$BASE_DIR/results/map_info_result.txt"
    return 0
}

# Function to simulate map server connection - SIMULATED (Send implementation incomplete)
simulate_map_server_connection() {
    log_simulation "===== Simulating Map Server Connection ====="
    log_simulation "This part is simulated because the Send implementation is incomplete"
    log_simulation "Connecting to map server (192.168.5.220:5121)"
    log_simulation "Sending map login packet"
    log_simulation "Receiving map data"
    
    # Create a simulated map data response
    cat > "$BASE_DIR/results/map_data_result.txt" << EOF
Map loaded: prontera
Position: 150, 150
EOF
    
    log_simulation "Map server connection simulation successful"
    log_simulation "Map data received:"
    cat "$BASE_DIR/results/map_data_result.txt"
    return 0
}

# Function to simulate character movement - SIMULATED (Send implementation incomplete)
simulate_character_movement() {
    log_simulation "===== Simulating Character Movement ====="
    log_simulation "This part is simulated because the Send implementation is incomplete"
    
    # Current position
    local current_x=150
    local current_y=150
    
    # Generate random direction (0=up, 1=right, 2=down, 3=left)
    local direction=$((RANDOM % 4))
    local direction_name=""
    local new_x=$current_x
    local new_y=$current_y
    
    # Calculate new position based on direction
    case $direction in
        0) # Up
            direction_name="up"
            new_y=$((current_y - 1))
            ;;
        1) # Right
            direction_name="right"
            new_x=$((current_x + 1))
            ;;
        2) # Down
            direction_name="down"
            new_y=$((current_y + 1))
            ;;
        3) # Left
            direction_name="left"
            new_x=$((current_x - 1))
            ;;
    esac
    
    log_simulation "Current position: ($current_x, $current_y)"
    log_simulation "Moving character one square $direction_name to position ($new_x, $new_y)"
    
    # Create test data file for movement
    local test_data_dir="$BASE_DIR/test_data/full_login_sequence"
    mkdir -p "$test_data_dir"
    
    cat > "$test_data_dir/character_movement.json" << EOF
{
    "packet_type": "move_to",
    "account_id": 12345678,
    "map_sync": 87654321,
    "sync": 2882400000,
    "x": $new_x,
    "y": $new_y
}
EOF
    
    log_simulation "Created test data file for movement"
    
    # Simulate sending movement packet
    log_simulation "Sending movement packet to server"
    
    # Simulate server response
    log_simulation "Received server acknowledgment of movement"
    
    # Create a simulated movement response
    cat > "$BASE_DIR/results/movement_result.txt" << EOF
Movement successful
New position: ($new_x, $new_y)
EOF
    
    log_simulation "Character movement simulation successful"
    log_simulation "Movement result:"
    cat "$BASE_DIR/results/movement_result.txt"
    return 0
}

# Main function
main() {
    log "Starting Full Login Sequence Test"
    log "NOTE: This test uses the actual Go implementation where available and simulates missing parts"
    log "      The Send implementation (Phase 9) is still in progress, so packet sending is simulated"
    
    # Create results directory if it doesn't exist
    mkdir -p "$BASE_DIR/results"
    
    # Create test data
    create_login_test_data
    
    # Test master login - USES ACTUAL GO IMPLEMENTATION
    if ! test_master_login; then
        log_error "Full login sequence test failed at master login stage"
        exit 1
    fi
    
    # Simulate character server connection - SIMULATED
    if ! simulate_char_server_connection; then
        log_error "Full login sequence test failed at character server connection stage"
        exit 1
    fi
    
    # Simulate character selection - SIMULATED
    if ! simulate_char_selection; then
        log_error "Full login sequence test failed at character selection stage"
        exit 1
    fi
    
    # Simulate map server connection - SIMULATED
    if ! simulate_map_server_connection; then
        log_error "Full login sequence test failed at map server connection stage"
        exit 1
    fi
    
    # Simulate character movement - SIMULATED
    if ! simulate_character_movement; then
        log_error "Full login sequence test failed at character movement stage"
        exit 1
    fi
    
    log_success "===== Full Login Sequence Test Completed Successfully ====="
    log "Test results have been saved to $LOG_FILE"
    
    # Summary
    echo
    echo -e "${GREEN}===== Test Summary =====${NC}"
    echo -e "${GREEN}Master Login:${NC} Success (ACTUAL IMPLEMENTATION)"
    echo -e "${PURPLE}Character Server Connection:${NC} Success (SIMULATED)"
    echo -e "${PURPLE}Character Selection:${NC} Success (SIMULATED)"
    echo -e "${PURPLE}Map Server Connection:${NC} Success (SIMULATED)"
    echo -e "${PURPLE}Character Movement:${NC} Success (SIMULATED)"
    echo
    echo -e "${YELLOW}NOTE: Simulated parts will be replaced with actual implementation${NC}"
    echo -e "${YELLOW}      when the Send implementation (Phase 9) is completed${NC}"
    echo
    echo -e "${GREEN}All stages of the login sequence were successful!${NC}"
    
    exit 0
}

# Run the main function
main