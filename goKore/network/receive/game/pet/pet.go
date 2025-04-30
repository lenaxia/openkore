package pet

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// PetInfo represents pet information
type PetInfo struct {
	ID         uint32
	Name       string
	RenameFlag bool
	Level      uint16
	Hungry     uint16
	Friendly   uint16
	Accessory  uint16
	Type       uint16
}

// PetManager handles pet-related packet handling
type PetManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for pet interactions
	petInfo *PetInfo
}

// NewPetManager creates a new pet manager
func NewPetManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *PetManager {
	return &PetManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		petInfo:     nil,
	}
}

// RegisterHandlers registers all pet-related packet handlers
func (pm *PetManager) RegisterHandlers() {
	// Register pet capture process handler
	pm.parser.RegisterHandlerFunc("019E", "pet_capture_process", "",
		[]string{}, pm.HandlePetCaptureProcess)

	// Register pet capture result handler
	pm.parser.RegisterHandlerFunc("01A0", "pet_capture_result", "B",
		[]string{"success"}, pm.HandlePetCaptureResult)

	// Register pet emotion handler
	pm.parser.RegisterHandlerFunc("01AA", "pet_emotion", "L B",
		[]string{"ID", "type"}, pm.HandlePetEmotion)

	// Register pet evolution result handler
	pm.parser.RegisterHandlerFunc("09FC", "pet_evolution_result", "B",
		[]string{"result"}, pm.HandlePetEvolutionResult)

	// Register pet food handler
	pm.parser.RegisterHandlerFunc("01A3", "pet_food", "B W",
		[]string{"success", "foodID"}, pm.HandlePetFood)

	// Register pet info handler
	pm.parser.RegisterHandlerFunc("01A2", "pet_info", "Z24 B W W W W W",
		[]string{"name", "renameflag", "level", "hungry", "friendly", "accessory", "type"}, pm.HandlePetInfo)

	// Register pet info2 handler
	pm.parser.RegisterHandlerFunc("01A4", "pet_info2", "B L L",
		[]string{"type", "ID", "value"}, pm.HandlePetInfo2)
}

// HandlePetCaptureProcess handles the pet_capture_process packet (lines 8985-8988)
func (pm *PetManager) HandlePetCaptureProcess(args map[string]interface{}) error {
	// Log message
	pm.logger.Info("Attempting to capture pet (slot machine).")

	// Call hook
	pm.hookManager.CallHook("pet_capture_process", map[string]interface{}{})

	return nil
}

// HandlePetCaptureResult handles the pet_capture_result packet (lines 8990-8997)
func (pm *PetManager) HandlePetCaptureResult(args map[string]interface{}) error {
	// Extract packet data
	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in pet_capture_result packet")
	}

	// Handle based on success
	if success != 0 {
		pm.logger.Info("Pet capture success")
	} else {
		pm.logger.Info("Pet capture failed")
	}

	// Call hook
	pm.hookManager.CallHook("pet_capture_result", map[string]interface{}{
		"success": success != 0,
	})

	return nil
}

// HandlePetEmotion handles the pet_emotion packet (lines 8999-9006)
func (pm *PetManager) HandlePetEmotion(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in pet_emotion packet")
	}

	emotionType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in pet_emotion packet")
	}

	// Get emotion display
	emotion := fmt.Sprintf("/e%d", emotionType) // Default if not found in emotions_lut

	// Log emotion if pet exists
	if pm.petInfo != nil && pm.petInfo.ID == id {
		pm.logger.Info("%s : %s", pm.petInfo.Name, emotion)
	}

	// Call hook
	pm.hookManager.CallHook("pet_emotion", map[string]interface{}{
		"ID":   id,
		"type": emotionType,
	})

	return nil
}

// HandlePetEvolutionResult handles the pet_evolution_result packet (lines 9008-9023)
func (pm *PetManager) HandlePetEvolutionResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in pet_evolution_result packet")
	}

	// Handle based on result
	switch result {
	case 0x0:
		pm.logger.Error("Pet evolution error.")
	case 0x1:
		pm.logger.Error("Pet evolution error: No pet summoned.")
	case 0x2:
		pm.logger.Error("Pet evolution error: No pet egg.")
	case 0x3:
		pm.logger.Error("Unequip pet accessories first to start evolution.")
	case 0x4:
		pm.logger.Error("Insufficient materials for evolution.")
	case 0x5:
		pm.logger.Error("Loyal Intimacy is required to evolve.")
	case 0x6:
		pm.logger.Success("Pet evolution success.")
	default:
		pm.logger.Error("Unknown pet evolution result: %d", result)
	}

	// Call hook
	pm.hookManager.CallHook("pet_evolution_result", map[string]interface{}{
		"result": result,
	})

	return nil
}

