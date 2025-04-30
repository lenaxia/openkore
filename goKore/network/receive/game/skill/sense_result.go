package skill

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SenseResultManager manages the sense_result packet handler
type SenseResultManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSenseResultManager creates a new sense result manager
func NewSenseResultManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SenseResultManager {
	return &SenseResultManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to sense result
func (m *SenseResultManager) RegisterHandlers() {
	// Register sense_result handler
	m.parser.RegisterHandlerFunc("0213", "sense_result", "v v C C v v C v v C C C C C C C C",
		[]string{"nameID", "level", "size", "race", "def", "mdef", "element", "hp",
			"ice", "earth", "fire", "wind", "poison", "holy", "dark", "spirit", "undead"},
		m.handleSenseResult)
}

// handleSenseResult handles the sense_result packet
// Packet format: 0213 <nameID>.W <level>.W <size>.B <race>.B <def>.W <mdef>.W <element>.B <hp>.W
//
//	<ice>.B <earth>.B <fire>.B <wind>.B <poison>.B <holy>.B <dark>.B <spirit>.B <undead>.B
func (m *SenseResultManager) handleSenseResult(args map[string]interface{}) error {
	// Process the packet
	result := m.processSenseResult(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.sense_result", result)
	}

	return nil
}

// processSenseResult processes the sense_result packet and returns a structured result
func (m *SenseResultManager) processSenseResult(args map[string]interface{}) map[string]interface{} {
	// Extract sense result information from args
	var nameID, level, def, mdef, hp uint16
	var size, race, element, ice, earth, fire, wind, poison, holy, dark, spirit, undead uint8

	// Extract nameID
	if val, ok := args["nameID"].(uint16); ok {
		nameID = val
	}

	// Extract level
	if val, ok := args["level"].(uint16); ok {
		level = val
	}

	// Extract size
	if val, ok := args["size"].(uint8); ok {
		size = val
	}

	// Extract race
	if val, ok := args["race"].(uint8); ok {
		race = val
	}

	// Extract def
	if val, ok := args["def"].(uint16); ok {
		def = val
	}

	// Extract mdef
	if val, ok := args["mdef"].(uint16); ok {
		mdef = val
	}

	// Extract element
	if val, ok := args["element"].(uint8); ok {
		element = val
	}

	// Extract hp
	if val, ok := args["hp"].(uint16); ok {
		hp = val
	}

	// Extract damage modifiers
	if val, ok := args["ice"].(uint8); ok {
		ice = val
	}

	if val, ok := args["earth"].(uint8); ok {
		earth = val
	}

	if val, ok := args["fire"].(uint8); ok {
		fire = val
	}

	if val, ok := args["wind"].(uint8); ok {
		wind = val
	}

	if val, ok := args["poison"].(uint8); ok {
		poison = val
	}

	if val, ok := args["holy"].(uint8); ok {
		holy = val
	}

	if val, ok := args["dark"].(uint8); ok {
		dark = val
	}

	if val, ok := args["spirit"].(uint8); ok {
		spirit = val
	}

	if val, ok := args["undead"].(uint8); ok {
		undead = val
	}

	// Get monster name, size name, race name, and element name
	monsterName := m.getMonsterName(nameID)
	sizeName := m.getSizeName(size)
	raceName := m.getRaceName(race)
	elementName := m.getElementName(element)

	// Generate message
	message := m.generateSenseResultMessage(monsterName, level, sizeName, raceName, def, mdef, elementName, hp,
		ice, earth, fire, wind, poison, holy, dark, spirit, undead)

	// Return structured result
	return map[string]interface{}{
		"nameID":      nameID,
		"monsterName": monsterName,
		"level":       level,
		"size":        size,
		"sizeName":    sizeName,
		"race":        race,
		"raceName":    raceName,
		"def":         def,
		"mdef":        mdef,
		"element":     element,
		"elementName": elementName,
		"hp":          hp,
		"ice":         ice,
		"earth":       earth,
		"fire":        fire,
		"wind":        wind,
		"poison":      poison,
		"holy":        holy,
		"dark":        dark,
		"spirit":      spirit,
		"undead":      undead,
		"message":     message,
	}
}

// getMonsterName returns the name of a monster by ID
// In a real implementation, this would look up the monster name from a database or map
func (m *SenseResultManager) getMonsterName(nameID uint16) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the monster name from a database or map
	return fmt.Sprintf("Monster_%d", nameID)
}

// getSizeName returns the name of a size by ID
func (m *SenseResultManager) getSizeName(size uint8) string {
	sizes := []string{"Small", "Medium", "Large"}
	if int(size) < len(sizes) {
		return sizes[size]
	}
	return "Unknown"
}

// getRaceName returns the name of a race by ID
func (m *SenseResultManager) getRaceName(race uint8) string {
	races := []string{
		"Formless", "Undead", "Beast", "Plant", "Insect", "Fish",
		"Demon", "Demi-Human", "Angel", "Dragon", "Boss", "Non-Boss",
	}
	if int(race) < len(races) {
		return races[race]
	}
	return "Unknown"
}

// getElementName returns the name of an element by ID
// In a real implementation, this would look up the element name from a database or map
func (m *SenseResultManager) getElementName(element uint8) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the element name from a database or map
	elements := []string{
		"Neutral", "Water", "Earth", "Fire", "Wind", "Poison", "Holy", "Dark", "Ghost", "Undead",
	}
	if int(element) < len(elements) {
		return elements[element]
	}
	return "Unknown"
}

// generateSenseResultMessage generates a message for the sense result
func (m *SenseResultManager) generateSenseResultMessage(monsterName string, level uint16, sizeName, raceName string,
	def, mdef uint16, elementName string, hp uint16,
	ice, earth, fire, wind, poison, holy, dark, spirit, undead uint8) string {

	var sb strings.Builder

	sb.WriteString("=====================Sense========================\n")
	sb.WriteString(fmt.Sprintf("Monster: %s Level: %d\n", monsterName, level))
	sb.WriteString(fmt.Sprintf("Size: %s Race: %s\n", sizeName, raceName))
	sb.WriteString(fmt.Sprintf("Def: %d MDef: %d\n", def, mdef))
	sb.WriteString(fmt.Sprintf("Element: %s HP: %d\n", elementName, hp))
	sb.WriteString("=================Damage Modifiers=================\n")
	sb.WriteString(fmt.Sprintf("Ice: %d     Earth: %d  Fire: %d  Wind: %d\n", ice, earth, fire, wind))
	sb.WriteString(fmt.Sprintf("Poison: %d  Holy: %d   Dark: %d  Spirit: %d\n", poison, holy, dark, spirit))
	sb.WriteString(fmt.Sprintf("Undead: %d\n", undead))
	sb.WriteString("==================================================")

	return sb.String()
}
