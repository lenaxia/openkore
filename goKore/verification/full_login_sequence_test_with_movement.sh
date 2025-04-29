#!/bin/bash

# Full Login Sequence Test with Movement
# This script tests the complete login sequence using the actual Go implementation
# It targets the renewal server with the login credentials botijo0/Melon.77
# and includes character movement verification

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$BASE_DIR/full_login_sequence_test_with_movement.log"

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

# Create test data files
create_test_data() {
    local test_data_dir="$BASE_DIR/test_data/full_login_sequence"
    mkdir -p "$test_data_dir"
    
    # Master login test data
    cat > "$test_data_dir/master_login.json" << EOF
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
    
    # Character server login test data
    cat > "$test_data_dir/char_server_login.json" << EOF
{
    "server_type": "0",
    "account_id": 0,
    "session_id1": 0,
    "session_id2": 0,
    "gender": 0
}
EOF
    
    # Character selection test data
    cat > "$test_data_dir/char_select.json" << EOF
{
    "char_num": 0
}
EOF
    
    # Map server login test data
    cat > "$test_data_dir/map_login.json" << EOF
{
    "account_id": 0,
    "char_id": 0,
    "session_id1": 0,
    "session_id2": 0,
    "gender": 0
}
EOF
    
    # Character movement test data
    cat > "$test_data_dir/movement.json" << EOF
{
    "from_x": 0,
    "from_y": 0,
    "to_x": 0,
    "to_y": 0
}
EOF
    
    log "Created test data files in $test_data_dir"
}

