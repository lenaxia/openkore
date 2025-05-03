package utils

// Essential login packet test data based on the essential_login_packets.md document

// Account Server Packets

// AccountServerLoginPacket represents the 0064 packet (Account Server Login)
var AccountServerLoginPacket = PacketTestCase{
	Name:     "Account Server Login",
	PacketID: "0064",
	RawHex:   "64001C00000062006F00740069006A006F00300000000000000000000000000000000000000000000000004D656C6F6E2E3737000000000000000000",
	ExpectedFields: map[string]interface{}{
		"username":       "botijo0",
		"password":       "",
		"version":        0,
		"master_version": 0,
	},
	Direction: "send",
}

// AccountInfoPacket represents the 0AC4 packet (Account Info With Server Info)
var AccountInfoPacket = PacketTestCase{
	Name:     "Account Info With Server Info",
	PacketID: "0AC4",
	RawHex:   "C40AE000E55DF6C182841E00012C9C53000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	ExpectedFields: map[string]interface{}{
		"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1},
		"accountID":  []byte{0x82, 0x84, 0x1E, 0x00},
		"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53},
	},
	Direction: "receive",
}

// Character Server Packets

// CharServerLoginPacket represents the 0065 packet (Character Server Login)
var CharServerLoginPacket = PacketTestCase{
	Name:     "Character Server Login",
	PacketID: "0065",
	RawHex:   "650082841E00E55DF6C1012C9C530000",
	ExpectedFields: map[string]interface{}{
		"accountID":  []byte{0x82, 0x84, 0x1E, 0x00},
		"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1},
		"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53},
		"accountSex": byte(0),
	},
	Direction: "send",
}

// CharLoginPacket represents the 0066 packet (Char Login)
var CharLoginPacket = PacketTestCase{
	Name:     "Char Login",
	PacketID: "0066",
	RawHex:   "660000",
	ExpectedFields: map[string]interface{}{
		"slot": byte(0),
	},
	Direction: "send",
}

// ReceivedCharactersInfoPacket represents the 082D packet (Received characters from Game Login Server)
var ReceivedCharactersInfoPacket = PacketTestCase{
	Name:     "Received Characters Info",
	PacketID: "082D",
	RawHex:   "2D081D000F0000090F000000000000000000000000000000",
	ExpectedFields: map[string]interface{}{
		"total_slot": byte(15),
	},
	Direction: "receive",
}

// ReceivedCharactersPacket represents the 006B packet (Received characters from Game Login Server)
var ReceivedCharactersPacket = PacketTestCase{
	Name:     "Received Characters",
	PacketID: "006B",
	RawHex:   "6B00B6000F0F0F000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	ExpectedFields: map[string]interface{}{
		"charInfo": []byte{},
	},
	Direction: "receive",
}

// PinCodeRequestPacket represents the 08B9 packet (PinCode Request)
var PinCodeRequestPacket = PacketTestCase{
	Name:     "PinCode Request",
	PacketID: "08B9",
	RawHex:   "B90837370000828400000000",
	ExpectedFields: map[string]interface{}{
		"accountID": []byte{0x82, 0x84, 0x1E, 0x00},
	},
	Direction: "receive",
}

// CharacterMapInfoPacket represents the 0AC5 packet (Received character ID and Map IP)
var CharacterMapInfoPacket = PacketTestCase{
	Name:     "Character Map Info",
	PacketID: "0AC5",
	RawHex:   "C50AF24902006765665F66696C64303700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	ExpectedFields: map[string]interface{}{
		"charID":  []byte{0xF2, 0x49, 0x02, 0x00},
		"mapName": "gef_fild07",
	},
	Direction: "receive",
}

// Map Server Packets

// MapLoginPacket represents the 0436 packet (Map Login)
var MapLoginPacket = PacketTestCase{
	Name:     "Map Login",
	PacketID: "0436",
	RawHex:   "360482841E00F24902E55DF6C1D65A0000",
	ExpectedFields: map[string]interface{}{
		"accountID": []byte{0x82, 0x84, 0x1E, 0x00},
		"charID":    []byte{0xF2, 0x49, 0x02, 0x00},
		"sessionID": []byte{0xE5, 0x5D, 0xF6, 0xC1},
	},
	Direction: "send",
}

// MapLoadedPacket represents the 007D packet (Map Loaded)
var MapLoadedPacket = PacketTestCase{
	Name:           "Map Loaded",
	PacketID:       "007D",
	RawHex:         "7D00",
	ExpectedFields: map[string]interface{}{},
	Direction:      "send",
}

// AccountIDPacket represents the 0283 packet (Account ID)
var AccountIDPacket = PacketTestCase{
	Name:     "Account ID",
	PacketID: "0283",
	RawHex:   "830200000000",
	ExpectedFields: map[string]interface{}{
		"accountID": []byte{0x82, 0x84, 0x1E, 0x00},
	},
	Direction: "receive",
}

// EnterMapPacket represents the 02EB packet (Enter Map)
var EnterMapPacket = PacketTestCase{
	Name:           "Enter Map",
	PacketID:       "02EB",
	RawHex:         "EB02C93E82023D8BF00505000000",
	ExpectedFields: map[string]interface{}{},
	Direction:      "receive",
}

// AllLoginPackets is a slice of all essential login packets
var AllLoginPackets = []PacketTestCase{
	AccountServerLoginPacket,
	AccountInfoPacket,
	CharServerLoginPacket,
	CharLoginPacket,
	ReceivedCharactersInfoPacket,
	ReceivedCharactersPacket,
	PinCodeRequestPacket,
	CharacterMapInfoPacket,
	MapLoginPacket,
	MapLoadedPacket,
	AccountIDPacket,
	EnterMapPacket,
}
