package actor

import (
	"reflect"
	"testing"
)

// TestHandlerHookSystem tests the hook system in more detail
func TestHandlerHookSystem(t *testing.T) {
	handler := NewHandler()

	// Track hook calls and their parameters
	hookCalls := make(map[string]int)
	hookData := make(map[string]interface{})

	// Register multiple hooks for the same event
	handler.RegisterHook("test_event", func(data interface{}) {
		hookCalls["hook1"]++
		hookData["hook1"] = data
	})

	handler.RegisterHook("test_event", func(data interface{}) {
		hookCalls["hook2"]++
		hookData["hook2"] = data
	})

	// Register hook for a different event
	handler.RegisterHook("other_event", func(data interface{}) {
		hookCalls["hook3"]++
		hookData["hook3"] = data
	})

	// Trigger the test event
	testData := map[string]interface{}{"key": "value"}
	handler.TriggerHook("test_event", testData)

	// Verify both hooks for test_event were called
	if hookCalls["hook1"] != 1 {
		t.Errorf("Expected hook1 to be called once, got %d", hookCalls["hook1"])
	}

	if hookCalls["hook2"] != 1 {
		t.Errorf("Expected hook2 to be called once, got %d", hookCalls["hook2"])
	}

	// Verify the other event hook was not called
	if hookCalls["hook3"] != 0 {
		t.Errorf("Expected hook3 not to be called, got %d", hookCalls["hook3"])
	}

	// Verify the hook data was passed correctly
	if !reflect.DeepEqual(hookData["hook1"], testData) {
		t.Errorf("Expected hook1 data to be %v, got %v", testData, hookData["hook1"])
	}

	if !reflect.DeepEqual(hookData["hook2"], testData) {
		t.Errorf("Expected hook2 data to be %v, got %v", testData, hookData["hook2"])
	}

	// Trigger the other event
	otherData := map[string]interface{}{"other": "data"}
	handler.TriggerHook("other_event", otherData)

	// Verify the other event hook was called
	if hookCalls["hook3"] != 1 {
		t.Errorf("Expected hook3 to be called once, got %d", hookCalls["hook3"])
	}

	// Verify the hook data was passed correctly
	if !reflect.DeepEqual(hookData["hook3"], otherData) {
		t.Errorf("Expected hook3 data to be %v, got %v", otherData, hookData["hook3"])
	}

	// Trigger a non-existent event
	handler.TriggerHook("non_existent_event", nil)

	// Verify no hooks were called
	if hookCalls["hook1"] != 1 || hookCalls["hook2"] != 1 || hookCalls["hook3"] != 1 {
		t.Errorf("Expected no hooks to be called for non-existent event")
	}
}

// TestHandlerActorDisplayComprehensive tests the HandleActorDisplay method with various packet types
func TestHandlerActorDisplayComprehensive(t *testing.T) {
	handler := NewHandler()

	// Test player appearance with actor_exists packet
	playerID := []byte{1, 2, 3, 4}
	playerExistsArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "0078", // actor_exists
	}

	handler.HandleActorDisplay(playerExistsArgs)

	// Verify player was added to the list
	player := handler.playersList.GetByID(playerID)
	if player == nil {
		t.Fatalf("Player was not added to the list")
	}

	// Test player appearance with actor_connected packet
	playerConnectedID := []byte{5, 6, 7, 8}
	playerConnectedArgs := map[string]interface{}{
		"ID":          playerConnectedID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "ConnectedPlayer",
		"coords":      []byte{200, 0, 250, 0, 1, 0},
		"switch":      "0079", // actor_connected
	}

	handler.HandleActorDisplay(playerConnectedArgs)

	// Verify player was added to the list
	connectedPlayer := handler.playersList.GetByID(playerConnectedID)
	if connectedPlayer == nil {
		t.Fatalf("Connected player was not added to the list")
	}

	// Test player appearance with actor_moved packet
	playerMovedArgs := map[string]interface{}{
		"ID":          playerID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{100, 0, 200, 0, 1, 0},
		"switch":      "007B", // actor_moved
	}

	handler.HandleActorDisplay(playerMovedArgs)

	// Test player appearance with actor_spawned packet
	playerSpawnedID := []byte{9, 10, 11, 12}
	playerSpawnedArgs := map[string]interface{}{
		"ID":          playerSpawnedID,
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "SpawnedPlayer",
		"coords":      []byte{100, 0, 150, 0, 1, 0},
		"switch":      "007C", // actor_spawned
	}

	handler.HandleActorDisplay(playerSpawnedArgs)

	// Verify player was added to the list
	spawnedPlayer := handler.playersList.GetByID(playerSpawnedID)
	if spawnedPlayer == nil {
		t.Fatalf("Spawned player was not added to the list")
	}
}

