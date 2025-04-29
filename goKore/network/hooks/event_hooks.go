// Package hooks provides functionality for event hooks in the OpenKore network stack.
// This file implements the event hook system for plugins to interact with the network stack.
package hooks

import (
	"fmt"
	"sync"
)

// HookCallback is a function that is called when a hook is triggered
type HookCallback func(hookName string, arg interface{}, userData interface{})

// HookEntry represents a registered hook callback
type HookEntry struct {
	Callback HookCallback
	UserData interface{}
}

// HookHandle is a handle to a registered hook that can be used to remove it
type HookHandle struct {
	HookName string
	Index    int
}

// HookManager manages event hooks for the network stack
// It implements the IHookManager interface
type HookManager struct {
	hooks     map[string][]HookEntry
	mutex     sync.RWMutex
	nextIndex map[string]int
}

// NewHookManager creates a new hook manager
func NewHookManager() *HookManager {
	return &HookManager{
		hooks:     make(map[string][]HookEntry),
		nextIndex: make(map[string]int),
	}
}

// AddHook adds a hook for the given hook name
// Returns a handle that can be used to remove the hook
func (m *HookManager) AddHook(hookName string, callback HookCallback, userData interface{}) *HookHandle {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create the hook list if it doesn't exist
	if _, ok := m.hooks[hookName]; !ok {
		m.hooks[hookName] = make([]HookEntry, 0)
		m.nextIndex[hookName] = 0
	}

	// Create the hook entry
	entry := HookEntry{
		Callback: callback,
		UserData: userData,
	}

	// Add the hook entry to the list
	index := m.nextIndex[hookName]
	m.nextIndex[hookName]++
	m.hooks[hookName] = append(m.hooks[hookName], entry)

	// Return a handle to the hook
	return &HookHandle{
		HookName: hookName,
		Index:    index,
	}
}

// AddHooks adds multiple hooks at once
// Returns a slice of handles that can be used to remove the hooks
func (m *HookManager) AddHooks(hooks []struct {
	HookName string
	Callback HookCallback
	UserData interface{}
}) []*HookHandle {
	handles := make([]*HookHandle, len(hooks))
	for i, hook := range hooks {
		handles[i] = m.AddHook(hook.HookName, hook.Callback, hook.UserData)
	}
	return handles
}

// DelHook removes a hook
func (m *HookManager) DelHook(handle *HookHandle) error {
	if handle == nil {
		return fmt.Errorf("invalid hook handle: nil")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if the hook exists
	hookList, ok := m.hooks[handle.HookName]
	if !ok {
		return fmt.Errorf("hook not found: %s", handle.HookName)
	}

	// Find the hook entry
	found := false
	for i := range hookList {
		if i == handle.Index {
			// Remove the hook entry
			m.hooks[handle.HookName] = append(hookList[:i], hookList[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("hook entry not found: %s[%d]", handle.HookName, handle.Index)
	}

	// Remove the hook list if it's empty
	if len(m.hooks[handle.HookName]) == 0 {
		delete(m.hooks, handle.HookName)
		delete(m.nextIndex, handle.HookName)
	}

	return nil
}

// DelHooks removes multiple hooks at once
func (m *HookManager) DelHooks(handles []*HookHandle) error {
	var lastErr error
	for _, handle := range handles {
		if err := m.DelHook(handle); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// CallHook calls all callbacks for the given hook name
func (m *HookManager) CallHook(hookName string, arg interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if the hook exists
	hookList, ok := m.hooks[hookName]
	if !ok {
		return
	}

	// Make a copy of the hook list to avoid issues if hooks are added or removed during the call
	hookListCopy := make([]HookEntry, len(hookList))
	copy(hookListCopy, hookList)

	// Call all callbacks
	for _, entry := range hookListCopy {
		entry.Callback(hookName, arg, entry.UserData)
	}
}

// HasHook checks if there are any hooks registered for the given hook name
func (m *HookManager) HasHook(hookName string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.hooks[hookName]
	return ok
}

// GetHookCount returns the number of hooks registered for the given hook name
func (m *HookManager) GetHookCount(hookName string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	hookList, ok := m.hooks[hookName]
	if !ok {
		return 0
	}
	return len(hookList)
}

// GetAllHookNames returns a list of all registered hook names
func (m *HookManager) GetAllHookNames() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.hooks))
	for name := range m.hooks {
		names = append(names, name)
	}
	return names
}

// Clear removes all hooks
func (m *HookManager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.hooks = make(map[string][]HookEntry)
	m.nextIndex = make(map[string]int)
}

// RegisterHook implements the IHookManager interface by wrapping AddHook
func (m *HookManager) RegisterHook(event string, fn func(interface{})) {
	m.AddHook(event, func(hookName string, arg interface{}, userData interface{}) {
		fn(arg)
	}, nil)
}

// NetworkHooks is a singleton instance of HookManager for the network stack
var NetworkHooks = NewHookManager()

// Common hook names used in the network stack
const (
	// Connection hooks
	HookConnecting    = "network.connecting"
	HookConnected     = "network.connected"
	HookDisconnecting = "network.disconnecting"
	HookDisconnected  = "network.disconnected"

	// Data hooks
	HookDataSent     = "network.data_sent"
	HookDataReceived = "network.data_received"

	// Packet hooks
	HookPacketSend     = "network.packet_send"
	HookPacketReceived = "network.packet_received"
	HookPacketParsed   = "network.packet_parsed"

	// Error hooks
	HookError = "network.error"

	// State hooks
	HookStateChange = "network.state_change"
)

// AddHook adds a hook for the given hook name to the global NetworkHooks instance
func AddHook(hookName string, callback HookCallback, userData interface{}) *HookHandle {
	return NetworkHooks.AddHook(hookName, callback, userData)
}

// AddHooks adds multiple hooks at once to the global NetworkHooks instance
func AddHooks(hooks []struct {
	HookName string
	Callback HookCallback
	UserData interface{}
}) []*HookHandle {
	return NetworkHooks.AddHooks(hooks)
}

// DelHook removes a hook from the global NetworkHooks instance
func DelHook(handle *HookHandle) error {
	return NetworkHooks.DelHook(handle)
}

// DelHooks removes multiple hooks at once from the global NetworkHooks instance
func DelHooks(handles []*HookHandle) error {
	return NetworkHooks.DelHooks(handles)
}

// CallHook calls all callbacks for the given hook name in the global NetworkHooks instance
func CallHook(hookName string, arg interface{}) {
	NetworkHooks.CallHook(hookName, arg)
}

// HasHook checks if there are any hooks registered for the given hook name in the global NetworkHooks instance
func HasHook(hookName string) bool {
	return NetworkHooks.HasHook(hookName)
}

// GetHookCount returns the number of hooks registered for the given hook name in the global NetworkHooks instance
func GetHookCount(hookName string) int {
	return NetworkHooks.GetHookCount(hookName)
}

// GetAllHookNames returns a list of all registered hook names in the global NetworkHooks instance
func GetAllHookNames() []string {
	return NetworkHooks.GetAllHookNames()
}

// Clear removes all hooks from the global NetworkHooks instance
func Clear() {
	NetworkHooks.Clear()
}
