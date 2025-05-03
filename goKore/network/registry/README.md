# Network Registry Adapter

This package provides adapters to integrate the send and receive registries with the gokore network stack.

## Overview

The registry adapter package solves the problem of integrating the send and receive registries with the network stack while avoiding import cycles. It provides adapters that implement the network interfaces required by the NetworkManager.

## Components

### SendRegistryAdapter

The `SendRegistryAdapter` adapts the `send.HandlerRegistry` to implement the `network.PacketSender` interface. It provides methods for:

- Sending packets with specific names and fields
- Getting various managers (CashShop, Misc, InfoChat)
- Initializing the send registry
- Setting the connection for sending packets

### ReceiveRegistryAdapter

The `ReceiveRegistryAdapter` adapts the `receive.HandlerRegistry` to implement the `network.PacketHandler` interface. It provides methods for:

- Handling received packets
- Initializing the receive registry

### NetworkRegistryIntegrator

The `NetworkRegistryIntegrator` integrates the send and receive registries with the network stack. It provides methods for:

- Creating a network manager with the integrated registries
- Getting the hook manager used by the integrator

## Usage

Here's a simple example of how to use the registry adapter:

```go
// Create a logger
logger := &YourLoggerImplementation{}

// Create a network integrator
integrator := registry.NewNetworkRegistryIntegrator(logger)

// Create a network interface
networkInterface := &YourNetworkInterfaceImplementation{}

// Create a network manager
manager := integrator.CreateNetworkManager(networkInterface)

// Connect to the server
err := manager.Connect()
if err != nil {
    logger.Error("Failed to connect: %v", err)
    return
}

// Send a packet
_, err = manager.Send("ping", map[string]interface{}{})
if err != nil {
    logger.Error("Failed to send ping: %v", err)
}

// Process a received packet
packet := []byte{...} // Your packet data
err = manager.HandlePacket(packet)
if err != nil {
    logger.Error("Failed to handle packet: %v", err)
}

// Disconnect from the server
err = manager.Disconnect()
if err != nil {
    logger.Error("Failed to disconnect: %v", err)
}
```

See `example_usage.go` for a more detailed example.

## Testing

The package includes tests for all components. Run the tests with:

```
go test -v
```

## Logger Interface

The package requires a logger that implements the following interface:

```go
type Logger interface {
    Debug(format string, args ...interface{})
    Info(format string, args ...interface{})
    Warning(format string, args ...interface{})
    Error(format string, args ...interface{})
    Success(format string, args ...interface{})
}
```

This logger is used by the adapters to log important events and errors.