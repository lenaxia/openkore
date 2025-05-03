package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lenaxia/goKore/cmd/common"
	networkCommon "github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	receiveFactory "github.com/lenaxia/goKore/network/receive/factory"
	receiveServers "github.com/lenaxia/goKore/network/receive/servers"
	sendFactory "github.com/lenaxia/goKore/network/send/factory"
	sendServers "github.com/lenaxia/goKore/network/send/servers"
)

// termios is a struct that describes terminal attributes
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

// makeRaw puts the terminal connected to the given file descriptor into raw mode
func makeRaw(fd int) (*termios, error) {
	var oldState termios

	// Get the current terminal state
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState)))
	if errno != 0 {
		return nil, errno
	}

	// Make a copy of the current state to modify
	newState := oldState

	// Turn off canonical mode and echo
	newState.Lflag &^= syscall.ICANON | syscall.ECHO

	// Set the terminal attributes
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)))
	if errno != 0 {
		return nil, errno
	}

	return &oldState, nil
}

// restoreTerminal restores the terminal connected to the given file descriptor to its previous state
func restoreTerminal(fd int, state *termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(state)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Column widths
const (
	colWidth      = 30
	arrowColWidth = 7
)

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

// PacketDebugger is the main struct for the packet debugger
type PacketDebugger struct {
	dumpFiles      []string
	selectedDump   *common.PacketDump
	currentPacket  int
	receiveFactory *receiveFactory.ReceiveFactory
	sendFactory    *sendFactory.SendFactory
	hookManager    *hooks.HookManager
	reader         *bufio.Reader
}

// NewPacketDebugger creates a new packet debugger
func NewPacketDebugger() *PacketDebugger {
	pd := &PacketDebugger{
		currentPacket: 0,
		hookManager:   hooks.NewHookManager(),
		reader:        bufio.NewReader(os.Stdin),
	}

	// Initialize factories
	pd.receiveFactory = receiveFactory.NewReceiveFactory()
	pd.receiveFactory.RegisterDefaultServerTypes()

	pd.sendFactory = sendFactory.NewSendFactoryAligned(pd.hookManager)
	pd.sendFactory.RegisterDefaultServerTypes()

	return pd
}

// loadDumpFiles loads the list of available dump files
func (pd *PacketDebugger) loadDumpFiles(dumpDir string) error {
	files, err := ioutil.ReadDir(dumpDir)
	if err != nil {
		return err
	}

	pd.dumpFiles = []string{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			pd.dumpFiles = append(pd.dumpFiles, file.Name())
		}
	}

	return nil
}

// loadDump loads a packet dump file
func (pd *PacketDebugger) loadDump(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var dump common.PacketDump
	if err := json.Unmarshal(data, &dump); err != nil {
		return err
	}

	pd.selectedDump = &dump
	pd.currentPacket = 0

	fmt.Printf("Loaded %s with %d packets.\n",
		filepath.Base(path), dump.PacketCount)

	return nil
}

// wrapString wraps a string to a specified length
func wrapString(s string, maxLen int) []string {
	if len(s) <= maxLen {
		return []string{s}
	}

	var lines []string
	for i := 0; i < len(s); i += maxLen {
		end := i + maxLen
		if end > len(s) {
			end = len(s)
		}
		lines = append(lines, s[i:end])
	}
	return lines
}

// formatBinaryData formats binary data in a readable way
func formatBinaryData(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Convert to hex representation
	hexStr := ""
	for _, b := range data {
		hexStr += fmt.Sprintf("%02x", b)
	}

	return hexStr
}

