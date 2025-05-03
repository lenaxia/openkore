package packet_validator

import (
	"encoding/hex"
	"testing"
)

// MockSendImplementation is a simple mock implementation for testing
type MockSendImplementation struct{}

// SendMove is a mock implementation of the sendMove method
func (m *MockSendImplementation) SendMove(x, y int) []byte {
	// This is where you would implement your actual packet construction logic
	// For this example, we'll just return a hardcoded packet that matches the expected output
	packet, _ := hex.DecodeString("8500190640")
	return packet
}

// SendChat is a mock implementation of the sendChat method
func (m *MockSendImplementation) SendChat(message string) []byte {
	// This is where you would implement your actual packet construction logic
	// For this example, we'll just return a hardcoded packet
	packet, _ := hex.DecodeString("8c000f00" + hex.EncodeToString([]byte(message)) + "00")
	return packet
}

// TestSendMoveWithValidationData demonstrates how to test your implementation against validation data
func TestSendMoveWithValidationData(t *testing.T) {
	// Load validation data for sendMove
	moveData, err := LoadPacketData("../testdata/sendMove.json")
	if err != nil {
		// Skip the test if the validation data doesn't exist
		t.Skip("Validation data not found. Run generate_all_packets.pl first.")
		return
	}

	// Create your implementation
	impl := &MockSendImplementation{}

	// Extract arguments from validation data
	x, ok1 := moveData.Args[0].(float64)
	y, ok2 := moveData.Args[1].(float64)
	if !ok1 || !ok2 {
		t.Fatalf("Invalid arguments in validation data: %v", moveData.Args)
	}

	// Call your implementation with the same arguments
	result := impl.SendMove(int(x), int(y))

	// Validate the result
	validationResult := ValidatePacket("sendMove", moveData.Args, result, moveData)

	// Check if the validation passed
	if !validationResult.IsValid {
		t.Errorf("SendMove output doesn't match validation data")
		for _, err := range validationResult.Errors {
			t.Errorf("  %s", err)
		}
		t.Errorf("Expected: %s", validationResult.ExpectedHex)
		t.Errorf("Got: %s", validationResult.ActualHex)
	}
}

// TestAllPackets demonstrates how to test all your implementations against all validation data
func TestAllPackets(t *testing.T) {
	// Create a packet generator function
	generator := func(methodName string, args []interface{}) ([]byte, error) {
		impl := &MockSendImplementation{}

		switch methodName {
		case "sendMove":
			x, ok1 := args[0].(float64)
			y, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, nil
			}
			return impl.SendMove(int(x), int(y)), nil

		case "sendChat":
			message, ok := args[0].(string)
			if !ok {
				return nil, nil
			}
			return impl.SendChat(message), nil

		default:
			// Skip methods we haven't implemented yet
			t.Skipf("Method %s not implemented", methodName)
			return nil, nil
		}
	}

	// Validate all packets
	results, err := ValidateAllPackets("../testdata", generator)
	if err != nil {
		t.Fatalf("Error validating packets: %v", err)
	}

	// Print results
	PrintValidationResults(results)
}
