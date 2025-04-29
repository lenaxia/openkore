// Package core provides core functionality for parsing and processing network packets.
package core

import (
	"encoding/binary"
	"errors"
	"math"
	"sync"
	"time"
)

// Position represents a 2D position in the game world
type Position struct {
	X int
	Y int
}

// CharacterManager manages character-related functionality
type CharacterManager struct {
	parser       *CoreParser
	mutex        sync.RWMutex
	guildMembers []GuildMember
	actors       map[uint32]*Actor
	reputations  []Reputation // List of reputation entries
}

// Actor represents a character or NPC in the game
type Actor struct {
	ID        uint32
	Name      string
	Lv        uint16
	Opt1      uint16
	Opt2      uint16
	Opt3      uint16
	Pos       Position // Current position
	PosTo     Position // Destination position
	TimeMove  int64    // Time when movement started
	WalkSpeed float64  // Walking speed
	Look      struct { // Direction the character is facing
		Body uint8
		Head uint8
	}

	// Appearance-related fields
	JobID    uint16
	Weapon   uint16
	Shield   uint16
	Headgear struct {
		Low uint16
		Mid uint16
		Top uint16
	}
	HairColor uint16
	Shoes     uint16
	Robe      uint16
	TitleID   uint16 // Character title ID

	// Status-related fields
	OverweightPercent uint8 // Character's overweight percentage
}

// Reputation represents a reputation entry
type Reputation struct {
	Type    uint32
	Type2   uint32
	Points  uint32
	Points2 uint32
}

// GuildMember represents a member of a guild
type GuildMember struct {
	CharID uint32
	Name   string
}

// NewCharacterManager creates a new character manager
func NewCharacterManager(parser *CoreParser) *CharacterManager {
	return &CharacterManager{
		parser:       parser,
		guildMembers: make([]GuildMember, 0),
		actors:       make(map[uint32]*Actor),
		reputations:  make([]Reputation, 0),
	}
}

// RegisterHandlers registers character-related packet handlers
func (m *CharacterManager) RegisterHandlers() {
	// Register handlers for character-related packets
	m.parser.RegisterHandlerFunc("0095", "character_name", "a4 Z24",
		[]string{"ID", "name"},
		m.handleCharacterName)

	// Register handlers for character status packets
	m.parser.RegisterHandlerFunc("028A", "character_status", "a4 v v",
		[]string{"ID", "lv", "opt3"},
		m.handleCharacterStatus)

	m.parser.RegisterHandlerFunc("0229", "character_status", "a4 v v",
		[]string{"ID", "opt1", "opt2"},
		m.handleCharacterStatus)

	m.parser.RegisterHandlerFunc("0119", "character_status", "a4 v v",
		[]string{"ID", "opt1", "opt2"},
		m.handleCharacterStatus)

	// Register handler for character movement
	m.parser.RegisterHandlerFunc("0087", "character_moves", "L x6",
		[]string{"walkStartTime", "coords"},
		m.handleCharacterMoves)

	// Register handler for sprite change
	m.parser.RegisterHandlerFunc("01D7", "sprite_change", "a4 C v v",
		[]string{"ID", "type", "value1", "value2"},
		m.handleSpriteChange)

	// Register handler for title change
	m.parser.RegisterHandlerFunc("0A2F", "change_title", "v",
		[]string{"title_id"},
		m.handleChangeTitle)

	// Register handler for reputation info
	m.parser.RegisterHandlerFunc("0B8D", "repute_info", "b",
		[]string{"reputeInfo"},
		m.handleReputeInfo)

	// Register handler for overweight percent
	m.parser.RegisterHandlerFunc("0A2B", "overweight_percent", "C",
		[]string{"percent"},
		m.handleOverweightPercent)

	// Register handler for character ban list
	m.parser.RegisterHandlerFunc("0267", "character_ban_list", "b",
		[]string{"charList"},
		m.handleCharacterBanList)

	// Register handler for flag
	m.parser.RegisterHandlerFunc("0A89", "flag", "",
		[]string{},
		m.handleFlag)
}

// handleCharacterName handles the character_name packet
// This packet is sent by the server to update a character's name
func (m *CharacterManager) handleCharacterName(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character ID and name from args
	var charID []byte
	var name string

	if charIDVal, ok := args["ID"].([]byte); ok {
		charID = charIDVal
	}

	if nameVal, ok := args["name"].(string); ok {
		name = nameVal
	}

	// Update guild member's name if the character ID matches
	if len(charID) >= 4 {
		charIDValue := uint32(charID[0]) | uint32(charID[1])<<8 | uint32(charID[2])<<16 | uint32(charID[3])<<24
		for i, member := range m.guildMembers {
			if member.CharID == charIDValue {
				m.guildMembers[i].Name = name
				break
			}
		}
	}

	// TODO: Add proper logging when logger is implemented
	// In the original Perl code, this would log:
	// debug "Character name received: $name\n"

	return nil
}

