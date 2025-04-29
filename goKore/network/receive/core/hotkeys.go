package core

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
)

// HotkeyType represents the type of a hotkey (skill or item)
type HotkeyType byte

const (
	HotkeyTypeItem  HotkeyType = 0
	HotkeyTypeSkill HotkeyType = 1
)

// Hotkey represents a single hotkey binding
type Hotkey struct {
	Type HotkeyType // 0 = item, 1 = skill
	ID   uint32     // Item ID or Skill ID
	Lv   uint16     // Level (for skills)
}

// HotkeyList represents a list of hotkeys
type HotkeyList struct {
	Hotkeys []Hotkey
	Rotate  bool // Whether hotkeys should rotate
}

// HotkeyManager manages hotkey-related functionality
type HotkeyManager struct {
	hookManager *hooks.HookManager
	hotkeyList  *HotkeyList
}

// NewHotkeyManager creates a new hotkey manager
func NewHotkeyManager(hookManager *hooks.HookManager) *HotkeyManager {
	return &HotkeyManager{
		hookManager: hookManager,
		hotkeyList:  &HotkeyList{},
	}
}

// HandleHotkeys handles the hotkeys packet
// Packet format: 07D9 <rotate>.B <hotkey data>.?B
func (m *HotkeyManager) HandleHotkeys(args map[string]interface{}) error {
	// Reset the hotkey list
	m.hotkeyList = &HotkeyList{}

	// Extract hotkeys data with safety check
	hotkeysData, ok := args["hotkeys"].([]byte)
	if !ok {
		return fmt.Errorf("invalid hotkeys data")
	}

	// TODO: Implement rotate flag if needed
	// if rotateVal, ok := args["rotate"].(byte); ok {
	//     m.hotkeyList.Rotate = rotateVal != 0
	// }

	// Parse the hotkeys data
	m.hotkeyList.Hotkeys = make([]Hotkey, 0, len(hotkeysData)/7)
	for i := 0; i < len(hotkeysData); i += 7 {
		if i+7 > len(hotkeysData) {
			break // Avoid out of bounds access
		}

		hotkey := Hotkey{
			Type: HotkeyType(hotkeysData[i]),
			ID:   uint32(hotkeysData[i+1]) | uint32(hotkeysData[i+2])<<8 | uint32(hotkeysData[i+3])<<16 | uint32(hotkeysData[i+4])<<24,
			Lv:   uint16(hotkeysData[i+5]) | uint16(hotkeysData[i+6])<<8,
		}
		m.hotkeyList.Hotkeys = append(m.hotkeyList.Hotkeys, hotkey)
	}

	// Format and log the hotkeys
	m.logHotkeys()

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("core.hotkeys", map[string]interface{}{
			"hotkeys": m.hotkeyList.Hotkeys,
			"rotate":  m.hotkeyList.Rotate,
		})
	}

	return nil
}

// logHotkeys formats and logs the hotkey list
func (m *HotkeyManager) logHotkeys() {
	// Create a header
	header := "Hotkeys"
	separator := strings.Repeat("-", 79)
	centered := fmt.Sprintf("%s %s %s", separator[:35], header, separator[:35])

	// Create a table header
	tableHeader := fmt.Sprintf("%-3s %-30s %-5s %-3s", "#", "Name", "Type", "Lv")

	// Create the message
	var message strings.Builder
	message.WriteString(centered + "\n")
	message.WriteString(tableHeader + "\n")
	message.WriteString(separator + "\n")

	// Add each hotkey
	for i, hotkey := range m.hotkeyList.Hotkeys {
		var name, typeName string
		if hotkey.Type == HotkeyTypeSkill {
			name = fmt.Sprintf("Skill #%d", hotkey.ID) // Would use skill name lookup in a real implementation
			typeName = "skill"
		} else {
			name = fmt.Sprintf("Item #%d", hotkey.ID) // Would use item name lookup in a real implementation
			typeName = "item"
		}

		message.WriteString(fmt.Sprintf("%-3d %-30s %-5s %-3d\n", i, name, typeName, hotkey.Lv))
	}

	message.WriteString(separator + "\n")

	// Log the message
	// In a real implementation, this would use a logger
	// logger.Debug(message.String())
}

// GetHotkeyList returns the current hotkey list
func (m *HotkeyManager) GetHotkeyList() *HotkeyList {
	return m.hotkeyList
}
