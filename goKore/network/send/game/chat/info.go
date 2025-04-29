// Package chat provides chat-related packet sending functionality.
package chat

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// InfoChatManager handles info chat-related packet sending.
type InfoChatManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewInfoChatManager creates a new info chat manager.
func NewInfoChatManager(baseSend core.Send) *InfoChatManager {
	return &InfoChatManager{
		baseSend: baseSend,
	}
}

// GetManagerName returns the name of the manager.
// This implements the ManagerProvider interface.
func (icm *InfoChatManager) GetManagerName() string {
	return "InfoChatManager"
}

// SendWho sends a request for the number of users online.
// This is equivalent to the sendWho function in Send.pm.
func (icm *InfoChatManager) SendWho() error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("request_user_count")
	if !exists {
		return fmt.Errorf("request_user_count packet ID not found")
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendClanChat sends a clan chat message.
// This is equivalent to the sendClanChat function in Send.pm.
func (icm *InfoChatManager) SendClanChat(message string, charName string) error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("clan_chat")
	if !exists {
		return fmt.Errorf("clan_chat packet ID not found")
	}

	// Format the message
	formattedMessage := charName + " : " + message

	// Create the arguments
	args := map[string]interface{}{
		"message": formattedMessage,
		"len":     uint16(len(formattedMessage) + 4),
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendGetClanInfo sends a request for clan information.
// This is equivalent to the sendGetClanInfo function in Send.pm.
func (icm *InfoChatManager) SendGetClanInfo() error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("request_clan_info")
	if !exists {
		return fmt.Errorf("request_clan_info packet ID not found")
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendClanLeave sends a request to leave a clan.
// This is equivalent to the sendClanLeave function in Send.pm.
func (icm *InfoChatManager) SendClanLeave() error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("clan_leave")
	if !exists {
		return fmt.Errorf("clan_leave packet ID not found")
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendClanMessage sends a clan message.
// This is equivalent to the sendClanMessage function in Send.pm.
func (icm *InfoChatManager) SendClanMessage(message string) error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("clan_message")
	if !exists {
		return fmt.Errorf("clan_message packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendGetPlayerInfo sends a request for player information.
// This is equivalent to the sendGetPlayerInfo function in Send.pm.
func (icm *InfoChatManager) SendGetPlayerInfo(ID uint32) error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("actor_info_request")
	if !exists {
		return fmt.Errorf("get_player_info packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendGetCharacterName sends a request for character name.
// This is equivalent to the sendGetCharacterName function in Send.pm.
func (icm *InfoChatManager) SendGetCharacterName(ID uint32) error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("actor_name_request")
	if !exists {
		return fmt.Errorf("get_character_name packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}

// SendBattlegroundChat sends a battleground chat message.
// This is equivalent to the sendBattlegroundChat function in Send.pm.
func (icm *InfoChatManager) SendBattlegroundChat(message string) error {
	// Get the packet ID
	packetID, exists := icm.baseSend.GetPacketID("battleground_chat")
	if !exists {
		return fmt.Errorf("battleground_chat packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": []byte(message),
	}

	// Construct and send the packet
	packet, err := icm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return icm.baseSend.SendToServer(packet)
}
