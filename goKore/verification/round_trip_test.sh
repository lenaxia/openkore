#!/bin/bash
# Script to perform round-trip testing for message ID encryption

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$BASE_DIR/test_data"
RESULTS_DIR="$BASE_DIR/results"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to run a round-trip test for message ID encryption
test_message_id_round_trip() {
    local test_file=$1
    local implementation=$2  # "perl" or "go"
    
    echo -e "${YELLOW}Running round-trip test for $implementation: $(basename "$test_file")${NC}"
    
    # Extract message ID from test data
    local message_id=""
    if command -v jq &> /dev/null; then
        message_id=$(jq -r '.message_id' "$test_file")
    else
        # Fallback to grep if jq is not available
        message_id=$(grep -o '"message_id": *"[^"]*"' "$test_file" | sed 's/"message_id": *"\(.*\)"/\1/')
    fi
    
    echo "Original message ID: $message_id"
    
    # Run encryption
    if [ "$implementation" == "perl" ]; then
        perl "$BASE_DIR/perl_test_harness.pl" --type="message_id_encryption" --input="$test_file" --format="hex" > "$RESULTS_DIR/${implementation}_encrypt.txt" 2> "$RESULTS_DIR/${implementation}_encrypt_error.txt"
    else
        go run "$BASE_DIR/go_test_harness.go" --type="message_id_encryption" --input="$test_file" --format="hex" > "$RESULTS_DIR/${implementation}_encrypt.txt" 2> "$RESULTS_DIR/${implementation}_encrypt_error.txt"
    fi
    
    local encrypt_exit_code=$?
    
    # Check for errors
    if [ $encrypt_exit_code -ne 0 ]; then
        echo -e "${RED}${implementation^} encryption failed with exit code $encrypt_exit_code${NC}"
        cat "$RESULTS_DIR/${implementation}_encrypt_error.txt"
        return 1
    fi
    
    # Extract encrypted message ID
    local encrypted_id=$(tail -n 1 "$RESULTS_DIR/${implementation}_encrypt.txt")
    echo "Encrypted message ID: $encrypted_id"
    
    # Extract crypt keys
    local key1=$(grep -o '"crypt_key_1": *"[^"]*"' "$test_file" | sed 's/"crypt_key_1": *"\(.*\)"/\1/' | head -1)
    local key2=$(grep -o '"crypt_key_2": *"[^"]*"' "$test_file" | sed 's/"crypt_key_2": *"\(.*\)"/\1/' | head -1)
    local key3=$(grep -o '"crypt_key_3": *"[^"]*"' "$test_file" | sed 's/"crypt_key_3": *"\(.*\)"/\1/' | head -1)
    
    # Create a temporary test file for decryption
    local temp_file="$RESULTS_DIR/temp_decrypt.json"
    cat > "$temp_file" << EOF
{
    "message_id": "$encrypted_id",
    "crypt_key_1": "$key1",
    "crypt_key_2": "$key2",
    "crypt_key_3": "$key3",
    "decrypt": true
}
EOF
    
    # Debug: Print the temporary file
    echo "Temporary file content:"
    cat "$temp_file"
    
    # Run decryption
    if [ "$implementation" == "perl" ]; then
        perl "$BASE_DIR/perl_test_harness.pl" --type="message_id_encryption" --input="$temp_file" --format="hex" > "$RESULTS_DIR/${implementation}_decrypt.txt" 2> "$RESULTS_DIR/${implementation}_decrypt_error.txt"
    else
        go run "$BASE_DIR/go_test_harness.go" --type="message_id_encryption" --input="$temp_file" --format="hex" > "$RESULTS_DIR/${implementation}_decrypt.txt" 2> "$RESULTS_DIR/${implementation}_decrypt_error.txt"
    fi
    
    local decrypt_exit_code=$?
    
    # Check for errors
    if [ $decrypt_exit_code -ne 0 ]; then
        echo -e "${RED}${implementation^} decryption failed with exit code $decrypt_exit_code${NC}"
        cat "$RESULTS_DIR/${implementation}_decrypt_error.txt"
        return 1
    fi
    
    # Extract decrypted message ID
    local decrypted_id=$(tail -n 1 "$RESULTS_DIR/${implementation}_decrypt.txt")
    echo "Decrypted message ID: $decrypted_id"
    
    # Compare original and decrypted message IDs (case-insensitive)
    if [ "${message_id,,}" == "${decrypted_id,,}" ]; then
        echo -e "${GREEN}Round-trip test passed: Original and decrypted message IDs match${NC}"
        return 0
    else
        echo -e "${RED}Round-trip test failed: Original and decrypted message IDs do not match${NC}"
        echo "Original: $message_id"
        echo "Decrypted: $decrypted_id"
        return 1
    fi
}

# Function to run round-trip tests for all message ID encryption test files
test_all_message_id_round_trips() {
    local implementation=$1  # "perl" or "go"
    local dir="$TEST_DATA_DIR/message_id_encryption"
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running all message ID encryption round-trip tests for $implementation implementation...${NC}"
    
    if [ ! -d "$dir" ]; then
        echo -e "${RED}Test directory not found: $dir${NC}"
        return 1
    fi
    
    for test_file in "$dir"/*.json; do
        if [ -f "$test_file" ]; then
            ((total++))
            if test_message_id_round_trip "$test_file" "$implementation"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Round-trip test summary for message ID encryption ($implementation):${NC}"
    echo -e "Total: $total, Passed: ${GREEN}$passed${NC}, Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Main function
main() {
    local command=$1
    shift
    
    case $command in
        test)
            local test_file=$1
            local implementation=$2
            
            if [ -z "$test_file" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 test <test_file> <implementation>"
                exit 1
            fi
            
            test_message_id_round_trip "$test_file" "$implementation"
            ;;
        test-all)
            local implementation=$1
            
            if [ -z "$implementation" ]; then
                echo "Usage: $0 test-all <implementation>"
                exit 1
            fi
            
            test_all_message_id_round_trips "$implementation"
            ;;
        *)
            echo "Usage: $0 <command> [args...]"
            echo "Commands:"
            echo "  test <test_file> <implementation>"
            echo "  test-all <implementation>"
            exit 1
            ;;
    esac
}

# Run the main function
main "$@"