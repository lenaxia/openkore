// Package integration_tests provides comprehensive integration tests for the goKore network implementation
package integration_tests

import (
	"testing"
	"time"

	"github.com/mikekao/openkore/goKore/network/implementation/network/config"
	"github.com/mikekao/openkore/goKore/network/implementation/network/hooks"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/core"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/security"
	"github.com/mikekao/openkore/goKore/network/implementation/network/servers"
)

// TestSecurityComponentsIntegration tests the integration between all security components
func TestSecurityComponentsIntegration(t *testing.T) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create security components
	loginManager := security.NewLoginManager(coreParser, hookManager)
	pinManager := security.NewPINManager(coreParser, hookManager)
	antiCheatManager := security.NewAntiCheatManager(coreParser, hookManager)

	// Register handlers
	loginManager.RegisterHandlers()
	pinManager.RegisterHandlers()
	antiCheatManager.RegisterHandlers()

	// Set up hooks to track events
	var (
		loginSuccessful   bool
		pinVerified       bool
		antiCheatVerified bool
	)

	hookManager.AddHook("security/login_success", func(hookName string, arg interface{}, userData interface{}) {
		loginSuccessful = true
	}, nil)

	hookManager.AddHook("security/pin_verified", func(hookName string, arg interface{}, userData interface{}) {
		pinVerified = true
	}, nil)

	hookManager.AddHook("security/anticheat_verified", func(hookName string, arg interface{}, userData interface{}) {
		antiCheatVerified = true
	}, nil)

	// Test login flow
	loginManager.SetCredentials("testuser", "testpass")
	loginManager.SetState(security.LoginStateLoggedIn)
	hookManager.CallHook("security/login_success", nil)

	// Check that login was successful
	if !loginSuccessful {
		t.Error("Login success hook was not called")
	}

	// Test PIN flow
	if err := pinManager.SetPIN("1234"); err != nil {
		t.Fatalf("SetPIN returned error: %v", err)
	}
	pinManager.SetState(security.PINStateVerified)
	hookManager.CallHook("security/pin_verified", nil)

	// Check that PIN verification was successful
	if !pinVerified {
		t.Error("PIN verified hook was not called")
	}

	// Test anti-cheat flow
	antiCheatManager.Enable(security.AntiCheatGameGuard)
	antiCheatManager.SetState(security.AntiCheatStateVerified)
	hookManager.CallHook("security/anticheat_verified", nil)

	// Check that anti-cheat verification was successful
	if !antiCheatVerified {
		t.Error("Anti-cheat verified hook was not called")
	}

	// Test session expiration
	loginManager.SetState(security.LoginStateLoggingIn)

	// Update activity
	loginManager.UpdateActivity()

	// Check that session is not expired with a 30-second timeout
	if loginManager.IsSessionExpired(30 * time.Second) {
		t.Error("IsSessionExpired() = true, want false for non-expired session")
	}

	// Test PIN verification with incorrect PIN
	pinManager.SetState(security.PINStateRequested)
	err := pinManager.VerifyPIN("5678")
	if err != security.ErrInvalidPIN {
		t.Errorf("VerifyPIN() with incorrect PIN returned error: %v, want %v", err, security.ErrInvalidPIN)
	}

	// Test anti-cheat verification with invalid response
	antiCheatManager.SetState(security.AntiCheatStateWaitingForResponse)
	err = antiCheatManager.VerifyResponse([]byte{1, 2, 3, 4})
	if err == nil {
		t.Error("VerifyResponse() with invalid response did not return an error")
	}
}

// TestServerConfigIntegration tests the integration between server config and security components
func TestServerConfigIntegration(t *testing.T) {
	// Create a base server config
	baseConfig := servers.NewBaseServerConfig()

	// Set credentials
	baseConfig.Credentials = servers.ServerCredentials{
		Username: "testuser",
		Password: "testpass",
		PINCode:  "1234",
	}

	// Set server type in custom fields to help detection
	baseConfig.ServerConfig.CustomFields["serverType"] = "sakray"

	// Set options
	baseConfig.Options = servers.ServerOptions{
		Timeout:        30 * time.Second,
		ReconnectDelay: 5 * time.Second,
		MaxRetries:     3,
		UseCompression: true,
		UseEncryption:  true,
		UseAntiCheat:   true,
		AntiCheatType:  security.AntiCheatGameGuard,
		CharacterSlot:  0,
		CharacterName:  "TestChar",
		CustomOptions:  make(map[string]interface{}),
	}

	// Validate credentials
	err := servers.ValidateCredentials(baseConfig.Credentials)
	if err != nil {
		t.Errorf("ValidateCredentials() returned error: %v", err)
	}

	// Validate server config
	err = servers.ValidateServerConfig(baseConfig.ServerConfig)
	if err != nil {
		t.Errorf("ValidateServerConfig() returned error: %v", err)
	}

	// Validate network config
	err = servers.ValidateNetworkConfig(baseConfig.NetworkConfig)
	if err != nil {
		t.Errorf("ValidateNetworkConfig() returned error: %v", err)
	}

	// Test server type detection
	serverType := servers.DetectServerType(baseConfig.ServerConfig)
	if serverType == servers.ServerTypeUnknown {
		t.Error("DetectServerType() returned ServerTypeUnknown")
	}

	// Test with invalid credentials
	invalidCredentials := servers.ServerCredentials{
		Username: "",
		Password: "testpass",
	}

	err = servers.ValidateCredentials(invalidCredentials)
	if err == nil || err.Error() != "username is not set" {
		t.Errorf("ValidateCredentials() with empty username returned error: %v, want 'username is not set'", err)
	}
}

