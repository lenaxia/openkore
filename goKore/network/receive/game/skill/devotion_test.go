package skill

import (
	"encoding/binary"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestDevotion(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the devotion manager
	manager := NewDevotionManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for devotion
	t.Run("Basic Devotion", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.devotion", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create target IDs
		targetIDs := make([]byte, 20)                        // 5 target IDs, 4 bytes each
		binary.LittleEndian.PutUint32(targetIDs[0:4], 1001)  // Target 1
		binary.LittleEndian.PutUint32(targetIDs[4:8], 1002)  // Target 2
		binary.LittleEndian.PutUint32(targetIDs[8:12], 1003) // Target 3
		binary.LittleEndian.PutUint32(targetIDs[12:16], 0)   // No target
		binary.LittleEndian.PutUint32(targetIDs[16:20], 0)   // No target

		// Create packet data
		args := map[string]interface{}{
			"switch":    "01CF",
			"sourceID":  uint32(2001), // Source ID
			"targetIDs": targetIDs,
			"range":     uint8(5), // Range
		}

		// Call the handler
		err := manager.handleDevotion(args)
		if err != nil {
			t.Errorf("handleDevotion() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the source ID
		if sourceID, ok := hookResult["sourceID"].(uint32); !ok || sourceID != 2001 {
			t.Errorf("Expected source ID 2001, got %v", sourceID)
		}

		// Check the range
		if devotionRange, ok := hookResult["range"].(uint8); !ok || devotionRange != 5 {
			t.Errorf("Expected range 5, got %v", devotionRange)
		}

		// Check the targets
		if targets, ok := hookResult["targets"].([]uint32); ok {
			if len(targets) != 3 {
				t.Errorf("Expected 3 targets, got %d", len(targets))
			} else {
				if targets[0] != 1001 {
					t.Errorf("Expected target 1 to be 1001, got %d", targets[0])
				}
				if targets[1] != 1002 {
					t.Errorf("Expected target 2 to be 1002, got %d", targets[1])
				}
				if targets[2] != 1003 {
					t.Errorf("Expected target 3 to be 1003, got %d", targets[2])
				}
			}
		} else {
			t.Error("targets not found in hook result or not a []uint32")
		}

		// Check the target indices
		if targetIndices, ok := hookResult["targetIndices"].(map[uint32]int); ok {
			if len(targetIndices) != 3 {
				t.Errorf("Expected 3 target indices, got %d", len(targetIndices))
			} else {
				if targetIndices[1001] != 0 {
					t.Errorf("Expected target index for 1001 to be 0, got %d", targetIndices[1001])
				}
				if targetIndices[1002] != 1 {
					t.Errorf("Expected target index for 1002 to be 1, got %d", targetIndices[1002])
				}
				if targetIndices[1003] != 2 {
					t.Errorf("Expected target index for 1003 to be 2, got %d", targetIndices[1003])
				}
			}
		} else {
			t.Error("targetIndices not found in hook result or not a map[uint32]int")
		}

		// Check that the message is not empty
		if message, ok := hookResult["message"].(string); !ok || message == "" {
			t.Error("message not found in hook result or empty")
		}
	})
}

// Test unhappy paths
func TestDevotionUnhappy(t *testing.T) {
	// Create a devotion manager
	manager := NewDevotionManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01CF",
			// Missing sourceID, targetIDs, and range
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleDevotion(args)
		if err != nil {
			t.Errorf("handleDevotion() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processDevotion(args)

		// Check that the source ID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected source ID 0, got %v", sourceID)
		}

		// Check that the range is zero
		if devotionRange, ok := result["range"].(uint8); !ok || devotionRange != 0 {
			t.Errorf("Expected range 0, got %v", devotionRange)
		}

		// Check that the targets is empty
		if targets, ok := result["targets"].([]uint32); !ok || len(targets) != 0 {
			t.Errorf("Expected empty targets, got %v", targets)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":    "01CF",
			"sourceID":  "not a uint32", // Wrong type
			"targetIDs": "not a []byte", // Wrong type
			"range":     "not a uint8",  // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleDevotion(args)
		if err != nil {
			t.Errorf("handleDevotion() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processDevotion(args)

		// Check that the source ID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected source ID 0, got %v", sourceID)
		}

		// Check that the range is zero
		if devotionRange, ok := result["range"].(uint8); !ok || devotionRange != 0 {
			t.Errorf("Expected range 0, got %v", devotionRange)
		}

		// Check that the targets is empty
		if targets, ok := result["targets"].([]uint32); !ok || len(targets) != 0 {
			t.Errorf("Expected empty targets, got %v", targets)
		}
	})

	// Test with invalid target IDs
	t.Run("InvalidTargetIDs", func(t *testing.T) {
		// Create target IDs with invalid data (too short)
		targetIDs := make([]byte, 2) // Not enough data for even one target ID

		// Create packet data
		args := map[string]interface{}{
			"switch":    "01CF",
			"sourceID":  uint32(2001),
			"targetIDs": targetIDs,
			"range":     uint8(5),
		}

		// This should not return an error, but the targets should be empty
		err := manager.handleDevotion(args)
		if err != nil {
			t.Errorf("handleDevotion() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processDevotion(args)

		// Check that the targets is empty
		if targets, ok := result["targets"].([]uint32); !ok || len(targets) != 0 {
			t.Errorf("Expected empty targets, got %v", targets)
		}
	})
}
