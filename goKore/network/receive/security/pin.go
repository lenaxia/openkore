// Package security provides security-related functionality for the network stack.
package security

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Errors
var (
	ErrInvalidPIN       = errors.New("invalid PIN code")
	ErrPINRequired      = errors.New("PIN code required")
	ErrPINAlreadySet    = errors.New("PIN code already set")
	ErrPINNotSet        = errors.New("PIN code not set")
	ErrPINLocked        = errors.New("PIN code locked")
	ErrPINWrongLength   = errors.New("PIN code must be 4 digits")
	ErrPINInvalidFormat = errors.New("PIN code must contain only digits")
)

// PINState represents the state of the PIN code
type PINState int

// PIN states
const (
	PINStateUnknown PINState = iota
	PINStateNotSet
	PINStateSet
	PINStateRequested
	PINStateVerifying
	PINStateVerified
	PINStateChanging
	PINStateLocked
)

// String returns the string representation of the PIN state
func (s PINState) String() string {
	switch s {
	case PINStateUnknown:
		return "Unknown"
	case PINStateNotSet:
		return "NotSet"
	case PINStateSet:
		return "Set"
	case PINStateRequested:
		return "Requested"
	case PINStateVerifying:
		return "Verifying"
	case PINStateVerified:
		return "Verified"
	case PINStateChanging:
		return "Changing"
	case PINStateLocked:
		return "Locked"
	default:
		return "Invalid"
	}
}

// PINManager manages PIN code-related functionality
type PINManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	state       PINState
	mutex       sync.RWMutex
	pin         string
	seed        uint32
	accountID   uint32
	flag        int
	lock        int
	attempts    int
	maxAttempts int
	lastAttempt time.Time
}

// NewPINManager creates a new PIN manager
func NewPINManager(parser *core.CoreParser, hookManager *hooks.HookManager) *PINManager {
	return &PINManager{
		parser:      parser,
		hookManager: hookManager,
		state:       PINStateUnknown,
		maxAttempts: 3,
		lastAttempt: time.Now(),
	}
}

// RegisterHandlers registers PIN-related packet handlers
func (m *PINManager) RegisterHandlers() {
	// Register handlers for PIN-related packets
	m.parser.RegisterHandlerFunc("02AD", "login_pin_code_request", "v V",
		[]string{"flag", "key"},
		m.handleLoginPinCodeRequest)

	m.parser.RegisterHandlerFunc("08B9", "login_pin_code_request", "V a4 v",
		[]string{"seed", "accountID", "flag"},
		m.handleLoginPinCodeRequest)

	m.parser.RegisterHandlerFunc("08BB", "login_pin_new_code_result", "v V",
		[]string{"flag", "seed"},
		m.handleLoginPinNewCodeResult)

	m.parser.RegisterHandlerFunc("0AE9", "login_pin_code_request", "V a4 v2",
		[]string{"seed", "accountID", "flag", "lock"},
		m.handleLoginPinCodeRequest)
}

// SetPIN sets the PIN code
func (m *PINManager) SetPIN(pin string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Validate PIN
	if err := m.validatePIN(pin); err != nil {
		return err
	}

	m.pin = pin
	m.state = PINStateSet

	return nil
}

// GetPIN returns the PIN code
func (m *PINManager) GetPIN() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.pin
}

// GetState returns the current PIN state
func (m *PINManager) GetState() PINState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state
}

// SetState sets the PIN state
func (m *PINManager) SetState(state PINState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.state = state

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_state_change", map[string]interface{}{
			"state": state,
		})
	}
}

// GetSeed returns the PIN seed
func (m *PINManager) GetSeed() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.seed
}

// GetAccountID returns the account ID
func (m *PINManager) GetAccountID() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.accountID
}

// GetFlag returns the PIN flag
func (m *PINManager) GetFlag() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.flag
}

// GetLock returns the PIN lock status
func (m *PINManager) GetLock() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.lock
}

// IsLocked returns whether the PIN is locked
func (m *PINManager) IsLocked() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state == PINStateLocked || m.lock > 0
}

// IsVerified returns whether the PIN has been verified
func (m *PINManager) IsVerified() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state == PINStateVerified
}

// IsRequired returns whether a PIN is required
func (m *PINManager) IsRequired() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state == PINStateRequested || m.flag == 1 || m.flag == 3 || m.flag == 5
}

