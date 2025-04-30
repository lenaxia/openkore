package skill

import (
	"encoding/binary"
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// AreaSpellManager manages the area_spell packet handler
type AreaSpellManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	spells      map[uint32]*AreaSpell // Map of spells by ID
}

// AreaSpell represents an area spell
type AreaSpell struct {
	ID          uint32
	SourceID    uint32
	X           uint16
	Y           uint16
	Type        uint16
	IsVisible   uint8
	BinID       int
	ScribbleMsg string // Only used for scribble messages (01C9)
}

// NewAreaSpellManager creates a new area spell manager
func NewAreaSpellManager(parser *core.CoreParser, hookManager *hooks.HookManager) *AreaSpellManager {
	return &AreaSpellManager{
		parser:      parser,
		hookManager: hookManager,
		spells:      make(map[uint32]*AreaSpell),
	}
}

// RegisterHandlers registers all handlers related to area spells
func (m *AreaSpellManager) RegisterHandlers() {
	// Register area_spell handler (011F)
	m.parser.RegisterHandlerFunc("011F", "area_spell", "a4 a4 v2 v C",
		[]string{"ID", "sourceID", "x", "y", "type", "isVisible"},
		m.handleAreaSpell)

	// Register area_spell handler (01C9) - with scribble message
	m.parser.RegisterHandlerFunc("01C9", "area_spell_scribble", "a4 a4 v2 v C Z80",
		[]string{"ID", "sourceID", "x", "y", "type", "isVisible", "scribbleMsg"},
		m.handleAreaSpell)

	// Register area_spell handler (08C7) - expanded version
	m.parser.RegisterHandlerFunc("08C7", "area_spell_expanded", "a4 a4 v2 v C",
		[]string{"ID", "sourceID", "x", "y", "type", "isVisible"},
		m.handleAreaSpell)

	// Register area_spell_multiple2 handler (099F)
	m.parser.RegisterHandlerFunc("099F", "area_spell_multiple2", "v Z*",
		[]string{"len", "spellInfo"},
		m.handleAreaSpellMultiple2)

	// Register area_spell_multiple3 handler (09CA)
	m.parser.RegisterHandlerFunc("09CA", "area_spell_multiple3", "v Z*",
		[]string{"len", "spellInfo"},
		m.handleAreaSpellMultiple3)

	// Register area_spell_disappears handler (0120)
	m.parser.RegisterHandlerFunc("0120", "area_spell_disappears", "a4",
		[]string{"ID"},
		m.handleAreaSpellDisappears)
}

// handleAreaSpell handles the area_spell packet
// Packet formats:
// 011F <ID>.L <sourceID>.L <x>.W <y>.W <type>.W <isVisible>.B
// 01C9 <ID>.L <sourceID>.L <x>.W <y>.W <type>.W <isVisible>.B <scribbleMsg>.80B
// 08C7 <ID>.L <sourceID>.L <x>.W <y>.W <type>.W <isVisible>.B
func (m *AreaSpellManager) handleAreaSpell(args map[string]interface{}) error {
	// Process the packet
	result := m.processAreaSpell(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.area_spell", result)
	}

	return nil
}

