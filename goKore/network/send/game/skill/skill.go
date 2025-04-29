// Package skill provides skill-related packet sending functionality.
package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// SkillManager handles skill-related packet sending.
type SkillManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewSkillManager creates a new skill manager.
func NewSkillManager(baseSend core.Send) *SkillManager {
	return &SkillManager{
		baseSend: baseSend,
	}
}

// SendSkillUse sends a request to use a skill on a target.
// This is equivalent to the sendSkillUse function in Send.pm.
func (sm *SkillManager) SendSkillUse(skillID, level uint16, targetID uint32) error {
	// In the original Perl code, there's a hook here, but we'll simplify for now
	// since the hooks system might be different in Go

	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("skill_use")
	if !exists {
		return fmt.Errorf("skill_use packet ID not found")
	}

	// Create the arguments
	packetArgs := map[string]interface{}{
		"skillID":  skillID,
		"lv":       level,
		"targetID": targetID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, packetArgs)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendSkillUseLoc sends a request to use a skill on a location.
// This is equivalent to the sendSkillUseLoc function in Send.pm.
func (sm *SkillManager) SendSkillUseLoc(skillID, level, x, y uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("skill_use_location")
	if !exists {
		return fmt.Errorf("skill_use_location packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID": skillID,
		"lv":      level,
		"x":       x,
		"y":       y,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendSkillUseLocInfo sends a request to use a skill on a location with additional info.
// This is equivalent to the sendSkillUseLocInfo function in Send.pm.
func (sm *SkillManager) SendSkillUseLocInfo(skillID, level, x, y uint16, info string) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("skill_use_location_text")
	if !exists {
		return fmt.Errorf("skill_use_location_text packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   skillID,
		"lvl":  level,
		"x":    x,
		"y":    y,
		"info": info,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendSkillSelect sends a request to select a skill.
// This is equivalent to the sendSkillSelect function in Send.pm.
func (sm *SkillManager) SendSkillSelect(skillID, why uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("skill_select")
	if !exists {
		return fmt.Errorf("skill_select packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID": skillID,
		"why":     why,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendStartSkillUse sends a request to start using a continuous skill.
// This is equivalent to the sendStartSkillUse function in Send.pm.
func (sm *SkillManager) SendStartSkillUse(skillID, level uint16, targetID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("start_skill_use")
	if !exists {
		return fmt.Errorf("start_skill_use packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID":  skillID,
		"lv":       level,
		"targetID": targetID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendStopSkillUse sends a request to stop using a continuous skill.
// This is equivalent to the sendStopSkillUse function in Send.pm.
func (sm *SkillManager) SendStopSkillUse(skillID uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("stop_skill_use")
	if !exists {
		return fmt.Errorf("stop_skill_use packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID": skillID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendAutoSpell sends a request to set auto-spell.
// This is equivalent to the sendAutoSpell function in Send.pm.
func (sm *SkillManager) SendAutoSpell(skillID uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("auto_spell")
	if !exists {
		return fmt.Errorf("auto_spell packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": skillID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}
