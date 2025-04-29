package core

import (
	"testing"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMapChangeCell(t *testing.T) {
	// Test case 1: Happy path - valid map cell change
	t.Run("ValidMapCellChange", func(t *testing.T) {
		parser := NewCoreParser("ServerType0", hooks.NewHookManager())
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)
		manager.session.MapName = "prontera"

		// Create test packet arguments
		args := map[string]interface{}{
			"x":        uint16(100),
			"y":        uint16(150),
			"type":     byte(1),    // Some cell type
			"map_name": "prontera", // Should match current map
		}

		// Call handler
		err := manager.handleMapChangeCell(args)
		if err != nil {
			t.Fatalf("handleMapChangeCell() returned error: %v", err)
		}

		// Check that session state remains unchanged
		if manager.GetNetworkState() != network.InGame {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
		}

		session := manager.GetSession()
		if session.State != AccountStateInGame {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
		}
		if session.MapName != "prontera" {
			t.Errorf("session.MapName = %v, want prontera", session.MapName)
		}
	})

	// Test case 2: Happy path - map cell change on different map
	t.Run("MapCellChangeOnDifferentMap", func(t *testing.T) {
		parser := NewCoreParser("ServerType0", hooks.NewHookManager())
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)
		manager.session.MapName = "prontera"

		// Create test packet arguments
		args := map[string]interface{}{
			"x":        uint16(200),
			"y":        uint16(250),
			"type":     byte(2), // Some cell type
			"map_name": "payon", // Different map
		}

		// Call handler
		err := manager.handleMapChangeCell(args)
		if err != nil {
			t.Fatalf("handleMapChangeCell() returned error: %v", err)
		}

		// Check that session state remains unchanged
		if manager.GetNetworkState() != network.InGame {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
		}

		session := manager.GetSession()
		if session.State != AccountStateInGame {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
		}
		// Map name should not change
		if session.MapName != "prontera" {
			t.Errorf("session.MapName = %v, want prontera", session.MapName)
		}
	})

	// Test case 3: Unhappy path - missing coordinates
	t.Run("MissingCoordinates", func(t *testing.T) {
		parser := NewCoreParser("ServerType0", hooks.NewHookManager())
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments with missing x,y
		args := map[string]interface{}{
			"type":     byte(1),
			"map_name": "prontera",
		}

		// Call handler
		err := manager.handleMapChangeCell(args)
		if err == nil {
			t.Fatalf("handleMapChangeCell() should return error for missing coordinates")
		}
	})

	// Test case 4: Unhappy path - missing map name
	t.Run("MissingMapName", func(t *testing.T) {
		parser := NewCoreParser("ServerType0", hooks.NewHookManager())
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments with missing map_name
		args := map[string]interface{}{
			"x":    uint16(100),
			"y":    uint16(150),
			"type": byte(1),
		}

		// Call handler
		err := manager.handleMapChangeCell(args)
		if err == nil {
			t.Fatalf("handleMapChangeCell() should return error for missing map name")
		}
	})
}
