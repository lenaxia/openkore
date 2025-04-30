package servers

import (
	"os"
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// setupTestBaseSend creates a configured BaseSend for testing
func setupTestBaseSend() *core.BaseSend {
	// Create a BaseSend
	baseSend := core.NewBaseSend(hooks.NewHookManager())

	// Configure with test packet constructions
	packetConstructions := map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
		"00E8": {
			ID:         "00E8",
			Name:       "party_organize",
			Format:     "Z24 C C",
			FieldNames: []string{"name", "share1", "share2"},
		},
		"0155": {
			ID:         "0155",
			Name:       "guild_member_positions",
			Format:     "v a*",
			FieldNames: []string{"len", "positions"},
		},
		"02AF": {
			ID:         "02AF",
			Name:       "message_id_encryption_initialized",
			Format:     "",
			FieldNames: []string{},
		},
		// Add more packet constructions as needed for server-specific tests
	}

	// Configure the BaseSend
	baseSend.Configure("ServerType0", packetConstructions)

	return baseSend
}

// TestRegisterServerType0Handlers tests that ServerType0 handlers are properly registered
func TestRegisterServerType0Handlers(t *testing.T) {
	// Create a configured BaseSend
	mockSend := setupTestBaseSend()

	// Register handlers
	RegisterServerType0Handlers(mockSend)

	// In a real implementation with actual handlers, we would verify that they were registered
	// For now, we just make sure the function doesn't panic
}

// TestServerType0Integration tests the integration between ServerType0 handlers and the BaseSend
func TestServerType0Integration(t *testing.T) {
	// Create a configured BaseSend
	baseSend := setupTestBaseSend()

	// Register handlers
	RegisterServerType0Handlers(baseSend)

	// In a real implementation with actual handlers, we would test them here
	// For now, we just make sure the setup works without errors
}

// TestShuffle tests the shuffle method
func TestShuffle(t *testing.T) {
	// Create a configured BaseSend
	baseSend := setupTestBaseSend()

	// Register handlers for the packets we want to test
	baseSend.RegisterHandler("login_request", func(args map[string]interface{}) ([]byte, error) {
		return baseSend.Reconstruct("0064", args)
	})
	baseSend.RegisterHandler("guild_member_positions", func(args map[string]interface{}) ([]byte, error) {
		return baseSend.Reconstruct("0155", args)
	})

	// Create a temporary shuffle.txt file
	tempFile, err := os.CreateTemp("", "shuffle.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the file
	testData := "0064 0065\n0155 0156\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Call the shuffle method
	err = Shuffle(baseSend, tempFile.Name())
	if err != nil {
		t.Fatalf("Shuffle failed: %v", err)
	}

	// Test with non-existent file
	err = Shuffle(baseSend, "non_existent_file.txt")
	if err == nil {
		t.Error("Expected error when shuffling with non-existent file, got nil")
	}

	// Test with invalid file format
	tempFile2, err := os.CreateTemp("", "invalid_shuffle.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile2.Name())

	// Write invalid data to the file
	invalidData := "0064\n0155\n"
	if _, err := tempFile2.Write([]byte(invalidData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile2.Close()

	// Call the shuffle method with invalid file
	err = Shuffle(baseSend, tempFile2.Name())
	if err == nil {
		t.Error("Expected error when shuffling with invalid file format, got nil")
	}
}
