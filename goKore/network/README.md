# Network Package

This package provides the core networking functionality for OpenKore. It defines connection states, interfaces, and error types used throughout the network stack.

## Structure

- `network.go`: Core interfaces and constants
- `common/`: Common utilities and shared functionality
- `config/`: Network configuration
- `connection/`: Connection management
- `hooks/`: Event hooks system
- `packets/`: Packet definitions and structures
- `protocol/`: Protocol implementations
- `proxy/`: Proxy functionality
- `receive/`: Packet receiving and handling
- `send/`: Packet construction and sending
- `servers/`: Server-specific implementations
- `state/`: Connection state management

## Future Improvements

### Domain-Organized Packet Handling

The current packet handling implementation in the `receive/` package uses a direct approach where handlers update state and often return nil. This approach works but has several limitations:

1. **Limited Decoupling**: Handlers are tightly coupled to their state objects
2. **Difficult Testing**: Hard to test handlers in isolation
3. **No Reaction Mechanism**: Other components can't easily react to packet events
4. **Redundant Code**: Similar patterns repeated across many handlers

#### Proposed Refactoring

A domain-organized approach that leverages the existing hooks system would provide better organization while maintaining Go's preference for simplicity:

```go
// Domain-specific manager for character behavior packets
type CharacterBehaviorManager struct {
    parser      *CoreParser
    hookManager *hooks.HookManager
}

// Create a new character behavior manager
func NewCharacterBehaviorManager(parser *CoreParser, hookManager *hooks.HookManager) *CharacterBehaviorManager {
    return &CharacterBehaviorManager{
        parser:      parser,
        hookManager: hookManager,
    }
}

// Register all handlers related to character behavior
func (m *CharacterBehaviorManager) RegisterHandlers() {
    // Register manner_message handler
    m.parser.RegisterHandlerFunc("0149", "manner_message", "C",
        []string{"flag"},
        m.handleMannerMessage)
        
    // Register hack_shield_alarm handler
    m.parser.RegisterHandlerFunc("08B3", "hack_shield_alarm", "",
        []string{},
        m.handleHackShieldAlarm)
}

// Handle manner_message packet
func (m *CharacterBehaviorManager) handleMannerMessage(args map[string]interface{}) error {
    // Process the packet
    result := m.processMannerMessage(args)
    
    // Notify through hooks system
    if m.hookManager != nil {
        m.hookManager.CallHook("character.manner_message", result)
    }
    
    return nil
}

// Process manner_message packet and return structured result
func (m *CharacterBehaviorManager) processMannerMessage(args map[string]interface{}) map[string]interface{} {
    // Process packet data and return structured result
    // ...
}
```

#### Benefits

1. **Domain-Driven Organization**: Group related packets by domain (character behavior, inventory, etc.)
2. **Improved Testability**: Domain-specific managers can be tested in isolation
3. **Consistent Hook Notifications**: Standard pattern for notifying other components
4. **Reduced Duplication**: Common patterns extracted to domain managers
5. **Simplicity**: Leverages existing hooks system without introducing new abstractions

#### Implementation Plan

1. Identify logical domains for packet handlers (e.g., character behavior, inventory, chat)
2. Create domain-specific manager types for each domain
3. Move related handlers to their respective managers
4. Update tests to use the new organization
5. Gradually migrate all packet handlers

This refactoring would make the codebase more maintainable and extensible while preserving the existing functionality and adhering to Go's philosophy of simplicity.