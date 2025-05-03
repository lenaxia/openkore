// Package main provides an example of how to register all handlers using the HandlerRegistry.
package main

import (
	"fmt"
	"os"

	"github.com/lenaxia/goKore/network/hooks"
	receivepkg "github.com/lenaxia/goKore/network/receive"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
)

// SimpleLogger implements the core.Logger interface
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

func main() {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a logger
	logger := &SimpleLogger{}

	// Create a receive factory
	receiveFactory := factory.NewReceiveFactory()

	// Register default server types
	receiveFactory.RegisterDefaultServerTypes()

	// Create a receive component for ServerType0
	receive, err := receiveFactory.CreateReceive("ServerType0", hookManager)
	if err != nil {
		logger.Error("Error creating receive component: %v", err)
		os.Exit(1)
	}

	// Get the core parser
	// In a real implementation, you would have access to the CoreParser directly
	// For this example, we'll create a new one
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create a handler registry
	registry := receivepkg.NewHandlerRegistry(parser, hookManager, receive, logger)

	// Register all handlers
	logger.Info("Registering all handlers...")
	registry.RegisterAllHandlers()
	logger.Success("All handlers registered successfully!")

	// Alternatively, you can register handlers by type
	logger.Info("Registering handlers with Receive interface...")
	registry.RegisterWithReceive()
	logger.Success("Receive handlers registered successfully!")

	logger.Info("Registering handlers with CoreParser...")
	registry.RegisterWithParser()
	logger.Success("CoreParser handlers registered successfully!")

	// Example packet data (account_server_info packet)
	// This is a simplified example - in a real implementation, this would come from the network
	packet := []byte{
		0x69, 0x00, // Packet ID (0069)
		0x4C, 0x00, // Length (76 bytes)
		0x01, 0x02, 0x03, 0x04, // Session ID
		0x05, 0x06, 0x07, 0x08, // Account ID
		0x09, 0x0A, 0x0B, 0x0C, // Session ID 2
		0x0D, 0x0E, 0x0F, 0x10, // Last Login IP
		// Last Login Time (26 bytes)
		0x32, 0x30, 0x32, 0x35, 0x2D, 0x30, 0x34, 0x2D, 0x32, 0x37, 0x20, 0x32, 0x30, 0x3A, 0x31, 0x34, 0x3A, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Account Sex (1 = male)
		// Server Info (remaining bytes)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
		0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E,
	}

	// Process the packet
	logger.Info("Processing account_server_info packet...")
	err = receive.Process(packet)
	if err != nil {
		logger.Error("Error processing packet: %v", err)
		os.Exit(1)
	}
	logger.Success("Packet processed successfully!")
}
