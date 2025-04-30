// Package config provides functionality for loading and managing server-specific configurations.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ServerType represents the type of server
type ServerType string

// Server types
const (
	ServerTypeUnknown ServerType = "unknown"
	ServerType0       ServerType = "ServerType0"
	ServerTypeSakray  ServerType = "Sakray"
	ServerTypeBRO     ServerType = "bRO"
	ServerTypeIRO     ServerType = "iRO"
	ServerTypeEURO    ServerType = "euRO"
)

// ServerConfig represents the configuration for a specific server
type ServerConfig struct {
	Name            string                 `json:"name"`
	Type            ServerType             `json:"type"`
	IP              string                 `json:"ip"`
	Port            int                    `json:"port"`
	MasterVersion   int                    `json:"master_version"`
	Version         int                    `json:"version"`
	ServerEncoding  string                 `json:"server_encoding"`
	PincodeStep     int                    `json:"pincode_step"`
	CharBlockSize   int                    `json:"char_block_size"`
	CharDeleteDate  int                    `json:"char_delete_date"`
	GameGuard       bool                   `json:"game_guard"`
	LoginPacketVer  int                    `json:"login_packet_ver"`
	RecvPacketKeys  []int                  `json:"recv_packet_keys"`
	SendPacketKeys  []int                  `json:"send_packet_keys"`
	PacketKeys      []int                  `json:"packet_keys"`
	PacketObfuscate bool                   `json:"packet_obfuscate"`
	TableFolders    []string               `json:"table_folders"`
	ServerTables    map[string]string      `json:"server_tables"`
	CustomFields    map[string]interface{} `json:"custom_fields"`
}

// ParseTableFolders parses a semicolon-delimited string of table folders
func (c *ServerConfig) ParseTableFolders(foldersStr string) {
	if foldersStr == "" {
		c.TableFolders = []string{}
		return
	}

	// Split by semicolon
	folders := strings.Split(foldersStr, ";")

	// Trim whitespace from each folder
	for i, folder := range folders {
		folders[i] = strings.TrimSpace(folder)
	}

	c.TableFolders = folders
}

// ServerConfigManager manages server configurations
type ServerConfigManager struct {
	configs map[string]*ServerConfig
}

// NewServerConfigManager creates a new server configuration manager
func NewServerConfigManager() *ServerConfigManager {
	return &ServerConfigManager{
		configs: make(map[string]*ServerConfig),
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for ServerConfig
func (c *ServerConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid infinite recursion
	type ServerConfigAlias ServerConfig

	// First, try to unmarshal normally
	alias := &struct {
		*ServerConfigAlias
		TableFolders    interface{} `json:"table_folders"`
		AddTableFolders string      `json:"addTableFolders"`
	}{
		ServerConfigAlias: (*ServerConfigAlias)(c),
	}

	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	// Handle table_folders field based on its type
	switch v := alias.TableFolders.(type) {
	case string:
		// If it's a string, parse it
		c.ParseTableFolders(v)
	case []interface{}:
		// If it's an array, convert to []string
		c.TableFolders = make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				c.TableFolders[i] = str
			}
		}
	}

	// Handle legacy addTableFolders field if table_folders is empty
	if len(c.TableFolders) == 0 && alias.AddTableFolders != "" {
		c.ParseTableFolders(alias.AddTableFolders)
	}

	return nil
}

// LoadServerConfig loads a server configuration from a file
func (m *ServerConfigManager) LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read server config file: %w", err)
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse server config: %w", err)
	}

	if err := m.ValidateServerConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid server config: %w", err)
	}

	// Store the config by name
	m.configs[config.Name] = &config

	return &config, nil
}

// LoadServerConfigs loads all server configurations from a directory
func (m *ServerConfigManager) LoadServerConfigs(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read server config directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		path := filepath.Join(dirPath, file.Name())
		_, err := m.LoadServerConfig(path)
		if err != nil {
			return fmt.Errorf("failed to load server config %s: %w", file.Name(), err)
		}
	}

	return nil
}

// GetServerConfig returns a server configuration by name
func (m *ServerConfigManager) GetServerConfig(name string) (*ServerConfig, bool) {
	config, exists := m.configs[name]
	return config, exists
}

// GetServerConfigs returns all server configurations
func (m *ServerConfigManager) GetServerConfigs() map[string]*ServerConfig {
	return m.configs
}

// DetectServerType detects the server type based on the configuration
func (m *ServerConfigManager) DetectServerType(config *ServerConfig) ServerType {
	// If the type is already set, return it
	if config.Type != ServerTypeUnknown && config.Type != "" {
		return config.Type
	}

	// Try to detect the server type based on the configuration
	// This is a simplified version, in a real implementation we would use more heuristics
	if strings.Contains(strings.ToLower(config.Name), "sakray") {
		return ServerTypeSakray
	} else if strings.Contains(strings.ToLower(config.Name), "bro") {
		return ServerTypeBRO
	} else if strings.Contains(strings.ToLower(config.Name), "iro") {
		return ServerTypeIRO
	} else if strings.Contains(strings.ToLower(config.Name), "euro") {
		return ServerTypeEURO
	}

	// Default to ServerType0
	return ServerType0
}

// ValidateServerConfig validates a server configuration
func (m *ServerConfigManager) ValidateServerConfig(config *ServerConfig) error {
	if config.Name == "" {
		return errors.New("server name is required")
	}

	if config.IP == "" {
		return errors.New("server IP is required")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return errors.New("invalid server port")
	}

	// If server type is not set, detect it
	if config.Type == ServerTypeUnknown || config.Type == "" {
		config.Type = m.DetectServerType(config)
	}

	// Set default values for optional fields
	if config.ServerEncoding == "" {
		config.ServerEncoding = "UTF-8"
	}

	if config.CharBlockSize == 0 {
		config.CharBlockSize = 106
	}

	return nil
}

// SaveServerConfig saves a server configuration to a file
func (m *ServerConfigManager) SaveServerConfig(config *ServerConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal server config: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write server config file: %w", err)
	}

	return nil
}

// CreateDefaultServerConfig creates a default server configuration
func (m *ServerConfigManager) CreateDefaultServerConfig(name string) *ServerConfig {
	config := &ServerConfig{
		Name:            name,
		Type:            ServerType0,
		IP:              "127.0.0.1",
		Port:            6900,
		MasterVersion:   1,
		Version:         22,
		ServerEncoding:  "UTF-8",
		PincodeStep:     0,
		CharBlockSize:   106,
		CharDeleteDate:  0,
		GameGuard:       false,
		LoginPacketVer:  1,
		RecvPacketKeys:  []int{},
		SendPacketKeys:  []int{},
		PacketKeys:      []int{},
		PacketObfuscate: false,
		TableFolders:    []string{},
		ServerTables:    make(map[string]string),
		CustomFields:    make(map[string]interface{}),
	}

	m.configs[name] = config
	return config
}

// GetServerConfigPath returns the path to a server configuration file
func GetServerConfigPath(configDir, serverName string) string {
	return filepath.Join(configDir, fmt.Sprintf("%s.json", serverName))
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
