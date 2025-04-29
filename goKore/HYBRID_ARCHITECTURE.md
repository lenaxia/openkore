# Hybrid Architecture Implementation Plan

## Introduction

This document outlines a plan to implement a hybrid architectural approach for the network package, combining the strengths of both factory and strategy patterns. The goal is to improve maintainability, scalability, and consistency across the codebase while handling the large number of packet types efficiently.

## Current Architecture

### Send Package
- Uses a factory pattern (`SendFactory`)
- Creates and configures `Send` implementations based on server types
- Well-structured but limited to packet construction

### Receive Package
- Uses a strategy pattern with a `CoreParser`
- Registers handlers for 400+ packet types
- Organized by domain but lacks centralized configuration

## Proposed Hybrid Architecture

The hybrid approach will:

1. Use factories to create and configure base components for different server types
2. Use the strategy pattern for registering packet handlers and constructors
3. Provide helper functions to organize related handlers by domain

### Architecture Diagram

```
┌─────────────────────┐
│   NetworkFactory    │
├─────────┬───────────┤
│SendFactory│ParserFactory│
└─────────┴───────────┘
      │           │
      ▼           ▼
┌─────────┐ ┌───────────┐
│  Send   │ │  Parser   │
└─────────┘ └───────────┘
      │           │
      ▼           ▼
┌─────────┐ ┌───────────┐
│ Packet  │ │  Packet   │
│Constructors│ Handlers  │
└─────────┘ └───────────┘
```

## Implementation Details

### 1. Parser Factory (New)

```go
// ParserCreator is a function that creates a Parser implementation
type ParserCreator func(cfg *config.ServerConfig, hookManager *hooks.HookManager) (core.Parser, error)

// ParserFactory creates Parser implementations based on server configuration
type ParserFactory struct {
    parserCreators map[string]ParserCreator
}

// NewParserFactory creates a new parser factory
func NewParserFactory() *ParserFactory {
    return &ParserFactory{
        parserCreators: make(map[string]ParserCreator),
    }
}

// RegisterParserType registers a parser creator function for a server type
func (pf *ParserFactory) RegisterParserType(serverType config.ServerType, creator ParserCreator) {
    pf.parserCreators[string(serverType)] = creator
}

// CreateParser creates a Parser implementation based on the server configuration
func (pf *ParserFactory) CreateParser(cfg *config.ServerConfig, hookManager *hooks.HookManager) (core.Parser, error) {
    creator, exists := pf.parserCreators[string(cfg.Type)]
    if !exists {
        return nil, fmt.Errorf("no parser implementation registered for server type: %s", cfg.Type)
    }
    return creator(cfg, hookManager)
}

// RegisterDefaultParserTypes registers the default parser types
func (pf *ParserFactory) RegisterDefaultParserTypes() {
    // Register ServerType0 (base implementation)
    pf.RegisterParserType(config.ServerType0, func(cfg *config.ServerConfig, hookManager *hooks.HookManager) (core.Parser, error) {
        parser := core.NewCoreParser(string(cfg.Type), hookManager)
        
        // Configure parser based on server type
        parser.SetDefaultState(0) // Default state
        
        return parser, nil
    })
    
    // Register Sakray server type
    pf.RegisterParserType(config.ServerTypeSakray, func(cfg *config.ServerConfig, hookManager *hooks.HookManager) (core.Parser, error) {
        parser := core.NewCoreParser(string(cfg.Type), hookManager)
        
        // Configure parser with Sakray-specific settings
        parser.SetDefaultState(1) // Different default state for Sakray
        
        return parser, nil
    })
}
```

### 2. Domain-Specific Handler Registration (New)

```go
// RegisterLoginHandlers registers all login-related handlers
func RegisterLoginHandlers(parser core.Parser) {
    // Register login handlers
    parser.RegisterHandlerFunc("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*",
        []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
        handleAccountServerInfo)
    
    // Register more login handlers...
}

// RegisterCharHandlers registers all character-related handlers
func RegisterCharHandlers(parser core.Parser) {
    // Register character handlers
    // ...
}

// RegisterMapHandlers registers all map-related handlers
func RegisterMapHandlers(parser core.Parser) {
    // Register map handlers
    // ...
}
```

### 3. Network Factory (New)

