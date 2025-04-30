package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillUseFailedManager manages the skill_use_failed packet handler
type SkillUseFailedManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillUseFailedManager creates a new skill use failed manager
func NewSkillUseFailedManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillUseFailedManager {
	return &SkillUseFailedManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to skill use failure
func (m *SkillUseFailedManager) RegisterHandlers() {
	// Register skill_use_failed handler
	m.parser.RegisterHandlerFunc("0110", "skill_use_failed", "v2 V2 C2",
		[]string{"skillID", "btype", "itemId", "flag", "cause", "unknown"},
		m.handleSkillUseFailed)

	// Register skill_use_failed handler (alternative packet)
	m.parser.RegisterHandlerFunc("01CD", "skill_use_failed_expanded", "v2 V2 C2",
		[]string{"skillID", "btype", "itemId", "flag", "cause", "unknown"},
		m.handleSkillUseFailed)
}

// handleSkillUseFailed handles the skill_use_failed packet
// Packet formats:
// 0110: <skill id>.W <btype>.W <item id>.L <flag>.L <cause>.B <unknown>.B
// 01CD: <skill id>.W <btype>.W <item id>.L <flag>.L <cause>.B <unknown>.B
func (m *SkillUseFailedManager) handleSkillUseFailed(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillUseFailed(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_use_failed", result)
	}

	return nil
}

// processSkillUseFailed processes the skill_use_failed packet and returns a structured result
func (m *SkillUseFailedManager) processSkillUseFailed(args map[string]interface{}) map[string]interface{} {
	// Extract skill use failed information from args
	var skillID uint16
	var btype uint16
	var itemID uint32
	var flag uint32
	var cause uint8

	// Extract skillID
	if val, ok := args["skillID"].(uint16); ok {
		skillID = val
	}

	// Extract btype
	if val, ok := args["btype"].(uint16); ok {
		btype = val
	}

	// Extract itemID
	if val, ok := args["itemId"].(uint32); ok {
		itemID = val
	}

	// Extract flag
	if val, ok := args["flag"].(uint32); ok {
		flag = val
	}

	// Extract cause
	if val, ok := args["cause"].(uint8); ok {
		cause = val
	}

	// Get the error message based on the cause
	errorMessage := m.getErrorMessage(skillID, btype, cause, itemID)

	// Handle special cases for homunculus skills
	isHomunculus := m.handleHomunculus(skillID, cause)

	// Return structured result
	return map[string]interface{}{
		"skillID":      skillID,
		"btype":        btype,
		"itemID":       itemID,
		"flag":         flag,
		"cause":        cause,
		"errorMessage": errorMessage,
		"isHomunculus": isHomunculus,
	}
}

// getErrorMessage returns the error message based on the cause
func (m *SkillUseFailedManager) getErrorMessage(skillID uint16, btype uint16, cause uint8, itemID uint32) string {
	// Base fail type messages
	baseFailType := map[uint16]string{
		0: "Skill failed",
		1: "No emotions",
		2: "No sit",
		3: "No chat",
		4: "No party",
		5: "No shout",
		6: "No PKing",
		7: "No aligning",
	}

	// Fail type messages
	failType := map[uint8]string{
		0:  "Basic",
		1:  "Insufficient SP",
		2:  "Insufficient HP",
		3:  "No Memo",
		4:  "Mid-Delay",
		5:  "No Zeny",
		6:  "Wrong Weapon Type",
		7:  "Red Gem Needed",
		8:  "Blue Gem Needed",
		9:  "90% Overweight",
		10: "Requirement",
		11: "Failed to use in Target",
		12: "Maximum Ancilla exceed",
		13: "Need this within the Holy water",
		14: "Missing Ancilla",
		19: "Full Amulet",
		24: "[Purchase Street Stall License] need 1",
		29: "Must have at least 1% of base XP",
		30: "Insufficient SP",
		33: "Failed to use Madogear",
		34: "Kunai is Required",
		37: "Canon ball is Required",
		43: "Failed to use Guillotine Poison",
		50: "Failed to use Madogear",
		71: "Missing Required Item",
		72: "Equipment is required",
		73: "Combo Skill Failed",
		76: "Too many HP",
		77: "Need Royal Guard Branding",
		78: "Required Equiped Weapon Class",
		83: "Location not allowed to create chatroom/market",
		84: "Need more bullet",
	}

	// Determine the error message
	var errorMessage string
	if skillID == 1 && cause == 0 && baseFailType[btype] != "" {
		errorMessage = baseFailType[btype]
	} else if failType[cause] != "" {
		errorMessage = failType[cause]
		if cause == 71 {
			errorMessage += fmt.Sprintf(" - item %d", itemID)
		}
	} else {
		errorMessage = "Unknown error"
	}

	return errorMessage
}

// handleHomunculus handles special cases for homunculus skills
func (m *SkillUseFailedManager) handleHomunculus(skillID uint16, cause uint8) bool {
	isHomunculus := false

	// Ressurect Homunculus failed - which means we have no dead homunculus
	if skillID == 247 && cause == 0 {
		isHomunculus = true
	}

	// Call Homunculus failed
	if skillID == 243 {
		if cause == 0 {
			// No vaporized homunculus
			isHomunculus = true
		} else if cause == 71 {
			// Missing item - which means we have a vaporized homunculus
			isHomunculus = true
		}
	}

	return isHomunculus
}
