// Package chat provides chat-related packet sending functionality.
package chat

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// ChannelChatManager handles channel chat-related packet sending.
type ChannelChatManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewChannelChatManager creates a new channel chat manager.
func NewChannelChatManager(baseSend core.Send) *ChannelChatManager {
	return &ChannelChatManager{
		baseSend: baseSend,
	}
}

// JoinChannel sends a request to join a chat channel.
func (ccm *ChannelChatManager) JoinChannel(channelName, password string) error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("join_channel")
	if !exists {
		return fmt.Errorf("join_channel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"channel_name": channelName,
		"password":     password,
	}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}

// LeaveChannel sends a request to leave a chat channel.
func (ccm *ChannelChatManager) LeaveChannel(channelID uint32) error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("leave_channel")
	if !exists {
		return fmt.Errorf("leave_channel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"channel_id": channelID,
	}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}

// SendChannelMessage sends a message to a chat channel.
func (ccm *ChannelChatManager) SendChannelMessage(channelID uint32, message string) error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("channel_message")
	if !exists {
		return fmt.Errorf("channel_message packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"channel_id": channelID,
		"message":    message,
	}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}

// ListChannels sends a request to list available chat channels.
func (ccm *ChannelChatManager) ListChannels() error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("list_channels")
	if !exists {
		return fmt.Errorf("list_channels packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}

// GetChannelUserList sends a request to get the list of users in a channel.
func (ccm *ChannelChatManager) GetChannelUserList(channelID uint32) error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("channel_user_list")
	if !exists {
		return fmt.Errorf("channel_user_list packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"channel_id": channelID,
	}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}

// SetChannelOption sends a request to set channel options.
func (ccm *ChannelChatManager) SetChannelOption(channelID uint32, option int) error {
	// Get the packet ID
	packetID, exists := ccm.baseSend.GetPacketID("channel_option")
	if !exists {
		return fmt.Errorf("channel_option packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"channel_id": channelID,
		"option":     option,
	}

	// Construct and send the packet
	packet, err := ccm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return ccm.baseSend.SendToServer(packet)
}