// VerifyPIN verifies the PIN code
func (m *PINManager) VerifyPIN(pin string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if PIN is locked
	if m.state == PINStateLocked {
		return ErrPINLocked
	}

	// Check if PIN is required
	if m.state != PINStateRequested && m.state != PINStateVerifying {
		return ErrPINNotSet
	}

	// Validate PIN
	if err := m.validatePIN(pin); err != nil {
		return err
	}

	// Check if PIN matches
	if m.pin != "" && m.pin != pin {
		m.attempts++
		m.lastAttempt = time.Now()

		// Check if max attempts reached
		if m.attempts >= m.maxAttempts {
			m.state = PINStateLocked

			// Call hook
			if m.hookManager != nil {
				m.hookManager.CallHook("security/pin_locked", map[string]interface{}{
					"attempts": m.attempts,
				})
			}

			return ErrPINLocked
		}

		return ErrInvalidPIN
	}

	// PIN verified
	m.attempts = 0
	m.state = PINStateVerified
	m.lastAttempt = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_verified", nil)
	}

	return nil
}

// ChangePIN changes the PIN code
func (m *PINManager) ChangePIN(oldPIN, newPIN string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if PIN is locked
	if m.state == PINStateLocked {
		return ErrPINLocked
	}

	// Check if old PIN matches
	if m.pin != "" && m.pin != oldPIN {
		m.attempts++
		m.lastAttempt = time.Now()

		// Check if max attempts reached
		if m.attempts >= m.maxAttempts {
			m.state = PINStateLocked

			// Call hook
			if m.hookManager != nil {
				m.hookManager.CallHook("security/pin_locked", map[string]interface{}{
					"attempts": m.attempts,
				})
			}

			return ErrPINLocked
		}

		return ErrInvalidPIN
	}

	// Validate new PIN
	if err := m.validatePIN(newPIN); err != nil {
		return err
	}

	// Change PIN
	m.pin = newPIN
	m.attempts = 0
	m.state = PINStateSet
	m.lastAttempt = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_changed", nil)
	}

	return nil
}

// ResetPIN resets the PIN code
func (m *PINManager) ResetPIN() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.pin = ""
	m.attempts = 0
	m.state = PINStateNotSet
	m.lastAttempt = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_reset", nil)
	}
}

// UnlockPIN unlocks the PIN
func (m *PINManager) UnlockPIN() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state == PINStateLocked {
		m.state = PINStateSet
		m.attempts = 0
		m.lastAttempt = time.Now()

		// Call hook
		if m.hookManager != nil {
			m.hookManager.CallHook("security/pin_unlocked", nil)
		}
	}
}

// Packet handlers

// handleLoginPinCodeRequest handles the login_pin_code_request packet
func (m *PINManager) handleLoginPinCodeRequest(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract seed, account ID, flag, and lock
	if seedVal, ok := args["seed"].(uint32); ok {
		m.seed = seedVal
	}

	if accountIDVal, ok := args["accountID"].([]byte); ok && len(accountIDVal) >= 4 {
		m.accountID = uint32(accountIDVal[0]) | uint32(accountIDVal[1])<<8 | uint32(accountIDVal[2])<<16 | uint32(accountIDVal[3])<<24
	}

	if flagVal, ok := args["flag"].(uint16); ok {
		m.flag = int(flagVal)
	} else if flagVal, ok := args["flag"].(byte); ok {
		m.flag = int(flagVal)
	}

	if lockVal, ok := args["lock"].(uint16); ok {
		m.lock = int(lockVal)
	}

	// Update state based on flag
	switch m.flag {
	case 0:
		// PIN is correct
		m.state = PINStateVerified
	case 1:
		// PIN is required
		m.state = PINStateRequested
	case 2:
		// PIN must be changed
		m.state = PINStateChanging
	case 3:
		// Create new PIN
		m.state = PINStateNotSet
	case 4:
		// PIN must be changed (expired)
		m.state = PINStateChanging
	case 5:
		// PIN is invalid
		m.state = PINStateRequested
		m.attempts++
	case 6:
		// PIN is correct (login state)
		m.state = PINStateVerified
	case 7:
		// PIN is correct (change state)
		m.state = PINStateVerified
	default:
		// Unknown flag
		m.state = PINStateUnknown
	}

	// Check if PIN is locked
	if m.lock > 0 {
		m.state = PINStateLocked
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_code_request", map[string]interface{}{
			"seed":      m.seed,
			"accountID": m.accountID,
			"flag":      m.flag,
			"lock":      m.lock,
			"state":     m.state,
		})
	}

	return nil
}

