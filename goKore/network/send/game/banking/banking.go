// Package banking provides banking-related packet sending functionality.
package banking

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// BankingManager handles banking-related packet sending.
type BankingManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewBankingManager creates a new banking manager.
func NewBankingManager(baseSend core.Send) *BankingManager {
	return &BankingManager{
		baseSend: baseSend,
	}
}

// SendBankingCheck sends a request to check banking data.
// This is equivalent to the sendBankingCheck function in Send.pm.
func (bm *BankingManager) SendBankingCheck(accountID uint32) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("banking_check_request")
	if !exists {
		return fmt.Errorf("banking_check_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendBankingWithdraw sends a request to withdraw zeny from the bank.
// This is equivalent to the sendBankingWithdraw function in Send.pm.
func (bm *BankingManager) SendBankingWithdraw(accountID, zeny uint32) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("banking_withdraw_request")
	if !exists {
		return fmt.Errorf("banking_withdraw_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"zeny":      zeny,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendBankingDeposit sends a request to deposit zeny to the bank.
// This is equivalent to the sendBankingDeposit function in Send.pm.
func (bm *BankingManager) SendBankingDeposit(accountID, zeny uint32) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("banking_deposit_request")
	if !exists {
		return fmt.Errorf("banking_deposit_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"zeny":      zeny,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}
