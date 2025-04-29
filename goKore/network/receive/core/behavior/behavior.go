// Package behavior provides handlers for character behavior-related packets.
package behavior

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CharacterBehaviorManager manages behavior-related packet handlers
type CharacterBehaviorManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewCharacterBehaviorManager creates a new character behavior manager
func NewCharacterBehaviorManager(parser *core.CoreParser, hookManager *hooks.HookManager) *CharacterBehaviorManager {
	return &CharacterBehaviorManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to character behavior
func (m *CharacterBehaviorManager) RegisterHandlers() {
	// Register manner_message handler
	m.parser.RegisterHandlerFunc("0149", "manner_message", "C",
		[]string{"flag"},
		m.handleMannerMessage)

	// Register hack_shield_alarm handler
	m.parser.RegisterHandlerFunc("08B3", "hack_shield_alarm", "",
		[]string{},
		m.handleHackShieldAlarm)
}

// handleMannerMessage handles the manner_message packet
// Packet format: 0149 <flag>.B
func (m *CharacterBehaviorManager) handleMannerMessage(args map[string]interface{}) error {
	// Process the packet
	result := m.processMannerMessage(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.manner_message", result)
	}

	return nil
}

// processMannerMessage processes the manner_message packet and returns a structured result
func (m *CharacterBehaviorManager) processMannerMessage(args map[string]interface{}) map[string]interface{} {
	var flag uint8
	var message string

	// Extract flag from args
	if flagVal, ok := args["flag"].(uint8); ok {
		flag = flagVal
	}

	// Process based on flag value
	switch flag {
	case 0:
		message = "A manner point has been successfully aligned."
	case 3:
		message = "Chat Block has been applied by GM due to your ill-mannerous action."
	case 4:
		message = "Automated Chat Block has been applied due to Anti-Spam System."
	case 5:
		message = "You got a good point."
	default:
		message = fmt.Sprintf("Unknown manner message result (flag: %d)", flag)
	}

	// Return structured result
	return map[string]interface{}{
		"flag":    flag,
		"message": message,
	}
}

// handleHackShieldAlarm handles the hack_shield_alarm packet
// Packet format: 08B3
func (m *CharacterBehaviorManager) handleHackShieldAlarm(args map[string]interface{}) error {
	// Process the packet
	result := m.processHackShieldAlarm()

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.hack_shield_alarm", result)
	}

	// In the original implementation, this would also run a command to relog
	// Commands::run('relog 100000000');
	// We'll need to implement this functionality elsewhere

	return nil
}

// processHackShieldAlarm processes the hack_shield_alarm packet and returns a structured result
func (m *CharacterBehaviorManager) processHackShieldAlarm() map[string]interface{} {
	message := "Error: You have been forced to disconnect by a Hack Shield. Please check Poseidon."

	// Return structured result
	return map[string]interface{}{
		"message": message,
	}
}
