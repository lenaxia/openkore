// Package field provides handlers for field-related packets.
package field

import (
	"fmt"
)

// RegisterBossHandlers registers all handlers related to boss monsters
func (m *FieldManager) RegisterBossHandlers() {
	// Register boss_map_info handler
	if m.parser != nil {
		m.parser.RegisterHandlerFunc("0293", "boss_map_info", "B Z24 v2 B2",
			[]string{"flag", "name", "x", "y", "hours", "minutes"},
			m.handleBossMapInfo)
	}
}

// handleBossMapInfo handles the boss_map_info packet
// Packet format: 0293 <flag>.B <name>.Z24 <x>.W <y>.W <hours>.B <minutes>.B
func (m *FieldManager) handleBossMapInfo(args map[string]interface{}) error {
	// Check required fields
	if _, ok := args["flag"].(byte); !ok {
		return fmt.Errorf("missing flag in boss_map_info packet")
	}

	// Process the packet
	result := m.processBossMapInfo(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("field.boss_map_info", result)
	}

	return nil
}

// processBossMapInfo processes the boss_map_info packet and returns a structured result
func (m *FieldManager) processBossMapInfo(args map[string]interface{}) map[string]interface{} {
	var flag byte
	var status string
	var bossName string

	// Extract flag from args
	if flagVal, ok := args["flag"].(byte); ok {
		flag = flagVal
	}

	// Extract boss name from args
	if nameBytes, ok := args["name"].([]byte); ok {
		// Convert bytes to string and trim null bytes
		bossName = bytesToString(nameBytes)
	}

	// Process based on flag value
	switch flag {
	case 0:
		status = "You cannot find any trace of a Boss Monster in this area."
	case 1:
		// Extract coordinates
		var x, y uint16
		if xVal, ok := args["x"].(uint16); ok {
			x = xVal
		}
		if yVal, ok := args["y"].(uint16); ok {
			y = yVal
		}
		status = fmt.Sprintf("MVP Boss %s is now on location: (%d, %d)", bossName, x, y)
	case 2:
		status = fmt.Sprintf("MVP Boss %s has been detected on this map!", bossName)
	case 3:
		// Extract respawn time
		var hours, minutes byte
		if hoursVal, ok := args["hours"].(byte); ok {
			hours = hoursVal
		}
		if minutesVal, ok := args["minutes"].(byte); ok {
			minutes = minutesVal
		}
		status = fmt.Sprintf("MVP Boss %s is dead, but will spawn again in %d hour(s) and %d minutes(s).", bossName, hours, minutes)
	default:
		status = fmt.Sprintf("Unknown boss_map_info result (flag: %d)", flag)
	}

	// Create the base result
	result := map[string]interface{}{
		"flag":     flag,
		"bossName": bossName,
		"status":   status,
	}

	// Add additional fields based on flag
	switch flag {
	case 1:
		// Add coordinates for flag 1
		if x, ok := args["x"].(uint16); ok {
			result["x"] = x
		}
		if y, ok := args["y"].(uint16); ok {
			result["y"] = y
		}
	case 3:
		// Add respawn time for flag 3
		if hours, ok := args["hours"].(byte); ok {
			result["hours"] = hours
		}
		if minutes, ok := args["minutes"].(byte); ok {
			result["minutes"] = minutes
		}
	}

	return result
}

// bytesToString converts a byte slice to a string, trimming null bytes
func bytesToString(b []byte) string {
	n := 0
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			break
		}
		n++
	}
	return string(b[:n])
}
