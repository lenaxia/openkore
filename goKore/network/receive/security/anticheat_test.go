package security

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestNewAntiCheatManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewAntiCheatManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewAntiCheatManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if manager.state != AntiCheatStateDisabled {
		t.Errorf("manager.state = %v, want %v", manager.state, AntiCheatStateDisabled)
	}

	if manager.antiCheatType != AntiCheatNone {
		t.Errorf("manager.antiCheatType = %v, want %v", manager.antiCheatType, AntiCheatNone)
	}

	if manager.timeout != 30*time.Second {
		t.Errorf("manager.timeout = %v, want %v", manager.timeout, 30*time.Second)
	}

	if manager.enabled {
		t.Error("manager.enabled = true, want false")
	}

	if manager.rng == nil {
		t.Error("manager.rng was not initialized")
	}
}

func TestAntiCheatRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewAntiCheatManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"gameguard_request",
		"gameguard_lingo",
		"gameguard_reply",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestEnableDisable(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Check initial state
	if manager.IsEnabled() {
		t.Error("IsEnabled() = true, want false initially")
	}

	// Enable
	manager.Enable(AntiCheatGameGuard)

	// Check state
	if !manager.IsEnabled() {
		t.Error("IsEnabled() = false, want true after Enable()")
	}

	if manager.GetAntiCheatType() != AntiCheatGameGuard {
		t.Errorf("GetAntiCheatType() = %v, want %v", manager.GetAntiCheatType(), AntiCheatGameGuard)
	}

	if manager.GetState() != AntiCheatStateInitializing {
		t.Errorf("GetState() = %v, want %v", manager.GetState(), AntiCheatStateInitializing)
	}

	// Disable
	manager.Disable()

	// Check state
	if manager.IsEnabled() {
		t.Error("IsEnabled() = true, want false after Disable()")
	}

	if manager.GetState() != AntiCheatStateDisabled {
		t.Errorf("GetState() = %v, want %v", manager.GetState(), AntiCheatStateDisabled)
	}
}

func TestAntiCheatSetGetState(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Set state
	manager.SetState(AntiCheatStateVerified)

	// Get state
	state := manager.GetState()
	if state != AntiCheatStateVerified {
		t.Errorf("GetState() = %v, want %v", state, AntiCheatStateVerified)
	}
}

func TestSetGetTimeout(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Set timeout
	timeout := 60 * time.Second
	manager.SetTimeout(timeout)

	// Get timeout
	if manager.GetTimeout() != timeout {
		t.Errorf("GetTimeout() = %v, want %v", manager.GetTimeout(), timeout)
	}
}

func TestAntiCheatIsVerified(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with state = AntiCheatStateVerified
	manager.state = AntiCheatStateVerified
	if !manager.IsVerified() {
		t.Error("IsVerified() = false, want true when state is AntiCheatStateVerified")
	}

	// Test with state != AntiCheatStateVerified
	manager.state = AntiCheatStateInitializing
	if manager.IsVerified() {
		t.Error("IsVerified() = true, want false when state != AntiCheatStateVerified")
	}
}

func TestIsRejected(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with state = AntiCheatStateRejected
	manager.state = AntiCheatStateRejected
	if !manager.IsRejected() {
		t.Error("IsRejected() = false, want true when state is AntiCheatStateRejected")
	}

	// Test with state != AntiCheatStateRejected
	manager.state = AntiCheatStateInitializing
	if manager.IsRejected() {
		t.Error("IsRejected() = true, want false when state != AntiCheatStateRejected")
	}
}

