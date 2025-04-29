package item

import (
	"testing"
)

// TestNewCartManager tests the NewCartManager function
func TestNewCartManager(t *testing.T) {
	mockSend := NewMockSendForInventory()
	cm := NewCartManager(mockSend)

	if cm == nil {
		t.Fatal("NewCartManager() returned nil")
	}

	if cm.baseSend == nil {
		t.Error("cm.baseSend was not set correctly")
	}
}

// TestSendCartAdd tests the SendCartAdd method
func TestSendCartAdd(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["cart_add"] = "0126"

	cm := NewCartManager(mockSend)

	// Test adding an item to cart
	index := uint16(1)
	amount := uint16(10)
	err := cm.SendCartAdd(index, amount)
	if err != nil {
		t.Fatalf("SendCartAdd() returned error: %v", err)
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
}

// TestSendCartGet tests the SendCartGet method
func TestSendCartGet(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["cart_get"] = "0127"

	cm := NewCartManager(mockSend)

	// Test getting an item from cart
	index := uint16(1)
	amount := uint16(10)
	err := cm.SendCartGet(index, amount)
	if err != nil {
		t.Fatalf("SendCartGet() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0127"]
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

// TestSendChangeCart tests the SendChangeCart method
func TestSendChangeCart(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["change_cart"] = "01B0"

	cm := NewCartManager(mockSend)

	// Test changing cart level
	level := 3 // 1..5
	err := cm.SendChangeCart(level)
	if err != nil {
		t.Fatalf("SendChangeCart() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["01B0"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["lvl"] != level {
		t.Errorf("args[\"lvl\"] = %v, want %v", args["lvl"], level)
	}
}