// HandlePetFood handles the pet_food packet (lines 9025-9032)
func (pm *PetManager) HandlePetFood(args map[string]interface{}) error {
	// Extract packet data
	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in pet_food packet")
	}

	foodID, ok := args["foodID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid foodID in pet_food packet")
	}

	// Handle based on success
	if success != 0 {
		pm.logger.Info("Fed pet with item ID %d", foodID)
	} else {
		pm.logger.Error("Failed to feed pet with item ID %d: no food in inventory.", foodID)
	}

	// Call hook
	pm.hookManager.CallHook("pet_food", map[string]interface{}{
		"success": success != 0,
		"foodID":  foodID,
	})

	return nil
}

// HandlePetInfo handles the pet_info packet (lines 9034-9044)
func (pm *PetManager) HandlePetInfo(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in pet_info packet")
	}

	// Create pet info if it doesn't exist
	if pm.petInfo == nil {
		pm.petInfo = &PetInfo{}
	}

	// Update pet info
	pm.petInfo.Name = name

	// Extract and update other properties
	if renameflag, ok := args["renameflag"].(uint8); ok {
		pm.petInfo.RenameFlag = renameflag != 0
	}
	if level, ok := args["level"].(uint16); ok {
		pm.petInfo.Level = level
	}
	if hungry, ok := args["hungry"].(uint16); ok {
		pm.petInfo.Hungry = hungry
	}
	if friendly, ok := args["friendly"].(uint16); ok {
		pm.petInfo.Friendly = friendly
	}
	if accessory, ok := args["accessory"].(uint16); ok {
		pm.petInfo.Accessory = accessory
	}
	if petType, ok := args["type"].(uint16); ok {
		pm.petInfo.Type = petType
	}

	// Debug log
	pm.logger.Debug("Pet status: name=%s name_set=%s level=%d hungry=%d intimacy=%d accessory=%d type=%d",
		pm.petInfo.Name,
		pm.petInfo.RenameFlag,
		pm.petInfo.Level,
		pm.petInfo.Hungry,
		pm.petInfo.Friendly,
		pm.petInfo.Accessory,
		pm.petInfo.Type)

	// Call hook
	pm.hookManager.CallHook("pet_info", map[string]interface{}{
		"petInfo": pm.petInfo,
	})

	return nil
}

// HandlePetInfo2 handles the pet_info2 packet (lines 9046-9098)
func (pm *PetManager) HandlePetInfo2(args map[string]interface{}) error {
	// Extract packet data
	infoType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in pet_info2 packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in pet_info2 packet")
	}

	value, ok := args["value"].(uint32)
	if !ok {
		return fmt.Errorf("invalid value in pet_info2 packet")
	}

	// Create pet info if it doesn't exist
	if pm.petInfo == nil {
		pm.petInfo = &PetInfo{}
	}

	// Handle based on type
	switch infoType {
	case 0:
		// You own no pet
		pm.petInfo.ID = 0
		pm.logger.Debug("You own no pet.")
	case 1:
		// Pet friendly
		pm.petInfo.Friendly = uint16(value)
		pm.logger.Debug("Pet friendly: %d", value)
	case 2:
		// Pet hungry
		pm.petInfo.Hungry = uint16(value)
		pm.logger.Debug("Pet hungry: %d", value)
	case 3:
		// Pet accessory
		pm.petInfo.Accessory = uint16(value)
		pm.logger.Debug("Pet accessory info: %d", value)
	case 4:
		// Pet performance
		pm.logger.Debug("Pet performance info: %d", value)
	case 5:
		// You own pet with this ID
		pm.petInfo.ID = id
		pm.logger.Debug("You own pet with ID: %d", id)
	default:
		pm.logger.Debug("Unknown pet info type: %d", infoType)
	}

	// Call hook
	pm.hookManager.CallHook("pet_info2", map[string]interface{}{
		"type":  infoType,
		"ID":    id,
		"value": value,
	})

	return nil
}
