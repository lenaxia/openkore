package item

import (
	"testing"
)

// TestSendEquip tests the SendEquip method
func TestSendEquip(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["equip"] = "00A9"

	im := NewInventoryManager(mockSend)

	// Test equipping an item
	inventoryIndex := uint32(1)
	equipType := uint32(2) // 2 = armor
	err := im.SendEquip(inventoryIndex, equipType)
	if err != nil {
		t.Fatalf("SendEquip() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00A9"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["inventory_index"] != inventoryIndex {
		t.Errorf("args[\"inventory_index\"] = %v, want %v", args["inventory_index"], inventoryIndex)
	}

	if args["equip_type"] != equipType {
		t.Errorf("args[\"equip_type\"] = %v, want %v", args["equip_type"], equipType)
	}
}

// TestSendUnequip tests the SendUnequip method
func TestSendUnequip(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["unequip"] = "00AB"

	im := NewInventoryManager(mockSend)

	// Test unequipping an item
	equipIndex := uint32(1)
	err := im.SendUnequip(equipIndex)
	if err != nil {
		t.Fatalf("SendUnequip() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00AB"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["equip_index"] != equipIndex {
		t.Errorf("args[\"equip_index\"] = %v, want %v", args["equip_index"], equipIndex)
	}
}

// TestSendEquipSwitchAdd tests the SendEquipSwitchAdd method
func TestSendEquipSwitchAdd(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["equip_switch_add"] = "0A97"

	im := NewInventoryManager(mockSend)

	// Test adding an item to the equip switch window
	inventoryIndex := uint32(1)
	position := uint32(2)
	err := im.SendEquipSwitchAdd(inventoryIndex, position)
	if err != nil {
		t.Fatalf("SendEquipSwitchAdd() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A97"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != inventoryIndex {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], inventoryIndex)
	}

	if args["position"] != position {
		t.Errorf("args[\"position\"] = %v, want %v", args["position"], position)
	}
}

// TestSendEquipSwitchRemove tests the SendEquipSwitchRemove method
func TestSendEquipSwitchRemove(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["equip_switch_remove"] = "0A99"

	im := NewInventoryManager(mockSend)

	// Test removing an item from the equip switch window
	inventoryIndex := uint32(1)
	position := uint32(2)
	err := im.SendEquipSwitchRemove(inventoryIndex, position)
	if err != nil {
		t.Fatalf("SendEquipSwitchRemove() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A99"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != inventoryIndex {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], inventoryIndex)
	}

	if args["position"] != position {
		t.Errorf("args[\"position\"] = %v, want %v", args["position"], position)
	}
}

// TestSendEquipSwitchRun tests the SendEquipSwitchRun method
func TestSendEquipSwitchRun(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["equip_switch_run"] = "0A9C"

	im := NewInventoryManager(mockSend)

	// Test running the equip switch
	err := im.SendEquipSwitchRun()
	if err != nil {
		t.Fatalf("SendEquipSwitchRun() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0A9C"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendEquipSwitchSingle tests the SendEquipSwitchSingle method
func TestSendEquipSwitchSingle(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["equip_switch_single"] = "0ACE"

	im := NewInventoryManager(mockSend)

	// Test running a single equip switch
	inventoryIndex := uint32(1)
	err := im.SendEquipSwitchSingle(inventoryIndex)
	if err != nil {
		t.Fatalf("SendEquipSwitchSingle() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0ACE"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != inventoryIndex {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], inventoryIndex)
	}
}

// TestSendChangeDress tests the SendChangeDress method
func TestSendChangeDress(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["change_dress"] = "0998"

	im := NewInventoryManager(mockSend)

	// Test changing dress
	err := im.SendChangeDress()
	if err != nil {
		t.Fatalf("SendChangeDress() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0998"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendInventoryExpansionRequest tests the SendInventoryExpansionRequest method
func TestSendInventoryExpansionRequest(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["inventory_expansion_request"] = "0B14"

	im := NewInventoryManager(mockSend)

	// Test requesting inventory expansion
	err := im.SendInventoryExpansionRequest()
	if err != nil {
		t.Fatalf("SendInventoryExpansionRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0B14"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendInventoryExpansionRejected tests the SendInventoryExpansionRejected method
func TestSendInventoryExpansionRejected(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["inventory_expansion_rejected"] = "0B19"

	im := NewInventoryManager(mockSend)

	// Test rejecting inventory expansion
	err := im.SendInventoryExpansionRejected()
	if err != nil {
		t.Fatalf("SendInventoryExpansionRejected() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0B19"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
