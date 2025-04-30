package skill

import (
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillCastManager manages the skill_cast packet handler
type SkillCastManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillCastManager creates a new skill cast manager
func NewSkillCastManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillCastManager {
	return &SkillCastManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// CastInfo represents information about a skill being cast
type CastInfo struct {
	SourceID  uint32    // ID of the actor casting the skill
	TargetID  uint32    // ID of the target (if any)
	SkillID   uint16    // ID of the skill being cast
	X         uint16    // X coordinate (if targeting a location)
	Y         uint16    // Y coordinate (if targeting a location)
	Type      uint8     // Type of cast
	Wait      uint32    // Cast time in milliseconds
	StartTime time.Time // Time when casting started
}

// RegisterHandlers registers all handlers related to skill casting
func (m *SkillCastManager) RegisterHandlers() {
	// Register skill_cast handler
	m.parser.RegisterHandlerFunc("013E", "skill_cast", "a4 a4 v2 v C V",
		[]string{"sourceID", "targetID", "x", "y", "skillID", "type", "wait"},
		m.handleSkillCast)

	// Register skill_cast handler (alternative packet)
	m.parser.RegisterHandlerFunc("07FB", "skill_cast_expanded", "a4 a4 v2 v C V2",
		[]string{"sourceID", "targetID", "x", "y", "skillID", "type", "wait", "unknown"},
		m.handleSkillCast)

	// Register skill_cast handler (another alternative packet)
	m.parser.RegisterHandlerFunc("0A1C", "skill_cast_nodamage", "a4 a4 v v C V v",
		[]string{"sourceID", "targetID", "skillID", "unknown", "type", "wait", "unknown2"},
		m.handleSkillCastNoDamage)
}

// handleSkillCast handles the skill_cast packet
// Packet formats:
// 013E: <src ID>.L <dst ID>.L <x>.W <y>.W <skill ID>.W <type>.B <wait time>.L
// 07FB: <src ID>.L <dst ID>.L <x>.W <y>.W <skill ID>.W <type>.B <wait time>.L <unknown>.L
func (m *SkillCastManager) handleSkillCast(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillCast(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_cast", result)
	}

	return nil
}

// handleSkillCastNoDamage handles the skill_cast_nodamage packet
// Packet format: 0A1C: <src ID>.L <dst ID>.L <skill ID>.W <unknown>.W <type>.B <wait time>.L <unknown>.W
func (m *SkillCastManager) handleSkillCastNoDamage(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillCast(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_cast", result)
	}

	return nil
}

// processSkillCast processes the skill_cast packet and returns a structured result
func (m *SkillCastManager) processSkillCast(args map[string]interface{}) map[string]interface{} {
	// Extract skill cast information from args
	var sourceID, targetID uint32
	var skillID, x, y uint16
	var castType uint8
	var wait uint32

	// Extract sourceID
	if val, ok := args["sourceID"].(uint32); ok {
		sourceID = val
	}

	// Extract targetID
	if val, ok := args["targetID"].(uint32); ok {
		targetID = val
	}

	// Extract skillID
	if val, ok := args["skillID"].(uint16); ok {
		skillID = val
	}

	// Extract x coordinate
	if val, ok := args["x"].(uint16); ok {
		x = val
	}

	// Extract y coordinate
	if val, ok := args["y"].(uint16); ok {
		y = val
	}

	// Extract cast type
	if val, ok := args["type"].(uint8); ok {
		castType = val
	}

	// Extract wait time
	if val, ok := args["wait"].(uint32); ok {
		wait = val
	}

	// Create cast info
	castInfo := CastInfo{
		SourceID:  sourceID,
		TargetID:  targetID,
		SkillID:   skillID,
		X:         x,
		Y:         y,
		Type:      castType,
		Wait:      wait,
		StartTime: time.Now(),
	}

	// Determine if this is a location-targeted skill or an actor-targeted skill
	isLocationTargeted := x != 0 || y != 0

	// Return structured result
	return map[string]interface{}{
		"castInfo":           castInfo,
		"isLocationTargeted": isLocationTargeted,
		"waitTimeSeconds":    float64(wait) / 1000.0,
	}
}