// processAreaSpell processes the area_spell packet and returns a structured result
func (m *AreaSpellManager) processAreaSpell(args map[string]interface{}) map[string]interface{} {
	// Extract area spell information from args
	var id, sourceID uint32
	var x, y, spellType uint16
	var isVisible uint8
	var scribbleMsg string
	var switchType string

	// Extract switch type
	if val, ok := args["switch"].(string); ok {
		switchType = val
	}

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		id = val
	}

	// Extract sourceID
	if val, ok := args["sourceID"].(uint32); ok {
		sourceID = val
	}

	// Extract x coordinate
	if val, ok := args["x"].(uint16); ok {
		x = val
	}

	// Extract y coordinate
	if val, ok := args["y"].(uint16); ok {
		y = val
	}

	// Extract spell type
	if val, ok := args["type"].(uint16); ok {
		spellType = val
	}

	// Extract isVisible
	if val, ok := args["isVisible"].(uint8); ok {
		isVisible = val
	}

	// Extract scribble message (only for 01C9)
	if switchType == "01C9" {
		if val, ok := args["scribbleMsg"].(string); ok {
			scribbleMsg = val
		}
	}

	// Find or create the spell
	var binID int
	spell, exists := m.spells[id]
	if exists && spell.SourceID == sourceID {
		// Spell already exists, update it
		binID = spell.BinID
	} else {
		// Create a new spell
		binID = m.getNextBinID()
		m.spells[id] = &AreaSpell{
			ID:       id,
			SourceID: sourceID,
			BinID:    binID,
		}
	}

	// Update the spell
	m.spells[id].X = x
	m.spells[id].Y = y
	m.spells[id].Type = spellType
	m.spells[id].IsVisible = isVisible
	if switchType == "01C9" {
		m.spells[id].ScribbleMsg = scribbleMsg
	}

	// Generate message
	message := m.generateAreaSpellMessage(sourceID, x, y, spellType, switchType, scribbleMsg)

	// Return structured result
	result := map[string]interface{}{
		"ID":         id,
		"sourceID":   sourceID,
		"x":          x,
		"y":          y,
		"type":       spellType,
		"isVisible":  isVisible,
		"binID":      binID,
		"message":    message,
		"spellName":  m.getSpellName(spellType),
		"sourceName": m.getActorName(sourceID),
	}

	// Add scribble message if present
	if switchType == "01C9" {
		result["scribbleMsg"] = scribbleMsg
	}

	return result
}

// getNextBinID returns the next available binID
func (m *AreaSpellManager) getNextBinID() int {
	// Find the highest binID
	highest := -1
	for _, spell := range m.spells {
		if spell.BinID > highest {
			highest = spell.BinID
		}
	}
	return highest + 1
}

// generateAreaSpellMessage generates a message for the area spell
func (m *AreaSpellManager) generateAreaSpellMessage(sourceID uint32, x, y, spellType uint16, switchType, scribbleMsg string) string {
	sourceName := m.getActorName(sourceID)
	spellName := m.getSpellName(spellType)

	if spellType == 0x81 {
		return fmt.Sprintf("%s opened Warp Portal on (%d, %d)", sourceName, x, y)
	} else if switchType == "01C9" {
		return fmt.Sprintf("%s has scribbled: %s on (%d, %d)", sourceName, scribbleMsg, x, y)
	} else {
		return fmt.Sprintf("Area effect %s from %s appeared on (%d, %d)", spellName, sourceName, x, y)
	}
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *AreaSpellManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}

// getSpellName returns the name of a spell by type
// In a real implementation, this would look up the spell name from a database or map
func (m *AreaSpellManager) getSpellName(spellType uint16) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the spell name from a database or map
	switch spellType {
	case 0x81:
		return "Warp Portal"
	default:
		return fmt.Sprintf("Spell_%d", spellType)
	}
}

// handleAreaSpellDisappears handles the area_spell_disappears packet
// Packet format: 0120 <ID>.L
func (m *AreaSpellManager) handleAreaSpellDisappears(args map[string]interface{}) error {
	// Process the packet
	result := m.processAreaSpellDisappears(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.area_spell_disappears", result)
	}

	return nil
}

