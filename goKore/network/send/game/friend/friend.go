// Package friend provides friend-related packet sending functionality.
package friend

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

// SendFriendListReply sends a reply to a friend request.
// This is equivalent to the sendFriendListReply function in Send.pm.
func (fm *FriendManager) SendFriendListReply(accountID, charID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("friend_response")
	if !exists {
		return fmt.Errorf("friend_response packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"friendAccountID": accountID,
		"friendCharID":    charID,
		"type":            flag,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// SendFriendRequest sends a friend request.
// This is equivalent to the sendFriendRequest function in Send.pm.
func (fm *FriendManager) SendFriendRequest(name string) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("friend_request")
	if !exists {
		return fmt.Errorf("friend_request packet ID not found")
	}

	// Convert name to bytes and pad to 24 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 24 {
		nameBytes = nameBytes[:24]
	}
	paddedName := make([]byte, 24)
	copy(paddedName, nameBytes)

	// Create the arguments
	args := map[string]interface{}{
		"username": paddedName,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// SendFriendRemove sends a request to remove a friend.
// This is equivalent to the sendFriendRemove function in Send.pm.
func (fm *FriendManager) SendFriendRemove(accountID, charID uint32) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("friend_remove")
	if !exists {
		return fmt.Errorf("friend_remove packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"charID":    charID,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// SendIgnore sends a request to ignore a player.
// This is equivalent to the sendIgnore function in Send.pm.
func (fm *FriendManager) SendIgnore(name string, flag int) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("ignore_player")
	if !exists {
		return fmt.Errorf("ignore_player packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": []byte(name), // stringToBytes in the original
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// SendIgnoreAll sends a request to ignore all players.
// This is equivalent to the sendIgnoreAll function in Send.pm.
func (fm *FriendManager) SendIgnoreAll(flag int) error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("ignore_all")
	if !exists {
		return fmt.Errorf("ignore_all packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := fm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return fm.baseSend.SendToServer(packet)
}

// SendGetIgnoreList sends a request to get the ignore list.
// This is equivalent to the sendGetIgnoreList function in Send.pm.
func (fm *FriendManager) SendGetIgnoreList() error {
	// Get the packet ID
	packetID, exists := fm.baseSend.GetPacketID("get_ignore_list")
	if !exists {
		return fmt.Errorf("get_ignore_list packet ID not found")
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
