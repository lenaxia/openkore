package servers

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/receive/security"
)

func TestServerTypeString(t *testing.T) {
	testCases := []struct {
		serverType ServerType
		expected   string
	}{
		{ServerTypeUnknown, "Unknown"},
		{ServerTypeOfficial, "Official"},
		{ServerTypeSakray, "Sakray"},
		{ServerTypeRenewal, "Renewal"},
		{ServerTypeClassic, "Classic"},
		{ServerTypePreRenewal, "PreRenewal"},
		{ServerTypePrivate, "Private"},
		{ServerTypeThor, "Thor"},
		{ServerTypeZero, "Zero"},
		{ServerTypeRagexe, "Ragexe"},
		{ServerTypeRagexeRE, "RagexeRE"},
		{ServerTypeCustom, "Custom"},
		{ServerType(99), "Invalid"},
	}

	for _, tc := range testCases {
		result := tc.serverType.String()
		if result != tc.expected {
			t.Errorf("ServerType(%d).String() = %s, want %s", tc.serverType, result, tc.expected)
		}
	}
}

func TestServerStateString(t *testing.T) {
	testCases := []struct {
		serverState ServerState
		expected    string
	}{
		{ServerStateDisconnected, "Disconnected"},
		{ServerStateConnecting, "Connecting"},
		{ServerStateHandshaking, "Handshaking"},
		{ServerStateLoggingIn, "LoggingIn"},
		{ServerStateLoggedIn, "LoggedIn"},
		{ServerStateInGame, "InGame"},
		{ServerStateDisconnecting, "Disconnecting"},
		{ServerState(99), "Invalid"},
	}

	for _, tc := range testCases {
		result := tc.serverState.String()
		if result != tc.expected {
			t.Errorf("ServerState(%d).String() = %s, want %s", tc.serverState, result, tc.expected)
		}
	}
}

func TestNewBaseServerConfig(t *testing.T) {
	config := NewBaseServerConfig()

	if config == nil {
		t.Fatal("NewBaseServerConfig() returned nil")
	}

	if config.NetworkConfig == nil {
		t.Error("NetworkConfig is nil")
	}

	if config.ServerConfig == nil {
		t.Error("ServerConfig is nil")
	}

	if config.Options.Timeout != 30*time.Second {
		t.Errorf("Options.Timeout = %v, want %v", config.Options.Timeout, 30*time.Second)
	}

	if config.Options.ReconnectDelay != 5*time.Second {
		t.Errorf("Options.ReconnectDelay = %v, want %v", config.Options.ReconnectDelay, 5*time.Second)
	}

	if config.Options.MaxRetries != 3 {
		t.Errorf("Options.MaxRetries = %d, want %d", config.Options.MaxRetries, 3)
	}

	if !config.Options.UseCompression {
		t.Error("Options.UseCompression = false, want true")
	}

	if !config.Options.UseEncryption {
		t.Error("Options.UseEncryption = false, want true")
	}

	if config.Options.UseAntiCheat {
		t.Error("Options.UseAntiCheat = true, want false")
	}

	if config.Options.AntiCheatType != security.AntiCheatNone {
		t.Errorf("Options.AntiCheatType = %v, want %v", config.Options.AntiCheatType, security.AntiCheatNone)
	}

	if config.Options.CustomOptions == nil {
		t.Error("Options.CustomOptions is nil")
	}

	if config.CustomConfig == nil {
		t.Error("CustomConfig is nil")
	}
}

func TestDetectServerType(t *testing.T) {
	// Create a server config manager
	manager := config.NewServerConfigManager()

	// Test with explicit server type
	serverConfig := manager.CreateDefaultServerConfig("test")
	serverConfig.CustomFields["serverType"] = "sakray"
	serverType := DetectServerType(serverConfig)
	if serverType != ServerTypeSakray {
		t.Errorf("DetectServerType() with serverType=sakray = %v, want %v", serverType, ServerTypeSakray)
	}

	// Test with server version
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverConfig.CustomFields["serverVersion"] = "sakray_1.0.0"
	serverType = DetectServerType(serverConfig)
	if serverType != ServerTypeSakray {
		t.Errorf("DetectServerType() with serverVersion=sakray_1.0.0 = %v, want %v", serverType, ServerTypeSakray)
	}

	// Test with packet version
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverConfig.CustomFields["packetVersion"] = 20180000
	serverType = DetectServerType(serverConfig)
	if serverType != ServerTypeRenewal {
		t.Errorf("DetectServerType() with packetVersion=20180000 = %v, want %v", serverType, ServerTypeRenewal)
	}

	// Test with unknown
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverType = DetectServerType(serverConfig)
	if serverType != ServerTypeUnknown {
		t.Errorf("DetectServerType() with no info = %v, want %v", serverType, ServerTypeUnknown)
	}
}