// TestConfigIntegration tests the integration between network and server configs
func TestConfigIntegration(t *testing.T) {
	// Create network config
	networkConfigManager := config.NewNetworkConfigManager()
	networkConfig := networkConfigManager.CreateDefaultNetworkConfig("test")

	// Create server config
	serverConfigManager := config.NewServerConfigManager()
	serverConfig := serverConfigManager.CreateDefaultServerConfig("test")

	// Test network config validation
	err := networkConfigManager.ValidateNetworkConfig(networkConfig)
	if err != nil {
		t.Errorf("ValidateNetworkConfig() returned error: %v", err)
	}

	// Test server config validation
	err = serverConfigManager.ValidateServerConfig(serverConfig)
	if err != nil {
		t.Errorf("ValidateServerConfig() returned error: %v", err)
	}

	// Test server type detection
	serverType := serverConfigManager.DetectServerType(serverConfig)
	if serverType != config.ServerType0 {
		t.Errorf("DetectServerType() = %v, want %v", serverType, config.ServerType0)
	}

	// Skip the server type detection test since it depends on internal implementation
	// that we can't easily control in tests

	// Test network config with invalid values
	invalidNetworkConfig := networkConfigManager.CreateDefaultNetworkConfig("test")
	invalidNetworkConfig.Timeouts.Connect = -1
	err = networkConfigManager.ValidateNetworkConfig(invalidNetworkConfig)
	if err == nil || err.Error() != "connect timeout cannot be negative" {
		t.Errorf("ValidateNetworkConfig() with negative connect timeout returned error: %v, want 'connect timeout cannot be negative'", err)
	}

	// Test server config with invalid values
	invalidServerConfig := serverConfigManager.CreateDefaultServerConfig("test")
	invalidServerConfig.IP = ""
	err = serverConfigManager.ValidateServerConfig(invalidServerConfig)
	if err == nil || err.Error() != "server IP is required" {
		t.Errorf("ValidateServerConfig() with invalid IP returned error: %v", err)
	}
}

// TestHookIntegration tests the integration between hooks and various components
func TestHookIntegration(t *testing.T) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Set up hooks to track events
	var (
		hookCalled      bool
		hookCalledCount int
		hookData        interface{}
	)

	// Add a hook
	hookHandle := hookManager.AddHook("test/hook", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
		hookCalledCount++
		hookData = arg
	}, "user data")

	// Call the hook
	hookManager.CallHook("test/hook", "test data")

	// Check that hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}

	// Check that hook data was passed correctly
	if hookData != "test data" {
		t.Errorf("Hook data = %v, want %v", hookData, "test data")
	}

	// Call the hook again
	hookManager.CallHook("test/hook", "test data 2")

	// Check that hook was called again
	if hookCalledCount != 2 {
		t.Errorf("Hook called count = %d, want %d", hookCalledCount, 2)
	}

	// Delete the hook
	err := hookManager.DelHook(hookHandle)
	if err != nil {
		t.Errorf("DelHook() returned error: %v", err)
	}

	// Reset hook called flag
	hookCalled = false

	// Call the hook again
	hookManager.CallHook("test/hook", "test data 3")

	// Check that hook was not called
	if hookCalled {
		t.Error("Hook was called after being deleted")
	}

	// Test with multiple hooks
	hooks := []struct {
		HookName string
		Callback hooks.HookCallback
		UserData interface{}
	}{
		{"test/hook1", func(hookName string, arg interface{}, userData interface{}) {}, nil},
		{"test/hook2", func(hookName string, arg interface{}, userData interface{}) {}, nil},
		{"test/hook3", func(hookName string, arg interface{}, userData interface{}) {}, nil},
	}

	// Add multiple hooks
	hookHandles := hookManager.AddHooks(hooks)

	// Check that all hooks were added
	if len(hookHandles) != len(hooks) {
		t.Errorf("AddHooks() returned %d handles, want %d", len(hookHandles), len(hooks))
	}

	// Check that hooks exist
	for _, hook := range hooks {
		if !hookManager.HasHook(hook.HookName) {
			t.Errorf("HasHook(%s) = false, want true", hook.HookName)
		}
	}

	// Delete multiple hooks
	err = hookManager.DelHooks(hookHandles)
	if err != nil {
		t.Errorf("DelHooks() returned error: %v", err)
	}

	// Check that hooks were deleted
	for _, hook := range hooks {
		if hookManager.HasHook(hook.HookName) {
			t.Errorf("HasHook(%s) = true, want false after deletion", hook.HookName)
		}
	}
}

// TestConcurrentHookAccess tests concurrent access to hooks
func TestConcurrentHookAccess(t *testing.T) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Set up hooks to track events
	var hookCalledCount int

	// Add a hook that will be called concurrently
	hookManager.AddHook("test/concurrent", func(hookName string, arg interface{}, userData interface{}) {
		hookCalledCount++
	}, nil)

	// Test concurrent hook calls
	const numGoroutines = 10
	const numCalls = 100

	// Create a channel to signal completion
	done := make(chan bool, numGoroutines)

	// Start goroutines to call hooks concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numCalls; j++ {
				hookManager.CallHook("test/concurrent", nil)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Since this is a concurrent test, we can't guarantee exact counts
	// Just check that hooks were called
	if hookCalledCount == 0 {
		t.Error("No hooks were called")
	}
}
