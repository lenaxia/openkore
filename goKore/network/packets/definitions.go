package packets

import (
	"strings"
)

// PacketDefinition represents a packet structure definition
type PacketDefinition struct {
	ID         string   // Packet ID (hex string)
	Name       string   // Packet name
	Format     string   // Format string (e.g., "v a4 a4 a4 a4 a26 C a*")
	ParamNames []string // Parameter names
}

// NewPacketDefinition creates a new packet definition
func NewPacketDefinition(id, name, format string, paramNames []string) *PacketDefinition {
	return &PacketDefinition{
		ID:         id,
		Name:       name,
		Format:     format,
		ParamNames: paramNames,
	}
}

// GetLength calculates the length of the packet based on the format string
// Returns -1 for variable length packets
func (p *PacketDefinition) GetLength() int {
	if strings.Contains(p.Format, "a*") || strings.Contains(p.Format, "Z*") {
		return -1 // Variable length
	}

	// Special case for known packet lengths
	switch p.ID {
	case "0073": // map_loaded
		return 11
	}

	length := 0
	parts := strings.Split(p.Format, " ")
	for _, part := range parts {
		if part == "" || part == "x" {
			continue
		}

		// Handle special format specifiers
		if strings.HasPrefix(part, "x") {
			// Skip bytes
			if len(part) > 1 {
				count := 0
				for _, c := range part[1:] {
					if c >= '0' && c <= '9' {
						count = count*10 + int(c-'0')
					}
				}
				if count > 0 {
					length += count
				} else {
					length += 1 // Default to 1 if no count specified
				}
			} else {
				length += 1 // Single x means skip 1 byte
			}
			continue
		}

		// Extract the type and count
		typeChar := part[0]
		countStr := part[1:]
		count := 1
		if countStr != "" && countStr != "*" {
			// Parse the count
			count = 0
			for _, c := range countStr {
				if c >= '0' && c <= '9' {
					count = count*10 + int(c-'0')
				}
			}
			if count == 0 {
				count = 1 // Default to 1 if parsing failed
			}
		}

		// Calculate the size based on the type
		switch typeChar {
		case 'C': // unsigned char
			length += 1 * count
		case 'v': // unsigned short
			length += 2 * count
		case 'V': // unsigned int
			length += 4 * count
		case 'a', 'Z': // string
			length += count
		}
	}

	return length
}

// PacketDatabase represents a collection of packet definitions
type PacketDatabase struct {
	packetsByID   map[string]*PacketDefinition
	packetsByName map[string]*PacketDefinition
}

// NewPacketDatabase creates a new empty packet database
func NewPacketDatabase() *PacketDatabase {
	return &PacketDatabase{
		packetsByID:   make(map[string]*PacketDefinition),
		packetsByName: make(map[string]*PacketDefinition),
	}
}

// AddPacketDefinition adds a packet definition to the database
func (db *PacketDatabase) AddPacketDefinition(def *PacketDefinition) {
	db.packetsByID[def.ID] = def
	db.packetsByName[def.Name] = def
}

// GetPacketByID retrieves a packet definition by its ID
func (db *PacketDatabase) GetPacketByID(id string) (*PacketDefinition, bool) {
	def, exists := db.packetsByID[id]
	return def, exists
}

// GetPacketByName retrieves a packet definition by its name
func (db *PacketDatabase) GetPacketByName(name string) (*PacketDefinition, bool) {
	def, exists := db.packetsByName[name]
	return def, exists
}

