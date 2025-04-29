package actor

import (
	"testing"
)

// TestSendSlaveAttack tests the SendSlaveAttack function
func TestSendSlaveAttack(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["slave_attack"] = "0233"
	am := NewActionManager(mockSend)

	// Test sending slave attack command
	slaveID := uint32(12345)
	targetID := uint32(67890)
	flag := 0
	err := am.SendSlaveAttack(slaveID, targetID, flag)
	if err != nil {
		t.Fatalf("SendSlaveAttack() returned error: %v", err)
	}

	// Check that the packet was sent
	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "slave_attack" {
		t.Errorf("Expected packet ID 'slave_attack', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["slaveID"] != slaveID {
		t.Errorf("Expected slaveID=%v, got %v", slaveID, args["slaveID"])
	}
	if args["targetID"] != targetID {
		t.Errorf("Expected targetID=%v, got %v", targetID, args["targetID"])
	}
	if args["flag"] != flag {
		t.Errorf("Expected flag=%v, got %v", flag, args["flag"])
	}
}

// TestSendSlaveStandBy tests the SendSlaveStandBy function
func TestSendSlaveStandBy(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["slave_move_to_master"] = "0234"
	am := NewActionManager(mockSend)

	// Test sending slave standby command
	slaveID := uint32(12345)
	err := am.SendSlaveStandBy(slaveID)
	if err != nil {
		t.Fatalf("SendSlaveStandBy() returned error: %v", err)
	}

	// Check that the packet was sent
	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "slave_move_to_master" {
		t.Errorf("Expected packet ID 'slave_move_to_master', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["slaveID"] != slaveID {
		t.Errorf("Expected slaveID=%v, got %v", slaveID, args["slaveID"])
	}
}

// TestSendAlignment tests the SendAlignment function
func TestSendAlignment(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["alignment"] = "0213"
	am := NewActionManager(mockSend)

	// Test sending alignment command
	targetID := uint32(12345)
	alignment := 1
	point := 500
	err := am.SendAlignment(targetID, alignment, point)
	if err != nil {
		t.Fatalf("SendAlignment() returned error: %v", err)
	}

	// Check that the packet was sent
	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "alignment" {
		t.Errorf("Expected packet ID 'alignment', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["targetID"] != targetID {
		t.Errorf("Expected targetID=%v, got %v", targetID, args["targetID"])
	}
	if args["type"] != alignment {
		t.Errorf("Expected type=%v, got %v", alignment, args["type"])
	}
	if args["point"] != point {
		t.Errorf("Expected point=%v, got %v", point, args["point"])
	}
}
