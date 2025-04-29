package behavior

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RegisterWithParser registers the behavior manager with the given parser and hook manager
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the behavior manager
	manager := NewCharacterBehaviorManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()
}

// RegisterWithBaseReceive registers the behavior manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the manner_message handler
	baseReceive.RegisterHandler("manner_message", func(args map[string]interface{}) error {
		// Create a behavior manager for this specific call
		manager := NewCharacterBehaviorManager(nil, nil)
		return manager.handleMannerMessage(args)
	})

	// Register the hack_shield_alarm handler
	baseReceive.RegisterHandler("hack_shield_alarm", func(args map[string]interface{}) error {
		// Create a behavior manager for this specific call
		manager := NewCharacterBehaviorManager(nil, nil)
		return manager.handleHackShieldAlarm(args)
	})
}
