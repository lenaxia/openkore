**Pet Capture & Evolution:**
- pet_capture_process() - Slot machine animation (lines 8985-8988)
  - Shows capture attempt message

- pet_capture_result() - Capture outcome (lines 8990-8997)
  - Handles success/failure cases
  - Displays appropriate message

- pet_evolution_result() - Evolution results (lines 9008-9023)
  - Processes various failure cases:
    * No pet called (1)
    * No egg (2)
    * Accessories equipped (3)
    * Insufficient materials (4)
    * Low intimacy (5)
  - Displays success message (6)

**Pet Interactions:**
- pet_emotion() - Pet emotes (lines 8999-9006)
  - Displays pet emotions
  - Maps emotion IDs to text

- pet_food() - Feeding results (lines 9025-9032)
  - Handles success/failure
  - Displays food item name

**Pet Information:**
- pet_info() - Basic pet stats (lines 9034-9044)
  - Updates:
    * Name and rename flag
    * Level and hunger
    * Friendliness
    * Accessory
    * Type
  - Logs full status

- pet_info2() - Detailed pet info (lines 9046-9098)
  - Handles various info types:
    * No pet (0)
    * Friendliness (1)
    * Hunger (2)
    * Accessory (3)
    * Performance (4)
    * Pet ID (5)

**Elemental Information:**
- elemental_info() - Elemental data (lines 9100-9111)
  - Updates elemental actor reference
  - Processes various stats
  - Maintains elemental object

**Pet System Handlers:**
- pet_info (lines 6377-6390)
  - Manages pet information
  - Processes:
    - Pet IDs
    - Name and level
    - Hunger and intimacy
    - Skill data
  - Maintains:
    - %petData hash
    - @petSkills array

- pet_status (lines 6392-6400)
  - Handles pet status updates

- pet_skill_update (lines 6402-6415)
  - Processes pet skill changes
  - Handles:
    - New skill acquisition
    - Skill level changes
    - Cooldown updates
  - Features:
    - Updates skill availability
    - Triggers 'pet_skill_change' hook
    - Maintains skill timers
  - Features:
    - Processes health changes
    - Updates buff/debuff effects
    - Triggers 'pet_status_change' hook
    - Manages auto-feed system