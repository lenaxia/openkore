package effects

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Effect Types
const (
	LEVELUP_EFFECT = 0x0
	// Add more effect constants as needed
)

// Sound Effect Actions
const (
	SOUND_PLAY_ONCE   = 0
	SOUND_PLAY_REPEAT = 1
	SOUND_STOP        = 2
)

// HatEffect represents a hat effect
type HatEffect struct {
	HatEFID uint16
	Handle  string
	Name    string
}

// EffectsManager handles effects-related packet handling
type EffectsManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// Maps for effect names and hat effect data
	effectNames     map[uint16]string
	hatEffectHandle map[uint16]string
	hatEffectName   map[string]string
	emotionsLut     map[uint8]map[string]string
}

// NewEffectsManager creates a new effects manager
func NewEffectsManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *EffectsManager {
	return &EffectsManager{
		parser:          parser,
		hookManager:     hookManager,
		logger:          logger,
		effectNames:     make(map[uint16]string),
		hatEffectHandle: make(map[uint16]string),
		hatEffectName:   make(map[string]string),
		emotionsLut:     make(map[uint8]map[string]string),
	}
}

// RegisterHandlers registers all effects-related packet handlers
func (em *EffectsManager) RegisterHandlers() {
	// Register misc effect handler
	em.parser.RegisterHandlerFunc("01F3", "misc_effect", "L L",
		[]string{"ID", "effect"}, em.HandleMiscEffect)

	// Register sound effect handler
	em.parser.RegisterHandlerFunc("01D3", "sound_effect", "Z24 B L L",
		[]string{"name", "type", "term", "ID"}, em.HandleSoundEffect)

	// Register hat effect handler
	em.parser.RegisterHandlerFunc("0A3B", "hat_effect", "L B a*",
		[]string{"ID", "flag", "effect"}, em.HandleHatEffect)

	// Register emoticon handler
	em.parser.RegisterHandlerFunc("00C0", "emoticon", "L B",
		[]string{"ID", "type"}, em.HandleEmoticon)
}

// SetEffectName sets the name for an effect ID
func (em *EffectsManager) SetEffectName(effectID uint16, name string) {
	em.effectNames[effectID] = name
}

// SetHatEffectHandle sets the handle for a hat effect ID
func (em *EffectsManager) SetHatEffectHandle(hatEFID uint16, handle string) {
	em.hatEffectHandle[hatEFID] = handle
}

// SetHatEffectName sets the name for a hat effect handle
func (em *EffectsManager) SetHatEffectName(handle string, name string) {
	em.hatEffectName[handle] = name
}

// SetEmotionLut sets the emotion lookup table entry
func (em *EffectsManager) SetEmotionLut(emotionType uint8, display string) {
	if em.emotionsLut[emotionType] == nil {
		em.emotionsLut[emotionType] = make(map[string]string)
	}
	em.emotionsLut[emotionType]["display"] = display
}

// GetActorName gets the name of an actor by its ID
func (em *EffectsManager) GetActorName(actorID uint32) string {
	// In a real implementation, this would look up the actor name from the actor list
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Actor#%d", actorID)
}

// GetActorType gets the type of an actor by its ID
func (em *EffectsManager) GetActorType(actorID uint32) string {
	// In a real implementation, this would look up the actor type from the actor list
	// For now, we'll just return a placeholder
	return "Unknown"
}

// HandleMiscEffect handles the misc_effect packet (lines 6861-6869)
func (em *EffectsManager) HandleMiscEffect(args map[string]interface{}) error {
	// Extract packet data
	actorID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in misc_effect packet")
	}

	effectID, ok := args["effect"].(uint32)
	if !ok {
		return fmt.Errorf("invalid effect in misc_effect packet")
	}

	// Get actor name
	actorName := em.GetActorName(actorID)

	// Get effect name
	effectName := "Unknown"
	if name, exists := em.effectNames[uint16(effectID)]; exists {
		effectName = name
	} else {
		effectName = fmt.Sprintf("Unknown #%d", effectID)
	}

	// Log message
	em.logger.Info("%s uses effect: %s", actorName, effectName)

	// Call hook
	em.hookManager.CallHook("misc_effect", map[string]interface{}{
		"ID":     actorID,
		"effect": effectID,
		"name":   effectName,
	})

	return nil
}