// processAreaSpellDisappears processes the area_spell_disappears packet and returns a structured result
func (m *AreaSpellManager) processAreaSpellDisappears(args map[string]interface{}) map[string]interface{} {
	// Extract area spell ID from args
	var id uint32

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		id = val
	}

	// Get the spell from the spells map
	spell, exists := m.spells[id]
	if !exists {
		// Spell not found, return minimal information
		return map[string]interface{}{
			"ID":      id,
			"message": fmt.Sprintf("Area spell %d disappeared (not found in spell list)", id),
		}
	}

	// Generate message
	message := m.generateAreaSpellDisappearsMessage(spell)

	// Create result
	result := map[string]interface{}{
		"ID":         id,
		"sourceID":   spell.SourceID,
		"x":          spell.X,
		"y":          spell.Y,
		"type":       spell.Type,
		"binID":      spell.BinID,
		"message":    message,
		"spellName":  m.getSpellName(spell.Type),
		"sourceName": m.getActorName(spell.SourceID),
	}

	// Remove the spell from the spells map
	delete(m.spells, id)

	return result
}

// generateAreaSpellDisappearsMessage generates a message for the area spell disappears
func (m *AreaSpellManager) generateAreaSpellDisappearsMessage(spell *AreaSpell) string {
	sourceName := m.getActorName(spell.SourceID)
	spellName := m.getSpellName(spell.Type)
	return fmt.Sprintf("Area effect %s from %s disappeared from (%d, %d)",
		spellName, sourceName, spell.X, spell.Y)
}

// handleAreaSpellMultiple2 handles the area_spell_multiple2 packet
// Packet format: 099F <len>.W <spellInfo>.B*len
func (m *AreaSpellManager) handleAreaSpellMultiple2(args map[string]interface{}) error {
	// Extract length and spell info
	var length uint16
	var spellInfo []byte

	// Extract length
	if val, ok := args["len"].(uint16); ok {
		length = val
	}

	// Extract spell info
	if val, ok := args["spellInfo"].([]byte); ok {
		spellInfo = val
	}

	// Process each spell in the packet
	for i := 0; i < int(length)-4; i += 18 {
		if i+18 > len(spellInfo) {
			break // Prevent out of bounds access
		}

		// Extract spell data
		spellData := spellInfo[i : i+18]
		if len(spellData) < 18 {
			continue // Skip if not enough data
		}

		// Create individual spell args
		spellArgs := map[string]interface{}{
			"switch":    "099F",
			"ID":        binary.LittleEndian.Uint32(spellData[0:4]),
			"sourceID":  binary.LittleEndian.Uint32(spellData[4:8]),
			"x":         binary.LittleEndian.Uint16(spellData[8:10]),
			"y":         binary.LittleEndian.Uint16(spellData[10:12]),
			"type":      binary.LittleEndian.Uint32(spellData[12:16]),
			"range":     uint8(spellData[16]),
			"isVisible": uint8(spellData[17]),
		}

		// Process the individual spell
		result := m.processAreaSpellMultiple(spellArgs)

		// Notify through hooks system
		if m.hookManager != nil {
			m.hookManager.CallHook("character.area_spell_multiple", result)
		}
	}

	return nil
}

// handleAreaSpellMultiple3 handles the area_spell_multiple3 packet
// Packet format: 09CA <len>.W <spellInfo>.B*len
func (m *AreaSpellManager) handleAreaSpellMultiple3(args map[string]interface{}) error {
	// Extract length and spell info
	var length uint16
	var spellInfo []byte

	// Extract length
	if val, ok := args["len"].(uint16); ok {
		length = val
	}

	// Extract spell info
	if val, ok := args["spellInfo"].([]byte); ok {
		spellInfo = val
	}

	// Process each spell in the packet
	for i := 0; i < int(length)-4; i += 19 {
		if i+19 > len(spellInfo) {
			break // Prevent out of bounds access
		}

		// Extract spell data
		spellData := spellInfo[i : i+19]
		if len(spellData) < 19 {
			continue // Skip if not enough data
		}

		// Create individual spell args
		spellArgs := map[string]interface{}{
			"switch":    "09CA",
			"ID":        binary.LittleEndian.Uint32(spellData[0:4]),
			"sourceID":  binary.LittleEndian.Uint32(spellData[4:8]),
			"x":         binary.LittleEndian.Uint16(spellData[8:10]),
			"y":         binary.LittleEndian.Uint16(spellData[10:12]),
			"type":      binary.LittleEndian.Uint32(spellData[12:16]),
			"range":     uint8(spellData[16]),
			"isVisible": uint8(spellData[17]),
			"level":     uint8(spellData[18]),
		}

		// Process the individual spell
		result := m.processAreaSpellMultiple(spellArgs)

		// Notify through hooks system
		if m.hookManager != nil {
			m.hookManager.CallHook("character.area_spell_multiple", result)
		}
	}

	return nil
}

