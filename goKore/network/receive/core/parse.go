// Package core provides core functionality for parsing and processing network packets.
package core

import (
	"errors"
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/protocol"
)

// Errors
var (
	ErrInvalidPacket = errors.New("invalid packet")
	ErrPacketIgnored = errors.New("packet ignored by hook")
)

// PacketHandler is a function that processes a parsed packet
type PacketHandler func(args map[string]interface{}) error

// CoreParser is responsible for parsing and processing packets
type CoreParser struct {
	parser       *protocol.PacketParser
	handlers     map[string]PacketHandler
	hookManager  *hooks.HookManager
	serverType   string
	defaultState int
}

// NewCoreParser creates a new core parser
func NewCoreParser(serverType string, hookManager *hooks.HookManager) *CoreParser {
	return &CoreParser{
		parser:      protocol.NewPacketParser(),
		handlers:    make(map[string]PacketHandler),
		hookManager: hookManager,
		serverType:  serverType,
	}
}

// RegisterHandler registers a handler for a specific packet
func (p *CoreParser) RegisterHandler(packetID, name, format string, paramNames []string, handler PacketHandler) {
	// Register with the underlying protocol parser
	p.parser.RegisterHandler(packetID, name, format, paramNames, nil)

	// Store the handler in our map
	p.handlers[name] = handler
}

// RegisterHandlerFunc registers a handler function for a specific packet
func (p *CoreParser) RegisterHandlerFunc(packetID, name, format string, paramNames []string, handler func(args map[string]interface{}) error) {
	p.RegisterHandler(packetID, name, format, paramNames, PacketHandler(handler))
}

// Parse parses a packet and returns the parsed arguments
func (p *CoreParser) Parse(packet []byte) (map[string]interface{}, error) {
	// Use the underlying protocol parser to parse the packet
	args, err := p.parser.Parse(packet)
	if err != nil {
		return nil, err
	}

	if args == nil {
		// Unknown packet
		return nil, nil
	}

	return args, nil
}

// Process processes a packet, calling the appropriate handler and hooks
func (p *CoreParser) Process(packet []byte) error {
	// Parse the packet
	args, err := p.Parse(packet)
	if err != nil {
		return err
	}

	if args == nil {
		// Unknown packet, but not an error
		return nil
	}

	// Get the packet ID and handler name
	packetID, ok := args["switch"].(string)
	if !ok {
		return ErrInvalidPacket
	}

	// Look up the handler info
	info, exists := p.parser.PacketList[packetID]
	if !exists {
		// Unknown packet, but not an error
		return nil
	}

	// Call pre-processing hooks
	hookName := fmt.Sprintf("receive/packet_pre/%s", info.Name)
	if p.hookManager != nil && p.hookManager.HasHook(hookName) {
		p.hookManager.CallHook(hookName, args)
		// Check if the packet should be ignored
		if returnVal, ok := args["return"].(bool); ok && returnVal {
			return ErrPacketIgnored
		}
	}

	// Call the handler if it exists
	if handler, exists := p.handlers[info.Name]; exists {
		if err := handler(args); err != nil {
			return err
		}
	}

	// Call post-processing hooks
	hookName = fmt.Sprintf("receive/packet/%s", info.Name)
	if p.hookManager != nil {
		p.hookManager.CallHook(hookName, args)
	}

	return nil
}

// ProcessBuffer processes all complete packets in a buffer
func (p *CoreParser) ProcessBuffer(tokenizer *protocol.Tokenizer) error {
	for {
		packet, msgType, err := tokenizer.ReadNext()
		if err == protocol.ErrIncompletePacket {
			// No more complete packets to process
			break
		}
		if err != nil {
			return err
		}

		switch msgType {
		case protocol.KnownMessage:
			if err := p.Process(packet); err != nil && err != ErrPacketIgnored {
				return err
			}
		case protocol.AccountID:
			// Handle account ID message
			if p.hookManager != nil {
				p.hookManager.CallHook("receive/account_id", map[string]interface{}{
					"accountID": packet,
				})
			}
		case protocol.UnknownMessage:
			// Handle unknown message
			if p.hookManager != nil {
				p.hookManager.CallHook("receive/unknown_packet", map[string]interface{}{
					"packet": packet,
				})
			}
		default:
			return fmt.Errorf("unknown message type: %d", msgType)
		}
	}

	return nil
}

// GetHandler returns a handler by name
func (p *CoreParser) GetHandler(name string) (PacketHandler, bool) {
	handler, exists := p.handlers[name]
	return handler, exists
}

// GetPacketID returns a packet ID by handler name
func (p *CoreParser) GetPacketID(handlerName string) (string, bool) {
	return p.parser.LookupPacketID(handlerName)
}

// SetDefaultState sets the default state for the parser
func (p *CoreParser) SetDefaultState(state int) {
	p.defaultState = state
}

// GetDefaultState returns the default state for the parser
func (p *CoreParser) GetDefaultState() int {
	return p.defaultState
}

// GetServerType returns the server type for the parser
func (p *CoreParser) GetServerType() string {
	return p.serverType
}
