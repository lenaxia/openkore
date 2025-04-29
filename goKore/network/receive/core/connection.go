// Package core provides core functionality for parsing and processing network packets.
package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lenaxia/goKore/network"
)

// handleConnectionRefused handles the connection_refused packet
func (m *AccountManager) handleConnectionRefused(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update state
	m.networkState = network.NotConnected
	m.session.State = AccountStateLoggedOut
	m.session.LastPacketTime = time.Now()

	// Log the error code if available
	// if errorCode, ok := args["error"].(byte); ok {
	// 	// TODO: Add proper logging when logger is implemented
	// }

	return nil
}

// handleMapLoadError handles the map_load_error packet
func (m *AccountManager) handleMapLoadError(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Keep the current network state, but update the last packet time
	m.session.LastPacketTime = time.Now()

	// Log the error code if available
	// if errorCode, ok := args["error"].(byte); ok {
	// 	// TODO: Add proper logging when logger is implemented
	// }

	return nil
}

// handleReceivedSync handles the received_sync packet
func (m *AccountManager) handleReceivedSync(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Log the sync time if available
	// if syncTime, ok := args["time"].(uint32); ok {
	// 	// TODO: Add proper logging when logger is implemented
	// }

	return nil
}

// handleActorMovementInterrupted handles the actor_movement_interrupted packet
func (m *AccountManager) handleActorMovementInterrupted(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Log the actor ID and position if available
	// if actorID, ok := args["ID"].(uint32); ok {
	// 	posX, posXOk := args["posX"].(uint16)
	// 	posY, posYOk := args["posY"].(uint16)
	// 	if posXOk && posYOk {
	// 		// TODO: Add proper logging when logger is implemented
	// 	} else {
	// 		// TODO: Add proper logging when logger is implemented
	// 	}
	// }

	return nil
}

// handleMapChange handles the map_change packet
func (m *AccountManager) handleMapChange(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update map name if available
	if mapName, ok := args["map"].(string); ok {
		m.session.MapName = mapName
	}

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Log the map change
	// TODO: Add proper logging when logger is implemented

	return nil
}

// handleMapChanged handles the map_changed packet
func (m *AccountManager) handleMapChanged(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update map name if available
	if mapName, ok := args["map"].(string); ok {
		m.session.MapName = mapName
	}

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Log the map change
	// TODO: Add proper logging when logger is implemented

	return nil
}

// handleSyncRequestEx handles the sync_request_ex packet
func (m *AccountManager) handleSyncRequestEx(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Get the packet ID from the args
	packetID, ok := args["switch"].(string)
	if !ok {
		return errors.New("invalid packet ID in sync_request_ex")
	}

	// Get the corresponding sync_ex_reply ID from the table
	syncID, ok := m.syncExReplyTable[packetID]
	if !ok {
		// If not found in the table, log a warning and return
		// TODO: Add proper logging when logger is implemented
		return nil
	}

	// Clean leading zeros from both IDs
	packetID = strings.TrimLeft(packetID, "0")
	syncID = strings.TrimLeft(syncID, "0")

	// Convert the sync ID to a hex number
	_, err := strconv.ParseInt(syncID, 16, 64)
	if err != nil {
		return fmt.Errorf("failed to parse sync ID: %v", err)
	}

	// Send a reply to the server with the sync ID
	// TODO: Implement sendReplySyncRequestEx in the messageSender
	// m.messageSender.sendReplySyncRequestEx(uint32(syncIDHex))

	return nil
}

// handleMoveInterrupt handles the move_interrupt packet
// This packet is sent when movement is interrupted by casting a skill, fleeing a mob, etc.
// Packet ID: 0AB8
func (m *AccountManager) handleMoveInterrupt(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Log the movement interruption
	// TODO: Add proper logging when logger is implemented
	// Debug message: "Movement interrupted by casting a skill/fleeing a mob/etc"

	return nil
}

// Map type constants
const (
	MapTypePvP          = 6
	MapTypeGvG          = 8
	MapTypeBattleground = 19
)

// PvP mode constants
const (
	PvPModeNone         = 0
	PvPModePvP          = 1
	PvPModeGvG          = 2
	PvPModeBattleground = 3
)

// Teleport error constants
const (
	TeleportErrorUnavailableArea = 0
	TeleportErrorUnavailableMemo = 1
)

// handleMapChangeCell handles the map_change_cell packet (0192)
// This packet is sent when a cell on the map changes its type
func (m *AccountManager) handleMapChangeCell(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract packet data
	x, xOk := args["x"].(uint16)
	y, yOk := args["y"].(uint16)
	cellType, _ := args["type"].(byte)
	mapName, mapNameOk := args["map_name"].(string)

	// Validate required fields
	if !xOk || !yOk {
		return errors.New("missing coordinates in map_change_cell packet")
	}
	if !mapNameOk {
		return errors.New("missing map name in map_change_cell packet")
	}

	// Log the cell change
	// TODO: Add proper logging when logger is implemented
	// Debug message: "Cell on (x, y) has been changed to type on map_name"

	// Notify through hooks system
	if m.parser != nil && m.parser.hookManager != nil {
		m.parser.hookManager.CallHook("map.cell_change", map[string]interface{}{
			"x":        x,
			"y":        y,
			"type":     cellType,
			"map_name": mapName,
		})
	}

	return nil
}

