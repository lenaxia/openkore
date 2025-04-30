package skill

import (
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CastCancelledManager manages the cast_cancelled packet handler
type CastCancelledManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewCastCancelledManager creates a new cast cancelled manager
func NewCastCancelledManager(parser *core.CoreParser, hookManager *hooks.HookManager) *CastCancelledManager {
	return &CastCancelledManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to cast cancellation
func (m *CastCancelledManager) RegisterHandlers() {
	// Register cast_cancelled handler
	m.parser.RegisterHandlerFunc("01B9", "cast_cancelled", "a4",
		[]string{"ID"},
		m.handleCastCancelled)

	// Register cast_cancelled handler (alternative packet)
	m.parser.RegisterHandlerFunc("08CD", "cast_cancelled_expanded", "a4",
		[]string{"ID"},
		m.handleCastCancelled)
}

// handleCastCancelled handles the cast_cancelled packet
// Packet formats:
// 01B9: <ID>.L
// 08CD: <ID>.L
func (m *CastCancelledManager) handleCastCancelled(args map[string]interface{}) error {
	// Process the packet
	result := m.processCastCancelled(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.cast_cancelled", result)
	}

	return nil
}

// processCastCancelled processes the cast_cancelled packet and returns a structured result
func (m *CastCancelledManager) processCastCancelled(args map[string]interface{}) map[string]interface{} {
	// Extract actor ID from args
	var actorID uint32

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		actorID = val
	}

	// Return structured result
	return map[string]interface{}{
		"actorID":        actorID,
		"cancelledTime":  time.Now(),
		"isOwnCharacter": false, // This would be determined by comparing with accountID in a real implementation
	}
}
