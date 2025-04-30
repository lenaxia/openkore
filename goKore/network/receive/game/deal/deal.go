package deal

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// DealItem represents an item in a deal
type DealItem struct {
	ID         uint32
	NameID     uint16
	Amount     uint16
	Identified bool
	Broken     bool
	Upgrade    uint8
	Cards      []uint32
	Options    []uint32
	Name       string
}

// DealState represents the current state of a deal
type DealState struct {
	ID             uint32
	Name           string
	OtherItems     map[uint16]*DealItem
	OtherZeny      uint32
	OtherFinalize  bool
	YourItems      map[uint16]*DealItem
	YourZeny       uint32
	YourFinalize   bool
	LastItemAmount uint16
}

// IncomingDeal represents an incoming deal request
type IncomingDeal struct {
	ID    uint32
	Name  string
	Level uint16
}

// OutgoingDeal represents an outgoing deal request
type OutgoingDeal struct {
	ID uint32
}

// DealManager handles deal-related packet handling
type DealManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for deal interactions
	currentDeal  *DealState
	incomingDeal *IncomingDeal
	outgoingDeal *OutgoingDeal
}

// NewDealManager creates a new deal manager
func NewDealManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *DealManager {
	return &DealManager{
		parser:       parser,
		hookManager:  hookManager,
		logger:       logger,
		currentDeal:  nil,
		incomingDeal: nil,
		outgoingDeal: nil,
	}
}

// RegisterHandlers registers all deal-related packet handlers
func (dm *DealManager) RegisterHandlers() {
	// Register deal add other handler
	dm.parser.RegisterHandlerFunc("00E9", "deal_add_other", "v V V C C C a8 a25",
		[]string{"amount", "nameID", "identified", "broken", "upgrade", "cards", "options"}, dm.HandleDealAddOther)

	// Register deal add you handler
	dm.parser.RegisterHandlerFunc("00E0", "deal_add_you", "C V",
		[]string{"fail", "ID"}, dm.HandleDealAddYou)

	// Register deal begin handler
	dm.parser.RegisterHandlerFunc("00E1", "deal_begin", "C",
		[]string{"type"}, dm.HandleDealBegin)

	// Register deal cancelled handler
	dm.parser.RegisterHandlerFunc("00E2", "deal_cancelled", "",
		[]string{}, dm.HandleDealCancelled)

	// Register deal complete handler
	dm.parser.RegisterHandlerFunc("00EF", "deal_complete", "",
		[]string{}, dm.HandleDealComplete)

	// Register deal finalize handler
	dm.parser.RegisterHandlerFunc("00EB", "deal_finalize", "C",
		[]string{"type"}, dm.HandleDealFinalize)

	// Register deal request handler
	dm.parser.RegisterHandlerFunc("00E5", "deal_request", "Z24 V W",
		[]string{"user", "ID", "level"}, dm.HandleDealRequest)
}

// GetItemName gets the name of an item by its nameID
func (dm *DealManager) GetItemName(nameID uint16) string {
	// In a real implementation, this would look up the item name from a database
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Item#%d", nameID)
}

// GetActorName gets the name of an actor by its ID
func (dm *DealManager) GetActorName(actorID uint32) string {
	// In a real implementation, this would look up the actor name from the actor list
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Actor#%d", actorID)
}

// FormatNumber formats a number with commas
func (dm *DealManager) FormatNumber(num uint32) string {
	// In a real implementation, this would format the number with commas
	// For now, we'll just return the number as a string
	return fmt.Sprintf("%d", num)
}

