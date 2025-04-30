package login

import (
	"bytes"
	"sync"
	"testing"
)

func TestSessionStore_InitialState(t *testing.T) {
	ss := NewSessionStore()

	// Check initial state
	sessionData := ss.GetSessionData()

	if sessionData.SessionID != nil {
		t.Errorf("Expected initial sessionID to be nil, got %v", sessionData.SessionID)
	}

	if sessionData.AccountID != nil {
		t.Errorf("Expected initial accountID to be nil, got %v", sessionData.AccountID)
	}

	if sessionData.SessionID2 != nil {
		t.Errorf("Expected initial sessionID2 to be nil, got %v", sessionData.SessionID2)
	}

	if sessionData.AccountSex != 0 {
		t.Errorf("Expected initial accountSex to be 0, got %v", sessionData.AccountSex)
	}
}

func TestSessionStore_UpdateFromAccountServerInfo(t *testing.T) {
	ss := NewSessionStore()

	// Test data
	sessionID := []byte{1, 2, 3, 4}
	accountID := []byte{5, 6, 7, 8}
	sessionID2 := []byte{9, 10, 11, 12}
	accountSex := 1

	// Update session data
	ss.UpdateFromAccountServerInfo(map[string]interface{}{
		"sessionID":  sessionID,
		"accountID":  accountID,
		"sessionID2": sessionID2,
		"accountSex": accountSex,
	})

	// Check updated state
	sessionData := ss.GetSessionData()

	if !bytes.Equal(sessionData.SessionID, sessionID) {
		t.Errorf("Expected sessionID to be %v, got %v", sessionID, sessionData.SessionID)
	}

	if !bytes.Equal(sessionData.AccountID, accountID) {
		t.Errorf("Expected accountID to be %v, got %v", accountID, sessionData.AccountID)
	}

	if !bytes.Equal(sessionData.SessionID2, sessionID2) {
		t.Errorf("Expected sessionID2 to be %v, got %v", sessionID2, sessionData.SessionID2)
	}

	if sessionData.AccountSex != accountSex {
		t.Errorf("Expected accountSex to be %v, got %v", accountSex, sessionData.AccountSex)
	}
}

func TestSessionStore_UpdateFromCharacterServerInfo(t *testing.T) {
	ss := NewSessionStore()

	// Test data
	charID := []byte{1, 2, 3, 4}
	mapName := "prontera"
	mapIP := "127.0.0.1"
	mapPort := 6900

	// Update session data
	ss.UpdateFromCharacterServerInfo(map[string]interface{}{
		"charID":  charID,
		"mapName": mapName,
		"mapIP":   mapIP,
		"mapPort": mapPort,
	})

	// Check updated state
	sessionData := ss.GetSessionData()

	if !bytes.Equal(sessionData.CharID, charID) {
		t.Errorf("Expected charID to be %v, got %v", charID, sessionData.CharID)
	}

	if sessionData.MapName != mapName {
		t.Errorf("Expected mapName to be %v, got %v", mapName, sessionData.MapName)
	}

	if sessionData.MapIP != mapIP {
		t.Errorf("Expected mapIP to be %v, got %v", mapIP, sessionData.MapIP)
	}

	if sessionData.MapPort != mapPort {
		t.Errorf("Expected mapPort to be %v, got %v", mapPort, sessionData.MapPort)
	}
}

func TestSessionStore_GetServerInfo(t *testing.T) {
	ss := NewSessionStore()

	// Test data
	serverInfo := []ServerInfo{
		{
			Name: "Server1",
			IP:   "127.0.0.1",
			Port: 6900,
		},
		{
			Name: "Server2",
			IP:   "127.0.0.2",
			Port: 6901,
		},
	}

	// Update server info
	ss.SetServerInfo(serverInfo)

	// Get server info by name
	server := ss.GetServerInfo("Server1")

	if server == nil {
		t.Fatalf("Expected server info for Server1, got nil")
	}

	if server.Name != "Server1" {
		t.Errorf("Expected server name to be Server1, got %v", server.Name)
	}

	if server.IP != "127.0.0.1" {
		t.Errorf("Expected server IP to be 127.0.0.1, got %v", server.IP)
	}

	if server.Port != 6900 {
		t.Errorf("Expected server port to be 6900, got %v", server.Port)
	}

	// Get non-existent server
	server = ss.GetServerInfo("NonExistentServer")

	if server != nil {
		t.Errorf("Expected nil for non-existent server, got %v", server)
	}
}

func TestSessionStore_ConcurrentAccess(t *testing.T) {
	ss := NewSessionStore()

	// Number of concurrent goroutines
	numGoroutines := 10

	// Use WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines that update session data
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine updates with different data
			sessionID := []byte{byte(id), byte(id + 1), byte(id + 2), byte(id + 3)}
			accountID := []byte{byte(id + 4), byte(id + 5), byte(id + 6), byte(id + 7)}

			ss.UpdateFromAccountServerInfo(map[string]interface{}{
				"sessionID": sessionID,
				"accountID": accountID,
			})

			// Read session data
			_ = ss.GetSessionData()
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Final state is non-deterministic due to race conditions,
	// but we can verify that the session store is still functional
	// Just get the data to ensure no panics occur
	_ = ss.GetSessionData()
}

func TestSessionStore_Reset(t *testing.T) {
	ss := NewSessionStore()

	// Update session data
	ss.UpdateFromAccountServerInfo(map[string]interface{}{
		"sessionID":  []byte{1, 2, 3, 4},
		"accountID":  []byte{5, 6, 7, 8},
		"sessionID2": []byte{9, 10, 11, 12},
		"accountSex": 1,
	})

	// Reset session data
	ss.Reset()

	// Check reset state
	sessionData := ss.GetSessionData()

	if sessionData.SessionID != nil {
		t.Errorf("Expected sessionID to be nil after reset, got %v", sessionData.SessionID)
	}

	if sessionData.AccountID != nil {
		t.Errorf("Expected accountID to be nil after reset, got %v", sessionData.AccountID)
	}

	if sessionData.SessionID2 != nil {
		t.Errorf("Expected sessionID2 to be nil after reset, got %v", sessionData.SessionID2)
	}

	if sessionData.AccountSex != 0 {
		t.Errorf("Expected accountSex to be 0 after reset, got %v", sessionData.AccountSex)
	}
}

func TestSessionStore_DeepCopy(t *testing.T) {
	ss := NewSessionStore()

	// Test data
	sessionID := []byte{1, 2, 3, 4}
	accountID := []byte{5, 6, 7, 8}

	// Update session data
	ss.UpdateFromAccountServerInfo(map[string]interface{}{
		"sessionID": sessionID,
		"accountID": accountID,
	})

	// Get session data
	sessionData := ss.GetSessionData()

	// Modify the original data
	sessionID[0] = 99
	accountID[0] = 99

	// Check that the returned data is a deep copy
	if sessionData.SessionID[0] == 99 {
		t.Errorf("Expected sessionID to be a deep copy, but it was modified")
	}

	if sessionData.AccountID[0] == 99 {
		t.Errorf("Expected accountID to be a deep copy, but it was modified")
	}
}
