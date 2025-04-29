// Package social provides social-related packet sending functionality.
package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// PartyManager handles party-related packet sending.
type PartyManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewPartyManager creates a new party manager.
func NewPartyManager(baseSend core.Send) *PartyManager {
	return &PartyManager{
		baseSend: baseSend,
	}
}

// CreateParty sends a request to create a party.
func (pm *PartyManager) CreateParty(name string, shareExp, shareItems bool) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("create_party")
	if !exists {
		return fmt.Errorf("create_party packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name":        name,
		"share_exp":   shareExp,
		"share_items": shareItems,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// JoinParty sends a request to join a party.
func (pm *PartyManager) JoinParty(partyID uint32) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("join_party")
	if !exists {
		return fmt.Errorf("join_party packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"party_id": partyID,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// LeaveParty sends a request to leave the current party.
func (pm *PartyManager) LeaveParty() error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("leave_party")
	if !exists {
		return fmt.Errorf("leave_party packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// InviteToParty sends a request to invite a player to the party.
func (pm *PartyManager) InviteToParty(playerName string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("invite_to_party")
	if !exists {
		return fmt.Errorf("invite_to_party packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"player_name": playerName,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// KickFromParty sends a request to kick a player from the party.
func (pm *PartyManager) KickFromParty(playerName string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("kick_from_party")
	if !exists {
		return fmt.Errorf("kick_from_party packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"player_name": playerName,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// PartyChat sends a message to the party chat.
func (pm *PartyManager) PartyChat(message string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_chat")
	if !exists {
		return fmt.Errorf("party_chat packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// ChangePartyLeader sends a request to change the party leader.
func (pm *PartyManager) ChangePartyLeader(playerName string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("change_party_leader")
	if !exists {
		return fmt.Errorf("change_party_leader packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"player_name": playerName,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// ChangePartyOption sends a request to change party options.
func (pm *PartyManager) ChangePartyOption(shareExp, shareItems bool) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("change_party_option")
	if !exists {
		return fmt.Errorf("change_party_option packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"share_exp":   shareExp,
		"share_items": shareItems,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// PartyBookingRegister sends a request to register a party in the party booking system.
func (pm *PartyManager) PartyBookingRegister(level int, job int, message string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_booking_register")
	if !exists {
		return fmt.Errorf("party_booking_register packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"level":   level,
		"job":     job,
		"message": message,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// PartyBookingSearch sends a request to search for parties in the party booking system.
func (pm *PartyManager) PartyBookingSearch(level int, job int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_booking_search")
	if !exists {
		return fmt.Errorf("party_booking_search packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"level": level,
		"job":   job,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}
