// Package actor provides actor-related packet sending functionality.
package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MovementManager handles actor movement-related packet sending.
type MovementManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMovementManager creates a new movement manager.
func NewMovementManager(baseSend core.Send) *MovementManager {
	return &MovementManager{
		baseSend: baseSend,
	}
}

// SendMove sends a move command for the player character.
func (mm *MovementManager) SendMove(x, y int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("move_to")
	if !exists {
		return fmt.Errorf("move_to packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"x": x,
		"y": y,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendSlaveMove sends a move command for a slave (homunculus, mercenary, etc.).
func (mm *MovementManager) SendSlaveMove(slaveID uint32, x, y int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("slave_move_to")
	if !exists {
		return fmt.Errorf("slave_move_to packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"slave_id": slaveID,
		"x":        x,
		"y":        y,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendActorMove sends a move command for any actor (used for pet, etc.).
func (mm *MovementManager) SendActorMove(actorID uint32, x, y int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("actor_move_to")
	if !exists {
		return fmt.Errorf("actor_move_to packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"actor_id": actorID,
		"x":        x,
		"y":        y,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
