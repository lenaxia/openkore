// Package core provides core functionality for parsing and processing network packets.
package core

// MovementManager manages movement-related functionality
type MovementManager struct {
	parser         *CoreParser
	accountManager *AccountManager
}

// NewMovementManager creates a new movement manager
func NewMovementManager(parser *CoreParser, accountManager *AccountManager) *MovementManager {
	return &MovementManager{
		parser:         parser,
		accountManager: accountManager,
	}
}

// RegisterHandlers registers movement-related packet handlers
func (m *MovementManager) RegisterHandlers() {
	// Register handler for move_interrupt
	m.parser.RegisterHandlerFunc("0AB8", "move_interrupt", "",
		[]string{},
		func(args map[string]interface{}) error {
			return m.accountManager.handleMoveInterrupt(args)
		})
}
