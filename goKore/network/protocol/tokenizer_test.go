package protocol

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func TestGetMessageID(t *testing.T) {
	tests := []struct {
		name   string
		packet []byte
		want   string
	}{
		{
			name:   "Valid packet",
			packet: []byte{0x01, 0x02, 0x03, 0x04},
			want:   "0201",
		},
		{
			name:   "Short packet",
			packet: []byte{0x01},
			want:   "",
		},
		{
			name:   "Empty packet",
			packet: []byte{},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMessageID(tt.packet); got != tt.want {
				t.Errorf("GetMessageID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizerAdd(t *testing.T) {
	tokenizer := NewTokenizer(nil)

	// Test adding data
	data1 := []byte{0x01, 0x02, 0x03}
	tokenizer.Add(data1)
	if !reflect.DeepEqual(tokenizer.buffer, data1) {
		t.Errorf("Add() failed, buffer = %v, want %v", tokenizer.buffer, data1)
	}

	// Test adding more data
	data2 := []byte{0x04, 0x05}
	tokenizer.Add(data2)
	want := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	if !reflect.DeepEqual(tokenizer.buffer, want) {
		t.Errorf("Add() failed, buffer = %v, want %v", tokenizer.buffer, want)
	}
}

func TestTokenizerClear(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	tokenizer.Add(data)

	tests := []struct {
		name       string
		size       int
		wantBuffer []byte
	}{
		{
			name:       "Clear partial",
			size:       3,
			wantBuffer: []byte{0x04, 0x05},
		},
		{
			name:       "Clear all",
			size:       5,
			wantBuffer: []byte{},
		},
		{
			name:       "Clear with size > buffer length",
			size:       10,
			wantBuffer: []byte{},
		},
		{
			name:       "Clear with size <= 0",
			size:       0,
			wantBuffer: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer.buffer = make([]byte, len(data))
			copy(tokenizer.buffer, data)

			tokenizer.Clear(tt.size)
			if !reflect.DeepEqual(tokenizer.buffer, tt.wantBuffer) {
				t.Errorf("Clear(%d) = %v, want %v", tt.size, tokenizer.buffer, tt.wantBuffer)
			}
		})
	}
}

func TestTokenizerGetBuffer(t *testing.T) {
	tokenizer := NewTokenizer(nil)
	data := []byte{0x01, 0x02, 0x03}
	tokenizer.Add(data)

	if got := tokenizer.GetBuffer(); !reflect.DeepEqual(got, data) {
		t.Errorf("GetBuffer() = %v, want %v", got, data)
	}
}

func TestTokenizerNextMessageMightBeAccountID(t *testing.T) {
	tokenizer := NewTokenizer(nil)

	// Initially should be false
	if tokenizer.nextMightBeAccountID {
		t.Errorf("Initial nextMightBeAccountID = true, want false")
	}

	// Set to true
	tokenizer.NextMessageMightBeAccountID()
	if !tokenizer.nextMightBeAccountID {
		t.Errorf("After NextMessageMightBeAccountID(), nextMightBeAccountID = false, want true")
	}
}

func TestTokenizerReadNext(t *testing.T) {
	// Create packet definitions for testing
	packetDefs := map[string]PacketLengthDef{
		"0102": {Length: 6, HasLength: false}, // Fixed length packet
		"0304": {Length: -1, HasLength: true}, // Variable length packet
		"0506": {Length: 0, HasLength: false}, // Invalid packet definition
	}

	tokenizer := NewTokenizer(packetDefs)

	// Test 1: Empty buffer
	_, _, err := tokenizer.ReadNext()
	if err != ErrIncompletePacket {
		t.Errorf("ReadNext() with empty buffer, err = %v, want %v", err, ErrIncompletePacket)
	}

	// Test 2: Account ID message
	tokenizer.buffer = []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	tokenizer.NextMessageMightBeAccountID()
	packet, msgType, err := tokenizer.ReadNext()
	if err != nil {
		t.Errorf("ReadNext() with account ID, err = %v, want nil", err)
	}
	if msgType != AccountID {
		t.Errorf("ReadNext() with account ID, msgType = %v, want %v", msgType, AccountID)
	}
	if !reflect.DeepEqual(packet, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Errorf("ReadNext() with account ID, packet = %v, want %v", packet, []byte{0x01, 0x02, 0x03, 0x04})
	}
	if !reflect.DeepEqual(tokenizer.buffer, []byte{0x05}) {
		t.Errorf("ReadNext() with account ID, buffer = %v, want %v", tokenizer.buffer, []byte{0x05})
	}

	// Test 3: Fixed length packet
	tokenizer.buffer = []byte{0x02, 0x01, 0x03, 0x04, 0x05, 0x06, 0x07}
	packet, msgType, err = tokenizer.ReadNext()
	if err != nil {
		t.Errorf("ReadNext() with fixed length packet, err = %v, want nil", err)
	}
	if msgType != KnownMessage {
		t.Errorf("ReadNext() with fixed length packet, msgType = %v, want %v", msgType, KnownMessage)
	}
	if !reflect.DeepEqual(packet, []byte{0x02, 0x01, 0x03, 0x04, 0x05, 0x06}) {
		t.Errorf("ReadNext() with fixed length packet, packet = %v, want %v", packet, []byte{0x02, 0x01, 0x03, 0x04, 0x05, 0x06})
	}
	if !reflect.DeepEqual(tokenizer.buffer, []byte{0x07}) {
		t.Errorf("ReadNext() with fixed length packet, buffer = %v, want %v", tokenizer.buffer, []byte{0x07})
	}

	// Test 4: Variable length packet
	tokenizer.buffer = []byte{0x04, 0x03}
	_, _, err = tokenizer.ReadNext()
	if err != ErrIncompletePacket {
		t.Errorf("ReadNext() with incomplete variable length packet, err = %v, want %v", err, ErrIncompletePacket)
	}

	// Add length bytes
	tokenizer.buffer = []byte{0x04, 0x03, 0x08, 0x00} // Length = 8
	_, _, err = tokenizer.ReadNext()
	if err != ErrIncompletePacket {
		t.Errorf("ReadNext() with incomplete variable length packet data, err = %v, want %v", err, ErrIncompletePacket)
	}

	// Complete variable length packet
	tokenizer.buffer = []byte{0x04, 0x03, 0x08, 0x00, 0x05, 0x06, 0x07, 0x08}
	packet, msgType, err = tokenizer.ReadNext()
	if err != nil {
		t.Errorf("ReadNext() with variable length packet, err = %v, want nil", err)
	}
	if msgType != KnownMessage {
		t.Errorf("ReadNext() with variable length packet, msgType = %v, want %v", msgType, KnownMessage)
	}
	if !reflect.DeepEqual(packet, []byte{0x04, 0x03, 0x08, 0x00, 0x05, 0x06, 0x07, 0x08}) {
		t.Errorf("ReadNext() with variable length packet, packet = %v, want %v", packet, []byte{0x04, 0x03, 0x08, 0x00, 0x05, 0x06, 0x07, 0x08})
	}
	if len(tokenizer.buffer) != 0 {
		t.Errorf("ReadNext() with variable length packet, buffer = %v, want empty", tokenizer.buffer)
	}

	// Test 5: Unknown packet type
	tokenizer.buffer = []byte{0xFF, 0xFF, 0x01, 0x02}
	packet, msgType, err = tokenizer.ReadNext()
	if err != nil {
		t.Errorf("ReadNext() with unknown packet type, err = %v, want nil", err)
	}
	if msgType != UnknownMessage {
		t.Errorf("ReadNext() with unknown packet type, msgType = %v, want %v", msgType, UnknownMessage)
	}
	if !reflect.DeepEqual(packet, []byte{0xFF, 0xFF, 0x01, 0x02}) {
		t.Errorf("ReadNext() with unknown packet type, packet = %v, want %v", packet, []byte{0xFF, 0xFF, 0x01, 0x02})
	}

	// Test 6: Invalid packet definition
	tokenizer.buffer = []byte{0x06, 0x05, 0x01, 0x02}
	_, _, err = tokenizer.ReadNext()
	if err != ErrInvalidPacket {
		t.Errorf("ReadNext() with invalid packet definition, err = %v, want %v", err, ErrInvalidPacket)
	}
}

func TestVariableLengthPacketWithDifferentSizes(t *testing.T) {
	// Create packet definitions for testing
	packetDefs := map[string]PacketLengthDef{
		"0304": {Length: -1, HasLength: true}, // Variable length packet
	}

	tokenizer := NewTokenizer(packetDefs)

	// Test with different packet sizes
	sizes := []int{8, 16, 32, 64, 128}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// Create a packet with the specified size
			packet := make([]byte, size)
			packet[0] = 0x04 // Packet ID (little endian)
			packet[1] = 0x03
			binary.LittleEndian.PutUint16(packet[2:4], uint16(size)) // Set length

			// Fill the rest with test data
			for i := 4; i < size; i++ {
				packet[i] = byte(i % 256)
			}

			tokenizer.buffer = make([]byte, size)
			copy(tokenizer.buffer, packet)

			result, msgType, err := tokenizer.ReadNext()
			if err != nil {
				t.Errorf("ReadNext() with size %d, err = %v, want nil", size, err)
			}
			if msgType != KnownMessage {
				t.Errorf("ReadNext() with size %d, msgType = %v, want %v", size, msgType, KnownMessage)
			}
			if !reflect.DeepEqual(result, packet) {
				t.Errorf("ReadNext() with size %d, result doesn't match expected packet", size)
			}
		})
	}
}
