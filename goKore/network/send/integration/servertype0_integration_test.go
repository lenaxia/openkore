package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/handlers/servers"
	"github.com/lenaxia/goKore/network/send/testing/mock"
)

// TestServerType0EndToEnd tests the end-to-end flow of ServerType0 packet sending and shuffling
func TestServerType0EndToEnd(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a mock connection
	mockConn := mock.NewMockConnection()

	// Create a BaseSend instance
	baseSend := core.NewBaseSend(hookManager)
	baseSend.SetConnection(mockConn)

	// Configure with packet constructions
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
	}
	baseSend.Configure("ServerType0", packetConstructions)

	// Register ServerType0 handlers
	servers.RegisterServerType0Handlers(baseSend)

	// Test sending a login request packet
	err := baseSend.SendPacket("login_request", map[string]interface{}{
		"version":    15,
		"username":   "testuser",
		"password":   "testpass",
		"clienttype": 0,
	})
	if err != nil {
		t.Fatalf("Failed to send login_request packet: %v", err)
	}

	// Verify that the packet was sent with the correct ID
	sentPackets := mockConn.GetSentPackets()
	if len(sentPackets) != 1 {
		t.Fatalf("Expected 1 sent packet, got %d", len(sentPackets))
	}

	// The first two bytes of the packet are the message ID in little-endian format
	messageID := sentPackets[0][0] | (sentPackets[0][1] << 8)
	if messageID != 0x0064 {
		t.Errorf("Expected message ID 0x0064, got 0x%04X", messageID)
	}

	// Now test the shuffle functionality
	// Create a temporary shuffle.txt file
	tempFile, err := os.CreateTemp("", "shuffle.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the file
	testData := "0064 0065\n00E8 00E9\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Apply the shuffle
	err = servers.Shuffle(baseSend, tempFile.Name())
	if err != nil {
		t.Fatalf("Shuffle failed: %v", err)
	}

	// Clear the sent packets
	mockConn.ClearSentPackets()

	// Send the login request packet again
	err = baseSend.SendPacket("login_request", map[string]interface{}{
		"version":    15,
		"username":   "testuser",
		"password":   "testpass",
		"clienttype": 0,
	})
	if err != nil {
		t.Fatalf("Failed to send login_request packet after shuffle: %v", err)
	}

	// Verify that the packet was sent with the shuffled ID
	sentPackets = mockConn.GetSentPackets()
	if len(sentPackets) != 1 {
		t.Fatalf("Expected 1 sent packet after shuffle, got %d", len(sentPackets))
	}

	// The first two bytes of the packet should now be the shuffled message ID
	messageID = sentPackets[0][0] | (sentPackets[0][1] << 8)
	if messageID != 0x0065 {
		t.Errorf("Expected shuffled message ID 0x0065, got 0x%04X", messageID)
	}

	// Test sending a party organize packet
	mockConn.ClearSentPackets()

	// Re-register the party_organize handler to use the shuffled ID
	baseSend.RegisterHandler("party_organize", func(args map[string]interface{}) ([]byte, error) {
		// Extract arguments
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid name parameter")
		}

		share1, ok := args["share1"].(int)
		if !ok {
			share1 = 1 // Default value
		}

		share2, ok := args["share2"].(int)
		if !ok {
			share2 = 1 // Default value
		}

		// Construct the packet using the packet definition
		return baseSend.Reconstruct("00E9", map[string]interface{}{
			"name":   name,
			"share1": share1,
			"share2": share2,
		})
	})

	err = baseSend.SendPacket("party_organize", map[string]interface{}{
		"name":   "TestParty",
		"share1": 1,
		"share2": 1,
	})
	if err != nil {
		t.Fatalf("Failed to send party_organize packet: %v", err)
	}

	// Verify that the packet was sent with the shuffled ID
	sentPackets = mockConn.GetSentPackets()
	if len(sentPackets) != 1 {
		t.Fatalf("Expected 1 sent packet for party_organize, got %d", len(sentPackets))
	}

	// The first two bytes of the packet should be the shuffled message ID
	messageID = sentPackets[0][0] | (sentPackets[0][1] << 8)
	if messageID != 0x00E9 {
		t.Errorf("Expected shuffled message ID 0x00E9, got 0x%04X", messageID)
	}

	// Verify the content of the party_organize packet
	// The packet should contain the party name (24 bytes) followed by share1 and share2
	if len(sentPackets[0]) < 26 {
		t.Fatalf("Party organize packet too short: %d bytes", len(sentPackets[0]))
	}

	// Print the entire packet for debugging
	t.Logf("Party organize packet: %v", sentPackets[0])

	// Check share1 and share2 values (should be at the end of the packet)
	// The packet format is: [message_id(2) name(24) share1(1) share2(1)]
	share1 := sentPackets[0][26]
	share2 := sentPackets[0][27]
	if share1 != 1 || share2 != 1 {
		t.Errorf("Expected share1=1, share2=1, got share1=%d, share2=%d", share1, share2)
	}
}
