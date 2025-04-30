package skill

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// DevotionManager manages the devotion packet handler
type DevotionManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewDevotionManager creates a new devotion manager
func NewDevotionManager(parser *core.CoreParser, hookManager *hooks.HookManager) *DevotionManager {
	return &DevotionManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to devotion
func (m *DevotionManager) RegisterHandlers() {
	// Register devotion handler
	m.parser.RegisterHandlerFunc("01CF", "devotion", "a4 a20 C",
		[]string{"sourceID", "targetIDs", "range"},
		m.handleDevotion)
}

// handleDevotion handles the devotion packet
// Packet format: 01CF <ID>.L {<target ID>.L}*5 <range>.B
func (m *DevotionManager) handleDevotion(args map[string]interface{}) error {
	// Process the packet
	result := m.processDevotion(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.devotion", result)
	}

	return nil
}

// processDevotion processes the devotion packet and returns a structured result
func (m *DevotionManager) processDevotion(args map[string]interface{}) map[string]interface{} {
	// Extract devotion information from args
	var sourceID uint32
	var targetIDs []byte
	var devotionRange uint8

	// Extract sourceID
	if val, ok := args["sourceID"].(uint32); ok {
		sourceID = val
	}

	// Extract targetIDs
	if val, ok := args["targetIDs"].([]byte); ok {
		targetIDs = val
	}

	// Extract range
	if val, ok := args["range"].(uint8); ok {
		devotionRange = val
	}

	// Parse target IDs
	targets := make([]uint32, 0, 5)
	targetIndices := make(map[uint32]int)

	for i := 0; i < 5; i++ {
		// Check if we have enough data
		if i*4+4 > len(targetIDs) {
			break
		}

		// Extract target ID
		targetID := binary.LittleEndian.Uint32(targetIDs[i*4 : i*4+4])

		// Skip if target ID is 0
		if targetID == 0 {
			break
		}

		// Add target ID to the list
		targets = append(targets, targetID)

		// Store the index
		targetIndices[targetID] = i
	}

	// Generate message
	message := m.generateDevotionMessage(sourceID, targets)

	// Return structured result
	return map[string]interface{}{
		"sourceID":      sourceID,
		"targets":       targets,
		"targetIndices": targetIndices,
		"range":         devotionRange,
		"message":       message,
	}
}

// generateDevotionMessage generates a message for the devotion skill
func (m *DevotionManager) generateDevotionMessage(sourceID uint32, targets []uint32) string {
	var sb strings.Builder

	// Get source actor name
	sourceName := m.getActorName(sourceID)

	// Generate message for each target
	for _, targetID := range targets {
		targetName := m.getActorName(targetID)
		sb.WriteString(fmt.Sprintf("%s has used Devotion on %s\n", sourceName, targetName))
	}

	return sb.String()
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *DevotionManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}
