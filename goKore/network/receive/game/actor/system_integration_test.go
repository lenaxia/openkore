package actor

import (
	"testing"
)

// TestActorSystemIntegration2 tests the interaction between all components of the actor system
func TestActorSystemIntegration2(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Register hooks to track events
	hookCalls := make(map[string]int)

	handler.RegisterHook("add_player_list", func(data interface{}) {
		hookCalls["add_player_list"]++
	})

	handler.RegisterHook("add_monster_list", func(data interface{}) {
		hookCalls["add_monster_list"]++
	})

	handler.RegisterHook("player_exist", func(data interface{}) {
		hookCalls["player_exist"]++
	})

	handler.RegisterHook("monster_exist", func(data interface{}) {
		hookCalls["monster_exist"]++
	})

	handler.RegisterHook("player_moved", func(data interface{}) {
		hookCalls["player_moved"]++
	})

	handler.RegisterHook("player_disappeared", func(data interface{}) {
		hookCalls["player_disappeared"]++
	})

	handler.RegisterHook("packet_attack", func(data interface{}) {
		hookCalls["packet_attack"]++
	})

	// Step 1: Create a player
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

	// Verify player was added to the list
	player := handler.playersList.GetByID(playerID)
	if player == nil {
		t.Fatalf("Player was not added to the list")
	}

	// Verify player hooks were called
	if hookCalls["add_player_list"] != 1 {
		t.Errorf("Expected add_player_list hook to be called once, got %d", hookCalls["add_player_list"])
	}

	if hookCalls["player_exist"] != 1 {
		t.Errorf("Expected player_exist hook to be called once, got %d", hookCalls["player_exist"])
	}

	// Step 2: Create a monster
	monsterID := []byte{5, 6, 7, 8}
	monster := NewMonster(monsterID)
	monster.SetName("TestMonster")
	monster.SetBinType(1001)
	handler.monstersList.Add(monster)

	// Verify monster was added to the list
	if handler.monstersList.Count() != 1 {
		t.Errorf("Expected 1 monster in the list, got %d", handler.monstersList.Count())
	}

	// Step 3: Move the player
	moveArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 200, 0, 1, 0},
		"switch":      "007B", // actor_moved
	}

	handler.HandleActorDisplay(moveArgs)

	// Verify player moved hook was called
	if hookCalls["player_moved"] != 1 {
		t.Errorf("Expected player_moved hook to be called once, got %d", hookCalls["player_moved"])
	}

	// Step 4: Player attacks monster
	attackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttack),
		"damage":            int32(100),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(attackArgs)

	// Verify attack hook was called
	if hookCalls["packet_attack"] != 1 {
		t.Errorf("Expected packet_attack hook to be called once, got %d", hookCalls["packet_attack"])
	}

	// Step 5: Update player info
	infoArgs := map[string]interface{}{
		"ID":         playerID,
		"name":       "UpdatedPlayer",
		"partyName":  "TestParty",
		"guildName":  "TestGuild",
		"guildTitle": "TestTitle",
	}

	handler.HandleActorInfo(infoArgs)

	// Verify player info was updated
	if player.Name() != "UpdatedPlayer" {
		t.Errorf("Expected player name 'UpdatedPlayer', got '%s'", player.Name())
	}

	if player.PartyName() != "TestParty" {
		t.Errorf("Expected party name 'TestParty', got '%s'", player.PartyName())
	}

	// Step 6: Monster dies
	monsterDeathArgs := map[string]interface{}{
		"ID":   monsterID,
		"type": byte(DisappearDied),
	}

	handler.HandleActorDiedOrDisappeared(monsterDeathArgs)

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

	// Step 7: Player disappears
	playerDisappearArgs := map[string]interface{}{
		"ID":   playerID,
		"type": byte(DisappearOutOfSight),
	}

	handler.HandleActorDiedOrDisappeared(playerDisappearArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerID) != nil {
		t.Errorf("Player was not removed from the list")
	}

	// Verify player disappear hook was called
	if hookCalls["player_disappeared"] != 1 {
		t.Errorf("Expected player_disappeared hook to be called once, got %d", hookCalls["player_disappeared"])
	}

	// Verify player was added to the old players map
	oldPlayer, exists := handler.playersOld[string(playerID)]
	if !exists {
		t.Errorf("Player was not added to the old players map")
	}

	if oldPlayer.Name() != "UpdatedPlayer" {
		t.Errorf("Expected old player name 'UpdatedPlayer', got '%s'", oldPlayer.Name())
	}

	if !oldPlayer.IsDisappeared() {
		t.Errorf("Expected old player to be marked as disappeared")
	}
}

