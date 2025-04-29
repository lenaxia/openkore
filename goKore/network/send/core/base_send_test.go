package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewBaseSend tests the NewBaseSend function
func TestNewBaseSend(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Check that the base send was created
	if baseSend == nil {
		t.Fatal("NewBaseSend() returned nil")
	}

	// Check that the hook manager was set
	if baseSend.hookManager != hookManager {
		t.Error("baseSend.hookManager was not set correctly")
	}

	// Check that the handlers map was initialized
	if baseSend.handlers == nil {
		t.Error("baseSend.handlers was not initialized")
	}

	// Check that the packet constructions map was initialized
	if baseSend.packetConstructions == nil {
		t.Error("baseSend.packetConstructions was not initialized")
	}

	// Check that the packet LUT map was initialized
	if baseSend.packetLUT == nil {
		t.Error("baseSend.packetLUT was not initialized")
	}

	// Check that the packet builder was initialized
	if baseSend.packetBuilder == nil {
		t.Error("baseSend.packetBuilder was not initialized")
	}
}

// TestBaseSendConfigure tests the Configure method
func TestBaseSendConfigure(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Create packet constructions
	packetConstructions := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
	}

	// Configure the base send
	err := baseSend.Configure("ServerType0", packetConstructions)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Check that the server type was set
	if baseSend.serverType != "ServerType0" {
		t.Errorf("baseSend.serverType = %v, want %v", baseSend.serverType, "ServerType0")
	}

	// Check that the packet constructions were set
	if len(baseSend.packetConstructions) != 1 {
		t.Errorf("len(baseSend.packetConstructions) = %v, want %v", len(baseSend.packetConstructions), 1)
	}

	// Check that the packet LUT was populated
	if len(baseSend.packetLUT) != 1 {
		t.Errorf("len(baseSend.packetLUT) = %v, want %v", len(baseSend.packetLUT), 1)
	}

	// Check that the packet ID was registered in the LUT
	packetID, exists := baseSend.packetLUT["login_request"]
	if !exists {
		t.Error("login_request was not registered in the packet LUT")
	} else if packetID != "0064" {
		t.Errorf("baseSend.packetLUT[\"login_request\"] = %v, want %v", packetID, "0064")
	}
}

// TestBaseSendRegisterHandler tests the RegisterHandler method
func TestBaseSendRegisterHandler(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Create a handler
	handler := func(args map[string]interface{}) ([]byte, error) {
		return []byte{0x01, 0x02, 0x03}, nil
	}

	// Register the handler
	baseSend.RegisterHandler("test_packet", handler)

	// Check that the handler was registered
	if len(baseSend.handlers) != 1 {
		t.Errorf("len(baseSend.handlers) = %v, want %v", len(baseSend.handlers), 1)
	}

	// Check that the handler can be retrieved
	registeredHandler, exists := baseSend.handlers["test_packet"]
	if !exists {
		t.Error("test_packet handler was not registered")
	} else {
		// Call the handler and check the result
		result, err := registeredHandler(nil)
		if err != nil {
			t.Errorf("handler() returned error: %v", err)
		} else if len(result) != 3 || result[0] != 0x01 || result[1] != 0x02 || result[2] != 0x03 {
			t.Errorf("handler() = %v, want %v", result, []byte{0x01, 0x02, 0x03})
		}
	}
}

// TestBaseSendConstructPacket tests the ConstructPacket method
func TestBaseSendConstructPacket(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Create packet constructions
	packetConstructions := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
	}

	// Configure the base send
	err := baseSend.Configure("ServerType0", packetConstructions)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Create a handler
	handler := func(args map[string]interface{}) ([]byte, error) {
		return []byte{0x01, 0x02, 0x03}, nil
	}

	// Register the handler
	baseSend.RegisterHandler("login_request", handler)

	// Test constructing a packet with a registered handler
	packet, err := baseSend.ConstructPacket("login_request", nil)
	if err != nil {
		t.Errorf("ConstructPacket() returned error: %v", err)
	} else if len(packet) != 3 || packet[0] != 0x01 || packet[1] != 0x02 || packet[2] != 0x03 {
		t.Errorf("ConstructPacket() = %v, want %v", packet, []byte{0x01, 0x02, 0x03})
	}

	// Test constructing a packet that doesn't exist
	_, err = baseSend.ConstructPacket("nonexistent_packet", nil)
	if err == nil {
		t.Error("ConstructPacket() did not return an error for a nonexistent packet")
	}
}

