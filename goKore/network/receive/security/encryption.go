// Package security provides security-related functionality for the network stack.
package security

import (
	"errors"
	"sync"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Errors
var (
	ErrInvalidEncryptionParam = errors.New("invalid encryption parameter")
)

// EncryptionManager manages message ID encryption-related functionality
type EncryptionManager struct {
	parser              *core.CoreParser
	hookManager         *hooks.HookManager
	mutex               sync.RWMutex
	messageIDEncryption string
	encVal1             uint32
	encVal2             uint32
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(parser *core.CoreParser, hookManager *hooks.HookManager) *EncryptionManager {
	return &EncryptionManager{
		parser:              parser,
		hookManager:         hookManager,
		messageIDEncryption: "0", // Default to no encryption
	}
}

// RegisterHandlers registers encryption-related packet handlers
func (m *EncryptionManager) RegisterHandlers() {
	// Register handlers for encryption-related packets
	m.parser.RegisterHandlerFunc("02F1", "initialize_message_id_encryption", "V V",
		[]string{"param1", "param2"},
		m.HandleInitializeMessageIDEncryption)
}

// HandleInitializeMessageIDEncryption handles the initialize_message_id_encryption packet
func (m *EncryptionManager) HandleInitializeMessageIDEncryption(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If messageIDEncryption is "0", do nothing
	if m.messageIDEncryption == "0" {
		return nil
	}

	// Get parameters
	param1, ok := args["param1"].(uint32)
	if !ok {
		return ErrInvalidEncryptionParam
	}

	param2, ok := args["param2"].(uint32)
	if !ok {
		return ErrInvalidEncryptionParam
	}

	// Extract digits from param1
	c := make([]uint32, 9) // 1-indexed array, c[0] is unused
	shtmp := param1
	for i := 8; i > 0; i-- {
		c[i] = shtmp & 0x0F
		shtmp >>= 4
	}

	// Calculate w and enc_val1
	w := (c[6] << 12) + (c[4] << 8) + (c[7] << 4) + c[1]
	m.encVal1 = (c[2] << 12) + (c[3] << 8) + (c[5] << 4) + c[8]

	// Calculate enc_val2
	// Note: This calculation will overflow uint32, so we use uint64 for intermediate calculations
	val1 := uint64((m.encVal1^0x0000F3AC)+w) << 16
	val2 := uint64((m.encVal1 ^ 0x000049DF) + w)
	m.encVal2 = uint32((val1 | val2) ^ uint64(param2))

	// TODO: Send message ID encryption initialized packet
	// m.messageSender.sendMessageIDEncryptionInitialized()

	return nil
}

// SetMessageIDEncryption sets the message ID encryption setting
func (m *EncryptionManager) SetMessageIDEncryption(value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.messageIDEncryption = value
}

// GetMessageIDEncryption gets the message ID encryption setting
func (m *EncryptionManager) GetMessageIDEncryption() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.messageIDEncryption
}

// SetEncVal1 sets the enc_val1 value
func (m *EncryptionManager) SetEncVal1(value uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.encVal1 = value
}

// GetEncVal1 gets the enc_val1 value
func (m *EncryptionManager) GetEncVal1() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.encVal1
}

// SetEncVal2 sets the enc_val2 value
func (m *EncryptionManager) SetEncVal2(value uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.encVal2 = value
}

// GetEncVal2 gets the enc_val2 value
func (m *EncryptionManager) GetEncVal2() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.encVal2
}
