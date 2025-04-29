# Verification Plan for Go Implementation

This document outlines the plan for verifying that the Go implementation of the OpenKore network stack is functionally identical to the original Perl implementation.

## Approach

To ensure functional equivalence between the Go and Perl implementations, we will:

1. Create side-by-side tests that run both implementations with the same inputs
2. Compare the outputs to verify they are identical
3. Focus on key functionality that must behave identically

## Key Areas for Verification

### 1. Packet Construction

- **Test Case**: Construct the same packet in both implementations
- **Verification Method**: Compare the binary output byte-by-byte
- **Perl Files**: `src/Network/Send.pm`, `src/Network/Send/ServerType0.pm`
- **Go Files**: `network/packets/constructor.go`
- **Test Data**: Create test data for common packet types:

#### 1.1 Login Packets
- **login_packet**: Send account name and password to login server
  - Parameters: username, password, version, clientHash
  - Expected: Properly formatted login packet with encrypted password

#### 1.2 Character Selection Packets
- **char_select**: Select a character to play
  - Parameters: slot number
  - Expected: Properly formatted character selection packet

#### 1.3 Movement Packets
- **move_to**: Move to coordinates
  - Parameters: x, y coordinates
  - Expected: Properly formatted movement packet with coordinates
- **sit_stand**: Sit or stand
  - Parameters: sit flag (0/1)
  - Expected: Properly formatted sit/stand packet

#### 1.4 Chat Packets
- **chat_message**: Send chat message
  - Parameters: message text
  - Expected: Properly formatted chat packet with message
- **private_message**: Send private message
  - Parameters: target name, message text
  - Expected: Properly formatted private message packet

#### 1.5 Item Packets
- **use_item**: Use an item from inventory
  - Parameters: inventory index
  - Expected: Properly formatted use item packet
- **drop_item**: Drop an item from inventory
  - Parameters: inventory index, amount
  - Expected: Properly formatted drop item packet

#### 1.6 Skill Packets
- **use_skill**: Use a skill
  - Parameters: skill ID, skill level, target ID
  - Expected: Properly formatted skill use packet
- **use_skill_location**: Use a ground-targeted skill
  - Parameters: skill ID, skill level, x, y coordinates
  - Expected: Properly formatted ground skill packet

### 2. Packet Parsing

- **Test Case**: Parse the same packet in both implementations
- **Verification Method**: Compare the extracted fields and values
- **Perl Files**: `src/Network/Receive.pm`, `src/Network/Receive/ServerType0.pm`
- **Go Files**: `network/protocol/parser.go`
- **Test Data**: Use packet captures for the following packet types:

#### 2.1 Server Status Packets
- **server_list**: List of available servers
  - Expected fields: server names, IPs, ports, user counts
- **login_result**: Result of login attempt
  - Expected fields: result code, error message (if any)

#### 2.2 Character Packets
- **character_list**: List of characters on account
  - Expected fields: character names, levels, classes, etc.
- **character_stats**: Character statistics
  - Expected fields: HP, SP, base stats, etc.

#### 2.3 Map Packets
- **map_change**: Map change notification
  - Expected fields: map name, coordinates
- **map_loaded**: Map loaded notification
  - Expected fields: confirmation code

#### 2.4 Entity Packets
- **actor_info**: Information about an actor (player, monster, NPC)
  - Expected fields: entity ID, position, name, etc.
- **actor_move**: Actor movement
  - Expected fields: entity ID, start coordinates, end coordinates

#### 2.5 Chat Packets
- **public_chat**: Public chat message
  - Expected fields: sender name, message
- **private_chat**: Private chat message
  - Expected fields: sender name, message

### 3. Message ID Encryption

- **Test Case**: Encrypt the same message ID with the same keys in both implementations
- **Verification Method**: Compare the encrypted message IDs
- **Perl Files**: `src/Network/Send.pm` (encryptMessageID method)
- **Go Files**: `network/packets/constructor.go` (EncryptMessageID method)
- **Test Data**: Various message IDs with different encryption keys:

