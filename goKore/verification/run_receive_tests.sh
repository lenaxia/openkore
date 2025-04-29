#!/bin/bash
# Script to run tests for Network::Receive functions

# Set the base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DATA_DIR="$BASE_DIR/test_data/receive_functions"
RESULTS_DIR="$BASE_DIR/results/receive_functions"

# Create results directory
mkdir -p "$RESULTS_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to run a test for a specific function
run_function_test() {
    local function_name=$1
    local implementation=$2  # "perl" or "go"
    local test_file="$TEST_DATA_DIR/${function_name}.json"
    
    if [ ! -f "$test_file" ]; then
        echo -e "${YELLOW}Test file not found for function: $function_name${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Running $implementation test for function: $function_name${NC}"
    
    # Create a temporary test file with the function name
    local temp_file="$RESULTS_DIR/temp_${function_name}.json"
    cp "$test_file" "$temp_file"
    
    echo "Test file: $test_file"
    echo "Temp file: $temp_file"
    echo "Results directory: $RESULTS_DIR"
    
    # Run the test
    if [ "$implementation" == "perl" ]; then
        echo "Running Perl test harness..."
        echo "Command: perl \"$BASE_DIR/perl_test_harness.pl\" --type=\"receive_function\" --input=\"$temp_file\" --format=\"raw\""
        
        # Check if the test harness file exists
        if [ ! -f "$BASE_DIR/perl_test_harness.pl" ]; then
            echo -e "${RED}Error: perl_test_harness.pl not found${NC}"
            return 1
        fi
        
        # Check if receive_function_test_harness.pl exists
        if [ ! -f "$BASE_DIR/receive_function_test_harness.pl" ]; then
            echo -e "${RED}Error: receive_function_test_harness.pl not found${NC}"
            return 1
        fi
        
        perl "$BASE_DIR/perl_test_harness.pl" --type="receive_function" --input="$temp_file" --format="raw" > "$RESULTS_DIR/${implementation}_${function_name}.txt" 2> "$RESULTS_DIR/${implementation}_${function_name}_error.txt"
    else
        echo "Running Go test harness..."
        echo "Command: go run \"$BASE_DIR/go_test_harness.go\" --type=\"receive_function\" --input=\"$temp_file\" --format=\"raw\""
        
        # Check if the test harness file exists
        if [ ! -f "$BASE_DIR/go_test_harness.go" ]; then
            echo -e "${RED}Error: go_test_harness.go not found${NC}"
            return 1
        fi
        
        # Check if receive_function_test_harness.go exists
        if [ ! -f "$BASE_DIR/receive_function_test_harness.go" ]; then
            echo -e "${RED}Error: receive_function_test_harness.go not found${NC}"
            return 1
        fi
        
        go run "$BASE_DIR/go_test_harness.go" --type="receive_function" --input="$temp_file" --format="raw" > "$RESULTS_DIR/${implementation}_${function_name}.txt" 2> "$RESULTS_DIR/${implementation}_${function_name}_error.txt"
    fi
    
    local exit_code=$?
    
    # Check for errors
    if [ $exit_code -ne 0 ]; then
        echo -e "${RED}${implementation^} test failed with exit code $exit_code${NC}"
        cat "$RESULTS_DIR/${implementation}_${function_name}_error.txt"
        return 1
    fi
    
    echo -e "${GREEN}Test completed successfully${NC}"
    return 0
}

# Function to compare results between Perl and Go implementations
compare_results() {
    local function_name=$1
    local perl_result_file="$RESULTS_DIR/perl_${function_name}.txt"
    local go_result_file="$RESULTS_DIR/go_${function_name}.txt"
    
    if [ ! -f "$perl_result_file" ] || [ ! -f "$go_result_file" ]; then
        echo -e "${YELLOW}Result files not found for function: $function_name${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Comparing results for function: $function_name${NC}"
    
    # Compare the results
    if diff -q "$perl_result_file" "$go_result_file" > /dev/null; then
        echo -e "${GREEN}Results match: Perl and Go implementations produce the same output${NC}"
        return 0
    else
        echo -e "${RED}Results differ: Perl and Go implementations produce different outputs${NC}"
        echo "Perl output:"
        cat "$perl_result_file"
        echo "Go output:"
        cat "$go_result_file"
        echo "Diff:"
        diff "$perl_result_file" "$go_result_file"
        return 1
    fi
}

