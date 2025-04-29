package core

import (
	"testing"
)

func TestChangeTitle(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	manager := NewCharacterManager(parser)

	// Create a test actor
	actorID := uint32(12345)
	actor := &Actor{
		ID:   actorID,
		Name: "TestActor",
	}
	manager.actors[actorID] = actor

	// Test case 1: Basic title change
	t.Run("BasicTitleChange", func(t *testing.T) {
		args := map[string]interface{}{
			"title_id": uint16(123), // Some title ID
		}

		// Call handler
		err := manager.handleChangeTitle(args)
		if err != nil {
			t.Fatalf("handleChangeTitle() returned error: %v", err)
		}

		// Check that the actor's title was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.TitleID != 123 {
			t.Errorf("updatedActor.TitleID = %d, want 123", updatedActor.TitleID)
		}
	})

	// Test case 2: Invalid title ID
	t.Run("InvalidTitleID", func(t *testing.T) {
		args := map[string]interface{}{
			"title_id": "invalid", // Invalid title ID type
		}

		// Call handler
		err := manager.handleChangeTitle(args)
		if err == nil {
			t.Fatal("handleChangeTitle() did not return error for invalid title ID")
		}
	})
}
