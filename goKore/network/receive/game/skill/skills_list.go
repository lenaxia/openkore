package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillsListManager manages the skills_list packet handler
type SkillsListManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillsListManager creates a new skills list manager
func NewSkillsListManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillsListManager {
	return &SkillsListManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to skills list
func (m *SkillsListManager) RegisterHandlers() {
	// Register skills_list handler for character skills
	m.parser.RegisterHandlerFunc("010F", "skills_list", "v",
		[]string{"RAW_MSG", "RAW_MSG_SIZE"},
		m.handleSkillsList)

	// Register skills_list handler for homunculus skills
	m.parser.RegisterHandlerFunc("0235", "homun_skills_list", "v",
		[]string{"RAW_MSG", "RAW_MSG_SIZE"},
		m.handleSkillsList)

	// Register skills_list handler for mercenary skills
	m.parser.RegisterHandlerFunc("029D", "merc_skills_list", "v",
		[]string{"RAW_MSG", "RAW_MSG_SIZE"},
		m.handleSkillsList)

	// Register skills_list handler for character skills (alternative packet)
	m.parser.RegisterHandlerFunc("0B32", "skills_list_short", "v",
		[]string{"RAW_MSG", "RAW_MSG_SIZE"},
		m.handleSkillsList)
}

// handleSkillsList handles the skills_list packet
// Packet formats:
// 010F/0235/029D: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <handle>.24B <up>.B
// 0B32: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <up>.B <lv2>.W
func (m *SkillsListManager) handleSkillsList(args map[string]interface{}) error {
	// Process the packet
	result, err := m.processSkillsList(args)
	if err != nil {
		return err
	}

	// Notify through hooks system
	if m.hookManager != nil {
		switch args["switch"].(string) {
		case "010F", "0B32":
			m.hookManager.CallHook("character.skills_list", result)
		case "0235":
			m.hookManager.CallHook("homunculus.skills_list", result)
		case "029D":
			m.hookManager.CallHook("mercenary.skills_list", result)
		}
	}

	return nil
}

// processSkillsList processes the skills_list packet and returns a structured result
func (m *SkillsListManager) processSkillsList(args map[string]interface{}) (map[string]interface{}, error) {
	// Get the packet switch
	packetSwitch, ok := args["switch"].(string)
	if !ok {
		return nil, fmt.Errorf("switch not found in args")
	}

	// Get the raw message
	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return nil, fmt.Errorf("RAW_MSG not found in args")
	}

	// Get the raw message size
	rawMsgSize, ok := args["RAW_MSG_SIZE"].(int)
	if !ok {
		return nil, fmt.Errorf("RAW_MSG_SIZE not found in args")
	}

	// Determine the skill info format based on the packet switch
	var skillInfoLen int
	var ownerType SkillOwnerType

	if packetSwitch == "0B32" {
		skillInfoLen = 15
		ownerType = OwnerChar
	} else {
		skillInfoLen = 37

		// Determine owner type based on packet switch
		switch packetSwitch {
		case "010F":
			ownerType = OwnerChar
		case "0235":
			ownerType = OwnerHomun
		case "029D":
			ownerType = OwnerMerc
		default:
			// Default to character owner type for unknown packet switches
			ownerType = OwnerChar
		}
	}

	// Parse the skills list
	skills := make([]SkillInfo, 0)
	for i := 4; i < rawMsgSize; i += skillInfoLen {
		// Skip if we don't have enough data
		if i+skillInfoLen > len(rawMsg) {
			break
		}

		// Parse the skill info
		skill := SkillInfo{}

		// Extract data based on the format
		if packetSwitch == "0B32" {
			// Format: v V v3 C v
			// ID, targetType, lv, sp, range, up, lv2
			skill.ID = uint16(rawMsg[i]) | uint16(rawMsg[i+1])<<8
			skill.TargetType = uint32(rawMsg[i+2]) | uint32(rawMsg[i+3])<<8 | uint32(rawMsg[i+4])<<16 | uint32(rawMsg[i+5])<<24
			skill.Level = uint16(rawMsg[i+6]) | uint16(rawMsg[i+7])<<8
			skill.SP = uint16(rawMsg[i+8]) | uint16(rawMsg[i+9])<<8
			skill.Range = uint16(rawMsg[i+10]) | uint16(rawMsg[i+11])<<8
			skill.Up = rawMsg[i+12]
			skill.Level2 = uint16(rawMsg[i+13]) | uint16(rawMsg[i+14])<<8
		} else {
			// Format: v V v3 Z24 C
			// ID, targetType, lv, sp, range, handle, up
			skill.ID = uint16(rawMsg[i]) | uint16(rawMsg[i+1])<<8
			skill.TargetType = uint32(rawMsg[i+2]) | uint32(rawMsg[i+3])<<8 | uint32(rawMsg[i+4])<<16 | uint32(rawMsg[i+5])<<24
			skill.Level = uint16(rawMsg[i+6]) | uint16(rawMsg[i+7])<<8
			skill.SP = uint16(rawMsg[i+8]) | uint16(rawMsg[i+9])<<8
			skill.Range = uint16(rawMsg[i+10]) | uint16(rawMsg[i+11])<<8

			// Extract handle (null-terminated string)
			handleBytes := rawMsg[i+12 : i+36]
			for j := 0; j < len(handleBytes); j++ {
				if handleBytes[j] == 0 {
					skill.Handle = string(handleBytes[:j])
					break
				}
			}
			if skill.Handle == "" {
				skill.Handle = string(handleBytes)
			}

			skill.Up = rawMsg[i+36]
		}

		skills = append(skills, skill)
	}

	// Return structured result
	return map[string]interface{}{
		"ownerType": ownerType,
		"skills":    skills,
	}, nil
}
