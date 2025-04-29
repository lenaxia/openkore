package core

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// RegisterAllHandlers registers all handlers in the core package with the given parser
func RegisterAllHandlers(parser *CoreParser, hookManager *hooks.HookManager) {
	// Create managers
	behaviorManager := NewBehaviorManager(parser)
	hotkeysManager := NewHotkeysManager(parser, hookManager)
	remainTimeManager := NewRemainTimeManager(hookManager)
	usersOnlineManager := NewUsersOnlineManager(hookManager)
	isvrDisconnectManager := NewISVRDisconnectManager(hookManager)

	// Create errors manager with default config
	config := map[string]interface{}{
		"dcOnDisconnect":     1,
		"dcOnServerShutDown": 0,
		"dcOnServerClose":    0,
		"dcOnDualLogin":      0,
	}
	errorsManager := NewErrorsManager(hookManager, NetworkStateDisconnected, config, "goKore")

	// Register handlers
	behaviorManager.RegisterHandlers()
	hotkeysManager.RegisterHandlers()
	remainTimeManager.RegisterHandlers(parser)
	usersOnlineManager.RegisterHandlers(parser)
	isvrDisconnectManager.RegisterHandlers(parser)
	errorsManager.RegisterHandlers(parser)
}

// RegisterPacketDefinitions registers all packet definitions in the core package
func RegisterPacketDefinitions(parser *CoreParser) {
	// Register packet definitions for remain time info
	parser.RegisterHandler("0C42", "remain_time_info", "W L L",
		[]string{"result", "expiration_date", "remain_time"},
		nil)

	// Register packet definitions for hotkeys
	parser.RegisterHandler("07D9", "hotkeys", "B a*",
		[]string{"rotate", "hotkeys"},
		nil)

	// Register packet definitions for users online
	parser.RegisterHandler("0AAC", "users_online", "L",
		[]string{"users"},
		nil)

	// Register packet definitions for ISVR disconnect
	parser.RegisterHandler("09CD", "isvr_disconnect", "",
		[]string{},
		nil)

	// Register packet definitions for errors
	parser.RegisterHandler("0081", "errors", "B",
		[]string{"type"},
		nil)

	// Other packet definitions are registered in their respective managers
}
