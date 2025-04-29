package core

import (
	"errors"
	"strings"
)

// BannedCharacter represents a character that has been banned
type BannedCharacter struct {
	Name string
}

// handleCharacterBanList handles the character_ban_list packet
// This packet is sent by the server to provide a list of banned characters
// Packet format: Header + Len + CharList[character_name(size:24)]
func (m *CharacterManager) handleCharacterBanList(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character list from args
	var charList []byte
	if charListVal, ok := args["charList"].([]byte); ok {
		charList = charListVal
	} else {
		return errors.New("invalid character list in character_ban_list packet")
	}

	// Check if the character list is valid
	if len(charList) < 1 {
		return errors.New("character list too short in character_ban_list packet")
	}

	// Get the number of entries
	numEntries := int(charList[0])

	// Check if the character list has enough bytes for the specified number of entries
	// Each entry is 24 bytes (character name)
	if len(charList) < 1+numEntries*24 {
		return errors.New("character list too short for specified number of entries in character_ban_list packet")
	}

	// Parse the character list
	banList := make([]string, 0, numEntries)
	for i := 0; i < numEntries; i++ {
		// Extract the character name
		start := 1 + i*24
		end := start + 24

		// Convert to string and trim null bytes
		name := strings.TrimRight(string(charList[start:end]), "\x00")

		// Add to the ban list
		banList = append(banList, name)
	}

	// Publish the ban list to hooks
	if m.parser != nil && m.parser.hookManager != nil {
		m.parser.hookManager.CallHook("character.ban_list", map[string]interface{}{
			"ban_list": banList,
		})
	}

	return nil
}

// RegisterCharacterBanListHandler registers the character_ban_list packet handler
func (m *CharacterManager) RegisterCharacterBanListHandler() {
	// Register handler for character ban list
	m.parser.RegisterHandlerFunc("0267", "character_ban_list", "b",
		[]string{"charList"},
		m.handleCharacterBanList)
}
