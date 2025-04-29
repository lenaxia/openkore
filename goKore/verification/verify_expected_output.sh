#!/bin/bash
# Script to verify test outputs against expected outputs defined in test data files

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$BASE_DIR/test_data"
RESULTS_DIR="$BASE_DIR/results"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to extract expected output from a test data file
extract_expected_output() {
    local test_file=$1
    
    # Use jq to extract the expected_output field if it exists
    if command -v jq &> /dev/null; then
        expected_output=$(jq -r '.expected_output // "N/A"' "$test_file")
    else
        # Fallback to grep if jq is not available
        expected_output=$(grep -o '"expected_output": *"[^"]*"' "$test_file" | sed 's/"expected_output": *"\(.*\)"/\1/')
        if [ -z "$expected_output" ]; then
            expected_output="N/A"
        fi
    fi
    
    echo "$expected_output"
}

# Function to extract validation type from a test data file
extract_validation_type() {
    local test_file=$1
    
    # Use jq to extract the validation.type field if it exists
    if command -v jq &> /dev/null; then
        validation_type=$(jq -r '.validation.type // "N/A"' "$test_file")
    else
        # Fallback to grep if jq is not available
        validation_type=$(grep -o '"type": *"[^"]*"' "$test_file" | head -1 | sed 's/"type": *"\(.*\)"/\1/')
        if [ -z "$validation_type" ]; then
            validation_type="N/A"
        fi
    fi
    
    echo "$validation_type"
}

# Function to verify test output against expected output
verify_expected_output() {
    local test_type=$1
    local test_file=$2
    local result_file=$3
    
    # Extract expected output from test data
    local expected_output=$(extract_expected_output "$test_file")
    local validation_type=$(extract_validation_type "$test_file")
    
    # If no expected output is defined, skip validation
    if [ "$expected_output" == "N/A" ]; then
        echo -e "${YELLOW}No expected output defined in test data, skipping validation${NC}"
        return 0
    fi
    
    # Extract actual output (last line of result file)
    local actual_output=$(tail -n 1 "$result_file")
    
    echo -e "${YELLOW}Validating output against expected result${NC}"
    echo "Expected: $expected_output"
    echo "Actual:   $actual_output"
    echo "Validation type: $validation_type"
    
    # Compare actual output with expected output
    if [ "$actual_output" == "$expected_output" ]; then
        echo -e "${GREEN}Validation passed: Output matches expected result${NC}"
        return 0
    else
        echo -e "${RED}Validation failed: Output does not match expected result${NC}"
        return 1
    fi
}

# Function to run a test case with validation
run_test_with_validation() {
    local test_type=$1
    local test_file=$2
    local implementation=$3  # "perl" or "go"
    local format=${4:-hex}
    
    echo -e "${YELLOW}Running $implementation test: $test_type - $(basename "$test_file")${NC}"
    
    # Run the test harness
    if [ "$implementation" == "perl" ]; then
        perl "$BASE_DIR/perl_test_harness.pl" --type="$test_type" --input="$test_file" --format="$format" > "$RESULTS_DIR/${implementation}_result.txt" 2> "$RESULTS_DIR/${implementation}_error.txt"
    else
        go run "$BASE_DIR/go_test_harness.go" --type="$test_type" --input="$test_file" --format="$format" > "$RESULTS_DIR/${implementation}_result.txt" 2> "$RESULTS_DIR/${implementation}_error.txt"
    fi
    
    local exit_code=$?
    
    # Check for errors
    if [ $exit_code -ne 0 ]; then
        echo -e "${RED}${implementation^} test harness failed with exit code $exit_code${NC}"
        cat "$RESULTS_DIR/${implementation}_error.txt"
        return 1
    fi
    
    # Verify output against expected output
    verify_expected_output "$test_type" "$test_file" "$RESULTS_DIR/${implementation}_result.txt"
    return $?
}

# Function to run all tests in a directory with validation
run_all_tests_with_validation() {
    local test_type=$1
    local implementation=$2  # "perl" or "go"
    local format=${3:-hex}
    local dir="$TEST_DATA_DIR/$test_type"
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running all $test_type tests for $implementation implementation...${NC}"
    
    if [ ! -d "$dir" ]; then
        echo -e "${RED}Test directory not found: $dir${NC}"
        return 1
    fi
    
    for test_file in "$dir"/*.json; do
        if [ -f "$test_file" ]; then
            ((total++))
            if run_test_with_validation "$test_type" "$test_file" "$implementation" "$format"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Test summary for $test_type ($implementation):${NC}"
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
        verify)
            local test_type=$1
            local test_file=$2
            local implementation=$3
            local format=${4:-hex}
            
            if [ -z "$test_type" ] || [ -z "$test_file" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 verify <test_type> <test_file> <implementation> [format]"
                exit 1
            fi
            
            run_test_with_validation "$test_type" "$test_file" "$implementation" "$format"
            ;;
        verify-all)
            local test_type=$1
            local implementation=$2
            local format=${3:-hex}
            
            if [ -z "$test_type" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 verify-all <test_type> <implementation> [format]"
                exit 1
            fi
            
            run_all_tests_with_validation "$test_type" "$implementation" "$format"
            ;;
        verify-all-types)
            local implementation=$1
            local format=${2:-hex}
            local failed=0
            
            if [ -z "$implementation" ]; then
                echo "Usage: $0 verify-all-types <implementation> [format]"
                exit 1
            fi
            
            for test_type in $(ls "$TEST_DATA_DIR"); do
                if [ -d "$TEST_DATA_DIR/$test_type" ]; then
                    if ! run_all_tests_with_validation "$test_type" "$implementation" "$format"; then
                        ((failed++))
                    fi
                    echo
                fi
            done
            
            if [ $failed -eq 0 ]; then
                echo -e "${GREEN}All tests passed validation!${NC}"
                exit 0
            else
                echo -e "${RED}$failed test types had validation failures${NC}"
                exit 1
            fi
            ;;
        *)
            echo "Usage: $0 <command> [args...]"
            echo "Commands:"
            echo "  verify <test_type> <test_file> <implementation> [format]"
            echo "  verify-all <test_type> <implementation> [format]"
            echo "  verify-all-types <implementation> [format]"
            exit 1
            ;;
    esac
}

# Run the main function
main "$@"