func TestIsTimedOut(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with state != AntiCheatStateWaitingForResponse
	manager.state = AntiCheatStateInitializing
	if manager.IsTimedOut() {
		t.Error("IsTimedOut() = true, want false when state != AntiCheatStateWaitingForResponse")
	}

	// Test with state = AntiCheatStateWaitingForResponse and not timed out
	manager.state = AntiCheatStateWaitingForResponse
	manager.lastChallenge = time.Now()
	manager.timeout = 30 * time.Second
	if manager.IsTimedOut() {
		t.Error("IsTimedOut() = true, want false when not timed out")
	}

	// Test with state = AntiCheatStateWaitingForResponse and timed out
	manager.lastChallenge = time.Now().Add(-60 * time.Second)
	if !manager.IsTimedOut() {
		t.Error("IsTimedOut() = false, want true when timed out")
	}
}

func TestGenerateChallenge(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with different anti-cheat types
	testCases := []struct {
		antiCheatType AntiCheatType
		expectedLen   int
	}{
		{AntiCheatGameGuard, 20},
		{AntiCheatXTrap, 16},
		{AntiCheatHShield, 32},
		{AntiCheatNProtect, 24},
		{AntiCheatNone, 16},
	}

	for _, tc := range testCases {
		manager.antiCheatType = tc.antiCheatType
		challenge := manager.GenerateChallenge()

		if len(challenge) != tc.expectedLen {
			t.Errorf("len(GenerateChallenge()) = %d, want %d for antiCheatType %v", len(challenge), tc.expectedLen, tc.antiCheatType)
		}

		if manager.state != AntiCheatStateWaitingForResponse {
			t.Errorf("manager.state = %v, want %v after GenerateChallenge()", manager.state, AntiCheatStateWaitingForResponse)
		}
	}
}

func TestVerifyResponse(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with anti-cheat disabled
	manager.enabled = false
	err := manager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrAntiCheatDisabled {
		t.Errorf("VerifyResponse() with anti-cheat disabled returned error: %v, want %v", err, ErrAntiCheatDisabled)
	}

	// Test with state != AntiCheatStateWaitingForResponse
	manager.enabled = true
	manager.state = AntiCheatStateInitializing
	err = manager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrInvalidResponse {
		t.Errorf("VerifyResponse() with state != AntiCheatStateWaitingForResponse returned error: %v, want %v", err, ErrInvalidResponse)
	}

	// Test with timed out response
	manager.state = AntiCheatStateWaitingForResponse
	manager.lastChallenge = time.Now().Add(-60 * time.Second)
	manager.timeout = 30 * time.Second
	err = manager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrAntiCheatTimeout {
		t.Errorf("VerifyResponse() with timed out response returned error: %v, want %v", err, ErrAntiCheatTimeout)
	}

	// Test with valid response
	manager.state = AntiCheatStateWaitingForResponse
	manager.lastChallenge = time.Now()
	manager.antiCheatType = AntiCheatNone
	manager.challenge = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	err = manager.VerifyResponse([]byte{0x00, 0x00, 0x00, 0x00})
	if err != nil {
		t.Errorf("VerifyResponse() with valid response returned error: %v", err)
	}

	if manager.state != AntiCheatStateVerified {
		t.Errorf("manager.state = %v, want %v after successful verification", manager.state, AntiCheatStateVerified)
	}

	// Test with invalid response
	manager.state = AntiCheatStateWaitingForResponse
	manager.lastChallenge = time.Now()
	manager.challenge = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	err = manager.VerifyResponse([]byte{0x01, 0x01, 0x01, 0x01})
	if err != ErrInvalidResponse {
		t.Errorf("VerifyResponse() with invalid response returned error: %v, want %v", err, ErrInvalidResponse)
	}

	if manager.state != AntiCheatStateRejected {
		t.Errorf("manager.state = %v, want %v after failed verification", manager.state, AntiCheatStateRejected)
	}
}

func TestHandleGameGuardRequest(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewAntiCheatManager(parser, hookManager)

	// Call handler
	err := manager.handleGameGuardRequest(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleGameGuardRequest() returned error: %v", err)
	}

	// Check state
	if !manager.enabled {
		t.Error("manager.enabled = false, want true after handling GameGuard request")
	}

	if manager.antiCheatType != AntiCheatGameGuard {
		t.Errorf("manager.antiCheatType = %v, want %v after handling GameGuard request", manager.antiCheatType, AntiCheatGameGuard)
	}

	if manager.state != AntiCheatStateWaitingForChallenge {
		t.Errorf("manager.state = %v, want %v after handling GameGuard request", manager.state, AntiCheatStateWaitingForChallenge)
	}
}

