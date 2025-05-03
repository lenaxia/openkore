package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestSpriteChange(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewCharacterManager(parser, hookManager, logger)

	// Create a test actor
	actorID := uint32(12345)
	actor := &Actor{
		ID:   actorID,
		Name: "TestActor",
	}
	manager.actors[actorID] = actor

	// Test case 1: Job change (type 0)
	t.Run("JobChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(0),                       // Job change
			"value1": uint16(4002),                   // New job ID
			"value2": uint16(0),                      // Not used for job change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's job was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.JobID != 4002 {
			t.Errorf("updatedActor.JobID = %d, want 4002", updatedActor.JobID)
		}
	})

	// Test case 2: Weapon and shield change (type 2)
	t.Run("WeaponShieldChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(2),                       // Weapon/shield change
			"value1": uint16(1201),                   // New weapon ID
			"value2": uint16(2101),                   // New shield ID
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's weapon and shield were updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Weapon != 1201 {
			t.Errorf("updatedActor.Weapon = %d, want 1201", updatedActor.Weapon)
		}

		if updatedActor.Shield != 2101 {
			t.Errorf("updatedActor.Shield = %d, want 2101", updatedActor.Shield)
		}
	})

	// Test case 3: Lower headgear change (type 3)
	t.Run("LowerHeadgearChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(3),                       // Lower headgear change
			"value1": uint16(2202),                   // New lower headgear ID
			"value2": uint16(0),                      // Not used for headgear change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's lower headgear was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Headgear.Low != 2202 {
			t.Errorf("updatedActor.Headgear.Low = %d, want 2202", updatedActor.Headgear.Low)
		}
	})

	// Test case 4: Upper headgear change (type 4)
	t.Run("UpperHeadgearChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(4),                       // Upper headgear change
			"value1": uint16(5001),                   // New upper headgear ID
			"value2": uint16(0),                      // Not used for headgear change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's upper headgear was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Headgear.Top != 5001 {
			t.Errorf("updatedActor.Headgear.Top = %d, want 5001", updatedActor.Headgear.Top)
		}
	})

	// Test case 5: Middle headgear change (type 5)
	t.Run("MiddleHeadgearChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(5),                       // Middle headgear change
			"value1": uint16(2501),                   // New middle headgear ID
			"value2": uint16(0),                      // Not used for headgear change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's middle headgear was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Headgear.Mid != 2501 {
			t.Errorf("updatedActor.Headgear.Mid = %d, want 2501", updatedActor.Headgear.Mid)
		}
	})

	// Test case 6: Hair color change (type 6)
	t.Run("HairColorChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(6),                       // Hair color change
			"value1": uint16(8),                      // New hair color ID
			"value2": uint16(0),                      // Not used for hair color change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's hair color was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.HairColor != 8 {
			t.Errorf("updatedActor.HairColor = %d, want 8", updatedActor.HairColor)
		}
	})

	// Test case 7: Shoes change (type 9)
	t.Run("ShoesChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(9),                       // Shoes change
			"value1": uint16(2412),                   // New shoes ID
			"value2": uint16(0),                      // Not used for shoes change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's shoes were updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Shoes != 2412 {
			t.Errorf("updatedActor.Shoes = %d, want 2412", updatedActor.Shoes)
		}
	})

	// Test case 8: Robe change (type 12)
	t.Run("RobeChange", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(12),                      // Robe change
			"value1": uint16(2510),                   // New robe ID
			"value2": uint16(0),                      // Not used for robe change
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err != nil {
			t.Fatalf("handleSpriteChange() returned error: %v", err)
		}

		// Check that the actor's robe was updated
		updatedActor, found := manager.GetActor(actorID)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if updatedActor.Robe != 2510 {
			t.Errorf("updatedActor.Robe = %d, want 2510", updatedActor.Robe)
		}
	})

	// Test case 9: Unknown sprite type
	t.Run("UnknownSpriteType", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x39, 0x30, 0x00, 0x00}, // 12345 in little-endian
			"type":   uint8(99),                      // Unknown sprite type
			"value1": uint16(1),
			"value2": uint16(2),
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err == nil {
			t.Fatal("handleSpriteChange() did not return error for unknown sprite type")
		}
	})

	// Test case 10: Invalid actor ID
	t.Run("InvalidActorID", func(t *testing.T) {
		args := map[string]interface{}{
			"ID":     []byte{0x42, 0x42, 0x00, 0x00}, // Some other actor ID
			"type":   uint8(0),
			"value1": uint16(1),
			"value2": uint16(2),
		}

		// Call handler
		err := manager.handleSpriteChange(args)
		if err == nil {
			t.Fatal("handleSpriteChange() did not return error for invalid actor ID")
		}
	})
}
