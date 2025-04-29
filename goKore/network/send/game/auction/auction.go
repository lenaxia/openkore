// Package auction provides auction-related packet sending functionality.
package auction

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// AuctionManager handles auction-related packet sending.
type AuctionManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewAuctionManager creates a new auction manager.
func NewAuctionManager(baseSend core.Send) *AuctionManager {
	return &AuctionManager{
		baseSend: baseSend,
	}
}

// SendAuctionAddItemCancel sends a request to cancel adding an item to the auction.
// This is equivalent to the sendAuctionAddItemCancel function in Send.pm.
func (am *AuctionManager) SendAuctionAddItemCancel(flag uint8) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_add_item_cancel")
	if !exists {
		return fmt.Errorf("auction_add_item_cancel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionAddItem sends a request to add an item to the auction.
// This is equivalent to the sendAuctionAddItem function in Send.pm.
func (am *AuctionManager) SendAuctionAddItem(itemID uint32, amount uint16) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_add_item")
	if !exists {
		return fmt.Errorf("auction_add_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     itemID,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionCreate sends a request to create an auction.
// This is equivalent to the sendAuctionCreate function in Send.pm.
func (am *AuctionManager) SendAuctionCreate(nowPrice, maxPrice, deleteTime uint32) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_create")
	if !exists {
		return fmt.Errorf("auction_create packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"now_price":   nowPrice,
		"max_price":   maxPrice,
		"delete_time": deleteTime,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionCancel sends a request to cancel an auction.
// This is equivalent to the sendAuctionCancel function in Send.pm.
func (am *AuctionManager) SendAuctionCancel(auctionID uint32) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_cancel")
	if !exists {
		return fmt.Errorf("auction_cancel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": auctionID,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionBuy sends a request to buy an item from the auction.
// This is equivalent to the sendAuctionBuy function in Send.pm.
func (am *AuctionManager) SendAuctionBuy(auctionID, price uint32) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_buy")
	if !exists {
		return fmt.Errorf("auction_buy packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":    auctionID,
		"price": price,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionItemSearch sends a request to search for items in the auction.
// This is equivalent to the sendAuctionItemSearch function in Send.pm.
// searchType:
// 0 => armor
// 1 => weapon
// 2 => card
// 3 => misc
// 4 => name
// 5 => auction id
func (am *AuctionManager) SendAuctionItemSearch(searchType uint8, price uint32, searchString string, page uint16) error {
	// Validate search type
	if searchType > 5 {
		return fmt.Errorf("invalid search type: %d", searchType)
	}

	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_search")
	if !exists {
		return fmt.Errorf("auction_search packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":          searchType,
		"price":         price,
		"search_string": searchString,
		"page":          page,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionReqMyInfo sends a request to get information about the player's auctions.
// This is equivalent to the sendAuctionReqMyInfo function in Send.pm.
func (am *AuctionManager) SendAuctionReqMyInfo(infoType uint8) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_info_self")
	if !exists {
		return fmt.Errorf("auction_info_self packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": infoType,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAuctionMySellStop sends a request to stop selling an item in the auction.
// This is equivalent to the sendAuctionMySellStop function in Send.pm.
func (am *AuctionManager) SendAuctionMySellStop(auctionID uint32) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("auction_sell_stop")
	if !exists {
		return fmt.Errorf("auction_sell_stop packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": auctionID,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}
