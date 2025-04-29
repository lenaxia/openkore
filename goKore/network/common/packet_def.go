// Package common provides shared types and utilities for both send and receive components.
package common

// PacketDef defines the structure of a packet for receiving
type PacketDef struct {
	// ID is the packet ID.
	ID string

	// Name is the name of the packet.
	Name string

	// Format is the format of the packet.
	Format string

	// FieldNames are the names of the fields in the packet.
	FieldNames []string
}

// PacketDefProvider provides packet definitions for a specific server type
type PacketDefProvider func() map[string]PacketDef

// PacketConstruction defines how to construct a packet for sending
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

// PacketConstructionProvider provides packet constructions for a specific server type
type PacketConstructionProvider func() map[string]PacketConstruction

// PacketHandler is a function that processes a parsed packet
type PacketHandler func(args map[string]interface{}) error
