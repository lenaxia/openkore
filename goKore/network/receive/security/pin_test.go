package security

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestNewPINManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewPINManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewPINManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if manager.state != PINStateUnknown {
		t.Errorf("manager.state = %v, want %v", manager.state, PINStateUnknown)
	}

	if manager.maxAttempts != 3 {
		t.Errorf("manager.maxAttempts = %d, want 3", manager.maxAttempts)
	}
}

func TestPINRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewPINManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"login_pin_code_request",
		"login_pin_new_code_result",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestSetGetPIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set PIN
	err := manager.SetPIN("1234")
	if err != nil {
		t.Fatalf("SetPIN() returned error: %v", err)
	}

	// Get PIN
	pin := manager.GetPIN()
	if pin != "1234" {
		t.Errorf("GetPIN() = %s, want 1234", pin)
	}

	// Check state
	if manager.GetState() != PINStateSet {
		t.Errorf("manager.GetState() = %v, want %v", manager.GetState(), PINStateSet)
	}

	// Test invalid PIN (wrong length)
	err = manager.SetPIN("123")
	if err != ErrPINWrongLength {
		t.Errorf("SetPIN() with wrong length returned error: %v, want %v", err, ErrPINWrongLength)
	}

	// Test invalid PIN (non-digits)
	err = manager.SetPIN("123a")
	if err != ErrPINInvalidFormat {
		t.Errorf("SetPIN() with non-digits returned error: %v, want %v", err, ErrPINInvalidFormat)
	}
}

func TestPINSetGetState(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set state
	manager.SetState(PINStateVerified)

	// Get state
	state := manager.GetState()
	if state != PINStateVerified {
		t.Errorf("GetState() = %v, want %v", state, PINStateVerified)
	}
}

func TestGetSeed(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set seed
	manager.seed = 12345

	// Get seed
	seed := manager.GetSeed()
	if seed != 12345 {
		t.Errorf("GetSeed() = %d, want 12345", seed)
	}
}

func TestGetAccountID(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set account ID
	manager.accountID = 67890

	// Get account ID
	accountID := manager.GetAccountID()
	if accountID != 67890 {
		t.Errorf("GetAccountID() = %d, want 67890", accountID)
	}
}

func TestGetFlag(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set flag
	manager.flag = 1

	// Get flag
	flag := manager.GetFlag()
	if flag != 1 {
		t.Errorf("GetFlag() = %d, want 1", flag)
	}
}

func TestGetLock(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set lock
	manager.lock = 2

	// Get lock
	lock := manager.GetLock()
	if lock != 2 {
		t.Errorf("GetLock() = %d, want 2", lock)
	}
}

func TestIsLocked(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Test with state = PINStateLocked
	manager.state = PINStateLocked
	manager.lock = 0
	if !manager.IsLocked() {
		t.Error("IsLocked() = false, want true when state is PINStateLocked")
	}

	// Test with lock > 0
	manager.state = PINStateSet
	manager.lock = 1
	if !manager.IsLocked() {
		t.Error("IsLocked() = false, want true when lock > 0")
	}

	// Test with state != PINStateLocked and lock = 0
	manager.state = PINStateSet
	manager.lock = 0
	if manager.IsLocked() {
		t.Error("IsLocked() = true, want false when state != PINStateLocked and lock = 0")
	}
}

func TestIsVerified(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Test with state = PINStateVerified
	manager.state = PINStateVerified
	if !manager.IsVerified() {
		t.Error("IsVerified() = false, want true when state is PINStateVerified")
	}

	// Test with state != PINStateVerified
	manager.state = PINStateSet
	if manager.IsVerified() {
		t.Error("IsVerified() = true, want false when state != PINStateVerified")
	}
}

