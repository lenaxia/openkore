package integration

import (
	"context"
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/login"
	"github.com/lenaxia/goKore/test/utils"
)

// TestLoginSequenceWithMocks tests the full login sequence using mock connections
func TestLoginSequenceWithMocks(t *testing.T) {

	// Create mock network manager
	networkManager := login.NewMockNetworkManager()

	// Create login config
	config := login.NewLoginConfig("botijo0", "password", "testserver")

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Start login process in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to signal when login is complete
	loginDone := make(chan error, 1)
	go func() {
		loginDone <- loginManager.Login(ctx)
	}()

	// Simulate server responses
	// The MockNetworkManager should already have handlers for these packets
	networkManager.SimulateReceivePacket("0AC4", []byte{0xC4, 0x0A, 0xE0, 0x00, 0xE5, 0x5D, 0xF6, 0xC1, 0x82, 0x84, 0x1E, 0x00, 0x01, 0x2C, 0x9C, 0x53})
	networkManager.SimulateReceivePacket("006B", []byte{0x6B, 0x00, 0xB6, 0x00, 0x0F, 0x0F, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	networkManager.SimulateReceivePacket("0AC5", []byte{0xC5, 0x0A, 0xF2, 0x49, 0x02, 0x00, 0x67, 0x65, 0x66, 0x5F, 0x66, 0x69, 0x6C, 0x64, 0x30, 0x37})
	networkManager.SimulateReceivePacket("02EB", []byte{0xEB, 0x02, 0xC9, 0x3E, 0x82, 0x02, 0x3D, 0x8B, 0xF0, 0x05, 0x05, 0x00, 0x00})

	// Wait for login to complete
	select {
	case err := <-loginDone:
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Login timed out")
	}

	// Verify session data
	sessionData := sessionStore.GetSessionData()

	// Verify account ID
	expectedAccountID := utils.AccountInfoPacket.ExpectedFields["accountID"].([]byte)
	if string(sessionData.AccountID) != string(expectedAccountID) {
		t.Errorf("Expected account ID %v, got %v", expectedAccountID, sessionData.AccountID)
	}

	// Verify char ID
	expectedCharID := utils.CharacterMapInfoPacket.ExpectedFields["charID"].([]byte)
	if string(sessionData.CharID) != string(expectedCharID) {
		t.Errorf("Expected char ID %v, got %v", expectedCharID, sessionData.CharID)
	}

	// Verify map name
	expectedMapName := utils.CharacterMapInfoPacket.ExpectedFields["mapName"].(string)
	if sessionData.MapName != expectedMapName {
		t.Errorf("Expected map name %v, got %v", expectedMapName, sessionData.MapName)
	}

	// Verify that the login process completed successfully
	// We can't directly check the sent packets without adding a GetSentPackets method to MockNetworkManager,
	// but we can verify that the session data was updated correctly, which indicates that
	// the packets were sent and processed correctly.
}
