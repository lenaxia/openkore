// Package chat provides chat-related packet sending functionality.
package chat

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// PublicChatManager handles public chat-related packet sending.
type PublicChatManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewPublicChatManager creates a new public chat manager.
func NewPublicChatManager(baseSend core.Send) *PublicChatManager {
	return &PublicChatManager{
		baseSend: baseSend,
	}
}

// SendChat sends a public chat message.
// Corresponds to sendChat in the original implementation.
func (pcm *PublicChatManager) SendChat(message string) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("public_chat")
	if !exists {
		return fmt.Errorf("public_chat packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// ParsePublicChat parses a public chat packet.
// Corresponds to parse_public_chat in the original implementation.
func (pcm *PublicChatManager) ParsePublicChat(args map[string]interface{}) error {
	// This would typically be handled by the receive component
	// but we include it here for completeness
	return nil
}

// ReconstructPublicChat reconstructs a public chat packet.
// Corresponds to reconstruct_public_chat in the original implementation.
func (pcm *PublicChatManager) ReconstructPublicChat(args map[string]interface{}) ([]byte, error) {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("public_chat")
	if !exists {
		return nil, fmt.Errorf("public_chat packet ID not found")
	}

	// Construct the packet
	return pcm.baseSend.Reconstruct(packetID, args)
}

// SendGMBroadcast sends a GM broadcast message.
// This corresponds to sendGMBroadcast in the original implementation.
func (pcm *PublicChatManager) SendGMBroadcast(message string) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("gm_broadcast")
	if !exists {
		return fmt.Errorf("gm_broadcast packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// SendLocalBroadcast sends a local broadcast message.
// This corresponds to sendGMBroadcastLocal in the original implementation.
func (pcm *PublicChatManager) SendLocalBroadcast(message string) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("local_broadcast")
	if !exists {
		return fmt.Errorf("local_broadcast packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}
