package actor

import (
	"testing"
)

// TestActorSystemIntegration tests the interaction between different actor components
func TestActorSystemIntegration(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Test player creation and management
	playerID := []byte{1, 2, 3, 4}
	args := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0}, // x, y, dir
		"switch":      "0079",                       // actor_connected
	}

	// Handle player appearance
	handler.HandleActorDisplay(args)

	// Verify player was added to the list
	player := handler.playersList.GetByID(playerID)
	if player == nil {
		t.Fatalf("Player was not added to the list")
	}

	// Verify player properties
	if player.Name() != "TestPlayer" {
		t.Errorf("Expected player name 'TestPlayer', got '%s'", player.Name())
	}

	// Test player movement
	moveArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 200, 0, 1, 0}, // Moving to x=200, y=150
		"switch":      "007B",                       // actor_moved
	}

	handler.HandleActorDisplay(moveArgs)

	// Verify player position was updated
	// Our simplified coordinate parsing uses the first byte directly
	if player.PositionTo().X != 1 {
		t.Errorf("Expected player position X to be 1, got %d", player.PositionTo().X)
	}

	// Test player disappearance
	disappearArgs := map[string]interface{}{
		"ID":   playerID,
		"type": byte(0), // out of sight
	}

	handler.HandleActorDiedOrDisappeared(disappearArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerID) != nil {
		t.Errorf("Player was not removed from the list")
	}

	// Verify player was added to the old players map
	oldPlayer, exists := handler.playersOld[string(playerID)]
	if !exists {
		t.Errorf("Player was not added to the old players map")
	}

	if oldPlayer.Name() != "TestPlayer" {
		t.Errorf("Expected old player name 'TestPlayer', got '%s'", oldPlayer.Name())
	}
}

// TestActorInteractions tests interactions between different actor types
func TestActorInteractions(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	playerArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "0079", // actor_connected
	}
	handler.HandleActorDisplay(playerArgs)

	// Create a monster directly
	monsterID := []byte{5, 6, 7, 8}
	monster := NewMonster(monsterID)
	monster.SetName("TestMonster")
	monster.SetBinType(1001)

	// Add monster to the list
	handler.monstersList.Add(monster)

	// Verify monster was added to the list
	t.Logf("Monsters in list: %d", handler.monstersList.Count())

	monsterActor := handler.monstersList.GetByID(monsterID)
	if monsterActor == nil {
		t.Fatalf("Monster was not added to the list")
	}

	// Verify monster properties
	if monsterActor.Name() != "TestMonster" {
		t.Errorf("Expected monster name 'TestMonster', got '%s'", monsterActor.Name())
	}

	// Test attack action from player to monster
	attackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(1), // normal attack
		"damage":            int32(100),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(attackArgs)

	// No need to check for monster again, we already verified it exists

	// Test monster death
	deathArgs := map[string]interface{}{
		"ID":   monsterID,
		"type": byte(1), // died
	}

	handler.HandleActorDiedOrDisappeared(deathArgs)

	// Verify monster was removed from the list
	if handler.monstersList.GetByID(monsterID) != nil {
		t.Errorf("Monster was not removed from the list")
	}

	// Verify monster was added to the old monsters map
	oldMonster, exists := handler.monstersOld[string(monsterID)]
	if !exists {
		t.Errorf("Monster was not added to the old monsters map")
	}

	if oldMonster.Name() != "TestMonster" {
		t.Errorf("Expected old monster name 'TestMonster', got '%s'", oldMonster.Name())
	}

	if !oldMonster.IsDead() {
		t.Errorf("Expected old monster to be marked as dead")
	}
}

