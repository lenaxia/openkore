**Minimap Handlers:**

- parse_minimap_indicator() - Processes minimap indicator packets (lines 2853-2865)
  - Gets actor using npcID via Actor::get()
  - Determines visibility (show = type != 2)
  - Sets default RGBA colors using QTYPE constant if not provided
  - Known issues:
    * FIXME: Missing coordinates when clearing indicators (packet 0144)
    * Wx depends on coordinates that may be missing

- reconstruct_minimap_indicator() - TODO placeholder (lines 2887-2889)
  - Currently empty implementation
  - Likely intended to rebuild minimap indicators after state changes