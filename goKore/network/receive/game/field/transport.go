// Package field provides handlers for field-related packets.
package field

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// FieldManager manages field-related packet handlers
type FieldManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewFieldManager creates a new field manager
func NewFieldManager(parser *core.CoreParser, hookManager *hooks.HookManager) *FieldManager {
	return &FieldManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterTransportHandlers registers all handlers related to field transportation
func (m *FieldManager) RegisterTransportHandlers() {
	// Register private_airship_type handler
	if m.parser != nil {
		m.parser.RegisterHandlerFunc("0A4B", "private_airship_type", "B",
			[]string{"fail"},
			m.handlePrivateAirshipType)
	}
}

// RegisterAllHandlers registers all field-related handlers
func (m *FieldManager) RegisterAllHandlers() {
	// Register transportation handlers
	m.RegisterTransportHandlers()

	// Register boss-related handlers
	m.RegisterBossHandlers()

	// Register navigation-related handlers
	m.RegisterNavigationHandlers()
}

// handlePrivateAirshipType handles the private_airship_type packet
// Packet format: 0A4B <fail>.B
func (m *FieldManager) handlePrivateAirshipType(args map[string]interface{}) error {
	// Process the packet
	result := m.processPrivateAirshipType(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("field.private_airship_type", result)
	}

	return nil
}

// processPrivateAirshipType processes the private_airship_type packet and returns a structured result
func (m *FieldManager) processPrivateAirshipType(args map[string]interface{}) map[string]interface{} {
	var failCode byte
	var status string

	// Extract fail code from args
	if failVal, ok := args["fail"].(byte); ok {
		failCode = failVal
	}

	// Process based on fail code value
	switch failCode {
	case 0:
		status = "Use Private Airship success."
	case 1:
		status = "Please try PivateAirship again."
	case 2:
		status = "You do not have enough Item to use PivateAirship."
	case 3:
		status = "Destination map is invalid."
	case 4:
		status = "Source map is invalid."
	case 5:
		status = "Item unavailable for use PivateAirship."
	default:
		status = fmt.Sprintf("Unknown Private Airship result: %d", failCode)
	}

	// Return structured result
	return map[string]interface{}{
		"result": failCode,
		"status": status,
	}
}

// RegisterWithParser registers the field manager with the given parser and hook manager
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the field manager
	manager := NewFieldManager(parser, hookManager)

	// Register all handlers
	manager.RegisterAllHandlers()
}

// RegisterWithBaseReceive registers the field manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the private_airship_type handler
	baseReceive.RegisterHandler("private_airship_type", func(args map[string]interface{}) error {
		// Create a field manager for this specific call
		manager := NewFieldManager(nil, nil)
		return manager.handlePrivateAirshipType(args)
	})

	// Register the boss_map_info handler
	baseReceive.RegisterHandler("boss_map_info", func(args map[string]interface{}) error {
		// Create a field manager for this specific call
		manager := NewFieldManager(nil, nil)
		return manager.handleBossMapInfo(args)
	})

	// Register the navigate_to handler
	baseReceive.RegisterHandler("navigate_to", func(args map[string]interface{}) error {
		// Create a field manager for this specific call
		manager := NewFieldManager(nil, nil)
		return manager.handleNavigateTo(args)
	})

	// Register the warp_portal_list handler
	baseReceive.RegisterHandler("warp_portal_list", func(args map[string]interface{}) error {
		// Create a field manager for this specific call
		manager := NewFieldManager(nil, nil)
		return manager.handleWarpPortalList(args)
	})
}
