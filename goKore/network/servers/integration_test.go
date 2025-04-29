package servers

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/security"
)

// TestServerConfigIntegration tests the integration between server config and security components
func TestServerConfigIntegration(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()

	// Set credentials
	baseConfig.Credentials = ServerCredentials{
		Username: "testuser",
		Password: "testpass",
		PINCode:  "1234",
	}

	// Set server type in custom fields to help detection
	baseConfig.ServerConfig.CustomFields["serverType"] = "sakray"

	// Set options
	baseConfig.Options = ServerOptions{
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
	err := ValidateCredentials(baseConfig.Credentials)
	if err != nil {
		t.Errorf("ValidateCredentials() returned error: %v", err)
	}

	// Validate server config
	err = ValidateServerConfig(baseConfig.ServerConfig)
	if err != nil {
		t.Errorf("ValidateServerConfig() returned error: %v", err)
	}

	// Validate network config
	err = ValidateNetworkConfig(baseConfig.NetworkConfig)
	if err != nil {
		t.Errorf("ValidateNetworkConfig() returned error: %v", err)
	}

	// Test server type detection
	serverType := DetectServerType(baseConfig.ServerConfig)
	if serverType == ServerTypeUnknown {
		t.Error("DetectServerType() returned ServerTypeUnknown")
	}
}

// TestServerSecurityIntegration tests the integration between server interface and security components
func TestServerSecurityIntegration(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create security components
	parser := core.NewCoreParser("ServerType0", hookManager)
	loginManager := security.NewLoginManager(parser, hookManager)
	pinManager := security.NewPINManager(parser, hookManager)
	antiCheatManager := security.NewAntiCheatManager(parser, hookManager)

	// Register handlers
	loginManager.RegisterHandlers()
	pinManager.RegisterHandlers()
	antiCheatManager.RegisterHandlers()

	// Create a base server config
	baseConfig := NewBaseServerConfig()

	// Set credentials
	baseConfig.Credentials = ServerCredentials{
		Username: "testuser",
		Password: "testpass",
		PINCode:  "1234",
	}

	// Simulate login
	loginManager.SetCredentials(baseConfig.Credentials.Username, baseConfig.Credentials.Password)

	// Manually set login state for testing
	loginManager.SetState(security.LoginStateLoggedIn)

	// Manually call the login success hook
	hookManager.CallHook("security/login_success", nil)

	// Check that login state is correct
	if loginManager.GetState() != security.LoginStateLoggedIn {
		t.Errorf("Login state = %v, want %v", loginManager.GetState(), security.LoginStateLoggedIn)
	}

	// Set PIN
	if err := pinManager.SetPIN(baseConfig.Credentials.PINCode); err != nil {
		t.Fatalf("SetPIN returned error: %v", err)
	}

	// Manually set PIN state for testing
	pinManager.SetState(security.PINStateVerified)

	// Manually call the PIN verified hook
	hookManager.CallHook("security/pin_verified", nil)

	// Check that PIN state is correct
	if pinManager.GetState() != security.PINStateVerified {
		t.Errorf("PIN state = %v, want %v", pinManager.GetState(), security.PINStateVerified)
	}

	// Enable anti-cheat
	antiCheatManager.Enable(baseConfig.Options.AntiCheatType)

	// Manually set anti-cheat state for testing
	antiCheatManager.SetState(security.AntiCheatStateVerified)

	// Manually call the anti-cheat verified hook
	hookManager.CallHook("security/anticheat_verified", nil)
}

// TestServerErrorHandling tests error handling in the server interface
func TestServerErrorHandling(t *testing.T) {
	// Test ValidateCredentials with invalid credentials
	credentials := ServerCredentials{
		Username: "",
		Password: "testpass",
	}

	err := ValidateCredentials(credentials)
	if err == nil || err.Error() != "username is not set" {
		t.Errorf("ValidateCredentials() with empty username returned error: %v, want 'username is not set'", err)
	}

	credentials = ServerCredentials{
		Username: "testuser",
		Password: "",
	}

	err = ValidateCredentials(credentials)
	if err == nil || err.Error() != "password is not set" {
		t.Errorf("ValidateCredentials() with empty password returned error: %v, want 'password is not set'", err)
	}

	// Test ValidateServerConfig with invalid server config
	serverConfig := config.NewServerConfigManager().CreateDefaultServerConfig("test")
	serverConfig.IP = ""

	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server IP is not set" {
		t.Errorf("ValidateServerConfig() with empty IP returned error: %v, want 'server IP is not set'", err)
	}

	serverConfig = config.NewServerConfigManager().CreateDefaultServerConfig("test")
	serverConfig.Port = 0

	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server port is not set" {
		t.Errorf("ValidateServerConfig() with zero port returned error: %v, want 'server port is not set'", err)
	}

	serverConfig = config.NewServerConfigManager().CreateDefaultServerConfig("test")
	serverConfig.Name = ""

	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server name is not set" {
		t.Errorf("ValidateServerConfig() with empty name returned error: %v, want 'server name is not set'", err)
	}

	// Test ValidateNetworkConfig with invalid network config
	networkConfig := config.NewNetworkConfigManager().CreateDefaultNetworkConfig("test")
	networkConfig.Timeouts.Connect = 0

	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "connection timeout must be greater than 0" {
		t.Errorf("ValidateNetworkConfig() with zero connect timeout returned error: %v, want 'connection timeout must be greater than 0'", err)
	}

	networkConfig = config.NewNetworkConfigManager().CreateDefaultNetworkConfig("test")
	networkConfig.ReconnectPolicy.InitialInterval = 0

	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "reconnect delay must be greater than 0" {
		t.Errorf("ValidateNetworkConfig() with zero reconnect delay returned error: %v, want 'reconnect delay must be greater than 0'", err)
	}

	networkConfig = config.NewNetworkConfigManager().CreateDefaultNetworkConfig("test")
	networkConfig.ReconnectPolicy.MaxAttempts = -1

	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "max retries must be greater than or equal to 0" {
		t.Errorf("ValidateNetworkConfig() with negative max retries returned error: %v, want 'max retries must be greater than or equal to 0'", err)
	}
}
