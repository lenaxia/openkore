// Package mercenary provides mercenary-related packet sending functionality.
package mercenary

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MercenaryManager handles mercenary-related packet sending.
type MercenaryManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMercenaryManager creates a new mercenary manager.
func NewMercenaryManager(baseSend core.Send) *MercenaryManager {
	return &MercenaryManager{
		baseSend: baseSend,
	}
}

// SendMercenaryCommand sends a command to control a mercenary.
// This is equivalent to the sendMercenaryCommand function in Send.pm.
// Command flags:
// 0 => COMMAND_REQ_NONE
// 1 => COMMAND_REQ_PROPERTY
// 2 => COMMAND_REQ_DELETE
func (mm *MercenaryManager) SendMercenaryCommand(command int) error {
	// Validate command
	if command < 0 || command > 2 {
		return fmt.Errorf("invalid mercenary command: %d", command)
	}

	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mercenary_command")
	if !exists {
		return fmt.Errorf("mercenary_command packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": uint8(command),
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendCompanionRelease sends a request to release a companion (Cart, Falcon or Pecopeco).
// This is equivalent to the sendCompanionRelease function in Send.pm.
func (mm *MercenaryManager) SendCompanionRelease() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("companion_release")
	if !exists {
		return fmt.Errorf("companion_release packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
