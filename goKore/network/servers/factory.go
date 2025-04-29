// Package servers provides server implementations for different Ragnarok Online server types.
package servers

// ServerFactory creates server instances based on server type
type ServerFactory struct {
	// Registry of server creators by type
	creators map[ServerType]func(*BaseServerConfig) (Server, error)
}

// NewServerFactory creates a new server factory
func NewServerFactory() *ServerFactory {
	factory := &ServerFactory{
		creators: make(map[ServerType]func(*BaseServerConfig) (Server, error)),
	}

	// Register default server creators
	factory.RegisterServerCreator(ServerTypeOfficial, createBaseServer)
	factory.RegisterServerCreator(ServerTypeSakray, createSakrayServer)
	// Add more server types as needed

	return factory
}

// RegisterServerCreator registers a creator function for a specific server type
func (f *ServerFactory) RegisterServerCreator(serverType ServerType, creator func(*BaseServerConfig) (Server, error)) {
	f.creators[serverType] = creator
}

// CreateServer creates a server instance based on the server type
func (f *ServerFactory) CreateServer(baseConfig *BaseServerConfig) (Server, error) {
	// If server type is not specified, detect it
	if baseConfig.Type == ServerTypeUnknown {
		baseConfig.Type = DetectServerType(baseConfig.ServerConfig)
	}

	// Look up the creator function
	creator, exists := f.creators[baseConfig.Type]
	if !exists {
		// Fall back to base server if no specific creator is registered
		return createBaseServer(baseConfig)
	}

	// Create the server
	return creator(baseConfig)
}

// createBaseServer creates a base server instance
func createBaseServer(baseConfig *BaseServerConfig) (Server, error) {
	return NewBaseServer(baseConfig)
}

// createSakrayServer creates a Sakray server instance
func createSakrayServer(baseConfig *BaseServerConfig) (Server, error) {
	return NewSakrayServer(baseConfig)
}
