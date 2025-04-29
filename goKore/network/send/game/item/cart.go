// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CartManager handles cart-related packet sending.
type CartManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCartManager creates a new cart manager.
func NewCartManager(baseSend core.Send) *CartManager {
	return &CartManager{
		baseSend: baseSend,
	}
}

// SendCartAdd sends a request to add an item to the cart.
// This is equivalent to the sendCartAdd function in Send.pm.
func (cm *CartManager) SendCartAdd(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("cart_add")
	if !exists {
		return fmt.Errorf("cart_add packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCartGet sends a request to get an item from the cart.
// This is equivalent to the sendCartGet function in Send.pm.
func (cm *CartManager) SendCartGet(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("cart_get")
	if !exists {
		return fmt.Errorf("cart_get packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendChangeCart sends a request to change the cart level.
// This is equivalent to the sendChangeCart function in Send.pm.
func (cm *CartManager) SendChangeCart(level int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("change_cart")
	if !exists {
		return fmt.Errorf("change_cart packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"lvl": level,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}
