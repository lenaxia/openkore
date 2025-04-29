package actor

import (
	"bytes"
	"testing"
)

// TestSendCharCreate tests the SendCharCreate function
func TestSendCharCreate(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["char_create"] = "0067"
	charManager := NewCharacterManager(mockSend)

	// Test with packet 0067 (full stats)
	err := charManager.SendCharCreate(0, "TestChar", 9, 9, 9, 9, 9, 9, 1, 1, 0, 0)
	if err != nil {
		t.Fatalf("SendCharCreate() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "char_create" {
		t.Errorf("Expected packet ID 'char_create', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["slot"].(int) != 0 {
		t.Errorf("Expected slot=0, got %v", args["slot"])
	}
	if args["name"].(string) != "TestChar" {
		t.Errorf("Expected name=TestChar, got %v", args["name"])
	}
	if args["str"].(int) != 9 {
		t.Errorf("Expected str=9, got %v", args["str"])
	}
	if args["agi"].(int) != 9 {
		t.Errorf("Expected agi=9, got %v", args["agi"])
	}
	if args["vit"].(int) != 9 {
		t.Errorf("Expected vit=9, got %v", args["vit"])
	}
	if args["int"].(int) != 9 {
		t.Errorf("Expected int=9, got %v", args["int"])
	}
	if args["dex"].(int) != 9 {
		t.Errorf("Expected dex=9, got %v", args["dex"])
	}
	if args["luk"].(int) != 9 {
		t.Errorf("Expected luk=9, got %v", args["luk"])
	}
	if args["hair_style"].(int) != 1 {
		t.Errorf("Expected hair_style=1, got %v", args["hair_style"])
	}
	if args["hair_color"].(int) != 1 {
		t.Errorf("Expected hair_color=1, got %v", args["hair_color"])
	}

	// Test with packet 0970 (simplified)
	mockSend.packetIDs["char_create"] = "0970"
	err = charManager.SendCharCreate(1, "TestChar2", 0, 0, 0, 0, 0, 0, 2, 2, 0, 0)
	if err != nil {
		t.Fatalf("SendCharCreate() returned error: %v", err)
	}

	// Check that the correct arguments were used
	args = mockSend.LastArgs()
	if args["slot"].(int) != 1 {
		t.Errorf("Expected slot=1, got %v", args["slot"])
	}
	if args["name"].(string) != "TestChar2" {
		t.Errorf("Expected name=TestChar2, got %v", args["name"])
	}
	if args["hair_style"].(int) != 2 {
		t.Errorf("Expected hair_style=2, got %v", args["hair_style"])
	}
	if args["hair_color"].(int) != 2 {
		t.Errorf("Expected hair_color=2, got %v", args["hair_color"])
	}

	// Test with packet 0A39 (with job_id and sex)
	mockSend.packetIDs["char_create"] = "0A39"
	err = charManager.SendCharCreate(2, "TestChar3", 0, 0, 0, 0, 0, 0, 3, 3, 1, 1)
	if err != nil {
		t.Fatalf("SendCharCreate() returned error: %v", err)
	}

	// Check that the correct arguments were used
	args = mockSend.LastArgs()
	if args["slot"].(int) != 2 {
		t.Errorf("Expected slot=2, got %v", args["slot"])
	}
	if args["name"].(string) != "TestChar3" {
		t.Errorf("Expected name=TestChar3, got %v", args["name"])
	}
	if args["hair_style"].(int) != 3 {
		t.Errorf("Expected hair_style=3, got %v", args["hair_style"])
	}
	if args["hair_color"].(int) != 3 {
		t.Errorf("Expected hair_color=3, got %v", args["hair_color"])
	}
	if args["job_id"].(int) != 1 {
		t.Errorf("Expected job_id=1, got %v", args["job_id"])
	}
	if args["sex"].(int) != 1 {
		t.Errorf("Expected sex=1, got %v", args["sex"])
	}
}

// TestSendCharDelete tests the SendCharDelete function
func TestSendCharDelete(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["char_delete"] = "0068"
	charManager := NewCharacterManager(mockSend)

	charID := []byte{1, 2, 3, 4}
	email := "test@example.com"

	err := charManager.SendCharDelete(charID, email)
	if err != nil {
		t.Fatalf("SendCharDelete() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "char_delete" {
		t.Errorf("Expected packet ID 'char_delete', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if !bytes.Equal(args["charID"].([]byte), charID) {
		t.Errorf("Expected charID=%v, got %v", charID, args["charID"])
	}
	if args["email"].(string) != email {
		t.Errorf("Expected email=%s, got %v", email, args["email"])
	}
}

// TestSendCharDelete2 tests the SendCharDelete2 function
func TestSendCharDelete2(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["char_delete2"] = "0827"
	charManager := NewCharacterManager(mockSend)

	charID := []byte{1, 2, 3, 4}

	err := charManager.SendCharDelete2(charID)
	if err != nil {
		t.Fatalf("SendCharDelete2() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "char_delete2" {
		t.Errorf("Expected packet ID 'char_delete2', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if !bytes.Equal(args["charID"].([]byte), charID) {
		t.Errorf("Expected charID=%v, got %v", charID, args["charID"])
	}
}

// TestSendCharDelete2Accept tests the SendCharDelete2Accept function
func TestSendCharDelete2Accept(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["char_delete2_accept"] = "0829"
	charManager := NewCharacterManager(mockSend)

	charID := []byte{1, 2, 3, 4}
	code := "123456"

	err := charManager.SendCharDelete2Accept(charID, code)
	if err != nil {
		t.Fatalf("SendCharDelete2Accept() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "char_delete2_accept" {
		t.Errorf("Expected packet ID 'char_delete2_accept', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if !bytes.Equal(args["charID"].([]byte), charID) {
		t.Errorf("Expected charID=%v, got %v", charID, args["charID"])
	}
	if args["code"].(string) != code {
		t.Errorf("Expected code=%s, got %v", code, args["code"])
	}
}

// TestSendCharDelete2Cancel tests the SendCharDelete2Cancel function
func TestSendCharDelete2Cancel(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetIDs["char_delete2_cancel"] = "082B"
	charManager := NewCharacterManager(mockSend)

	charID := []byte{1, 2, 3, 4}

	err := charManager.SendCharDelete2Cancel(charID)
	if err != nil {
		t.Fatalf("SendCharDelete2Cancel() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "char_delete2_cancel" {
		t.Errorf("Expected packet ID 'char_delete2_cancel', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if !bytes.Equal(args["charID"].([]byte), charID) {
		t.Errorf("Expected charID=%v, got %v", charID, args["charID"])
	}
}