func TestGenerateResponse(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with different anti-cheat types
	testCases := []struct {
		antiCheatType AntiCheatType
		expectedLen   int
	}{
		{AntiCheatGameGuard, 16},
		{AntiCheatXTrap, 12},
		{AntiCheatHShield, 24},
		{AntiCheatNProtect, 20},
		{AntiCheatNone, 4},
	}

	challenge := []byte{1, 2, 3, 4}

	for _, tc := range testCases {
		manager.antiCheatType = tc.antiCheatType
		response := manager.GenerateResponse(challenge)

		if tc.antiCheatType == AntiCheatNone {
			// For AntiCheatNone, the response length should match the challenge length
			if len(response) != len(challenge) {
				t.Errorf("len(GenerateResponse()) = %d, want %d for antiCheatType %v", len(response), len(challenge), tc.antiCheatType)
			}
		} else {
			// For other types, check the expected length
			if len(response) != tc.expectedLen {
				t.Errorf("len(GenerateResponse()) = %d, want %d for antiCheatType %v", len(response), tc.expectedLen, tc.antiCheatType)
			}
		}
	}
}

func TestCalculateChecksum(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewAntiCheatManager(parser, nil)

	// Test with a simple buffer
	buffer := []byte("test")
	checksum := manager.CalculateChecksum(buffer)

	// MD5 of "test" is "098f6bcd4621d373cade4e832627b4f6"
	expectedChecksum := []byte{0x09, 0x8f, 0x6b, 0xcd, 0x46, 0x21, 0xd3, 0x73, 0xca, 0xde, 0x4e, 0x83, 0x26, 0x27, 0xb4, 0xf6}

	if len(checksum) != len(expectedChecksum) {
		t.Fatalf("len(CalculateChecksum()) = %d, want %d", len(checksum), len(expectedChecksum))
	}

	for i := range checksum {
		if checksum[i] != expectedChecksum[i] {
			t.Errorf("CalculateChecksum()[%d] = %02x, want %02x", i, checksum[i], expectedChecksum[i])
		}
	}
}

func TestAntiCheatTypeString(t *testing.T) {
	testCases := []struct {
		antiCheatType AntiCheatType
		want          string
	}{
		{AntiCheatNone, "None"},
		{AntiCheatGameGuard, "GameGuard"},
		{AntiCheatXTrap, "XTrap"},
		{AntiCheatHShield, "HShield"},
		{AntiCheatNProtect, "NProtect"},
		{AntiCheatCustom, "Custom"},
		{AntiCheatType(99), "Unknown"},
	}

	for _, tc := range testCases {
		got := tc.antiCheatType.String()
		if got != tc.want {
			t.Errorf("AntiCheatType(%d).String() = %s, want %s", tc.antiCheatType, got, tc.want)
		}
	}
}

func TestAntiCheatStateString(t *testing.T) {
	testCases := []struct {
		antiCheatState AntiCheatState
		want           string
	}{
		{AntiCheatStateDisabled, "Disabled"},
		{AntiCheatStateInitializing, "Initializing"},
		{AntiCheatStateWaitingForChallenge, "WaitingForChallenge"},
		{AntiCheatStateWaitingForResponse, "WaitingForResponse"},
		{AntiCheatStateVerified, "Verified"},
		{AntiCheatStateRejected, "Rejected"},
		{AntiCheatState(99), "Unknown"},
	}

	for _, tc := range testCases {
		got := tc.antiCheatState.String()
		if got != tc.want {
			t.Errorf("AntiCheatState(%d).String() = %s, want %s", tc.antiCheatState, got, tc.want)
		}
	}
}
