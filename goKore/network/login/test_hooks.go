package login

import (
	"log"
)

// TestHookType defines the type of test hook
type TestHookType string

const (
	// TestHookBeforeLogin is triggered before the login process starts
	TestHookBeforeLogin TestHookType = "test_before_login"

	// TestHookAfterLogin is triggered after the login process completes
	TestHookAfterLogin TestHookType = "test_after_login"

	// TestHookBeforeSendPacket is triggered before sending a packet
	TestHookBeforeSendPacket TestHookType = "test_before_send_packet"

	// TestHookAfterSendPacket is triggered after sending a packet
	TestHookAfterSendPacket TestHookType = "test_after_send_packet"

	// TestHookBeforeReceivePacket is triggered before receiving a packet
	TestHookBeforeReceivePacket TestHookType = "test_before_receive_packet"

	// TestHookAfterReceivePacket is triggered after receiving a packet
	TestHookAfterReceivePacket TestHookType = "test_after_receive_packet"

	// TestHookStateChange is triggered when the login state changes
	TestHookStateChange TestHookType = "test_state_change"
)

// TestHookCallback is the callback function type for test hooks
type TestHookCallback func(hookType TestHookType, data interface{})

// RegisterTestHook registers a test hook with the login manager
func (lm *LoginManager) RegisterTestHook(hookType TestHookType, callback TestHookCallback) {
	// Get the hook manager
	hookManager, ok := lm.networkManager.GetHookManager().(HookManager)
	if !ok {
		log.Println("Failed to get hook manager for test hook")
		return
	}

	// Register the test hook
	hookManager.Register(string(hookType), func(hookName string, arg interface{}, userData interface{}) {
		callback(hookType, arg)
	})
}

// UnregisterTestHook unregisters a test hook from the login manager
func (lm *LoginManager) UnregisterTestHook(hookType TestHookType, callback TestHookCallback) {
	// Get the hook manager
	hookManager, ok := lm.networkManager.GetHookManager().(HookManager)
	if !ok {
		log.Println("Failed to get hook manager for test hook")
		return
	}

	// Unregister the test hook
	hookManager.Unregister(string(hookType), func(hookName string, arg interface{}, userData interface{}) {
		callback(hookType, arg)
	})
}

// TriggerTestHook triggers a test hook with the given data
func (lm *LoginManager) TriggerTestHook(hookType TestHookType, data interface{}) {
	// Get the hook manager
	hookManager, ok := lm.networkManager.GetHookManager().(HookManager)
	if !ok {
		log.Println("Failed to get hook manager for test hook")
		return
	}

	// Trigger the test hook by calling the hook
	// Try to use the CallHook method if available
	if mockHM, ok := hookManager.(interface{ CallHook(string, interface{}) }); ok {
		mockHM.CallHook(string(hookType), data)
	} else {
		log.Println("Failed to call test hook: hook manager doesn't support CallHook")
	}
}

// WithTestHooks is a convenience function for registering multiple test hooks
func (lm *LoginManager) WithTestHooks(hooks map[TestHookType]TestHookCallback) *LoginManager {
	for hookType, callback := range hooks {
		lm.RegisterTestHook(hookType, callback)
	}
	return lm
}
