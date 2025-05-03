package common

// PacketDump represents the structure of the packet dump JSON files
type PacketDump struct {
	FileName    string   `json:"file_name"`
	PacketCount int      `json:"packet_count"`
	Packets     []Packet `json:"packets"`
}

// Packet represents a single packet in the dump
type Packet struct {
	Direction   string      `json:"direction"`
	PacketID    string      `json:"packet_id"`
	Description string      `json:"description"`
	Size        int         `json:"size"`
	Timestamp   string      `json:"timestamp"`
	RawData     []RawData   `json:"raw_data"`
	ParsedData  interface{} `json:"parsed_data"`
	ServerType  string      `json:"server_type"`
}

// RawData represents the raw data of a packet
type RawData struct {
	HexBytes            []string `json:"hex_bytes"`
	AsciiRepresentation string   `json:"ascii_representation"`
	BinaryBase64        string   `json:"binary_base64"`
}