// TestHandlerActorDiedOrDisappearedComprehensive tests the HandleActorDiedOrDisappeared method with various disappearance types
func TestHandlerActorDiedOrDisappearedComprehensive(t *testing.T) {
	handler := NewHandler()

	// Create players with different disappearance types
	playerDiedID := []byte{1, 2, 3, 4}
	playerDied := NewPlayer(playerDiedID)
	playerDied.SetName("DiedPlayer")
	handler.playersList.Add(playerDied)

	playerDisappearedID := []byte{5, 6, 7, 8}
	playerDisappeared := NewPlayer(playerDisappearedID)
	playerDisappeared.SetName("DisappearedPlayer")
	handler.playersList.Add(playerDisappeared)

	playerLoggedOutID := []byte{9, 10, 11, 12}
	playerLoggedOut := NewPlayer(playerLoggedOutID)
	playerLoggedOut.SetName("LoggedOutPlayer")
	handler.playersList.Add(playerLoggedOut)

	playerTeleportedID := []byte{13, 14, 15, 16}
	playerTeleported := NewPlayer(playerTeleportedID)
	playerTeleported.SetName("TeleportedPlayer")
	handler.playersList.Add(playerTeleported)

	// Test player death
	diedArgs := map[string]interface{}{
		"ID":   playerDiedID,
		"type": byte(DisappearDied),
	}

	handler.HandleActorDiedOrDisappeared(diedArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerDiedID) != nil {
		t.Errorf("Died player was not removed from the list")
	}

	// Verify player was added to the old players map
	oldDiedPlayer, exists := handler.playersOld[string(playerDiedID)]
	if !exists {
		t.Errorf("Died player was not added to the old players map")
	}

	if oldDiedPlayer.Name() != "DiedPlayer" {
		t.Errorf("Expected old died player name 'DiedPlayer', got '%s'", oldDiedPlayer.Name())
	}

	if !oldDiedPlayer.IsDead() {
		t.Errorf("Expected old died player to be marked as dead")
	}

	// Test player disappearance
	disappearedArgs := map[string]interface{}{
		"ID":   playerDisappearedID,
		"type": byte(DisappearOutOfSight),
	}

	handler.HandleActorDiedOrDisappeared(disappearedArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerDisappearedID) != nil {
		t.Errorf("Disappeared player was not removed from the list")
	}

	// Verify player was added to the old players map
	oldDisappearedPlayer, exists := handler.playersOld[string(playerDisappearedID)]
	if !exists {
		t.Errorf("Disappeared player was not added to the old players map")
	}

	if oldDisappearedPlayer.Name() != "DisappearedPlayer" {
		t.Errorf("Expected old disappeared player name 'DisappearedPlayer', got '%s'", oldDisappearedPlayer.Name())
	}

	if !oldDisappearedPlayer.IsDisappeared() {
		t.Errorf("Expected old disappeared player to be marked as disappeared")
	}

	// Test player logged out
	loggedOutArgs := map[string]interface{}{
		"ID":   playerLoggedOutID,
		"type": byte(DisappearLoggedOut),
	}

	handler.HandleActorDiedOrDisappeared(loggedOutArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerLoggedOutID) != nil {
		t.Errorf("Logged out player was not removed from the list")
	}

	// Test player teleported
	teleportedArgs := map[string]interface{}{
		"ID":   playerTeleportedID,
		"type": byte(DisappearTeleport),
	}

	handler.HandleActorDiedOrDisappeared(teleportedArgs)

	// Verify player was removed from the list
	if handler.playersList.GetByID(playerTeleportedID) != nil {
		t.Errorf("Teleported player was not removed from the list")
	}

	// Verify player was added to the old players map
	oldTeleportedPlayer, exists := handler.playersOld[string(playerTeleportedID)]
	if !exists {
		t.Errorf("Teleported player was not added to the old players map")
	}

	if oldTeleportedPlayer.Name() != "TeleportedPlayer" {
		t.Errorf("Expected old teleported player name 'TeleportedPlayer', got '%s'", oldTeleportedPlayer.Name())
	}

	if !oldTeleportedPlayer.IsTeleported() {
		t.Errorf("Expected old teleported player to be marked as teleported")
	}
}

