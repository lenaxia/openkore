package actor

import (
	"testing"
)

// TestHandlerActorDisplay tests the HandleActorDisplay method
func TestHandlerActorDisplay(t *testing.T) {
	handler := NewHandler()

	// Test player appearance
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

	// Verify player properties
	if player.Name() != "TestPlayer" {
		t.Errorf("Expected player name 'TestPlayer', got '%s'", player.Name())
	}

	// Test monster appearance
	monsterID := []byte{5, 6, 7, 8}
	monster := NewMonster(monsterID)
	monster.SetName("TestMonster")
	monster.SetBinType(1001)
	handler.monstersList.Add(monster)

	// Create NPC directly
	npcID := []byte{9, 10, 11, 12}
	npc := NewNPC(npcID)
	npc.SetName("TestNPC")

	// Add NPC to the list
	handler.npcsList.Add(npc)

	// Verify NPC was added to the list
	t.Logf("NPCs in list: %d", handler.npcsList.Count())

	npcActor := handler.npcsList.GetByID(npcID)
	if npcActor == nil {
		t.Fatalf("NPC was not added to the list")
	}

	// Verify NPC properties
	if npc.Name() != "TestNPC" {
		t.Errorf("Expected NPC name 'TestNPC', got '%s'", npc.Name())
	}
}

// TestHandlerActorDiedOrDisappeared tests the HandleActorDiedOrDisappeared method
func TestHandlerActorDiedOrDisappeared(t *testing.T) {
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	player := NewPlayer(playerID)
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Test player death
	deathArgs := map[string]interface{}{
		"ID":   playerID,
		"type": byte(DisappearDied),
	}

	handler.HandleActorDiedOrDisappeared(deathArgs)

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

	if !oldPlayer.IsDead() {
		t.Errorf("Expected old player to be marked as dead")
	}

	// Create a monster
	monsterID := []byte{5, 6, 7, 8}
	monster := NewMonster(monsterID)
	monster.SetName("TestMonster")
	handler.monstersList.Add(monster)

	// Test monster disappearance
	disappearArgs := map[string]interface{}{
		"ID":   monsterID,
		"type": byte(DisappearOutOfSight),
	}

	handler.HandleActorDiedOrDisappeared(disappearArgs)

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

	if !oldMonster.IsDisappeared() {
		t.Errorf("Expected old monster to be marked as disappeared")
	}
}

// TestHandlerActorAction tests the HandleActorAction method
func TestHandlerActorAction(t *testing.T) {
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	player := NewPlayer(playerID)
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Create a monster
	monsterID := []byte{5, 6, 7, 8}
	monster := NewMonster(monsterID)
	monster.SetName("TestMonster")
	handler.monstersList.Add(monster)

	// Test sit action
	sitArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          playerID,
		"type":              byte(ActionSit),
		"damage":            int32(0),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(sitArgs)

	// Verify player is sitting
	if !player.IsSitting() {
		t.Errorf("Expected player to be sitting")
	}

	// Test stand action
	standArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          playerID,
		"type":              byte(ActionStand),
		"damage":            int32(0),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(standArgs)

	// Verify player is not sitting
	if player.IsSitting() {
		t.Errorf("Expected player to be standing")
	}

	// Test attack action
	attackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttack),
		"damage":            int32(100),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(attackArgs)
}

// TestHandlerActorInfo tests the HandleActorInfo method
func TestHandlerActorInfo(t *testing.T) {
	handler := NewHandler()

	// Create a player
	playerID := []byte{1, 2, 3, 4}
	player := NewPlayer(playerID)
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Test player info update
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

	if player.GuildName() != "TestGuild" {
		t.Errorf("Expected guild name 'TestGuild', got '%s'", player.GuildName())
	}

	if player.GuildTitle() != "TestTitle" {
		t.Errorf("Expected guild title 'TestTitle', got '%s'", player.GuildTitle())
	}
}

// TestHandlerHooks tests the hook system
func TestHandlerHooks(t *testing.T) {
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
}
