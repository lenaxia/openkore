// Package factory provides factories for creating network components based on server type.
package factory

import (
	"fmt"
	"os"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/protocol"
)

// ServerTypeFactory creates network components based on server type
type ServerTypeFactory struct {
	// Add any factory-specific fields here
}

// NewServerTypeFactory creates a new server type factory
func NewServerTypeFactory() *ServerTypeFactory {
	return &ServerTypeFactory{}
}

// LoadPacketDefinitions loads packet definitions based on server type and table folders
func (f *ServerTypeFactory) LoadPacketDefinitions(basePath string, serverConfig *config.ServerConfig) (map[string]protocol.PacketLengthDef, error) {
	// Check if any of the table folders exist
	if len(serverConfig.TableFolders) > 0 {
		folderExists := false
		for _, folder := range serverConfig.TableFolders {
			path := fmt.Sprintf("%s/%s", basePath, folder)
			if _, err := os.Stat(path); err == nil {
				folderExists = true
				break
			}
		}

		if !folderExists {
			return nil, fmt.Errorf("no valid table folders found")
		}
	}

	// Load packet definitions
	packetDefs, err := protocol.LoadRecvPackets(basePath, serverConfig.TableFolders)
	if err != nil {
		return nil, fmt.Errorf("failed to load packet definitions: %w", err)
	}

	return packetDefs, nil
}

// CreateTokenizer creates a tokenizer with packet definitions for the specified server type
func (f *ServerTypeFactory) CreateTokenizer(basePath string, serverConfig *config.ServerConfig) (*protocol.Tokenizer, error) {
	// Load packet definitions
	packetDefs, err := f.LoadPacketDefinitions(basePath, serverConfig)
	if err != nil {
		return nil, err
	}

	// Create tokenizer
	tokenizer := protocol.NewTokenizer(packetDefs)

	return tokenizer, nil
}
