package actor

import (
	"sync"
	"testing"
)

// TestActorListConcurrency tests concurrent access to the ActorList
func TestActorListConcurrency(t *testing.T) {
	list := NewActorList()

	// Create a large number of actors
	numActors := 100
	actors := make([]Actor, numActors)
	for i := 0; i < numActors; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		actors[i] = NewBaseActor(id, "Actor")
	}

	// Add actors concurrently
	var wg sync.WaitGroup
	wg.Add(numActors)
	for i := 0; i < numActors; i++ {
		go func(i int) {
			defer wg.Done()
			list.Add(actors[i])
		}(i)
	}
	wg.Wait()

	// Verify all actors were added
	if list.Count() != numActors {
		t.Errorf("Expected %d actors, got %d", numActors, list.Count())
	}

	// Get actors concurrently
	wg.Add(numActors)
	for i := 0; i < numActors; i++ {
		go func(i int) {
			defer wg.Done()
			actor := list.GetByID(actors[i].ID())
			if actor == nil {
				t.Errorf("Actor %d not found", i)
			}
		}(i)
	}
	wg.Wait()

	// Remove actors concurrently
	wg.Add(numActors)
	for i := 0; i < numActors; i++ {
		go func(i int) {
			defer wg.Done()
			list.Remove(actors[i])
		}(i)
	}
	wg.Wait()

	// Verify all actors were removed
	if list.Count() != 0 {
		t.Errorf("Expected 0 actors, got %d", list.Count())
	}
}

// TestActorListForEach tests the ForEach method
func TestActorListForEach(t *testing.T) {
	list := NewActorList()

	// Create and add actors
	numActors := 10
	for i := 0; i < numActors; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		actor := NewBaseActor(id, "Actor")
		list.Add(actor)
	}

	// Count actors using ForEach
	count := 0
	list.ForEach(func(actor Actor) {
		count++
	})

	if count != numActors {
		t.Errorf("Expected ForEach to visit %d actors, visited %d", numActors, count)
	}

	// Sum actor IDs using ForEach
	sum := 0
	list.ForEach(func(actor Actor) {
		sum += int(actor.ID()[0])
	})

	expectedSum := 0
	for i := 0; i < numActors; i++ {
		expectedSum += i
	}

	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}
}

// TestActorListFilter tests the Filter method
func TestActorListFilter(t *testing.T) {
	list := NewActorList()

	// Create and add actors
	numActors := 10
	for i := 0; i < numActors; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		actor := NewBaseActor(id, "Actor")
		list.Add(actor)
	}

	// Filter even-numbered actors
	evenList := list.Filter(func(actor Actor) bool {
		return actor.ID()[0]%2 == 0
	})

	if evenList.Count() != numActors/2 {
		t.Errorf("Expected %d even-numbered actors, got %d", numActors/2, evenList.Count())
	}

	// Filter odd-numbered actors
	oddList := list.Filter(func(actor Actor) bool {
		return actor.ID()[0]%2 == 1
	})

	if oddList.Count() != numActors/2 {
		t.Errorf("Expected %d odd-numbered actors, got %d", numActors/2, oddList.Count())
	}

	// Filter actors with ID[0] < 5
	lessThan5List := list.Filter(func(actor Actor) bool {
		return actor.ID()[0] < 5
	})

	if lessThan5List.Count() != 5 {
		t.Errorf("Expected 5 actors with ID[0] < 5, got %d", lessThan5List.Count())
	}

	// Filter actors with ID[0] >= 5
	greaterThanEqualTo5List := list.Filter(func(actor Actor) bool {
		return actor.ID()[0] >= 5
	})

	if greaterThanEqualTo5List.Count() != 5 {
		t.Errorf("Expected 5 actors with ID[0] >= 5, got %d", greaterThanEqualTo5List.Count())
	}
}

