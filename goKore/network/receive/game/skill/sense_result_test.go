package skill

import (
	"strings"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSenseResult(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the sense result manager
	manager := NewSenseResultManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for sense result
	t.Run("Basic Sense Result", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.sense_result", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch":  "0213",
			"nameID":  uint16(1002),
			"level":   uint16(50),
			"size":    uint8(1), // Medium
			"race":    uint8(2), // Beast
			"def":     uint16(100),
			"mdef":    uint16(50),
			"element": uint8(3), // Fire
			"hp":      uint16(5000),
			"ice":     uint8(150),
			"earth":   uint8(100),
			"fire":    uint8(50),
			"wind":    uint8(75),
			"poison":  uint8(25),
			"holy":    uint8(200),
			"dark":    uint8(0),
			"spirit":  uint8(100),
			"undead":  uint8(150),
		}

		// Call the handler
		err := manager.handleSenseResult(args)
		if err != nil {
			t.Errorf("handleSenseResult() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the nameID
		if nameID, ok := hookResult["nameID"].(uint16); !ok || nameID != 1002 {
			t.Errorf("Expected nameID 1002, got %v", nameID)
		}

		// Check the level
		if level, ok := hookResult["level"].(uint16); !ok || level != 50 {
			t.Errorf("Expected level 50, got %v", level)
		}

		// Check the size
		if size, ok := hookResult["size"].(uint8); !ok || size != 1 {
			t.Errorf("Expected size 1, got %v", size)
		}

		// Check the size name
		if sizeName, ok := hookResult["sizeName"].(string); !ok || sizeName != "Medium" {
			t.Errorf("Expected size name 'Medium', got %q", sizeName)
		}

		// Check the race
		if race, ok := hookResult["race"].(uint8); !ok || race != 2 {
			t.Errorf("Expected race 2, got %v", race)
		}

		// Check the race name
		if raceName, ok := hookResult["raceName"].(string); !ok || raceName != "Beast" {
			t.Errorf("Expected race name 'Beast', got %q", raceName)
		}

		// Check the def
		if def, ok := hookResult["def"].(uint16); !ok || def != 100 {
			t.Errorf("Expected def 100, got %v", def)
		}

		// Check the mdef
		if mdef, ok := hookResult["mdef"].(uint16); !ok || mdef != 50 {
			t.Errorf("Expected mdef 50, got %v", mdef)
		}

		// Check the element
		if element, ok := hookResult["element"].(uint8); !ok || element != 3 {
			t.Errorf("Expected element 3, got %v", element)
		}

		// Check the element name
		if elementName, ok := hookResult["elementName"].(string); !ok || elementName != "Fire" {
			t.Errorf("Expected element name 'Fire', got %q", elementName)
		}

		// Check the hp
		if hp, ok := hookResult["hp"].(uint16); !ok || hp != 5000 {
			t.Errorf("Expected hp 5000, got %v", hp)
		}

		// Check the damage modifiers
		if ice, ok := hookResult["ice"].(uint8); !ok || ice != 150 {
			t.Errorf("Expected ice 150, got %v", ice)
		}

		if earth, ok := hookResult["earth"].(uint8); !ok || earth != 100 {
			t.Errorf("Expected earth 100, got %v", earth)
		}

		if fire, ok := hookResult["fire"].(uint8); !ok || fire != 50 {
			t.Errorf("Expected fire 50, got %v", fire)
		}

		if wind, ok := hookResult["wind"].(uint8); !ok || wind != 75 {
			t.Errorf("Expected wind 75, got %v", wind)
		}

		if poison, ok := hookResult["poison"].(uint8); !ok || poison != 25 {
			t.Errorf("Expected poison 25, got %v", poison)
		}

		if holy, ok := hookResult["holy"].(uint8); !ok || holy != 200 {
			t.Errorf("Expected holy 200, got %v", holy)
		}

		if dark, ok := hookResult["dark"].(uint8); !ok || dark != 0 {
			t.Errorf("Expected dark 0, got %v", dark)
		}

		if spirit, ok := hookResult["spirit"].(uint8); !ok || spirit != 100 {
			t.Errorf("Expected spirit 100, got %v", spirit)
		}

		if undead, ok := hookResult["undead"].(uint8); !ok || undead != 150 {
			t.Errorf("Expected undead 150, got %v", undead)
		}

		// Check the message
		if message, ok := hookResult["message"].(string); !ok {
			t.Error("Message not found in hook result")
		} else {
			// Check that the message contains the expected information
			expectedStrings := []string{
				"Monster: Monster_1002",
				"Level: 50",
				"Size: Medium",
				"Race: Beast",
				"Def: 100",
				"MDef: 50",
				"Element: Fire",
				"HP: 5000",
				"Ice: 150",
				"Earth: 100",
				"Fire: 50",
				"Wind: 75",
				"Poison: 25",
				"Holy: 200",
				"Dark: 0",
				"Spirit: 100",
				"Undead: 150",
			}

			for _, expected := range expectedStrings {
				if !strings.Contains(message, expected) {
					t.Errorf("Expected message to contain %q, but it doesn't", expected)
				}
			}
		}
	})
}

// Test unhappy paths
func TestSenseResultUnhappy(t *testing.T) {
	// Create a sense result manager
	manager := NewSenseResultManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0213",
			// Missing all other fields
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSenseResult(args)
		if err != nil {
			t.Errorf("handleSenseResult() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSenseResult(args)

		// Check that the nameID is zero
		if nameID, ok := result["nameID"].(uint16); !ok || nameID != 0 {
			t.Errorf("Expected nameID 0, got %v", nameID)
		}

		// Check that the level is zero
		if level, ok := result["level"].(uint16); !ok || level != 0 {
			t.Errorf("Expected level 0, got %v", level)
		}

		// Check that the size is zero
		if size, ok := result["size"].(uint8); !ok || size != 0 {
			t.Errorf("Expected size 0, got %v", size)
		}

		// Check that the race is zero
		if race, ok := result["race"].(uint8); !ok || race != 0 {
			t.Errorf("Expected race 0, got %v", race)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "0213",
			"nameID":  "not a uint16", // Wrong type
			"level":   "not a uint16", // Wrong type
			"size":    "not a uint8",  // Wrong type
			"race":    "not a uint8",  // Wrong type
			"def":     "not a uint16", // Wrong type
			"mdef":    "not a uint16", // Wrong type
			"element": "not a uint8",  // Wrong type
			"hp":      "not a uint16", // Wrong type
			"ice":     "not a uint8",  // Wrong type
			"earth":   "not a uint8",  // Wrong type
			"fire":    "not a uint8",  // Wrong type
			"wind":    "not a uint8",  // Wrong type
			"poison":  "not a uint8",  // Wrong type
			"holy":    "not a uint8",  // Wrong type
			"dark":    "not a uint8",  // Wrong type
			"spirit":  "not a uint8",  // Wrong type
			"undead":  "not a uint8",  // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSenseResult(args)
		if err != nil {
			t.Errorf("handleSenseResult() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSenseResult(args)

		// Check that the nameID is zero
		if nameID, ok := result["nameID"].(uint16); !ok || nameID != 0 {
			t.Errorf("Expected nameID 0, got %v", nameID)
		}

		// Check that the level is zero
		if level, ok := result["level"].(uint16); !ok || level != 0 {
			t.Errorf("Expected level 0, got %v", level)
		}

		// Check that the size is zero
		if size, ok := result["size"].(uint8); !ok || size != 0 {
			t.Errorf("Expected size 0, got %v", size)
		}

		// Check that the race is zero
		if race, ok := result["race"].(uint8); !ok || race != 0 {
			t.Errorf("Expected race 0, got %v", race)
		}
	})
}
