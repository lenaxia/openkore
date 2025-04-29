// Package config provides functionality for loading and managing network configurations.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

// ProxyType represents the type of proxy
type ProxyType string

// Proxy types
const (
	ProxyTypeNone   ProxyType = "none"
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeSOCKS4 ProxyType = "socks4"
	ProxyTypeSOCKS5 ProxyType = "socks5"
)

// ProxyConfig represents the configuration for a proxy
type ProxyConfig struct {
	Type     ProxyType `json:"type"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
}

// ReconnectPolicy represents the policy for reconnection attempts
type ReconnectPolicy struct {
	MaxAttempts     int           `json:"max_attempts"`
	InitialInterval time.Duration `json:"initial_interval"`
	MaxInterval     time.Duration `json:"max_interval"`
	Multiplier      float64       `json:"multiplier"`
	RandomFactor    float64       `json:"random_factor"`
}

// TimeoutConfig represents the configuration for network timeouts
type TimeoutConfig struct {
	Connect    time.Duration `json:"connect"`
	Read       time.Duration `json:"read"`
	Write      time.Duration `json:"write"`
	Idle       time.Duration `json:"idle"`
	KeepAlive  time.Duration `json:"keep_alive"`
	TLSTimeout time.Duration `json:"tls_timeout"`
}

// NetworkConfig represents the configuration for network behavior
type NetworkConfig struct {
	ServerName      string                 `json:"server_name"`
	Proxy           ProxyConfig            `json:"proxy"`
	Timeouts        TimeoutConfig          `json:"timeouts"`
	ReconnectPolicy ReconnectPolicy        `json:"reconnect_policy"`
	PacketDelay     time.Duration          `json:"packet_delay"`
	MaxPacketSize   int                    `json:"max_packet_size"`
	BufferSize      int                    `json:"buffer_size"`
	EnableTLS       bool                   `json:"enable_tls"`
	VerifyCert      bool                   `json:"verify_cert"`
	CustomFields    map[string]interface{} `json:"custom_fields"`
}

// NetworkConfigManager manages network configurations
type NetworkConfigManager struct {
	configs map[string]*NetworkConfig
}

// NewNetworkConfigManager creates a new network configuration manager
func NewNetworkConfigManager() *NetworkConfigManager {
	return &NetworkConfigManager{
		configs: make(map[string]*NetworkConfig),
	}
}

// LoadNetworkConfig loads a network configuration from a file
func (m *NetworkConfigManager) LoadNetworkConfig(path string) (*NetworkConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read network config file: %w", err)
	}

	var config NetworkConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse network config: %w", err)
	}

	if err := m.ValidateNetworkConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid network config: %w", err)
	}

	// Store the config by server name
	m.configs[config.ServerName] = &config

	return &config, nil
}

// GetNetworkConfig returns a network configuration by server name
func (m *NetworkConfigManager) GetNetworkConfig(serverName string) (*NetworkConfig, bool) {
	config, exists := m.configs[serverName]
	return config, exists
}

// GetNetworkConfigs returns all network configurations
func (m *NetworkConfigManager) GetNetworkConfigs() map[string]*NetworkConfig {
	return m.configs
}

// ValidateNetworkConfig validates a network configuration
func (m *NetworkConfigManager) ValidateNetworkConfig(config *NetworkConfig) error {
	if config.ServerName == "" {
		return errors.New("server name is required")
	}

	// Validate proxy configuration
	if config.Proxy.Type != ProxyTypeNone {
		if config.Proxy.Host == "" {
			return errors.New("proxy host is required when proxy type is not 'none'")
		}
		if config.Proxy.Port <= 0 || config.Proxy.Port > 65535 {
			return errors.New("invalid proxy port")
		}
	}

	// Validate reconnect policy
	if config.ReconnectPolicy.MaxAttempts < 0 {
		return errors.New("max reconnect attempts cannot be negative")
	}
	if config.ReconnectPolicy.InitialInterval < 0 {
		return errors.New("initial reconnect interval cannot be negative")
	}
	if config.ReconnectPolicy.MaxInterval < 0 {
		return errors.New("max reconnect interval cannot be negative")
	}
	if config.ReconnectPolicy.Multiplier <= 0 {
		return errors.New("reconnect multiplier must be positive")
	}
	if config.ReconnectPolicy.RandomFactor < 0 || config.ReconnectPolicy.RandomFactor > 1 {
		return errors.New("reconnect random factor must be between 0 and 1")
	}

	// Validate timeouts
	if config.Timeouts.Connect < 0 {
		return errors.New("connect timeout cannot be negative")
	}
	if config.Timeouts.Read < 0 {
		return errors.New("read timeout cannot be negative")
	}
	if config.Timeouts.Write < 0 {
		return errors.New("write timeout cannot be negative")
	}
	if config.Timeouts.Idle < 0 {
		return errors.New("idle timeout cannot be negative")
	}
	if config.Timeouts.KeepAlive < 0 {
		return errors.New("keep-alive timeout cannot be negative")
	}
	if config.Timeouts.TLSTimeout < 0 {
		return errors.New("TLS timeout cannot be negative")
	}

	// Validate other fields
	if config.PacketDelay < 0 {
		return errors.New("packet delay cannot be negative")
	}
	if config.MaxPacketSize <= 0 {
		return errors.New("max packet size must be positive")
	}
	if config.BufferSize <= 0 {
		return errors.New("buffer size must be positive")
	}

	return nil
}

// SaveNetworkConfig saves a network configuration to a file
func (m *NetworkConfigManager) SaveNetworkConfig(config *NetworkConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network config: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write network config file: %w", err)
	}

	return nil
}

// CreateDefaultNetworkConfig creates a default network configuration
func (m *NetworkConfigManager) CreateDefaultNetworkConfig(serverName string) *NetworkConfig {
	config := &NetworkConfig{
		ServerName: serverName,
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
		CustomFields:  make(map[string]interface{}),
	}

	m.configs[serverName] = config
	return config
}

// MergeNetworkConfigs merges two network configurations
// The values from the second config override the values from the first config
func MergeNetworkConfigs(base, override *NetworkConfig) *NetworkConfig {
	result := &NetworkConfig{
		ServerName:    override.ServerName,
		PacketDelay:   override.PacketDelay,
		MaxPacketSize: override.MaxPacketSize,
		BufferSize:    override.BufferSize,
		EnableTLS:     override.EnableTLS,
		VerifyCert:    override.VerifyCert,
	}

	// Merge proxy config
	if override.Proxy.Type != ProxyTypeNone {
		result.Proxy = override.Proxy
	} else {
		result.Proxy = base.Proxy
	}

	// Merge timeouts
	result.Timeouts = TimeoutConfig{
		Connect:    getOverrideValue(base.Timeouts.Connect, override.Timeouts.Connect),
		Read:       getOverrideValue(base.Timeouts.Read, override.Timeouts.Read),
		Write:      getOverrideValue(base.Timeouts.Write, override.Timeouts.Write),
		Idle:       getOverrideValue(base.Timeouts.Idle, override.Timeouts.Idle),
		KeepAlive:  getOverrideValue(base.Timeouts.KeepAlive, override.Timeouts.KeepAlive),
		TLSTimeout: getOverrideValue(base.Timeouts.TLSTimeout, override.Timeouts.TLSTimeout),
	}

	// Merge reconnect policy
	result.ReconnectPolicy = ReconnectPolicy{
		MaxAttempts:     getOverrideValue(base.ReconnectPolicy.MaxAttempts, override.ReconnectPolicy.MaxAttempts),
		InitialInterval: getOverrideValue(base.ReconnectPolicy.InitialInterval, override.ReconnectPolicy.InitialInterval),
		MaxInterval:     getOverrideValue(base.ReconnectPolicy.MaxInterval, override.ReconnectPolicy.MaxInterval),
		Multiplier:      getOverrideValue(base.ReconnectPolicy.Multiplier, override.ReconnectPolicy.Multiplier),
		RandomFactor:    getOverrideValue(base.ReconnectPolicy.RandomFactor, override.ReconnectPolicy.RandomFactor),
	}

	// Merge custom fields
	result.CustomFields = make(map[string]interface{})
	for k, v := range base.CustomFields {
		result.CustomFields[k] = v
	}
	for k, v := range override.CustomFields {
		result.CustomFields[k] = v
	}

	return result
}

// Helper function to get the override value if it's not the zero value
func getOverrideValue[T comparable](base, override T) T {
	var zero T
	if override != zero {
		return override
	}
	return base
}