// handleLoginPinNewCodeResult handles the login_pin_new_code_result packet
func (m *PINManager) handleLoginPinNewCodeResult(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var flag int
	var seed uint32

	if flagVal, ok := args["flag"].(uint16); ok {
		flag = int(flagVal)
	}

	if seedVal, ok := args["seed"].(uint32); ok {
		seed = seedVal
	}

	// Update state based on flag
	switch flag {
	case 0:
		// PIN successfully changed
		m.state = PINStateSet
	case 1:
		// PIN change failed
		// State remains unchanged
	}

	// Update seed
	m.seed = seed

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_new_code_result", map[string]interface{}{
			"flag": flag,
			"seed": seed,
		})
	}

	return nil
}

// Helper functions

// validatePIN validates a PIN code
func (m *PINManager) validatePIN(pin string) error {
	// Check PIN length
	if len(pin) != 4 {
		return ErrPINWrongLength
	}

	// Check PIN format (must be digits only)
	for _, c := range pin {
		if c < '0' || c > '9' {
			return ErrPINInvalidFormat
		}
	}

	return nil
}

// EncryptPIN encrypts a PIN code using the seed
func (m *PINManager) EncryptPIN(pin string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Validate PIN
	if err := m.validatePIN(pin); err != nil {
		return "", err
	}

	// Simple XOR encryption with seed
	// In a real implementation, this would be more complex
	encrypted := ""
	seed := fmt.Sprintf("%08x", m.seed) // Convert seed to hex string

	for i, c := range pin {
		// XOR the PIN digit with a byte from the seed
		seedByte := seed[i%len(seed)]
		encryptedByte := byte(c) ^ byte(seedByte)
		encrypted += fmt.Sprintf("%02x", encryptedByte)
	}

	return encrypted, nil
}

// DecryptPIN decrypts an encrypted PIN code using the seed
func (m *PINManager) DecryptPIN(encryptedPIN string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check encrypted PIN length
	if len(encryptedPIN) != 8 { // 4 bytes, each represented as 2 hex chars
		return "", errors.New("invalid encrypted PIN length")
	}

	// Simple XOR decryption with seed
	// In a real implementation, this would be more complex
	decrypted := ""
	seed := fmt.Sprintf("%08x", m.seed) // Convert seed to hex string

	for i := 0; i < len(encryptedPIN); i += 2 {
		// Convert hex to byte
		var encryptedByte byte
		fmt.Sscanf(encryptedPIN[i:i+2], "%02x", &encryptedByte)

		// XOR the encrypted byte with a byte from the seed
		seedByte := seed[(i/2)%len(seed)]
		decryptedByte := encryptedByte ^ byte(seedByte)

		// Check if decrypted byte is a digit
		if decryptedByte < '0' || decryptedByte > '9' {
			return "", ErrPINInvalidFormat
		}

		decrypted += string(decryptedByte)
	}

	return decrypted, nil
}

// QueryLoginPinCode requests a PIN code from the user
// Returns the PIN code, or an empty string if cancelled
func (m *PINManager) QueryLoginPinCode(message string) (string, error) {
	// In the Go implementation, we'll use hooks to request the PIN from the UI
	if m.hookManager == nil {
		return "", errors.New("hook manager is not initialized")
	}

	// Create a channel to receive the PIN
	pinChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Call the hook to request the PIN
	m.hookManager.CallHook("security/request_pin", map[string]interface{}{
		"message": message,
		"callback": func(pin string, err error) {
			pinChan <- pin
			errChan <- err
		},
	})

	// Wait for the PIN
	pin := <-pinChan
	err := <-errChan

	if err != nil {
		return "", err
	}

	// Validate the PIN
	if err := m.validatePIN(pin); err != nil {
		return "", err
	}

	return pin, nil
}

// QueryAndSaveLoginPinCode requests a PIN code from the user and saves it
// Returns true if the PIN was successfully saved, false otherwise
func (m *PINManager) QueryAndSaveLoginPinCode(message string) (bool, error) {
	pin, err := m.QueryLoginPinCode(message)
	if err != nil {
		return false, err
	}

	// Save the PIN
	if err := m.SetPIN(pin); err != nil {
		return false, err
	}

	// Call hook to save PIN in configuration
	if m.hookManager != nil {
		m.hookManager.CallHook("security/save_pin", map[string]interface{}{
			"pin": pin,
		})
	}

	return true, nil
}