#### 3.1 Basic Encryption Tests
- **Standard IDs**: Common packet IDs (0x007D, 0x0089, etc.)
  - With standard keys (0x12345678, 0x87654321, 0xABCDEF01)
  - Expected: Correctly encrypted IDs

#### 3.2 Edge Cases
- **Zero ID**: Message ID 0x0000
  - Expected: Correctly encrypted ID
- **Max ID**: Message ID 0xFFFF
  - Expected: Correctly encrypted ID

#### 3.3 Different Key Combinations
- **All Zero Keys**: Keys (0x00000000, 0x00000000, 0x00000000)
  - Expected: Correctly encrypted ID
- **All Max Keys**: Keys (0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF)
  - Expected: Correctly encrypted ID
- **Mixed Keys**: Various combinations of keys
  - Expected: Correctly encrypted ID

### 4. Padded Packets

- **Test Case**: Generate padded packets in both implementations
- **Verification Method**: Compare the generated packets
- **Perl Files**: `src/Network/PacketParser.pm`, `src/Network/PaddedPackets.pm`
- **Go Files**: `network/protocol/padding.go`
- **Test Data**: Different packet types that require padding:

#### 4.1 Sit/Stand Packets
- **Sit**: Generate sit packet with padding
  - Parameters: account ID, map sync, sync, sit=true
  - Expected: Correctly padded sit packet
- **Stand**: Generate stand packet with padding
  - Parameters: account ID, map sync, sync, sit=false
  - Expected: Correctly padded stand packet

#### 4.2 Attack Packets
- **Normal Attack**: Generate attack packet with padding
  - Parameters: account ID, map sync, sync, target ID, flag=0
  - Expected: Correctly padded attack packet
- **Continuous Attack**: Generate continuous attack packet with padding
  - Parameters: account ID, map sync, sync, target ID, flag=7
  - Expected: Correctly padded continuous attack packet

#### 4.3 Skill Use Packets
- **Target Skill**: Generate skill use packet with padding
  - Parameters: account ID, map sync, sync, skill ID, skill level, target ID
  - Expected: Correctly padded skill use packet
- **Ground Skill**: Generate ground skill packet with padding
  - Parameters: account ID, map sync, sync, skill ID, skill level, x, y
  - Expected: Correctly padded ground skill packet

#### 4.4 Hash Generation
- **Hash Algorithm**: Test hash generation with different inputs
  - Parameters: various account IDs, map syncs, syncs
  - Expected: Correctly generated hash values

### 5. PIN Encoding

- **Test Case**: Encode the same PIN with the same seed in both implementations
- **Verification Method**: Compare the encoded PINs
- **Perl Files**: `src/Network/Send.pm` (pinEncode method)
- **Go Files**: `network/packets/constructor.go` (PinEncode method)
- **Test Data**: Various PINs and seeds:

#### 5.1 Standard PINs
- **4-Digit PIN**: Common 4-digit PINs (1234, 9999, etc.)
  - With various seeds
  - Expected: Correctly encoded PIN

#### 5.2 Edge Cases
- **All Zeros**: PIN 0000
  - Expected: Correctly encoded PIN
- **All Nines**: PIN 9999
  - Expected: Correctly encoded PIN
- **Sequential**: PIN 1234
  - Expected: Correctly encoded PIN
- **Repeated**: PIN 1111
  - Expected: Correctly encoded PIN

#### 5.3 Different Seeds
- **Zero Seed**: Seed 0
  - Expected: Correctly encoded PIN
- **Large Seed**: Seed 2147483647 (max int32)
  - Expected: Correctly encoded PIN
- **Negative Seed**: Seed -1234567890
  - Expected: Correctly encoded PIN

### 6. Tokenizer

- **Test Case**: Process the same packet stream in both implementations
- **Verification Method**: Compare the extracted packets
- **Perl Files**: `src/Network/MessageTokenizer.pm`
- **Go Files**: `network/protocol/tokenizer.go`
- **Test Data**: Various packet streams:

#### 6.1 Simple Packets
- **Fixed-Length Packets**: Stream of fixed-length packets
  - Expected: Correctly extracted packets
