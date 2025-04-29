# OpenKore Network Verification Test Harness

This directory contains a comprehensive test harness for verifying that the Go implementation of OpenKore's network functionality matches the Perl implementation.

## Overview

The test harness consists of:

1. **Perl Test Harness**: Tests the Perl implementation of OpenKore's network functionality
2. **Go Test Harness**: Tests the Go implementation of OpenKore's network functionality
3. **Test Runner**: Runs tests and compares the results between Perl and Go implementations
4. **Test Data**: JSON files containing test inputs and expected outputs

## Test Types

The test harness supports the following test types:

- **packet_construction**: Tests packet construction
- **packet_parsing**: Tests packet parsing
- **message_id_encryption**: Tests message ID encryption
- **padded_packets**: Tests padded packets
- **pin_encode**: Tests PIN encoding
- **network_stack**: Tests network stack instantiation
- **server_connection**: Tests server connection
- **actor_handling**: Tests actor handling
- **field_handling**: Tests field handling
- **event_hooks**: Tests event hooks
- **server_config**: Tests server configuration
- **connection_management**: Tests connection management
- **receive_function**: Tests individual functions in Network::Receive

## Network::Receive Function Testing

The test harness has been expanded to test all functions in Network::Receive.pm. This ensures that the Go implementation matches the Perl implementation for all network receive functionality.

### Test Data Generation

To generate test data for all Network::Receive functions:

```bash
./generate_receive_tests.sh
```

This will create test data files for each function in Network::Receive.pm in the `test_data/receive_functions` directory.

### Running Tests

To run tests for all Network::Receive functions:

```bash
./run_receive_tests.sh test-all
```

This will:
1. Run tests for all functions in the Perl implementation
2. Run tests for all functions in the Go implementation
3. Compare the results between the Perl and Go implementations

### Running Tests for a Specific Function

To run tests for a specific function:

```bash
./run_receive_tests.sh run <function_name> <implementation>
```

For example:

```bash
./run_receive_tests.sh run exp perl
./run_receive_tests.sh run exp go
```

### Comparing Results

To compare the results between Perl and Go implementations for a specific function:

```bash
./run_receive_tests.sh compare <function_name>
```

For example:

```bash
./run_receive_tests.sh compare exp
```

## Adding New Tests

To add a new test:

1. Create a JSON file in the appropriate test_data directory
2. Add the necessary test data and expected output
3. Run the test using the test runner

## Test Data Structure

Test data files are JSON files with the following structure:

```json
{
    "function_name": "function_name",
    "args": {
        "param1": "value1",
        "param2": "value2"
    },
    "expected_output": "expected result",
    "validation": {
        "type": "exact_match",
        "description": "Test description",
        "requirements": [
            "Requirement 1",
            "Requirement 2"
        ]
    }
}
```

## Validation Types

The test harness supports the following validation types:

- **exact_match**: Exact comparison of output with expected result
- **round_trip**: Test bidirectional operations (e.g., encrypt/decrypt)
- **protocol_compliance**: Verify compliance with protocol specifications
- **boundary_test**: Test behavior at boundary conditions
- **performance_test**: Test performance characteristics
- **security_test**: Test security properties
- **concurrency_test**: Test behavior under concurrent operations
- **integration_test**: Test interactions between components
- **error_handling**: Test response to error conditions
- **compatibility_test**: Test compatibility with different versions

## Scripts

- **test_runner.sh**: Main script for running tests
- **verify_expected_output.sh**: Script for verifying test outputs against expected outputs
- **round_trip_test.sh**: Script for testing bidirectional operations
- **generate_receive_tests.sh**: Script for generating test data for Network::Receive functions
- **run_receive_tests.sh**: Script for running tests for Network::Receive functions

## Files

- **perl_test_harness.pl**: Perl test harness
- **perl_test_harness_extensions.pl**: Extensions to the Perl test harness
- **go_test_harness.go**: Go test harness
- **go_test_harness_extensions.go**: Extensions to the Go test harness
- **receive_function_test_harness.pl**: Test harness for Network::Receive functions in Perl
- **receive_function_test_harness.go**: Test harness for Network::Receive functions in Go
- **server_connection.go**: Go implementation of server connection functionality