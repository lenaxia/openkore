package login

import (
	"sync"
	"testing"
	"time"
)

// TestRegisterTestHook tests the RegisterTestHook method
func TestRegisterTestHook(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Create a variable to track if the hook was called
	var hookCalled bool
	var hookData interface{}
	var mu sync.Mutex

	// Register a test hook
	loginManager.RegisterTestHook(TestHookBeforeLogin, func(hookType TestHookType, data interface{}) {
		mu.Lock()
		defer mu.Unlock()
		hookCalled = true
		hookData = data
	})

	// Trigger the test hook
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeLogin), "test data")

	// Wait a short time for the hook to be called
	time.Sleep(10 * time.Millisecond)

	// Verify that the hook was called
	mu.Lock()
	defer mu.Unlock()
	if !hookCalled {
		t.Error("Expected test hook to be called")
	}
	if hookData != "test data" {
		t.Errorf("Expected hook data to be 'test data', got %v", hookData)
	}
}

// TestUnregisterTestHook tests the UnregisterTestHook method
func TestUnregisterTestHook(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Create a variable to track if the hook was called
	var hookCalled bool
	var mu sync.Mutex

	// Create a hook callback
	callback := func(hookType TestHookType, data interface{}) {
		mu.Lock()
		defer mu.Unlock()
		hookCalled = true
	}

	// Register a test hook
	loginManager.RegisterTestHook(TestHookBeforeLogin, callback)

	// Unregister the test hook
	loginManager.UnregisterTestHook(TestHookBeforeLogin, callback)

	// Trigger the test hook
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeLogin), nil)

	// Wait a short time for the hook to be called
	time.Sleep(10 * time.Millisecond)

	// Verify that the hook was not called
	mu.Lock()
	defer mu.Unlock()
	if hookCalled {
		t.Error("Expected test hook not to be called after unregistering")
	}
}

// TestWithTestHooks tests the WithTestHooks method
func TestWithTestHooks(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Create variables to track if the hooks were called
	var beforeLoginCalled, afterLoginCalled bool
	var mu sync.Mutex

	// Register multiple test hooks
	loginManager.WithTestHooks(map[TestHookType]TestHookCallback{
		TestHookBeforeLogin: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			beforeLoginCalled = true
		},
		TestHookAfterLogin: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			afterLoginCalled = true
		},
	})

	// Trigger the test hooks
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeLogin), nil)
	mockNetworkManager.hookManager.CallHook(string(TestHookAfterLogin), nil)

	// Wait a short time for the hooks to be called
	time.Sleep(10 * time.Millisecond)

	// Verify that the hooks were called
	mu.Lock()
	defer mu.Unlock()
	if !beforeLoginCalled {
		t.Error("Expected before login hook to be called")
	}
	if !afterLoginCalled {
		t.Error("Expected after login hook to be called")
	}
}

// TestIntegrationWithLoginProcess tests that the test hooks integrate with the login process
func TestIntegrationWithLoginProcess(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Create variables to track the login process
	var (
		beforeLoginCalled   bool
		beforeSendCalled    bool
		afterSendCalled     bool
		beforeReceiveCalled bool
		afterReceiveCalled  bool
		stateChangeCalled   bool
		afterLoginCalled    bool
		packetName          string
		stateChangeOldState int
		stateChangeNewState int
		mu                  sync.Mutex
	)

	// Register test hooks
	loginManager.WithTestHooks(map[TestHookType]TestHookCallback{
		TestHookBeforeLogin: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			beforeLoginCalled = true
		},
		TestHookBeforeSendPacket: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			beforeSendCalled = true
			if args, ok := data.(map[string]interface{}); ok {
				if name, ok := args["name"].(string); ok {
					packetName = name
				}
			}
		},
		TestHookAfterSendPacket: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			afterSendCalled = true
		},
		TestHookBeforeReceivePacket: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			beforeReceiveCalled = true
		},
		TestHookAfterReceivePacket: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			afterReceiveCalled = true
		},
		TestHookStateChange: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			stateChangeCalled = true
			if args, ok := data.(map[string]interface{}); ok {
				if oldState, ok := args["oldState"].(int); ok {
					stateChangeOldState = oldState
				}
				if newState, ok := args["newState"].(int); ok {
					stateChangeNewState = newState
				}
			}
		},
		TestHookAfterLogin: func(hookType TestHookType, data interface{}) {
			mu.Lock()
			defer mu.Unlock()
			afterLoginCalled = true
		},
	})

	// Simulate the login process
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeLogin), nil)
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeSendPacket), map[string]interface{}{
		"name": "master_login",
	})
	mockNetworkManager.hookManager.CallHook(string(TestHookAfterSendPacket), nil)
	mockNetworkManager.hookManager.CallHook(string(TestHookBeforeReceivePacket), nil)
	mockNetworkManager.hookManager.CallHook(string(TestHookAfterReceivePacket), nil)
	mockNetworkManager.hookManager.CallHook(string(TestHookStateChange), map[string]interface{}{
		"oldState": 0,
		"newState": 1,
	})
	mockNetworkManager.hookManager.CallHook(string(TestHookAfterLogin), nil)

	// Wait a short time for the hooks to be called
	time.Sleep(10 * time.Millisecond)

	// Verify that the hooks were called
	mu.Lock()
	defer mu.Unlock()
	if !beforeLoginCalled {
		t.Error("Expected before login hook to be called")
	}
	if !beforeSendCalled {
		t.Error("Expected before send hook to be called")
	}
	if packetName != "master_login" {
		t.Errorf("Expected packet name to be 'master_login', got %s", packetName)
	}
	if !afterSendCalled {
		t.Error("Expected after send hook to be called")
	}
	if !beforeReceiveCalled {
		t.Error("Expected before receive hook to be called")
	}
	if !afterReceiveCalled {
		t.Error("Expected after receive hook to be called")
	}
	if !stateChangeCalled {
		t.Error("Expected state change hook to be called")
	}
	if stateChangeOldState != 0 || stateChangeNewState != 1 {
		t.Errorf("Expected state change from 0 to 1, got %d to %d", stateChangeOldState, stateChangeNewState)
	}
	if !afterLoginCalled {
		t.Error("Expected after login hook to be called")
	}
}