// TestActorSystemConcurrency tests the actor system under concurrent access
func TestActorSystemConcurrency(t *testing.T) {
	// Create a handler
	handler := NewHandler()

	// Create multiple actors concurrently
	numActors := 10

	// Create players
	for i := 0; i < numActors; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		playerArgs := map[string]interface{}{
			"ID":          id,
			"object_type": byte(0),
			"type":        uint16(0),
			"name":        "Player",
			"coords":      []byte{100, 0, 150, 0, 1, 0},
			"switch":      "0078", // actor_exists
		}

		handler.HandleActorDisplay(playerArgs)
	}

	// Verify all players were added
	if handler.playersList.Count() != numActors {
		t.Errorf("Expected %d players, got %d", numActors, handler.playersList.Count())
	}

	// Create monsters
	for i := 0; i < numActors; i++ {
		id := []byte{byte(i + 100), byte((i + 100) >> 8), byte((i + 100) >> 16), byte((i + 100) >> 24)}
		monster := NewMonster(id)
		monster.SetName("Monster")
		monster.SetBinType(1001)
		handler.monstersList.Add(monster)
	}

	// Verify all monsters were added
	if handler.monstersList.Count() != numActors {
		t.Errorf("Expected %d monsters, got %d", numActors, handler.monstersList.Count())
	}

	// Make players attack monsters
	for i := 0; i < numActors; i++ {
		playerID := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		monsterID := []byte{byte(i + 100), byte((i + 100) >> 8), byte((i + 100) >> 16), byte((i + 100) >> 24)}

		attackArgs := map[string]interface{}{
			"sourceID":          playerID,
			"targetID":          monsterID,
			"type":              byte(ActionAttack),
			"damage":            int32(100),
			"dual_wield_damage": int32(0),
		}

		handler.HandleActorAction(attackArgs)
	}

	// Make monsters die
	for i := 0; i < numActors; i++ {
		monsterID := []byte{byte(i + 100), byte((i + 100) >> 8), byte((i + 100) >> 16), byte((i + 100) >> 24)}

		monsterDeathArgs := map[string]interface{}{
			"ID":   monsterID,
			"type": byte(DisappearDied),
		}

		handler.HandleActorDiedOrDisappeared(monsterDeathArgs)
	}

	// Verify all monsters were removed
	if handler.monstersList.Count() != 0 {
		t.Errorf("Expected 0 monsters after death, got %d", handler.monstersList.Count())
	}

	// Verify all monsters were added to the old monsters map
	if len(handler.monstersOld) != numActors {
		t.Errorf("Expected %d monsters in old monsters map, got %d", numActors, len(handler.monstersOld))
	}

	// Make players disappear
	for i := 0; i < numActors; i++ {
		playerID := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}

		playerDisappearArgs := map[string]interface{}{
			"ID":   playerID,
			"type": byte(DisappearOutOfSight),
		}

		handler.HandleActorDiedOrDisappeared(playerDisappearArgs)
	}

	// Verify all players were removed
	if handler.playersList.Count() != 0 {
		t.Errorf("Expected 0 players after disappearance, got %d", handler.playersList.Count())
	}

	// Verify all players were added to the old players map
	if len(handler.playersOld) != numActors {
		t.Errorf("Expected %d players in old players map, got %d", numActors, len(handler.playersOld))
	}
}
