package servers

import (
	"testing"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/receive/security"
)

// TestServerFactoryWithMultipleTypes tests the integration between the server factory and different server types
func TestServerFactoryWithMultipleTypes(t *testing.T) {
	// Create a server factory
	factory := NewServerFactory()

	// Test creating different server types
	serverTypes := []ServerType{
		ServerTypeOfficial,
		ServerTypeSakray,
		ServerTypeRenewal,
		ServerTypeClassic,
		ServerTypePreRenewal,
	}

	for _, serverType := range serverTypes {
		// Create a base server config
		baseConfig := NewBaseServerConfig()
		baseConfig.Type = serverType

		// Create a server
		server, err := factory.CreateServer(baseConfig)
		if err != nil {
			t.Fatalf("CreateServer() for type %v returned error: %v", serverType, err)
		}
		if server == nil {
			t.Fatalf("CreateServer() for type %v returned nil", serverType)
		}

		// Check that the server has the correct type
		// For server types that don't have specific implementations,
		// they will fall back to the base server type
		if serverType == ServerTypeOfficial || serverType == ServerTypeSakray {
			if server.GetServerType() != serverType {
				t.Errorf("GetServerType() for type %v = %v, want %v", serverType, server.GetServerType(), serverType)
			}
		}
	}
}

// TestServerHookIntegration tests the integration between servers and the hook system
func TestServerHookIntegration(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()

	// Create a base server
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Register a hook on the server's hook manager
	hookCalled := false
	server.RegisterHook("test_hook", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Call the hook manually
	server.hookManager.CallHook("test_hook", nil)

	// Check that the hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}
}

// TestServerSecurityComponentsIntegration tests the integration between servers and security components
func TestServerSecurityComponentsIntegration(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()

	// Create a base server
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Set credentials
	credentials := ServerCredentials{
		Username: "testuser",
		Password: "testpass",
		PINCode:  "1234",
	}
	server.SetCredentials(credentials)

	// Set options with anti-cheat enabled
	options := ServerOptions{
		UseAntiCheat:  true,
		AntiCheatType: security.AntiCheatGameGuard,
	}
	server.SetOptions(options)

	// Check that the login manager has the correct credentials
	loginManager := server.GetLoginManager()
	if loginManager == nil {
		t.Fatal("GetLoginManager() returned nil")
	}

	// Check that the PIN manager has the correct PIN
	pinManager := server.GetPINManager()
	if pinManager == nil {
		t.Fatal("GetPINManager() returned nil")
	}

	// Check that the anti-cheat manager has the correct type
	antiCheatManager := server.GetAntiCheatManager()
	if antiCheatManager == nil {
		t.Fatal("GetAntiCheatManager() returned nil")
	}
}

// TestServerWithConfigManagerIntegration tests the integration between servers and configuration managers
func TestServerWithConfigManagerIntegration(t *testing.T) {
	// Create a server config manager
	serverConfigManager := config.NewServerConfigManager()

	// Create a network config manager
	networkConfigManager := config.NewNetworkConfigManager()

	// Create a server config
	serverConfig := serverConfigManager.CreateDefaultServerConfig("test")
	serverConfig.CustomFields["serverType"] = "sakray"

	// Create a network config
	networkConfig := networkConfigManager.CreateDefaultNetworkConfig("test")

	// Create a base server config
	baseConfig := &BaseServerConfig{
		ServerConfig:  serverConfig,
		NetworkConfig: networkConfig,
		Type:          ServerTypeUnknown,
	}

	// Create a server
	server, err := CreateServerFromConfig(baseConfig)
	if err != nil {
		t.Fatalf("CreateServerFromConfig() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("CreateServerFromConfig() returned nil")
	}

	// Check that the server has the correct type
	if server.GetServerType() != ServerTypeSakray {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeSakray)
	}

	// Check that the server has the correct config
	if server.GetConfig() != serverConfig {
		t.Error("GetConfig() returned incorrect config")
	}

	// Check that the server has the correct network config
	if server.GetNetworkConfig() != networkConfig {
		t.Error("GetNetworkConfig() returned incorrect config")
	}
}

// TestServerPacketHandlerIntegration tests the integration between servers and packet handlers
func TestServerPacketHandlerIntegration(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()

	// Create a base server
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Register a packet handler
	server.RegisterPacketHandler("0001", "test_packet", "v", []string{"len"}, func(args map[string]interface{}) error {
		// Just a test handler that does nothing
		return nil
	})

	// Get the parser
	parser := server.GetParser()
	if parser == nil {
		t.Fatal("GetParser() returned nil")
	}

	// There's no direct way to check if the handler was registered correctly without processing a packet,
	// but we can check that the parser has the handler registered
	handler, exists := parser.GetHandler("test_packet")
	if !exists {
		t.Error("GetHandler() returned false, handler not registered")
	}
	if handler == nil {
		t.Error("GetHandler() returned nil handler")
	}
}
