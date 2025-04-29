package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestServerConfigManager_CreateDefaultServerConfig(t *testing.T) {
	manager := NewServerConfigManager()
	config := manager.CreateDefaultServerConfig("TestServer")

	if config.Name != "TestServer" {
		t.Errorf("Expected server name to be 'TestServer', got '%s'", config.Name)
	}

	if config.Type != ServerType0 {
		t.Errorf("Expected server type to be 'ServerType0', got '%s'", config.Type)
	}

	if config.IP != "127.0.0.1" {
		t.Errorf("Expected server IP to be '127.0.0.1', got '%s'", config.IP)
	}

	if config.Port != 6900 {
		t.Errorf("Expected server port to be 6900, got %d", config.Port)
	}

	// Check if the config was added to the manager
	storedConfig, exists := manager.GetServerConfig("TestServer")
	if !exists {
		t.Error("Expected server config to be stored in the manager")
	}

	if storedConfig != config {
		t.Error("Expected stored config to be the same as the created config")
	}
}

func TestServerConfigManager_ValidateServerConfig(t *testing.T) {
	manager := NewServerConfigManager()

	// Test valid config
	validConfig := &ServerConfig{
		Name: "TestServer",
		IP:   "127.0.0.1",
		Port: 6900,
	}
	err := manager.ValidateServerConfig(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}

	// Test server type detection
	if validConfig.Type != ServerType0 {
		t.Errorf("Expected server type to be detected as 'ServerType0', got '%s'", validConfig.Type)
	}

	// Test default values
	if validConfig.ServerEncoding != "UTF-8" {
		t.Errorf("Expected server encoding to default to 'UTF-8', got '%s'", validConfig.ServerEncoding)
	}

	if validConfig.CharBlockSize != 106 {
		t.Errorf("Expected char block size to default to 106, got %d", validConfig.CharBlockSize)
	}

	// Test invalid configs
	testCases := []struct {
		name   string
		config *ServerConfig
	}{
		{
			name:   "Missing name",
			config: &ServerConfig{IP: "127.0.0.1", Port: 6900},
		},
		{
			name:   "Missing IP",
			config: &ServerConfig{Name: "TestServer", Port: 6900},
		},
		{
			name:   "Invalid port (negative)",
			config: &ServerConfig{Name: "TestServer", IP: "127.0.0.1", Port: -1},
		},
		{
			name:   "Invalid port (too large)",
			config: &ServerConfig{Name: "TestServer", IP: "127.0.0.1", Port: 70000},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.ValidateServerConfig(tc.config)
			if err == nil {
				t.Error("Expected error for invalid config, got nil")
			}
		})
	}
}

func TestServerConfigManager_DetectServerType(t *testing.T) {
	manager := NewServerConfigManager()

	testCases := []struct {
		name     string
		config   *ServerConfig
		expected ServerType
	}{
		{
			name:     "Explicit type",
			config:   &ServerConfig{Name: "TestServer", Type: ServerTypeSakray},
			expected: ServerTypeSakray,
		},
		{
			name:     "Sakray in name",
			config:   &ServerConfig{Name: "SakrayServer"},
			expected: ServerTypeSakray,
		},
		{
			name:     "bRO in name",
			config:   &ServerConfig{Name: "bRO_Server"},
			expected: ServerTypeBRO,
		},
		{
			name:     "iRO in name",
			config:   &ServerConfig{Name: "iRO_Server"},
			expected: ServerTypeIRO,
		},
		{
			name:     "euRO in name",
			config:   &ServerConfig{Name: "euRO_Server"},
			expected: ServerTypeEURO,
		},
		{
			name:     "Default type",
			config:   &ServerConfig{Name: "GenericServer"},
			expected: ServerType0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serverType := manager.DetectServerType(tc.config)
			if serverType != tc.expected {
				t.Errorf("Expected server type to be '%s', got '%s'", tc.expected, serverType)
			}
		})
	}
}

