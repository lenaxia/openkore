package actor

import (
	"time"
)

// HandleOfflineCloneFound handles the offline_clone_found packet
// Packet format: 0A7B <ID>.L <name>.24B <job>.W <coords>.3B <robe>.W <clothes_color>.W <lowhead>.W <midhead>.W <tophead>.W <weapon>.W <shield>.W <sex>.B <hair_color>.W
func (h *Handler) HandleOfflineCloneFound(args map[string]interface{}) error {
	// Extract ID with safety check
	idVal, ok := args["ID"]
	if !ok {
		return nil // Silently ignore if ID is missing
	}
	id, ok := idVal.([]byte)
	if !ok {
		return nil // Silently ignore if ID is not a byte slice
	}

	// Check if the player already exists
	player := h.playersList.GetByID(id)
	if player == nil {
		// Create a new player
		player = NewPlayer(id)
		player.SetClone(true) // Mark as an offline shop clone

		// Extract name with safety check
		if nameVal, ok := args["name"]; ok {
			if name, ok := nameVal.([]byte); ok {
				player.SetName(bytesToString(name))
			}
		}

		// Extract job ID with safety check
		if jobIDVal, ok := args["jobID"]; ok {
			if jobID, ok := jobIDVal.(uint16); ok {
				player.SetJob(jobID)
			}
		}

		// Extract position with safety check
		var x, y int
		if coordXVal, ok := args["coord_x"]; ok {
			if coordX, ok := coordXVal.(int16); ok {
				x = int(coordX)
			}
		}
		if coordYVal, ok := args["coord_y"]; ok {
			if coordY, ok := coordYVal.(int16); ok {
				y = int(coordY)
			}
		}
		player.SetPosition(&Position{X: x, Y: y})
		player.SetPositionTo(&Position{X: x, Y: y})

		// Extract appearance with safety check
		var robe, clothesColor, lowhead, midhead, tophead, weapon, shield uint16
		if robeVal, ok := args["robe"]; ok {
			if val, ok := robeVal.(uint16); ok {
				robe = val
			}
		}
		if clothesColorVal, ok := args["clothes_color"]; ok {
			if val, ok := clothesColorVal.(uint16); ok {
				clothesColor = val
			}
		}
		if lowheadVal, ok := args["lowhead"]; ok {
			if val, ok := lowheadVal.(uint16); ok {
				lowhead = val
			}
		}
		if midheadVal, ok := args["midhead"]; ok {
			if val, ok := midheadVal.(uint16); ok {
				midhead = val
			}
		}
		if topheadVal, ok := args["tophead"]; ok {
			if val, ok := topheadVal.(uint16); ok {
				tophead = val
			}
		}
		if weaponVal, ok := args["weapon"]; ok {
			if val, ok := weaponVal.(uint16); ok {
				weapon = val
			}
		}
		if shieldVal, ok := args["shield"]; ok {
			if val, ok := shieldVal.(uint16); ok {
				shield = val
			}
		}
		player.SetAppearance(0, clothesColor, clothesColor) // Use clothes_color for hair_style and hair_color
		player.SetHeadgear(tophead, midhead, lowhead)
		player.SetEquipment(weapon, shield, robe)

		// Extract sex with safety check
		if sexVal, ok := args["sex"]; ok {
			if sex, ok := sexVal.(byte); ok {
				player.SetSex(sex)
			}
		}

		// Extract hair color with safety check
		if hairColorVal, ok := args["hair_color"]; ok {
			if hairColor, ok := hairColorVal.(uint16); ok {
				player.SetAppearance(0, hairColor, clothesColor)
			}
		}

		// Add the player to the list
		h.playersList.Add(player)

		// Trigger hooks
		if h.hookManager != nil {
			h.hookManager.CallHook("add_player_list", player)
			h.hookManager.CallHook("player", map[string]interface{}{"player": player})
			h.hookManager.CallHook("player_exist", map[string]interface{}{"player": player})
		}
	}

	return nil
}

// HandleOfflineCloneLost handles the offline_clone_lost packet
// Packet format: 0A7C <ID>.L
func (h *Handler) HandleOfflineCloneLost(args map[string]interface{}) error {
	// Extract ID with safety check
	idVal, ok := args["ID"]
	if !ok {
		return nil // Silently ignore if ID is missing
	}
	id, ok := idVal.([]byte)
	if !ok {
		return nil // Silently ignore if ID is not a byte slice
	}

	// Check if the player exists
	player := h.playersList.GetByID(id)
	if player != nil {
		// Set the gone time
		player.SetGoneTime(time.Now())

		// Store a deep copy in the old players map
		h.playersOld[string(id)] = player.DeepCopy().(*Player)

		// Trigger hook
		if h.hookManager != nil {
			h.hookManager.CallHook("player_disappeared", map[string]interface{}{"player": player})
		}

		// Remove the player from the list
		h.playersList.Remove(player)

		// Try to remove from vender list and buyer list
		// This would be handled by other components in the Go implementation
	}

	return nil
}

// bytesToString converts a byte slice to a string, stopping at the first null byte
func bytesToString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
