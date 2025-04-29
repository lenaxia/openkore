package core

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
)

// UsersOnlineManager manages users online-related functionality
type UsersOnlineManager struct {
	hookManager *hooks.HookManager
}

// NewUsersOnlineManager creates a new users online manager
func NewUsersOnlineManager(hookManager *hooks.HookManager) *UsersOnlineManager {
	return &UsersOnlineManager{
		hookManager: hookManager,
	}
}

// HandleUsersOnline handles the users_online packet
// Packet format: 0AAC <users>.L
func (m *UsersOnlineManager) HandleUsersOnline(args map[string]interface{}) error {
	// Extract users count with safety check
	var users uint32
	if usersVal, ok := args["users"].(uint32); ok {
		users = usersVal
	}

	// Format the message
	message := fmt.Sprintf("There are currently %d users online", users)

	// Log the message
	// In a real implementation, this would use a proper logger
	// logger.Info(message)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("core.users_online", map[string]interface{}{
			"users":   users,
			"message": message,
		})
	}

	return nil
}

// RegisterHandlers registers users online-related packet handlers with the given parser
func (m *UsersOnlineManager) RegisterHandlers(parser *CoreParser) {
	parser.RegisterHandlerFunc("0AAC", "users_online", "L",
		[]string{"users"},
		m.HandleUsersOnline)
}
