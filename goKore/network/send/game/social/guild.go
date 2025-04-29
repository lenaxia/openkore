// Package social provides social-related packet sending functionality.
package social

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

// CreateGuild sends a request to create a guild.
func (gm *GuildManager) CreateGuild(name string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("create_guild")
	if !exists {
		return fmt.Errorf("create_guild packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": name,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// JoinGuild sends a request to join a guild.
func (gm *GuildManager) JoinGuild(guildID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("join_guild")
	if !exists {
		return fmt.Errorf("join_guild packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guild_id": guildID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// LeaveGuild sends a request to leave a guild.
func (gm *GuildManager) LeaveGuild(guildID, accountID, charID uint32, reason string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("leave_guild")
	if !exists {
		return fmt.Errorf("leave_guild packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guild_id":   guildID,
		"account_id": accountID,
		"char_id":    charID,
		"reason":     reason,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// InviteToGuild sends a request to invite a player to the guild.
func (gm *GuildManager) InviteToGuild(accountID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("invite_to_guild")
	if !exists {
		return fmt.Errorf("invite_to_guild packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id": accountID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// GuildChat sends a message to the guild chat.
func (gm *GuildManager) GuildChat(message string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("guild_chat")
	if !exists {
		return fmt.Errorf("guild_chat packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"message": message,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// ChangeGuildPositionInfo sends a request to change guild position information.
func (gm *GuildManager) ChangeGuildPositionInfo(positionID uint32, name string, mode, ranking uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("change_guild_position_info")
	if !exists {
		return fmt.Errorf("change_guild_position_info packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"position_id": positionID,
		"name":        name,
		"mode":        mode,
		"ranking":     ranking,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// ChangeGuildMemberPosition sends a request to change a guild member's position.
func (gm *GuildManager) ChangeGuildMemberPosition(accountID, charID, positionID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("change_guild_member_position")
	if !exists {
		return fmt.Errorf("change_guild_member_position packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"account_id":  accountID,
		"char_id":     charID,
		"position_id": positionID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// BreakGuild sends a request to disband a guild.
func (gm *GuildManager) BreakGuild(guildName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("break_guild")
	if !exists {
		return fmt.Errorf("break_guild packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guild_name": guildName,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// UpdateGuildNotice sends a request to update the guild notice.
func (gm *GuildManager) UpdateGuildNotice(guildID uint32, notice, introduction string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("update_guild_notice")
	if !exists {
		return fmt.Errorf("update_guild_notice packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guild_id":     guildID,
		"notice":       notice,
		"introduction": introduction,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// RequestGuildInfo sends a request to get guild information.
func (gm *GuildManager) RequestGuildInfo(guildID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("request_guild_info")
	if !exists {
		return fmt.Errorf("request_guild_info packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"guild_id": guildID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}
