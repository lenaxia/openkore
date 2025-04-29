package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestSendPing tests the SendPing method
func TestSendPing(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Create a mock connection
	mockConn := &MockConnection{
		sentPackets: make([][]byte, 0),
	}

	// Set the connection
	baseSend.SetConnection(mockConn)

	// Register the ping packet with "x0" format (0 bytes of padding, effectively no data)
	baseSend.RegisterPacketHandler("0B1C", "ping", "x0", []string{}, nil)

	// Send a ping packet
	err := baseSend.SendPing()
	if err != nil {
		t.Fatalf("SendPing() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockConn.sentPackets) != 1 {
		t.Fatalf("len(mockConn.sentPackets) = %v, want %v", len(mockConn.sentPackets), 1)
	}

	// Check that the packet has the correct ID (0B1C)
	sentPacket := mockConn.sentPackets[0]
	if len(sentPacket) < 2 || sentPacket[0] != 0x1C || sentPacket[1] != 0x0B {
		t.Errorf("Incorrect packet ID: %v", sentPacket[:2])
	}
}

// TestSendPingNoConnection tests the SendPing method with no connection
func TestSendPingNoConnection(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send without a connection
	baseSend := NewBaseSend(hookManager)

	// Register the ping packet with "x0" format (0 bytes of padding, effectively no data)
	baseSend.RegisterPacketHandler("0B1C", "ping", "x0", []string{}, nil)

	// Send a ping packet
	err := baseSend.SendPing()

	// Check that an error was returned
	if err == nil {
		t.Fatal("SendPing() did not return an error with no connection")
	}

	// Check that the error is ErrNoConnection
	if err != ErrNoConnection {
		t.Errorf("SendPing() returned error %v, want %v", err, ErrNoConnection)
	}
}