func TestServerConfigManager_LoadSaveServerConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "server_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewServerConfigManager()

	// Create a test config
	config := &ServerConfig{
		Name:           "TestServer",
		Type:           ServerTypeSakray,
		IP:             "192.168.1.1",
		Port:           6900,
		MasterVersion:  1,
		Version:        22,
		ServerEncoding: "UTF-8",
		RecvPacketKeys: []int{1, 2, 3},
		SendPacketKeys: []int{4, 5, 6},
		ServerTables:   map[string]string{"items": "items.txt"},
		CustomFields:   map[string]interface{}{"custom": "value"},
	}

	// Save the config to a file
	configPath := filepath.Join(tempDir, "test_server.json")
	err = manager.SaveServerConfig(config, configPath)
	if err != nil {
		t.Fatalf("Failed to save server config: %v", err)
	}

	// Load the config from the file
	loadedConfig, err := manager.LoadServerConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load server config: %v", err)
	}

	// Check if the loaded config matches the original
	if loadedConfig.Name != config.Name {
		t.Errorf("Expected loaded config name to be '%s', got '%s'", config.Name, loadedConfig.Name)
	}

	if loadedConfig.Type != config.Type {
		t.Errorf("Expected loaded config type to be '%s', got '%s'", config.Type, loadedConfig.Type)
	}

	if loadedConfig.IP != config.IP {
		t.Errorf("Expected loaded config IP to be '%s', got '%s'", config.IP, loadedConfig.IP)
	}

	if loadedConfig.Port != config.Port {
		t.Errorf("Expected loaded config port to be %d, got %d", config.Port, loadedConfig.Port)
	}

	// Check if the config was added to the manager
	storedConfig, exists := manager.GetServerConfig("TestServer")
	if !exists {
		t.Error("Expected loaded config to be stored in the manager")
	}

	if storedConfig != loadedConfig {
		t.Error("Expected stored config to be the same as the loaded config")
	}
}

func TestServerConfigManager_LoadServerConfigs(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "server_configs_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configs
	configs := []*ServerConfig{
		{
			Name: "Server1",
			IP:   "192.168.1.1",
			Port: 6900,
		},
		{
			Name: "Server2",
			IP:   "192.168.1.2",
			Port: 6901,
		},
		{
			Name: "Server3",
			IP:   "192.168.1.3",
			Port: 6902,
		},
	}

	// Save configs to files
	for i, config := range configs {
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal config %d: %v", i, err)
		}

		configPath := filepath.Join(tempDir, config.Name+".json")
		err = ioutil.WriteFile(configPath, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write config %d: %v", i, err)
		}
	}

	// Create a non-JSON file to test filtering
	nonJsonPath := filepath.Join(tempDir, "not_a_config.txt")
	err = ioutil.WriteFile(nonJsonPath, []byte("This is not a JSON file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write non-JSON file: %v", err)
	}

	// Create a subdirectory to test directory filtering
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Load configs from the directory
	manager := NewServerConfigManager()
	err = manager.LoadServerConfigs(tempDir)
	if err != nil {
		t.Fatalf("Failed to load server configs: %v", err)
	}

	// Check if all configs were loaded
	for _, config := range configs {
		loadedConfig, exists := manager.GetServerConfig(config.Name)
		if !exists {
			t.Errorf("Expected config '%s' to be loaded", config.Name)
			continue
		}

		if loadedConfig.Name != config.Name {
			t.Errorf("Expected loaded config name to be '%s', got '%s'", config.Name, loadedConfig.Name)
		}

		if loadedConfig.IP != config.IP {
			t.Errorf("Expected loaded config IP to be '%s', got '%s'", config.IP, loadedConfig.IP)
		}

		if loadedConfig.Port != config.Port {
			t.Errorf("Expected loaded config port to be %d, got %d", config.Port, loadedConfig.Port)
		}
	}

	// Check the total number of loaded configs
	if len(manager.GetServerConfigs()) != len(configs) {
		t.Errorf("Expected %d configs to be loaded, got %d", len(configs), len(manager.GetServerConfigs()))
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "file_exists_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test existing file
	if !FileExists(tempFile.Name()) {
		t.Errorf("Expected FileExists to return true for existing file")
	}

	// Test non-existing file
	if FileExists(tempFile.Name() + ".nonexistent") {
		t.Errorf("Expected FileExists to return false for non-existing file")
	}
}

func TestGetServerConfigPath(t *testing.T) {
	configDir := "/path/to/configs"
	serverName := "TestServer"
	expected := filepath.Join(configDir, "TestServer.json")

	path := GetServerConfigPath(configDir, serverName)
	if path != expected {
		t.Errorf("Expected path to be '%s', got '%s'", expected, path)
	}
}
