package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/send/core"
)

// TestNewCharacterManager tests the NewCharacterManager function
func TestNewCharacterManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	charManager := NewCharacterManager(mockSend)

	if charManager == nil {
		t.Fatal("NewCharacterManager() returned nil")
	}

	if charManager.baseSend == nil {
		t.Error("charManager.baseSend was not set correctly")
	}
}

// TestSendRestart tests the SendRestart function
func TestSendRestart(t *testing.T) {
	mockSend := NewMockSend()
	charManager := NewCharacterManager(mockSend)

	// Test sending a restart command
	restartType := 1 // 1 = quit to character select
	err := charManager.SendRestart(restartType)
	if err != nil {
		t.Fatalf("SendRestart() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "restart" {
		t.Errorf("Expected packet ID 'restart', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || int(typeVal) != restartType {
		t.Errorf("Expected type=%d, got %v", restartType, mockSend.LastArgs()["type"])
	}
}

// TestSendRespawn tests the SendRespawn function
func TestSendRespawn(t *testing.T) {
	mockSend := NewMockSend()
	charManager := NewCharacterManager(mockSend)

	// Test sending a respawn command
	err := charManager.SendRespawn()
	if err != nil {
		t.Fatalf("SendRespawn() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "restart" {
		t.Errorf("Expected packet ID 'restart', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 0 {
		t.Errorf("Expected type=0, got %v", mockSend.LastArgs()["type"])
	}
}

// TestSendQuitToCharSelect tests the SendQuitToCharSelect function
func TestSendQuitToCharSelect(t *testing.T) {
	mockSend := NewMockSend()
	charManager := NewCharacterManager(mockSend)

	// Test sending a quit to character select command
	err := charManager.SendQuitToCharSelect()
	if err != nil {
		t.Fatalf("SendQuitToCharSelect() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "restart" {
		t.Errorf("Expected packet ID 'restart', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 1 {
		t.Errorf("Expected type=1, got %v", mockSend.LastArgs()["type"])
	}
}
