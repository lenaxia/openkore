# Send.pm Wrapper for Validation Data Generation

This directory contains tools for generating validation data from the original OpenKore Send.pm module. These tools help create test data for validating Go reimplementations of the network protocol.

## Overview

The tools in this directory analyze the Send.pm module to determine what arguments each method expects, then generate appropriate test data by calling each method and capturing the output packets. The generated data can be used to validate that your Go reimplementation produces the same packets as the original Perl code.

## Scripts

### analyze_send_methods.pl

This script analyzes the Send.pm module to determine:

1. What parameters each method expects
2. What types those parameters are likely to be (based on naming conventions)
3. Which methods are "complex" and need special handling
4. For methods that take a generic `$args` parameter, what keys are accessed in that hash

The script saves its analysis to two JSON files:
- `analysis/send_methods_analysis.json`: Contains the full analysis of all methods
- `analysis/complex_methods.json`: Contains only the complex methods that need special handling

#### Usage

```bash
./analyze_send_methods.pl
```

### generate_from_analysis.pl

This script uses the analysis from `analyze_send_methods.pl` to generate test data:

1. It loads the analysis results
2. For each method that isn't marked as "complex":
   - It generates appropriate arguments based on parameter types
   - It calls the method with those arguments
   - It captures the output packets
   - It saves the input arguments and output packets to a JSON file

Complex methods are skipped and listed in a separate file for later processing by an LLM or manual intervention.

#### Usage

```bash
./generate_from_analysis.pl
```

## Output Format

The generated test data is saved in JSON format with the following structure:

```json
{
  "method": "sendMove",
  "args": [100, 100],
  "packets": [
    {
      "hex": "85640000640000a9860100",
      "messageID": "0085",
      "bytes": [133, 100, 0, 0, 100, 0, 0, 169, 134, 1, 0]
    }
  ]
}
```

Each JSON file contains:
- The method name
- The arguments used
- An array of packets, each with:
  - Hex representation
  - Message ID
  - Array of byte values

## Parameter Type Heuristics

The scripts use the following heuristics to determine parameter types:

1. Parameters with names like `$name`, `$message`, `$title`, etc. are treated as strings
2. Parameters with names like `$ID`, `$accountID`, `$charID`, etc. are treated as numeric IDs
3. Parameters with names like `$amount`, `$flag`, `$type`, etc. are treated as numeric values
4. Parameters with names like `$x`, `$y` are treated as coordinates
5. Parameters with names like `$r_message`, `$r_array` (with `r_` prefix) are treated as references
6. Parameters named `$args` are treated as hash references

## Handling Complex Methods

Some methods are marked as "complex" and skipped by the automatic generation:

1. Methods that take reference parameters (`$r_message`, etc.)
2. Methods that use complex data structures
3. Methods that use special Perl features like dereferencing (`$$var`)

These methods are listed in `testdata/skipped_methods.txt` and can be processed separately by an LLM or manually.

## Using the Test Data

The generated test data can be used to validate your Go reimplementation:

1. Load the test data in your Go code
2. Extract the arguments
3. Call your implementation with the same arguments
4. Compare the output with the expected output from the test data

This ensures that your Go implementation produces the same packets as the original Perl code.

## Example Go Validation Code

```go
func TestSendMove(t *testing.T) {
    // Load test data
    data, err := loadTestData("testdata/sendMove.json")
    if err != nil {
        t.Fatal(err)
    }
    
    // Extract arguments
    x := data.Args[0].(float64)
    y := data.Args[1].(float64)
    
    // Call your implementation
    packet := yourImplementation.SendMove(int(x), int(y))
    
    // Compare with expected output
    expectedHex := data.Packets[0].Hex
    actualHex := hex.EncodeToString(packet)
    
    if expectedHex != actualHex {
        t.Errorf("Expected %s, got %s", expectedHex, actualHex)
    }
}
```

## Troubleshooting

If you encounter errors:

1. Check that the paths to the OpenKore source files are correct
2. Ensure all required dependencies are mocked properly
3. For methods that require specific state or configuration, you may need to add additional mock data to the script
4. If a method fails to generate packets, check the method implementation to see if it requires specific conditions or state