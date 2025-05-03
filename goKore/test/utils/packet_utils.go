package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// HexToBytes converts a hex string to a byte slice
func HexToBytes(hexStr string) ([]byte, error) {
	// Remove spaces and newlines
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	hexStr = strings.ReplaceAll(hexStr, "\n", "")

	return hex.DecodeString(hexStr)
}

// BytesToHex converts a byte slice to a hex string
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// FormatPacket formats a packet for display
func FormatPacket(packet []byte) string {
	if len(packet) < 2 {
		return fmt.Sprintf("Invalid packet (length %d)", len(packet))
	}

	packetID := fmt.Sprintf("%02X%02X", packet[1], packet[0])
	hexStr := BytesToHex(packet)

	return fmt.Sprintf("Packet %s [%d bytes]: %s", packetID, len(packet), hexStr)
}

// GetPacketID extracts the packet ID from a packet
func GetPacketID(packet []byte) string {
	if len(packet) < 2 {
		return ""
	}

	return fmt.Sprintf("%02X%02X", packet[1], packet[0])
}

// PacketTestCase represents a test case for packet handling
type PacketTestCase struct {
	Name           string
	PacketID       string
	RawHex         string
	RawData        []byte
	ExpectedFields map[string]interface{}
	Direction      string // "send" or "receive"
}

// NewPacketTestCase creates a new packet test case from a hex string
func NewPacketTestCase(name, packetID, rawHex string, expectedFields map[string]interface{}, direction string) (PacketTestCase, error) {
	rawData, err := HexToBytes(rawHex)
	if err != nil {
		return PacketTestCase{}, err
	}

	return PacketTestCase{
		Name:           name,
		PacketID:       packetID,
		RawHex:         rawHex,
		RawData:        rawData,
		ExpectedFields: expectedFields,
		Direction:      direction,
	}, nil
}
