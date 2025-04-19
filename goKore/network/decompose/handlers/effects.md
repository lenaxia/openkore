**Effect Handlers:**

- hat_effect (lines 7375-7396)
  - Manages visual hat effects
  - Processes ZC_HAT_EFFECT packets
  - Features:
    - Supports multiple simultaneous effects
    - Uses hat effect lookup tables
    - Handles unknown effect IDs
    - Provides formatted display messages
    - Maintains actor effect state

- misc_effect (lines 6861-6869)
  - Displays special effects (NPCs, weather, etc)
  - Handles ZC_NOTIFY_EFFECT2 packets
  - Features:
    - Supports actor-based effects
    - Shows effect names from lookup table
    - Handles unknown effect IDs
    - Includes proper verb conjugation

- sound_effect (lines 6882-6896)
  - Manages wave sound playback
  - Handles ZC_SOUND packets
  - Features:
    - Supports play/stop actions
    - Handles actor-based sounds
    - Differentiates between sound types
    - Includes proper verb conjugation
    - Supports continuous sound effects