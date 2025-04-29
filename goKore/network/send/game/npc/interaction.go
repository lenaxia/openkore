// Package npc provides NPC-related packet sending functionality.
package npc

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// InteractionManager handles NPC interaction-related packet sending.
type InteractionManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewInteractionManager creates a new NPC interaction manager.
func NewInteractionManager(baseSend core.Send) *InteractionManager {
	return &InteractionManager{
		baseSend: baseSend,
	}
}

// SendTalk sends a talk request to an NPC.
// Corresponds to sendTalk in the original implementation.
func (im *InteractionManager) SendTalk(npcID uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk")
	if !exists {
		return fmt.Errorf("npc_talk packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   npcID,
		"type": 1,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTalkCancel sends a cancel talk request to an NPC.
// Corresponds to sendTalkCancel in the original implementation.
func (im *InteractionManager) SendTalkCancel(npcID uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk_cancel")
	if !exists {
		return fmt.Errorf("npc_talk_cancel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": npcID,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTalkContinue sends a continue talk request to an NPC.
// Corresponds to sendTalkContinue in the original implementation.
func (im *InteractionManager) SendTalkContinue(npcID uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk_continue")
	if !exists {
		return fmt.Errorf("npc_talk_continue packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": npcID,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTalkResponse sends a response to an NPC dialog.
// Corresponds to sendTalkResponse in the original implementation.
func (im *InteractionManager) SendTalkResponse(npcID uint32, response int) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk_response")
	if !exists {
		return fmt.Errorf("npc_talk_response packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":       npcID,
		"response": response,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTalkNumber sends a number response to an NPC dialog.
// Corresponds to sendTalkNumber in the original implementation.
func (im *InteractionManager) SendTalkNumber(npcID uint32, number int) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk_number")
	if !exists {
		return fmt.Errorf("npc_talk_number packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     npcID,
		"number": number,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTalkText sends a text response to an NPC dialog.
// Corresponds to sendTalkText in the original implementation.
func (im *InteractionManager) SendTalkText(npcID uint32, text string) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("npc_talk_text")
	if !exists {
		return fmt.Errorf("npc_talk_text packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   npcID,
		"text": text,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}