// padString pads a string to a specified length
func padString(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// centerString centers a string in a field of the specified length
func centerString(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}

	leftPad := (length - len(s)) / 2
	rightPad := length - len(s) - leftPad

	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// displayCurrentPacket displays the current packet in the console
func (pd *PacketDebugger) displayCurrentPacket() {
	if pd.selectedDump == nil || pd.currentPacket >= len(pd.selectedDump.Packets) {
		return
	}

	packet := pd.selectedDump.Packets[pd.currentPacket]

	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Print header
	fmt.Printf("=== Packet %d/%d: %s (%s) - %s ===\n\n",
		pd.currentPacket+1, len(pd.selectedDump.Packets), packet.PacketID, packet.Description, packet.Timestamp)

	// Get packet definition and decoded fields
	def, exists := getPacketDefinition(packet.PacketID, packet.Direction)
	var decodedFields map[string]interface{}

	if exists {
		rawBytes := rawDataToBytes(packet.RawData)
		decodedFields = decodePacketFields(rawBytes, def.Format, def.FieldNames)
	}

	// Print column headers
	fmt.Println("┌─" + strings.Repeat("─", colWidth) + "─┬─" + strings.Repeat("─", colWidth) + "─┬─" + strings.Repeat("─", arrowColWidth-2) + "─┬─" + strings.Repeat("─", colWidth) + "─┬─" + strings.Repeat("─", colWidth) + "─┐")
	fmt.Println("│ " + padString("CLIENT (Decoded)", colWidth) + " │ " + padString("SEND PACKETS", colWidth) + " │ " + padString("DIR", arrowColWidth-2) + " │ " + padString("RECEIVE PACKETS", colWidth) + " │ " + padString("SERVER (Encoded)", colWidth) + " │")
	fmt.Println("├─" + strings.Repeat("─", colWidth) + "─┼─" + strings.Repeat("─", colWidth) + "─┼─" + strings.Repeat("─", arrowColWidth-2) + "─┼─" + strings.Repeat("─", colWidth) + "─┼─" + strings.Repeat("─", colWidth) + "─┤")

	// Display packet information based on direction
	if packet.Direction == "sent" {
		// For sent packets, display in columns 2 and 4
		pd.displaySentPacket(packet, decodedFields, def, exists)
	} else if packet.Direction == "received" {
		// For received packets, display in columns 1 and 3
		pd.displayReceivedPacket(packet, decodedFields, def, exists)
	}

	fmt.Println("\nCommands: n=next packet, p=previous packet, q=quit")
}

// displaySentPacket displays a sent packet in the console
func (pd *PacketDebugger) displaySentPacket(packet common.Packet, decodedFields map[string]interface{}, def networkCommon.PacketConstruction, exists bool) {
	// Print packet info
	fmt.Println("│ " + padString("", colWidth) + " │ " + padString("ID: "+packet.PacketID, colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString("ID: "+packet.PacketID, colWidth) + " │")

	// Handle description and format wrapping
	descLines := wrapString("Desc: "+packet.Description, colWidth)
	formatLines := wrapString("Format: "+def.Format, colWidth)

	// Print first line of description and format
	fmt.Println("│ " + padString("", colWidth) + " │ " + padString(descLines[0], colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString(formatLines[0], colWidth) + " │")

	// Print any additional lines
	maxLines := len(descLines)
	if len(formatLines) > maxLines {
		maxLines = len(formatLines)
	}

	for i := 1; i < maxLines; i++ {
		descLine := ""
		if i < len(descLines) {
			descLine = descLines[i]
		}

		formatLine := ""
		if i < len(formatLines) {
			formatLine = formatLines[i]
		}

		fmt.Println("│ " + padString("", colWidth) + " │ " + padString(descLine, colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString(formatLine, colWidth) + " │")
	}

	// Print decoded fields
	fmt.Println("│ " + padString("", colWidth) + " │ " + padString("--- Decoded Fields ---", colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString("--- Encoded Fields ---", colWidth) + " │")

	if exists {
		// Print field values
		lineCount := 0
		for name, value := range decodedFields {
			var valueStr string

			// Format binary data specially
			switch v := value.(type) {
			case string:
				if len(v) > 0 && !isPrintable(v) {
					valueStr = formatBinaryData([]byte(v))
				} else {
					valueStr = v
				}
			default:
				valueStr = fmt.Sprintf("%v", value)
			}

			// Wrap long values
			fieldLines := wrapString(name+": "+valueStr, colWidth)

			for _, line := range fieldLines {
				fmt.Println("│ " + padString("", colWidth) + " │ " + padString(line, colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString(line, colWidth) + " │")
				lineCount++
				if lineCount > 15 {
					fmt.Println("│ " + padString("", colWidth) + " │ " + padString("...", colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString("...", colWidth) + " │")
					break
				}
			}

			if lineCount > 15 {
				break
			}
		}
	} else {
		// Wrap "No packet definition found" message
		noDefLines := wrapString("No packet definition found", colWidth)
		for _, line := range noDefLines {
			fmt.Println("│ " + padString("", colWidth) + " │ " + padString(line, colWidth) + " │ " + centerString("→", arrowColWidth) + " │ " + padString("", colWidth) + " │ " + padString(line, colWidth) + " │")
		}
	}

	// Close table
	fmt.Println("└─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", arrowColWidth-2) + "─┴─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", colWidth) + "─┘")
}

// displayReceivedPacket displays a received packet in the console
func (pd *PacketDebugger) displayReceivedPacket(packet common.Packet, decodedFields map[string]interface{}, def networkCommon.PacketConstruction, exists bool) {
	// Print packet info
	fmt.Println("│ " + padString("ID: "+packet.PacketID, colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString("ID: "+packet.PacketID, colWidth) + " │ " + padString("", colWidth) + " │")

	// Handle format and description wrapping
	formatLines := wrapString("Format: "+def.Format, colWidth)
	descLines := wrapString("Desc: "+packet.Description, colWidth)

	// Print first line of format and description
	fmt.Println("│ " + padString(formatLines[0], colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString(descLines[0], colWidth) + " │ " + padString("", colWidth) + " │")

	// Print any additional lines
	maxLines := len(formatLines)
	if len(descLines) > maxLines {
		maxLines = len(descLines)
	}

	for i := 1; i < maxLines; i++ {
		formatLine := ""
		if i < len(formatLines) {
			formatLine = formatLines[i]
		}

		descLine := ""
		if i < len(descLines) {
			descLine = descLines[i]
		}

		fmt.Println("│ " + padString(formatLine, colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString(descLine, colWidth) + " │ " + padString("", colWidth) + " │")
	}

	// Print decoded fields
	fmt.Println("│ " + padString("--- Decoded Fields ---", colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString("--- Raw Fields ---", colWidth) + " │ " + padString("", colWidth) + " │")

	if exists {
		// Print field values
		lineCount := 0
		for name, value := range decodedFields {
			var valueStr string

			// Format binary data specially
			switch v := value.(type) {
			case string:
				if len(v) > 0 && !isPrintable(v) {
					valueStr = formatBinaryData([]byte(v))
				} else {
					valueStr = v
				}
			default:
				valueStr = fmt.Sprintf("%v", value)
			}

			// Wrap long values
			fieldLines := wrapString(name+": "+valueStr, colWidth)

			for _, line := range fieldLines {
				fmt.Println("│ " + padString(line, colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString(line, colWidth) + " │ " + padString("", colWidth) + " │")
				lineCount++
				if lineCount > 15 {
					fmt.Println("│ " + padString("...", colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString("...", colWidth) + " │ " + padString("", colWidth) + " │")
					break
				}
			}

			if lineCount > 15 {
				break
			}
		}
	} else {
		// Wrap "No packet definition found" message
		noDefLines := wrapString("No packet definition found", colWidth)
		for _, line := range noDefLines {
			fmt.Println("│ " + padString(line, colWidth) + " │ " + padString("", colWidth) + " │ " + centerString("←", arrowColWidth) + " │ " + padString(line, colWidth) + " │ " + padString("", colWidth) + " │")
		}
	}

	// Close table
	fmt.Println("└─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", arrowColWidth-2) + "─┴─" + strings.Repeat("─", colWidth) + "─┴─" + strings.Repeat("─", colWidth) + "─┘")
}

// hexBytesToBytes converts a slice of hex strings to a byte slice
func hexBytesToBytes(hexBytes []string) []byte {
	bytes := make([]byte, len(hexBytes))
	for i, hex := range hexBytes {
		b, _ := strconv.ParseUint(hex, 16, 8)
		bytes[i] = byte(b)
	}
	return bytes
}

// bytesToHexStrings converts a byte slice to a slice of hex strings
func bytesToHexStrings(bytes []byte) []string {
	hexStrings := make([]string, len(bytes))
	for i, b := range bytes {
		hexStrings[i] = fmt.Sprintf("%02x", b)
	}
	return hexStrings
}

// rawDataToBytes converts raw data to a byte slice
func rawDataToBytes(rawData []common.RawData) []byte {
	var allBytes []byte
	for _, data := range rawData {
		allBytes = append(allBytes, hexBytesToBytes(data.HexBytes)...)
	}
	return allBytes
}

// formatRawData formats the raw data of a packet for display
func formatRawData(rawData []common.RawData) string {
	var result strings.Builder

	for i, data := range rawData {
		result.WriteString(fmt.Sprintf("Chunk %d:\n", i+1))

		// Format hex bytes in rows of 16
		for j := 0; j < len(data.HexBytes); j += 16 {
			end := j + 16
			if end > len(data.HexBytes) {
				end = len(data.HexBytes)
			}

			// Write hex values
			for k := j; k < end; k++ {
				result.WriteString(data.HexBytes[k] + " ")
			}

			result.WriteString("\n")
		}

		result.WriteString("\n")
	}

	return result.String()
}

// getPacketDefinition gets the packet definition based on direction
func getPacketDefinition(packetID string, direction string) (networkCommon.PacketConstruction, bool) {
	if direction == "sent" {
		// For sent packets, try to get from send servers
		sendDefs := sendServers.ServerType0PacketConstructions()
		if def, exists := sendDefs[packetID]; exists {
			return def, true
		}
	} else {
		// For received packets, get from receive servers
		receiveDefs := receiveServers.ServerType0PacketConstructions()
		if def, exists := receiveDefs[packetID]; exists {
			return def, true
		}
	}

	// If not found in the appropriate factory, try the other one as fallback
	if direction == "sent" {
		receiveDefs := receiveServers.ServerType0PacketConstructions()
		def, exists := receiveDefs[packetID]
		return def, exists
	} else {
		sendDefs := sendServers.ServerType0PacketConstructions()
		def, exists := sendDefs[packetID]
		return def, exists
	}
}

// decodePacketFields decodes packet fields based on format string
func decodePacketFields(rawBytes []byte, format string, fieldNames []string) map[string]interface{} {
	fields := make(map[string]interface{})

	// This is a simplified decoder that handles some basic format specifiers
	// In a real implementation, you would use a more robust decoder

	offset := 0
	fieldIndex := 0

	for i := 0; i < len(format) && fieldIndex < len(fieldNames); i++ {
		if offset >= len(rawBytes) {
			break
		}

		switch format[i] {
		case 'C': // unsigned char (1 byte)
			if offset < len(rawBytes) {
				fields[fieldNames[fieldIndex]] = int(rawBytes[offset])
				offset += 1
				fieldIndex++
			}
		case 'v': // unsigned short (2 bytes, little-endian)
			if offset+1 < len(rawBytes) {
				value := uint16(rawBytes[offset]) | uint16(rawBytes[offset+1])<<8
				fields[fieldNames[fieldIndex]] = int(value)
				offset += 2
				fieldIndex++
			}
		case 'V': // unsigned int (4 bytes, little-endian)
			if offset+3 < len(rawBytes) {
				value := uint32(rawBytes[offset]) |
					uint32(rawBytes[offset+1])<<8 |
					uint32(rawBytes[offset+2])<<16 |
					uint32(rawBytes[offset+3])<<24
				fields[fieldNames[fieldIndex]] = int(value)
				offset += 4
				fieldIndex++
			}
		case 'a': // fixed-length string
			// Look for a number after 'a'
			j := i + 1
			numStr := ""
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				numStr += string(format[j])
				j++
			}

			if numStr != "" {
				length, _ := strconv.Atoi(numStr)
				if offset+length <= len(rawBytes) {
					// Extract the string and trim null bytes
					str := string(rawBytes[offset : offset+length])
					str = strings.TrimRight(str, "\x00")
					fields[fieldNames[fieldIndex]] = str
					offset += length
					fieldIndex++
				}
				i = j - 1 // Skip the digits we've processed
			}
		case 'Z': // null-terminated string
			// Find the null terminator
			end := offset
			for end < len(rawBytes) && rawBytes[end] != 0 {
				end++
			}

			if end < len(rawBytes) {
				str := string(rawBytes[offset:end])
				fields[fieldNames[fieldIndex]] = str
				offset = end + 1
				fieldIndex++
			}
		case ' ', '*', 'x': // Ignore spaces, wildcards, and padding
			continue
		}
	}

	return fields
}

// nextPacket moves to the next packet
func (pd *PacketDebugger) nextPacket() bool {
	if pd.selectedDump == nil {
		return false
	}

	pd.currentPacket++
	if pd.currentPacket >= len(pd.selectedDump.Packets) {
		pd.currentPacket = 0
	}

	pd.displayCurrentPacket()
	return true
}

// previousPacket moves to the previous packet
func (pd *PacketDebugger) previousPacket() bool {
	if pd.selectedDump == nil {
		return false
	}

	pd.currentPacket--
	if pd.currentPacket < 0 {
		pd.currentPacket = len(pd.selectedDump.Packets) - 1
	}

	pd.displayCurrentPacket()
	return true
}

// Run starts the packet debugger
func (pd *PacketDebugger) Run(dumpDir string) error {
	// Load dump files
	if err := pd.loadDumpFiles(dumpDir); err != nil {
		return err
	}

	// Display available dump files
	fmt.Println("Available packet dumps:")
	for i, file := range pd.dumpFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	// Get user selection
	fmt.Print("\nSelect a dump file (enter number): ")
	input, _ := pd.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(pd.dumpFiles) {
		return fmt.Errorf("invalid selection")
	}

	// Load selected dump
	dumpPath := filepath.Join(dumpDir, pd.dumpFiles[index-1])
	if err := pd.loadDump(dumpPath); err != nil {
		return err
	}

	// Display first packet
	pd.displayCurrentPacket()

	// Process user input
	fmt.Println("\nPress Enter for next packet, p for previous packet, q to quit")

	// Set terminal to raw mode to capture single keystrokes
	oldState, err := makeRaw(0)
	if err != nil {
		return fmt.Errorf("could not set terminal to raw mode: %v", err)
	}
	defer restoreTerminal(0, oldState)

	for {
		// Read a single byte
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return err
		}

		switch b[0] {
		case 'q', 'Q':
			return nil
		case 'n', 'N', '\r', '\n': // n, N, or Enter for next packet
			pd.nextPacket()
		case 'p', 'P': // p or P for previous packet
			pd.previousPacket()
		}
	}
}

// isPrintable checks if a string contains only printable ASCII characters
func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func main() {
	// Define paths
	dumpDir := "/home/mikekao/personal/openkore/goKore/verification/PacketAnalysis/extracteddata"

	// Create and run the packet debugger
	pd := NewPacketDebugger()
	if err := pd.Run(dumpDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