// TestBaseSendEncryptMessageID tests the encryptMessageID method
func TestBaseSendEncryptMessageID(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := NewBaseSend(hookManager)

	// Set encryption keys
	baseSend.SetEncryptionKeys(0x12345678, 0x87654321, 0x11223344)

	// Create a packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}
	packetCopy := make([]byte, len(packet))
	copy(packetCopy, packet)

	// Encrypt the message ID
	err := baseSend.encryptMessageID(&packetCopy)
	if err != nil {
		t.Fatalf("encryptMessageID() returned error: %v", err)
	}

	// Check that the packet was modified
	if packetCopy[0] == packet[0] && packetCopy[1] == packet[1] {
		t.Error("encryptMessageID() did not modify the packet")
	}
}

// MockConnection is a mock implementation of the connection interface
type MockConnection struct {
	sentPackets [][]byte
}

// Send implements the Send method for the connection interface
func (mc *MockConnection) Send(packet []byte) error {
	mc.sentPackets = append(mc.sentPackets, packet)
	return nil
}

// TestBaseSendSendToServer tests the SendToServer method
func TestBaseSendSendToServer(t *testing.T) {
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

	// Create a packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}

	// Send the packet
	err := baseSend.SendToServer(packet)
	if err != nil {
		t.Fatalf("SendToServer() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockConn.sentPackets) != 1 {
		t.Fatalf("len(mockConn.sentPackets) = %v, want %v", len(mockConn.sentPackets), 1)
	}

	// Check that the packet was sent correctly
	sentPacket := mockConn.sentPackets[0]
	if len(sentPacket) != len(packet) {
		t.Fatalf("len(sentPacket) = %v, want %v", len(sentPacket), len(packet))
	}

	for i := 0; i < len(packet); i++ {
		if sentPacket[i] != packet[i] {
			t.Errorf("sentPacket[%d] = %v, want %v", i, sentPacket[i], packet[i])
		}
	}
}

// TestBaseSendSendPacket tests the SendPacket method
func TestBaseSendSendPacket(t *testing.T) {
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

	// Create packet constructions
	packetConstructions := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
	}

	// Configure the base send
	err := baseSend.Configure("ServerType0", packetConstructions)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Create a handler
	handler := func(args map[string]interface{}) ([]byte, error) {
		return []byte{0x01, 0x02, 0x03, 0x04}, nil
	}

	// Register the handler
	baseSend.RegisterHandler("login_request", handler)

	// Send a packet
	err = baseSend.SendPacket("login_request", nil)
	if err != nil {
		t.Fatalf("SendPacket() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockConn.sentPackets) != 1 {
		t.Fatalf("len(mockConn.sentPackets) = %v, want %v", len(mockConn.sentPackets), 1)
	}

	// Check that the packet was sent correctly
	sentPacket := mockConn.sentPackets[0]
	expectedPacket := []byte{0x01, 0x02, 0x03, 0x04}
	if len(sentPacket) != len(expectedPacket) {
		t.Fatalf("len(sentPacket) = %v, want %v", len(sentPacket), len(expectedPacket))
	}

	for i := 0; i < len(expectedPacket); i++ {
		if sentPacket[i] != expectedPacket[i] {
			t.Errorf("sentPacket[%d] = %v, want %v", i, sentPacket[i], expectedPacket[i])
		}
	}
}
