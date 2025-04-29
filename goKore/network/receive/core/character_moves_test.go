package core

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestCharacterMoves(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	// Extend the Actor struct to include position information
	// This should be done in the actual implementation

	// Test case 1: Basic movement
	t.Run("BasicMovement", func(t *testing.T) {
		// Create a test actor
		actorID := uint32(12345)
		actor := &Actor{
			ID:   actorID,
			Name: "TestActor",
		}

		// Add position fields to the actor
		manager.actors[actorID] = actor

		// Create test packet arguments
		// The packet format is: 0087 <walk start time>.L <walk data>.6B
		args := map[string]interface{}{
			"walkStartTime": uint32(time.Now().Unix()),
			"coords":        []byte{10, 10, 15, 5, 0, 0}, // From (10,10) to (15,5) - Northeast
		}

		// Call handler
		err := manager.handleCharacterMoves(args)
		if err != nil {
			t.Fatalf("handleCharacterMoves() returned error: %v", err)
		}

		// Check that the actor's position was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		// Check that the position fields were added and updated correctly
		if updatedActor.Pos.X != 10 {
			t.Errorf("updatedActor.Pos.X = %d, want 10", updatedActor.Pos.X)
		}
		if updatedActor.Pos.Y != 10 {
			t.Errorf("updatedActor.Pos.Y = %d, want 10", updatedActor.Pos.Y)
		}
		if updatedActor.PosTo.X != 15 {
			t.Errorf("updatedActor.PosTo.X = %d, want 15", updatedActor.PosTo.X)
		}
		if updatedActor.PosTo.Y != 5 {
			t.Errorf("updatedActor.PosTo.Y = %d, want 5", updatedActor.PosTo.Y)
		}

		// Check that the time_move field was set
		if updatedActor.TimeMove == 0 {
			t.Error("updatedActor.TimeMove is zero")
		}

		// Check that the direction was updated
		if updatedActor.Look.Body != 1 {
			t.Errorf("updatedActor.Look.Body = %d, want 1 (northeast)", updatedActor.Look.Body)
		}
	})

	// Test case 2: Invalid coordinates
	t.Run("InvalidCoordinates", func(t *testing.T) {
		// Create a test actor
		actorID := uint32(67890)
		actor := &Actor{
			ID:   actorID,
			Name: "TestActor2",
		}

		// Add position fields to the actor
		manager.actors[actorID] = actor

		// Create test packet arguments with invalid coordinates
		args := map[string]interface{}{
			"walkStartTime": uint32(time.Now().Unix()),
			"coords":        []byte{1}, // Invalid coordinates (too short)
		}

		// Call handler
		err := manager.handleCharacterMoves(args)
		if err == nil {
			t.Fatal("handleCharacterMoves() did not return error for invalid coordinates")
		}

		// Actor's position should not be updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}
		if updatedActor.TimeMove != 0 {
			t.Errorf("updatedActor.TimeMove = %d, want 0", updatedActor.TimeMove)
		}
	})

	// Test case 3: Missing actor
	t.Run("MissingActor", func(t *testing.T) {
		// Clear the actors map
		manager.actors = make(map[uint32]*Actor)

		// Create test packet arguments
		args := map[string]interface{}{
			"walkStartTime": uint32(time.Now().Unix()),
			"coords":        []byte{20, 20, 25, 25, 0, 0}, // From (20,20) to (25,25)
		}

		// Call handler
		err := manager.handleCharacterMoves(args)
		if err == nil {
			t.Fatal("handleCharacterMoves() did not return error for missing actor")
		}
	})
}
