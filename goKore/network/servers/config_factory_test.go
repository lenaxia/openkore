package servers

import (
	"testing"
)

// TestGetServerType0Config tests the GetServerType0Config function
func TestGetServerType0Config(t *testing.T) {
	config := GetServerType0Config()

	// Verify that the config is not nil
	if config == nil {
		t.Fatal("Expected non-nil config")
	}

	// Verify that the send config contains expected packets
	if len(config.SendConfig) == 0 {
		t.Error("Expected non-empty send config")
	}

	// Check for a specific packet that should be in the send config
	masterLogin, exists := config.SendConfig["0064"]
	if !exists {
		t.Error("Expected master_login packet in send config")
	} else if masterLogin.Name != "master_login" {
		t.Errorf("Expected packet name to be master_login, got %s", masterLogin.Name)
	}

	// Verify that the receive config contains expected packets
	if len(config.ReceiveConfig) == 0 {
		t.Error("Expected non-empty receive config")
	}

	// Check for a specific packet that should be in the receive config
	accountServerInfo, exists := config.ReceiveConfig["0069"]
	if !exists {
		t.Error("Expected account_server_info packet in receive config")
	} else if accountServerInfo.Name != "account_server_info" {
		t.Errorf("Expected packet name to be account_server_info, got %s", accountServerInfo.Name)
	}
}

// TestGetServerType0SendConfig tests the GetServerType0SendConfig function
func TestGetServerType0SendConfig(t *testing.T) {
	config := GetServerType0SendConfig()

	// Verify that the config is not nil
	if config == nil {
		t.Fatal("Expected non-nil config")
	}

	// Verify that the config contains expected packets
	if len(config) == 0 {
		t.Error("Expected non-empty config")
	}

	// Check for a specific packet that should be in the config
	masterLogin, exists := config["0064"]
	if !exists {
		t.Error("Expected master_login packet in config")
	} else if masterLogin.Name != "master_login" {
		t.Errorf("Expected packet name to be master_login, got %s", masterLogin.Name)
	}
}

// TestGetServerType0ReceiveConfig tests the GetServerType0ReceiveConfig function
func TestGetServerType0ReceiveConfig(t *testing.T) {
	config := GetServerType0ReceiveConfig()

	// Verify that the config is not nil
	if config == nil {
		t.Fatal("Expected non-nil config")
	}

	// Verify that the config contains expected packets
	if len(config) == 0 {
		t.Error("Expected non-empty config")
	}

	// Check for a specific packet that should be in the config
	accountServerInfo, exists := config["0069"]
	if !exists {
		t.Error("Expected account_server_info packet in config")
	} else if accountServerInfo.Name != "account_server_info" {
		t.Errorf("Expected packet name to be account_server_info, got %s", accountServerInfo.Name)
	}
}

// TestGetServerConfigByType tests the GetServerConfigByType function
func TestGetServerConfigByType(t *testing.T) {
	// Test with ServerTypeOfficial
	officialConfig := GetServerConfigByType(ServerTypeOfficial)
	if officialConfig == nil {
		t.Fatal("Expected non-nil config for ServerTypeOfficial")
	}

	// Test with ServerTypeSakray
	sakrayConfig := GetServerConfigByType(ServerTypeSakray)
	if sakrayConfig == nil {
		t.Fatal("Expected non-nil config for ServerTypeSakray")
	}

	// Test with unknown server type
	unknownConfig := GetServerConfigByType(ServerType(999))
	if unknownConfig == nil {
		t.Fatal("Expected non-nil config for unknown server type")
	}

	// Verify that the unknown server type returns the default config (ServerType0)
	masterLogin, exists := unknownConfig.SendConfig["0064"]
	if !exists {
		t.Error("Expected master_login packet in default config")
	} else if masterLogin.Name != "master_login" {
		t.Errorf("Expected packet name to be master_login, got %s", masterLogin.Name)
	}
}

// TestConfigIntegration tests that the configs can be used with other components
func TestConfigIntegration(t *testing.T) {
	// Get the server config
	config := GetServerType0Config()

	// Create a mock packet builder
	mockBuilder := &MockPacketBuilder{}

	// Register the packet formats with the mock builder
	for id, construction := range config.SendConfig {
		mockBuilder.RegisterPacket(id, construction.Name, construction.Format, construction.FieldNames)
	}

	// Verify that the mock builder has the expected packets
	if len(mockBuilder.packets) == 0 {
		t.Error("Expected non-empty packets in mock builder")
	}

	// Check for a specific packet that should be registered
	packet, exists := mockBuilder.packets["0064"]
	if !exists {
		t.Error("Expected master_login packet to be registered")
	} else if packet.name != "master_login" {
		t.Errorf("Expected packet name to be master_login, got %s", packet.name)
	}
}

// MockPacketBuilder is a mock implementation of a packet builder
type MockPacketBuilder struct {
	packets map[string]struct {
		name       string
		format     string
		fieldNames []string
	}
}

// RegisterPacket registers a packet with the mock builder
func (m *MockPacketBuilder) RegisterPacket(id, name, format string, fieldNames []string) {
	if m.packets == nil {
		m.packets = make(map[string]struct {
			name       string
			format     string
			fieldNames []string
		})
	}
	m.packets[id] = struct {
		name       string
		format     string
		fieldNames []string
	}{
		name:       name,
		format:     format,
		fieldNames: fieldNames,
	}
}
