// Package guild provides guild-related packet sending functionality.
package guild

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// GuildManager handles guild-related packet sending.
type GuildManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewGuildManager creates a new guild manager.
func NewGuildManager(baseSend core.Send) *GuildManager {
	return &GuildManager{
		baseSend: baseSend,
	}
}

// SendGuildMasterMemberCheck sends a request to check guild master/member status.
// This is equivalent to the sendGuildMasterMemberCheck function in Send.pm.
func (gm *GuildManager) SendGuildMasterMemberCheck() error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_check")
	if !exists {
		return fmt.Errorf("guild_check packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildRequestInfo sends a request for guild information.
// This is equivalent to the sendGuildRequestInfo function in Send.pm.
func (gm *GuildManager) SendGuildRequestInfo(page int) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_info_request")
	if !exists {
		return fmt.Errorf("guild_info_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": page,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildAlly sends a response to a guild alliance request.
// This is equivalent to the sendGuildAlly function in Send.pm.
func (gm *GuildManager) SendGuildAlly(ID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_alliance_reply")
	if !exists {
		return fmt.Errorf("guild_alliance_reply packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   ID,
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildRequestEmblem sends a request for a guild emblem.
// This is equivalent to the sendGuildRequestEmblem function in Send.pm.
func (gm *GuildManager) SendGuildRequestEmblem(guildID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_emblem_request")
	if !exists {
		return fmt.Errorf("guild_emblem_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guildID": guildID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildBreak sends a request to disband a guild.
// This is equivalent to the sendGuildBreak function in Send.pm.
func (gm *GuildManager) SendGuildBreak(guildName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_break")
	if !exists {
		return fmt.Errorf("guild_break packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guildName": []byte(guildName), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildLeave sends a request to leave a guild.
// This is equivalent to the sendGuildLeave function in Send.pm.
func (gm *GuildManager) SendGuildLeave(reason string, guildID, accountID, charID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_leave")
	if !exists {
		return fmt.Errorf("guild_leave packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guildID":   guildID,
		"accountID": accountID,
		"charID":    charID,
		"reason":    []byte(reason), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildMemberKick sends a request to kick a member from a guild.
// This is equivalent to the sendGuildMemberKick function in Send.pm.
func (gm *GuildManager) SendGuildMemberKick(guildID, accountID, charID uint32, reason string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_kick")
	if !exists {
		return fmt.Errorf("guild_kick packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guildID":   guildID,
		"accountID": accountID,
		"charID":    charID,
		"reason":    []byte(reason), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildCreate sends a request to create a guild.
// This is equivalent to the sendGuildCreate function in Send.pm.
func (gm *GuildManager) SendGuildCreate(name string, charID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_create")
	if !exists {
		return fmt.Errorf("guild_create packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"charID":    charID,
		"guildName": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildJoin sends a response to a guild join request.
// This is equivalent to the sendGuildJoin function in Send.pm.
func (gm *GuildManager) SendGuildJoin(ID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_join")
	if !exists {
		return fmt.Errorf("guild_join packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   ID,
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildJoinRequest sends a request to join a guild.
// This is equivalent to the sendGuildJoinRequest function in Send.pm.
func (gm *GuildManager) SendGuildJoinRequest(ID, accountID, charID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_join_request")
	if !exists {
		return fmt.Errorf("guild_join_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":        ID,
		"accountID": accountID,
		"charID":    charID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildSetAlly sends a request to set an alliance with another guild.
// This is equivalent to the sendGuildSetAlly function in Send.pm.
func (gm *GuildManager) SendGuildSetAlly(targetAID, myAID, charID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_alliance_request")
	if !exists {
		return fmt.Errorf("guild_alliance_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"targetAccountID": targetAID,
		"accountID":       myAID,
		"charID":          charID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGuildNotice sends a request to update the guild notice.
// This is equivalent to the sendGuildNotice function in Send.pm.
func (gm *GuildManager) SendGuildNotice(guildID uint32, name, notice string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_notice")
	if !exists {
		return fmt.Errorf("guild_notice packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guildID": guildID,
		"name":    []byte(name),   // stringToBytes in the original
		"notice":  []byte(notice), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}