// TestSpecializedLists tests the specialized actor lists
func TestSpecializedLists(t *testing.T) {
	// Test PlayersList
	playersList := NewPlayersList()

	// Create and add players
	numPlayers := 5
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		players[i] = NewPlayer(id)
		players[i].SetName("Player" + string(byte('0'+i)))
		playersList.Add(players[i])
	}

	// Verify players were added
	if playersList.Count() != numPlayers {
		t.Errorf("Expected %d players, got %d", numPlayers, playersList.Count())
	}

	// Get players by ID
	for i := 0; i < numPlayers; i++ {
		player := playersList.GetByID(players[i].ID())
		if player == nil {
			t.Errorf("Player %d not found", i)
		}
		if player.Name() != "Player"+string(byte('0'+i)) {
			t.Errorf("Expected player name 'Player%c', got '%s'", byte('0'+i), player.Name())
		}
	}

	// Test MonstersList
	monstersList := NewMonstersList()

	// Create and add monsters
	numMonsters := 5
	monsters := make([]*Monster, numMonsters)
	for i := 0; i < numMonsters; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		monsters[i] = NewMonster(id)
		monsters[i].SetName("Monster" + string(byte('0'+i)))
		monstersList.Add(monsters[i])
	}

	// Verify monsters were added
	if monstersList.Count() != numMonsters {
		t.Errorf("Expected %d monsters, got %d", numMonsters, monstersList.Count())
	}

	// Get monsters by ID
	for i := 0; i < numMonsters; i++ {
		monster := monstersList.GetByID(monsters[i].ID())
		if monster == nil {
			t.Errorf("Monster %d not found", i)
		}
		if monster.Name() != "Monster"+string(byte('0'+i)) {
			t.Errorf("Expected monster name 'Monster%c', got '%s'", byte('0'+i), monster.Name())
		}
	}

	// Test NPCsList
	npcsList := NewNPCsList()

	// Create and add NPCs
	numNPCs := 5
	npcs := make([]*NPC, numNPCs)
	for i := 0; i < numNPCs; i++ {
		id := []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
		npcs[i] = NewNPC(id)
		npcs[i].SetName("NPC" + string(byte('0'+i)))
		npcsList.Add(npcs[i])
	}

	// Verify NPCs were added
	if npcsList.Count() != numNPCs {
		t.Errorf("Expected %d NPCs, got %d", numNPCs, npcsList.Count())
	}

	// Get NPCs by ID
	for i := 0; i < numNPCs; i++ {
		npc := npcsList.GetByID(npcs[i].ID())
		if npc == nil {
			t.Errorf("NPC %d not found", i)
		}
		if npc.Name() != "NPC"+string(byte('0'+i)) {
			t.Errorf("Expected NPC name 'NPC%c', got '%s'", byte('0'+i), npc.Name())
		}
	}
}

// TestActorListEdgeCases tests edge cases for the ActorList
func TestActorListEdgeCases(t *testing.T) {
	list := NewActorList()

	// Test with nil actor
	list.Add(nil)
	if list.Count() != 0 {
		t.Errorf("Expected 0 actors after adding nil, got %d", list.Count())
	}

	// Test with nil ID
	actor := NewBaseActor(nil, "NilIDActor")
	list.Add(actor)

	// Test removing non-existent actor
	nonExistentActor := NewBaseActor([]byte{99, 99, 99, 99}, "NonExistentActor")
	list.Remove(nonExistentActor)

	// Test removing by non-existent ID
	list.RemoveByID([]byte{99, 99, 99, 99})

	// Test getting by non-existent ID
	nonExistentActorResult := list.GetByID([]byte{99, 99, 99, 99})
	if nonExistentActorResult != nil {
		t.Errorf("Expected nil for non-existent actor ID, got %v", nonExistentActorResult)
	}

	// Test ForEach with empty list
	list.Clear()
	count := 0
	list.ForEach(func(actor Actor) {
		count++
	})
	if count != 0 {
		t.Errorf("Expected ForEach to visit 0 actors in empty list, visited %d", count)
	}

	// Test Filter with empty list
	filteredList := list.Filter(func(actor Actor) bool {
		return true
	})
	if filteredList.Count() != 0 {
		t.Errorf("Expected filtered list to have 0 actors, got %d", filteredList.Count())
	}
}
