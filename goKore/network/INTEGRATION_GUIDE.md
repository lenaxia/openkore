# Network Stack Integration Guide

This guide explains how to properly integrate the send and receive stacks with the NetworkManager in the goKore network package.

## Overview

The network stack in goKore consists of three main components:

1. **NetworkManager** (`network.go`) - The central coordinator that connects the send and receive stacks
2. **Send Stack** (`network/send`) - Responsible for constructing and sending packets to the server
3. **Receive Stack** (`network/receive`) - Responsible for parsing incoming packets and triggering appropriate handlers

## Component Interfaces

### NetworkManager

The NetworkManager requires implementations of the following interfaces:

- `NetworkInterface` - Handles the actual network connection
- `PacketSender` - Constructs and sends packets
- `PacketHandler` - Processes received packets

```go
// NetworkManager manages the network connections and packet handling
type NetworkManager struct {
    // Current connection state
    state int

    // Callback for state changes
    stateChangeCallback StateChangeCallback

    // Network interface
    networkInterface NetworkInterface

    // Packet sender
    packetSender PacketSender

    // Packet handler
    packetHandler PacketHandler
}
```

### Send Stack

The send stack implements the `PacketSender` interface:

```go
// PacketSender defines the interface for packet sending
type PacketSender interface {
    // Send constructs and sends a packet with the given name and fields
    Send(packetName string, fields map[string]interface{}) ([]byte, error)

    // GetCashShopManager returns a manager for cash shop-related operations
    GetCashShopManager() interface{}

    // GetMiscManager returns a manager for miscellaneous operations
    GetMiscManager() interface{}

    // GetInfoChatManager returns a manager for chat info operations
    GetInfoChatManager() interface{}
}
```

### Receive Stack

The receive stack implements the `PacketHandler` interface:

```go
// PacketHandler defines the interface for packet handling
type PacketHandler interface {
    // Handle processes a packet and returns an error if processing fails
    Handle(packet []byte) error
}
```

## Integration Steps

### 1. Create the Hook Manager

The hook manager is shared between the send and receive stacks to allow for event-based communication:

```go
import "github.com/lenaxia/goKore/network/hooks"

// Create a hook manager
hookManager := hooks.NewHookManager()
```

### 2. Set Up the Send Stack

```go
import (
    "github.com/lenaxia/goKore/network/common"
    "github.com/lenaxia/goKore/network/send"
)

// Create a base send instance
baseSend := send.NewBaseSend(hookManager)

// Configure with packet constructions
packetConstructions := map[string]common.PacketConstruction{
    "0064": {
        ID:         "0064",
        Name:       "login_request",
        Format:     "v a24 a24 C",
        FieldNames: []string{"version", "username", "password", "clienttype"},
    },
    // Add more packet constructions as needed
}
baseSend.Configure("YourServerType", packetConstructions)

// Register custom handlers for specific packets
baseSend.RegisterHandler("login_request", func(args map[string]interface{}) ([]byte, error) {
    // Custom logic to construct login_request packet
    return baseSend.Reconstruct("0064", args)
})
```

### 3. Set Up the Receive Stack

```go
import (
    "github.com/lenaxia/goKore/network/common"
    "github.com/lenaxia/goKore/network/receive/base"
)

// Create a base receive instance
baseReceive := base.NewBaseReceive(hookManager)

// Configure with packet definitions
packetDefs := map[string]common.PacketDef{
    "0073": {
        Name:       "server_connected",
        Format:     "C a4 a4 v C",
        FieldNames: []string{"result", "sessionID", "accountID", "sessionID2", "sex"},
    },
    // Add more packet definitions as needed
}
baseReceive.Configure("YourServerType", packetDefs)

// Register handlers for specific packets
baseReceive.RegisterHandler("server_connected", func(args map[string]interface{}) error {
    // Handle the server_connected packet
    fmt.Printf("Connected to server with account ID: %v\n", args["accountID"])
    return nil
})
```

### 4. Create Adapter Classes

Since the base implementations don't directly implement the required interfaces, we need to create adapter classes:

