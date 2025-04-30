package skill

import (
	"strconv"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillDelayManager manages the skill_post_delay and skill_post_delaylist packet handlers
type SkillDelayManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillDelayManager creates a new skill delay manager
func NewSkillDelayManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillDelayManager {
	return &SkillDelayManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// SkillDelayInfo represents information about a skill delay
type SkillDelayInfo struct {
	ID         uint16 // Skill ID
	RemainTime uint32 // Remaining time in milliseconds
	TotalTime  uint32 // Total time in milliseconds (only for 0985 packet)
	SkillName  string // Skill name (looked up from ID)
}

// RegisterHandlers registers all handlers related to skill delays
func (m *SkillDelayManager) RegisterHandlers() {
	// Register skill_post_delay handler
	m.parser.RegisterHandlerFunc("043D", "skill_post_delay", "v V",
		[]string{"ID", "time"},
		m.handleSkillPostDelay)

	// Register skill_post_delaylist handler
	m.parser.RegisterHandlerFunc("043E", "skill_post_delaylist", "v",
		[]string{"skill_list"},
		m.handleSkillPostDelayList)

	// Register skill_post_delaylist handler (alternative packet)
	m.parser.RegisterHandlerFunc("0985", "skill_post_delaylist_expanded", "v",
		[]string{"skill_list"},
		m.handleSkillPostDelayList)
}

// handleSkillPostDelay handles the skill_post_delay packet
// Packet format: 043D <skill ID>.W <tick>.L
func (m *SkillDelayManager) handleSkillPostDelay(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillPostDelay(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_post_delay", result)
	}

	return nil
}

// processSkillPostDelay processes the skill_post_delay packet and returns a structured result
func (m *SkillDelayManager) processSkillPostDelay(args map[string]interface{}) map[string]interface{} {
	// Extract skill delay information from args
	var skillID uint16
	var time uint32

	// Extract skillID
	if val, ok := args["ID"].(uint16); ok {
		skillID = val
	}

	// Extract time
	if val, ok := args["time"].(uint32); ok {
		time = val
	}

	// Create skill delay info
	skillDelay := SkillDelayInfo{
		ID:         skillID,
		RemainTime: time,
		// In a real implementation, we would look up the skill name from the ID
		SkillName: m.getSkillName(skillID),
	}

	// Return structured result
	return map[string]interface{}{
		"skillDelay": skillDelay,
	}
}

// handleSkillPostDelayList handles the skill_post_delaylist packet
// Packet formats:
// 043E <len>.w { <skill ID>.W <tick>.L }*
// 0985 <len>.w { <skill ID>.W <total time>.L <tick>.L }*
func (m *SkillDelayManager) handleSkillPostDelayList(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillPostDelayList(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_post_delaylist", result)
	}

	return nil
}

// processSkillPostDelayList processes the skill_post_delaylist packet and returns a structured result
func (m *SkillDelayManager) processSkillPostDelayList(args map[string]interface{}) map[string]interface{} {
	// Get the packet switch
	packetSwitch, ok := args["switch"].(string)
	if !ok {
		packetSwitch = "043E" // Default to the standard packet
	}

	// Get the skill list
	skillList, ok := args["skill_list"].([]byte)
	if !ok {
		return map[string]interface{}{
			"skillDelays": []SkillDelayInfo{},
		}
	}

	// Determine the format based on the packet switch
	var itemLen int
	var hasTotal bool

	if packetSwitch == "0985" {
		itemLen = 10 // ID (2) + total_time (4) + remain_time (4)
		hasTotal = true
	} else {
		itemLen = 6 // ID (2) + remain_time (4)
		hasTotal = false
	}

	// Parse the skill list
	skillDelays := []SkillDelayInfo{}
	for i := 0; i < len(skillList); i += itemLen {
		// Skip if we don't have enough data
		if i+itemLen > len(skillList) {
			break
		}

		// Extract skill ID
		skillID := uint16(skillList[i]) | uint16(skillList[i+1])<<8

		// Extract remain time
		var remainTime uint32
		if hasTotal {
			remainTime = uint32(skillList[i+6]) | uint32(skillList[i+7])<<8 | uint32(skillList[i+8])<<16 | uint32(skillList[i+9])<<24
		} else {
			remainTime = uint32(skillList[i+2]) | uint32(skillList[i+3])<<8 | uint32(skillList[i+4])<<16 | uint32(skillList[i+5])<<24
		}

		// Extract total time (if available)
		var totalTime uint32
		if hasTotal {
			totalTime = uint32(skillList[i+2]) | uint32(skillList[i+3])<<8 | uint32(skillList[i+4])<<16 | uint32(skillList[i+5])<<24
		}

		// Create skill delay info
		skillDelay := SkillDelayInfo{
			ID:         skillID,
			RemainTime: remainTime,
			TotalTime:  totalTime,
			// In a real implementation, we would look up the skill name from the ID
			SkillName: m.getSkillName(skillID),
		}

		skillDelays = append(skillDelays, skillDelay)
	}

	// Return structured result
	return map[string]interface{}{
		"skillDelays": skillDelays,
	}
}

// getSkillName returns the skill name for a given skill ID
// In a real implementation, this would look up the skill name from a database or map
func (m *SkillDelayManager) getSkillName(skillID uint16) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the skill name from a database or map
	return "Skill_" + strconv.Itoa(int(skillID))
}
