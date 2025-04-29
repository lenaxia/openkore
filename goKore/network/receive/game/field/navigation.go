// Package field provides handlers for field-related packets.
package field

import (
	"fmt"
	"strings"
)

// RegisterNavigationHandlers registers all handlers related to navigation
func (m *FieldManager) RegisterNavigationHandlers() {
	// Register navigate_to handler
	if m.parser != nil {
		m.parser.RegisterHandlerFunc("08E2", "navigate_to", "B3 Z16 v3",
			[]string{"type", "flag", "hide", "map", "x", "y", "mob_id"},
			m.handleNavigateTo)

		// Register warp_portal_list handler
		m.parser.RegisterHandlerFunc("0B1B", "warp_portal_list", "B Z16 Z16 Z16 Z16",
			[]string{"type", "memo1", "memo2", "memo3", "memo4"},
			m.handleWarpPortalList)
	}
}

// handleNavigateTo handles the navigate_to packet
// Packet format: 08E2 <type>.B <flag>.B <hide>.B <map>.16B <x pos>.W <y pos>.W <mob id>.W
func (m *FieldManager) handleNavigateTo(args map[string]interface{}) error {
	// Process the packet
	result := m.processNavigateTo(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("field.navigate_to", result)
	}

	return nil
}

// processNavigateTo processes the navigate_to packet and returns a structured result
func (m *FieldManager) processNavigateTo(args map[string]interface{}) map[string]interface{} {
	var status string
	var mapName string
	var x, y uint16
	var mobID uint16
	var hasMobID bool

	// Extract type, flag, hide from args (not used in the original implementation)
	// var typeVal, flagVal, hideVal byte
	// if val, ok := args["type"].(byte); ok {
	// 	typeVal = val
	// }
	// if val, ok := args["flag"].(byte); ok {
	// 	flagVal = val
	// }
	// if val, ok := args["hide"].(byte); ok {
	// 	hideVal = val
	// }

	// Extract map name from args
	if mapBytes, ok := args["map"].([]byte); ok {
		// Convert bytes to string and trim null bytes
		mapName = bytesToString(mapBytes)
	}

	// Extract coordinates from args
	if xVal, ok := args["x"].(uint16); ok {
		x = xVal
	}
	if yVal, ok := args["y"].(uint16); ok {
		y = yVal
	}

	// Extract mob ID from args
	if mobIDVal, ok := args["mob_id"].(uint16); ok {
		mobID = mobIDVal
		hasMobID = mobID > 0
	}

	// Generate status message based on whether we have a mob ID
	if hasMobID {
		status = fmt.Sprintf("Server asked us to navigate to %s map and look for monster with ID %d", mapName, mobID)
	} else {
		status = fmt.Sprintf("Server asked us to navigate to %s (%d,%d)", mapName, x, y)
	}

	// Create the result
	result := map[string]interface{}{
		"map":    mapName,
		"x":      x,
		"y":      y,
		"status": status,
	}

	// Add mob_id if present
	if hasMobID {
		result["mob_id"] = mobID
	}

	return result
}

// handleWarpPortalList handles the warp_portal_list packet
// Packet format: 0B1B <type>.B <memo1>.Z16 <memo2>.Z16 <memo3>.Z16 <memo4>.Z16
func (m *FieldManager) handleWarpPortalList(args map[string]interface{}) error {
	// Process the packet
	result := m.processWarpPortalList(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("warp.portal_list", result)

		// Update config based on warp type
		if warpType, ok := args["type"].(byte); ok {
			configUpdate := map[string]interface{}{
				"key": "saveMap",
			}

			if warpType == 26 && args["memo2"] != nil {
				// Teleport skill
				if memo2, ok := args["memo2"].(string); ok && memo2 != "" {
					configUpdate["value"] = stripGatExtension(memo2)
				}
			} else if warpType == 27 && args["memo1"] != nil {
				// Butterfly Wing
				if memo1, ok := args["memo1"].(string); ok && memo1 != "" {
					configUpdate["value"] = stripGatExtension(memo1)
				}
			}

			// Call config.update hook if we have a value
			if _, ok := configUpdate["value"]; ok {
				m.hookManager.CallHook("config.update", configUpdate)
			}
		}
	}

	return nil
}

// processWarpPortalList processes the warp_portal_list packet and returns a structured result
func (m *FieldManager) processWarpPortalList(args map[string]interface{}) map[string]interface{} {
	var warpType byte
	var memoList []string

	// Extract warp type from args
	if typeVal, ok := args["type"].(byte); ok {
		warpType = typeVal
	} else {
		return nil
	}

	// Process memo fields
	for i := 1; i <= 4; i++ {
		memoKey := fmt.Sprintf("memo%d", i)
		if memo, ok := args[memoKey].(string); ok && memo != "" {
			// Strip .gat extension
			memoList = append(memoList, stripGatExtension(memo))
		}
	}

	// Create the result
	result := map[string]interface{}{
		"type":      warpType,
		"memo_list": memoList,
	}

	return result
}

// stripGatExtension removes the .gat extension from a map name
func stripGatExtension(mapName string) string {
	if idx := strings.Index(mapName, ".gat"); idx != -1 {
		return mapName[:idx]
	}
	return mapName
}

// Note: bytesToString is defined in boss.go
