package login

import (
	"reflect"
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// setupTestBaseSend creates a configured BaseSend for testing
func setupTestBaseSend() *core.BaseSend {
	// Create a BaseSend
	baseSend := core.NewBaseSend(hooks.NewHookManager())

	// Configure with test packet constructions
	packetConstructions := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
	}

	// Configure the BaseSend
	baseSend.Configure("TestServer", packetConstructions)

	return baseSend
}

// TestRegisterHandlers tests that handlers are properly registered
func TestRegisterHandlers(t *testing.T) {
	// Create a configured BaseSend
	mockSend := setupTestBaseSend()

	// Register handlers
	RegisterHandlers(mockSend)

	// Verify that the login_request handler was registered by trying to construct a packet
	_, err := mockSend.ConstructPacket("login_request", map[string]interface{}{})
	if err != nil {
		t.Error("login_request handler was not registered properly")
	}
}

// TestHandleLoginRequest tests the login_request handler
func TestHandleLoginRequest(t *testing.T) {
	// Test with empty args
	packet, err := handleLoginRequest(map[string]interface{}{})
	if err != nil {
		t.Errorf("handleLoginRequest returned error: %v", err)
	}

	// Verify the packet structure
	expectedPacket := []byte{0x64, 0x00, 0x01, 0x02, 0x03, 0x04}
	if !reflect.DeepEqual(packet, expectedPacket) {
		t.Errorf("handleLoginRequest returned %v, want %v", packet, expectedPacket)
	}

	// In a real implementation, we would test with various input arguments
	// and verify that the packet is constructed correctly
}

// TestLoginIntegration tests the integration between the login handlers and the BaseSend
func TestLoginIntegration(t *testing.T) {
	// Create a configured BaseSend
	baseSend := setupTestBaseSend()

	// Register handlers
	RegisterHandlers(baseSend)

	// Send a login request
	packet, err := baseSend.ConstructPacket("login_request", map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
		"version":  1234,
	})

	if err != nil {
		t.Errorf("ConstructPacket returned error: %v", err)
	}

	// In a real implementation, we would verify that the packet is constructed correctly
	// based on the input arguments
	if len(packet) == 0 {
		t.Error("ConstructPacket returned empty packet")
	}
}
