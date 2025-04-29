// Package npc provides NPC-related packet sending functionality.
package npc

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// NPCManager handles NPC-related packet sending.
type NPCManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewNPCManager creates a new NPC manager.
func NewNPCManager(baseSend core.Send) *NPCManager {
	return &NPCManager{
		baseSend: baseSend,
	}
}

// SendNPCBuySellList sends a request to get the buy/sell list from an NPC.
// This is equivalent to the sendNPCBuySellList function in Send.pm.
func (nm *NPCManager) SendNPCBuySellList(ID uint32, type_ int) error {
	// Get the packet ID
	packetID, exists := nm.baseSend.GetPacketID("request_buy_sell_list")
	if !exists {
		return fmt.Errorf("request_buy_sell_list packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   ID,
		"type": type_,
	}

	// Construct and send the packet
	packet, err := nm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return nm.baseSend.SendToServer(packet)
}

// SendNPCCreateRequest sends a request to create an NPC.
// This is equivalent to the sendNPCCreateRequest function in Send.pm.
func (nm *NPCManager) SendNPCCreateRequest(name string) error {
	// Get the packet ID
	packetID, exists := nm.baseSend.GetPacketID("dynamicnpc_create_request")
	if !exists {
		return fmt.Errorf("dynamicnpc_create_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": name,
	}

	// Construct and send the packet
	packet, err := nm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return nm.baseSend.SendToServer(packet)
}