func TestIsRequired(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Test with state = PINStateRequested
	manager.state = PINStateRequested
	manager.flag = 0
	if !manager.IsRequired() {
		t.Error("IsRequired() = false, want true when state is PINStateRequested")
	}

	// Test with flag = 1
	manager.state = PINStateSet
	manager.flag = 1
	if !manager.IsRequired() {
		t.Error("IsRequired() = false, want true when flag = 1")
	}

	// Test with flag = 3
	manager.state = PINStateSet
	manager.flag = 3
	if !manager.IsRequired() {
		t.Error("IsRequired() = false, want true when flag = 3")
	}

	// Test with flag = 5
	manager.state = PINStateSet
	manager.flag = 5
	if !manager.IsRequired() {
		t.Error("IsRequired() = false, want true when flag = 5")
	}

	// Test with state != PINStateRequested and flag not in [1, 3, 5]
	manager.state = PINStateSet
	manager.flag = 0
	if manager.IsRequired() {
		t.Error("IsRequired() = true, want false when state != PINStateRequested and flag not in [1, 3, 5]")
	}
}

func TestVerifyPIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set PIN and state
	manager.pin = "1234"
	manager.state = PINStateRequested

	// Test with correct PIN
	err := manager.VerifyPIN("1234")
	if err != nil {
		t.Errorf("VerifyPIN() with correct PIN returned error: %v", err)
	}

	// Check state
	if manager.state != PINStateVerified {
		t.Errorf("manager.state = %v, want %v after successful verification", manager.state, PINStateVerified)
	}

	// Reset state
	manager.state = PINStateRequested
	manager.attempts = 0

	// Test with incorrect PIN
	err = manager.VerifyPIN("5678")
	if err != ErrInvalidPIN {
		t.Errorf("VerifyPIN() with incorrect PIN returned error: %v, want %v", err, ErrInvalidPIN)
	}

	// Check attempts
	if manager.attempts != 1 {
		t.Errorf("manager.attempts = %d, want 1 after failed verification", manager.attempts)
	}

	// Test with locked PIN
	manager.state = PINStateLocked
	err = manager.VerifyPIN("1234")
	if err != ErrPINLocked {
		t.Errorf("VerifyPIN() with locked PIN returned error: %v, want %v", err, ErrPINLocked)
	}

	// Test with PIN not required
	manager.state = PINStateSet
	err = manager.VerifyPIN("1234")
	if err != ErrPINNotSet {
		t.Errorf("VerifyPIN() with PIN not required returned error: %v, want %v", err, ErrPINNotSet)
	}

	// Test with invalid PIN format
	manager.state = PINStateRequested
	err = manager.VerifyPIN("123a")
	if err != ErrPINInvalidFormat {
		t.Errorf("VerifyPIN() with invalid PIN format returned error: %v, want %v", err, ErrPINInvalidFormat)
	}

	// Test with max attempts reached
	manager.state = PINStateRequested
	manager.attempts = 2
	manager.maxAttempts = 3
	err = manager.VerifyPIN("5678")
	if err != ErrPINLocked {
		t.Errorf("VerifyPIN() with max attempts reached returned error: %v, want %v", err, ErrPINLocked)
	}

	// Check state
	if manager.state != PINStateLocked {
		t.Errorf("manager.state = %v, want %v after max attempts reached", manager.state, PINStateLocked)
	}
}

func TestChangePIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set PIN and state
	manager.pin = "1234"
	manager.state = PINStateSet

	// Test with correct old PIN and valid new PIN
	err := manager.ChangePIN("1234", "5678")
	if err != nil {
		t.Errorf("ChangePIN() with correct old PIN and valid new PIN returned error: %v", err)
	}

	// Check PIN
	if manager.pin != "5678" {
		t.Errorf("manager.pin = %s, want 5678 after successful change", manager.pin)
	}

	// Test with incorrect old PIN
	manager.pin = "5678"
	manager.attempts = 0
	err = manager.ChangePIN("1234", "9012")
	if err != ErrInvalidPIN {
		t.Errorf("ChangePIN() with incorrect old PIN returned error: %v, want %v", err, ErrInvalidPIN)
	}

	// Check attempts
	if manager.attempts != 1 {
		t.Errorf("manager.attempts = %d, want 1 after failed change", manager.attempts)
	}

	// Test with locked PIN
	manager.state = PINStateLocked
	err = manager.ChangePIN("5678", "9012")
	if err != ErrPINLocked {
		t.Errorf("ChangePIN() with locked PIN returned error: %v, want %v", err, ErrPINLocked)
	}

	// Test with invalid new PIN
	manager.state = PINStateSet
	err = manager.ChangePIN("5678", "901")
	if err != ErrPINWrongLength {
		t.Errorf("ChangePIN() with invalid new PIN returned error: %v, want %v", err, ErrPINWrongLength)
	}

	// Test with max attempts reached
	manager.state = PINStateSet
	manager.attempts = 2
	manager.maxAttempts = 3
	err = manager.ChangePIN("1234", "9012")
	if err != ErrPINLocked {
		t.Errorf("ChangePIN() with max attempts reached returned error: %v, want %v", err, ErrPINLocked)
	}

	// Check state
	if manager.state != PINStateLocked {
		t.Errorf("manager.state = %v, want %v after max attempts reached", manager.state, PINStateLocked)
	}
}

func TestResetPIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set PIN, state, and attempts
	manager.pin = "1234"
	manager.state = PINStateSet
	manager.attempts = 2

	// Reset PIN
	manager.ResetPIN()

	// Check PIN
	if manager.pin != "" {
		t.Errorf("manager.pin = %s, want empty string after reset", manager.pin)
	}

	// Check state
	if manager.state != PINStateNotSet {
		t.Errorf("manager.state = %v, want %v after reset", manager.state, PINStateNotSet)
	}

	// Check attempts
	if manager.attempts != 0 {
		t.Errorf("manager.attempts = %d, want 0 after reset", manager.attempts)
	}
}

func TestUnlockPIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set state and attempts
	manager.state = PINStateLocked
	manager.attempts = 3

	// Unlock PIN
	manager.UnlockPIN()

	// Check state
	if manager.state != PINStateSet {
		t.Errorf("manager.state = %v, want %v after unlock", manager.state, PINStateSet)
	}

	// Check attempts
	if manager.attempts != 0 {
		t.Errorf("manager.attempts = %d, want 0 after unlock", manager.attempts)
	}

	// Test with non-locked state
	manager.state = PINStateSet
	manager.attempts = 1
	manager.UnlockPIN()

	// Check state (should remain unchanged)
	if manager.state != PINStateSet {
		t.Errorf("manager.state = %v, want %v after unlock with non-locked state", manager.state, PINStateSet)
	}

	// Check attempts (should remain unchanged)
	if manager.attempts != 1 {
		t.Errorf("manager.attempts = %d, want 1 after unlock with non-locked state", manager.attempts)
	}
}

func TestHandleLoginPinCodeRequest(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewPINManager(parser, hookManager)

	// Test with flag = 0 (PIN is correct)
	args := map[string]interface{}{
		"seed":      uint32(12345),
		"accountID": []byte{1, 2, 3, 4},
		"flag":      uint16(0),
		"lock":      uint16(0),
	}

	err := manager.handleLoginPinCodeRequest(args)
	if err != nil {
		t.Fatalf("handleLoginPinCodeRequest() returned error: %v", err)
	}

	// Check state
	if manager.state != PINStateVerified {
		t.Errorf("manager.state = %v, want %v after handling flag = 0", manager.state, PINStateVerified)
	}

	// Test with flag = 1 (PIN is required)
	args["flag"] = uint16(1)
	err = manager.handleLoginPinCodeRequest(args)
	if err != nil {
		t.Fatalf("handleLoginPinCodeRequest() returned error: %v", err)
	}

	// Check state
	if manager.state != PINStateRequested {
		t.Errorf("manager.state = %v, want %v after handling flag = 1", manager.state, PINStateRequested)
	}

	// Test with lock > 0
	args["flag"] = uint16(0)
	args["lock"] = uint16(1)
	err = manager.handleLoginPinCodeRequest(args)
	if err != nil {
		t.Fatalf("handleLoginPinCodeRequest() returned error: %v", err)
	}

	// Check state
	if manager.state != PINStateLocked {
		t.Errorf("manager.state = %v, want %v after handling lock > 0", manager.state, PINStateLocked)
	}
}

