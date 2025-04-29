package core

import (
	"testing"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMapProperty3(t *testing.T) {
	// Test case 1: Happy path - PvP map type
	t.Run("PvPMapType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		pvpModeCalled := false
		pvpModeType := 0
		hookManager.RegisterHook("pvp_mode", func(args interface{}) {
			pvpModeCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if pvp, ok := hookArgs["pvp"].(int); ok {
					pvpModeType = pvp
				}
			}
		})

		// Create test packet arguments for PvP map (type 6)
		args := map[string]interface{}{
			"type":       byte(6),            // PvP map type
			"info_table": []byte{1, 0, 0, 0}, // Some info table data
		}

		// Call handler
		err := manager.handleMapProperty3(args)
		if err != nil {
			t.Fatalf("handleMapProperty3() returned error: %v", err)
		}

		// Check that PvP mode was set correctly
		if !pvpModeCalled {
			t.Error("PvP mode hook was not called")
		}
		if pvpModeType != 1 {
			t.Errorf("PvP mode = %d, want 1", pvpModeType)
		}
	})

	// Test case 2: Happy path - GvG map type
	t.Run("GvGMapType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		pvpModeCalled := false
		pvpModeType := 0
		hookManager.RegisterHook("pvp_mode", func(args interface{}) {
			pvpModeCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if pvp, ok := hookArgs["pvp"].(int); ok {
					pvpModeType = pvp
				}
			}
		})

		// Create test packet arguments for GvG map (type 8)
		args := map[string]interface{}{
			"type":       byte(8),            // GvG map type
			"info_table": []byte{1, 0, 0, 0}, // Some info table data
		}

		// Call handler
		err := manager.handleMapProperty3(args)
		if err != nil {
			t.Fatalf("handleMapProperty3() returned error: %v", err)
		}

		// Check that PvP mode was set correctly
		if !pvpModeCalled {
			t.Error("PvP mode hook was not called")
		}
		if pvpModeType != 2 {
			t.Errorf("PvP mode = %d, want 2", pvpModeType)
		}
	})

	// Test case 3: Happy path - Battleground map type
	t.Run("BattlegroundMapType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		pvpModeCalled := false
		pvpModeType := 0
		hookManager.RegisterHook("pvp_mode", func(args interface{}) {
			pvpModeCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if pvp, ok := hookArgs["pvp"].(int); ok {
					pvpModeType = pvp
				}
			}
		})

		// Create test packet arguments for Battleground map (type 19)
		args := map[string]interface{}{
			"type":       byte(19),           // Battleground map type
			"info_table": []byte{1, 0, 0, 0}, // Some info table data
		}

		// Call handler
		err := manager.handleMapProperty3(args)
		if err != nil {
			t.Fatalf("handleMapProperty3() returned error: %v", err)
		}

		// Check that PvP mode was set correctly
		if !pvpModeCalled {
			t.Error("PvP mode hook was not called")
		}
		if pvpModeType != 3 {
			t.Errorf("PvP mode = %d, want 3", pvpModeType)
		}
	})

	// Test case 4: Happy path - Normal map type (non-PvP)
	t.Run("NormalMapType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		pvpModeCalled := false
		hookManager.RegisterHook("pvp_mode", func(args interface{}) {
			pvpModeCalled = true
		})

		// Create test packet arguments for normal map (type 0)
		args := map[string]interface{}{
			"type":       byte(0),            // Normal map type
			"info_table": []byte{1, 0, 0, 0}, // Some info table data
		}

		// Call handler
		err := manager.handleMapProperty3(args)
		if err != nil {
			t.Fatalf("handleMapProperty3() returned error: %v", err)
		}

		// Check that PvP mode hook was not called
		if pvpModeCalled {
			t.Error("PvP mode hook was called for non-PvP map")
		}
	})

	// Test case 5: Unhappy path - Missing type
	t.Run("MissingType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments with missing type
		args := map[string]interface{}{
			"info_table": []byte{1, 0, 0, 0}, // Some info table data
		}

		// Call handler
		err := manager.handleMapProperty3(args)
		if err == nil {
			t.Fatalf("handleMapProperty3() should return error for missing type")
		}
	})
}
