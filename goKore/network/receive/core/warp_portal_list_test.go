package core

import (
	"testing"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleWarpPortalList(t *testing.T) {
	// Test case 1: Happy path - Teleport skill (type 26)
	t.Run("TeleportSkill", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		warpListCalled := false
		var warpType byte
		var memoList []string
		hookManager.RegisterHook("warp.portal_list", func(args interface{}) {
			warpListCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if t, ok := hookArgs["type"].(byte); ok {
					warpType = t
				}
				if list, ok := hookArgs["memo_list"].([]string); ok {
					memoList = list
				}
			}
		})

		// Create a channel to capture config updates
		configUpdateCalled := false
		var configKey, configValue string
		hookManager.RegisterHook("config.update", func(args interface{}) {
			configUpdateCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if k, ok := hookArgs["key"].(string); ok {
					configKey = k
				}
				if v, ok := hookArgs["value"].(string); ok {
					configValue = v
				}
			}
		})

		// Create test packet arguments
		args := map[string]interface{}{
			"type":  byte(26), // Teleport skill
			"memo1": "prontera.gat",
			"memo2": "geffen.gat",
			"memo3": "morocc.gat",
			"memo4": "alberta.gat",
		}

		// Call handler
		err := manager.handleWarpPortalList(args)
		if err != nil {
			t.Fatalf("handleWarpPortalList() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !warpListCalled {
			t.Error("warp.portal_list hook was not called")
		}
		if warpType != 26 {
			t.Errorf("warp type = %d, want 26", warpType)
		}
		if len(memoList) != 4 {
			t.Errorf("memo list length = %d, want 4", len(memoList))
		}
		expectedMemos := []string{"prontera", "geffen", "morocc", "alberta"}
		for i, memo := range expectedMemos {
			if i < len(memoList) && memoList[i] != memo {
				t.Errorf("memo[%d] = %q, want %q", i, memoList[i], memo)
			}
		}

		// Check that config was updated
		if !configUpdateCalled {
			t.Error("config.update hook was not called")
		}
		if configKey != "saveMap" {
			t.Errorf("config key = %q, want %q", configKey, "saveMap")
		}
		if configValue != "geffen" {
			t.Errorf("config value = %q, want %q", configValue, "geffen")
		}
	})

	// Test case 2: Happy path - Butterfly Wing (type 27)
	t.Run("ButterflyWing", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		warpListCalled := false
		var warpType byte
		var memoList []string
		hookManager.RegisterHook("warp.portal_list", func(args interface{}) {
			warpListCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if t, ok := hookArgs["type"].(byte); ok {
					warpType = t
				}
				if list, ok := hookArgs["memo_list"].([]string); ok {
					memoList = list
				}
			}
		})

		// Create a channel to capture config updates
		configUpdateCalled := false
		var configKey, configValue string
		hookManager.RegisterHook("config.update", func(args interface{}) {
			configUpdateCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if k, ok := hookArgs["key"].(string); ok {
					configKey = k
				}
				if v, ok := hookArgs["value"].(string); ok {
					configValue = v
				}
			}
		})

		// Create test packet arguments
		args := map[string]interface{}{
			"type":  byte(27), // Butterfly Wing
			"memo1": "prontera.gat",
			"memo2": "geffen.gat",
			"memo3": "morocc.gat",
			"memo4": "alberta.gat",
		}

		// Call handler
		err := manager.handleWarpPortalList(args)
		if err != nil {
			t.Fatalf("handleWarpPortalList() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !warpListCalled {
			t.Error("warp.portal_list hook was not called")
		}
		if warpType != 27 {
			t.Errorf("warp type = %d, want 27", warpType)
		}
		if len(memoList) != 4 {
			t.Errorf("memo list length = %d, want 4", len(memoList))
		}
		expectedMemos := []string{"prontera", "geffen", "morocc", "alberta"}
		for i, memo := range expectedMemos {
			if i < len(memoList) && memoList[i] != memo {
				t.Errorf("memo[%d] = %q, want %q", i, memoList[i], memo)
			}
		}

		// Check that config was updated
		if !configUpdateCalled {
			t.Error("config.update hook was not called")
		}
		if configKey != "saveMap" {
			t.Errorf("config key = %q, want %q", configKey, "saveMap")
		}
		if configValue != "prontera" {
			t.Errorf("config value = %q, want %q", configValue, "prontera")
		}
	})

	// Test case 3: Empty memo fields
	t.Run("EmptyMemoFields", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		warpListCalled := false
		var memoList []string
		hookManager.RegisterHook("warp.portal_list", func(args interface{}) {
			warpListCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if list, ok := hookArgs["memo_list"].([]string); ok {
					memoList = list
				}
			}
		})

		// Create test packet arguments with empty memo fields
		args := map[string]interface{}{
			"type":  byte(26),
			"memo1": "prontera.gat",
			"memo2": "",
			"memo3": "morocc.gat",
			"memo4": "",
		}

		// Call handler
		err := manager.handleWarpPortalList(args)
		if err != nil {
			t.Fatalf("handleWarpPortalList() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !warpListCalled {
			t.Error("warp.portal_list hook was not called")
		}
		if len(memoList) != 2 {
			t.Errorf("memo list length = %d, want 2", len(memoList))
		}
		expectedMemos := []string{"prontera", "morocc"}
		for i, memo := range expectedMemos {
			if i < len(memoList) && memoList[i] != memo {
				t.Errorf("memo[%d] = %q, want %q", i, memoList[i], memo)
			}
		}
	})

	// Test case 4: Unhappy path - Missing type
	t.Run("MissingType", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		manager := NewAccountManager(parser)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments with missing type
		args := map[string]interface{}{
			"memo1": "prontera.gat",
			"memo2": "geffen.gat",
			"memo3": "morocc.gat",
			"memo4": "alberta.gat",
		}

		// Call handler
		err := manager.handleWarpPortalList(args)
		if err == nil {
			t.Fatalf("handleWarpPortalList() should return error for missing type")
		}
	})
}
