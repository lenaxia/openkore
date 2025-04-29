package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNetworkConfigManager_CreateDefaultNetworkConfig(t *testing.T) {
	manager := NewNetworkConfigManager()
	config := manager.CreateDefaultNetworkConfig("TestServer")

	if config.ServerName != "TestServer" {
		t.Errorf("Expected server name to be 'TestServer', got '%s'", config.ServerName)
	}

	if config.Proxy.Type != ProxyTypeNone {
		t.Errorf("Expected proxy type to be 'none', got '%s'", config.Proxy.Type)
	}

	if config.Timeouts.Connect != 10*time.Second {
		t.Errorf("Expected connect timeout to be 10s, got %v", config.Timeouts.Connect)
	}

	if config.ReconnectPolicy.MaxAttempts != 5 {
		t.Errorf("Expected max reconnect attempts to be 5, got %d", config.ReconnectPolicy.MaxAttempts)
	}

	if config.PacketDelay != 50*time.Millisecond {
		t.Errorf("Expected packet delay to be 50ms, got %v", config.PacketDelay)
	}

	if config.MaxPacketSize != 16384 {
		t.Errorf("Expected max packet size to be 16384, got %d", config.MaxPacketSize)
	}

	if config.BufferSize != 8192 {
		t.Errorf("Expected buffer size to be 8192, got %d", config.BufferSize)
	}

	if config.EnableTLS != false {
		t.Errorf("Expected enable TLS to be false, got %v", config.EnableTLS)
	}

	if config.VerifyCert != true {
		t.Errorf("Expected verify cert to be true, got %v", config.VerifyCert)
	}

	// Check if the config was added to the manager
	storedConfig, exists := manager.GetNetworkConfig("TestServer")
	if !exists {
		t.Error("Expected network config to be stored in the manager")
	}

	if storedConfig != config {
		t.Error("Expected stored config to be the same as the created config")
	}
}

