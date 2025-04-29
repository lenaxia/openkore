package item

import (
	"bytes"
	"testing"
)

// TestSendStoragePassword tests the SendStoragePassword method
func TestSendStoragePassword(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["storage_password"] = "023B"

	sm := NewStorageManager(mockSend)

	// Test sending storage password
	pass := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	type_ := 3 // 3 = check password
	err := sm.SendStoragePassword(pass, type_)
	if err != nil {
		t.Fatalf("SendStoragePassword() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["023B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if !bytes.Equal(args["pass"].([]byte), pass) {
		t.Errorf("args[\"pass\"] = %v, want %v", args["pass"], pass)
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}
}

// TestReconstructStoragePassword tests the ReconstructStoragePassword function
func TestReconstructStoragePassword(t *testing.T) {
	// Test with type = 3 (check password)
	args := map[string]interface{}{
		"type": 3,
		"pass": []byte{0x01, 0x02, 0x03, 0x04},
	}

	ReconstructStoragePassword(args)

	// Check that the data field was set
	if args["data"] == nil {
		t.Fatal("ReconstructStoragePassword did not set data field")
	}

	// Test with type = 2 (change password)
	args = map[string]interface{}{
		"type": 2,
		"pass": []byte{0x01, 0x02, 0x03, 0x04},
	}

	ReconstructStoragePassword(args)

	// Check that the data field was set
	if args["data"] == nil {
		t.Fatal("ReconstructStoragePassword did not set data field")
	}

	// Test with invalid type
	args = map[string]interface{}{
		"type": 1,
		"pass": []byte{0x01, 0x02, 0x03, 0x04},
	}

	err := ReconstructStoragePassword(args)
	if err == nil {
		t.Fatal("ReconstructStoragePassword did not return error for invalid type")
	}
}

// TestSendStorageGetToCart tests the SendStorageGetToCart method
func TestSendStorageGetToCart(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["storage_to_cart"] = "0126"
	mockSend.packetLUT["guild_storage_to_cart"] = "0127"

	sm := NewStorageManager(mockSend)

	// Test with normal storage
	index := uint16(1)
	amount := uint16(10)
	err := sm.SendStorageGetToCart(index, amount, false)
	if err != nil {
		t.Fatalf("SendStorageGetToCart() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0126"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != index {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}

	// Test with guild storage
	mockSend = NewMockSendForStorage()
	mockSend.packetLUT["storage_to_cart"] = "0126"
	mockSend.packetLUT["guild_storage_to_cart"] = "0127"

	sm = NewStorageManager(mockSend)

	err = sm.SendStorageGetToCart(index, amount, true)
	if err != nil {
		t.Fatalf("SendStorageGetToCart() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists = mockSend.reconstructArgs["0127"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != index {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestSendStorageAddFromCart tests the SendStorageAddFromCart method
func TestSendStorageAddFromCart(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["cart_to_storage"] = "0128"
	mockSend.packetLUT["cart_to_guild_storage"] = "0129"

	sm := NewStorageManager(mockSend)

	// Test with normal storage
	index := uint16(1)
	amount := uint16(10)
	err := sm.SendStorageAddFromCart(index, amount, false)
	if err != nil {
		t.Fatalf("SendStorageAddFromCart() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0128"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != index {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}

	// Test with guild storage
	mockSend = NewMockSendForStorage()
	mockSend.packetLUT["cart_to_storage"] = "0128"
	mockSend.packetLUT["cart_to_guild_storage"] = "0129"

	sm = NewStorageManager(mockSend)

	err = sm.SendStorageAddFromCart(index, amount, true)
	if err != nil {
		t.Fatalf("SendStorageAddFromCart() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists = mockSend.reconstructArgs["0129"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != index {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}
