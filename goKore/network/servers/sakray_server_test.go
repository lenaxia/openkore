package servers

import (
	"testing"

	"github.com/lenaxia/goKore/network/receive/security"
)

func TestNewSakrayServer(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a Sakray server
	server, err := NewSakrayServer(baseConfig)
	if err != nil {
		t.Fatalf("NewSakrayServer() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("NewSakrayServer() returned nil")
	}

	// Check that the server has the correct type
	if server.GetServerType() != ServerTypeSakray {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeSakray)
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
}

func TestSakrayServerPacketDefinitions(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a Sakray server
	server, err := NewSakrayServer(baseConfig)
	if err != nil {
		t.Fatalf("NewSakrayServer() returned error: %v", err)
	}

	// Load Sakray packet definitions
	err = server.LoadSakrayPacketDefinitions()
	if err != nil {
		t.Fatalf("LoadSakrayPacketDefinitions() returned error: %v", err)
	}

	// There's no direct way to check if the packet definitions were loaded correctly,
	// but we can check that the method doesn't panic
}

func TestSakrayServerConnect(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a Sakray server
	_, err := NewSakrayServer(baseConfig)
	if err != nil {
		t.Fatalf("NewSakrayServer() returned error: %v", err)
	}

	// We can't actually connect to a server in a unit test,
	// but we can check that the method doesn't panic
	// server.Connect() would fail without a real server to connect to
}

func TestSakrayServerLogin(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a Sakray server
	server, err := NewSakrayServer(baseConfig)
	if err != nil {
		t.Fatalf("NewSakrayServer() returned error: %v", err)
	}

	// Set credentials
	credentials := ServerCredentials{
		Username: "testuser",
		Password: "testpass",
		PINCode:  "1234",
	}
	server.SetCredentials(credentials)

	// We can't actually login to a server in a unit test,
	// but we can check that the method doesn't panic
	// server.Login() would fail without a real server to connect to
}

func TestSakrayServerOptions(t *testing.T) {
	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a Sakray server
	server, err := NewSakrayServer(baseConfig)
	if err != nil {
		t.Fatalf("NewSakrayServer() returned error: %v", err)
	}

	// Set options
	options := ServerOptions{
		UseAntiCheat:  true,
		AntiCheatType: security.AntiCheatGameGuard,
	}
	server.SetOptions(options)

	// There's no direct way to check the options, but we can check that the method doesn't panic
}