// HandleDealAddOther handles the deal_add_other packet (lines 5825-5844)
func (dm *DealManager) HandleDealAddOther(args map[string]interface{}) error {
	// Extract packet data
	nameID, ok := args["nameID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid nameID in deal_add_other packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in deal_add_other packet")
	}

	identified, ok := args["identified"].(uint32)
	if !ok {
		return fmt.Errorf("invalid identified in deal_add_other packet")
	}

	broken, ok := args["broken"].(uint8)
	if !ok {
		return fmt.Errorf("invalid broken in deal_add_other packet")
	}

	upgrade, ok := args["upgrade"].(uint8)
	if !ok {
		return fmt.Errorf("invalid upgrade in deal_add_other packet")
	}

	cards, ok := args["cards"].([]byte)
	if !ok {
		return fmt.Errorf("invalid cards in deal_add_other packet")
	}

	options, ok := args["options"].([]byte)
	if !ok {
		return fmt.Errorf("invalid options in deal_add_other packet")
	}

	// Check if we're in a deal
	if dm.currentDeal == nil {
		return fmt.Errorf("received deal_add_other packet but not in a deal")
	}

	// Handle adding item or zeny
	if nameID > 0 {
		// Add item to deal
		item, exists := dm.currentDeal.OtherItems[uint16(nameID)]
		if !exists {
			item = &DealItem{
				NameID:     uint16(nameID),
				Amount:     0,
				Identified: identified != 0,
				Broken:     broken != 0,
				Upgrade:    upgrade,
				Cards:      make([]uint32, 0),
				Options:    make([]uint32, 0),
				Name:       dm.GetItemName(uint16(nameID)),
			}
			dm.currentDeal.OtherItems[uint16(nameID)] = item
		}
		item.Amount += amount

		// Parse cards
		for i := 0; i < len(cards); i += 4 {
			if i+4 > len(cards) {
				break
			}
			cardID := uint32(cards[i]) | uint32(cards[i+1])<<8 | uint32(cards[i+2])<<16 | uint32(cards[i+3])<<24
			if cardID > 0 {
				item.Cards = append(item.Cards, cardID)
			}
		}

		// Parse options
		for i := 0; i < len(options); i += 5 {
			if i+5 > len(options) {
				break
			}
			optionID := uint32(options[i]) | uint32(options[i+1])<<8 | uint32(options[i+2])<<16 | uint32(options[i+3])<<24
			optionValue := uint32(options[i+4])
			if optionID > 0 {
				item.Options = append(item.Options, optionID)
				item.Options = append(item.Options, optionValue)
			}
		}

		// Log message
		dm.logger.Info("%s added Item to Deal: %s x %d", dm.currentDeal.Name, item.Name, amount)
	} else if amount > 0 {
		// Add zeny to deal
		dm.currentDeal.OtherZeny += uint32(amount)
		formattedAmount := dm.FormatNumber(uint32(amount))
		dm.logger.Info("%s added %s z to Deal", dm.currentDeal.Name, formattedAmount)
	}

	// Call hook
	dm.hookManager.CallHook("deal_add_other", map[string]interface{}{
		"nameID":     nameID,
		"amount":     amount,
		"identified": identified != 0,
		"broken":     broken != 0,
		"upgrade":    upgrade,
		"cards":      cards,
		"options":    options,
	})

	return nil
}

// HandleDealAddYou handles the deal_add_you packet (lines 7828-7862)
func (dm *DealManager) HandleDealAddYou(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in deal_add_you packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in deal_add_you packet")
	}

	// Handle failure cases
	if fail == 1 {
		dm.logger.Error("That person is overweight; you cannot trade.")
		return nil
	} else if fail == 2 {
		dm.logger.Error("This item cannot be traded.")
		return nil
	} else if fail == 192 {
		dm.logger.Debug("Unknown status (success).")
	} else if fail != 0 {
		dm.logger.Error("You cannot trade (fail code %d).", fail)
		return nil
	}

	// Check if we're in a deal
	if dm.currentDeal == nil {
		return fmt.Errorf("received deal_add_you packet but not in a deal")
	}

	// Check if ID is valid
	if id == 0 {
		return nil
	}

	// In a real implementation, we would get the item from inventory
	// For now, we'll just create a placeholder
	itemNameID := uint16(id)
	itemName := dm.GetItemName(itemNameID)

	// Add item to deal
	item, exists := dm.currentDeal.YourItems[itemNameID]
	if !exists {
		item = &DealItem{
			NameID: itemNameID,
			Amount: 0,
			Name:   itemName,
		}
		dm.currentDeal.YourItems[itemNameID] = item
	}
	item.Amount += dm.currentDeal.LastItemAmount

	// Log message
	dm.logger.Info("You added Item to Deal: %s x %d", itemName, dm.currentDeal.LastItemAmount)

	// Call hook
	dm.hookManager.CallHook("deal_you_added", map[string]interface{}{
		"id":     id,
		"nameID": itemNameID,
		"name":   itemName,
		"amount": dm.currentDeal.LastItemAmount,
	})

	return nil
}

