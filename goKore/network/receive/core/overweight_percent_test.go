package core

import (
	"testing"
)

func TestOverweightPercent(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	manager := NewCharacterManager(parser)

	// Create a test actor
	actorID := uint32(12345)
	actor := &Actor{
		ID:   actorID,
		Name: "TestActor",
	}
	manager.actors[actorID] = actor

	// Test case 1: Basic overweight percent
	t.Run("BasicOverweightPercent", func(t *testing.T) {
		args := map[string]interface{}{
			"percent": uint8(75), // 75% overweight
		}

		// Call handler
		err := manager.handleOverweightPercent(args)
		if err != nil {
			t.Fatalf("handleOverweightPercent() returned error: %v", err)
		}

		// Check that the actor's overweight percent was updated
		// In a real implementation, we would need to get the main character's ID
		// For testing purposes, we'll use the actor that's already in the manager
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.OverweightPercent != 75 {
			t.Errorf("updatedActor.OverweightPercent = %d, want 75", updatedActor.OverweightPercent)
		}
	})

	// Test case 2: Invalid percent
	t.Run("InvalidPercent", func(t *testing.T) {
		args := map[string]interface{}{
			"percent": "invalid", // Invalid percent type
		}

		// Call handler
		err := manager.handleOverweightPercent(args)
		if err == nil {
			t.Fatal("handleOverweightPercent() did not return error for invalid percent")
		}
	})

	// Test case 3: Missing percent
	t.Run("MissingPercent", func(t *testing.T) {
		args := map[string]interface{}{}

		// Call handler
		err := manager.handleOverweightPercent(args)
		if err == nil {
			t.Fatal("handleOverweightPercent() did not return error for missing percent")
		}
	})
}