// processAreaSpellMultiple processes an individual area spell from a multiple spell packet
func (m *AreaSpellManager) processAreaSpellMultiple(args map[string]interface{}) map[string]interface{} {
	// Extract area spell information from args
	var id, sourceID uint32
	var x, y uint16
	var spellType uint32
	var isVisible, spellRange, level uint8
	var switchType string

	// Extract switch type
	if val, ok := args["switch"].(string); ok {
		switchType = val
	}

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		id = val
	}

	// Extract sourceID
	if val, ok := args["sourceID"].(uint32); ok {
		sourceID = val
	}

	// Extract x coordinate
	if val, ok := args["x"].(uint16); ok {
		x = val
	}

	// Extract y coordinate
	if val, ok := args["y"].(uint16); ok {
		y = val
	}

	// Extract spell type
	if val, ok := args["type"].(uint32); ok {
		spellType = val
	}

	// Extract range
	if val, ok := args["range"].(uint8); ok {
		spellRange = val
	}

	// Extract isVisible
	if val, ok := args["isVisible"].(uint8); ok {
		isVisible = val
	}

	// Extract level (only for 09CA)
	if switchType == "09CA" {
		if val, ok := args["level"].(uint8); ok {
			level = val
		}
	}

	// Find or create the spell
	var binID int
	spell, exists := m.spells[id]
	if exists && spell.SourceID == sourceID {
		// Spell already exists, update it
		binID = spell.BinID
	} else {
		// Create a new spell
		binID = m.getNextBinID()
		m.spells[id] = &AreaSpell{
			ID:       id,
			SourceID: sourceID,
			BinID:    binID,
		}
	}

	// Update the spell
	m.spells[id].X = x
	m.spells[id].Y = y
	m.spells[id].Type = uint16(spellType)
	m.spells[id].IsVisible = isVisible

	// Generate message
	message := m.generateAreaSpellMultipleMessage(sourceID, x, y, uint16(spellType), spellRange, isVisible, level, switchType)

	// Return structured result
	result := map[string]interface{}{
		"ID":         id,
		"sourceID":   sourceID,
		"x":          x,
		"y":          y,
		"type":       spellType,
		"range":      spellRange,
		"isVisible":  isVisible,
		"binID":      binID,
		"message":    message,
		"spellName":  m.getSpellName(uint16(spellType)),
		"sourceName": m.getActorName(sourceID),
	}

	// Add level if present
	if switchType == "09CA" {
		result["level"] = level
	}

	return result
}

// generateAreaSpellMultipleMessage generates a message for the area spell multiple
func (m *AreaSpellManager) generateAreaSpellMultipleMessage(sourceID uint32, x, y, spellType uint16, spellRange, isVisible, level uint8, switchType string) string {
	sourceName := m.getActorName(sourceID)
	spellName := m.getSpellName(spellType)

	if spellType == 0x81 {
		return fmt.Sprintf("%s opened Warp Portal on (%d, %d)", sourceName, x, y)
	} else if switchType == "09CA" {
		return fmt.Sprintf("Area effect %s (level %d) from %s appeared on (%d, %d), range = %d", spellName, level, sourceName, x, y, spellRange)
	} else {
		return fmt.Sprintf("Area effect %s from %s appeared on (%d, %d), range = %d", spellName, sourceName, x, y, spellRange)
	}
}
