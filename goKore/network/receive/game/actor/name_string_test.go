package actor

import (
	"testing"
)

// TestBaseActorNameString tests the NameString method of BaseActor
func TestBaseActorNameString(t *testing.T) {
	// Create a base actor
	id := []byte{1, 2, 3, 4}
	actor := NewBaseActor(id, "TestActor")
	actor.SetName("BaseActorName")

	// Test NameString
	expected := "TestActor BaseActorName"
	if actor.NameString() != expected {
		t.Errorf("Expected NameString to be '%s', got '%s'", expected, actor.NameString())
	}

	// Test with empty name
	actor.SetName("")
	expected = "TestActor Unknown #67305985"
	if actor.NameString() != expected {
		t.Errorf("Expected NameString with empty name to be '%s', got '%s'", expected, actor.NameString())
	}

	// Test with special characters in name
	actor.SetName("Special!@#$%&*()")
	expected = "TestActor Special!@#$%&*()"
	if actor.NameString() != expected {
		t.Errorf("Expected NameString with special characters to be '%s', got '%s'", expected, actor.NameString())
	}
}

// TestPlayerNameString tests the NameString method of Player
func TestPlayerNameString(t *testing.T) {
	// Create a player
	id := []byte{1, 2, 3, 4}
	player := NewPlayer(id)
	player.SetName("PlayerName")

	// Test NameString with job
	expected := "Player PlayerName (Job 0)"
	if player.NameString() != expected {
		t.Errorf("Expected NameString to be '%s', got '%s'", expected, player.NameString())
	}

	// Test with different job ID
	player.SetJob(10)
	expected = "Player PlayerName (Job 10)"
	if player.NameString() != expected {
		t.Errorf("Expected NameString with job to be '%s', got '%s'", expected, player.NameString())
	}
}

// TestMonsterNameString tests the NameString method of Monster
func TestMonsterNameString(t *testing.T) {
	// Create a monster
	id := []byte{1, 2, 3, 4}
	monster := NewMonster(id)
	monster.SetName("MonsterName")

	// Test NameString
	expected := "Monster MonsterName"
	if monster.NameString() != expected {
		t.Errorf("Expected NameString to be '%s', got '%s'", expected, monster.NameString())
	}
}

// TestNPCNameString tests the NameString method of NPC
func TestNPCNameString(t *testing.T) {
	// Create an NPC
	id := []byte{1, 2, 3, 4}
	npc := NewNPC(id)
	npc.SetName("NPCName")

	// Test NameString
	expected := "NPC NPCName"
	if npc.NameString() != expected {
		t.Errorf("Expected NameString to be '%s', got '%s'", expected, npc.NameString())
	}
}
