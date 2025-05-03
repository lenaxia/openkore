package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lenaxia/goKore/cmd/common"
)

func TestPacketDumpLoading(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_dump.json")

	// Create test data
	testDump := common.PacketDump{
		FileName:    "TEST_DUMP",
		PacketCount: 2,
		Packets: []common.Packet{
			{
				Direction:   "sent",
				PacketID:    "0064",
				Description: "Account Server Login",
				Size:        55,
				Timestamp:   "2025.04.30 16:24:58",
				RawData: []common.RawData{
					{
						HexBytes: []string{
							"64", "00", "1C", "00", "00", "00",
							"62", "6F", "74", "69", "6A", "6F",
						},
						AsciiRepresentation: "",
						BinaryBase64:        "ZAAcAAAAYm90aWpv",
					},
				},
				ServerType: "Account Server",
			},
			{
				Direction:   "received",
				PacketID:    "0AC4",
				Description: "Account Info With Server Info",
				Size:        224,
				Timestamp:   "2025.04.30 16:24:58",
				RawData: []common.RawData{
					{
						HexBytes: []string{
							"C4", "0A", "E0", "00", "E5", "5D",
							"F6", "C1", "82", "84",
						},
						AsciiRepresentation: "",
						BinaryBase64:        "xADgAOVd9sGChA==",
					},
				},
				ServerType: "Account Server",
			},
		},
	}

	// Write test data to file
	data, err := json.MarshalIndent(testDump, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(testFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create packet debugger
	pd := NewPacketDebugger()

	// Test loadDumpFiles
	if err := pd.loadDumpFiles(tempDir); err != nil {
		t.Fatalf("Failed to load dump files: %v", err)
	}
	if len(pd.dumpFiles) != 1 {
		t.Errorf("Expected 1 dump file, got %d", len(pd.dumpFiles))
	}

	// Test loadDump
	if err := pd.loadDump(testFile); err != nil {
		t.Fatalf("Failed to load dump: %v", err)
	}
	if pd.selectedDump == nil {
		t.Fatal("Selected dump is nil")
	}
	if pd.selectedDump.PacketCount != 2 {
		t.Errorf("Expected packet count 2, got %d", pd.selectedDump.PacketCount)
	}
	if len(pd.selectedDump.Packets) != 2 {
		t.Errorf("Expected 2 packets, got %d", len(pd.selectedDump.Packets))
	}

	// Test packet data
	if pd.selectedDump.Packets[0].PacketID != "0064" {
		t.Errorf("Expected packet ID 0064, got %s", pd.selectedDump.Packets[0].PacketID)
	}
	if pd.selectedDump.Packets[1].PacketID != "0AC4" {
		t.Errorf("Expected packet ID 0AC4, got %s", pd.selectedDump.Packets[1].PacketID)
	}

	// Test nextPacket
	initialPacket := pd.currentPacket
	pd.nextPacket()
	if pd.currentPacket != initialPacket+1 {
		t.Errorf("Expected current packet to be %d, got %d", initialPacket+1, pd.currentPacket)
	}
}

func TestHexBytesConversion(t *testing.T) {
	// Test hexBytesToBytes
	hexBytes := []string{"64", "00", "1C", "FF"}
	bytes := hexBytesToBytes(hexBytes)
	if len(bytes) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(bytes))
	}
	if bytes[0] != 0x64 || bytes[1] != 0x00 || bytes[2] != 0x1C || bytes[3] != 0xFF {
		t.Errorf("Incorrect byte conversion: %v", bytes)
	}

	// Test bytesToHexStrings
	hexStrings := bytesToHexStrings(bytes)
	if len(hexStrings) != 4 {
		t.Errorf("Expected 4 hex strings, got %d", len(hexStrings))
	}
	if hexStrings[0] != "64" || hexStrings[1] != "00" || hexStrings[2] != "1c" || hexStrings[3] != "ff" {
		t.Errorf("Incorrect hex string conversion: %v", hexStrings)
	}
}

func TestRawDataFormatting(t *testing.T) {
	// Test formatRawData
	rawData := []common.RawData{
		{
			HexBytes: []string{
				"64", "00", "1C", "00", "00", "00",
				"62", "6F", "74", "69", "6A", "6F",
				"30", "00", "00", "00", "00", "00",
			},
			AsciiRepresentation: "",
			BinaryBase64:        "ZAAcAAAAYm90aWpvMAAAAA==",
		},
	}

	formatted := formatRawData(rawData)
	if formatted == "" {
		t.Error("Expected non-empty formatted raw data")
	}
}
