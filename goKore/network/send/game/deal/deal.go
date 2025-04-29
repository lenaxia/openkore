// Package deal provides deal-related packet sending functionality.
package deal

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// DealManager handles deal-related packet sending.
type DealManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewDealManager creates a new deal manager.
func NewDealManager(baseSend core.Send) *DealManager {
	return &DealManager{
		baseSend: baseSend,
	}
}

// SendDealAddItem sends a request to add an item to the current deal.
// This is equivalent to the sendDealAddItem function in Send.pm.
func (dm *DealManager) SendDealAddItem(itemID []byte, amount uint16) error {
	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_item_add")
	if !exists {
		return fmt.Errorf("deal_item_add packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     itemID,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}

// SendDeal sends a request to initiate a deal with another player.
// This is equivalent to the sendDeal function in Send.pm.
func (dm *DealManager) SendDeal(playerID uint32) error {
	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_initiate")
	if !exists {
		return fmt.Errorf("deal_initiate packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": playerID,
	}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}

// SendDealReply sends a reply to a deal request.
// This is equivalent to the sendDealReply function in Send.pm.
// action:
//
//	3 = Accept
//	4 = Cancel
func (dm *DealManager) SendDealReply(action uint8) error {
	// Validate action
	if action != 3 && action != 4 {
		return fmt.Errorf("invalid deal reply action: %d (must be 3 or 4)", action)
	}

	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_reply")
	if !exists {
		return fmt.Errorf("deal_reply packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"action": action,
	}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}

// SendDealFinalize sends a request to finalize the current deal.
// This is equivalent to the sendDealFinalize function in Send.pm.
func (dm *DealManager) SendDealFinalize() error {
	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_finalize")
	if !exists {
		return fmt.Errorf("deal_finalize packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}

// SendCurrentDealCancel sends a request to cancel the current deal.
// This is equivalent to the sendCurrentDealCancel function in Send.pm.
func (dm *DealManager) SendCurrentDealCancel() error {
	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_cancel")
	if !exists {
		return fmt.Errorf("deal_cancel packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}

// SendDealTrade sends a request to complete the current deal.
// This is equivalent to the sendDealTrade function in Send.pm.
func (dm *DealManager) SendDealTrade() error {
	// Get the packet ID
	packetID, exists := dm.baseSend.GetPacketID("deal_trade")
	if !exists {
		return fmt.Errorf("deal_trade packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := dm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return dm.baseSend.SendToServer(packet)
}