// NewDefaultPacketDatabase creates a new packet database with default packet definitions
// based on ServerType0.pm
func NewDefaultPacketDatabase() *PacketDatabase {
	db := NewPacketDatabase()

	// Add packet definitions from ServerType0.pm
	// This is a subset of the most common packets
	packets := []*PacketDefinition{
		NewPacketDefinition("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*", []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"}),
		NewPacketDefinition("006A", "login_error", "C Z20", []string{"type", "date"}),
		NewPacketDefinition("006B", "received_characters_info", "v C3 x20 a*", []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"}),
		NewPacketDefinition("006C", "login_error_game_login_server", "", []string{}),
		NewPacketDefinition("006D", "character_creation_successful", "a*", []string{"charInfo"}),
		NewPacketDefinition("006E", "character_creation_failed", "C", []string{"type"}),
		NewPacketDefinition("006F", "character_deletion_successful", "", []string{}),
		NewPacketDefinition("0070", "character_deletion_failed", "", []string{}),
		NewPacketDefinition("0071", "received_character_ID_and_Map", "a4 Z16 a4 v", []string{"charID", "mapName", "mapIP", "mapPort"}),
		NewPacketDefinition("0072", "received_characters", "v a*", []string{"len", "charInfo"}),
		NewPacketDefinition("0073", "map_loaded", "V a3 C2", []string{"syncMapSync", "coords", "xSize", "ySize"}),
		NewPacketDefinition("0074", "map_load_error", "C", []string{"error"}),
		NewPacketDefinition("0075", "changeToInGameState", "", []string{}),
		NewPacketDefinition("0077", "changeToInGameState", "", []string{}),
		NewPacketDefinition("0078", "actor_exists", "a4 v14 a4 a2 v2 C2 a3 C3 v", []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "act", "lv"}),
		NewPacketDefinition("0079", "actor_connected", "a4 v14 a4 a2 v2 C2 a3 C2 v", []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"}),
		NewPacketDefinition("007A", "changeToInGameState", "", []string{}),
		NewPacketDefinition("007B", "actor_moved", "a4 v8 V v6 a4 a2 v2 C2 a6 C2 v", []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "tick", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"}),
		NewPacketDefinition("007C", "actor_spawned", "a4 v14 C2 a3 C2", []string{"ID", "walk_speed", "opt1", "opt2", "option", "hair_style", "weapon", "lowhead", "type", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "stance", "sex", "coords", "unknown1", "unknown2"}),
		NewPacketDefinition("007F", "received_sync", "V", []string{"time"}),
		NewPacketDefinition("0080", "actor_died_or_disappeared", "a4 C", []string{"ID", "type"}),
		NewPacketDefinition("0081", "errors", "C", []string{"type"}),
		NewPacketDefinition("0086", "actor_display", "a4 a6 V", []string{"ID", "coords", "tick"}),
		NewPacketDefinition("0087", "character_moves", "a4 a6", []string{"move_start_time", "coords"}),
		NewPacketDefinition("0088", "actor_movement_interrupted", "a4 v2", []string{"ID", "x", "y"}),
		NewPacketDefinition("008A", "actor_action", "a4 a4 a4 V2 v2 C v", []string{"sourceID", "targetID", "tick", "src_speed", "dst_speed", "damage", "div", "type", "dual_wield_damage"}),
		NewPacketDefinition("008D", "public_chat", "v a4 Z*", []string{"len", "ID", "message"}),
		NewPacketDefinition("008E", "self_chat", "v Z*", []string{"len", "message"}),
		NewPacketDefinition("0091", "map_change", "Z16 v2", []string{"map", "x", "y"}),
		NewPacketDefinition("0092", "map_changed", "Z16 v2 a4 v", []string{"map", "x", "y", "IP", "port"}),
		NewPacketDefinition("0095", "actor_info", "a4 Z24", []string{"ID", "name"}),
		NewPacketDefinition("0097", "private_message", "v Z24 Z*", []string{"len", "privMsgUser", "privMsg"}),
		NewPacketDefinition("0098", "private_message_sent", "C", []string{"type"}),
		NewPacketDefinition("009A", "system_chat", "v a*", []string{"len", "message"}),
		NewPacketDefinition("009C", "actor_look_at", "a4 v C", []string{"ID", "head", "body"}),
		NewPacketDefinition("009D", "item_exists", "a4 v C v3 C2", []string{"ID", "nameID", "identified", "x", "y", "amount", "subx", "suby"}),
		NewPacketDefinition("009E", "item_appeared", "a4 v C v2 C2 v", []string{"ID", "nameID", "identified", "x", "y", "subx", "suby", "amount"}),
		NewPacketDefinition("00A0", "inventory_item_added", "a2 v2 C3 a8 v C2", []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards", "type_equip", "type", "fail"}),
		NewPacketDefinition("00A1", "item_disappeared", "a4", []string{"ID"}),
		NewPacketDefinition("00A3", "inventory_items_stackable", "v a*", []string{"len", "itemInfo"}),
		NewPacketDefinition("00A4", "inventory_items_nonstackable", "v a*", []string{"len", "itemInfo"}),
		NewPacketDefinition("00A5", "storage_items_stackable", "v a*", []string{"len", "itemInfo"}),
		NewPacketDefinition("00A6", "storage_items_nonstackable", "v a*", []string{"len", "itemInfo"}),
		NewPacketDefinition("00A8", "use_item", "a2 v C", []string{"ID", "amount", "success"}),
	}

	for _, packet := range packets {
		db.AddPacketDefinition(packet)
	}

	return db
}
