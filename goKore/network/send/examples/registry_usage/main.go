// Package main demonstrates how to use the send registry
package main

import (
	"fmt"
	"os"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send"
	"github.com/lenaxia/goKore/network/send/core"
)

// SimpleLogger is a simple implementation of the core.Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *SimpleLogger) Warning(format string, args ...interface{}) {
	fmt.Printf("[WARNING] "+format+"\n", args...)
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

func (l *SimpleLogger) Success(format string, args ...interface{}) {
	fmt.Printf("[SUCCESS] "+format+"\n", args...)
}

// MockConnection is a mock implementation of a connection
type MockConnection struct{}

func (mc *MockConnection) Send(data []byte) error {
	fmt.Printf("Sending packet: %v\n", data)
	return nil
}

func main() {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a logger
	logger := &SimpleLogger{}

	// Create a base send instance
	baseSend := core.NewBaseSend(hookManager)

	// Set debug mode
	baseSend.SetDebugMode(true)

	// Create a registry
	registry := send.NewHandlerRegistry(baseSend, hookManager, logger)

	// Register all handlers
	registry.RegisterAllHandlers()

	// Get packet definitions
	packetDefs := registry.GetPacketDefinitions()

	// Configure server type
	err := registry.ConfigureServerType("ServerType0", packetDefs)
	if err != nil {
		logger.Error("Failed to configure server type: %v", err)
		os.Exit(1)
	}

	// Set a mock connection
	baseSend.SetConnection(&MockConnection{})

	// Example: Send a card merge request
	err = baseSend.SendPacket("card_merge_request", map[string]interface{}{
		"cardID": uint32(12345),
	})
	if err != nil {
		logger.Error("Failed to send card_merge_request: %v", err)
		os.Exit(1)
	}

	// Example: Send a card merge
	err = baseSend.SendPacket("card_merge", map[string]interface{}{
		"cardID": uint32(12345),
		"itemID": uint32(67890),
	})
	if err != nil {
		logger.Error("Failed to send card_merge: %v", err)
		os.Exit(1)
	}

	logger.Success("Successfully sent packets")
}