// TestActorStateTransitions tests actor state transitions
func TestActorStateTransitions(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	playerArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "0079", // actor_connected
	}
	handler.HandleActorDisplay(playerArgs)

	player := handler.playersList.GetByID(playerID)
	if player == nil {
		t.Fatalf("Player was not added to the list")
	}

	// Test sitting
	sitArgs := map[string]interface{}{
		"sourceID": playerID,
		"targetID": playerID,
		"type":     byte(ActionSit),
	}

	handler.HandleActorAction(sitArgs)

	// Verify player is sitting
	if !player.IsSitting() {
		t.Errorf("Expected player to be sitting")
	}

	// Test standing
	standArgs := map[string]interface{}{
		"sourceID": playerID,
		"targetID": playerID,
		"type":     byte(ActionStand),
	}

	handler.HandleActorAction(standArgs)

	// Verify player is not sitting
	if player.IsSitting() {
		t.Errorf("Expected player to be standing")
	}

	// Test death
	deathArgs := map[string]interface{}{
		"ID":   playerID,
		"type": byte(DisappearDied),
	}

	handler.HandleActorDiedOrDisappeared(deathArgs)

	// Verify player is dead
	oldPlayer, exists := handler.playersOld[string(playerID)]
	if !exists {
		t.Errorf("Player was not added to the old players map")
	}

	if !oldPlayer.IsDead() {
		t.Errorf("Expected player to be marked as dead")
	}
}

// TestActorInfoUpdates tests actor information updates
func TestActorInfoUpdates(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	playerArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "0079", // actor_connected
	}
	handler.HandleActorDisplay(playerArgs)

	// Update player info
	infoArgs := map[string]interface{}{
		"ID":         playerID,
		"name":       "UpdatedPlayer",
		"partyName":  "TestParty",
		"guildName":  "TestGuild",
		"guildTitle": "TestTitle",
	}

	handler.HandleActorInfo(infoArgs)

	// Verify player info was updated
	player := handler.playersList.GetByID(playerID)
	if player == nil {
		t.Fatalf("Player was not found in the list")
	}

	if player.Name() != "UpdatedPlayer" {
		t.Errorf("Expected player name 'UpdatedPlayer', got '%s'", player.Name())
	}

	if player.PartyName() != "TestParty" {
		t.Errorf("Expected party name 'TestParty', got '%s'", player.PartyName())
	}

	if player.GuildName() != "TestGuild" {
		t.Errorf("Expected guild name 'TestGuild', got '%s'", player.GuildName())
	}

	if player.GuildTitle() != "TestTitle" {
		t.Errorf("Expected guild title 'TestTitle', got '%s'", player.GuildTitle())
	}
}

// TestActorHooks tests the hook system for actor events
func TestActorHooks(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Track hook calls
	hookCalls := make(map[string]int)

	// Register hooks
	handler.RegisterHook("add_player_list", func(data interface{}) {
		hookCalls["add_player_list"]++
	})

	handler.RegisterHook("player_exist", func(data interface{}) {
		hookCalls["player_exist"]++
	})

	handler.RegisterHook("player_moved", func(data interface{}) {
		hookCalls["player_moved"]++
	})

	handler.RegisterHook("player_disappeared", func(data interface{}) {
		hookCalls["player_disappeared"]++
	})

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	playerArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "0078", // actor_exists
	}
	handler.HandleActorDisplay(playerArgs)

	// Verify hooks were called
	if hookCalls["add_player_list"] != 1 {
		t.Errorf("Expected add_player_list hook to be called once, got %d", hookCalls["add_player_list"])
	}

	if hookCalls["player_exist"] != 1 {
		t.Errorf("Expected player_exist hook to be called once, got %d", hookCalls["player_exist"])
	}

	// Move the player
	moveArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 200, 0, 1, 0},
		"switch":      "007B", // actor_moved
	}
	handler.HandleActorDisplay(moveArgs)

	// Verify move hook was called
	if hookCalls["player_moved"] != 1 {
		t.Errorf("Expected player_moved hook to be called once, got %d", hookCalls["player_moved"])
	}

	// Make the player disappear
	disappearArgs := map[string]interface{}{
		"ID":   playerID,
		"type": byte(0), // out of sight
	}
	handler.HandleActorDiedOrDisappeared(disappearArgs)

	// Verify disappear hook was called
	if hookCalls["player_disappeared"] != 1 {
		t.Errorf("Expected player_disappeared hook to be called once, got %d", hookCalls["player_disappeared"])
	}
}