// HandleSoundEffect handles the sound_effect packet (lines 6882-6896)
func (em *EffectsManager) HandleSoundEffect(args map[string]interface{}) error {
	// Extract packet data
	soundName, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in sound_effect packet")
	}

	soundType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in sound_effect packet")
	}

	term, ok := args["term"].(uint32)
	if !ok {
		return fmt.Errorf("invalid term in sound_effect packet")
	}

	// Actor ID is optional
	var actorID uint32
	var actorName string
	if id, ok := args["ID"].(uint32); ok {
		actorID = id
		actorName = em.GetActorName(actorID)
	}

	// Log message based on sound type and actor presence
	if actorID != 0 {
		switch soundType {
		case SOUND_PLAY_ONCE:
			em.logger.Info("%s plays: %s", actorName, soundName)
		case SOUND_PLAY_REPEAT:
			em.logger.Info("%s is now playing: %s", actorName, soundName)
		case SOUND_STOP:
			em.logger.Info("%s stopped playing: %s", actorName, soundName)
		default:
			em.logger.Info("%s plays unknown sound type %d: %s", actorName, soundType, soundName)
		}
	} else {
		em.logger.Info("Now playing: %s", soundName)
	}

	// Call hook
	em.hookManager.CallHook("sound_effect", map[string]interface{}{
		"name": soundName,
		"type": soundType,
		"term": term,
		"ID":   actorID,
	})

	return nil
}

// ParseHatEffect parses the hat effect data
func (em *EffectsManager) ParseHatEffect(effect []byte) []HatEffect {
	hatEffects := make([]HatEffect, 0)

	// Parse hat effects from the effect data
	// Each hat effect is 2 bytes (uint16)
	for i := 0; i < len(effect); i += 2 {
		if i+2 > len(effect) {
			break
		}

		hatEFID := uint16(effect[i]) | uint16(effect[i+1])<<8

		hatEffect := HatEffect{
			HatEFID: hatEFID,
		}

		// Set handle and name if available
		if handle, exists := em.hatEffectHandle[hatEFID]; exists {
			hatEffect.Handle = handle
			if name, exists := em.hatEffectName[handle]; exists {
				hatEffect.Name = name
			} else {
				hatEffect.Name = handle
			}
		} else {
			hatEffect.Name = fmt.Sprintf("Unknown #%d", hatEFID)
		}

		hatEffects = append(hatEffects, hatEffect)
	}

	return hatEffects
}

// HandleHatEffect handles the hat_effect packet (lines 7375-7406)
func (em *EffectsManager) HandleHatEffect(args map[string]interface{}) error {
	// Extract packet data
	actorID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in hat_effect packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in hat_effect packet")
	}

	effect, ok := args["effect"].([]byte)
	if !ok {
		return fmt.Errorf("invalid effect in hat_effect packet")
	}

	// Get actor name
	actorName := em.GetActorName(actorID)

	// Parse hat effects
	hatEffects := em.ParseHatEffect(effect)

	// Build hat effect name string
	hatName := ""
	for i, hatEffect := range hatEffects {
		if i > 0 {
			hatName += ", "
		}
		hatName += hatEffect.Name
	}

	// Log message based on flag
	if flag == 1 {
		em.logger.Info("%s uses effect: %s", actorName, hatName)
	} else {
		em.logger.Info("%s is no longer: %s", actorName, hatName)
	}

	// Call hook
	em.hookManager.CallHook("hat_effect", map[string]interface{}{
		"ID":      actorID,
		"flag":    flag,
		"effects": hatEffects,
	})

	return nil
}

// HandleEmoticon handles the emoticon packet (lines 5961-6024)
func (em *EffectsManager) HandleEmoticon(args map[string]interface{}) error {
	// Extract packet data
	actorID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in emoticon packet")
	}

	emotionType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in emoticon packet")
	}

	// Get emotion display
	emotionDisplay := fmt.Sprintf("<emotion #%d>", emotionType)
	if em.emotionsLut[emotionType] != nil {
		if display, exists := em.emotionsLut[emotionType]["display"]; exists {
			emotionDisplay = display
		}
	}

	// Get actor name and type
	actorName := em.GetActorName(actorID)
	actorType := em.GetActorType(actorID)

	// Log message
	em.logger.Info("%s %s: %s", actorType, actorName, emotionDisplay)

	// Call hook
	em.hookManager.CallHook("packet_emotion", map[string]interface{}{
		"emotion": emotionDisplay,
		"ID":      actorID,
	})

	return nil
}
