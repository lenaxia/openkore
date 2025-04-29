package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestOfflineCloneHandlers(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to capture hook calls
	addPlayerListChan := make(chan interface{}, 1)
	playerExistChan := make(chan interface{}, 1)
	playerDisappearedChan := make(chan interface{}, 1)

	// Register hooks to capture results
	hookManager.AddHook("add_player_list", func(hookName string, arg interface{}, userData interface{}) {
		addPlayerListChan <- arg
	}, nil)
	hookManager.AddHook("player_exist", func(hookName string, arg interface{}, userData interface{}) {
		playerExistChan <- arg
	}, nil)
	hookManager.AddHook("player_disappeared", func(hookName string, arg interface{}, userData interface{}) {
		playerDisappearedChan <- arg
	}, nil)

	// Create a handler with the hook manager
	handler := NewHandler()
	handler.hookManager = hookManager

	// Test offline_clone_found
	testID := []byte{0x01, 0x02, 0x03, 0x04}
	testArgs := map[string]interface{}{
		"ID":            testID,
		"name":          []byte("TestPlayer"),
		"jobID":         uint16(1),
		"coord_x":       int16(100),
		"coord_y":       int16(200),
		"robe":          uint16(3),
		"clothes_color": uint16(4),
		"lowhead":       uint16(5),
		"midhead":       uint16(6),
		"tophead":       uint16(7),
		"weapon":        uint16(8),
		"shield":        uint16(9),
		"sex":           byte(1),
		"hair_color":    uint16(10),
	}

	// Call the handler
	err := handler.HandleOfflineCloneFound(testArgs)
	if err != nil {
		t.Errorf("HandleOfflineCloneFound returned an error: %v", err)
	}

	// Check that the player was added to the list
	player := handler.playersList.GetByID(testID)
	if player == nil {
		t.Errorf("Player was not added to the list")
	} else {
		// Verify player properties
		if player.Name() != "TestPlayer" {
			t.Errorf("Expected name %s, got %s", "TestPlayer", player.Name())
		}
		if player.jobID != 1 {
			t.Errorf("Expected jobID %d, got %d", 1, player.jobID)
		}
		if player.Position().X != 100 || player.Position().Y != 200 {
			t.Errorf("Expected position (%d, %d), got (%d, %d)", 100, 200, player.Position().X, player.Position().Y)
		}
		if !player.IsClone() {
			t.Errorf("Expected player to be a clone")
		}
	}

	// Test offline_clone_lost
	testArgs = map[string]interface{}{
		"ID": testID,
	}

	// Call the handler
	err = handler.HandleOfflineCloneLost(testArgs)
	if err != nil {
		t.Errorf("HandleOfflineCloneLost returned an error: %v", err)
	}

	// Check that the player was removed from the list
	if handler.playersList.GetByID(testID) != nil {
		t.Errorf("Player was not removed from the list")
	}

	// Check that the player was added to the old players map
	oldPlayer := handler.playersOld[string(testID)]
	if oldPlayer == nil {
		t.Errorf("Player was not added to the old players map")
	} else {
		// Verify player properties
		if oldPlayer.Name() != "TestPlayer" {
			t.Errorf("Expected name %s, got %s", "TestPlayer", oldPlayer.Name())
		}
		if !oldPlayer.IsClone() {
			t.Errorf("Expected player to be a clone")
		}
		if oldPlayer.GoneTime().IsZero() {
			t.Errorf("Expected gone_time to be set")
		}
	}
}
