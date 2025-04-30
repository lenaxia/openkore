package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RegisterWithParser registers the skill handlers with the given parser and hook manager
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the skills list manager
	skillsListManager := NewSkillsListManager(parser, hookManager)

	// Create the skill add manager
	skillAddManager := NewSkillAddManager(parser, hookManager)

	// Create the skill update manager
	skillUpdateManager := NewSkillUpdateManager(parser, hookManager)

	// Create the skill delete manager
	skillDeleteManager := NewSkillDeleteManager(parser, hookManager)

	// Create the skill message manager
	skillMsgManager := NewSkillMsgManager(parser, hookManager)

	// Create the skill cast manager
	skillCastManager := NewSkillCastManager(parser, hookManager)

	// Create the cast cancelled manager
	castCancelledManager := NewCastCancelledManager(parser, hookManager)

	// Create the skill use failed manager
	skillUseFailedManager := NewSkillUseFailedManager(parser, hookManager)

	// Create the skill delay manager
	skillDelayManager := NewSkillDelayManager(parser, hookManager)

	// Create the gospel buff manager
	gospelBuffManager := NewGospelBuffManager(parser, hookManager)

	// Create the combo delay manager
	comboDelayManager := NewComboDelayManager(parser, hookManager)

	// Create the attack range manager
	attackRangeManager := NewAttackRangeManager(parser, hookManager)

	// Create the sage autospell manager
	sageAutospellManager := NewSageAutospellManager(parser, hookManager)

	// Create the skill exchange item manager
	skillExchangeItemManager := NewSkillExchangeItemManager(parser, hookManager)

	// Create the devotion manager
	devotionManager := NewDevotionManager(parser, hookManager)

	// Create the blade stop manager
	bladeStopManager := NewBladeStopManager(parser, hookManager)

	// Create the high jump manager
	highJumpManager := NewHighJumpManager(parser, hookManager)

	// Create the resurrection manager
	resurrectionManager := NewResurrectionManager(parser, hookManager)

	// Create the sense result manager
	senseResultManager := NewSenseResultManager(parser, hookManager)

	// Create the area spell manager
	areaSpellManager := NewAreaSpellManager(parser, hookManager)

	// Create the starplace manager
	starplaceManager := NewStarplaceManager(parser, hookManager)

	// Create the talkie box manager
	talkieBoxManager := NewTalkieBoxManager(parser, hookManager)

	// Register handlers
	skillsListManager.RegisterHandlers()
	skillAddManager.RegisterHandlers()
	skillUpdateManager.RegisterHandlers()
	skillDeleteManager.RegisterHandlers()
	skillMsgManager.RegisterHandlers()
	skillCastManager.RegisterHandlers()
	castCancelledManager.RegisterHandlers()
	skillUseFailedManager.RegisterHandlers()
	skillDelayManager.RegisterHandlers()
	gospelBuffManager.RegisterHandlers()
	comboDelayManager.RegisterHandlers()
	attackRangeManager.RegisterHandlers()
	sageAutospellManager.RegisterHandlers()
	skillExchangeItemManager.RegisterHandlers()
	devotionManager.RegisterHandlers()
	bladeStopManager.RegisterHandlers()
	highJumpManager.RegisterHandlers()
	resurrectionManager.RegisterHandlers()
	senseResultManager.RegisterHandlers()
	areaSpellManager.RegisterHandlers()
	starplaceManager.RegisterHandlers()
	talkieBoxManager.RegisterHandlers()
}

