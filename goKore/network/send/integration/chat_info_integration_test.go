package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/send/game/chat"
)

// TestChatInfoIntegration tests the integration of the chat info system
func TestChatInfoIntegration(t *testing.T) {
	mockSend := NewMockSend()
	// Update the packet IDs for the chat info system
	mockSend.packetIDs["actor_info_request"] = "0A90"
	mockSend.packetIDs["actor_name_request"] = "0095"
	mockSend.packetIDs["request_user_count"] = "01C1"
	mockSend.packetIDs["battleground_chat"] = "02DB"
	mockSend.packetIDs["clan_chat"] = "0B01"

	infoChatManager := chat.NewInfoChatManager(mockSend)

	// Test a sequence of chat info operations

	// 1. Send get player info
	ID := uint32(12345)
	err := infoChatManager.SendGetPlayerInfo(ID)
	if err != nil {
		t.Fatalf("SendGetPlayerInfo() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "actor_info_request" {
		t.Errorf("Expected packet ID 'actor_info_request', got '%s'", packetID)
	}

	// 2. Send get character name
	err = infoChatManager.SendGetCharacterName(ID)
	if err != nil {
		t.Fatalf("SendGetCharacterName() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "actor_name_request" {
		t.Errorf("Expected packet ID 'actor_name_request', got '%s'", packetID)
	}

	// 3. Send who
	err = infoChatManager.SendWho()
	if err != nil {
		t.Fatalf("SendWho() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_user_count" {
		t.Errorf("Expected packet ID 'request_user_count', got '%s'", packetID)
	}

	// 4. Send battleground chat
	message := "Hello, battleground!"
	err = infoChatManager.SendBattlegroundChat(message)
	if err != nil {
		t.Fatalf("SendBattlegroundChat() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "battleground_chat" {
		t.Errorf("Expected packet ID 'battleground_chat', got '%s'", packetID)
	}

	// 5. Send clan chat
	charName := "TestChar"
	err = infoChatManager.SendClanChat(message, charName)
	if err != nil {
		t.Fatalf("SendClanChat() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "clan_chat" {
		t.Errorf("Expected packet ID 'clan_chat', got '%s'", packetID)
	}

	// Verify the message format for clan chat
	expectedMessage := charName + " : " + message
	if messageVal, ok := mockSend.LastArgs()["message"].(string); !ok || messageVal != expectedMessage {
		t.Errorf("Expected message=%s, got %v", expectedMessage, mockSend.LastArgs()["message"])
	}
}
