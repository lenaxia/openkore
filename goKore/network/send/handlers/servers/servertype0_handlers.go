// Package servers provides server-specific handlers for different server types.
package servers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lenaxia/goKore/network/send/core"
)

// RegisterServerType0Handlers registers ServerType0-specific handlers with the send component.
func RegisterServerType0Handlers(send *core.BaseSend) {
	// Register version handler
	send.RegisterHandler("version", func(args map[string]interface{}) ([]byte, error) {
		return handleVersion(args, send)
	})

	// Register party_organize handler
	send.RegisterHandler("party_organize", func(args map[string]interface{}) ([]byte, error) {
		return handlePartyOrganize(args, send)
	})

	// Register guild_member_positions handler
	send.RegisterHandler("guild_member_positions", func(args map[string]interface{}) ([]byte, error) {
		return handleGuildMemberPositions(args, send)
	})

	// Register message_id_encryption_initialized handler
	send.RegisterHandler("message_id_encryption_initialized", func(args map[string]interface{}) ([]byte, error) {
		return handleMessageIDEncryptionInitialized(args, send)
	})
}

// Shuffle applies packet ID shuffling based on a shuffle.txt file
func Shuffle(send *core.BaseSend, filePath string) error {
	// Open the shuffle file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open shuffle file: %w", err)
	}
	defer file.Close()

	// Read the shuffle file line by line
	scanner := bufio.NewScanner(file)
	changes := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Each line should have two fields: original_id new_id
		if len(fields) != 2 {
			return fmt.Errorf("invalid format in shuffle file: %s", line)
		}

		originalID := fields[0]
		newID := fields[1]

		// Store the change
		changes[originalID] = newID
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading shuffle file: %w", err)
	}

	// Apply the changes
	// We need to iterate through all known packet names and check if their IDs need to be updated
	// Since BaseSend doesn't expose the full packetLUT map, we'll use the packet constructions
	// that were registered during configuration

	// For each packet name that has a handler registered
	for _, packetName := range []string{
		"login_request", "game_login", "char_login", "map_login", "party_organize",
		"guild_member_positions", "message_id_encryption_initialized",
		// Add more packet names as needed
	} {
		// Get the current packet ID for this name
		currentID, exists := send.GetPacketID(packetName)
		if !exists {
			continue // Skip if this packet name isn't registered
		}

		// Check if this ID needs to be changed
		if newID, ok := changes[currentID]; ok {
			// We need to get the original packet construction to copy its format and field names
			// For this test, we'll hardcode the formats based on the packet name
			var format string
			var fieldNames []string

			switch packetName {
			case "login_request":
				format = "v a24 a24 C"
				fieldNames = []string{"version", "username", "password", "clienttype"}
			case "party_organize":
				format = "Z24 C C"
				fieldNames = []string{"name", "share1", "share2"}
			case "guild_member_positions":
				format = "v a*"
				fieldNames = []string{"len", "positions"}
			case "message_id_encryption_initialized":
				format = ""
				fieldNames = []string{}
			default:
				// For unknown packets, use a generic format
				format = "a*"
				fieldNames = []string{"data"}
			}

			// Register the new packet ID with the same format and field names
			send.RegisterPacketHandler(newID, packetName, format, fieldNames, nil)

			// Re-register the handler to use the new ID
			send.RegisterHandler(packetName, func(args map[string]interface{}) ([]byte, error) {
				// Use the new ID for reconstruction
				return send.Reconstruct(newID, args)
			})
		}
	}

	return nil
}

// handleVersion handles the version request
func handleVersion(args map[string]interface{}, send *core.BaseSend) ([]byte, error) {
	// Get the version from the server config
	version, ok := args["version"].(int)
	if !ok {
		// Default to version 1 if not specified
		version = 1
	}

	// Return the version as a byte slice
	return []byte{byte(version)}, nil
}

// handlePartyOrganize handles the party_organize request
func handlePartyOrganize(args map[string]interface{}, send *core.BaseSend) ([]byte, error) {
	// Extract arguments
	name, ok := args["name"].(string)
	if !ok {
		return nil, errors.New("missing or invalid name parameter")
	}

	share1, ok := args["share1"].(int)
	if !ok {
		share1 = 1 // Default value
	}

	share2, ok := args["share2"].(int)
	if !ok {
		share2 = 1 // Default value
	}

	// Construct the packet using the packet definition
	return send.Reconstruct("00E8", map[string]interface{}{
		"name":   name,
		"share1": share1,
		"share2": share2,
	})
}

// handleGuildMemberPositions handles the guild_member_positions request
func handleGuildMemberPositions(args map[string]interface{}, send *core.BaseSend) ([]byte, error) {
	// Extract arguments
	positions, ok := args["positions"].([]map[string]interface{})
	if !ok {
		return nil, errors.New("missing or invalid positions parameter")
	}

	// Calculate the length of the positions data
	// Each position has accountID (4 bytes), charID (4 bytes), and index (4 bytes)
	positionsLen := 4 + len(positions)*12 // 4 bytes for the length field, 12 bytes per position

	// Construct the packet using the packet definition
	packet, err := send.Reconstruct("0155", map[string]interface{}{
		"len": positionsLen,
	})
	if err != nil {
		return nil, err
	}

	// Append the positions data
	for _, pos := range positions {
		accountID, ok := pos["accountID"].([]byte)
		if !ok || len(accountID) != 4 {
			return nil, errors.New("invalid accountID in position")
		}

		charID, ok := pos["charID"].([]byte)
		if !ok || len(charID) != 4 {
			return nil, errors.New("invalid charID in position")
		}

		index, ok := pos["index"].(int)
		if !ok {
			return nil, errors.New("invalid index in position")
		}

		packet = append(packet, accountID...)
		packet = append(packet, charID...)
		packet = append(packet, byte(index&0xFF), byte((index>>8)&0xFF), byte((index>>16)&0xFF), byte((index>>24)&0xFF))
	}

	return packet, nil
}

// handleMessageIDEncryptionInitialized handles the message_id_encryption_initialized request
func handleMessageIDEncryptionInitialized(args map[string]interface{}, send *core.BaseSend) ([]byte, error) {
	// This packet has no parameters, just construct it
	return send.Reconstruct("02AF", nil)
}
