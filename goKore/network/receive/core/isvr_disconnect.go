package core

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// ISVRDisconnectManager manages ISVR disconnect-related functionality
type ISVRDisconnectManager struct {
	hookManager *hooks.HookManager
}

// NewISVRDisconnectManager creates a new ISVR disconnect manager
func NewISVRDisconnectManager(hookManager *hooks.HookManager) *ISVRDisconnectManager {
	return &ISVRDisconnectManager{
		hookManager: hookManager,
	}
}

// HandleISVRDisconnect handles the isvr_disconnect packet
// Packet format: 09CD
func (m *ISVRDisconnectManager) HandleISVRDisconnect(args map[string]interface{}) error {
	// Log the message
	// In a real implementation, this would use a proper logger
	// logger.Debug("Received the package 'isvr_disconnect'")

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("core.isvr_disconnect", map[string]interface{}{
			"message": "Received the package 'isvr_disconnect'",
		})
	}

	return nil
}

// RegisterHandlers registers ISVR disconnect-related packet handlers with the given parser
func (m *ISVRDisconnectManager) RegisterHandlers(parser *CoreParser) {
	parser.RegisterHandlerFunc("09CD", "isvr_disconnect", "",
		[]string{},
		m.HandleISVRDisconnect)
}
