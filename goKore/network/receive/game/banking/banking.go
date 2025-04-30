package banking

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// BankingManager handles banking-related packet handling
type BankingManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for banking interactions
	bankingOpened bool
	bankingZeny   uint32
	charZeny      uint32
}

// NewBankingManager creates a new banking manager
func NewBankingManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *BankingManager {
	return &BankingManager{
		parser:        parser,
		hookManager:   hookManager,
		logger:        logger,
		bankingOpened: false,
		bankingZeny:   0,
		charZeny:      0,
	}
}

// RegisterHandlers registers all banking-related packet handlers
func (bm *BankingManager) RegisterHandlers() {
	// Register banking check handler
	bm.parser.RegisterHandlerFunc("09A6", "banking_check", "L",
		[]string{"zeny"}, bm.HandleBankingCheck)

	// Register banking deposit handler
	bm.parser.RegisterHandlerFunc("09A8", "banking_deposit", "W Q L",
		[]string{"reason", "money", "balance"}, bm.HandleBankingDeposit)

	// Register banking withdraw handler
	bm.parser.RegisterHandlerFunc("09AA", "banking_withdraw", "W Q L",
		[]string{"reason", "money", "balance"}, bm.HandleBankingWithdraw)
}

// HandleBankingCheck handles the banking_check packet (lines 11976-11988)
func (bm *BankingManager) HandleBankingCheck(args map[string]interface{}) error {
	// Extract packet data
	zeny, ok := args["zeny"].(uint32)
	if !ok {
		return fmt.Errorf("invalid zeny in banking_check packet")
	}

	// Update banking state
	bm.bankingOpened = true
	bm.bankingZeny = zeny

	// Log banking information
	bm.logger.Info("===== [Zeny Storage (Bank)] =====")
	bm.logger.Info("In Bank : %d z", zeny)
	bm.logger.Info("On Hand : %d z", bm.charZeny)
	bm.logger.Info("=============================")

	// Call hook
	bm.hookManager.CallHook("banking_opened", map[string]interface{}{
		"zeny": zeny,
	})

	return nil
}

// HandleBankingDeposit handles the banking_deposit packet (lines 11997-12013)
func (bm *BankingManager) HandleBankingDeposit(args map[string]interface{}) error {
	// Extract packet data
	reason, ok := args["reason"].(uint16)
	if !ok {
		return fmt.Errorf("invalid reason in banking_deposit packet")
	}

	money, ok := args["money"].(uint64)
	if !ok {
		return fmt.Errorf("invalid money in banking_deposit packet")
	}

	balance, ok := args["balance"].(uint32)
	if !ok {
		return fmt.Errorf("invalid balance in banking_deposit packet")
	}

	// Handle based on reason code
	switch reason {
	case 0: // BDA_SUCCESS
		bm.logger.Success("Bank: Deposit Success.")
		bm.charZeny = balance // Update character zeny

		// Call success hook
		bm.hookManager.CallHook("banking_deposit_success", map[string]interface{}{
			"money":   money,
			"balance": balance,
		})
	case 1: // BDA_ERROR
		bm.logger.Error("Bank: Deposit Error (Try it again).")

		// Call failed hook
		bm.hookManager.CallHook("banking_deposit_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	case 2: // BDA_NO_MONEY
		bm.logger.Error("Bank: No Money For Deposit.")

		// Call failed hook
		bm.hookManager.CallHook("banking_deposit_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	case 3: // BDA_OVERFLOW
		bm.logger.Error("Bank: Money in the bank overflow.")

		// Call failed hook
		bm.hookManager.CallHook("banking_deposit_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	default:
		bm.logger.Error("Bank: Unknown deposit error (code: %d).", reason)

		// Call failed hook
		bm.hookManager.CallHook("banking_deposit_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	}

	return nil
}

// HandleBankingWithdraw handles the banking_withdraw packet (lines 12021-12035)
func (bm *BankingManager) HandleBankingWithdraw(args map[string]interface{}) error {
	// Extract packet data
	reason, ok := args["reason"].(uint16)
	if !ok {
		return fmt.Errorf("invalid reason in banking_withdraw packet")
	}

	money, ok := args["money"].(uint64)
	if !ok {
		return fmt.Errorf("invalid money in banking_withdraw packet")
	}

	balance, ok := args["balance"].(uint32)
	if !ok {
		return fmt.Errorf("invalid balance in banking_withdraw packet")
	}

	// Handle based on reason code
	switch reason {
	case 0: // BWA_SUCCESS
		bm.logger.Success("Bank: Withdraw Success.")
		bm.charZeny = balance // Update character zeny

		// Call success hook
		bm.hookManager.CallHook("banking_withdraw_success", map[string]interface{}{
			"money":   money,
			"balance": balance,
		})
	case 1: // BWA_NO_MONEY
		bm.logger.Error("Bank: No Money for Withdraw.")

		// Call failed hook
		bm.hookManager.CallHook("banking_withdraw_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	case 2: // BWA_UNKNOWN_ERROR
		bm.logger.Error("Bank: Money in the bank overflow.")

		// Call failed hook
		bm.hookManager.CallHook("banking_withdraw_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	default:
		bm.logger.Error("Bank: Unknown withdraw error (code: %d).", reason)

		// Call failed hook
		bm.hookManager.CallHook("banking_withdraw_failed", map[string]interface{}{
			"reason": reason,
			"money":  money,
		})
	}

	return nil
}

// SetCharZeny sets the character's current zeny
func (bm *BankingManager) SetCharZeny(zeny uint32) {
	bm.charZeny = zeny
}
