// Package actor provides actor-related packet sending functionality.
package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// SyncManager handles actor synchronization-related packet sending.
type SyncManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewSyncManager creates a new sync manager.
func NewSyncManager(baseSend core.Send) *SyncManager {
	return &SyncManager{
		baseSend: baseSend,
	}
}

// SendSync sends a sync packet to keep the connection alive.
// This is equivalent to the sendSync function in Send.pm.
func (sm *SyncManager) SendSync(initialSync bool) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("sync")
	if !exists {
		return fmt.Errorf("sync packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"time": sm.baseSend.GetTime(),
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendCharacterMove sends a character move packet.
// This is equivalent to the sendMove function in Send.pm.
func (sm *SyncManager) SendCharacterMove(x, y int) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("character_move")
	if !exists {
		return fmt.Errorf("character_move packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"x":    uint16(x),
		"y":    uint16(y),
		"time": sm.baseSend.GetTime(),
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}
