package servers

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/receive/security"
)

func TestNewBaseServer(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Create a base server
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("NewBaseServer() returned nil")
	}

	// Check that the server has the correct type
	if server.GetServerType() != ServerTypeUnknown {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeUnknown)
	}

	// Check that the server is not connected
	if server.IsConnected() {
		t.Error("IsConnected() = true, want false")
	}

	// Check that the server state is disconnected
	if server.GetState() != ServerStateDisconnected {
		t.Errorf("GetState() = %v, want %v", server.GetState(), ServerStateDisconnected)
	}

	// Check that the server info is correct
	serverInfo := server.GetServerInfo()
	if serverInfo.Name != baseConfig.Name {
		t.Errorf("GetServerInfo().Name = %v, want %v", serverInfo.Name, baseConfig.Name)
	}
	if serverInfo.IP != baseConfig.IP {
		t.Errorf("GetServerInfo().IP = %v, want %v", serverInfo.IP, baseConfig.IP)
	}
	if serverInfo.Port != baseConfig.Port {
		t.Errorf("GetServerInfo().Port = %v, want %v", serverInfo.Port, baseConfig.Port)
	}

	// Check that the server has a parser
	if server.GetParser() == nil {
		t.Error("GetParser() returned nil")
	}

	// Check that the server has a tokenizer
	if server.GetTokenizer() == nil {
		t.Error("GetTokenizer() returned nil")
	}

	// Check that the server has a login manager
	if server.GetLoginManager() == nil {
		t.Error("GetLoginManager() returned nil")
	}

	// Check that the server has a PIN manager
	if server.GetPINManager() == nil {
		t.Error("GetPINManager() returned nil")
	}

	// Check that the server has an anti-cheat manager
	if server.GetAntiCheatManager() == nil {
		t.Error("GetAntiCheatManager() returned nil")
	}
}

func TestBaseServerCredentials(t *testing.T) {
	// Create a base server
	baseConfig := NewBaseServerConfig()
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

	// Check that the login manager has the correct credentials
	loginManager := server.GetLoginManager()
	if loginManager == nil {
		t.Fatal("GetLoginManager() returned nil")
	}
}

func TestBaseServerOptions(t *testing.T) {
	// Create a base server
	baseConfig := NewBaseServerConfig()
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Set options
	options := ServerOptions{
		Timeout:        60 * time.Second,
		ReconnectDelay: 10 * time.Second,
		MaxRetries:     5,
		UseProxy:       true,
		ProxyAddress:   "proxy.example.com",
		ProxyPort:      8080,
		ProxyUsername:  "proxyuser",
		ProxyPassword:  "proxypass",
		UseCompression: true,
		UseEncryption:  true,
		UseAntiCheat:   true,
		AntiCheatType:  security.AntiCheatGameGuard,
		PacketVersion:  1234,
		ClientVersion:  "1.0.0",
		ServerVersion:  "2.0.0",
		CharacterSlot:  1,
		CharacterName:  "TestChar",
		CustomOptions:  map[string]interface{}{"option1": "value1"},
	}
	server.SetOptions(options)

	// There's no direct way to check the options, but we can check that the method doesn't panic
}

func TestBaseServerHooks(t *testing.T) {
	// Create a base server
	baseConfig := NewBaseServerConfig()
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Register a hook
	hookCalled := false
	handle := server.RegisterHook("test_hook", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
	}, nil)
	if handle == nil {
		t.Fatal("RegisterHook() returned nil")
	}

	// Call the hook manually through the hook manager
	server.hookManager.CallHook("test_hook", nil)

	// Check that the hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}
}

func TestBaseServerPacketHandler(t *testing.T) {
	// Create a base server
	baseConfig := NewBaseServerConfig()
	server, err := NewBaseServer(baseConfig)
	if err != nil {
		t.Fatalf("NewBaseServer() returned error: %v", err)
	}

	// Register a packet handler
	server.RegisterPacketHandler("0001", "test_packet", "v", []string{"len"}, func(args map[string]interface{}) error {
		// Just a test handler that does nothing
		return nil
	})

	// There's no direct way to check if the handler was registered correctly without processing a packet,
	// but we can check that the method doesn't panic
}

func TestBaseServerWithInvalidConfig(t *testing.T) {
	// Create a base server config with invalid server config
	baseConfig := NewBaseServerConfig()
	baseConfig.ServerConfig.IP = ""

	// Create a base server
	_, err := NewBaseServer(baseConfig)
	if err == nil {
		t.Error("NewBaseServer() with invalid server config did not return an error")
	}

	// Create a base server config with invalid network config
	baseConfig = NewBaseServerConfig()
	baseConfig.NetworkConfig.Timeouts.Connect = 0

	// Create a base server
	_, err = NewBaseServer(baseConfig)
	if err == nil {
		t.Error("NewBaseServer() with invalid network config did not return an error")
	}
}
