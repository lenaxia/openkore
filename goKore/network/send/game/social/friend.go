// Package social provides social-related packet sending functionality.
package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// FriendManager handles friend-related packet sending.
type FriendManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewFriendManager creates a new friend manager.
func NewFriendManager(baseSend core.Send) *FriendManager {
	return &FriendManager{
		baseSend: baseSend,
	}
}

// RequestFriendList sends a request to get the friend list.
func (fm *FriendManager) RequestFriendList() error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("request_friend_list")
	if !exists {
		return fmt.Errorf("request_friend_list packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// AddFriend sends a request to add a friend.
func (fm *FriendManager) AddFriend(name string) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("add_friend")
	if !exists {
		return fmt.Errorf("add_friend packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": name,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// RemoveFriend sends a request to remove a friend.
func (fm *FriendManager) RemoveFriend(accountID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("remove_friend")
	if !exists {
		return fmt.Errorf("remove_friend packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// AcceptFriendRequest sends a request to accept a friend request.
func (fm *FriendManager) AcceptFriendRequest(accountID, charID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("accept_friend_request")
	if !exists {
		return fmt.Errorf("accept_friend_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
		"char_id":    charID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// RejectFriendRequest sends a request to reject a friend request.
func (fm *FriendManager) RejectFriendRequest(accountID, charID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("reject_friend_request")
	if !exists {
		return fmt.Errorf("reject_friend_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
		"char_id":    charID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// BlockFriend sends a request to block a friend.
func (fm *FriendManager) BlockFriend(accountID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("block_friend")
	if !exists {
		return fmt.Errorf("block_friend packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// UnblockFriend sends a request to unblock a friend.
func (fm *FriendManager) UnblockFriend(accountID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("unblock_friend")
	if !exists {
		return fmt.Errorf("unblock_friend packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// RequestBlockList sends a request to get the block list.
func (fm *FriendManager) RequestBlockList() error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("request_block_list")
	if !exists {
		return fmt.Errorf("request_block_list packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}