// RegisterWithBaseReceive registers the skill handlers with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the skills_list handler
	baseReceive.RegisterHandler("skills_list", func(args map[string]interface{}) error {
		// Create a skills list manager for this specific call
		manager := NewSkillsListManager(nil, nil)
		return manager.handleSkillsList(args)
	})

	// Register the homun_skills_list handler
	baseReceive.RegisterHandler("homun_skills_list", func(args map[string]interface{}) error {
		// Create a skills list manager for this specific call
		manager := NewSkillsListManager(nil, nil)
		return manager.handleSkillsList(args)
	})

	// Register the merc_skills_list handler
	baseReceive.RegisterHandler("merc_skills_list", func(args map[string]interface{}) error {
		// Create a skills list manager for this specific call
		manager := NewSkillsListManager(nil, nil)
		return manager.handleSkillsList(args)
	})

	// Register the skills_list_short handler
	baseReceive.RegisterHandler("skills_list_short", func(args map[string]interface{}) error {
		// Create a skills list manager for this specific call
		manager := NewSkillsListManager(nil, nil)
		return manager.handleSkillsList(args)
	})

	// Register the skill_add handler
	baseReceive.RegisterHandler("skill_add", func(args map[string]interface{}) error {
		// Create a skill add manager for this specific call
		manager := NewSkillAddManager(nil, nil)
		return manager.handleSkillAdd(args)
	})

	// Register the skill_add_new handler
	baseReceive.RegisterHandler("skill_add_new", func(args map[string]interface{}) error {
		// Create a skill add manager for this specific call
		manager := NewSkillAddManager(nil, nil)
		return manager.handleSkillAdd(args)
	})

	// Register the skill_update handler
	baseReceive.RegisterHandler("skill_update", func(args map[string]interface{}) error {
		// Create a skill update manager for this specific call
		manager := NewSkillUpdateManager(nil, nil)
		return manager.handleSkillUpdate(args)
	})

	// Register the skill_delete handler
	baseReceive.RegisterHandler("skill_delete", func(args map[string]interface{}) error {
		// Create a skill delete manager for this specific call
		manager := NewSkillDeleteManager(nil, nil)
		return manager.handleSkillDelete(args)
	})

	// Register the skill_msg handler
	baseReceive.RegisterHandler("skill_msg", func(args map[string]interface{}) error {
		// Create a skill message manager for this specific call
		manager := NewSkillMsgManager(nil, nil)
		return manager.handleSkillMsg(args)
	})

	// Register the skill_cast handler
	baseReceive.RegisterHandler("skill_cast", func(args map[string]interface{}) error {
		// Create a skill cast manager for this specific call
		manager := NewSkillCastManager(nil, nil)
		return manager.handleSkillCast(args)
	})

	// Register the skill_cast_expanded handler
	baseReceive.RegisterHandler("skill_cast_expanded", func(args map[string]interface{}) error {
		// Create a skill cast manager for this specific call
		manager := NewSkillCastManager(nil, nil)
		return manager.handleSkillCast(args)
	})

	// Register the skill_cast_nodamage handler
	baseReceive.RegisterHandler("skill_cast_nodamage", func(args map[string]interface{}) error {
		// Create a skill cast manager for this specific call
		manager := NewSkillCastManager(nil, nil)
		return manager.handleSkillCastNoDamage(args)
	})

	// Register the cast_cancelled handler
	baseReceive.RegisterHandler("cast_cancelled", func(args map[string]interface{}) error {
		// Create a cast cancelled manager for this specific call
		manager := NewCastCancelledManager(nil, nil)
		return manager.handleCastCancelled(args)
	})

	// Register the cast_cancelled_expanded handler
	baseReceive.RegisterHandler("cast_cancelled_expanded", func(args map[string]interface{}) error {
		// Create a cast cancelled manager for this specific call
		manager := NewCastCancelledManager(nil, nil)
		return manager.handleCastCancelled(args)
	})

	// Register the skill_use_failed handler
	baseReceive.RegisterHandler("skill_use_failed", func(args map[string]interface{}) error {
		// Create a skill use failed manager for this specific call
		manager := NewSkillUseFailedManager(nil, nil)
		return manager.handleSkillUseFailed(args)
	})

	// Register the skill_use_failed_expanded handler
	baseReceive.RegisterHandler("skill_use_failed_expanded", func(args map[string]interface{}) error {
		// Create a skill use failed manager for this specific call
		manager := NewSkillUseFailedManager(nil, nil)
		return manager.handleSkillUseFailed(args)
	})

	// Register the skill_post_delay handler
	baseReceive.RegisterHandler("skill_post_delay", func(args map[string]interface{}) error {
		// Create a skill delay manager for this specific call
		manager := NewSkillDelayManager(nil, nil)
		return manager.handleSkillPostDelay(args)
	})

	// Register the skill_post_delaylist handler
	baseReceive.RegisterHandler("skill_post_delaylist", func(args map[string]interface{}) error {
		// Create a skill delay manager for this specific call
		manager := NewSkillDelayManager(nil, nil)
		return manager.handleSkillPostDelayList(args)
	})

	// Register the skill_post_delaylist_expanded handler
	baseReceive.RegisterHandler("skill_post_delaylist_expanded", func(args map[string]interface{}) error {
		// Create a skill delay manager for this specific call
		manager := NewSkillDelayManager(nil, nil)
		return manager.handleSkillPostDelayList(args)
	})

	// Register the gospel_buff_aligned handler
	baseReceive.RegisterHandler("gospel_buff_aligned", func(args map[string]interface{}) error {
		// Create a gospel buff manager for this specific call
		manager := NewGospelBuffManager(nil, nil)
		return manager.handleGospelBuffAligned(args)
	})

	// Register the combo_delay handler
	baseReceive.RegisterHandler("combo_delay", func(args map[string]interface{}) error {
		// Create a combo delay manager for this specific call
		manager := NewComboDelayManager(nil, nil)
		return manager.handleComboDelay(args)
	})

	// Register the attack_range handler
	baseReceive.RegisterHandler("attack_range", func(args map[string]interface{}) error {
		// Create an attack range manager for this specific call
		manager := NewAttackRangeManager(nil, nil)
		return manager.handleAttackRange(args)
	})

	// Register the sage_autospell handler
	baseReceive.RegisterHandler("sage_autospell", func(args map[string]interface{}) error {
		// Create a sage autospell manager for this specific call
		manager := NewSageAutospellManager(nil, nil)
		return manager.handleSageAutospell(args)
	})

	// Register the sage_autospell_shadow handler
	baseReceive.RegisterHandler("sage_autospell_shadow", func(args map[string]interface{}) error {
		// Create a sage autospell manager for this specific call
		manager := NewSageAutospellManager(nil, nil)
		return manager.handleSageAutospell(args)
	})

	// Register the skill_exchange_item handler
	baseReceive.RegisterHandler("skill_exchange_item", func(args map[string]interface{}) error {
		// Create a skill exchange item manager for this specific call
		manager := NewSkillExchangeItemManager(nil, nil)
		return manager.handleSkillExchangeItem(args)
	})

	// Register the devotion handler
	baseReceive.RegisterHandler("devotion", func(args map[string]interface{}) error {
		// Create a devotion manager for this specific call
		manager := NewDevotionManager(nil, nil)
		return manager.handleDevotion(args)
	})

	// Register the blade_stop handler
	baseReceive.RegisterHandler("blade_stop", func(args map[string]interface{}) error {
		// Create a blade stop manager for this specific call
		manager := NewBladeStopManager(nil, nil)
		return manager.handleBladeStop(args)
	})

	// Register the high_jump handler
	baseReceive.RegisterHandler("high_jump", func(args map[string]interface{}) error {
		// Create a high jump manager for this specific call
		manager := NewHighJumpManager(nil, nil)
		return manager.handleHighJump(args)
	})

	// Register the resurrection handler
	baseReceive.RegisterHandler("resurrection", func(args map[string]interface{}) error {
		// Create a resurrection manager for this specific call
		manager := NewResurrectionManager(nil, nil)
		return manager.handleResurrection(args)
	})

	// Register the sense_result handler
	baseReceive.RegisterHandler("sense_result", func(args map[string]interface{}) error {
		// Create a sense result manager for this specific call
		manager := NewSenseResultManager(nil, nil)
		return manager.handleSenseResult(args)
	})

	// Register the area_spell handler
	baseReceive.RegisterHandler("area_spell", func(args map[string]interface{}) error {
		// Create an area spell manager for this specific call
		manager := NewAreaSpellManager(nil, nil)
		return manager.handleAreaSpell(args)
	})

	// Register the area_spell_scribble handler
	baseReceive.RegisterHandler("area_spell_scribble", func(args map[string]interface{}) error {
		// Create an area spell manager for this specific call
		manager := NewAreaSpellManager(nil, nil)
		return manager.handleAreaSpell(args)
	})

	// Register the area_spell_expanded handler
	baseReceive.RegisterHandler("area_spell_expanded", func(args map[string]interface{}) error {
		// Create an area spell manager for this specific call
		manager := NewAreaSpellManager(nil, nil)
		return manager.handleAreaSpell(args)
	})

	// Register the area_spell_disappears handler
	baseReceive.RegisterHandler("area_spell_disappears", func(args map[string]interface{}) error {
		// Create an area spell manager for this specific call
		manager := NewAreaSpellManager(nil, nil)
		return manager.handleAreaSpellDisappears(args)
	})

	// Register the starplace handler
	baseReceive.RegisterHandler("starplace", func(args map[string]interface{}) error {
		// Create a starplace manager for this specific call
		manager := NewStarplaceManager(nil, nil)
		return manager.handleStarplace(args)
	})

	// Register the talkie_box handler
	baseReceive.RegisterHandler("talkie_box", func(args map[string]interface{}) error {
		// Create a talkie box manager for this specific call
		manager := NewTalkieBoxManager(nil, nil)
		return manager.handleTalkieBox(args)
	})
}
