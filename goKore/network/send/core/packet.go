package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidFormat is returned when a packet format is invalid.
	ErrInvalidFormat = errors.New("invalid packet format")

	// ErrInvalidArgument is returned when an argument is invalid.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrMissingArgument is returned when a required argument is missing.
	ErrMissingArgument = errors.New("missing argument")
)

// PacketBuilder is responsible for constructing packets.
type PacketBuilder struct {
	// Packet definitions.
	packetDefinitions map[string]PacketDefinition

	// Packet ID lookup table.
	packetLUT map[string]string
}

// NewPacketBuilder creates a new packet builder.
func NewPacketBuilder() *PacketBuilder {
	return &PacketBuilder{
		packetDefinitions: make(map[string]PacketDefinition),
		packetLUT:         make(map[string]string),
	}
}

// RegisterPacket registers a packet definition.
func (pb *PacketBuilder) RegisterPacket(id, name, format string, keys []string) {
	pb.packetDefinitions[id] = PacketDefinition{
		ID:     id,
		Name:   name,
		Format: format,
		Keys:   keys,
	}
	pb.packetLUT[name] = id
}

// GetPacketID returns the packet ID for a given packet name.
func (pb *PacketBuilder) GetPacketID(name string) (string, bool) {
	id, exists := pb.packetLUT[name]
	return id, exists
}

// BuildPacket constructs a packet from a packet ID and arguments.
func (pb *PacketBuilder) BuildPacket(packetID string, args map[string]interface{}) ([]byte, error) {
	// Check if the packet is registered
	def, exists := pb.packetDefinitions[packetID]
	if !exists {
		return nil, fmt.Errorf("%w: packet %s not registered", ErrPacketNotRegistered, packetID)
	}

	// Parse the format string
	formatParts := parseFormat(def.Format)
	if len(formatParts) == 0 {
		return nil, fmt.Errorf("%w: empty format", ErrInvalidFormat)
	}

	// Create the packet
	packet := make([]byte, 0)

	// Add the packet ID
	if len(packetID) != 4 {
		return nil, fmt.Errorf("%w: packet ID must be 4 characters", ErrInvalidPacketID)
	}

	id1, err := strconv.ParseUint(packetID[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPacketID, err)
	}

	id2, err := strconv.ParseUint(packetID[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPacketID, err)
	}

	packet = append(packet, byte(id2), byte(id1))

	// Process each format part
	keyIndex := 0
	for _, part := range formatParts {
		// Skip empty parts
		if part == "" {
			continue
		}

		// Get the format type and count
		formatType, count := parseFormatPart(part)
		if formatType == "" {
			return nil, fmt.Errorf("%w: invalid format part %s", ErrInvalidFormat, part)
		}

		// Get the key for this format part
		var key string
		if keyIndex < len(def.Keys) {
			key = def.Keys[keyIndex]
			keyIndex++
		} else {
			// If there are more format parts than keys, use a generated key
			key = fmt.Sprintf("arg%d", keyIndex)
			keyIndex++
		}

		// Special case for padding format type which doesn't need an argument
		if formatType == "x" {
			// For padding, we don't need an argument
			packet = append(packet, make([]byte, count)...)
			continue
		}

		// Get the argument value for all other format types
		value, exists := args[key]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrMissingArgument, key)
		}

		// Process the value based on the format type
		switch formatType {
		case "C": // unsigned char (1 byte)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				packet = append(packet, byte(val))
			}
		case "v": // unsigned short (2 bytes)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 2)
				binary.LittleEndian.PutUint16(buf, uint16(val))
				packet = append(packet, buf...)
			}
		case "V": // unsigned int (4 bytes)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(val))
				packet = append(packet, buf...)
			}
		case "a": // array of bytes (n bytes)
			data, err := getBytesValue(value)
			if err != nil {
				return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
			}
			if len(data) > count {
				data = data[:count]
			} else if len(data) < count {
				// Pad with zeros
				data = append(data, make([]byte, count-len(data))...)
			}
			packet = append(packet, data...)
		case "Z": // null-terminated string (n bytes)
			data, err := getStringValue(value)
			if err != nil {
				return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
			}
			// Convert to bytes
			bytes := []byte(data)
			if len(bytes) > count-1 {
				bytes = bytes[:count-1]
			}
			// Pad with zeros
			bytes = append(bytes, 0)
			if len(bytes) < count {
				bytes = append(bytes, make([]byte, count-len(bytes))...)
			}
			packet = append(packet, bytes...)
		case "b": // signed char (1 byte)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				packet = append(packet, byte(val))
			}
		case "w": // signed short (2 bytes)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 2)
				binary.LittleEndian.PutUint16(buf, uint16(val))
				packet = append(packet, buf...)
			}
		case "l": // signed int (4 bytes)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, uint32(val))
				packet = append(packet, buf...)
			}
		case "q": // signed long long (8 bytes)
			for i := 0; i < count; i++ {
				val, err := getIntValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 8)
				binary.LittleEndian.PutUint64(buf, uint64(val))
				packet = append(packet, buf...)
			}
		case "f": // float (4 bytes)
			for i := 0; i < count; i++ {
				val, err := getFloatValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 4)
				binary.LittleEndian.PutUint32(buf, float32ToUint32(float32(val)))
				packet = append(packet, buf...)
			}
		case "d": // double (8 bytes)
			for i := 0; i < count; i++ {
				val, err := getFloatValue(value, i)
				if err != nil {
					return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
				}
				buf := make([]byte, 8)
				binary.LittleEndian.PutUint64(buf, float64ToUint64(val))
				packet = append(packet, buf...)
			}
		case "*": // variable length string
			data, err := getStringValue(value)
			if err != nil {
				return nil, fmt.Errorf("%w: %s: %v", ErrInvalidArgument, key, err)
			}
			// Convert to bytes
			bytes := []byte(data)
			// Add null terminator
			bytes = append(bytes, 0)
			packet = append(packet, bytes...)
		default:
			return nil, fmt.Errorf("%w: unknown format type %s", ErrInvalidFormat, formatType)
		}
	}

	return packet, nil
}

