// Package protocol provides functionality for handling the Ragnarok Online network protocol.
// This file implements the packet parser which interprets messages sent by the RO server.
package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// HandlerFunc is a function that processes a packet
type HandlerFunc func(packet []byte) error

// PacketInfo contains information about a packet's structure and handler
type PacketInfo struct {
	Name       string      // Name of the handler function
	Format     string      // Format string for unpacking the packet
	ParamNames []string    // Names of the parameters in the packet
	Handler    HandlerFunc // Function to handle the packet
}

// PacketParser is responsible for parsing and constructing packets
type PacketParser struct {
	PacketList map[string]PacketInfo // Maps packet IDs to packet info
	PacketLUT  map[string]string     // Maps handler names to packet IDs
	HookPrefix string                // Prefix for hooks
}

// NewPacketParser creates a new packet parser
func NewPacketParser() *PacketParser {
	return &PacketParser{
		PacketList: make(map[string]PacketInfo),
		PacketLUT:  make(map[string]string),
	}
}

// RegisterHandler registers a handler for a specific packet ID
func (p *PacketParser) RegisterHandler(packetID, name, format string, paramNames []string, handler HandlerFunc) {
	p.PacketList[packetID] = PacketInfo{
		Name:       name,
		Format:     format,
		ParamNames: paramNames,
		Handler:    handler,
	}
	p.PacketLUT[name] = packetID
}

// LookupPacketID looks up a packet ID by handler name
func (p *PacketParser) LookupPacketID(handlerName string) (string, bool) {
	packetID, exists := p.PacketLUT[handlerName]
	return packetID, exists
}

// Parse parses a packet into a map of arguments
func (p *PacketParser) Parse(packet []byte) (map[string]interface{}, error) {
	if len(packet) < 2 {
		return nil, errors.New("packet too short")
	}

	// Extract the packet ID (switch)
	packetID := fmt.Sprintf("%02X%02X", packet[1], packet[0])

	// Look up the handler
	info, exists := p.PacketList[packetID]
	if !exists {
		// Unknown packet
		return nil, nil
	}

	// Create the arguments map
	args := make(map[string]interface{})
	args["switch"] = packetID
	args["RAW_MSG"] = packet
	args["RAW_MSG_SIZE"] = len(packet)

	// Parse the packet data according to the format string
	if info.Format != "" && len(info.ParamNames) > 0 {
		offset := 2 // Skip the packet ID
		for i, paramName := range info.ParamNames {
			switch info.Format[i*3 : i*3+2] {
			case "v1": // uint16 (2 bytes)
				if offset+2 <= len(packet) {
					args[paramName] = binary.LittleEndian.Uint16(packet[offset : offset+2])
					offset += 2
				}
			case "V1": // uint32 (4 bytes)
				if offset+4 <= len(packet) {
					args[paramName] = binary.LittleEndian.Uint32(packet[offset : offset+4])
					offset += 4
				}
			case "C1": // uint8 (1 byte)
				if offset+1 <= len(packet) {
					args[paramName] = packet[offset]
					offset += 1
				}
			case "a*": // remaining bytes as string
				if offset < len(packet) {
					args[paramName] = packet[offset:]
					offset = len(packet)
				}
				// Add more format types as needed
			}
		}
	}

	// Call the handler if provided
	if info.Handler != nil {
		if err := info.Handler(packet); err != nil {
			return args, err
		}
	}

	return args, nil
}

// Reconstruct reconstructs a packet from arguments
func (p *PacketParser) Reconstruct(args map[string]interface{}) ([]byte, error) {
	switchVal, ok := args["switch"].(string)
	if !ok {
		return nil, errors.New("switch not provided or not a string")
	}

	// Check if the switch is a handler name
	if !isHexString(switchVal) {
		// Look up the packet ID by handler name
		var exists bool
		switchVal, exists = p.LookupPacketID(switchVal)
		if !exists {
			return nil, fmt.Errorf("unknown handler: %s", args["switch"])
		}
	}

	// Look up the packet info
	info, exists := p.PacketList[switchVal]
	if !exists {
		return nil, fmt.Errorf("unknown packet: %s", switchVal)
	}

	// Create the packet
	// First 2 bytes are the packet ID (switch)
	packet := make([]byte, 2)
	packet[0] = byte(hexToByte(switchVal[2:4]))
	packet[1] = byte(hexToByte(switchVal[0:2]))

	// Add the packet data according to the format string
	if info.Format != "" && len(info.ParamNames) > 0 {
		for i, paramName := range info.ParamNames {
			paramValue, exists := args[paramName]
			if !exists {
				return nil, fmt.Errorf("missing parameter: %s", paramName)
			}

			switch info.Format[i*3 : i*3+2] {
			case "v1": // uint16 (2 bytes)
				val, ok := paramValue.(uint16)
				if !ok {
					return nil, fmt.Errorf("parameter %s is not uint16", paramName)
				}
				buf := make([]byte, 2)
				binary.LittleEndian.PutUint16(buf, val)
				packet = append(packet, buf...)
			case "V1": // uint32 (4 bytes)
				val, ok := paramValue.(uint32)
				if !ok {
					return nil, fmt.Errorf("parameter %s is not uint32", paramName)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, val)
				packet = append(packet, buf...)
			case "C1": // uint8 (1 byte)
				val, ok := paramValue.(uint8)
				if !ok {
					return nil, fmt.Errorf("parameter %s is not uint8", paramName)
				}
				packet = append(packet, val)
			case "a*": // string
				val, ok := paramValue.([]byte)
				if !ok {
					return nil, fmt.Errorf("parameter %s is not []byte", paramName)
				}
				packet = append(packet, val...)
				// Add more format types as needed
			}
		}
	}

	return packet, nil
}

// Process processes packets from a tokenizer
func (p *PacketParser) Process(tokenizer *Tokenizer) error {
	for {
		message, msgType, err := tokenizer.ReadNext()
		if err == ErrIncompletePacket {
			// No more complete packets to process
			break
		}
		if err != nil {
			return err
		}

		switch msgType {
		case KnownMessage:
			args, err := p.Parse(message)
			if err != nil {
				return err
			}
			if args == nil {
				// Unknown packet, but not an error
				continue
			}

			// Get the packet ID
			packetID, ok := args["switch"].(string)
			if !ok {
				continue
			}

			// Look up the handler and call it if it exists
			if info, exists := p.PacketList[packetID]; exists && info.Handler != nil {
				if err := info.Handler(message); err != nil {
					return err
				}
			}
		case AccountID:
			// Handle account ID message
		case UnknownMessage:
			// Handle unknown message
		default:
			return fmt.Errorf("unknown message type: %d", msgType)
		}
	}

	return nil
}

// Helper functions

// isHexString checks if a string is a valid hex string
func isHexString(s string) bool {
	if len(s) != 4 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// hexToByte converts a hex string to a byte
func hexToByte(s string) byte {
	var b byte
	fmt.Sscanf(s, "%02x", &b)
	return b
}
