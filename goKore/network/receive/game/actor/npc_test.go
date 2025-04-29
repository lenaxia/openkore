package actor

import (
	"testing"
	"time"
)

// TestNPCMethods tests all methods of the NPC struct
func TestNPCMethods(t *testing.T) {
	// Create a new NPC
	id := []byte{1, 2, 3, 4}
	npc := NewNPC(id)

	// Test GuildID and SetGuildID
	if npc.GuildID() != 0 {
		t.Errorf("Expected default GuildID to be 0, got %d", npc.GuildID())
	}

	npc.SetGuildID(12345)
	if npc.GuildID() != 12345 {
		t.Errorf("Expected GuildID to be 12345, got %d", npc.GuildID())
	}

	// Test EmblemID and SetEmblemID
	if npc.EmblemID() != 0 {
		t.Errorf("Expected default EmblemID to be 0, got %d", npc.EmblemID())
	}

	npc.SetEmblemID(6789)
	if npc.EmblemID() != 6789 {
		t.Errorf("Expected EmblemID to be 6789, got %d", npc.EmblemID())
	}

	// Test IsDisappeared and SetDisappeared
	if npc.IsDisappeared() {
		t.Errorf("Expected default IsDisappeared to be false")
	}

	npc.SetDisappeared(true)
	if !npc.IsDisappeared() {
		t.Errorf("Expected IsDisappeared to be true")
	}

	// Test GoneTime and SetGoneTime
	zeroTime := time.Time{}
	if !npc.GoneTime().Equal(zeroTime) {
		t.Errorf("Expected default GoneTime to be zero time")
	}

	now := time.Now()
	npc.SetGoneTime(now)
	if !npc.GoneTime().Equal(now) {
		t.Errorf("Expected GoneTime to be %v, got %v", now, npc.GoneTime())
	}

	// Test SendTalk
	// Just call it to ensure it doesn't panic
	npc.SendTalk()

	// Test NameString
	npc.SetName("TestNPC")
	expected := "NPC TestNPC"
	if npc.NameString() != expected {
		t.Errorf("Expected NameString to be '%s', got '%s'", expected, npc.NameString())
	}

	// Test DeepCopy
	npc.SetName("OriginalNPC")
	npc.SetPosition(&Position{X: 10, Y: 20})
	npc.SetGuildID(12345)
	npc.SetEmblemID(6789)
	npc.SetDisappeared(true)

	copy := npc.DeepCopy().(*NPC)

	// Verify copy has the same values
	if copy.Name() != "OriginalNPC" {
		t.Errorf("Expected copy name to be 'OriginalNPC', got '%s'", copy.Name())
	}

	if copy.Position().X != 10 || copy.Position().Y != 20 {
		t.Errorf("Expected copy position to be {10, 20}, got {%d, %d}", copy.Position().X, copy.Position().Y)
	}

	if copy.GuildID() != 12345 {
		t.Errorf("Expected copy GuildID to be 12345, got %d", copy.GuildID())
	}

	if copy.EmblemID() != 6789 {
		t.Errorf("Expected copy EmblemID to be 6789, got %d", copy.EmblemID())
	}

	if !copy.IsDisappeared() {
		t.Errorf("Expected copy IsDisappeared to be true")
	}

	// Modify original and verify copy is not affected
	npc.SetName("ModifiedNPC")
	npc.SetPosition(&Position{X: 30, Y: 40})
	npc.SetGuildID(54321)
	npc.SetEmblemID(9876)
	npc.SetDisappeared(false)

	if copy.Name() != "OriginalNPC" {
		t.Errorf("Expected copy name to remain 'OriginalNPC', got '%s'", copy.Name())
	}

	if copy.Position().X != 10 || copy.Position().Y != 20 {
		t.Errorf("Expected copy position to remain {10, 20}, got {%d, %d}", copy.Position().X, copy.Position().Y)
	}

	if copy.GuildID() != 12345 {
		t.Errorf("Expected copy GuildID to remain 12345, got %d", copy.GuildID())
	}

	if copy.EmblemID() != 6789 {
		t.Errorf("Expected copy EmblemID to remain 6789, got %d", copy.EmblemID())
	}

	if !copy.IsDisappeared() {
		t.Errorf("Expected copy IsDisappeared to remain true")
	}
}
