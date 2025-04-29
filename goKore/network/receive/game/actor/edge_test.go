package actor

import (
	"bytes"
	"testing"
)

// TestActorEdgeCases tests edge cases for the Actor interface and implementations
func TestActorEdgeCases(t *testing.T) {
	// Test with empty ID
	emptyID := []byte{}
	actor := NewBaseActor(emptyID, "EmptyActor")

	if actor.NameID() != 0 {
		t.Errorf("Expected NameID 0 for empty ID, got %d", actor.NameID())
	}

	// Test with nil ID
	nilActor := NewBaseActor(nil, "NilActor")
	if nilActor.NameID() != 0 {
		t.Errorf("Expected NameID 0 for nil ID, got %d", nilActor.NameID())
	}

	// Test with empty name
	emptyNameActor := NewBaseActor([]byte{1, 2, 3, 4}, "")
	if emptyNameActor.Name() == "" {
		t.Errorf("Expected default name for empty name, got empty string")
	}

	// Test with special characters in name
	specialNameActor := NewBaseActor([]byte{1, 2, 3, 4}, "Special!@#$%&*()")
	specialNameActor.SetName("Special!@#$%&*()")
	if specialNameActor.Name() != "Special!@#$%&*()" {
		t.Errorf("Expected name 'Special!@#$%%&*()', got '%s'", specialNameActor.Name())
	}

	// Test with very long name
	longName := string(make([]byte, 1000))
	longNameActor := NewBaseActor([]byte{1, 2, 3, 4}, "LongName")
	longNameActor.SetName(longName)
	if longNameActor.Name() != longName {
		t.Errorf("Expected long name to be preserved")
	}
}

// TestPlayerEdgeCases tests edge cases for the Player implementation
func TestPlayerEdgeCases(t *testing.T) {
	// Test with zero HP/SP
	player := NewPlayer([]byte{1, 2, 3, 4})
	player.SetHP(0)
	player.SetMaxHP(100)
	player.SetSP(0)
	player.SetMaxSP(100)

	if player.HPPercent() != 0 {
		t.Errorf("Expected HP percent 0, got %d", player.HPPercent())
	}

	if player.SPPercent() != 0 {
		t.Errorf("Expected SP percent 0, got %d", player.SPPercent())
	}

	// Test with zero max HP/SP
	player.SetHP(50)
	player.SetMaxHP(0)
	player.SetSP(50)
	player.SetMaxSP(0)

	if player.HPPercent() != 0 {
		t.Errorf("Expected HP percent 0 for zero max HP, got %d", player.HPPercent())
	}

	if player.SPPercent() != 0 {
		t.Errorf("Expected SP percent 0 for zero max SP, got %d", player.SPPercent())
	}

	// Test damage tracking with empty source/target
	player.AddDamageTaken("", 100)
	player.AddDamageDone("", 200)

	if player.GetDamageTaken("") != 100 {
		t.Errorf("Expected damage taken 100 from empty source, got %d", player.GetDamageTaken(""))
	}

	if player.GetDamageDone("") != 200 {
		t.Errorf("Expected damage done 200 to empty target, got %d", player.GetDamageDone(""))
	}

	// Test damage tracking with non-existent source/target
	if player.GetDamageTaken("nonexistent") != 0 {
		t.Errorf("Expected damage taken 0 from nonexistent source, got %d", player.GetDamageTaken("nonexistent"))
	}

	if player.GetDamageDone("nonexistent") != 0 {
		t.Errorf("Expected damage done 0 to nonexistent target, got %d", player.GetDamageDone("nonexistent"))
	}
}

// TestMonsterEdgeCases tests edge cases for the Monster implementation
func TestMonsterEdgeCases(t *testing.T) {
	// Test with zero HP
	monster := NewMonster([]byte{1, 2, 3, 4})
	monster.SetHP(0)
	monster.SetMaxHP(100)

	if monster.HPPercent() != 0 {
		t.Errorf("Expected HP percent 0, got %d", monster.HPPercent())
	}

	// Test with zero max HP
	monster.SetHP(50)
	monster.SetMaxHP(0)

	if monster.HPPercent() != 0 {
		t.Errorf("Expected HP percent 0 for zero max HP, got %d", monster.HPPercent())
	}

	// Test damage tracking with empty source
	monster.AddDamageFrom("", 100)

	if monster.GetDamageFrom("") != 100 {
		t.Errorf("Expected damage 100 from empty source, got %d", monster.GetDamageFrom(""))
	}

	// Test damage tracking with non-existent source
	if monster.GetDamageFrom("nonexistent") != 0 {
		t.Errorf("Expected damage 0 from nonexistent source, got %d", monster.GetDamageFrom("nonexistent"))
	}
}

// TestActorIDComparison tests ID comparison between actors
func TestActorIDComparison(t *testing.T) {
	// Create two actors with the same ID
	id1 := []byte{1, 2, 3, 4}
	id2 := []byte{1, 2, 3, 4}

	actor1 := NewBaseActor(id1, "Actor1")
	actor2 := NewBaseActor(id2, "Actor2")

	// IDs should be equal even though they are different byte slices
	if !bytes.Equal(actor1.ID(), actor2.ID()) {
		t.Errorf("Expected actor IDs to be equal")
	}

	// Create an actor with a different ID
	id3 := []byte{5, 6, 7, 8}
	actor3 := NewBaseActor(id3, "Actor3")

	// IDs should not be equal
	if bytes.Equal(actor1.ID(), actor3.ID()) {
		t.Errorf("Expected actor IDs to be different")
	}
}

// TestActorDeepCopy tests deep copying of actors
func TestActorDeepCopy(t *testing.T) {
	// Create a base actor
	id := []byte{1, 2, 3, 4}
	actor := NewBaseActor(id, "Actor")
	actor.SetName("TestActor")
	actor.SetPosition(&Position{X: 10, Y: 20})
	actor.SetPositionTo(&Position{X: 30, Y: 40})
	actor.SetAvoid(true)
	actor.SetWalkSpeed(0.5)

	// Create a deep copy
	copy := actor.DeepCopy().(*BaseActor)

	// Modify the original actor
	actor.SetName("ModifiedActor")
	actor.SetPosition(&Position{X: 50, Y: 60})
	actor.SetPositionTo(&Position{X: 70, Y: 80})
	actor.SetAvoid(false)
	actor.SetWalkSpeed(0.8)

	// Verify the copy was not affected
	if copy.Name() != "TestActor" {
		t.Errorf("Expected copy name 'TestActor', got '%s'", copy.Name())
	}

	if copy.Position().X != 10 || copy.Position().Y != 20 {
		t.Errorf("Expected copy position {10, 20}, got {%d, %d}", copy.Position().X, copy.Position().Y)
	}

	if copy.PositionTo().X != 30 || copy.PositionTo().Y != 40 {
		t.Errorf("Expected copy positionTo {30, 40}, got {%d, %d}", copy.PositionTo().X, copy.PositionTo().Y)
	}

	if !copy.IsAvoid() {
		t.Errorf("Expected copy avoid to be true")
	}

	if copy.WalkSpeed() != 0.5 {
		t.Errorf("Expected copy walk speed 0.5, got %f", copy.WalkSpeed())
	}

	// Test ID deep copy
	// Modify the original ID
	actor.ID()[0] = 99

	// Verify the copy's ID was not affected
	if copy.ID()[0] != 1 {
		t.Errorf("Expected copy ID[0] to be 1, got %d", copy.ID()[0])
	}
}