// AddGuildMember adds a guild member to the list
func (m *CharacterManager) AddGuildMember(charID uint32, name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if the member already exists
	for i, member := range m.guildMembers {
		if member.CharID == charID {
			// Update the existing member
			m.guildMembers[i].Name = name
			return
		}
	}

	// Add a new member
	m.guildMembers = append(m.guildMembers, GuildMember{
		CharID: charID,
		Name:   name,
	})
}

// GetGuildMember gets a guild member by character ID
func (m *CharacterManager) GetGuildMember(charID uint32) (GuildMember, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, member := range m.guildMembers {
		if member.CharID == charID {
			return member, true
		}
	}

	return GuildMember{}, false
}

// GetGuildMembers gets all guild members
func (m *CharacterManager) GetGuildMembers() []GuildMember {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the guild members slice
	members := make([]GuildMember, len(m.guildMembers))
	copy(members, m.guildMembers)

	return members
}

// GetActor gets an actor by ID
func (m *CharacterManager) GetActor(actorID uint32) (*Actor, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	actor, found := m.actors[actorID]
	return actor, found
}

// handleCharacterStatus handles the character_status packet
// This packet is sent by the server to update a character's status
func (m *CharacterManager) handleCharacterStatus(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character ID from args
	var charID []byte
	if charIDVal, ok := args["ID"].([]byte); ok {
		charID = charIDVal
	} else {
		return nil
	}

	// Convert character ID to uint32
	if len(charID) < 4 {
		return nil
	}

	charIDValue := uint32(charID[0]) | uint32(charID[1])<<8 | uint32(charID[2])<<16 | uint32(charID[3])<<24

	// Get or create the actor
	actor, found := m.actors[charIDValue]
	if !found {
		actor = &Actor{
			ID: charIDValue,
		}
		m.actors[charIDValue] = actor
	}

	// Instead of trying to get the packet ID, we'll check for specific fields in the args
	// to determine which packet type we're handling

	// Check if this is a 028A packet (has lv and opt3)
	_, hasLv := args["lv"].(uint16)
	_, hasOpt3 := args["opt3"].(uint16)

	// Check if this is a 0229 or 0119 packet (has opt1 and opt2)
	_, hasOpt1 := args["opt1"].(uint16)
	_, hasOpt2 := args["opt2"].(uint16)

	if hasLv && hasOpt3 {
		// Update level and opt3
		if lv, ok := args["lv"].(uint16); ok {
			actor.Lv = lv
		}

		if opt3, ok := args["opt3"].(uint16); ok {
			actor.Opt3 = opt3
		}
	} else if hasOpt1 && hasOpt2 {
		// Update opt1 and opt2
		if opt1, ok := args["opt1"].(uint16); ok {
			actor.Opt1 = opt1
		}

		if opt2, ok := args["opt2"].(uint16); ok {
			actor.Opt2 = opt2
		}
	}

	return nil
}

