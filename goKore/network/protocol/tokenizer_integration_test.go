package protocol

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTokenizerWithLoadedPacketDefs(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test recvpackets.txt file
	testFilePath := filepath.Join(tempDir, "recvpackets.txt")
	testContent := `# Test recvpackets.txt
0064 55
0065 17
0066 3
0069 -1
006A 23
006B -1
0073 11
`
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Parse recvpackets.txt
	packetDefs, err := ParseRecvPackets(testFilePath)
	if err != nil {
		t.Fatalf("Failed to parse recvpackets.txt: %v", err)
	}

	// Create tokenizer with loaded packet definitions
	tokenizer := NewTokenizer(packetDefs)

	// Test fixed-length packet
	t.Run("Fixed-length packet", func(t *testing.T) {
		// Create a test packet with ID 0064 (length 55)
		packet := make([]byte, 55)
		packet[0] = 0x64 // Little-endian byte order
		packet[1] = 0x00

		// Add packet to tokenizer
		tokenizer.Add(packet)

		// Read packet
		result, msgType, err := tokenizer.ReadNext()
		if err != nil {
			t.Errorf("ReadNext() error = %v", err)
			return
		}

		if msgType != KnownMessage {
			t.Errorf("ReadNext() msgType = %v, want %v", msgType, KnownMessage)
		}

		if len(result) != 55 {
			t.Errorf("ReadNext() result length = %d, want %d", len(result), 55)
		}

		// Buffer should be empty now
		if len(tokenizer.GetBuffer()) != 0 {
			t.Errorf("Buffer not empty after reading packet")
		}
	})

	// Test variable-length packet
	t.Run("Variable-length packet", func(t *testing.T) {
		// Create a test packet with ID 0069 (variable length)
		packetLength := 20
		packet := make([]byte, packetLength)
		packet[0] = 0x69 // Little-endian byte order
		packet[1] = 0x00
		packet[2] = byte(packetLength) // Length bytes
		packet[3] = 0x00

		// Add packet to tokenizer
		tokenizer.Add(packet)

		// Read packet
		result, msgType, err := tokenizer.ReadNext()
		if err != nil {
			t.Errorf("ReadNext() error = %v", err)
			return
		}

		if msgType != KnownMessage {
			t.Errorf("ReadNext() msgType = %v, want %v", msgType, KnownMessage)
		}

		if len(result) != packetLength {
			t.Errorf("ReadNext() result length = %d, want %d", len(result), packetLength)
		}

		// Buffer should be empty now
		if len(tokenizer.GetBuffer()) != 0 {
			t.Errorf("Buffer not empty after reading packet")
		}
	})

	// Test incomplete packet
	t.Run("Incomplete packet", func(t *testing.T) {
		// Create a partial packet with ID 0064 (length 55)
		packet := make([]byte, 30) // Only 30 bytes of a 55-byte packet
		packet[0] = 0x64           // Little-endian byte order
		packet[1] = 0x00

		// Add packet to tokenizer
		tokenizer.Add(packet)

		// Try to read packet
		_, _, err := tokenizer.ReadNext()
		if err != ErrIncompletePacket {
			t.Errorf("ReadNext() error = %v, want %v", err, ErrIncompletePacket)
		}

		// Buffer should still contain the partial packet
		if len(tokenizer.GetBuffer()) != 30 {
			t.Errorf("Buffer length = %d, want %d", len(tokenizer.GetBuffer()), 30)
		}

		// Clear the buffer
		tokenizer.Clear(30)
	})

	// Test unknown packet
	t.Run("Unknown packet", func(t *testing.T) {
		// Create a packet with unknown ID 0099
		packet := make([]byte, 10)
		packet[0] = 0x99 // Little-endian byte order
		packet[1] = 0x00

		// Add packet to tokenizer
		tokenizer.Add(packet)

		// Read packet
		result, msgType, err := tokenizer.ReadNext()
		if err != nil {
			t.Errorf("ReadNext() error = %v", err)
			return
		}

		if msgType != UnknownMessage {
			t.Errorf("ReadNext() msgType = %v, want %v", msgType, UnknownMessage)
		}

		if len(result) != 10 {
			t.Errorf("ReadNext() result length = %d, want %d", len(result), 10)
		}

		// Buffer should be empty now
		if len(tokenizer.GetBuffer()) != 0 {
			t.Errorf("Buffer not empty after reading packet")
		}
	})

	// Test multiple packets
	t.Run("Multiple packets", func(t *testing.T) {
		// Create two packets: 0064 (length 55) and 0065 (length 17)
		packet1 := make([]byte, 55)
		packet1[0] = 0x64 // Little-endian byte order
		packet1[1] = 0x00

		packet2 := make([]byte, 17)
		packet2[0] = 0x65 // Little-endian byte order
		packet2[1] = 0x00

		// Add both packets to tokenizer
		tokenizer.Add(packet1)
		tokenizer.Add(packet2)

		// Read first packet
		result1, msgType1, err := tokenizer.ReadNext()
		if err != nil {
			t.Errorf("ReadNext() error = %v", err)
			return
		}

		if msgType1 != KnownMessage {
			t.Errorf("ReadNext() msgType = %v, want %v", msgType1, KnownMessage)
		}

		if len(result1) != 55 {
			t.Errorf("ReadNext() result length = %d, want %d", len(result1), 55)
		}

		// Read second packet
		result2, msgType2, err := tokenizer.ReadNext()
		if err != nil {
			t.Errorf("ReadNext() error = %v", err)
			return
		}

		if msgType2 != KnownMessage {
			t.Errorf("ReadNext() msgType = %v, want %v", msgType2, KnownMessage)
		}

		if len(result2) != 17 {
			t.Errorf("ReadNext() result length = %d, want %d", len(result2), 17)
		}

		// Buffer should be empty now
		if len(tokenizer.GetBuffer()) != 0 {
			t.Errorf("Buffer not empty after reading packets")
		}
	})
}
