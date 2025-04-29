// Package types provides type definitions for the send component.
package types

// SendHandler is a function that constructs a packet
type SendHandler func(args map[string]interface{}) ([]byte, error)

// PacketConstruction defines how to construct a packet
type PacketConstruction struct {
	// ID is the packet ID.
	ID string

	// Name is the name of the packet.
	Name string

	// Format is the format of the packet.
	Format string

	// FieldNames are the names of the fields in the packet.
	FieldNames []string
}

// Send interface defines the methods that a send implementation must provide
type Send interface {
	// RegisterHandler registers a handler for a specific packet
	RegisterHandler(packetName string, handler SendHandler)

	// ConstructPacket constructs a packet from a packet name and arguments
	ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error)

	// SendPacket constructs and sends a packet
	SendPacket(packetName string, args map[string]interface{}) error

	// SendToServer sends a raw packet to the server
	SendToServer(packet []byte) error

	// Configure configures the send component with server-specific packet constructions
	Configure(serverType string, packetConstructions map[string]PacketConstruction) error
}
