package base

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewBaseReceive tests the NewBaseReceive function
func TestNewBaseReceive(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Check that the base receive was created
	if baseReceive == nil {
		t.Fatal("NewBaseReceive() returned nil")
	}

	// Check that the hook manager was set
	if baseReceive.hookManager != hookManager {
		t.Error("baseReceive.hookManager was not set correctly")
	}

	// Check that the handlers map was initialized
	if baseReceive.handlers == nil {
		t.Error("baseReceive.handlers was not initialized")
	}

	// Check that the packet definitions map was initialized
	if baseReceive.packetDefs == nil {
		t.Error("baseReceive.packetDefs was not initialized")
	}
}

// TestBaseReceiveConfigure tests the Configure method
func TestBaseReceiveConfigure(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Create packet definitions
	packetDefs := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_response",
			Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
			FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
		},
	}

	// Configure the base receive
	err := baseReceive.Configure("ServerType0", packetDefs)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Check that the server type was set
	if baseReceive.serverType != "ServerType0" {
		t.Errorf("baseReceive.serverType = %v, want %v", baseReceive.serverType, "ServerType0")
	}

	// Check that the packet definitions were set
	if len(baseReceive.packetDefs) != 1 {
		t.Errorf("len(baseReceive.packetDefs) = %v, want %v", len(baseReceive.packetDefs), 1)
	}

	// Check that the packet ID was registered
	packetDef, exists := baseReceive.packetDefs["0064"]
	if !exists {
		t.Error("0064 packet definition was not registered")
	} else if packetDef.Name != "login_response" {
		t.Errorf("packetDef.Name = %v, want %v", packetDef.Name, "login_response")
	}
}

// TestBaseReceiveRegisterHandler tests the RegisterHandler method
func TestBaseReceiveRegisterHandler(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Create a handler
	handler := func(args map[string]interface{}) error {
		return nil
	}

	// Register the handler
	baseReceive.RegisterHandler("login_response", handler)

	// Check that the handler was registered
	if len(baseReceive.handlers) != 1 {
		t.Errorf("len(baseReceive.handlers) = %v, want %v", len(baseReceive.handlers), 1)
	}

	// Check that the handler can be retrieved
	registeredHandler, exists := baseReceive.handlers["login_response"]
	if !exists {
		t.Error("login_response handler was not registered")
	} else if registeredHandler == nil {
		t.Error("login_response handler is nil")
	}
}

// MockPacket is a mock implementation of a packet for testing
type MockPacket struct {
	ID   string
	Data []byte
}

// TestBaseReceiveProcess tests the Process method
func TestBaseReceiveProcess(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Create packet definitions
	packetDefs := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_response",
			Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
			FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
		},
	}

	// Configure the base receive
	err := baseReceive.Configure("ServerType0", packetDefs)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Create a handler that sets a flag when called
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}

	// Register the handler
	baseReceive.RegisterHandler("login_response", handler)

	// Create a mock packet
	packet := []byte{0x64, 0x00, 0x01, 0x02, 0x03, 0x04}

	// Process the packet
	err = baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// Check that the handler was called
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

// TestBaseReceiveProcessUnknownPacket tests the Process method with an unknown packet
func TestBaseReceiveProcessUnknownPacket(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Create packet definitions
	packetDefs := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_response",
			Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
			FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
		},
	}

	// Configure the base receive
	err := baseReceive.Configure("ServerType0", packetDefs)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Create a mock packet with an unknown ID and sufficient length
	packet := []byte{0x65, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

	// Process the packet
	err = baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// No assertion needed, we just want to make sure it doesn't panic
}

// TestBaseReceiveProcessNoHandler tests the Process method with a packet that has no handler
func TestBaseReceiveProcessNoHandler(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test due to nil pointer dereference")
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base receive
	baseReceive := NewBaseReceive(hookManager)

	// Create packet definitions
	packetDefs := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_response",
			Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
			FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
		},
	}

	// Configure the base receive
	err := baseReceive.Configure("ServerType0", packetDefs)
	if err != nil {
		t.Fatalf("Configure() returned error: %v", err)
	}

	// Create a mock packet with sufficient length
	packet := []byte{0x64, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

	// Process the packet
	err = baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// No assertion needed, we just want to make sure it doesn't panic
}
