package storage

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestGuildStorageLog(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("storage.guild_storage_log", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewStorageManager(nil, hookManager)

	// Test case 1: Log with items
	t.Run("LogWithItems", func(t *testing.T) {
		// Create test data
		// Create a log with 2 items
		log := make([]byte, 0)

		// Item 1
		item1 := []byte{
			// ID (4 bytes)
			0x01, 0x00, 0x00, 0x00,
			// nameID (2 bytes)
			0x02, 0x00,
			// amount (4 bytes)
			0x03, 0x00, 0x00, 0x00,
			// action (1 byte)
			0x00,
			// upgrade (4 bytes)
			0x04, 0x00, 0x00, 0x00,
			// uniqueID (8 bytes)
			0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// identified (1 byte)
			0x01,
			// type_equip (2 bytes)
			0x06, 0x00,
			// cards (8 bytes)
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// charName (24 bytes)
			'P', 'l', 'a', 'y', 'e', 'r', '1', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// time (24 bytes)
			'2', '0', '2', '5', '-', '0', '4', '-', '2', '8', ' ', '2', '1', ':', '0', '0', ':', '0', '0', 0x00, 0x00, 0x00, 0x00, 0x00,
			// attribute (1 byte)
			0x00,
		}
		log = append(log, item1...)

		// Item 2
		item2 := []byte{
			// ID (4 bytes)
			0x07, 0x00, 0x00, 0x00,
			// nameID (2 bytes)
			0x08, 0x00,
			// amount (4 bytes)
			0x09, 0x00, 0x00, 0x00,
			// action (1 byte)
			0x01,
			// upgrade (4 bytes)
			0x0A, 0x00, 0x00, 0x00,
			// uniqueID (8 bytes)
			0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// identified (1 byte)
			0x01,
			// type_equip (2 bytes)
			0x0C, 0x00,
			// cards (8 bytes)
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// charName (24 bytes)
			'P', 'l', 'a', 'y', 'e', 'r', '2', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// time (24 bytes)
			'2', '0', '2', '5', '-', '0', '4', '-', '2', '8', ' ', '2', '1', ':', '0', '5', ':', '0', '0', 0x00, 0x00, 0x00, 0x00, 0x00,
			// attribute (1 byte)
			0x00,
		}
		log = append(log, item2...)

		args := map[string]interface{}{
			"result": byte(0),
			"log":    log,
		}

		// Call the handler
		err := manager.HandleGuildStorageLog(args)
		if err != nil {
			t.Errorf("HandleGuildStorageLog returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["result"] != byte(0) {
			t.Errorf("Expected result 0, got %d", result["result"])
		}

		items := result["items"].([]map[string]interface{})
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}

		// Check the first item
		item := items[0]
		if item["ID"] != uint32(1) {
			t.Errorf("Expected ID 1, got %d", item["ID"])
		}
		if item["nameID"] != uint16(2) {
			t.Errorf("Expected nameID 2, got %d", item["nameID"])
		}
		if item["amount"] != uint32(3) {
			t.Errorf("Expected amount 3, got %d", item["amount"])
		}
		if item["action"] != byte(0) {
			t.Errorf("Expected action 0, got %d", item["action"])
		}
		if item["charName"] != "Player1" {
			t.Errorf("Expected charName 'Player1', got '%s'", item["charName"])
		}
		if item["time"] != "2025-04-28 21:00:00" {
			t.Errorf("Expected time '2025-04-28 21:00:00', got '%s'", item["time"])
		}

		// Check the second item
		item = items[1]
		if item["ID"] != uint32(7) {
			t.Errorf("Expected ID 7, got %d", item["ID"])
		}
		if item["nameID"] != uint16(8) {
			t.Errorf("Expected nameID 8, got %d", item["nameID"])
		}
		if item["amount"] != uint32(9) {
			t.Errorf("Expected amount 9, got %d", item["amount"])
		}
		if item["action"] != byte(1) {
			t.Errorf("Expected action 1, got %d", item["action"])
		}
		if item["charName"] != "Player2" {
			t.Errorf("Expected charName 'Player2', got '%s'", item["charName"])
		}
		if item["time"] != "2025-04-28 21:05:00" {
			t.Errorf("Expected time '2025-04-28 21:05:00', got '%s'", item["time"])
		}

		// Check the status message
		status := result["status"].(string)
		if len(status) == 0 {
			t.Errorf("Status message is empty")
		}
	})

	// Test case 2: Empty storage
	t.Run("EmptyStorage", func(t *testing.T) {
		args := map[string]interface{}{
			"result": byte(2),
		}

		// Call the handler
		err := manager.HandleGuildStorageLog(args)
		if err != nil {
			t.Errorf("HandleGuildStorageLog returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["result"] != byte(2) {
			t.Errorf("Expected result 2, got %d", result["result"])
		}

		// Check the status message
		status := result["status"].(string)
		if status != "Guild Storage empty." {
			t.Errorf("Expected status 'Guild Storage empty.', got '%s'", status)
		}
	})

	// Test case 3: Not using guild storage
	t.Run("NotUsingGuildStorage", func(t *testing.T) {
		args := map[string]interface{}{
			"result": byte(3),
		}

		// Call the handler
		err := manager.HandleGuildStorageLog(args)
		if err != nil {
			t.Errorf("HandleGuildStorageLog returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["result"] != byte(3) {
			t.Errorf("Expected result 3, got %d", result["result"])
		}

		// Check the status message
		status := result["status"].(string)
		if status != "You are not currently using Guild Storage. Please try later." {
			t.Errorf("Expected status 'You are not currently using Guild Storage. Please try later.', got '%s'", status)
		}
	})
}
