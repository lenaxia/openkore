package actor

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// StylistManager manages stylist-related functionality
type StylistManager struct {
	hookManager *hooks.HookManager
}

// NewStylistManager creates a new stylist manager
func NewStylistManager(hookManager *hooks.HookManager) *StylistManager {
	return &StylistManager{
		hookManager: hookManager,
	}
}

// HandleStylistRes handles the stylist_res packet
// Packet format: 0A46 <result>.B
func (m *StylistManager) HandleStylistRes(args map[string]interface{}) error {
	// Extract result with safety check
	var result byte
	if resultVal, ok := args["result"].(byte); ok {
		result = resultVal
	}

	var message string
	if result != 0 {
		message = "[Stylist UI] Success."
	} else {
		message = "[Stylist UI] Fail."
	}

	// Log the message
	// In a real implementation, this would use a proper logger
	// if result != 0 {
	//     logger.Info(message)
	// } else {
	//     logger.Error(message)
	// }

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("game.actor.stylist_res", map[string]interface{}{
			"result":  result,
			"success": result != 0,
			"message": message,
		})
	}

	return nil
}

// RegisterHandlers registers stylist-related packet handlers with the given parser
func (m *StylistManager) RegisterHandlers(parser interface{}) {
	if p, ok := parser.(interface {
		RegisterHandlerFunc(id, name, format string, fieldNames []string, handler interface{})
	}); ok {
		p.RegisterHandlerFunc("0A46", "stylist_res", "B",
			[]string{"result"},
			m.HandleStylistRes)
	}
}