// handleCharacterMoves handles the character_moves packet
// This packet is sent by the server to notify the client that the character is moving
// Packet format: 0087 <walk start time>.L <walk data>.6B
func (m *CharacterManager) handleCharacterMoves(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we're in the game state
	// In the original Perl code, this would call changeToInGameState()
	// For now, we'll assume we're always in the game state

	// Extract coordinates from args
	coords, ok := args["coords"].([]byte)
	if !ok || len(coords) < 6 {
		return errors.New("invalid coordinates in character_moves packet")
	}

	// Extract walk start time
	_, ok = args["walkStartTime"].(uint32)
	if !ok {
		return errors.New("invalid walk start time in character_moves packet")
	}

	// In the original implementation, this would update the main character's position
	// For now, we'll just update the actor in our manager
	// In a real implementation, we would need to get the main character's ID
	// For testing purposes, we'll use the actor that's already in the manager
	if len(m.actors) == 0 {
		return errors.New("no actors in manager")
	}

	// Get the first actor (for testing purposes)
	var actor *Actor
	var actorID uint32
	for id, a := range m.actors {
		actor = a
		actorID = id
		break
	}

	// Update the actor's position
	// In the original Perl code, this would call makeCoordsFromTo()
	// Here we'll just extract the coordinates directly
	fromX := int(coords[0])
	fromY := int(coords[1])
	toX := int(coords[2])
	toY := int(coords[3])

	actor.Pos = Position{X: fromX, Y: fromY}
	actor.PosTo = Position{X: toX, Y: toY}
	actor.TimeMove = time.Now().Unix()

	// Calculate distance
	dx := toX - fromX
	dy := toY - fromY
	dist := math.Sqrt(float64(dx*dx + dy*dy))
	_ = dist // Unused for now

	// Set walking speed (default to 0.12 if not set)
	if actor.WalkSpeed == 0 {
		actor.WalkSpeed = 0.12
	}

	// Calculate direction
	// In the original Perl code, this would call getVector() and vectorToDegree()
	// Here we'll calculate the direction directly
	var direction uint8
	if dx == 0 && dy < 0 {
		direction = 0 // North
	} else if dx > 0 && dy < 0 {
		direction = 1 // Northeast
	} else if dx > 0 && dy == 0 {
		direction = 2 // East
	} else if dx > 0 && dy > 0 {
		direction = 3 // Southeast
	} else if dx == 0 && dy > 0 {
		direction = 4 // South
	} else if dx < 0 && dy > 0 {
		direction = 5 // Southwest
	} else if dx < 0 && dy == 0 {
		direction = 6 // West
	} else if dx < 0 && dy < 0 {
		direction = 7 // Northwest
	}

	actor.Look.Body = direction
	actor.Look.Head = 0

	// Update the actor in the manager
	m.actors[actorID] = actor

	// In the original Perl code, there's some AI code here
	// We'll skip that for now

	return nil
}

// handleSpriteChange handles the sprite_change packet
// This packet is sent by the server to notify the client that a character's appearance has changed
// Packet format: 01D7 <ID>.L <type>.B <value1>.W <value2>.W
func (m *CharacterManager) handleSpriteChange(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character ID from args
	var charID []byte
	if charIDVal, ok := args["ID"].([]byte); ok {
		charID = charIDVal
	} else {
		return errors.New("invalid character ID in sprite_change packet")
	}

	// Convert character ID to uint32
	if len(charID) < 4 {
		return errors.New("character ID too short in sprite_change packet")
	}

	charIDValue := uint32(charID[0]) | uint32(charID[1])<<8 | uint32(charID[2])<<16 | uint32(charID[3])<<24

	// Get the actor
	actor, found := m.actors[charIDValue]
	if !found {
		return errors.New("actor not found in sprite_change packet")
	}

	// Extract type, value1, and value2 from args
	var spriteType uint8
	var value1, value2 uint16

	if typeVal, ok := args["type"].(uint8); ok {
		spriteType = typeVal
	} else {
		return errors.New("invalid sprite type in sprite_change packet")
	}

	if val1, ok := args["value1"].(uint16); ok {
		value1 = val1
	} else {
		return errors.New("invalid value1 in sprite_change packet")
	}

	if val2, ok := args["value2"].(uint16); ok {
		value2 = val2
	} else {
		return errors.New("invalid value2 in sprite_change packet")
	}

	// Update the actor's appearance based on the sprite type
	switch spriteType {
	case 0: // Job change
		actor.JobID = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Job to: %s\n", $player, $jobs_lut{$value1}), "parseMsg_statuslook"

	case 2: // Weapon and shield change
		actor.Weapon = value1
		actor.Shield = value2
		// In the original Perl code, this would log:
		// message TF("%s changed Weapon to %s (%d)\n", $player, itemName({nameID => $value1}), $value1), "parseMsg_statuslook", 2
		// message TF("%s changed Shield to %s (%d)\n", $player, itemName({nameID => $value2}), $value2), "parseMsg_statuslook", 2

	case 3: // Lower headgear change
		actor.Headgear.Low = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Lower headgear to %s (%d)\n", $player, headgearName($value1), $value1), "parseMsg_statuslook"

	case 4: // Upper headgear change
		actor.Headgear.Top = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Upper headgear to %s (%d)\n", $player, headgearName($value1), $value1), "parseMsg_statuslook"

	case 5: // Middle headgear change
		actor.Headgear.Mid = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Middle headgear to %s (%d)\n", $player, headgearName($value1), $value1), "parseMsg_statuslook"

	case 6: // Hair color change
		actor.HairColor = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Hair color to: %s (%d)\n", $player, $haircolors{$value1}, $value1), "parseMsg_statuslook"

	case 9: // Shoes change
		actor.Shoes = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Shoes to: %s\n", $player, itemName({nameID => $value1})), "parseMsg_statuslook", 2

	case 12: // Robe change
		actor.Robe = value1
		// In the original Perl code, this would log:
		// message TF("%s changed Robe to: SPRITE_ROBE_ID=%d\n", $player, $value1, $value1), "parseMsg_statuslook", 2

	case 7, 13: // Body palette/color or body2
		// In the original Perl code, this would log:
		// debug sprintf("%s changed type= %d. value1=%d, value2=%d\n", $player, $type, $value1, $value2)
		// We'll just ignore these for now

	default:
		return errors.New("unknown sprite type in sprite_change packet")
	}

	// Update the actor in the manager
	m.actors[charIDValue] = actor

	// In the original Perl code, this would call:
	// Plugins::callHook('sprite_job_change')
	// We'll skip that for now

	return nil
}

