# Unified Handler Registration Pattern for Send

This directory contains a template for implementing the unified handler registration pattern in the send side of goKore. This pattern standardizes how packet handlers are registered across different packages in the codebase.

## Overview

The unified registration pattern consists of:

1. A `register.go` file that exposes standard registration functions
2. A package-specific implementation file (e.g., `template.go`) that contains the actual handler logic

## Standard Registration Functions

Each package should implement these standard registration functions:

### RegisterWithSend

```go
// RegisterWithSend registers all handlers with the send component
func RegisterWithSend(send *core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) {
    // Create the manager
    manager := NewManager(send, hookManager, logger)
    
    // Register handlers
    manager.RegisterHandlers()
}
```

This function is used to register handlers with the send component. It creates a manager and calls its `RegisterHandlers` method.

### GetPacketDefinitions

```go
// GetPacketDefinitions returns the packet definitions for this package
func GetPacketDefinitions() map[string]common.PacketConstruction {
    return map[string]common.PacketConstruction{
        "0000": {
            Name:       "example_handler",
            Format:     "v",
            FieldNames: []string{"example_field"},
        },
    }
}
```

This function returns the packet definitions for the package.

## Manager Implementation

Each package should implement a manager that handles the actual packet processing:

```go
// Manager handles template-related packets
type Manager struct {
    send        *core.BaseSend
    hookManager *hooks.HookManager
    logger      core.Logger
}

// NewManager creates a new template manager
func NewManager(send *core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) *Manager {
    return &Manager{
        send:        send,
        hookManager: hookManager,
        logger:      logger,
    }
}

// RegisterHandlers registers all template-related packet handlers
func (m *Manager) RegisterHandlers() {
    // Register example handler
    if m.send != nil {
        m.send.RegisterHandler("example_handler", m.handleExample)
    }
}

// handleExample handles the example packet
func (m *Manager) handleExample(args map[string]interface{}) ([]byte, error) {
    // Implementation
    return []byte{0x00, 0x00, 0x01}, nil
}
```

## How to Use This Template

1. Copy the `template` directory to create a new package
2. Rename `template.go` to match your package name
3. Update the package name in all files
4. Implement your packet handlers in the manager
5. Update the registration functions to register your handlers

## Benefits of the Unified Pattern

- Consistent interface across all packages
- Clear separation of registration and implementation
- Easy to add new packet handlers
- Simplified integration with the registry
- Better logging support

## Integration with HandlerRegistry

The `HandlerRegistry` in `goKore/network/send/registry.go` provides a central place for registering all send handlers. It uses the unified registration pattern to register handlers from different packages.

To integrate your package with the registry, update the `registerGameHandlers` method in the registry to call your package's `RegisterWithSend` function:

```go
func (hr *HandlerRegistry) registerGameHandlers() {
    // Register game handlers
    game.RegisterHandlers(hr.baseSend)
    
    // Register your package handlers
    yourpackage.RegisterWithSend(hr.baseSend, hr.hookManager, hr.logger)
    
    // Log registration
    hr.logger.Debug("Registered game send handlers")
}