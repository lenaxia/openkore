package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/send/game/misc"
)

// TestMiscIntegration tests the integration of the miscellaneous system
func TestMiscIntegration(t *testing.T) {
	mockSend := NewMockSend()
	// Update the packet IDs for the misc system
	mockSend.packetIDs["token_login"] = "0825"
	mockSend.packetIDs["request_remain_time"] = "0A37"
	mockSend.packetIDs["blocking_play_cancel"] = "0447"
	mockSend.packetIDs["recall_sso"] = "0842"
	mockSend.packetIDs["remove_aid_sso"] = "0843"
	mockSend.packetIDs["starplace_agree"] = "0B0D"
	mockSend.packetIDs["sync_request_ex"] = "09F1"

	miscManager := misc.NewMiscManager(mockSend)

	// Test a sequence of miscellaneous operations

	// 1. Send token to server
	username := "testuser"
	password := "testpass"
	masterVersion := uint32(1)
	version := uint32(2)
	token := "testtoken"
	length := uint16(10)
	otpIP := "127.0.0.1"
	otpPort := uint16(6900)
	err := miscManager.SendTokenToServer(username, password, masterVersion, version, token, length, otpIP, otpPort)
	if err != nil {
		t.Fatalf("SendTokenToServer() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "token_login" {
		t.Errorf("Expected packet ID 'token_login', got '%s'", packetID)
	}

	// 2. Send request remain time
	err = miscManager.SendReqRemainTime()
	if err != nil {
		t.Fatalf("SendReqRemainTime() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_remain_time" {
		t.Errorf("Expected packet ID 'request_remain_time', got '%s'", packetID)
	}

	// 3. Send blocking player cancel
	err = miscManager.SendBlockingPlayerCancel()
	if err != nil {
		t.Fatalf("SendBlockingPlayerCancel() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "blocking_play_cancel" {
		t.Errorf("Expected packet ID 'blocking_play_cancel', got '%s'", packetID)
	}

	// 4. Send recall SSO
	accountID := uint32(12345)
	err = miscManager.SendRecallSso(accountID)
	if err != nil {
		t.Fatalf("SendRecallSso() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "recall_sso" {
		t.Errorf("Expected packet ID 'recall_sso', got '%s'", packetID)
	}

	// 5. Send remove AID SSO
	err = miscManager.SendRemoveAidSso(accountID)
	if err != nil {
		t.Fatalf("SendRemoveAidSso() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "remove_aid_sso" {
		t.Errorf("Expected packet ID 'remove_aid_sso', got '%s'", packetID)
	}

	// 6. Send feel save ok
	flag := uint8(1)
	err = miscManager.SendFeelSaveOk(flag)
	if err != nil {
		t.Fatalf("SendFeelSaveOk() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "starplace_agree" {
		t.Errorf("Expected packet ID 'starplace_agree', got '%s'", packetID)
	}

	// 7. Send reply sync request ex
	syncID := uint16(12345)
	err = miscManager.SendReplySyncRequestEx(syncID)
	if err != nil {
		t.Fatalf("SendReplySyncRequestEx() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "sync_request_ex" {
		t.Errorf("Expected packet ID 'sync_request_ex', got '%s'", packetID)
	}
}
