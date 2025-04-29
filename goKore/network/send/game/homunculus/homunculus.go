// Package homunculus provides homunculus-related packet sending functionality.
package homunculus

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// HomunculusManager handles homunculus-related packet sending.
type HomunculusManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewHomunculusManager creates a new homunculus manager.
func NewHomunculusManager(baseSend core.Send) *HomunculusManager {
	return &HomunculusManager{
		baseSend: baseSend,
	}
}

// SendHomunculusName sends a request to rename a homunculus.
// This is equivalent to the sendHomunculusName function in Send.pm.
func (hm *HomunculusManager) SendHomunculusName(name string) error {
	// Get the packet ID
	packetID, exists := hm.baseSend.GetPacketID("homunculus_name")
	if !exists {
		return fmt.Errorf("homunculus_name packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := hm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return hm.baseSend.SendToServer(packet)
}

// SendHomunculusCommand sends a command to a homunculus.
// This is equivalent to the sendHomunculusCommand function in Send.pm.
// command can be 0:get stats, 1:feed or 2:fire
func (hm *HomunculusManager) SendHomunculusCommand(command, commandType uint8) error {
	// Get the packet ID
	packetID, exists := hm.baseSend.GetPacketID("homunculus_command")
	if !exists {
		return fmt.Errorf("homunculus_command packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"commandID":   command,
		"commandType": commandType,
	}

	// Construct and send the packet
	packet, err := hm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return hm.baseSend.SendToServer(packet)
}
