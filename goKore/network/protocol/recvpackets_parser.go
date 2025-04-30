// Package protocol provides functionality for handling the Ragnarok Online network protocol.
// This file implements the parser for recvpackets.txt files.
package protocol

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lenaxia/goKore/network/common"
)

// ParseRecvPackets parses a recvpackets.txt file and returns a map of packet IDs to their definitions
func ParseRecvPackets(filePath string) (map[string]PacketLengthDef, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open recvpackets.txt: %w", err)
	}
	defer file.Close()

	packetDefs := make(map[string]PacketLengthDef)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comments and empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into fields
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid format at line %d: %s", lineNum, line)
		}

		// Extract packet ID and length
		packetID := fields[0]
		lengthStr := fields[1]

		// Parse length
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid packet length at line %d: %s", lineNum, lengthStr)
		}

		// Determine if the packet has a length field
		hasLength := length < 0

		// Add to packet definitions
		packetDefs[packetID] = PacketLengthDef{
			Length:    length,
			HasLength: hasLength,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading recvpackets.txt: %w", err)
	}

	return packetDefs, nil
}

// ConvertToTokenizerPacketDefs converts a map of packet lengths to tokenizer packet definitions
func ConvertToTokenizerPacketDefs(recvPackets map[string]struct{ Length int }) map[string]PacketLengthDef {
	packetDefs := make(map[string]PacketLengthDef)

	for id, packet := range recvPackets {
		hasLength := packet.Length < 0
		packetDefs[id] = PacketLengthDef{
			Length:    packet.Length,
			HasLength: hasLength,
		}
	}

	return packetDefs
}

// LoadRecvPackets loads packet definitions from recvpackets.txt using the table loader
func LoadRecvPackets(basePath string, tableFolders []string) (map[string]PacketLengthDef, error) {
	// Create table loader
	loader := common.NewTableLoader(basePath, tableFolders)

	// Find recvpackets.txt
	recvpacketsPath, err := loader.FindTableFile("recvpackets.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to find recvpackets.txt: %w", err)
	}

	// Parse recvpackets.txt
	packetDefs, err := ParseRecvPackets(recvpacketsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recvpackets.txt: %w", err)
	}

	return packetDefs, nil
}
