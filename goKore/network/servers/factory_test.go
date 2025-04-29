package servers

import (
	"testing"
)

func TestServerFactory(t *testing.T) {
	// Create a server factory
	factory := NewServerFactory()
	if factory == nil {
		t.Fatal("NewServerFactory() returned nil")
	}

	// Create a base server config
	baseConfig := NewBaseServerConfig()
	if baseConfig == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	// Set server type to Sakray
	baseConfig.Type = ServerTypeSakray

	// Create a server
	server, err := factory.CreateServer(baseConfig)
	if err != nil {
		t.Fatalf("CreateServer() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("CreateServer() returned nil")
	}

	// Check that the server is of the correct type
	if server.GetServerType() != ServerTypeSakray {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeSakray)
	}

	// Test with official server type
	baseConfig.Type = ServerTypeOfficial
	server, err = factory.CreateServer(baseConfig)
	if err != nil {
		t.Fatalf("CreateServer() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("CreateServer() returned nil")
	}
	if server.GetServerType() != ServerTypeOfficial {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeOfficial)
	}

	// Test with unknown server type (should default to base server)
	baseConfig.Type = ServerTypeUnknown
	baseConfig.ServerConfig.CustomFields = map[string]interface{}{
		"serverType": "sakray",
	}
	server, err = factory.CreateServer(baseConfig)
	if err != nil {
		t.Fatalf("CreateServer() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("CreateServer() returned nil")
	}
	if server.GetServerType() != ServerTypeSakray {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeSakray)
	}

	// Test CreateServerFromConfig
	baseConfig = NewBaseServerConfig()
	baseConfig.Type = ServerTypeSakray
	server, err = CreateServerFromConfig(baseConfig)
	if err != nil {
		t.Fatalf("CreateServerFromConfig() returned error: %v", err)
	}
	if server == nil {
		t.Fatal("CreateServerFromConfig() returned nil")
	}
	if server.GetServerType() != ServerTypeSakray {
		t.Errorf("GetServerType() = %v, want %v", server.GetServerType(), ServerTypeSakray)
	}
}