- **Variable-Length Packets**: Stream of variable-length packets
  - Expected: Correctly extracted packets

#### 6.2 Mixed Packets
- **Mixed Packet Types**: Stream of different packet types
  - Expected: Correctly extracted packets
- **Fragmented Packets**: Stream with packet fragments
  - Expected: Correctly handled fragments and extracted packets

#### 6.3 Edge Cases
- **Empty Stream**: Empty packet stream
  - Expected: No packets extracted
- **Incomplete Packets**: Stream ending with incomplete packet
  - Expected: Correctly extracted complete packets, buffer incomplete packet

## Implementation Plan

### 1. Create Test Harness for Perl Implementation

Create a Perl script that:
- Takes input parameters (test case type, input data)
- Calls the appropriate OpenKore functions
- Outputs the results in a standardized format (e.g., hex dump for binary data)

```perl
# Example for packet construction
use Network::Send;
use Network::Send::ServerType0;

my $args = { ... }; # Test input
my $send = Network::Send::ServerType0->new();
my $packet = $send->construct_packet("packet_name", $args);
print unpack("H*", $packet); # Output as hex string
```

### 2. Create Test Harness for Go Implementation

Create a Go program that:
- Takes the same input parameters
- Calls the equivalent Go functions
- Outputs the results in the same standardized format

```go
// Example for packet construction
package main

import (
    "encoding/hex"
    "fmt"
    "github.com/lenaxia/goKore/network/packets"
)

func main() {
    db := packets.NewDefaultPacketDatabase()
    constructor := packets.NewPacketConstructor(db)
    
    args := map[string]interface{}{...} // Same test input
    packet, _ := constructor.ConstructPacket("packet_name", args)
    fmt.Println(hex.EncodeToString(packet)) // Output as hex string
}
```

### 3. Create Test Runner

Create a script that:
- Runs both test harnesses with the same inputs
- Compares the outputs
- Reports any differences

```bash
#!/bin/bash

# Run test case
perl_output=$(perl perl_test_harness.pl "$test_case" "$input_data")
go_output=$(go run go_test_harness.go "$test_case" "$input_data")

# Compare outputs
if [ "$perl_output" == "$go_output" ]; then
    echo "Test passed: outputs match"
else
    echo "Test failed: outputs differ"
    echo "Perl output: $perl_output"
    echo "Go output: $go_output"
fi
```

### 4. Create Test Data

Create test data files for each test case, including:
- Input parameters
- Expected output (optional, for additional verification)

### 5. Run Tests and Document Results

Run the tests for each key area and document the results:
- Which tests pass (outputs match)
- Which tests fail (outputs differ)
- For failed tests, analyze the differences and fix the Go implementation

## Test Case Examples

### Packet Construction Test

**Input:**
```json
{
  "packet_name": "actor_action",
  "args": {
    "targetID": "12345678",
    "type": 0
  }
}
```

**Expected Perl Output:**
```
89123456780000
```

**Expected Go Output:**
```
89123456780000
```

### Message ID Encryption Test

**Input:**
```json
{
  "message_id": "0089",
  "crypt_key_1": "0x12345678",
  "crypt_key_2": "0x87654321",
  "crypt_key_3": "0xABCDEF01"
}
```

**Expected Perl Output:**
```
7a89
```

**Expected Go Output:**
```
7a89
```

## Timeline

1. Create test harnesses for Perl and Go implementations (1 week)
2. Create test data for each key area (1 week)
3. Run tests and document results (1 week)
4. Fix any discrepancies in the Go implementation (2 weeks)
5. Re-run tests to verify fixes (1 week)

## Success Criteria

The Go implementation will be considered functionally identical to the Perl implementation when:

1. All test cases produce identical outputs in both implementations
2. Any differences are documented and justified (e.g., intentional improvements or bug fixes)
3. The Go implementation passes all existing tests in the OpenKore test suite

## Conclusion

This verification plan provides a systematic approach to ensuring that the Go implementation of the OpenKore network stack is functionally identical to the original Perl implementation. By focusing on key areas and using side-by-side testing, we can identify and fix any discrepancies to ensure compatibility and correctness.