// Package cashshop provides cash shop-related packet sending functionality.
package cashshop

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CashShopManager handles cash shop-related packet sending.
type CashShopManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCashShopManager creates a new cash shop manager.
func NewCashShopManager(baseSend core.Send) *CashShopManager {
	return &CashShopManager{
		baseSend: baseSend,
	}
}

// GetManagerName returns the name of the manager.
// This implements the ManagerProvider interface.
func (csm *CashShopManager) GetManagerName() string {
	return "CashShopManager"
}

// BuyBulk sends a buy bulk request to the cash shop.
// This is equivalent to the sendCashBuyBulk function in Send.pm.
func (csm *CashShopManager) BuyBulk(kafraPoints int, items []int) error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_buy_bulk")
	if !exists {
		return fmt.Errorf("cash_buy_bulk packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"kafra_points": kafraPoints,
		"items":        items,
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// OpenShop sends a request to open the cash shop.
// This is equivalent to the sendCashShopOpen function in Send.pm.
func (csm *CashShopManager) OpenShop() error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_open")
	if !exists {
		return fmt.Errorf("cash_shop_open packet ID not found")
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// CloseShop sends a request to close the cash shop.
// This is equivalent to the sendCashShopClose function in Send.pm.
func (csm *CashShopManager) CloseShop() error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_close")
	if !exists {
		return fmt.Errorf("cash_shop_close packet ID not found")
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// List sends a request to list cash shop items.
// This is equivalent to the sendCashList function in Send.pm.
func (csm *CashShopManager) List(tab int) error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_list")
	if !exists {
		return fmt.Errorf("cash_shop_list packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"tab": tab,
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// Buy sends a request to buy a cash shop item.
// This is equivalent to the sendCashBuy function in Send.pm.
func (csm *CashShopManager) Buy(itemID int, amount int, kafraPoints int) error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_buy")
	if !exists {
		return fmt.Errorf("cash_shop_buy packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"item_id":      itemID,
		"amount":       amount,
		"kafra_points": kafraPoints,
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// RequestPoints sends a request to get cash shop points.
// This is equivalent to the sendCashRequest function in Send.pm.
func (csm *CashShopManager) RequestPoints() error {
	// Try to get the packet ID using the first name
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_request_points")
	if !exists {
		// If the first name doesn't exist, try the second name
		packetID, exists = csm.baseSend.GetPacketID("request_cashitems")
		if !exists {
			return fmt.Errorf("cash_shop_request_points or request_cashitems packet ID not found")
		}
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}

// CheckCoupon sends a request to check a cash shop coupon.
// This is equivalent to the sendCashCheckCoupon function in Send.pm.
func (csm *CashShopManager) CheckCoupon(couponCode string) error {
	// Get the packet ID
	packetID, exists := csm.baseSend.GetPacketID("cash_shop_check_coupon")
	if !exists {
		return fmt.Errorf("cash_shop_check_coupon packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"coupon_code": couponCode,
	}

	// Construct and send the packet
	packet, err := csm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return csm.baseSend.SendToServer(packet)
}