func TestValidateServerConfig(t *testing.T) {
	// Create a server config manager
	manager := config.NewServerConfigManager()

	// Test with valid config
	serverConfig := manager.CreateDefaultServerConfig("test")
	err := ValidateServerConfig(serverConfig)
	if err != nil {
		t.Errorf("ValidateServerConfig() with valid config returned error: %v", err)
	}

	// Test with empty IP
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverConfig.IP = ""
	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server IP is not set" {
		t.Errorf("ValidateServerConfig() with empty IP returned error: %v, want 'server IP is not set'", err)
	}

	// Test with zero port
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverConfig.Port = 0
	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server port is not set" {
		t.Errorf("ValidateServerConfig() with zero port returned error: %v, want 'server port is not set'", err)
	}

	// Test with empty name
	serverConfig = manager.CreateDefaultServerConfig("test")
	serverConfig.Name = ""
	err = ValidateServerConfig(serverConfig)
	if err == nil || err.Error() != "server name is not set" {
		t.Errorf("ValidateServerConfig() with empty name returned error: %v, want 'server name is not set'", err)
	}
}

func TestValidateNetworkConfig(t *testing.T) {
	// Create a network config manager
	manager := config.NewNetworkConfigManager()

	// Test with valid config
	networkConfig := manager.CreateDefaultNetworkConfig("test")
	err := ValidateNetworkConfig(networkConfig)
	if err != nil {
		t.Errorf("ValidateNetworkConfig() with valid config returned error: %v", err)
	}

	// Test with zero connect timeout
	networkConfig = manager.CreateDefaultNetworkConfig("test")
	networkConfig.Timeouts.Connect = 0
	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "connection timeout must be greater than 0" {
		t.Errorf("ValidateNetworkConfig() with zero connect timeout returned error: %v, want 'connection timeout must be greater than 0'", err)
	}

	// Test with zero reconnect delay
	networkConfig = manager.CreateDefaultNetworkConfig("test")
	networkConfig.ReconnectPolicy.InitialInterval = 0
	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "reconnect delay must be greater than 0" {
		t.Errorf("ValidateNetworkConfig() with zero reconnect delay returned error: %v, want 'reconnect delay must be greater than 0'", err)
	}

	// Test with negative max retries
	networkConfig = manager.CreateDefaultNetworkConfig("test")
	networkConfig.ReconnectPolicy.MaxAttempts = -1
	err = ValidateNetworkConfig(networkConfig)
	if err == nil || err.Error() != "max retries must be greater than or equal to 0" {
		t.Errorf("ValidateNetworkConfig() with negative max retries returned error: %v, want 'max retries must be greater than or equal to 0'", err)
	}
}

func TestValidateCredentials(t *testing.T) {
	// Test with valid credentials
	credentials := ServerCredentials{
		Username: "test",
		Password: "password",
	}
	err := ValidateCredentials(credentials)
	if err != nil {
		t.Errorf("ValidateCredentials() with valid credentials returned error: %v", err)
	}

	// Test with empty username
	credentials = ServerCredentials{
		Username: "",
		Password: "password",
	}
	err = ValidateCredentials(credentials)
	if err == nil || err.Error() != "username is not set" {
		t.Errorf("ValidateCredentials() with empty username returned error: %v, want 'username is not set'", err)
	}

	// Test with empty password
	credentials = ServerCredentials{
		Username: "test",
		Password: "",
	}
	err = ValidateCredentials(credentials)
	if err == nil || err.Error() != "password is not set" {
		t.Errorf("ValidateCredentials() with empty password returned error: %v, want 'password is not set'", err)
	}
}
