// Package chat provides chat-related packet sending functionality.
package chat

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// PrivateChatManager handles private chat-related packet sending.
type PrivateChatManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewPrivateChatManager creates a new private chat manager.
func NewPrivateChatManager(baseSend core.Send) *PrivateChatManager {
	return &PrivateChatManager{
		baseSend: baseSend,
	}
}

// ParsePrivateMessage parses a private message packet.
// Corresponds to parse_private_message in the original implementation.
func (pcm *PrivateChatManager) ParsePrivateMessage(args map[string]interface{}) error {
	// This would typically be handled by the receive component
	// but we include it here for completeness
	return nil
}

// ReconstructPrivateMessage reconstructs a private message packet.
// Corresponds to reconstruct_private_message in the original implementation.
func (pcm *PrivateChatManager) ReconstructPrivateMessage(args map[string]interface{}) ([]byte, error) {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("private_message")
	if !exists {
		return nil, fmt.Errorf("private_message packet ID not found")
	}

	// Construct the packet
	return pcm.baseSend.Reconstruct(packetID, args)
}

// SendPrivateMessage sends a private message to another player.
// This is an alias for SendPrivateMsg to match the test function name.
func (pcm *PrivateChatManager) SendPrivateMessage(target string, message string) error {
	return pcm.SendPrivateMsg(target, message)
}

// SendPrivateMsg sends a private message to another player.
// Corresponds to sendPrivateMsg in the original implementation.
func (pcm *PrivateChatManager) SendPrivateMsg(target string, message string) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("private_message")
	if !exists {
		return fmt.Errorf("private_message packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"target":  target,
		"message": message,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// SendWhisperResponse sends a response to a whisper.
func (pcm *PrivateChatManager) SendWhisperResponse(target string, response int) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("whisper_response")
	if !exists {
		return fmt.Errorf("whisper_response packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"target":   target,
		"response": response,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// SendIgnorePlayer sends a request to ignore or unignore a player.
func (pcm *PrivateChatManager) SendIgnorePlayer(target string, flag int) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("ignore_player")
	if !exists {
		return fmt.Errorf("ignore_player packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"target": target,
		"flag":   flag,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// SendTalk sends a talk request to an NPC.
func (pcm *PrivateChatManager) SendTalk(npcID uint32) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("talk")
	if !exists {
		return fmt.Errorf("talk packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"npc_id": npcID,
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}

// SendTalkResponse sends a response to an NPC dialog.
func (pcm *PrivateChatManager) SendTalkResponse(npcID uint32, response int) error {
	// Get the packet ID
	packetID, exists := pcm.baseSend.GetPacketID("talk_response")
	if !exists {
		return fmt.Errorf("talk_response packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"npc_id":   npcID,
		"response": uint8(response),
	}

	// Construct and send the packet
	packet, err := pcm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pcm.baseSend.SendToServer(packet)
}
