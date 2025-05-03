package main

import (
	"bytes"
	"testing"
)

// This is a simple example of how to use the validation data in Go tests

// MockSendImplementation is a simple mock implementation of the Send interface
type MockSendImplementation struct {
	// Add any fields needed for your implementation
}

// SendMove is a mock implementation of the sendMove method
func (m *MockSendImplementation) SendMove(x, y int) []byte {
	// This is where you would implement your actual packet construction logic
	// For this example, we'll just return a hardcoded packet
	return []byte{0x85, 0x64, 0x00, 0x00, 0x64, 0x00, 0x00, 0xA9, 0x86, 0x01, 0x00}
}

// SendChat is a mock implementation of the sendChat method
func (m *MockSendImplementation) SendChat(message string) []byte {
	// This is where you would implement your actual packet construction logic
	// For this example, we'll just return a hardcoded packet
	return []byte{0x8C, 0x00, 0x00, 0x00, 0x00, 0x00}
}

// TestSendMoveWithValidationData demonstrates how to test your implementation against validation data
func TestSendMoveWithValidationData(t *testing.T) {
	// Load validation data
	data, err := loadValidationData("../validation_data/packet_output_sendMove.json")
	if err != nil {
		// Skip the test if the validation data doesn't exist
		// This allows the test to run even if the validation data hasn't been generated yet
		t.Skip("Validation data not found. Run generate_validation_data.sh first.")
		return
	}

	// Create your implementation
	impl := &MockSendImplementation{}

	// Extract arguments from validation data
	x, ok1 := data.Args[0].(float64)
	y, ok2 := data.Args[1].(float64)
	if !ok1 || !ok2 {
		t.Fatalf("Invalid arguments in validation data: %v", data.Args)
	}

	// Call your implementation with the same arguments
	result := impl.SendMove(int(x), int(y))

	// Convert validation data bytes to a byte slice
	expectedBytes := make([]byte, len(data.Packets[0].Bytes))
	for i, b := range data.Packets[0].Bytes {
		expectedBytes[i] = byte(b)
	}

	// Compare the results
	if !bytes.Equal(result, expectedBytes) {
		t.Errorf("SendMove output doesn't match validation data")
		t.Errorf("Expected: %v", expectedBytes)
		t.Errorf("Got: %v", result)
	}
}

// TestSendChatWithValidationData demonstrates how to test your implementation against validation data
func TestSendChatWithValidationData(t *testing.T) {
	// Load validation data
	data, err := loadValidationData("../validation_data/packet_output_sendChat.json")
	if err != nil {
		// Skip the test if the validation data doesn't exist
		t.Skip("Validation data not found. Run generate_validation_data.sh first.")
		return
	}

	// Create your implementation
	impl := &MockSendImplementation{}

	// Extract arguments from validation data
	message, ok := data.Args[0].(string)
	if !ok {
		t.Fatalf("Invalid arguments in validation data: %v", data.Args)
	}

	// Call your implementation with the same arguments
	result := impl.SendChat(message)

	// Convert validation data bytes to a byte slice
	expectedBytes := make([]byte, len(data.Packets[0].Bytes))
	for i, b := range data.Packets[0].Bytes {
		expectedBytes[i] = byte(b)
	}

	// Compare the results
	if !bytes.Equal(result, expectedBytes) {
		t.Errorf("SendChat output doesn't match validation data")
		t.Errorf("Expected: %v", expectedBytes)
		t.Errorf("Got: %v", result)
	}
}

// Note: We're using the loadValidationData function from main.go
// No need to redefine it here
