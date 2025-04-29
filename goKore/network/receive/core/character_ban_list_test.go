package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleCharacterBanList(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a character manager
	parser := NewCoreParser("ServerType0", hookManager)
	charManager := NewCharacterManager(parser)

	// Register a hook to capture the ban list
	var capturedBanList []string
	hookManager.AddHook("character.ban_list", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		capturedBanList = result["ban_list"].([]string)
	}, nil)

	// Test case 1: Happy path - Valid ban list with multiple entries
	t.Run("Valid ban list with multiple entries", func(t *testing.T) {
		// Reset captured ban list
		capturedBanList = nil

		// Create mock packet data
		// Format: Header + Len + CharList[character_name(size:24)]
		// Let's say we have 2 banned characters: "TestChar1" and "TestChar2"
		mockData := make([]byte, 0)
		mockData = append(mockData, 0x02) // Number of entries

		// Add first character name (padded to 24 bytes)
		char1 := "TestChar1"
		char1Bytes := make([]byte, 24)
		copy(char1Bytes, []byte(char1))
		mockData = append(mockData, char1Bytes...)

		// Add second character name (padded to 24 bytes)
		char2 := "TestChar2"
		char2Bytes := make([]byte, 24)
		copy(char2Bytes, []byte(char2))
		mockData = append(mockData, char2Bytes...)

		// Call the handler
		err := charManager.handleCharacterBanList(map[string]interface{}{
			"charList": mockData,
		})

		// Verify results
		if err != nil {
			t.Fatalf("handleCharacterBanList() returned error: %v", err)
		}

		if len(capturedBanList) != 2 {
			t.Errorf("len(capturedBanList) = %d, want 2", len(capturedBanList))
		}

		if capturedBanList[0] != "TestChar1" {
			t.Errorf("capturedBanList[0] = %s, want TestChar1", capturedBanList[0])
		}

		if capturedBanList[1] != "TestChar2" {
			t.Errorf("capturedBanList[1] = %s, want TestChar2", capturedBanList[1])
		}
	})

	// Test case 2: Happy path - Empty ban list
	t.Run("Empty ban list", func(t *testing.T) {
		// Reset captured ban list
		capturedBanList = nil

		// Create mock packet data with 0 entries
		mockData := []byte{0x00}

		// Call the handler
		err := charManager.handleCharacterBanList(map[string]interface{}{
			"charList": mockData,
		})

		// Verify results
		if err != nil {
			t.Fatalf("handleCharacterBanList() returned error: %v", err)
		}

		if len(capturedBanList) != 0 {
			t.Errorf("len(capturedBanList) = %d, want 0", len(capturedBanList))
		}
	})

	// Test case 3: Unhappy path - Invalid data (nil)
	t.Run("Invalid data (nil)", func(t *testing.T) {
		// Reset captured ban list
		capturedBanList = nil

		// Call the handler with nil data
		err := charManager.handleCharacterBanList(map[string]interface{}{
			"charList": nil,
		})

		// Verify results
		if err == nil {
			t.Fatal("handleCharacterBanList() did not return error for nil data")
		}
	})

	// Test case 4: Unhappy path - Invalid data (truncated)
	t.Run("Invalid data (truncated)", func(t *testing.T) {
		// Reset captured ban list
		capturedBanList = nil

		// Create truncated mock data (says 1 entry but doesn't have enough bytes)
		mockData := []byte{0x01, 0x01, 0x02} // Should have 24 bytes for the character name

		// Call the handler
		err := charManager.handleCharacterBanList(map[string]interface{}{
			"charList": mockData,
		})

		// Verify results
		if err == nil {
			t.Fatal("handleCharacterBanList() did not return error for truncated data")
		}
	})
}
