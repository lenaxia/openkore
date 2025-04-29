package card

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestNewCardManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCardManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewCardManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if len(manager.cardMergeItemsID) != 0 {
		t.Errorf("manager.cardMergeItemsID = %v, want empty slice", manager.cardMergeItemsID)
	}
}

func TestCardRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCardManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"card_merge_list",
		"card_merge_status",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestClearCardMergeItems(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCardManager(parser, nil)

	// Set some items
	manager.cardMergeItemsID = []int{1, 2, 3}

	// Clear items
	manager.ClearCardMergeItems()

	// Check if items were cleared
	if len(manager.cardMergeItemsID) != 0 {
		t.Errorf("manager.cardMergeItemsID = %v, want empty slice", manager.cardMergeItemsID)
	}
}

func TestGetCardMergeItems(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCardManager(parser, nil)

	// Set some items
	manager.cardMergeItemsID = []int{1, 2, 3}

	// Get items
	items := manager.GetCardMergeItems()

	// Check if items match
	if len(items) != 3 {
		t.Errorf("len(items) = %d, want 3", len(items))
	}

	for i, item := range items {
		if item != manager.cardMergeItemsID[i] {
			t.Errorf("items[%d] = %d, want %d", i, item, manager.cardMergeItemsID[i])
		}
	}

	// Verify that the returned slice is a copy
	items[0] = 999
	if manager.cardMergeItemsID[0] == 999 {
		t.Error("GetCardMergeItems() should return a copy of the slice, not a reference")
	}
}

func TestHandleCardMergeList(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCardManager(parser, nil)

	// Test with valid item list
	itemList := []byte{1, 0, 2, 0, 3, 0} // Item IDs 1, 2, 3
	args := map[string]interface{}{
		"item_list": itemList,
	}

	err := manager.handleCardMergeList(args)
	if err != nil {
		t.Fatalf("handleCardMergeList() returned error: %v", err)
	}

	// Check if items were added
	items := manager.GetCardMergeItems()
	if len(items) != 3 {
		t.Errorf("len(items) = %d, want 3", len(items))
	}

	expectedItems := []int{1, 2, 3}
	for i, item := range items {
		if item != expectedItems[i] {
			t.Errorf("items[%d] = %d, want %d", i, item, expectedItems[i])
		}
	}

	// Test with empty item list
	manager.ClearCardMergeItems()
	args = map[string]interface{}{
		"item_list": []byte{},
	}

	err = manager.handleCardMergeList(args)
	if err != nil {
		t.Fatalf("handleCardMergeList() with empty list returned error: %v", err)
	}

	// Check if items list is empty
	items = manager.GetCardMergeItems()
	if len(items) != 0 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}
}

func TestHandleCardMergeStatus(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCardManager(parser, nil)

	// Set some items
	manager.cardMergeItemsID = []int{1, 2, 3}

	// Test with success (fail = 0)
	args := map[string]interface{}{
		"fail":       byte(0),
		"item_index": byte(5),
		"card_index": byte(10),
	}

	err := manager.handleCardMergeStatus(args)
	if err != nil {
		t.Fatalf("handleCardMergeStatus() returned error: %v", err)
	}

	// Check if items were cleared
	items := manager.GetCardMergeItems()
	if len(items) != 0 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}

	// Test with failure (fail = 1)
	manager.cardMergeItemsID = []int{1, 2, 3}
	args = map[string]interface{}{
		"fail":       byte(1),
		"item_index": byte(5),
		"card_index": byte(10),
	}

	err = manager.handleCardMergeStatus(args)
	if err != nil {
		t.Fatalf("handleCardMergeStatus() with fail=1 returned error: %v", err)
	}

	// Check if items were cleared
	items = manager.GetCardMergeItems()
	if len(items) != 0 {
		t.Errorf("len(items) = %d, want 0", len(items))
	}
}
