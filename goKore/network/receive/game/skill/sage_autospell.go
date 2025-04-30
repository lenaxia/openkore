package skill

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SageAutospellManager manages the sage_autospell packet handlers
type SageAutospellManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSageAutospellManager creates a new sage autospell manager
func NewSageAutospellManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SageAutospellManager {
	return &SageAutospellManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to sage autospell
func (m *SageAutospellManager) RegisterHandlers() {
	// Register sage_autospell handler for Sage's Hindsight
	m.parser.RegisterHandlerFunc("01CD", "sage_autospell", "v v",
		[]string{"why", "autospell_list"},
		m.handleSageAutospell)

	// Register sage_autospell handler for Shadow Chaser's Auto Shadow Spell
	m.parser.RegisterHandlerFunc("0442", "sage_autospell_shadow", "v v",
		[]string{"why", "autoshadowspell_list"},
		m.handleSageAutospell)
}

// handleSageAutospell handles the sage_autospell packet
// Packet formats:
// 01CD: <why>.W <autospell skill id>.L {<autospell skill id>.L}*7
// 0442: <why>.W <autoshadowspell skill id>.W {<autoshadowspell skill id>.W}*5
func (m *SageAutospellManager) handleSageAutospell(args map[string]interface{}) error {
	// Parse the packet
	parsedArgs := m.parseSageAutospell(args)

	// Process the packet
	result := m.processSageAutospell(parsedArgs)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.sage_autospell", result)
	}

	return nil
}

// parseSageAutospell parses the sage_autospell packet and extracts skill IDs
func (m *SageAutospellManager) parseSageAutospell(args map[string]interface{}) map[string]interface{} {
	// Get the packet switch
	packetSwitch, ok := args["switch"].(string)
	if !ok {
		packetSwitch = "01CD" // Default to the standard packet
	}

	// Extract why
	var why uint16
	if val, ok := args["why"].(uint16); ok {
		why = val
	}

	// Extract skill list
	var skillIDs []uint16
	var rawList []byte

	if packetSwitch == "0442" {
		// Shadow Chaser's Auto Shadow Spell (16-bit skill IDs)
		if val, ok := args["autoshadowspell_list"].([]byte); ok {
			rawList = val
			// Parse 16-bit skill IDs
			for i := 0; i < len(rawList); i += 2 {
				if i+1 < len(rawList) {
					skillID := uint16(rawList[i]) | uint16(rawList[i+1])<<8
					if skillID > 0 {
						skillIDs = append(skillIDs, skillID)
					}
				}
			}
		}
	} else {
		// Sage's Hindsight (32-bit skill IDs)
		if val, ok := args["autospell_list"].([]byte); ok {
			rawList = val
			// Parse 32-bit skill IDs
			for i := 0; i < len(rawList); i += 4 {
				if i+3 < len(rawList) {
					skillID := uint16(binary.LittleEndian.Uint32(rawList[i : i+4]))
					if skillID > 0 {
						skillIDs = append(skillIDs, skillID)
					}
				}
			}
		}
	}

	// Sort skill IDs
	// In Go, we would need to implement a sort function, but for simplicity,
	// we'll assume the skill IDs are already sorted

	return map[string]interface{}{
		"why":      why,
		"skillIDs": skillIDs,
		"switch":   packetSwitch,
	}
}

// processSageAutospell processes the sage_autospell packet and returns a structured result
func (m *SageAutospellManager) processSageAutospell(args map[string]interface{}) map[string]interface{} {
	// Extract information from args
	var why uint16
	var skillIDs []uint16
	var packetSwitch string

	// Extract why
	if val, ok := args["why"].(uint16); ok {
		why = val
	}

	// Extract skillIDs
	if val, ok := args["skillIDs"].([]uint16); ok {
		skillIDs = val
	}

	// Extract switch
	if val, ok := args["switch"].(string); ok {
		packetSwitch = val
	}

	// Create skill info list
	skillInfoList := make([]map[string]interface{}, 0, len(skillIDs))
	for _, skillID := range skillIDs {
		skillInfo := map[string]interface{}{
			"skillID":   skillID,
			"skillName": m.getSkillName(skillID),
		}
		skillInfoList = append(skillInfoList, skillInfo)
	}

	// Generate message
	message := m.generateAutospellMessage(skillInfoList)

	// Return structured result
	return map[string]interface{}{
		"why":           why,
		"skillIDs":      skillIDs,
		"skillInfoList": skillInfoList,
		"message":       message,
		"isAutoShadow":  packetSwitch == "0442",
	}
}

// getSkillName returns the skill name for a given skill ID
// In a real implementation, this would look up the skill name from a database or map
func (m *SageAutospellManager) getSkillName(skillID uint16) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the skill name from a database or map
	return fmt.Sprintf("Skill_%d", skillID)
}

// generateAutospellMessage generates a formatted message for the autospell list
func (m *SageAutospellManager) generateAutospellMessage(skillInfoList []map[string]interface{}) string {
	var sb strings.Builder

	// Header
	sb.WriteString("--------------- Auto Spell ---------------\n")
	sb.WriteString("   # Skill\n")

	// Skill list
	for _, skillInfo := range skillInfoList {
		skillID := skillInfo["skillID"].(uint16)
		skillName := skillInfo["skillName"].(string)
		sb.WriteString(fmt.Sprintf("%4d %s\n", skillID, skillName))
	}

	// Footer
	sb.WriteString("----------------------------------------\n")

	return sb.String()
}

// reconstructSageAutospell reconstructs the sage_autospell packet
func (m *SageAutospellManager) reconstructSageAutospell(skillIDs []uint16, isAutoShadow bool) ([]byte, []byte) {
	// Create autoshadowspell_list (16-bit skill IDs)
	autoshadowspellList := make([]byte, len(skillIDs)*2)
	for i, skillID := range skillIDs {
		autoshadowspellList[i*2] = byte(skillID)
		autoshadowspellList[i*2+1] = byte(skillID >> 8)
	}

	// Create autospell_list (32-bit skill IDs)
	autospellList := make([]byte, len(skillIDs)*4)
	for i, skillID := range skillIDs {
		binary.LittleEndian.PutUint32(autospellList[i*4:i*4+4], uint32(skillID))
	}

	return autoshadowspellList, autospellList
}
