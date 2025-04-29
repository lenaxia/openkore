package game

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
		"0085": {
			ID:         "0085",
			Name:       "move_to",
			Format:     "v3",
			FieldNames: []string{"x", "y", "unknown"},
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

	// Verify that the move_to handler was registered by trying to construct a packet
	_, err := mockSend.ConstructPacket("move_to", map[string]interface{}{})
	if err != nil {
		t.Error("move_to handler was not registered properly")
	}
}

// TestHandleMoveTo tests the move_to handler
func TestHandleMoveTo(t *testing.T) {
	// Test with empty args
	packet, err := handleMoveTo(map[string]interface{}{})
	if err != nil {
		t.Errorf("handleMoveTo returned error: %v", err)
	}

	// Verify the packet structure
	expectedPacket := []byte{0x85, 0x00, 0x01, 0x02, 0x03, 0x04}
	if !reflect.DeepEqual(packet, expectedPacket) {
		t.Errorf("handleMoveTo returned %v, want %v", packet, expectedPacket)
	}

	// In a real implementation, we would test with various input arguments
	// and verify that the packet is constructed correctly
}

// TestGameIntegration tests the integration between the game handlers and the BaseSend
func TestGameIntegration(t *testing.T) {
	// Create a configured BaseSend
	baseSend := setupTestBaseSend()

	// Register handlers
	RegisterHandlers(baseSend)

	// Send a move_to request
	packet, err := baseSend.ConstructPacket("move_to", map[string]interface{}{
		"x": 100,
		"y": 200,
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
