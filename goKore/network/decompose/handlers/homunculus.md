**Homunculus System Handlers:**
- homunculus_info (lines 6417-6430)
  - Manages homunculus information
  - Processes:
    - Homunculus ID and type
    - Level and stats
    - Hunger and intimacy
    - Skill data
  - Maintains:
    - %homunculusData hash
    - @homunculusSkills array

- homunculus_status (lines 6432-6440)
  - Handles homunculus status updates

- homunculus_skill_update (lines 6442-6455)
  - Processes homunculus skill changes
  - Handles:
    - New skill acquisition
    - Skill level changes
    - Cooldown updates
  - Features:
    - Updates skill availability
    - Triggers 'homunculus_skill_change' hook
    - Maintains skill timers
  - Features:
    - Processes health changes
    - Updates buff/debuff effects
    - Triggers 'homunculus_status_change' hook
    - Manages auto-feed system