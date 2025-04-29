// Package market provides market-related packet sending functionality.
package market

import (
	"encoding/binary"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MarketManager handles market-related packet sending.
type MarketManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMarketManager creates a new market manager.
func NewMarketManager(baseSend core.Send) *MarketManager {
	return &MarketManager{
		baseSend: baseSend,
	}
}

// ReconstructBuyBulkMarket reconstructs the buy bulk market packet.
// This is equivalent to the reconstruct_buy_bulk_market function in Send.pm.
func ReconstructBuyBulkMarket(args map[string]interface{}) {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return
	}

	buyInfo := []byte{}
	for _, item := range items {
		itemID, ok1 := item["itemID"].(uint16)
		amount, ok2 := item["amount"].(uint32)
		if !ok1 || !ok2 {
			continue
		}

		itemIDBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(itemIDBytes, itemID)

		amountBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(amountBytes, amount)

		buyInfo = append(buyInfo, itemIDBytes...)
		buyInfo = append(buyInfo, amountBytes...)
	}

	args["buyInfo"] = buyInfo
}

// SendBuyBulkMarket sends a request to buy items from the market.
// This is equivalent to the sendBuyBulkMarket function in Send.pm.
func (mm *MarketManager) SendBuyBulkMarket(items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("buy_bulk_market")
	if !exists {
		return fmt.Errorf("buy_bulk_market packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the buy info
	ReconstructBuyBulkMarket(args)

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMarketClose sends a request to close the market.
// This is equivalent to the sendMarketClose function in Send.pm.
func (mm *MarketManager) SendMarketClose() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("market_close")
	if !exists {
		return fmt.Errorf("market_close packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendSellBuyComplete sends a confirmation that the sell/buy process is complete.
// This is equivalent to the sendSellBuyComplete function in Send.pm.
func (mm *MarketManager) SendSellBuyComplete() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("sell_buy_complete")
	if !exists {
		return fmt.Errorf("sell_buy_complete packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
