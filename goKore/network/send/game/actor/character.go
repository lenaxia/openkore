// Package actor provides functionality for character-related packets
package actor

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// Manager handles character-related packets
type Manager struct {
	send        core.BaseSend
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewManager creates a new character manager
func NewManager(send core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) *Manager {
	return &Manager{
		send:        send,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all character-related packet handlers
func (m *Manager) RegisterHandlers() {
	// Register character handlers
	if m.send != nil {
		m.send.RegisterHandler("char_create", m.HandleCharCreate)
		m.send.RegisterHandler("char_delete", m.HandleCharDelete)
	}
}

// HandleCharCreate handles the char_create packet
// This is a placeholder for the already implemented function
func (m *Manager) HandleCharCreate(args map[string]interface{}) ([]byte, error) {
	// Log the event
	if m.logger != nil {
		m.logger.Debug("Handling char_create packet")
	}

	// Call hook if needed
	if m.hookManager != nil {
		m.hookManager.CallHook("send/actor/char_create", args)
	}

	// Use the send component to construct the packet
	return m.send.ConstructPacket("char_create", args)
}

// HandleCharDelete handles the char_delete packet
// Implementation based on the Perl function:
//
//	sub sendCharDelete {
//	  my ($self, $charID, $email) = @_;
//	  $self->sendToServer($self->reconstruct({
//	    switch => 'char_delete',
//	    charID => $charID,
//	    email => stringToBytes($email),
//	  }));
//	  debug "Sent Char Delete\n", "sendPacket", 2;
//	}
func (m *Manager) HandleCharDelete(args map[string]interface{}) ([]byte, error) {
	// Log the event
	if m.logger != nil {
		m.logger.Debug("Handling char_delete packet")
	}

	// Call hook if needed
	if m.hookManager != nil {
		m.hookManager.CallHook("send/actor/char_delete", args)
	}

	// Use the send component to construct the packet
	return m.send.ConstructPacket("char_delete", args)
}
