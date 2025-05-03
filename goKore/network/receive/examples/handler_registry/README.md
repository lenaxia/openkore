# Handler Registry Example

This example demonstrates how to use the `HandlerRegistry` to register all packet handlers in the goKore network receive system.

## Overview

The `HandlerRegistry` provides a centralized way to register all packet handlers in the goKore network receive system. It supports registering handlers with both the `Receive` interface and the `CoreParser` directly.

## Key Components

1. **HandlerRegistry**: Manages the registration of all packet handlers
2. **CoreParser**: Parses and processes packets
3. **HookManager**: Manages hooks for packet processing
4. **Logger**: Provides logging functionality

## Usage

### Creating a Handler Registry

```go
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
parser := core.NewCoreParser("ServerType0", hookManager)

// Create a handler registry
registry := receivepkg.NewHandlerRegistry(parser, hookManager, receive, logger)
```

### Registering All Handlers

```go
// Register all handlers
registry.RegisterAllHandlers()
```

### Registering Handlers by Type

```go
// Register handlers with the Receive interface
registry.RegisterWithReceive()

// Register handlers with the CoreParser
registry.RegisterWithParser()
```

## Running the Example

To run this example:

```bash
cd goKore/network/receive/examples/handler_registry
go run main.go
```

## Expected Output

The example will:

1. Register all handlers
2. Process a sample account_server_info packet
3. Log the results

You should see output similar to:

```
[INFO] Registering all handlers...
[SUCCESS] All handlers registered successfully!
[INFO] Registering handlers with Receive interface...
[SUCCESS] Receive handlers registered successfully!
[INFO] Registering handlers with CoreParser...
[SUCCESS] CoreParser handlers registered successfully!
[INFO] Processing account_server_info packet...
[SUCCESS] Packet processed successfully!
```

## Extending the Handler Registry

To add new handlers to the registry:

1. Create a new manager for your handlers
2. Implement the `RegisterHandlers()` method
3. Add the manager to the appropriate registration method in the `HandlerRegistry`