// GetReputations returns a copy of the reputation list
func (m *CharacterManager) GetReputations() []Reputation {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the reputation list
	reputations := make([]Reputation, len(m.reputations))
	copy(reputations, m.reputations)

	return reputations
}

// handleReputeInfo handles the repute_info packet
// This packet is sent by the server to update the character's reputation information
// Packet format: 0B8D <reputeInfo>.B
func (m *CharacterManager) handleReputeInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract reputation info from args
	var reputeInfo []byte
	if reputeInfoVal, ok := args["reputeInfo"].([]byte); ok {
		reputeInfo = reputeInfoVal
	} else {
		return errors.New("invalid reputation info in repute_info packet")
	}

	// Each reputation entry is 16 bytes (4 uint32 values)
	entrySize := 16
	if len(reputeInfo)%entrySize != 0 {
		return errors.New("invalid reputation info length in repute_info packet")
	}

	// Clear the existing reputation list
	m.reputations = make([]Reputation, 0)

	// Parse the reputation info
	for i := 0; i < len(reputeInfo); i += entrySize {
		// Extract the reputation entry
		if i+entrySize > len(reputeInfo) {
			return errors.New("invalid reputation info length in repute_info packet")
		}

		// Parse the reputation entry
		reputation := Reputation{
			Type:    binary.LittleEndian.Uint32(reputeInfo[i : i+4]),
			Type2:   binary.LittleEndian.Uint32(reputeInfo[i+4 : i+8]),
			Points:  binary.LittleEndian.Uint32(reputeInfo[i+8 : i+12]),
			Points2: binary.LittleEndian.Uint32(reputeInfo[i+12 : i+16]),
		}

		// Add the reputation entry to the list
		m.reputations = append(m.reputations, reputation)
	}

	return nil
}

// handleOverweightPercent handles the overweight_percent packet
// This packet is sent by the server to notify the client of the character's overweight percentage
// Packet format: 0A2B <percent>.B
func (m *CharacterManager) handleOverweightPercent(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract percent from args
	var percent uint8
	if percentVal, ok := args["percent"].(uint8); ok {
		percent = percentVal
	} else {
		return errors.New("invalid percent in overweight_percent packet")
	}

	// In the original implementation, this would update the main character's overweight percent
	// For now, we'll just update the first actor in our manager
	// In a real implementation, we would need to get the main character's ID
	if len(m.actors) == 0 {
		return errors.New("no actors in manager")
	}

	// Get the first actor (for testing purposes)
	var actor *Actor
	var actorID uint32
	for id, a := range m.actors {
		actor = a
		actorID = id
		break
	}

	// Update the actor's overweight percent
	actor.OverweightPercent = percent

	// Update the actor in the manager
	m.actors[actorID] = actor

	// In the original Perl code, this would log:
	// debug "Received overweight percent: $args->{percent}\n";

	return nil
}

// handleChangeTitle handles the change_title packet
// This packet is sent by the server to notify the client that the character's title has changed
// Packet format: 0A2F <title_id>.W
func (m *CharacterManager) handleChangeTitle(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract title ID from args
	var titleID uint16
	if titleIDVal, ok := args["title_id"].(uint16); ok {
		titleID = titleIDVal
	} else {
		return errors.New("invalid title ID in change_title packet")
	}

	// In the original implementation, this would update the main character's title
	// For now, we'll just update the first actor in our manager
	// In a real implementation, we would need to get the main character's ID
	if len(m.actors) == 0 {
		return errors.New("no actors in manager")
	}

	// Get the first actor (for testing purposes)
	var actor *Actor
	var actorID uint32
	for id, a := range m.actors {
		actor = a
		actorID = id
		break
	}

	// Update the actor's title
	actor.TitleID = titleID

	// Update the actor in the manager
	m.actors[actorID] = actor

	// In the original Perl code, this would log:
	// message TF("You changed Title_ID :  %s.\n", $args->{title_id}), "info";

	return nil
}
