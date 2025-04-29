// Package security provides security-related packet sending functionality.
package security

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// PINManager handles PIN code-related packet sending.
type PINManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewPINManager creates a new PIN manager.
func NewPINManager(baseSend core.Send) *PINManager {
	return &PINManager{
		baseSend: baseSend,
	}
}

// SendPINCode sends a PIN code to the server.
func (pm *PINManager) SendPINCode(pin, seed int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("send_pin_code")
	if !exists {
		return fmt.Errorf("send_pin_code packet ID not found")
	}

	// Encode the PIN
	encodedPin := pm.baseSend.PinEncode(seed, pin)

	// Create the arguments
	args := map[string]interface{}{
		"pin": encodedPin,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPINCodeState sends the PIN code state to the server.
// State values:
// 0 = PIN system disabled
// 1 = PIN system enabled
// 2 = PIN code requested
func (pm *PINManager) SendPINCodeState(state int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("send_pin_state")
	if !exists {
		return fmt.Errorf("send_pin_state packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"state": state,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendLoginPinCode sends a PIN code to the server for login authentication.
// type: 0 = send_pin_password, 1 = new_pin_password
func (pm *PINManager) SendLoginPinCode(seed, pinType int) error {
	// Get the config PIN code
	// In the original Perl code, this is $config{loginPinCode}
	// For now, we'll use a hardcoded value of 1234
	// TODO: Get this from the config
	pin := 1234

	// Encode the PIN
	encodedPin := pm.baseSend.PinEncode(seed, pin)

	// Get the packet ID based on the type
	var packetName string
	switch pinType {
	case 0:
		packetName = "send_pin_password"
	case 1:
		packetName = "new_pin_password"
	default:
		return fmt.Errorf("invalid PIN type: %d", pinType)
	}

	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID(packetName)
	if !exists {
		return fmt.Errorf("%s packet ID not found", packetName)
	}

	// Create the arguments
	// In the original Perl code, this includes accountID
	// TODO: Get accountID from the global state
	args := map[string]interface{}{
		"pin": encodedPin,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPINCodeChange sends a request to change the PIN code.
func (pm *PINManager) SendPINCodeChange(oldPin, newPin, seed int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("send_pin_change")
	if !exists {
		return fmt.Errorf("send_pin_change packet ID not found")
	}

	// Encode the PINs
	encodedOldPin := pm.baseSend.PinEncode(seed, oldPin)
	encodedNewPin := pm.baseSend.PinEncode(seed, newPin)

	// Create the arguments
	args := map[string]interface{}{
		"old_pin": encodedOldPin,
		"new_pin": encodedNewPin,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}