func TestNetworkConfigManager_ValidateNetworkConfig(t *testing.T) {
	manager := NewNetworkConfigManager()

	// Test valid config
	validConfig := &NetworkConfig{
		ServerName: "TestServer",
		Proxy: ProxyConfig{
			Type: ProxyTypeNone,
		},
		Timeouts: TimeoutConfig{
			Connect:    10 * time.Second,
			Read:       30 * time.Second,
			Write:      30 * time.Second,
			Idle:       60 * time.Second,
			KeepAlive:  60 * time.Second,
			TLSTimeout: 10 * time.Second,
		},
		ReconnectPolicy: ReconnectPolicy{
			MaxAttempts:     5,
			InitialInterval: 1 * time.Second,
			MaxInterval:     30 * time.Second,
			Multiplier:      2.0,
			RandomFactor:    0.5,
		},
		PacketDelay:   50 * time.Millisecond,
		MaxPacketSize: 16384,
		BufferSize:    8192,
		EnableTLS:     false,
		VerifyCert:    true,
	}
	err := manager.ValidateNetworkConfig(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}

	// Test invalid configs
	testCases := []struct {
		name   string
		config *NetworkConfig
	}{
		{
			name: "Missing server name",
			config: &NetworkConfig{
				Proxy: ProxyConfig{
					Type: ProxyTypeNone,
				},
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid proxy config",
			config: &NetworkConfig{
				ServerName: "TestServer",
				Proxy: ProxyConfig{
					Type: ProxyTypeHTTP,
					Port: -1,
				},
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid reconnect policy",
			config: &NetworkConfig{
				ServerName: "TestServer",
				Proxy: ProxyConfig{
					Type: ProxyTypeNone,
				},
				ReconnectPolicy: ReconnectPolicy{
					MaxAttempts:     -1,
					InitialInterval: 1 * time.Second,
					MaxInterval:     30 * time.Second,
					Multiplier:      2.0,
					RandomFactor:    0.5,
				},
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid timeouts",
			config: &NetworkConfig{
				ServerName: "TestServer",
				Proxy: ProxyConfig{
					Type: ProxyTypeNone,
				},
				Timeouts: TimeoutConfig{
					Connect:    -10 * time.Second,
					Read:       30 * time.Second,
					Write:      30 * time.Second,
					Idle:       60 * time.Second,
					KeepAlive:  60 * time.Second,
					TLSTimeout: 10 * time.Second,
				},
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid packet delay",
			config: &NetworkConfig{
				ServerName:    "TestServer",
				PacketDelay:   -50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid max packet size",
			config: &NetworkConfig{
				ServerName:    "TestServer",
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 0,
				BufferSize:    8192,
			},
		},
		{
			name: "Invalid buffer size",
			config: &NetworkConfig{
				ServerName:    "TestServer",
				PacketDelay:   50 * time.Millisecond,
				MaxPacketSize: 16384,
				BufferSize:    0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.ValidateNetworkConfig(tc.config)
			if err == nil {
				t.Error("Expected error for invalid config, got nil")
			}
		})
	}
}

func TestNetworkConfigManager_LoadSaveNetworkConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "network_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewNetworkConfigManager()

	// Create a test config
	config := &NetworkConfig{
		ServerName: "TestServer",
		Proxy: ProxyConfig{
			Type:     ProxyTypeHTTP,
			Host:     "proxy.example.com",
			Port:     8080,
			Username: "user",
			Password: "pass",
		},
		Timeouts: TimeoutConfig{
			Connect:    10 * time.Second,
			Read:       30 * time.Second,
			Write:      30 * time.Second,
			Idle:       60 * time.Second,
			KeepAlive:  60 * time.Second,
			TLSTimeout: 10 * time.Second,
		},
		ReconnectPolicy: ReconnectPolicy{
			MaxAttempts:     5,
			InitialInterval: 1 * time.Second,
			MaxInterval:     30 * time.Second,
			Multiplier:      2.0,
			RandomFactor:    0.5,
		},
		PacketDelay:   50 * time.Millisecond,
		MaxPacketSize: 16384,
		BufferSize:    8192,
		EnableTLS:     true,
		VerifyCert:    true,
		CustomFields: map[string]interface{}{
			"custom": "value",
		},
	}

	// Save the config to a file
	configPath := filepath.Join(tempDir, "test_network.json")
	err = manager.SaveNetworkConfig(config, configPath)
	if err != nil {
		t.Fatalf("Failed to save network config: %v", err)
	}

	// Load the config from the file
	loadedConfig, err := manager.LoadNetworkConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load network config: %v", err)
	}

	// Check if the loaded config matches the original
	if loadedConfig.ServerName != config.ServerName {
		t.Errorf("Expected loaded config server name to be '%s', got '%s'", config.ServerName, loadedConfig.ServerName)
	}

	if loadedConfig.Proxy.Type != config.Proxy.Type {
		t.Errorf("Expected loaded config proxy type to be '%s', got '%s'", config.Proxy.Type, loadedConfig.Proxy.Type)
	}

	if loadedConfig.Proxy.Host != config.Proxy.Host {
		t.Errorf("Expected loaded config proxy host to be '%s', got '%s'", config.Proxy.Host, loadedConfig.Proxy.Host)
	}

	if loadedConfig.Proxy.Port != config.Proxy.Port {
		t.Errorf("Expected loaded config proxy port to be %d, got %d", config.Proxy.Port, loadedConfig.Proxy.Port)
	}

	if loadedConfig.EnableTLS != config.EnableTLS {
		t.Errorf("Expected loaded config enable TLS to be %v, got %v", config.EnableTLS, loadedConfig.EnableTLS)
	}

	// Check if the config was added to the manager
	storedConfig, exists := manager.GetNetworkConfig("TestServer")
	if !exists {
		t.Error("Expected loaded config to be stored in the manager")
	}

	if storedConfig != loadedConfig {
		t.Error("Expected stored config to be the same as the loaded config")
	}
}

func TestMergeNetworkConfigs(t *testing.T) {
	// Create base config
	base := &NetworkConfig{
		ServerName: "BaseServer",
		Proxy: ProxyConfig{
			Type: ProxyTypeHTTP,
			Host: "base-proxy.example.com",
			Port: 8080,
		},
		Timeouts: TimeoutConfig{
			Connect:    10 * time.Second,
			Read:       30 * time.Second,
			Write:      30 * time.Second,
			Idle:       60 * time.Second,
			KeepAlive:  60 * time.Second,
			TLSTimeout: 10 * time.Second,
		},
		ReconnectPolicy: ReconnectPolicy{
			MaxAttempts:     5,
			InitialInterval: 1 * time.Second,
			MaxInterval:     30 * time.Second,
			Multiplier:      2.0,
			RandomFactor:    0.5,
		},
		PacketDelay:   50 * time.Millisecond,
		MaxPacketSize: 16384,
		BufferSize:    8192,
		EnableTLS:     false,
		VerifyCert:    true,
		CustomFields: map[string]interface{}{
			"base":   "value",
			"shared": "base-value",
		},
	}

	// Create override config
	override := &NetworkConfig{
		ServerName: "OverrideServer",
		Proxy: ProxyConfig{
			Type: ProxyTypeSOCKS5,
			Host: "override-proxy.example.com",
			Port: 1080,
		},
		Timeouts: TimeoutConfig{
			Connect: 5 * time.Second,
			// Other timeouts are zero values
		},
		ReconnectPolicy: ReconnectPolicy{
			MaxAttempts: 3,
			// Other reconnect policy values are zero values
		},
		PacketDelay:   100 * time.Millisecond,
		MaxPacketSize: 32768,
		BufferSize:    16384,
		EnableTLS:     true,
		VerifyCert:    false,
		CustomFields: map[string]interface{}{
			"override": "value",
			"shared":   "override-value",
		},
	}

	// Merge configs
	merged := MergeNetworkConfigs(base, override)

	// Check merged values
	if merged.ServerName != override.ServerName {
		t.Errorf("Expected merged server name to be '%s', got '%s'", override.ServerName, merged.ServerName)
	}

	if merged.Proxy.Type != override.Proxy.Type {
		t.Errorf("Expected merged proxy type to be '%s', got '%s'", override.Proxy.Type, merged.Proxy.Type)
	}

	if merged.Proxy.Host != override.Proxy.Host {
		t.Errorf("Expected merged proxy host to be '%s', got '%s'", override.Proxy.Host, merged.Proxy.Host)
	}

	if merged.Timeouts.Connect != override.Timeouts.Connect {
		t.Errorf("Expected merged connect timeout to be %v, got %v", override.Timeouts.Connect, merged.Timeouts.Connect)
	}

	if merged.Timeouts.Read != base.Timeouts.Read {
		t.Errorf("Expected merged read timeout to be %v, got %v", base.Timeouts.Read, merged.Timeouts.Read)
	}

	if merged.ReconnectPolicy.MaxAttempts != override.ReconnectPolicy.MaxAttempts {
		t.Errorf("Expected merged max attempts to be %d, got %d", override.ReconnectPolicy.MaxAttempts, merged.ReconnectPolicy.MaxAttempts)
	}

	if merged.ReconnectPolicy.InitialInterval != base.ReconnectPolicy.InitialInterval {
		t.Errorf("Expected merged initial interval to be %v, got %v", base.ReconnectPolicy.InitialInterval, merged.ReconnectPolicy.InitialInterval)
	}

	if merged.PacketDelay != override.PacketDelay {
		t.Errorf("Expected merged packet delay to be %v, got %v", override.PacketDelay, merged.PacketDelay)
	}

	if merged.MaxPacketSize != override.MaxPacketSize {
		t.Errorf("Expected merged max packet size to be %d, got %d", override.MaxPacketSize, merged.MaxPacketSize)
	}

	if merged.EnableTLS != override.EnableTLS {
		t.Errorf("Expected merged enable TLS to be %v, got %v", override.EnableTLS, merged.EnableTLS)
	}

	// Check custom fields
	if merged.CustomFields["base"] != base.CustomFields["base"] {
		t.Errorf("Expected merged custom field 'base' to be '%v', got '%v'", base.CustomFields["base"], merged.CustomFields["base"])
	}

	if merged.CustomFields["override"] != override.CustomFields["override"] {
		t.Errorf("Expected merged custom field 'override' to be '%v', got '%v'", override.CustomFields["override"], merged.CustomFields["override"])
	}

	if merged.CustomFields["shared"] != override.CustomFields["shared"] {
		t.Errorf("Expected merged custom field 'shared' to be '%v', got '%v'", override.CustomFields["shared"], merged.CustomFields["shared"])
	}
}

func TestGetOverrideValue(t *testing.T) {
	// Test with int
	baseInt := 10
	overrideInt := 20
	zeroInt := 0

	if getOverrideValue(baseInt, overrideInt) != overrideInt {
		t.Errorf("Expected override value to be %d, got %d", overrideInt, getOverrideValue(baseInt, overrideInt))
	}

	if getOverrideValue(baseInt, zeroInt) != baseInt {
		t.Errorf("Expected override value to be %d, got %d", baseInt, getOverrideValue(baseInt, zeroInt))
	}

	// Test with string
	baseString := "base"
	overrideString := "override"
	zeroString := ""

	if getOverrideValue(baseString, overrideString) != overrideString {
		t.Errorf("Expected override value to be '%s', got '%s'", overrideString, getOverrideValue(baseString, overrideString))
	}

	if getOverrideValue(baseString, zeroString) != baseString {
		t.Errorf("Expected override value to be '%s', got '%s'", baseString, getOverrideValue(baseString, zeroString))
	}

	// Test with time.Duration
	baseDuration := 10 * time.Second
	overrideDuration := 20 * time.Second
	zeroDuration := time.Duration(0)

	if getOverrideValue(baseDuration, overrideDuration) != overrideDuration {
		t.Errorf("Expected override value to be %v, got %v", overrideDuration, getOverrideValue(baseDuration, overrideDuration))
	}

	if getOverrideValue(baseDuration, zeroDuration) != baseDuration {
		t.Errorf("Expected override value to be %v, got %v", baseDuration, getOverrideValue(baseDuration, zeroDuration))
	}
}