// HandleDealBegin handles the deal_begin packet (lines 5846-5880)
func (dm *DealManager) HandleDealBegin(args map[string]interface{}) error {
	// Extract packet data
	dealType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in deal_begin packet")
	}

	// Handle based on deal type
	switch dealType {
	case 0:
		dm.logger.Error("That person is too far from you to trade.")
		dm.hookManager.CallHook("error_deal", map[string]interface{}{"type": dealType})
	case 2:
		dm.logger.Error("That person is in another deal.")
		dm.hookManager.CallHook("error_deal", map[string]interface{}{"type": dealType})
	case 3:
		// Initialize current deal
		dm.currentDeal = &DealState{
			OtherItems:    make(map[uint16]*DealItem),
			OtherZeny:     0,
			OtherFinalize: false,
			YourItems:     make(map[uint16]*DealItem),
			YourZeny:      0,
			YourFinalize:  false,
		}

		// Set deal name based on incoming or outgoing deal
		if dm.incomingDeal != nil {
			dm.currentDeal.Name = dm.incomingDeal.Name
			dm.incomingDeal = nil
		} else if dm.outgoingDeal != nil {
			dm.currentDeal.ID = dm.outgoingDeal.ID
			dm.currentDeal.Name = dm.GetActorName(dm.outgoingDeal.ID)
			dm.outgoingDeal = nil
		}

		dm.logger.Info("Engaged Deal with %s", dm.currentDeal.Name)
		dm.hookManager.CallHook("engaged_deal", map[string]interface{}{"name": dm.currentDeal.Name})
	case 5:
		dm.logger.Error("That person is opening storage.")
		dm.hookManager.CallHook("error_deal", map[string]interface{}{"type": dealType})
	default:
		dm.logger.Error("Deal request failed (unknown error %d).", dealType)
		dm.hookManager.CallHook("error_deal", map[string]interface{}{"type": dealType})
	}

	return nil
}

// HandleDealCancelled handles the deal_cancelled packet (lines 5882-5888)
func (dm *DealManager) HandleDealCancelled(args map[string]interface{}) error {
	// Reset deal state
	dm.incomingDeal = nil
	dm.outgoingDeal = nil
	dm.currentDeal = nil

	// Log message
	dm.logger.Info("Deal Cancelled")

	// Call hook
	dm.hookManager.CallHook("cancelled_deal", map[string]interface{}{})

	return nil
}

// HandleDealComplete handles the deal_complete packet (lines 5890-5896)
func (dm *DealManager) HandleDealComplete(args map[string]interface{}) error {
	// Reset deal state
	dm.outgoingDeal = nil
	dm.incomingDeal = nil
	dm.currentDeal = nil

	// Log message
	dm.logger.Info("Deal Complete")

	// Call hook
	dm.hookManager.CallHook("complete_deal", map[string]interface{}{})

	return nil
}

// HandleDealFinalize handles the deal_finalize packet (lines 5898-5911)
func (dm *DealManager) HandleDealFinalize(args map[string]interface{}) error {
	// Extract packet data
	finalizeType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in deal_finalize packet")
	}

	// Check if we're in a deal
	if dm.currentDeal == nil {
		return fmt.Errorf("received deal_finalize packet but not in a deal")
	}

	// Handle based on finalize type
	if finalizeType == 1 {
		dm.currentDeal.OtherFinalize = true
		dm.logger.Info("%s finalized the Deal", dm.currentDeal.Name)
		dm.hookManager.CallHook("finalized_deal", map[string]interface{}{"name": dm.currentDeal.Name})
	} else {
		dm.currentDeal.YourFinalize = true
		dm.logger.Info("You finalized the Deal")
	}

	return nil
}

// HandleDealRequest handles the deal_request packet (lines 5913-5927)
func (dm *DealManager) HandleDealRequest(args map[string]interface{}) error {
	// Extract packet data
	user, ok := args["user"].(string)
	if !ok {
		return fmt.Errorf("invalid user in deal_request packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in deal_request packet")
	}

	level, ok := args["level"].(uint16)
	if !ok {
		return fmt.Errorf("invalid level in deal_request packet")
	}

	// Create incoming deal
	dm.incomingDeal = &IncomingDeal{
		ID:    id,
		Name:  user,
		Level: level,
	}

	// Log message
	dm.logger.Info("%s (level %d) Requests a Deal", user, level)
	dm.logger.Info("Type 'deal' to start dealing, or 'deal no' to deny the deal.")

	// Call hook
	dm.hookManager.CallHook("incoming_deal", map[string]interface{}{
		"name":  user,
		"level": level,
		"ID":    id,
	})

	return nil
}

// SetLastItemAmount sets the last item amount for the current deal
func (dm *DealManager) SetLastItemAmount(amount uint16) {
	if dm.currentDeal != nil {
		dm.currentDeal.LastItemAmount = amount
	}
}

// SetOutgoingDealID sets the outgoing deal ID
func (dm *DealManager) SetOutgoingDealID(id uint32) {
	dm.outgoingDeal = &OutgoingDeal{
		ID: id,
	}
}
