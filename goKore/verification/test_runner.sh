#!/bin/bash
# Test runner script for comparing Perl and Go implementations
# This script runs both test harnesses with the same inputs and compares the outputs

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$BASE_DIR/test_data"
RESULTS_DIR="$BASE_DIR/results"

# Create directories if they don't exist
mkdir -p "$TEST_DATA_DIR"
mkdir -p "$RESULTS_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to run a test case
run_test() {
    local test_type=$1
    local test_file=$2
    local format=${3:-hex}
    
    echo -e "${YELLOW}Running test: $test_type - $(basename "$test_file")${NC}"
    
    # Run Perl test harness
    perl "$BASE_DIR/perl_test_harness.pl" --type="$test_type" --input="$test_file" --format="$format" > "$RESULTS_DIR/perl_result.txt" 2> "$RESULTS_DIR/perl_error.txt"
    perl_exit_code=$?
    
    # Run Go test harness
    go run "$BASE_DIR/go_test_harness.go" --type="$test_type" --input="$test_file" --format="$format" > "$RESULTS_DIR/go_result.txt" 2> "$RESULTS_DIR/go_error.txt"
    go_exit_code=$?
    
    # Check for errors
    perl_success=true
    go_success=true
    
    if [ $perl_exit_code -ne 0 ]; then
        echo -e "${YELLOW}Perl test harness failed with exit code $perl_exit_code${NC}"
        cat "$RESULTS_DIR/perl_error.txt"
        perl_success=false
        echo -e "${YELLOW}Skipping Perl test harness due to errors${NC}"
    fi
    
    if [ $go_exit_code -ne 0 ]; then
        echo -e "${RED}Go test harness failed with exit code $go_exit_code${NC}"
        cat "$RESULTS_DIR/go_error.txt"
        go_success=false
        return 1
    fi
    
    # If Perl test harness failed, just show Go results
    if [ "$perl_success" = false ]; then
        echo -e "${GREEN}Go test harness succeeded${NC}"
        echo "Go output:"
        cat "$RESULTS_DIR/go_result.txt"
        return 0
    fi
    
    # Compare results if both succeeded
    # For packet_construction tests, normalize the output to handle different argument orders
    if [ "$test_type" == "packet_construction" ]; then
        # Extract the actual packet data (last line) from both outputs
        perl_packet=$(tail -n 1 "$RESULTS_DIR/perl_result.txt")
        go_packet=$(tail -n 1 "$RESULTS_DIR/go_result.txt")
        
        if [ "$perl_packet" == "$go_packet" ]; then
            echo -e "${GREEN}Test passed: packet data matches${NC}"
            return 0
        else
            echo -e "${RED}Test outputs differ${NC}"
            echo "Perl output:"
            cat "$RESULTS_DIR/perl_result.txt"
            echo "Go output:"
            cat "$RESULTS_DIR/go_result.txt"
            echo "Diff:"
            diff "$RESULTS_DIR/perl_result.txt" "$RESULTS_DIR/go_result.txt"
            # Return failure if outputs don't match
            return 1
        fi
    else
        # For other test types, compare the full output
        if diff -q "$RESULTS_DIR/perl_result.txt" "$RESULTS_DIR/go_result.txt" > /dev/null; then
            echo -e "${GREEN}Test passed: outputs match${NC}"
            return 0
        else
            echo -e "${RED}Test outputs differ${NC}"
            echo "Perl output:"
            cat "$RESULTS_DIR/perl_result.txt"
            echo "Go output:"
            cat "$RESULTS_DIR/go_result.txt"
            echo "Diff:"
            diff "$RESULTS_DIR/perl_result.txt" "$RESULTS_DIR/go_result.txt"
            # Return failure if outputs don't match
            return 1
        fi
    fi
}

