// Package movement provides movement-related packet sending functionality.
package movement

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// TeleportManager handles teleportation-related packet sending.
type TeleportManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewTeleportManager creates a new teleport manager.
func NewTeleportManager(baseSend core.Send) *TeleportManager {
	return &TeleportManager{
		baseSend: baseSend,
	}
}

// SendWarpTele sends a warp teleport request.
// This is equivalent to the sendWarpTele function in Send.pm.
// skillID:
// 26 => Teleport (Respawn/Random)
// 27 => Open Warp
func (tm *TeleportManager) SendWarpTele(skillID int, mapName string) error {
	// Get the packet ID
	packetID, exists := tm.baseSend.GetPacketID("warp_select")
	if !exists {
		return fmt.Errorf("warp_select packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID": skillID,
		"mapName": mapName,
	}

	// Construct and send the packet
	packet, err := tm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return tm.baseSend.SendToServer(packet)
}

// SendPrivateAirshipRequest sends a private airship request.
// This is equivalent to the sendPrivateAirshipRequest function in Send.pm.
func (tm *TeleportManager) SendPrivateAirshipRequest(mapName string, nameID int) error {
	// Get the packet ID
	packetID, exists := tm.baseSend.GetPacketID("private_airship_request")
	if !exists {
		return fmt.Errorf("private_airship_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"map_name": mapName,
		"nameID":   nameID,
	}

	// Construct and send the packet
	packet, err := tm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return tm.baseSend.SendToServer(packet)
}
