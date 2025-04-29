// Package macro provides macro-related packet sending functionality.
package macro

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MacroManager handles macro-related packet sending.
type MacroManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMacroManager creates a new macro manager.
func NewMacroManager(baseSend core.Send) *MacroManager {
	return &MacroManager{
		baseSend: baseSend,
	}
}

// SendMacroStart sends a request to start a macro.
// This is equivalent to the sendMacroStart function in Send.pm.
func (mm *MacroManager) SendMacroStart() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("macro_start")
	if !exists {
		return fmt.Errorf("macro_start packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMacroStop sends a request to stop a macro.
// This is equivalent to the sendMacroStop function in Send.pm.
func (mm *MacroManager) SendMacroStop() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("macro_stop")
	if !exists {
		return fmt.Errorf("macro_stop packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMacroDetectorDownload sends a request to download the macro detector.
// This is equivalent to the sendMacroDetectorDownload function in Send.pm.
func (mm *MacroManager) SendMacroDetectorDownload() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("macro_detector_download")
	if !exists {
		return fmt.Errorf("macro_detector_download packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMacroDetectorAnswer sends an answer to the macro detector.
// This is equivalent to the sendMacroDetectorAnswer function in Send.pm.
func (mm *MacroManager) SendMacroDetectorAnswer(answer string) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("macro_detector_answer")
	if !exists {
		return fmt.Errorf("macro_detector_answer packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"answer": []byte(answer),
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendReqCashTabCode sends a request for a cash tab code.
// This is equivalent to the sendReqCashTabCode function in Send.pm.
func (mm *MacroManager) SendReqCashTabCode(tabID uint16) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("req_cash_tabcode")
	if !exists {
		return fmt.Errorf("req_cash_tabcode packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": tabID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
