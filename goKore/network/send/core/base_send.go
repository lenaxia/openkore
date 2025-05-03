// Package core provides core functionality for the send component.
package core

// SendHandler is a function that handles a send request.
type SendHandler func(args map[string]interface{}) ([]byte, error)

// BaseSend is the interface for the base send functionality.
type BaseSend interface {
	// RegisterHandler registers a handler for a packet type.
	RegisterHandler(packetName string, handler SendHandler)

	// ConstructPacket constructs a packet for the given packet name and arguments.
	ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error)

	// GetServerType returns the server type.
	GetServerType() string
}

// Logger is the interface for logging.
type Logger interface {
	// Debug logs a debug message.
	Debug(format string, args ...interface{})

	// Info logs an info message.
	Info(format string, args ...interface{})

	// Warning logs a warning message.
	Warning(format string, args ...interface{})

	// Error logs an error message.
	Error(format string, args ...interface{})
}
