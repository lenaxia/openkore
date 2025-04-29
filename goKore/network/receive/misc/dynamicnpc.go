// Package misc provides handlers for miscellaneous packets that don't fit into other categories.
package misc

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Constants for dynamicnpc_create_result
const (
	DYNAMICNPC_RESULT_SUCCESS    = 0x0
	DYNAMICNPC_RESULT_UNKNOWN    = 0x1
	DYNAMICNPC_RESULT_UNKNOWNNPC = 0x2
	DYNAMICNPC_RESULT_DUPLICATE  = 0x3
	DYNAMICNPC_RESULT_OUTOFTIME  = 0x4
)

// MiscManager manages miscellaneous packet handlers
type MiscManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewMiscManager creates a new miscellaneous packet manager
func NewMiscManager(parser *core.CoreParser, hookManager *hooks.HookManager) *MiscManager {
	return &MiscManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to miscellaneous functionality
func (m *MiscManager) RegisterHandlers() {
	// Register dynamicnpc_create_result handler
	m.parser.RegisterHandlerFunc("0A17", "dynamicnpc_create_result", "B",
		[]string{"result"},
		m.handleDynamicNPCCreateResult)
}

// handleDynamicNPCCreateResult handles the dynamicnpc_create_result packet
// Packet format: 0A17 <result>.B
func (m *MiscManager) handleDynamicNPCCreateResult(args map[string]interface{}) error {
	// Process the packet
	result := m.processDynamicNPCCreateResult(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("misc.dynamicnpc_create_result", result)
	}

	return nil
}

// processDynamicNPCCreateResult processes the dynamicnpc_create_result packet and returns a structured result
func (m *MiscManager) processDynamicNPCCreateResult(args map[string]interface{}) map[string]interface{} {
	var resultCode uint8
	var status string

	// Extract result from args
	if resultVal, ok := args["result"].(uint8); ok {
		resultCode = resultVal
	}

	// Process based on result value
	switch resultCode {
	case DYNAMICNPC_RESULT_SUCCESS:
		status = "Success"
	case DYNAMICNPC_RESULT_UNKNOWN:
		status = "Unknown"
	case DYNAMICNPC_RESULT_UNKNOWNNPC:
		status = "Unknown NPC"
	case DYNAMICNPC_RESULT_DUPLICATE:
		status = "Duplicate"
	case DYNAMICNPC_RESULT_OUTOFTIME:
		status = "Out of time"
	default:
		status = fmt.Sprintf("Unknown Result: %d", resultCode)
	}

	// Log the message
	// In the original implementation, this would call message()
	// We'll use the hook system to handle this

	// Return structured result
	return map[string]interface{}{
		"result":  resultCode,
		"status":  status,
		"message": fmt.Sprintf("Dynamic NPC create result - Status: %s", status),
	}
}

// RegisterWithParser registers the misc manager with the given parser and hook manager
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the misc manager
	manager := NewMiscManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()
}

// RegisterWithBaseReceive registers the misc manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the dynamicnpc_create_result handler
	baseReceive.RegisterHandler("dynamicnpc_create_result", func(args map[string]interface{}) error {
		// Create a misc manager for this specific call
		manager := NewMiscManager(nil, nil)
		return manager.handleDynamicNPCCreateResult(args)
	})
}
