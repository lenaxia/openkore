**Movement Handlers:**

- starplace() - Star Gladiator feeling map handler (lines 12201-12206)
  - Handles star gladiator map confirmation prompts
  * Displays which map is being confirmed
  * Processes ZC_STARPLACE packets

- high_jump (lines 7134-7155)
  - Handles instant movement skills (Leap, Snap, Back Slide)
  - Processes ZC_HIGHJUMP packets
  - Features:
    - Validates movement success/failure
    - Updates actor position
    - Maintains movement timing
    - Provides debug messages

- map_change (lines 7180-7243)
  - Manages local map changes/teleports
  - Processes ZC_NPCACK_MAPMOVE/ZC_AIRSHIP_MAPMOVE packets
  - Features:
    - Handles field transitions
    - Maintains character position
    - Clears AI queues
    - Fires map change hooks
    - Supports instance maps

- map_changed (lines 7248-7361)
  - Manages server map changes
  - Processes ZC_NPCACK_SERVERMOVE* packets
  - Features:
    - Handles server transitions
    - Updates connection info
    - Resets character state
    - Clears temporary effects
    - Maintains guild/cart state
    - Fires map change hooks