```go
// NetworkFactory creates and configures both send and receive components
type NetworkFactory struct {
    sendFactory   *factory.SendFactory
    parserFactory *ParserFactory
}

// NewNetworkFactory creates a new network factory
func NewNetworkFactory() *NetworkFactory {
    return &NetworkFactory{
        sendFactory:   factory.NewSendFactory(),
        parserFactory: NewParserFactory(),
    }
}

// Initialize registers default types for both factories
func (nf *NetworkFactory) Initialize() {
    nf.sendFactory.RegisterDefaultSendTypes()
    nf.parserFactory.RegisterDefaultParserTypes()
}

// CreateNetworkStack creates a complete network stack for a server
func (nf *NetworkFactory) CreateNetworkStack(cfg *config.ServerConfig, hookManager *hooks.HookManager) (*NetworkStack, error) {
    // Create send component
    send, err := nf.sendFactory.CreateSend(cfg)
    if err != nil {
        return nil, err
    }
    
    // Create parser component
    parser, err := nf.parserFactory.CreateParser(cfg, hookManager)
    if err != nil {
        return nil, err
    }
    
    // Register domain-specific handlers based on server type
    switch cfg.Type {
    case config.ServerTypeLogin:
        RegisterLoginHandlers(parser)
    case config.ServerTypeChar:
        RegisterCharHandlers(parser)
    case config.ServerTypeMap:
        RegisterMapHandlers(parser)
    }
    
    // Create and return the network stack
    return &NetworkStack{
        Send:   send,
        Parser: parser,
    }, nil
}
```

### 4. Network Stack (New)

```go
// NetworkStack represents a complete network stack for a server
type NetworkStack struct {
    Send   core.Send
    Parser core.Parser
}
```

## Migration Strategy

### Phase 1: Preparation (Week 1-2)

1. **Create New Components**
   - Implement `ParserFactory`
   - Implement domain-specific handler registration functions
   - Implement `NetworkFactory`
   - Create unit tests for new components

2. **Document API**
   - Document the new API and usage patterns
   - Create examples for common use cases

### Phase 2: Integration (Week 3-4)

1. **Create Adapter Layer**
   - Implement adapters to make new components work with existing code
   - Ensure backward compatibility

2. **Update Test Infrastructure**
   - Update test infrastructure to support both old and new approaches
   - Create integration tests for the hybrid approach

### Phase 3: Migration (Week 5-8)

1. **Migrate Core Components**
   - Update core components to use the new factories
   - Migrate handler registration to domain-specific functions

2. **Migrate Client Code**
   - Update client code to use the `NetworkFactory`
   - Remove direct usage of old components

3. **Validate**
   - Run all tests to ensure functionality is preserved
   - Perform manual testing of key scenarios

### Phase 4: Cleanup (Week 9-10)

1. **Remove Deprecated Code**
   - Remove adapter layer
   - Remove deprecated components and functions

2. **Finalize Documentation**
   - Update all documentation to reflect the new architecture
   - Create migration guides for any remaining legacy code

## Testing Approach

### Unit Tests

1. **Factory Tests**
   - Test creation of parsers for different server types
   - Test configuration of parsers

2. **Handler Registration Tests**
   - Test registration of handlers
   - Test handler invocation

3. **Integration Tests**
   - Test end-to-end packet handling
   - Test interaction between send and receive components

### Performance Tests

1. **Memory Usage**
   - Compare memory usage between old and new approaches
   - Ensure no significant increase in memory usage

2. **Processing Speed**
   - Compare processing speed for packet handling
   - Ensure no significant decrease in performance

## Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Breaking changes | High | Medium | Create adapter layer, thorough testing |
| Performance regression | High | Low | Performance testing, optimization |
| Increased complexity | Medium | Medium | Clear documentation, examples |
| Migration delays | Medium | Medium | Phased approach, prioritize critical components |
| Team resistance | Medium | Low | Clear communication, training |

## Benefits

1. **Improved Maintainability**
   - Centralized configuration
   - Organized handler registration
   - Consistent interfaces

2. **Better Scalability**
   - Factory pattern for server types
   - Strategy pattern for packet handlers
   - Easier to add new server types and packet handlers

3. **Enhanced Consistency**
   - Consistent approach across send and receive packages
   - Clear separation of concerns
   - Better organization of code

4. **Easier Testing**
   - More modular components
   - Clearer interfaces
   - Better separation of concerns

## Conclusion

The hybrid approach combines the strengths of both factory and strategy patterns, providing a more maintainable and scalable architecture for the network package. By centralizing configuration while maintaining flexibility for packet handlers, we can better manage the complexity of supporting multiple server types and hundreds of packet types.

The phased migration approach ensures minimal disruption to existing code while gradually moving towards the new architecture. With proper testing and documentation, the transition should be smooth and result in a more maintainable codebase.