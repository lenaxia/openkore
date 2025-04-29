package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/send/core"
)

// TestSendSync tests the SendSync function
func TestSendSync(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	syncManager := NewSyncManager(mockSend)

	err := syncManager.SendSync(false)
	if err != nil {
		t.Fatalf("SendSync() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "sync" {
		t.Errorf("Expected packet ID 'sync', got '%s'", packetID)
	}

	// Check that the time argument was set
	if _, ok := mockSend.LastArgs()["time"]; !ok {
		t.Error("Expected time argument to be set")
	}
}

// TestSendCharacterMove tests the SendCharacterMove function
func TestSendCharacterMove(t *testing.T) {
	mockSend := NewMockSend()
	syncManager := NewSyncManager(mockSend)

	x, y := 100, 200
	err := syncManager.SendCharacterMove(x, y)
	if err != nil {
		t.Fatalf("SendCharacterMove() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "character_move" {
		t.Errorf("Expected packet ID 'character_move', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if xVal, ok := mockSend.LastArgs()["x"].(uint16); !ok || int(xVal) != x {
		t.Errorf("Expected x=%d, got %v", x, mockSend.LastArgs()["x"])
	}

	if yVal, ok := mockSend.LastArgs()["y"].(uint16); !ok || int(yVal) != y {
		t.Errorf("Expected y=%d, got %v", y, mockSend.LastArgs()["y"])
	}

	if _, ok := mockSend.LastArgs()["time"]; !ok {
		t.Error("Expected time argument to be set")
	}
}