# Function to run tests for all functions
run_all_function_tests() {
    local implementation=$1  # "perl" or "go"
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Running all function tests for $implementation implementation...${NC}"
    
    # Check if test data directory exists
    if [ ! -d "$TEST_DATA_DIR" ]; then
        echo -e "${RED}Test data directory not found: $TEST_DATA_DIR${NC}"
        echo -e "${YELLOW}Run ./generate_receive_tests.sh to create test data files${NC}"
        return 1
    fi
    
    # Run tests for each function
    for test_file in "$TEST_DATA_DIR"/*.json; do
        if [ -f "$test_file" ]; then
            local function_name=$(basename "$test_file" .json)
            ((total++))
            if run_function_test "$function_name" "$implementation"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Test summary for $implementation:${NC}"
    echo -e "Total: $total, Passed: ${GREEN}$passed${NC}, Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to compare all results
compare_all_results() {
    local failed=0
    local passed=0
    local total=0
    
    echo -e "${YELLOW}Comparing all results between Perl and Go implementations...${NC}"
    
    # Compare results for each function
    for test_file in "$TEST_DATA_DIR"/*.json; do
        if [ -f "$test_file" ]; then
            local function_name=$(basename "$test_file" .json)
            ((total++))
            if compare_results "$function_name"; then
                ((passed++))
            else
                ((failed++))
            fi
            echo
        fi
    done
    
    echo -e "${YELLOW}Comparison summary:${NC}"
    echo -e "Total: $total, Matched: ${GREEN}$passed${NC}, Differed: ${RED}$failed${NC}"
    
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
        run)
            local function_name=$1
            local implementation=$2
            
            if [ -z "$function_name" ] || [ -z "$implementation" ]; then
                echo "Usage: $0 run <function_name> <implementation>"
                exit 1
            fi
            
            run_function_test "$function_name" "$implementation"
            ;;
        run-all)
            local implementation=$1
            
            if [ -z "$implementation" ]; then
                echo "Usage: $0 run-all <implementation>"
                exit 1
            fi
            
            run_all_function_tests "$implementation"
            ;;
        compare)
            local function_name=$1
            
            if [ -z "$function_name" ]; then
                echo "Usage: $0 compare <function_name>"
                exit 1
            fi
            
            compare_results "$function_name"
            ;;
        compare-all)
            compare_all_results
            ;;
        test-all)
            echo -e "${YELLOW}Running tests for both implementations and comparing results...${NC}"
            
            # Run tests for Perl implementation
            run_all_function_tests "perl"
            local perl_result=$?
            
            # Run tests for Go implementation
            run_all_function_tests "go"
            local go_result=$?
            
            # Compare results
            compare_all_results
            local compare_result=$?
            
            # Check overall result
            if [ $perl_result -eq 0 ] && [ $go_result -eq 0 ] && [ $compare_result -eq 0 ]; then
                echo -e "${GREEN}All tests passed and results match!${NC}"
                exit 0
            else
                echo -e "${RED}Some tests failed or results differ${NC}"
                exit 1
            fi
            ;;
        *)
            echo "Usage: $0 <command> [args...]"
            echo "Commands:"
            echo "  run <function_name> <implementation>  Run test for a specific function"
            echo "  run-all <implementation>             Run tests for all functions"
            echo "  compare <function_name>              Compare results for a specific function"
            echo "  compare-all                          Compare results for all functions"
            echo "  test-all                             Run tests for both implementations and compare results"
            exit 1
            ;;
    esac
}

# Run the main function
main "$@"