# Function to run the master login test
run_master_login() {
    log "===== Testing Master Login ====="
    log "Using actual Go implementation for master login"
    
    # Create results directory if it doesn't exist
    mkdir -p "$BASE_DIR/results"
    
    # Run the test using the Go implementation
    cd "$BASE_DIR/../.."
    
    # Use the actual Go implementation to connect to the server and login
    log "Executing master login with actual Go implementation"
    go run ./cmd/login/main.go \
        --server 192.168.5.219 \
        --port 6900 \
        --username botijo0 \
        --password Melon.77 \
        --verbose > "$BASE_DIR/results/master_login_result.txt" 2>&1
    
    # Check if the login was successful
    if grep -q "Login successful" "$BASE_DIR/results/master_login_result.txt"; then
        log_success "Master login successful"
        
        # Extract account ID and session IDs from the output
        ACCOUNT_ID=$(grep "Account ID:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $3}')
        SESSION_ID1=$(grep "Session ID 1:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $4}')
        SESSION_ID2=$(grep "Session ID 2:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $4}')
        GENDER=$(grep "Gender:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $2}')
        
        log "Account ID: $ACCOUNT_ID"
        log "Session ID 1: $SESSION_ID1"
        log "Session ID 2: $SESSION_ID2"
        log "Gender: $GENDER"
        
        # Update the character server login test data with the extracted values
        sed -i "s/\"account_id\": 0/\"account_id\": $ACCOUNT_ID/" "$BASE_DIR/test_data/full_login_sequence/char_server_login.json"
        sed -i "s/\"session_id1\": 0/\"session_id1\": $SESSION_ID1/" "$BASE_DIR/test_data/full_login_sequence/char_server_login.json"
        sed -i "s/\"session_id2\": 0/\"session_id2\": $SESSION_ID2/" "$BASE_DIR/test_data/full_login_sequence/char_server_login.json"
        sed -i "s/\"gender\": 0/\"gender\": $GENDER/" "$BASE_DIR/test_data/full_login_sequence/char_server_login.json"
        
        # Also update the map login test data
        sed -i "s/\"account_id\": 0/\"account_id\": $ACCOUNT_ID/" "$BASE_DIR/test_data/full_login_sequence/map_login.json"
        sed -i "s/\"session_id1\": 0/\"session_id1\": $SESSION_ID1/" "$BASE_DIR/test_data/full_login_sequence/map_login.json"
        sed -i "s/\"session_id2\": 0/\"session_id2\": $SESSION_ID2/" "$BASE_DIR/test_data/full_login_sequence/map_login.json"
        sed -i "s/\"gender\": 0/\"gender\": $GENDER/" "$BASE_DIR/test_data/full_login_sequence/map_login.json"
        
        # Extract server list
        SERVER_IP=$(grep "Server IP:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $3}')
        SERVER_PORT=$(grep "Server Port:" "$BASE_DIR/results/master_login_result.txt" | awk '{print $3}')
        
        log "Character Server IP: $SERVER_IP"
        log "Character Server Port: $SERVER_PORT"
        
        return 0
    else
        log_error "Master login failed"
        cat "$BASE_DIR/results/master_login_result.txt"
        return 1
    fi
}

# Function to run the character server login test
run_char_server_login() {
    log "===== Testing Character Server Login ====="
    log "Using actual Go implementation for character server login"
    
    # Run the test using the Go implementation
    cd "$BASE_DIR/../.."
    
    # Use the actual Go implementation to connect to the character server
    log "Executing character server login with actual Go implementation"
    go run ./cmd/charlogin/main.go \
        --server "$SERVER_IP" \
        --port "$SERVER_PORT" \
        --account-id "$ACCOUNT_ID" \
        --session-id1 "$SESSION_ID1" \
        --session-id2 "$SESSION_ID2" \
        --gender "$GENDER" \
        --verbose > "$BASE_DIR/results/char_server_login_result.txt" 2>&1
    
    # Check if the login was successful
    if grep -q "Character server login successful" "$BASE_DIR/results/char_server_login_result.txt"; then
        log_success "Character server login successful"
        
        # Extract character list
        log "Character list:"
        grep "Character:" "$BASE_DIR/results/char_server_login_result.txt" | while read -r line; do
            log "  $line"
        done
        
        # Extract character ID for the first character
        CHAR_ID=$(grep "Character ID:" "$BASE_DIR/results/char_server_login_result.txt" | head -1 | awk '{print $3}')
        log "Selected Character ID: $CHAR_ID"
        
        # Update the map login test data with the character ID
        sed -i "s/\"char_id\": 0/\"char_id\": $CHAR_ID/" "$BASE_DIR/test_data/full_login_sequence/map_login.json"
        
        return 0
    else
        log_error "Character server login failed"
        cat "$BASE_DIR/results/char_server_login_result.txt"
        return 1
    fi
}

# Function to run the character selection test
run_char_selection() {
    log "===== Testing Character Selection ====="
    log "Using actual Go implementation for character selection"
    
    # Run the test using the Go implementation
    cd "$BASE_DIR/../.."
    
    # Use the actual Go implementation to select a character
    log "Executing character selection with actual Go implementation"
    go run ./cmd/charselect/main.go \
        --server "$SERVER_IP" \
        --port "$SERVER_PORT" \
        --char-num 0 \
        --verbose > "$BASE_DIR/results/char_selection_result.txt" 2>&1
    
    # Check if the selection was successful
    if grep -q "Character selection successful" "$BASE_DIR/results/char_selection_result.txt"; then
        log_success "Character selection successful"
        
        # Extract map server information
        MAP_SERVER_IP=$(grep "Map Server IP:" "$BASE_DIR/results/char_selection_result.txt" | awk '{print $4}')
        MAP_SERVER_PORT=$(grep "Map Server Port:" "$BASE_DIR/results/char_selection_result.txt" | awk '{print $4}')
        
        log "Map Server IP: $MAP_SERVER_IP"
        log "Map Server Port: $MAP_SERVER_PORT"
        
        return 0
    else
        log_error "Character selection failed"
        cat "$BASE_DIR/results/char_selection_result.txt"
        return 1
    fi
}

# Function to run the map server login test
run_map_server_login() {
    log "===== Testing Map Server Login ====="
    log "Using actual Go implementation for map server login"
    
    # Run the test using the Go implementation
    cd "$BASE_DIR/../.."
    
    # Use the actual Go implementation to connect to the map server
    log "Executing map server login with actual Go implementation"
    go run ./cmd/maplogin/main.go \
        --server "$MAP_SERVER_IP" \
        --port "$MAP_SERVER_PORT" \
        --account-id "$ACCOUNT_ID" \
        --char-id "$CHAR_ID" \
        --session-id1 "$SESSION_ID1" \
        --session-id2 "$SESSION_ID2" \
        --gender "$GENDER" \
        --verbose > "$BASE_DIR/results/map_server_login_result.txt" 2>&1
    
    # Check if the login was successful
    if grep -q "Map login successful" "$BASE_DIR/results/map_server_login_result.txt"; then
        log_success "Map server login successful"
        
        # Extract position information
        POSITION_X=$(grep "Position:" "$BASE_DIR/results/map_server_login_result.txt" | awk '{print $2}' | tr -d '(,')
        POSITION_Y=$(grep "Position:" "$BASE_DIR/results/map_server_login_result.txt" | awk '{print $3}' | tr -d ')')
        MAP_NAME=$(grep "Map:" "$BASE_DIR/results/map_server_login_result.txt" | awk '{print $2}')
        
        log "Map: $MAP_NAME"
        log "Position: ($POSITION_X, $POSITION_Y)"
        
        # Update the movement test data with the position
        sed -i "s/\"from_x\": 0/\"from_x\": $POSITION_X/" "$BASE_DIR/test_data/full_login_sequence/movement.json"
        sed -i "s/\"from_y\": 0/\"from_y\": $POSITION_Y/" "$BASE_DIR/test_data/full_login_sequence/movement.json"
        
        # Calculate new position (one square in a random direction)
        DIRECTION=$((RANDOM % 4))
        NEW_X=$POSITION_X
        NEW_Y=$POSITION_Y
        
        case $DIRECTION in
            0) # Up
                NEW_Y=$((POSITION_Y - 1))
                DIRECTION_NAME="up"
                ;;
            1) # Right
                NEW_X=$((POSITION_X + 1))
                DIRECTION_NAME="right"
                ;;
            2) # Down
                NEW_Y=$((POSITION_Y + 1))
                DIRECTION_NAME="down"
                ;;
            3) # Left
                NEW_X=$((POSITION_X - 1))
                DIRECTION_NAME="left"
                ;;
        esac
        
        log "Moving character one square $DIRECTION_NAME to position ($NEW_X, $NEW_Y)"
        
        # Update the movement test data with the new position
        sed -i "s/\"to_x\": 0/\"to_x\": $NEW_X/" "$BASE_DIR/test_data/full_login_sequence/movement.json"
        sed -i "s/\"to_y\": 0/\"to_y\": $NEW_Y/" "$BASE_DIR/test_data/full_login_sequence/movement.json"
        
        return 0
    else
        log_error "Map server login failed"
        cat "$BASE_DIR/results/map_server_login_result.txt"
        return 1
    fi
}

# Function to run the character movement test
run_character_movement() {
    log "===== Testing Character Movement ====="
    log "Using actual Go implementation for character movement"
    
    # Run the test using the Go implementation
    cd "$BASE_DIR/../.."
    
    # Use the actual Go implementation to move the character
    log "Executing character movement with actual Go implementation"
    go run ./cmd/move/main.go \
        --server "$MAP_SERVER_IP" \
        --port "$MAP_SERVER_PORT" \
        --x "$NEW_X" \
        --y "$NEW_Y" \
        --verbose > "$BASE_DIR/results/character_movement_result.txt" 2>&1
    
    # Check if the movement was successful
    if grep -q "Movement successful" "$BASE_DIR/results/character_movement_result.txt"; then
        log_success "Character movement successful"
        log "Character moved from ($POSITION_X, $POSITION_Y) to ($NEW_X, $NEW_Y)"
        return 0
    else
        log_error "Character movement failed"
        cat "$BASE_DIR/results/character_movement_result.txt"
        return 1
    fi
}

# Main function
main() {
    log "Starting Full Login Sequence Test with Movement"
    log "Using actual Go implementation for all steps"
    
    # Create test data
    create_test_data
    
    # Run master login test
    if ! run_master_login; then
        log_error "Full login sequence test failed at master login stage"
        exit 1
    fi
    
    # Run character server login test
    if ! run_char_server_login; then
        log_error "Full login sequence test failed at character server login stage"
        exit 1
    fi
    
    # Run character selection test
    if ! run_char_selection; then
        log_error "Full login sequence test failed at character selection stage"
        exit 1
    fi
    
    # Run map server login test
    if ! run_map_server_login; then
        log_error "Full login sequence test failed at map server login stage"
        exit 1
    fi
    
    # Run character movement test
    if ! run_character_movement; then
        log_error "Full login sequence test failed at character movement stage"
        exit 1
    fi
    
    log_success "===== Full Login Sequence Test with Movement Completed Successfully ====="
    log "Test results have been saved to $LOG_FILE"
    
    # Summary
    echo
    echo -e "${GREEN}===== Test Summary =====${NC}"
    echo -e "${GREEN}Master Login:${NC} Success"
    echo -e "${GREEN}Character Server Login:${NC} Success"
    echo -e "${GREEN}Character Selection:${NC} Success"
    echo -e "${GREEN}Map Server Login:${NC} Success"
    echo -e "${GREEN}Character Movement:${NC} Success"
    echo
    echo -e "${GREEN}All stages of the login sequence were successful!${NC}"
    
    exit 0
}

# Run the main function
main