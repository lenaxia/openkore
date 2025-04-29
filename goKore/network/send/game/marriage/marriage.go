// Package marriage provides marriage-related packet sending functionality.
package marriage

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MarriageManager handles marriage-related packet sending.
type MarriageManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMarriageManager creates a new marriage manager.
func NewMarriageManager(baseSend core.Send) *MarriageManager {
	return &MarriageManager{
		baseSend: baseSend,
	}
}

// SendAdoptRequest sends a request to adopt a player.
// This is equivalent to the sendAdoptRequest function in Send.pm.
func (mm *MarriageManager) SendAdoptRequest(playerID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("adopt_request")
	if !exists {
		return fmt.Errorf("adopt_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": playerID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendAdoptReply sends a reply to an adoption request.
// This is equivalent to the sendAdoptReply function in Send.pm.
// result: 1 = accept, 0 = reject
func (mm *MarriageManager) SendAdoptReply(parentID1, parentID2 uint32, result uint8) error {
	// Validate result
	if result > 1 {
		return fmt.Errorf("invalid adoption reply result: %d (must be 0 or 1)", result)
	}

	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("adopt_reply_request")
	if !exists {
		return fmt.Errorf("adopt_reply_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"parentID1": parentID1,
		"parentID2": parentID2,
		"result":    result,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
