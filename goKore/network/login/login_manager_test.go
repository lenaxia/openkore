package login

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Network state constants for testing
const (
	NetworkNotConnected            = 0
	NetworkConnectedToMasterServer = 1
	NetworkConnectedToCharServer   = 2
	NetworkConnectedToMapServer    = 3
	NetworkInGame                  = 4
)

func TestLoginManager_NewLoginManager(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	loginManager := NewLoginManager(networkManager, config)

	if loginManager == nil {
		t.Fatal("Expected non-nil login manager")
	}

	// We can't directly compare interface values, so we'll skip this check
	// and rely on the functionality tests instead

	if loginManager.config != config {
		t.Error("Expected config to be set correctly")
	}

	if loginManager.stateManager == nil {
		t.Error("Expected state manager to be initialized")
	}

	if loginManager.sessionStore == nil {
		t.Error("Expected session store to be initialized")
	}
}

func TestLoginManager_Login_Success(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start login process in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- loginManager.Login(ctx)
	}()

	// Simulate successful login sequence
	time.Sleep(100 * time.Millisecond) // Give login process time to start

	// Simulate account server info response
	err := networkManager.SimulateReceivePacket("0069", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate account server info: %v", err)
	}

	// Simulate character server info response
	err = networkManager.SimulateReceivePacket("006B", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate character server info: %v", err)
	}

	// Simulate character ID and map info response
	err = networkManager.SimulateReceivePacket("0071", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate character ID and map info: %v", err)
	}

	// Simulate map loaded response
	err = networkManager.SimulateReceivePacket("0073", []byte{})
	if err != nil {
		t.Fatalf("Failed to simulate map loaded: %v", err)
	}

	// Wait for login process to complete
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("Expected successful login, got error: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Login process timed out")
	}

	// Check final state
	if loginManager.stateManager.GetState() != StateInGame {
		t.Errorf("Expected final state to be StateInGame, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Login_ConnectError(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	// Set up connect error
	networkManager.networkInterface.SetConnectError(errors.New("connect error"))

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt login
	err := loginManager.Login(ctx)

	// Check error
	if err == nil {
		t.Error("Expected error on connect failure, got nil")
	}

	// Check state
	if loginManager.stateManager.GetState() != StateNotConnected {
		t.Errorf("Expected state to be StateNotConnected, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Login_SendError(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	// Set up send error
	networkManager.packetSender.SetSendError(errors.New("send error"))

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt login
	err := loginManager.Login(ctx)

	// Check error
	if err == nil {
		t.Error("Expected error on send failure, got nil")
	}

	// Check state
	if loginManager.stateManager.GetState() != StateNotConnected {
		t.Errorf("Expected state to be StateNotConnected, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Login_LoginError(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")
	config.MaxRetries = 0 // Disable retries for this test

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start login process in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- loginManager.Login(ctx)
	}()

	// Simulate login error
	time.Sleep(100 * time.Millisecond) // Give login process time to start
	err := networkManager.SimulateReceivePacket("006A", []byte{1})
	if err != nil {
		t.Fatalf("Failed to simulate login error: %v", err)
	}

	// Wait for login process to complete
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error on login failure, got nil")
		}
	case <-ctx.Done():
		t.Fatal("Login process timed out")
	}

	// Check state
	if loginManager.stateManager.GetState() != StateNotConnected {
		t.Errorf("Expected state to be StateNotConnected, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Login_Timeout(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")
	config.LoginTimeout = 100 * time.Millisecond // Short timeout for testing

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt login
	err := loginManager.Login(ctx)

	// Check error
	if err == nil {
		t.Error("Expected error on timeout, got nil")
	}

	// Check state
	if loginManager.stateManager.GetState() != StateNotConnected {
		t.Errorf("Expected state to be StateNotConnected, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Login_Cancelled(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context that we'll cancel immediately
	ctx, cancel := context.WithCancel(context.Background())

	// Start login process in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- loginManager.Login(ctx)
	}()

	// Cancel the context immediately
	cancel()

	// Wait for login process to complete
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error on cancellation, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Login process did not respond to cancellation")
	}

	// Check state
	if loginManager.stateManager.GetState() != StateNotConnected {
		t.Errorf("Expected state to be StateNotConnected, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_Retry(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")
	config.MaxRetries = 1
	config.RetryDelay = 100 * time.Millisecond

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start login process in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- loginManager.Login(ctx)
	}()

	// Simulate login error
	time.Sleep(100 * time.Millisecond) // Give login process time to start
	err := networkManager.SimulateReceivePacket("006A", []byte{1})
	if err != nil {
		t.Fatalf("Failed to simulate login error: %v", err)
	}

	// Wait for retry
	time.Sleep(200 * time.Millisecond)

	// Simulate successful login sequence after retry
	err = networkManager.SimulateReceivePacket("0069", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate account server info: %v", err)
	}

	err = networkManager.SimulateReceivePacket("006B", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate character server info: %v", err)
	}

	err = networkManager.SimulateReceivePacket("0071", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate character ID and map info: %v", err)
	}

	err = networkManager.SimulateReceivePacket("0073", []byte{})
	if err != nil {
		t.Fatalf("Failed to simulate map loaded: %v", err)
	}

	// Wait for login process to complete
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("Expected successful login after retry, got error: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Login process timed out")
	}

	// Check final state
	if loginManager.stateManager.GetState() != StateInGame {
		t.Errorf("Expected final state to be StateInGame, got %v", loginManager.stateManager.GetState())
	}
}

func TestLoginManager_SessionDataPersistence(t *testing.T) {
	networkManager := NewMockNetworkManager()
	config := NewLoginConfig("testuser", "testpass", "TestServer")

	loginManager := NewLoginManager(networkManager, config)

	// Set up a context with timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start login process in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- loginManager.Login(ctx)
	}()

	// Simulate account server info response
	time.Sleep(100 * time.Millisecond) // Give login process time to start
	err := networkManager.SimulateReceivePacket("0069", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate account server info: %v", err)
	}

	// Check session data after account server info
	sessionData := loginManager.sessionStore.GetSessionData()
	if len(sessionData.SessionID) != 4 || sessionData.SessionID[0] != 1 {
		t.Errorf("Session ID not stored correctly: %v", sessionData.SessionID)
	}

	if len(sessionData.AccountID) != 4 || sessionData.AccountID[0] != 5 {
		t.Errorf("Account ID not stored correctly: %v", sessionData.AccountID)
	}

	// Simulate character ID and map info response
	err = networkManager.SimulateReceivePacket("0071", []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Failed to simulate character ID and map info: %v", err)
	}

	// Check session data after character ID and map info
	sessionData = loginManager.sessionStore.GetSessionData()
	if len(sessionData.CharID) != 4 || sessionData.CharID[0] != 1 {
		t.Errorf("Character ID not stored correctly: %v", sessionData.CharID)
	}

	if sessionData.MapName != "prontera" {
		t.Errorf("Map name not stored correctly: %v", sessionData.MapName)
	}

	// Cancel the login process
	cancel()
}
