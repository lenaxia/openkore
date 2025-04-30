package marriage

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MarriageManager handles marriage-related packet handling
type MarriageManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
	charName    string
}

// NewMarriageManager creates a new marriage manager
func NewMarriageManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MarriageManager {
	return &MarriageManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		charName:    "",
	}
}

// RegisterHandlers registers all marriage-related packet handlers
func (mm *MarriageManager) RegisterHandlers() {
	// Register married handler
	mm.parser.RegisterHandlerFunc("01E6", "married", "L",
		[]string{"ID"}, mm.HandleMarried)

	// Register divorced handler
	mm.parser.RegisterHandlerFunc("0205", "divorced", "Z24",
		[]string{"name"}, mm.HandleDivorced)

	// Register marriage partner name handler
	mm.parser.RegisterHandlerFunc("01E4", "marriage_partner_name", "Z24",
		[]string{"name"}, mm.HandleMarriagePartnerName)

	// Register adopt request handler
	mm.parser.RegisterHandlerFunc("01F6", "adopt_request", "Z24",
		[]string{"name"}, mm.HandleAdoptRequest)

	// Register adopt reply handler
	mm.parser.RegisterHandlerFunc("0216", "adopt_reply", "C",
		[]string{"type"}, mm.HandleAdoptReply)
}

// SetCharName sets the character's name
func (mm *MarriageManager) SetCharName(name string) {
	mm.charName = name
}

// HandleMarried handles the married packet (lines 7004-7009)
func (mm *MarriageManager) HandleMarried(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in married packet")
	}

	// Log marriage
	mm.logger.Info("Actor ID %d got married!", id)

	// Call hook
	mm.hookManager.CallHook("married", map[string]interface{}{
		"ID": id,
	})

	return nil
}

// HandleDivorced handles the divorced packet (lines 10442-10445)
func (mm *MarriageManager) HandleDivorced(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in divorced packet")
	}

	// Log divorce
	mm.logger.Info("%s and %s have divorced from each other.", mm.charName, name)

	// Call hook
	mm.hookManager.CallHook("divorced", map[string]interface{}{
		"name": name,
	})

	return nil
}

// HandleMarriagePartnerName handles the marriage_partner_name packet (lines 4112-4116)
func (mm *MarriageManager) HandleMarriagePartnerName(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in marriage_partner_name packet")
	}

	// Log partner name
	mm.logger.Info("Marriage partner name: %s", name)

	// Call hook
	mm.hookManager.CallHook("marriage_partner_name", map[string]interface{}{
		"name": name,
	})

	return nil
}

// HandleAdoptRequest handles the adopt_request packet (lines 10117-10120)
func (mm *MarriageManager) HandleAdoptRequest(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in adopt_request packet")
	}

	// Log adopt request
	mm.logger.Info("%s wishes to adopt you. Do you accept?", name)

	// Call hook
	mm.hookManager.CallHook("adopt_request", map[string]interface{}{
		"name": name,
	})

	return nil
}

// HandleAdoptReply handles the adopt_reply packet (lines 9557-9566)
func (mm *MarriageManager) HandleAdoptReply(args map[string]interface{}) error {
	// Extract packet data
	replyType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in adopt_reply packet")
	}

	// Handle based on reply type
	switch replyType {
	case 0:
		mm.logger.Info("You cannot adopt more than 1 child.")
	case 1:
		mm.logger.Info("You must be at least character level 70 in order to adopt someone.")
	case 2:
		mm.logger.Info("You cannot adopt a married person.")
	default:
		mm.logger.Warning("Unknown adopt reply type: %d", replyType)
	}

	// Call hook
	mm.hookManager.CallHook("adopt_reply", map[string]interface{}{
		"type": replyType,
	})

	return nil
}
