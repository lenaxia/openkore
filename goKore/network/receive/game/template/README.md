# Unified Handler Registration Pattern

This directory contains a template for implementing the unified handler registration pattern in goKore. This pattern standardizes how packet handlers are registered across different packages in the codebase.

## Overview

The unified registration pattern consists of:

1. A `register.go` file that exposes standard registration functions
2. A package-specific implementation file (e.g., `template.go`) that contains the actual handler logic

## Standard Registration Functions

Each package should implement these standard registration functions:

### RegisterWithParser

```go
// RegisterWithParser registers all handlers with the parser
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) {
    // Create the manager
    manager := NewManager(parser, hookManager, logger)
    
    // Register handlers
    manager.RegisterHandlers()
}
```

This function is used to register handlers with the core parser. It creates a manager and calls its `RegisterHandlers` method.

### RegisterWithReceive

```go
// RegisterWithReceive registers all handlers with the receive interface
func RegisterWithReceive(receive types.Receive) {
    // Register handlers directly with the receive interface
    receive.RegisterHandler("example_handler", func(args map[string]interface{}) error {
        // Create a manager for this specific call
        manager := NewManager(nil, nil, nil)
        return manager.handleExample(args)
    })
}
```

This function is used to register handlers with the receive interface. It registers handlers directly with the receive interface.

### GetPacketDefinitions

```go
// GetPacketDefinitions returns the packet definitions for this package
func GetPacketDefinitions() map[string]common.PacketConstruction {
    return map[string]common.PacketConstruction{
        "0000": {
            ID:         "0000",
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
    parser      *core.CoreParser
    hookManager *hooks.HookManager
    logger      core.Logger
}

// NewManager creates a new template manager
func NewManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *Manager {
    return &Manager{
        parser:      parser,
        hookManager: hookManager,
        logger:      logger,
    }
}

// RegisterHandlers registers all template-related packet handlers
func (m *Manager) RegisterHandlers() {
    // Register example handler
    if m.parser != nil {
        m.parser.RegisterHandlerFunc("0000", "example_handler", "v",
            []string{"example_field"},
            m.handleExample)
    }
}

// handleExample handles the example packet
func (m *Manager) handleExample(args map[string]interface{}) error {
    // Implementation
    return nil
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