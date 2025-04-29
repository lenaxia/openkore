// Package core provides core functionality for parsing and processing network packets.
package core

import (
	"fmt"
	"sync"
)

// PacketEvent represents an event triggered by a packet
type PacketEvent struct {
	Type    string
	Data    interface{}
	RawData []byte
}

// EventBus handles event publishing and subscription
type EventBus struct {
	handlers map[string][]func(PacketEvent) error
	mutex    sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func(PacketEvent) error),
	}
}

// Subscribe registers a handler for a specific event type
func (b *EventBus) Subscribe(eventType string, handler func(PacketEvent) error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = make([]func(PacketEvent) error, 0)
	}
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all subscribers
func (b *EventBus) Publish(event PacketEvent) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if handlers, exists := b.handlers[event.Type]; exists {
		for _, handler := range handlers {
			// In a production system, we might want to handle errors or use goroutines
			_ = handler(event)
		}
	}
}

// MannerMessageProcessor processes manner_message packets
type MannerMessageProcessor struct {
	eventBus *EventBus
}

// NewMannerMessageProcessor creates a new manner message processor
func NewMannerMessageProcessor(eventBus *EventBus) *MannerMessageProcessor {
	return &MannerMessageProcessor{
		eventBus: eventBus,
	}
}

// HandlePacket processes a manner_message packet and returns a result
func (p *MannerMessageProcessor) HandlePacket(packetType string, data map[string]interface{}) (interface{}, error) {
	var flag uint8
	var message string

	// Extract flag from data
	if flagVal, ok := data["flag"].(uint8); ok {
		flag = flagVal
	}

	// Process based on flag value
	switch flag {
	case 0:
		message = "A manner point has been successfully aligned."
	case 3:
		message = "Chat Block has been applied by GM due to your ill-mannerous action."
	case 4:
		message = "Automated Chat Block has been applied due to Anti-Spam System."
	case 5:
		message = "You got a good point."
	default:
		message = fmt.Sprintf("Unknown manner message result (flag: %d)", flag)
	}

	// Publish event
	p.eventBus.Publish(PacketEvent{
		Type: "manner_message",
		Data: map[string]interface{}{
			"flag":    flag,
			"message": message,
		},
	})

	return message, nil
}

// HackShieldAlarmProcessor processes hack_shield_alarm packets
type HackShieldAlarmProcessor struct {
	eventBus *EventBus
}

// NewHackShieldAlarmProcessor creates a new hack shield alarm processor
func NewHackShieldAlarmProcessor(eventBus *EventBus) *HackShieldAlarmProcessor {
	return &HackShieldAlarmProcessor{
		eventBus: eventBus,
	}
}

// HandlePacket processes a hack_shield_alarm packet and returns a result
func (p *HackShieldAlarmProcessor) HandlePacket(packetType string, data map[string]interface{}) (interface{}, error) {
	message := "Error: You have been forced to disconnect by a Hack Shield. Please check Poseidon."

	// Publish event
	p.eventBus.Publish(PacketEvent{
		Type: "hack_shield_alarm",
		Data: map[string]interface{}{
			"message": message,
		},
	})

	// In the original implementation, this would also run a command to relog
	// Commands::run('relog 100000000');
	// We'll need to implement this functionality elsewhere

	return message, nil
}

// BehaviorManager manages behavior-related functionality
type BehaviorManager struct {
	parser   *CoreParser
	eventBus *EventBus
}

// NewBehaviorManager creates a new behavior manager
func NewBehaviorManager(parser *CoreParser) *BehaviorManager {
	return &BehaviorManager{
		parser:   parser,
		eventBus: NewEventBus(),
	}
}

// RegisterHandlers registers behavior-related packet handlers
func (m *BehaviorManager) RegisterHandlers() {
	// Create processors
	mannerProcessor := NewMannerMessageProcessor(m.eventBus)
	hackShieldProcessor := NewHackShieldAlarmProcessor(m.eventBus)

	// Register handler for manner_message
	m.parser.RegisterHandlerFunc("0149", "manner_message", "C",
		[]string{"flag"},
		func(args map[string]interface{}) error {
			result, err := mannerProcessor.HandlePacket("manner_message", args)
			if err != nil {
				return err
			}

			// Log the result
			// In a real implementation, we would use a proper logger
			_ = result

			return nil
		})

	// Register handler for hack_shield_alarm
	m.parser.RegisterHandlerFunc("08B3", "hack_shield_alarm", "",
		[]string{},
		func(args map[string]interface{}) error {
			result, err := hackShieldProcessor.HandlePacket("hack_shield_alarm", args)
			if err != nil {
				return err
			}

			// Log the result
			// In a real implementation, we would use a proper logger
			_ = result

			return nil
		})
}

// GetEventBus returns the event bus
func (m *BehaviorManager) GetEventBus() *EventBus {
	return m.eventBus
}