```go
// SendAdapter adapts the BaseSend to implement the network.PacketSender interface
type SendAdapter struct {
    *send.BaseSend
}

func NewSendAdapter(baseSend *send.BaseSend) *SendAdapter {
    return &SendAdapter{BaseSend: baseSend}
}

// Send implements the network.PacketSender interface
func (sa *SendAdapter) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
    // Construct the packet
    packet, err := sa.BaseSend.ConstructPacket(packetName, fields)
    if err != nil {
        return nil, err
    }
    
    // Send the packet
    err = sa.BaseSend.SendToServer(packet)
    return packet, err
}

// GetCashShopManager implements the network.PacketSender interface
func (sa *SendAdapter) GetCashShopManager() interface{} {
    return nil
}

// GetMiscManager implements the network.PacketSender interface
func (sa *SendAdapter) GetMiscManager() interface{} {
    return nil
}

// GetInfoChatManager implements the network.PacketSender interface
func (sa *SendAdapter) GetInfoChatManager() interface{} {
    return nil
}

// ReceiveAdapter adapts the BaseReceive to implement the network.PacketHandler interface
type ReceiveAdapter struct {
    *base.BaseReceive
}

func NewReceiveAdapter(baseReceive *base.BaseReceive) *ReceiveAdapter {
    return &ReceiveAdapter{BaseReceive: baseReceive}
}

// Handle implements the network.PacketHandler interface
func (ra *ReceiveAdapter) Handle(packet []byte) error {
    return ra.Process(packet)
}
```

### 5. Create the Network Interface

Implement the NetworkInterface to handle the actual network connection:

```go
import "github.com/lenaxia/goKore/network/connection"

// Create a direct connection
config := &connection.ConnectionConfig{
    Host:        "your.server.com",
    Port:        6900,
    Timeout:     5 * time.Second,
    RecvTimeout: 5 * time.Second,
    SendTimeout: 5 * time.Second,
    ServerType:  "YourServerType",
}

conn := connection.NewDirectConnection(config)

// Create a connection manager
networkInterface := connection.NewConnectionManager(conn)
```

### 6. Create the NetworkManager

Finally, create the NetworkManager with all the components:

```go
import "github.com/lenaxia/goKore/network"

// Create adapters
sendAdapter := NewSendAdapter(baseSend)
receiveAdapter := NewReceiveAdapter(baseReceive)

// Create the network manager
networkManager := network.NewNetworkManager(networkInterface, sendAdapter, receiveAdapter)

// Set a state change callback if needed
networkManager.SetStateChangeCallback(func(oldState, newState int) {
    fmt.Printf("Connection state changed from %d to %d\n", oldState, newState)
})
```

### 7. Use the NetworkManager

Now you can use the NetworkManager to connect, send packets, and handle received packets:

```go
// Connect to the server
err := networkManager.Connect()
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Send a packet
_, err = networkManager.Send("login_request", map[string]interface{}{
    "version":    15,
    "username":   "user",
    "password":   "pass",
    "clienttype": 0,
})
if err != nil {
    log.Fatalf("Failed to send login request: %v", err)
}

// When you receive a packet from the network interface, process it
// This would typically be in a receive loop
for {
    receivedPacket, err := networkInterface.Receive()
    if err != nil {
        if err == network.ErrTimeout {
            continue
        }
        log.Fatalf("Failed to receive packet: %v", err)
    }

    err = networkManager.HandlePacket(receivedPacket)
    if err != nil {
        log.Printf("Failed to handle packet: %v", err)
    }
}
```

## Complete Example

For a complete example of how to integrate the send and receive stacks with the NetworkManager, see the `send_receive_integration_test.go` file in the network package.

## Best Practices

1. **Share the Hook Manager**: Use a single hook manager for both the send and receive stacks to ensure proper event handling.

2. **Register Handlers Early**: Register packet handlers before connecting to ensure all packets can be properly processed.

3. **Use State Management**: Leverage the state management capabilities of the NetworkManager to track connection status.

4. **Error Handling**: Always check for errors when sending packets or handling received packets.

5. **Packet Definitions**: Keep packet definitions up-to-date with the server protocol to ensure proper packet construction and parsing.

6. **Testing**: Write comprehensive tests for your network integration to ensure reliability.