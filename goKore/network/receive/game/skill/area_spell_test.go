package skill

import (
	"fmt"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestAreaSpell(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the area spell manager
	manager := NewAreaSpellManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different area spell types
	testCases := []struct {
		name            string
		switchType      string
		id              uint32
		sourceID        uint32
		x               uint16
		y               uint16
		spellType       uint16
		isVisible       uint8
		scribbleMsg     string
		expectedBinID   int
		expectedMessage string
	}{
		{
			name:            "Basic Area Spell (011F)",
			switchType:      "011F",
			id:              1001,
			sourceID:        2001,
			x:               150,
			y:               100,
			spellType:       123,
			isVisible:       1,
			scribbleMsg:     "",
			expectedBinID:   0,
			expectedMessage: "Area effect Spell_123 from Actor_2001 appeared on (150, 100)",
		},
		{
			name:            "Warp Portal (011F)",
			switchType:      "011F",
			id:              1002,
			sourceID:        2002,
			x:               200,
			y:               150,
			spellType:       0x81,
			isVisible:       1,
			scribbleMsg:     "",
			expectedBinID:   1,
			expectedMessage: "Actor_2002 opened Warp Portal on (200, 150)",
		},
		{
			name:            "Scribble Message (01C9)",
			switchType:      "01C9",
			id:              1003,
			sourceID:        2003,
			x:               250,
			y:               200,
			spellType:       456,
			isVisible:       1,
			scribbleMsg:     "Hello World!",
			expectedBinID:   2,
			expectedMessage: "Actor_2003 has scribbled: Hello World! on (250, 200)",
		},
		{
			name:            "Expanded Area Spell (08C7)",
			switchType:      "08C7",
			id:              1004,
			sourceID:        2004,
			x:               300,
			y:               250,
			spellType:       789,
			isVisible:       0,
			scribbleMsg:     "",
			expectedBinID:   3,
			expectedMessage: "Area effect Spell_789 from Actor_2004 appeared on (300, 250)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.area_spell", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch":    tc.switchType,
				"ID":        tc.id,
				"sourceID":  tc.sourceID,
				"x":         tc.x,
				"y":         tc.y,
				"type":      tc.spellType,
				"isVisible": tc.isVisible,
			}

			// Add scribble message if needed
			if tc.switchType == "01C9" {
				args["scribbleMsg"] = tc.scribbleMsg
			}

			// Call the handler
			err := manager.handleAreaSpell(args)
			if err != nil {
				t.Errorf("handleAreaSpell() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the ID
			if id, ok := hookResult["ID"].(uint32); !ok || id != tc.id {
				t.Errorf("Expected ID %d, got %v", tc.id, id)
			}

			// Check the sourceID
			if sourceID, ok := hookResult["sourceID"].(uint32); !ok || sourceID != tc.sourceID {
				t.Errorf("Expected sourceID %d, got %v", tc.sourceID, sourceID)
			}

			// Check the x coordinate
			if x, ok := hookResult["x"].(uint16); !ok || x != tc.x {
				t.Errorf("Expected x %d, got %v", tc.x, x)
			}

			// Check the y coordinate
			if y, ok := hookResult["y"].(uint16); !ok || y != tc.y {
				t.Errorf("Expected y %d, got %v", tc.y, y)
			}

			// Check the spell type
			if spellType, ok := hookResult["type"].(uint16); !ok || spellType != tc.spellType {
				t.Errorf("Expected type %d, got %v", tc.spellType, spellType)
			}

			// Check the isVisible flag
			if isVisible, ok := hookResult["isVisible"].(uint8); !ok || isVisible != tc.isVisible {
				t.Errorf("Expected isVisible %d, got %v", tc.isVisible, isVisible)
			}

			// Check the binID
			if binID, ok := hookResult["binID"].(int); !ok || binID != tc.expectedBinID {
				t.Errorf("Expected binID %d, got %v", tc.expectedBinID, binID)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}

			// Check the scribble message if applicable
			if tc.switchType == "01C9" {
				if scribbleMsg, ok := hookResult["scribbleMsg"].(string); !ok || scribbleMsg != tc.scribbleMsg {
					t.Errorf("Expected scribbleMsg %q, got %q", tc.scribbleMsg, scribbleMsg)
				}
			}
		})
	}
}

// Test unhappy paths
func TestAreaSpellUnhappy(t *testing.T) {
	// Create an area spell manager
	manager := NewAreaSpellManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "011F",
			// Missing ID, sourceID, x, y, type, and isVisible
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleAreaSpell(args)
		if err != nil {
			t.Errorf("handleAreaSpell() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processAreaSpell(args)

		// Check that the ID is zero
		if id, ok := result["ID"].(uint32); !ok || id != 0 {
			t.Errorf("Expected ID 0, got %v", id)
		}

		// Check that the sourceID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected sourceID 0, got %v", sourceID)
		}

		// Check that the x coordinate is zero
		if x, ok := result["x"].(uint16); !ok || x != 0 {
			t.Errorf("Expected x 0, got %v", x)
		}

		// Check that the y coordinate is zero
		if y, ok := result["y"].(uint16); !ok || y != 0 {
			t.Errorf("Expected y 0, got %v", y)
		}

		// Check that the spell type is zero
		if spellType, ok := result["type"].(uint16); !ok || spellType != 0 {
			t.Errorf("Expected type 0, got %v", spellType)
		}

		// Check that the isVisible flag is zero
		if isVisible, ok := result["isVisible"].(uint8); !ok || isVisible != 0 {
			t.Errorf("Expected isVisible 0, got %v", isVisible)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":    "011F",
			"ID":        "not a uint32", // Wrong type
			"sourceID":  "not a uint32", // Wrong type
			"x":         "not a uint16", // Wrong type
			"y":         "not a uint16", // Wrong type
			"type":      "not a uint16", // Wrong type
			"isVisible": "not a uint8",  // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleAreaSpell(args)
		if err != nil {
			t.Errorf("handleAreaSpell() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processAreaSpell(args)

		// Check that the ID is zero
		if id, ok := result["ID"].(uint32); !ok || id != 0 {
			t.Errorf("Expected ID 0, got %v", id)
		}

		// Check that the sourceID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected sourceID 0, got %v", sourceID)
		}

		// Check that the x coordinate is zero
		if x, ok := result["x"].(uint16); !ok || x != 0 {
			t.Errorf("Expected x 0, got %v", x)
		}

		// Check that the y coordinate is zero
		if y, ok := result["y"].(uint16); !ok || y != 0 {
			t.Errorf("Expected y 0, got %v", y)
		}

		// Check that the spell type is zero
		if spellType, ok := result["type"].(uint16); !ok || spellType != 0 {
			t.Errorf("Expected type 0, got %v", spellType)
		}

		// Check that the isVisible flag is zero
		if isVisible, ok := result["isVisible"].(uint8); !ok || isVisible != 0 {
			t.Errorf("Expected isVisible 0, got %v", isVisible)
		}
	})
}

// TestAreaSpellDisappears tests the area_spell_disappears handler
func TestAreaSpellDisappears(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the area spell manager
	manager := NewAreaSpellManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// First, add a spell to the manager's spell map
	spellID := uint32(1001)
	sourceID := uint32(2001)
	x := uint16(150)
	y := uint16(100)
	spellType := uint16(123)
	isVisible := uint8(1)
	binID := 0

	// Create the spell
	manager.spells[spellID] = &AreaSpell{
		ID:        spellID,
		SourceID:  sourceID,
		X:         x,
		Y:         y,
		Type:      spellType,
		IsVisible: isVisible,
		BinID:     binID,
	}

	// Test case for area spell disappears
	t.Run("Basic Area Spell Disappears", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.area_spell_disappears", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch": "0120",
			"ID":     spellID,
		}

		// Call the handler
		err := manager.handleAreaSpellDisappears(args)
		if err != nil {
			t.Errorf("handleAreaSpellDisappears() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the ID
		if id, ok := hookResult["ID"].(uint32); !ok || id != spellID {
			t.Errorf("Expected ID %d, got %v", spellID, id)
		}

		// Check the sourceID
		if sid, ok := hookResult["sourceID"].(uint32); !ok || sid != sourceID {
			t.Errorf("Expected sourceID %d, got %v", sourceID, sid)
		}

		// Check the x coordinate
		if xCoord, ok := hookResult["x"].(uint16); !ok || xCoord != x {
			t.Errorf("Expected x %d, got %v", x, xCoord)
		}

		// Check the y coordinate
		if yCoord, ok := hookResult["y"].(uint16); !ok || yCoord != y {
			t.Errorf("Expected y %d, got %v", y, yCoord)
		}

		// Check the spell type
		if sType, ok := hookResult["type"].(uint16); !ok || sType != spellType {
			t.Errorf("Expected type %d, got %v", spellType, sType)
		}

		// Check the binID
		if bid, ok := hookResult["binID"].(int); !ok || bid != binID {
			t.Errorf("Expected binID %d, got %v", binID, bid)
		}

		// Check the message
		expectedMessage := fmt.Sprintf("Area effect Spell_%d from Actor_%d disappeared from (%d, %d)",
			spellType, sourceID, x, y)
		if message, ok := hookResult["message"].(string); !ok || message != expectedMessage {
			t.Errorf("Expected message %q, got %q", expectedMessage, message)
		}

		// Check that the spell was removed from the spells map
		if _, exists := manager.spells[spellID]; exists {
			t.Error("Spell was not removed from the spells map")
		}
	})

	// Test case for area spell disappears with unknown spell ID
	t.Run("Unknown Spell ID", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.area_spell_disappears", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data with an unknown spell ID
		unknownSpellID := uint32(9999)
		args := map[string]interface{}{
			"switch": "0120",
			"ID":     unknownSpellID,
		}

		// Call the handler
		err := manager.handleAreaSpellDisappears(args)
		if err != nil {
			t.Errorf("handleAreaSpellDisappears() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the ID
		if id, ok := hookResult["ID"].(uint32); !ok || id != unknownSpellID {
			t.Errorf("Expected ID %d, got %v", unknownSpellID, id)
		}

		// Check the message
		expectedMessage := fmt.Sprintf("Area spell %d disappeared (not found in spell list)", unknownSpellID)
		if message, ok := hookResult["message"].(string); !ok || message != expectedMessage {
			t.Errorf("Expected message %q, got %q", expectedMessage, message)
		}
	})
}
