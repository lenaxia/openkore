#!/bin/bash
# Script to run advanced test scenarios

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$BASE_DIR/test_data"
RESULTS_DIR="$BASE_DIR/results"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create results directory if it doesn't exist
mkdir -p "$RESULTS_DIR"

# Function to run a specific test scenario
run_test_scenario() {
    local test_type=$1
    local test_file=$2
    local implementation=$3  # "perl" or "go"
    
    echo -e "${BLUE}Running $test_type test: $(basename "$test_file") on $implementation implementation${NC}"
    
    # First verify against expected output
    "$BASE_DIR/verify_expected_output.sh" verify "$test_type" "$test_file" "$implementation"
    local verify_result=$?
    
    # For message_id_encryption, also run round-trip test
    if [ "$test_type" == "message_id_encryption" ]; then
        echo -e "${BLUE}Running round-trip test for: $(basename "$test_file")${NC}"
        "$BASE_DIR/round_trip_test.sh" test "$test_file" "$implementation"
        local round_trip_result=$?
        
        # Both tests must pass
        if [ $verify_result -eq 0 ] && [ $round_trip_result -eq 0 ]; then
            return 0
        else
            return 1
        fi
    else
        return $verify_result
    fi
}

# Function to run all advanced test scenarios for a specific type
run_advanced_tests() {
    local test_type=$1
    local implementation=$2  # "perl" or "go"
    local dir="$TEST_DATA_DIR/$test_type"
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running advanced tests for $test_type on $implementation implementation...${NC}"
    
    if [ ! -d "$dir" ]; then
        echo -e "${RED}Test directory not found: $dir${NC}"
        return 1
    fi
    
    # Find test files with advanced validation criteria
    for test_file in "$dir"/*.json; do
        if [ -f "$test_file" ]; then
            # Check if the file contains validation criteria
            if grep -q '"validation"' "$test_file"; then
                ((total++))
                if run_test_scenario "$test_type" "$test_file" "$implementation"; then
                    echo -e "${GREEN}Test passed: $(basename "$test_file")${NC}"
                    ((passed++))
                else
                    echo -e "${RED}Test failed: $(basename "$test_file")${NC}"
                    ((failed++))
                fi
                echo
            fi
        fi
    done
    
    echo -e "${YELLOW}Advanced test summary for $test_type ($implementation):${NC}"
    echo -e "Total: $total, Passed: ${GREEN}$passed${NC}, Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to run all advanced tests for all types
run_all_advanced_tests() {
    local implementation=$1  # "perl" or "go"
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running all advanced tests for $implementation implementation...${NC}"
    
    # Run tests for each test type
    for test_type in $(ls "$TEST_DATA_DIR"); do
        if [ -d "$TEST_DATA_DIR/$test_type" ]; then
            ((total++))
            if run_advanced_tests "$test_type" "$implementation"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Overall advanced test summary ($implementation):${NC}"
    echo -e "Total test types: $total, Passed: ${GREEN}$passed${NC}, Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}All advanced tests passed!${NC}"
        return 0
    else
        echo -e "${RED}$failed test types had failures${NC}"
        return 1
    fi
}

# Function to compare implementations
compare_implementations() {
    echo -e "${YELLOW}Comparing Perl and Go implementations...${NC}"
    
    # Run advanced tests for both implementations
    echo -e "${BLUE}Running advanced tests for Perl implementation...${NC}"
    run_all_advanced_tests "perl"
    local perl_result=$?
    
    echo -e "${BLUE}Running advanced tests for Go implementation...${NC}"
    run_all_advanced_tests "go"
    local go_result=$?
    
    # Compare results
    if [ $perl_result -eq 0 ] && [ $go_result -eq 0 ]; then
        echo -e "${GREEN}Both implementations passed all advanced tests!${NC}"
        return 0
    elif [ $perl_result -eq 0 ]; then
        echo -e "${RED}Perl implementation passed but Go implementation failed${NC}"
        return 1
    elif [ $go_result -eq 0 ]; then
        echo -e "${RED}Go implementation passed but Perl implementation failed${NC}"
        return 1
    else
        echo -e "${RED}Both implementations had failures${NC}"
        return 1
    fi
}

# Main function
main() {
    local command=$1
    shift
    
    case $command in
        run)
            local test_type=$1
            local test_file=$2
            local implementation=$3
            
            if [ -z "$test_type" ] || [ -z "$test_file" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 run <test_type> <test_file> <implementation>"
                exit 1
            fi
            
            run_test_scenario "$test_type" "$test_file" "$implementation"
            ;;
        run-type)
            local test_type=$1
            local implementation=$2
            
            if [ -z "$test_type" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 run-type <test_type> <implementation>"
                exit 1
            fi
            
            run_advanced_tests "$test_type" "$implementation"
            ;;
        run-all)
            local implementation=$1
            
            if [ -z "$implementation" ]; then
                echo "Usage: $0 run-all <implementation>"
                exit 1
            fi
            
            run_all_advanced_tests "$implementation"
            ;;
        compare)
            compare_implementations
            ;;
        *)
            echo "Usage: $0 <command> [args...]"
            echo "Commands:"
            echo "  run <test_type> <test_file> <implementation>"
            echo "  run-type <test_type> <implementation>"
            echo "  run-all <implementation>"
            echo "  compare"
            exit 1
            ;;
    esac
}

# Run the main function
main "$@"