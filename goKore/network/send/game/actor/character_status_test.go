package actor

import (
	"testing"
)

// TestSendAddStatusPoint tests the SendAddStatusPoint function
func TestSendAddStatusPoint(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["send_add_status_point"] = "00BB"
	charManager := NewCharacterManager(mockSend)

	statusID := 1 // STR
	amount := 1

	err := charManager.SendAddStatusPoint(statusID, amount)
	if err != nil {
		t.Fatalf("SendAddStatusPoint() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "send_add_status_point" {
		t.Errorf("Expected packet ID 'send_add_status_point', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["statusID"].(int) != statusID {
		t.Errorf("Expected statusID=%d, got %v", statusID, args["statusID"])
	}
	if args["Amount"].(int) != amount {
		t.Errorf("Expected Amount=%d, got %v", amount, args["Amount"])
	}
}

// TestSendAddSkillPoint tests the SendAddSkillPoint function
func TestSendAddSkillPoint(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["send_add_skill_point"] = "0112"
	charManager := NewCharacterManager(mockSend)

	skillID := 123

	err := charManager.SendAddSkillPoint(skillID)
	if err != nil {
		t.Fatalf("SendAddSkillPoint() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "send_add_skill_point" {
		t.Errorf("Expected packet ID 'send_add_skill_point', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["skillID"].(int) != skillID {
		t.Errorf("Expected skillID=%d, got %v", skillID, args["skillID"])
	}
}

// TestSendHotKeyChange tests the SendHotKeyChange function
func TestSendHotKeyChange(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["hotkey_change"] = "0212"
	charManager := NewCharacterManager(mockSend)

	idx := 0
	type_ := 1
	id := 123
	lvl := 5

	err := charManager.SendHotKeyChange(idx, type_, id, lvl)
	if err != nil {
		t.Fatalf("SendHotKeyChange() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "hotkey_change" {
		t.Errorf("Expected packet ID 'hotkey_change', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["idx"].(int) != idx {
		t.Errorf("Expected idx=%d, got %v", idx, args["idx"])
	}
	if args["type"].(int) != type_ {
		t.Errorf("Expected type=%d, got %v", type_, args["type"])
	}
	if args["id"].(int) != id {
		t.Errorf("Expected id=%d, got %v", id, args["id"])
	}
	if args["lvl"].(int) != lvl {
		t.Errorf("Expected lvl=%d, got %v", lvl, args["lvl"])
	}
}

// TestSendChangeTitle tests the SendChangeTitle function
func TestSendChangeTitle(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["send_change_title"] = "0A2E"
	charManager := NewCharacterManager(mockSend)

	titleID := 123

	err := charManager.SendChangeTitle(titleID)
	if err != nil {
		t.Fatalf("SendChangeTitle() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "send_change_title" {
		t.Errorf("Expected packet ID 'send_change_title', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["ID"].(int) != titleID {
		t.Errorf("Expected ID=%d, got %v", titleID, args["ID"])
	}
}

// TestSendQuit tests the SendQuit function
func TestSendQuit(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["quit_request"] = "00B3" // Use the correct packet ID that doesn't conflict with "restart"
	charManager := NewCharacterManager(mockSend)

	err := charManager.SendQuit()
	if err != nil {
		t.Fatalf("SendQuit() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "quit_request" {
		t.Errorf("Expected packet ID 'quit_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["type"].(int) != 0 {
		t.Errorf("Expected type=0, got %v", args["type"])
	}
}

// TestSendAutoRevive tests the SendAutoRevive function
func TestSendAutoRevive(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["auto_revive"] = "0292"
	charManager := NewCharacterManager(mockSend)

	err := charManager.SendAutoRevive()
	if err != nil {
		t.Fatalf("SendAutoRevive() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auto_revive" {
		t.Errorf("Expected packet ID 'auto_revive', got '%s'", packetID)
	}
}
