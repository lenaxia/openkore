**Method Implementations:**
- parse_npc_image - NPC image parser (lines 3102-3106)
  - Converts NPC image bytes to string
  - Handles raw image data from packets

- reconstruct_npc_image - NPC image reconstructor (lines 3108-3112)
  - Converts NPC image string back to bytes
  - Prepares image data for sending

- npc_image - NPC image handler (lines 3117-3133)
  - Handles ZC_SHOW_IMAGE and ZC_SHOW_IMAGE2 packets
  - Displays/hides NPC illustrations
  - Supports different image types (type=2 for show, 255 for hide)
  - Manages talk{image} state
  - Logs image operations with debug messages