// handleNoTeleport handles the no_teleport packet
// This packet is sent when a teleport or memo action fails
func (m *AccountManager) handleNoTeleport(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract packet data
	failCode, failCodeOk := args["fail"].(byte)
	if !failCodeOk {
		return errors.New("missing fail code in no_teleport packet")
	}

	var errorMessage string
	var clearAI bool

	// Determine error message based on fail code
	switch failCode {
	case TeleportErrorUnavailableArea:
		errorMessage = "Unavailable Area To Teleport"
		clearAI = true
	case TeleportErrorUnavailableMemo:
		errorMessage = "Unavailable Area To Memo"
		clearAI = false
	default:
		errorMessage = fmt.Sprintf("Unavailable Area To Teleport (fail code %d)", failCode)
		clearAI = false
	}

	// Log the error
	// TODO: Add proper logging when logger is implemented

	// Notify through hooks system
	if m.parser != nil && m.parser.hookManager != nil {
		// Send teleport error message
		m.parser.hookManager.CallHook("teleport.error", map[string]interface{}{
			"message": errorMessage,
			"code":    failCode,
		})

		// Clear AI teleport action if needed
		if clearAI {
			m.parser.hookManager.CallHook("ai.clear", map[string]interface{}{
				"action": "teleport",
			})
		}
	}

	return nil
}

// handleWarpPortalList handles the warp_portal_list packet
// This packet is sent when a warp portal list is displayed
func (m *AccountManager) handleWarpPortalList(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract packet data
	warpType, warpTypeOk := args["type"].(byte)
	if !warpTypeOk {
		return errors.New("missing type in warp_portal_list packet")
	}

	// Strip .gat extension from memo map names and create memo list
	memoList := []string{}

	// Process memo1
	if memo1, ok := args["memo1"].(string); ok && memo1 != "" {
		// Strip .gat extension
		if gatIndex := strings.Index(memo1, ".gat"); gatIndex != -1 {
			memo1 = memo1[:gatIndex]
		}
		memoList = append(memoList, memo1)
	}

	// Process memo2
	if memo2, ok := args["memo2"].(string); ok && memo2 != "" {
		// Strip .gat extension
		if gatIndex := strings.Index(memo2, ".gat"); gatIndex != -1 {
			memo2 = memo2[:gatIndex]
		}
		memoList = append(memoList, memo2)
	}

	// Process memo3
	if memo3, ok := args["memo3"].(string); ok && memo3 != "" {
		// Strip .gat extension
		if gatIndex := strings.Index(memo3, ".gat"); gatIndex != -1 {
			memo3 = memo3[:gatIndex]
		}
		memoList = append(memoList, memo3)
	}

	// Process memo4
	if memo4, ok := args["memo4"].(string); ok && memo4 != "" {
		// Strip .gat extension
		if gatIndex := strings.Index(memo4, ".gat"); gatIndex != -1 {
			memo4 = memo4[:gatIndex]
		}
		memoList = append(memoList, memo4)
	}

	// Auto-detect saveMap based on warp type
	if m.parser != nil && m.parser.hookManager != nil {
		if warpType == 26 {
			// For teleport skill, use memo2 as saveMap
			if len(memoList) >= 2 {
				m.parser.hookManager.CallHook("config.update", map[string]interface{}{
					"key":   "saveMap",
					"value": memoList[1],
				})
			}
		} else if warpType == 27 {
			// For butterfly wing, use memo1 as saveMap
			if len(memoList) >= 1 {
				m.parser.hookManager.CallHook("config.update", map[string]interface{}{
					"key":   "saveMap",
					"value": memoList[0],
				})
			}
		}
	}

	// Notify through hooks system
	if m.parser != nil && m.parser.hookManager != nil {
		m.parser.hookManager.CallHook("warp.portal_list", map[string]interface{}{
			"type":      warpType,
			"memo_list": memoList,
		})
	}

	// If teleport skill is used and in teleport AI queue, send warp teleport request
	// This would be handled by the AI system in the original code
	// For now, we'll just log it
	// TODO: Implement AI integration when AI system is implemented

	return nil
}

// handleMapProperty3 handles the map_property3 packet (099B)
// This packet is sent when entering a map with special properties (PvP, GvG, etc.)
func (m *AccountManager) handleMapProperty3(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract packet data
	mapType, mapTypeOk := args["type"].(byte)
	if !mapTypeOk {
		return errors.New("missing map type in map_property3 packet")
	}

	// Process info_table if available
	if infoTable, ok := args["info_table"].([]byte); ok && len(infoTable) >= 4 {
		// In the original Perl code, this would set character statuses based on the info_table
		// For now, we'll just log that we received the info_table
		// TODO: Implement character status updates when character status system is implemented
	}

	// Set PvP mode based on map type
	var pvpMode int
	switch mapType {
	case MapTypePvP:
		pvpMode = PvPModePvP
	case MapTypeGvG:
		pvpMode = PvPModeGvG
	case MapTypeBattleground:
		pvpMode = PvPModeBattleground
	default:
		pvpMode = PvPModeNone
	}

	// Call PvP mode hook if applicable
	if pvpMode > 0 && m.parser != nil && m.parser.hookManager != nil {
		m.parser.hookManager.CallHook("pvp_mode", map[string]interface{}{
			"pvp": pvpMode,
		})
	}

	return nil
}