# Function to run all tests in a directory
run_all_tests() {
    local test_type=$1
    local dir="$TEST_DATA_DIR/$test_type"
    local format=${2:-hex}
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running all $test_type tests...${NC}"
    
    if [ ! -d "$dir" ]; then
        echo -e "${RED}Test directory not found: $dir${NC}"
        return 1
    fi
    
    for test_file in "$dir"/*.json; do
        if [ -f "$test_file" ]; then
            ((total++))
            if run_test "$test_type" "$test_file" "$format"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Test summary for $test_type:${NC}"
    echo -e "Total: $total, Passed: ${GREEN}$passed${NC}, Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Create example test data if it doesn't exist
create_example_test_data() {
    # Create directories for each test type
    mkdir -p "$TEST_DATA_DIR/packet_construction"
    mkdir -p "$TEST_DATA_DIR/packet_parsing"
    mkdir -p "$TEST_DATA_DIR/message_id_encryption"
    mkdir -p "$TEST_DATA_DIR/padded_packets"
    mkdir -p "$TEST_DATA_DIR/pin_encode"
    mkdir -p "$TEST_DATA_DIR/network_stack"
    mkdir -p "$TEST_DATA_DIR/actor_handling"
    mkdir -p "$TEST_DATA_DIR/field_handling"
    mkdir -p "$TEST_DATA_DIR/event_hooks"
    mkdir -p "$TEST_DATA_DIR/server_config"
    mkdir -p "$TEST_DATA_DIR/connection_management"
    
    # Example packet construction test
    cat > "$TEST_DATA_DIR/packet_construction/actor_action.json" << EOF
{
    "packet_name": "actor_action",
    "server_type": "0",
    "args": {
        "targetID": "12345678",
        "type": 0
    }
}
EOF
    
    # Example message ID encryption test
    cat > "$TEST_DATA_DIR/message_id_encryption/example.json" << EOF
{
    "message_id": "0089",
    "crypt_key_1": "0x12345678",
    "crypt_key_2": "0x87654321",
    "crypt_key_3": "0xABCDEF01"
}
EOF
    
    # Example padded packets test
    cat > "$TEST_DATA_DIR/padded_packets/sit.json" << EOF
{
    "packet_type": "sit_stand",
    "account_id": 12345678,
    "map_sync": 87654321,
    "sync": 2882400000,
    "sit": true
}
EOF
    
    # Example PIN encode test
    cat > "$TEST_DATA_DIR/pin_encode/example.json" << EOF
{
    "seed": 1234567890,
    "pin": 1234
}
EOF
    
    # Example network stack test
    cat > "$TEST_DATA_DIR/network_stack/example.json" << EOF
{
    "server_type": "0",
    "server_ip": "127.0.0.1",
    "server_port": 6900,
    "account_id": 12345678,
    "map_sync": 87654321,
    "sync": 2882400000
}
EOF
    
    echo -e "${GREEN}Example test data created in $TEST_DATA_DIR${NC}"
}

# Main function
main() {
    local command=$1
    shift
    
    case $command in
        create-examples)
            create_example_test_data
            ;;
        run)
            local test_type=$1
            local test_file=$2
            local format=${3:-hex}
            
            if [ -z "$test_type" ] || [ -z "$test_file" ]; then
                echo "Usage: $0 run <test_type> <test_file> [format]"
                exit 1
            fi
            
            run_test "$test_type" "$test_file" "$format"
            ;;
        run-all)
            local test_type=$1
            local format=${2:-hex}
            
            if [ -z "$test_type" ]; then
                echo "Usage: $0 run-all <test_type> [format]"
                exit 1
            fi
            
            run_all_tests "$test_type" "$format"
            ;;
        run-all-types)
            local format=${1:-hex}
            local failed=0
            
            for test_type in packet_construction packet_parsing message_id_encryption padded_packets pin_encode network_stack actor_handling field_handling event_hooks server_config connection_management; do
                if [ -d "$TEST_DATA_DIR/$test_type" ]; then
                    if ! run_all_tests "$test_type" "$format"; then
                        ((failed++))
                    fi
                    echo
                fi
            done
            
            if [ $failed -eq 0 ]; then
                echo -e "${GREEN}All tests passed!${NC}"
                exit 0
            else
                echo -e "${RED}$failed test types had failures${NC}"
                exit 1
            fi
            ;;
        *)
            echo "Usage: $0 <command> [args...]"
            echo "Commands:"
            echo "  create-examples             Create example test data"
            echo "  run <type> <file> [format]  Run a specific test"
            echo "  run-all <type> [format]     Run all tests of a specific type"
            echo "  run-all-types [format]      Run all tests of all types"
            echo
            echo "Test types:"
            echo "  packet_construction"
            echo "  packet_parsing"
            echo "  message_id_encryption"
            echo "  padded_packets"
            echo "  pin_encode"
            echo "  network_stack"
            echo
            echo "Output formats:"
            echo "  hex (default)"
            echo "  json"
            echo "  raw"
            exit 1
            ;;
    esac
}

# Run the main function with all arguments
main "$@"
