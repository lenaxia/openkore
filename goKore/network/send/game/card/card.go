// Package card provides card-related packet sending functionality.
package card

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CardManager handles card-related packet sending.
type CardManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCardManager creates a new card manager.
func NewCardManager(baseSend core.Send) *CardManager {
	return &CardManager{
		baseSend: baseSend,
	}
}

// SendCardMergeRequest sends a request to merge a card.
// This is equivalent to the sendCardMergeRequest function in Send.pm.
func (cm *CardManager) SendCardMergeRequest(cardID uint32) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("card_merge_request")
	if !exists {
		return fmt.Errorf("card_merge_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"cardID": cardID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCardMerge sends a request to merge a card with an item.
// This is equivalent to the sendCardMerge function in Send.pm.
func (cm *CardManager) SendCardMerge(cardID, itemID uint32) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("card_merge")
	if !exists {
		return fmt.Errorf("card_merge packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"cardID": cardID,
		"itemID": itemID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}
