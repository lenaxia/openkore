//go:build ignore
// +build ignore

// This file demonstrates how to use the registry adapter in a real application.
// It is not meant to be compiled as part of the package, but rather serves as documentation.
package main

import (
	"fmt"
	"log"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/registry"
)

// SimpleLogger implements the registry.Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *SimpleLogger) Warning(format string, args ...interface{}) {
	log.Printf("[WARNING] "+format, args...)
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (l *SimpleLogger) Success(format string, args ...interface{}) {
	log.Printf("[SUCCESS] "+format, args...)
}

// SimpleNetworkInterface implements the network.NetworkInterface interface
type SimpleNetworkInterface struct {
	connected bool
	state     int
}

func (s *SimpleNetworkInterface) Connect() error {
	s.connected = true
	s.state = network.ConnectedToMasterServer
	return nil
}

func (s *SimpleNetworkInterface) Disconnect() error {
	s.connected = false
	s.state = network.NotConnected
	return nil
}

func (s *SimpleNetworkInterface) IsConnected() bool {
	return s.connected
}

func (s *SimpleNetworkInterface) GetState() int {
	return s.state
}

func (s *SimpleNetworkInterface) SetState(state int) {
	s.state = state
}

func (s *SimpleNetworkInterface) Send(data []byte) error {
	fmt.Printf("Sending data: %v\n", data)
	return nil
}

func (s *SimpleNetworkInterface) Receive() ([]byte, error) {
	return []byte{}, nil
}

func main() {
	// Create a logger
	logger := &SimpleLogger{}

	// Create a network integrator
	integrator := registry.NewNetworkRegistryIntegrator(logger)

	// Create a network interface
	networkInterface := &SimpleNetworkInterface{}

	// Create a network manager
	manager := integrator.CreateNetworkManager(networkInterface)

	// Connect to the server
	err := manager.Connect()
	if err != nil {
		logger.Error("Failed to connect: %v", err)
		return
	}

	// Send a ping packet
	_, err = manager.Send("ping", map[string]interface{}{})
	if err != nil {
		logger.Error("Failed to send ping: %v", err)
	}

	// Process a received packet
	packet := []byte{0x01, 0x02, 0x03} // Example packet
	err = manager.HandlePacket(packet)
	if err != nil {
		logger.Error("Failed to handle packet: %v", err)
	}

	// Disconnect from the server
	err = manager.Disconnect()
	if err != nil {
		logger.Error("Failed to disconnect: %v", err)
	}
}