// TestHandlerActorActionComprehensive tests the HandleActorAction method with various action types
func TestHandlerActorActionComprehensive(t *testing.T) {
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

	// Test critical attack action
	criticalAttackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttackCritical),
		"damage":            int32(200),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(criticalAttackArgs)

	// Test dual wield attack action
	dualWieldAttackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttack),
		"damage":            int32(100),
		"dual_wield_damage": int32(50),
	}

	handler.HandleActorAction(dualWieldAttackArgs)

	// Test miss attack action
	missAttackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttack),
		"damage":            int32(0),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(missAttackArgs)

	// Test lucky dodge attack action
	luckyDodgeAttackArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          monsterID,
		"type":              byte(ActionAttackLucky),
		"damage":            int32(0),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(luckyDodgeAttackArgs)

	// Test item pickup action
	itemID := []byte{9, 10, 11, 12}
	itemPickupArgs := map[string]interface{}{
		"sourceID":          playerID,
		"targetID":          itemID,
		"type":              byte(ActionItemPickup),
		"damage":            int32(0),
		"dual_wield_damage": int32(0),
	}

	handler.HandleActorAction(itemPickupArgs)
}

// TestHandlerActorInfoComprehensive tests the HandleActorInfo method with various actor types
func TestHandlerActorInfoComprehensive(t *testing.T) {
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

	// Create an NPC
	npcID := []byte{9, 10, 11, 12}
	npc := NewNPC(npcID)
	npc.SetName("TestNPC")
	handler.npcsList.Add(npc)

	// Test player info update
	playerInfoArgs := map[string]interface{}{
		"ID":         playerID,
		"name":       "UpdatedPlayer",
		"partyName":  "TestParty",
		"guildName":  "TestGuild",
		"guildTitle": "TestTitle",
	}

	handler.HandleActorInfo(playerInfoArgs)

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

	// Test monster info update
	monsterInfoArgs := map[string]interface{}{
		"ID":   monsterID,
		"name": "UpdatedMonster",
	}

	handler.HandleActorInfo(monsterInfoArgs)

	// Verify monster info was updated
	if monster.Name() != "UpdatedMonster" {
		t.Errorf("Expected monster name 'UpdatedMonster', got '%s'", monster.Name())
	}

	// Test NPC info update
	npcInfoArgs := map[string]interface{}{
		"ID":   npcID,
		"name": "UpdatedNPC",
	}

	handler.HandleActorInfo(npcInfoArgs)

	// Verify NPC info was updated
	if npc.Name() != "UpdatedNPC" {
		t.Errorf("Expected NPC name 'UpdatedNPC', got '%s'", npc.Name())
	}
}

// TestHandlerEdgeCases tests edge cases for the Handler
func TestHandlerEdgeCases(t *testing.T) {
	handler := NewHandler()

	// Test HandleActorDisplay with missing fields
	missingFieldsArgs := map[string]interface{}{
		"ID": []byte{1, 2, 3, 4},
		// Missing object_type, type, name, coords, switch
	}

	// This should not panic
	handler.HandleActorDisplay(missingFieldsArgs)

	// Test HandleActorDiedOrDisappeared with non-existent actor
	nonExistentActorArgs := map[string]interface{}{
		"ID":   []byte{99, 99, 99, 99},
		"type": byte(DisappearDied),
	}

	// This should not panic
	handler.HandleActorDiedOrDisappeared(nonExistentActorArgs)

	// Test HandleActorAction with non-existent source actor
	nonExistentSourceArgs := map[string]interface{}{
		"sourceID":          []byte{99, 99, 99, 99},
		"targetID":          []byte{1, 2, 3, 4},
		"type":              byte(ActionAttack),
		"damage":            int32(100),
		"dual_wield_damage": int32(0),
	}

	// This should not panic
	handler.HandleActorAction(nonExistentSourceArgs)

	// Test HandleActorInfo with non-existent actor
	nonExistentInfoArgs := map[string]interface{}{
		"ID":   []byte{99, 99, 99, 99},
		"name": "NonExistentActor",
	}

	// This should not panic
	handler.HandleActorInfo(nonExistentInfoArgs)
}
