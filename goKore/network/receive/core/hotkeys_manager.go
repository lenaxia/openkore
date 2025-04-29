package core

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// HotkeysManager manages hotkey-related functionality
type HotkeysManager struct {
	parser      *CoreParser
	hookManager *hooks.HookManager
	hotkeyMgr   *HotkeyManager
}

// NewHotkeysManager creates a new hotkeys manager
func NewHotkeysManager(parser *CoreParser, hookManager *hooks.HookManager) *HotkeysManager {
	return &HotkeysManager{
		parser:      parser,
		hookManager: hookManager,
		hotkeyMgr:   NewHotkeyManager(hookManager),
	}
}

// RegisterHandlers registers hotkey-related packet handlers
func (m *HotkeysManager) RegisterHandlers() {
	// Register handler for hotkeys
	m.parser.RegisterHandlerFunc("07D9", "hotkeys", "B a*",
		[]string{"rotate", "hotkeys"},
		func(args map[string]interface{}) error {
			return m.hotkeyMgr.HandleHotkeys(args)
		})
}

// GetHotkeyManager returns the hotkey manager
func (m *HotkeysManager) GetHotkeyManager() *HotkeyManager {
	return m.hotkeyMgr
}