// parseFormat parses a format string into its parts.
func parseFormat(format string) []string {
	// Split the format string into parts
	parts := strings.Fields(format)
	return parts
}

// parseFormatPart parses a format part into its type and count.
func parseFormatPart(part string) (string, int) {
	// Check if the part is empty
	if part == "" {
		return "", 0
	}

	// Get the format type (first character)
	formatType := part[0:1]

	// Get the count (rest of the string)
	count := 1
	if len(part) > 1 {
		// Try to parse the count
		var err error
		count, err = strconv.Atoi(part[1:])
		if err != nil {
			// If the count is not a number, assume it's 1
			count = 1
		}
	}

	return formatType, count
}

// getIntValue gets an integer value from an interface{}.
func getIntValue(value interface{}, index int) (int64, error) {
	// Check if the value is an array or slice
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
		// Check if the index is valid
		if index >= rv.Len() {
			return 0, fmt.Errorf("index %d out of range", index)
		}
		// Get the value at the index
		value = rv.Index(index).Interface()
	} else if index > 0 {
		return 0, fmt.Errorf("index %d out of range", index)
	}

	// Convert the value to an integer
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		// Try to parse the string as an integer
		i, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to integer: %v", err)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to integer", value)
	}
}

// getFloatValue gets a float value from an interface{}.
func getFloatValue(value interface{}, index int) (float64, error) {
	// Check if the value is an array or slice
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
		// Check if the index is valid
		if index >= rv.Len() {
			return 0, fmt.Errorf("index %d out of range", index)
		}
		// Get the value at the index
		value = rv.Index(index).Interface()
	} else if index > 0 {
		return 0, fmt.Errorf("index %d out of range", index)
	}

	// Convert the value to a float
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		// Try to parse the string as a float
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to float: %v", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float", value)
	}
}

// getBytesValue gets a byte slice from an interface{}.
func getBytesValue(value interface{}) ([]byte, error) {
	// Convert the value to a byte slice
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case []rune:
		return []byte(string(v)), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to byte slice", value)
	}
}

// getStringValue gets a string from an interface{}.
func getStringValue(value interface{}) (string, error) {
	// Convert the value to a string
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case []rune:
		return string(v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", value)
	}
}

// float32ToUint32 converts a float32 to a uint32.
func float32ToUint32(f float32) uint32 {
	return binary.LittleEndian.Uint32([]byte{
		byte(0),
		byte(0),
		byte(0),
		byte(0),
	})
}

// float64ToUint64 converts a float64 to a uint64.
func float64ToUint64(f float64) uint64 {
	return binary.LittleEndian.Uint64([]byte{
		byte(0),
		byte(0),
		byte(0),
		byte(0),
		byte(0),
		byte(0),
		byte(0),
		byte(0),
	})
}
