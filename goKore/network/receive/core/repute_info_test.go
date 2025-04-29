package core

import (
	"testing"
)

func TestReputeInfo(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	manager := NewCharacterManager(parser)

	// Test case 1: Single reputation entry
	t.Run("SingleReputationEntry", func(t *testing.T) {
		// Create test packet arguments
		// Each reputation entry is 16 bytes (4 uint32 values)
		reputeInfo := []byte{
			0x01, 0x00, 0x00, 0x00, // type = 1
			0x02, 0x00, 0x00, 0x00, // type2 = 2
			0x64, 0x00, 0x00, 0x00, // points = 100
			0xC8, 0x00, 0x00, 0x00, // points2 = 200
		}

		args := map[string]interface{}{
			"reputeInfo": reputeInfo,
		}

		// Call handler
		err := manager.handleReputeInfo(args)
		if err != nil {
			t.Fatalf("handleReputeInfo() returned error: %v", err)
		}

		// Check that the reputation list was updated
		reputations := manager.GetReputations()
		if len(reputations) != 1 {
			t.Fatalf("len(reputations) = %d, want 1", len(reputations))
		}

		if reputations[0].Type != 1 {
			t.Errorf("reputations[0].Type = %d, want 1", reputations[0].Type)
		}

		if reputations[0].Type2 != 2 {
			t.Errorf("reputations[0].Type2 = %d, want 2", reputations[0].Type2)
		}

		if reputations[0].Points != 100 {
			t.Errorf("reputations[0].Points = %d, want 100", reputations[0].Points)
		}

		if reputations[0].Points2 != 200 {
			t.Errorf("reputations[0].Points2 = %d, want 200", reputations[0].Points2)
		}
	})

	// Test case 2: Multiple reputation entries
	t.Run("MultipleReputationEntries", func(t *testing.T) {
		// Create test packet arguments
		// Each reputation entry is 16 bytes (4 uint32 values)
		reputeInfo := []byte{
			0x01, 0x00, 0x00, 0x00, // type = 1
			0x02, 0x00, 0x00, 0x00, // type2 = 2
			0x64, 0x00, 0x00, 0x00, // points = 100
			0xC8, 0x00, 0x00, 0x00, // points2 = 200

			0x03, 0x00, 0x00, 0x00, // type = 3
			0x04, 0x00, 0x00, 0x00, // type2 = 4
			0x2C, 0x01, 0x00, 0x00, // points = 300
			0x90, 0x01, 0x00, 0x00, // points2 = 400
		}

		args := map[string]interface{}{
			"reputeInfo": reputeInfo,
		}

		// Call handler
		err := manager.handleReputeInfo(args)
		if err != nil {
			t.Fatalf("handleReputeInfo() returned error: %v", err)
		}

		// Check that the reputation list was updated
		reputations := manager.GetReputations()
		if len(reputations) != 2 {
			t.Fatalf("len(reputations) = %d, want 2", len(reputations))
		}

		// Check first reputation entry
		if reputations[0].Type != 1 {
			t.Errorf("reputations[0].Type = %d, want 1", reputations[0].Type)
		}

		if reputations[0].Type2 != 2 {
			t.Errorf("reputations[0].Type2 = %d, want 2", reputations[0].Type2)
		}

		if reputations[0].Points != 100 {
			t.Errorf("reputations[0].Points = %d, want 100", reputations[0].Points)
		}

		if reputations[0].Points2 != 200 {
			t.Errorf("reputations[0].Points2 = %d, want 200", reputations[0].Points2)
		}

		// Check second reputation entry
		if reputations[1].Type != 3 {
			t.Errorf("reputations[1].Type = %d, want 3", reputations[1].Type)
		}

		if reputations[1].Type2 != 4 {
			t.Errorf("reputations[1].Type2 = %d, want 4", reputations[1].Type2)
		}

		if reputations[1].Points != 300 {
			t.Errorf("reputations[1].Points = %d, want 300", reputations[1].Points)
		}

		if reputations[1].Points2 != 400 {
			t.Errorf("reputations[1].Points2 = %d, want 400", reputations[1].Points2)
		}
	})

	// Test case 3: Invalid reputation info
	t.Run("InvalidReputationInfo", func(t *testing.T) {
		// Create test packet arguments with invalid reputation info
		args := map[string]interface{}{
			"reputeInfo": "invalid",
		}

		// Call handler
		err := manager.handleReputeInfo(args)
		if err == nil {
			t.Fatal("handleReputeInfo() did not return error for invalid reputation info")
		}
	})

	// Test case 4: Missing reputation info
	t.Run("MissingReputationInfo", func(t *testing.T) {
		// Create test packet arguments with missing reputation info
		args := map[string]interface{}{}

		// Call handler
		err := manager.handleReputeInfo(args)
		if err == nil {
			t.Fatal("handleReputeInfo() did not return error for missing reputation info")
		}
	})
}
