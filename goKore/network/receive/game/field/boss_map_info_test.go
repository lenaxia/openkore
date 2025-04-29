package field

import (
	"fmt"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleBossMapInfo(t *testing.T) {
	// Test case 1: No boss monster found (flag 0)
	t.Run("NoBossFound", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create a channel to capture hook calls
		resultChan := make(chan map[string]interface{}, 1)

		// Register a hook to capture the result
		hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
			result := arg.(map[string]interface{})
			resultChan <- result
		}, nil)

		// Create test packet arguments
		args := map[string]interface{}{
			"flag": byte(0),
			"name": []byte("Baphomet"),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err != nil {
			t.Fatalf("handleBossMapInfo returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(0) {
			t.Errorf("Expected flag 0, got %v", result["flag"])
		}
		if result["status"] != "You cannot find any trace of a Boss Monster in this area." {
			t.Errorf("Expected status about no boss, got %v", result["status"])
		}
	})

	// Test case 2: Boss monster found at specific coordinates (flag 1)
	t.Run("BossFoundAtCoordinates", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create a channel to capture hook calls
		resultChan := make(chan map[string]interface{}, 1)

		// Register a hook to capture the result
		hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
			result := arg.(map[string]interface{})
			resultChan <- result
		}, nil)

		// Create test packet arguments
		args := map[string]interface{}{
			"flag": byte(1),
			"name": []byte("Baphomet"),
			"x":    uint16(150),
			"y":    uint16(200),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err != nil {
			t.Fatalf("handleBossMapInfo returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(1) {
			t.Errorf("Expected flag 1, got %v", result["flag"])
		}
		if result["bossName"] != "Baphomet" {
			t.Errorf("Expected boss name Baphomet, got %v", result["bossName"])
		}
		if result["x"] != uint16(150) {
			t.Errorf("Expected x 150, got %v", result["x"])
		}
		if result["y"] != uint16(200) {
			t.Errorf("Expected y 200, got %v", result["y"])
		}
		expectedStatus := "MVP Boss Baphomet is now on location: (150, 200)"
		if result["status"] != expectedStatus {
			t.Errorf("Expected status %q, got %q", expectedStatus, result["status"])
		}
	})

	// Test case 3: Boss monster detected on the map (flag 2)
	t.Run("BossDetectedOnMap", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create a channel to capture hook calls
		resultChan := make(chan map[string]interface{}, 1)

		// Register a hook to capture the result
		hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
			result := arg.(map[string]interface{})
			resultChan <- result
		}, nil)

		// Create test packet arguments
		args := map[string]interface{}{
			"flag": byte(2),
			"name": []byte("Baphomet"),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err != nil {
			t.Fatalf("handleBossMapInfo returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(2) {
			t.Errorf("Expected flag 2, got %v", result["flag"])
		}
		if result["bossName"] != "Baphomet" {
			t.Errorf("Expected boss name Baphomet, got %v", result["bossName"])
		}
		expectedStatus := "MVP Boss Baphomet has been detected on this map!"
		if result["status"] != expectedStatus {
			t.Errorf("Expected status %q, got %q", expectedStatus, result["status"])
		}
	})

	// Test case 4: Boss monster is dead, will respawn after a certain time (flag 3)
	t.Run("BossDeadWillRespawn", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create a channel to capture hook calls
		resultChan := make(chan map[string]interface{}, 1)

		// Register a hook to capture the result
		hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
			result := arg.(map[string]interface{})
			resultChan <- result
		}, nil)

		// Create test packet arguments
		args := map[string]interface{}{
			"flag":    byte(3),
			"name":    []byte("Baphomet"),
			"hours":   byte(2),
			"minutes": byte(30),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err != nil {
			t.Fatalf("handleBossMapInfo returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(3) {
			t.Errorf("Expected flag 3, got %v", result["flag"])
		}
		if result["bossName"] != "Baphomet" {
			t.Errorf("Expected boss name Baphomet, got %v", result["bossName"])
		}
		if result["hours"] != byte(2) {
			t.Errorf("Expected hours 2, got %v", result["hours"])
		}
		if result["minutes"] != byte(30) {
			t.Errorf("Expected minutes 30, got %v", result["minutes"])
		}
		expectedStatus := "MVP Boss Baphomet is dead, but will spawn again in 2 hour(s) and 30 minutes(s)."
		if result["status"] != expectedStatus {
			t.Errorf("Expected status %q, got %q", expectedStatus, result["status"])
		}
	})

	// Test case 5: Unknown flag
	t.Run("UnknownFlag", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create a channel to capture hook calls
		resultChan := make(chan map[string]interface{}, 1)

		// Register a hook to capture the result
		hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
			result := arg.(map[string]interface{})
			resultChan <- result
		}, nil)

		// Create test packet arguments
		args := map[string]interface{}{
			"flag": byte(4), // Unknown flag
			"name": []byte("Baphomet"),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err != nil {
			t.Fatalf("handleBossMapInfo returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(4) {
			t.Errorf("Expected flag 4, got %v", result["flag"])
		}
		if result["bossName"] != "Baphomet" {
			t.Errorf("Expected boss name Baphomet, got %v", result["bossName"])
		}
		expectedStatus := fmt.Sprintf("Unknown boss_map_info result (flag: %d)", 4)
		if result["status"] != expectedStatus {
			t.Errorf("Expected status %q, got %q", expectedStatus, result["status"])
		}
	})

	// Test case 6: Missing flag
	t.Run("MissingFlag", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		manager := NewFieldManager(nil, hookManager)

		// Create test packet arguments with missing flag
		args := map[string]interface{}{
			"name": []byte("Baphomet"),
		}

		// Call handler
		err := manager.handleBossMapInfo(args)
		if err == nil {
			t.Fatalf("handleBossMapInfo should return error for missing flag")
		}
	})
}
