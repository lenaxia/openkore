package auction

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// AuctionManager handles auction-related packet handling
type AuctionManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewAuctionManager creates a new auction manager
func NewAuctionManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *AuctionManager {
	return &AuctionManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all auction-related packet handlers
func (am *AuctionManager) RegisterHandlers() {
	// Register auction my sell stop handler
	am.parser.RegisterHandlerFunc("0B33", "auction_my_sell_stop", "B",
		[]string{"flag"}, am.HandleAuctionMySellStop)

	// Register auction windows handler
	am.parser.RegisterHandlerFunc("025D", "auction_windows", "B",
		[]string{"flag"}, am.HandleAuctionWindows)

	// Register auction add item handler
	am.parser.RegisterHandlerFunc("0256", "auction_add_item", "B L",
		[]string{"fail", "ID"}, am.HandleAuctionAddItem)

	// Register auction result handler
	am.parser.RegisterHandlerFunc("0250", "auction_result", "B",
		[]string{"flag"}, am.HandleAuctionResult)
}

// HandleAuctionMySellStop handles the auction_my_sell_stop packet (lines 10239-10252)
func (am *AuctionManager) HandleAuctionMySellStop(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in auction_my_sell_stop packet")
	}

	// Handle based on flag
	switch flag {
	case 0:
		am.logger.Info("You have ended the auction.")
	case 1:
		am.logger.Info("You cannot end the auction.")
	case 2:
		am.logger.Info("Bid number is incorrect.")
	default:
		am.logger.Warning("Unknown results in auction_my_sell_stop (flag: %d)", flag)
	}

	// Call hook
	am.hookManager.CallHook("auction_my_sell_stop", map[string]interface{}{
		"flag": flag,
	})

	return nil
}

// HandleAuctionWindows handles the auction_windows packet (lines 10254-10262)
func (am *AuctionManager) HandleAuctionWindows(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in auction_windows packet")
	}

	// Handle based on flag
	if flag != 0 {
		am.logger.Info("Auction window is now closed.")
	} else {
		am.logger.Info("Auction window is now opened.")
	}

	// Call hook
	am.hookManager.CallHook("auction_windows", map[string]interface{}{
		"flag": flag,
	})

	return nil
}

// HandleAuctionAddItem handles the auction_add_item packet (lines 10264-10272)
func (am *AuctionManager) HandleAuctionAddItem(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in auction_add_item packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in auction_add_item packet")
	}

	// Handle based on fail
	if fail != 0 {
		am.logger.Info("Failed (note: usable items can't be auctioned) to add item with index: %d.", id)
	} else {
		am.logger.Info("Succeeded to add item with index: %d.", id)
	}

	// Call hook
	am.hookManager.CallHook("auction_add_item", map[string]interface{}{
		"fail": fail != 0,
		"ID":   id,
	})

	return nil
}

// HandleAuctionResult handles the auction_result packet (lines 10325-10352)
func (am *AuctionManager) HandleAuctionResult(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in auction_result packet")
	}

	// Handle based on flag
	switch flag {
	case 0:
		am.logger.Info("You have failed to bid into the auction.")
	case 1:
		am.logger.Info("You have successfully bid in the auction.")
	case 2:
		am.logger.Info("The auction has been canceled.")
	case 3:
		am.logger.Info("An auction with at least one bidder cannot be canceled.")
	case 4:
		am.logger.Info("You cannot register more than 5 items in an auction at a time.")
	case 5:
		am.logger.Info("You do not have enough Zeny to pay the Auction Fee.")
	case 6:
		am.logger.Info("You have won the auction.")
	case 7:
		am.logger.Info("You have failed to win the auction.")
	case 8:
		am.logger.Info("You do not have enough Zeny.")
	case 9:
		am.logger.Info("You cannot place more than 5 bids at a time.")
	default:
		am.logger.Warning("Unknown results in auction_result (flag: %d)", flag)
	}

	// Call hook
	am.hookManager.CallHook("auction_result", map[string]interface{}{
		"flag": flag,
	})

	return nil
}