func TestHandleLoginPinNewCodeResult(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewPINManager(parser, hookManager)

	// Test with flag = 0 (PIN successfully changed)
	args := map[string]interface{}{
		"flag": uint16(0),
		"seed": uint32(12345),
	}

	err := manager.handleLoginPinNewCodeResult(args)
	if err != nil {
		t.Fatalf("handleLoginPinNewCodeResult() returned error: %v", err)
	}

	// Check state
	if manager.state != PINStateSet {
		t.Errorf("manager.state = %v, want %v after handling flag = 0", manager.state, PINStateSet)
	}

	// Check seed
	if manager.seed != 12345 {
		t.Errorf("manager.seed = %d, want 12345 after handling", manager.seed)
	}

	// Test with flag = 1 (PIN change failed)
	manager.state = PINStateChanging
	args["flag"] = uint16(1)
	err = manager.handleLoginPinNewCodeResult(args)
	if err != nil {
		t.Fatalf("handleLoginPinNewCodeResult() returned error: %v", err)
	}

	// Check state (should remain unchanged)
	if manager.state != PINStateChanging {
		t.Errorf("manager.state = %v, want %v after handling flag = 1", manager.state, PINStateChanging)
	}
}

func TestValidatePIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Test valid PIN
	err := manager.validatePIN("1234")
	if err != nil {
		t.Errorf("validatePIN() with valid PIN returned error: %v", err)
	}

	// Test PIN with wrong length
	err = manager.validatePIN("123")
	if err != ErrPINWrongLength {
		t.Errorf("validatePIN() with wrong length returned error: %v, want %v", err, ErrPINWrongLength)
	}

	// Test PIN with non-digits
	err = manager.validatePIN("123a")
	if err != ErrPINInvalidFormat {
		t.Errorf("validatePIN() with non-digits returned error: %v, want %v", err, ErrPINInvalidFormat)
	}
}

func TestEncryptDecryptPIN(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewPINManager(parser, nil)

	// Set seed
	manager.seed = 12345

	// Test encrypt/decrypt
	pin := "1234"
	encrypted, err := manager.EncryptPIN(pin)
	if err != nil {
		t.Fatalf("EncryptPIN() returned error: %v", err)
	}

	decrypted, err := manager.DecryptPIN(encrypted)
	if err != nil {
		t.Fatalf("DecryptPIN() returned error: %v", err)
	}

	if decrypted != pin {
		t.Errorf("DecryptPIN(EncryptPIN(%s)) = %s, want %s", pin, decrypted, pin)
	}

	// Test encrypt with invalid PIN
	_, err = manager.EncryptPIN("123")
	if err != ErrPINWrongLength {
		t.Errorf("EncryptPIN() with invalid PIN returned error: %v, want %v", err, ErrPINWrongLength)
	}

	// Test decrypt with invalid encrypted PIN
	_, err = manager.DecryptPIN("123")
	if err == nil || err.Error() != "invalid encrypted PIN length" {
		t.Errorf("DecryptPIN() with invalid encrypted PIN returned error: %v, want 'invalid encrypted PIN length'", err)
	}
}

func TestPINStateString(t *testing.T) {
	testCases := []struct {
		state PINState
		want  string
	}{
		{PINStateUnknown, "Unknown"},
		{PINStateNotSet, "NotSet"},
		{PINStateSet, "Set"},
		{PINStateRequested, "Requested"},
		{PINStateVerifying, "Verifying"},
		{PINStateVerified, "Verified"},
		{PINStateChanging, "Changing"},
		{PINStateLocked, "Locked"},
		{PINState(99), "Invalid"},
	}

	for _, tc := range testCases {
		got := tc.state.String()
		if got != tc.want {
			t.Errorf("PINState(%d).String() = %s, want %s", tc.state, got, tc.want)
		}
	}
}
