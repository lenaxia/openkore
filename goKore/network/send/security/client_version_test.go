package security

import (
	"testing"
)

// TestSendClientVersion tests the SendClientVersion method
func TestSendClientVersion(t *testing.T) {
	mockSend := NewMockSendForPin()
	mockSend.packetLUT["client_version"] = "044A"

	// Create a login manager
	lm := NewLoginManager(mockSend)

	// Test sending client version
	version := 12345
	err := lm.SendClientVersion(version)
	if err != nil {
		t.Fatalf("SendClientVersion() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["044A"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["clientVersion"] != version {
		t.Errorf("args[\"clientVersion\"] = %v, want %v", args["clientVersion"], version)
	}
}
