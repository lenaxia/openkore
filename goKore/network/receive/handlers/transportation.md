# Transportation Related Handlers

**Method Implementations:**
- private_airship_type - Private airship usage handler (lines 9501-9516)
  - Processes private airship usage result notifications
  - Handles multiple result codes:
    * 0: Success - "Use Private Airship success"
    * 1: Retry - "Please try PivateAirship again"
    * 2: Item missing - "You do not have enough Item to use PivateAirship"
    * 3: Invalid destination - "Destination map is invalid"
    * 4: Invalid source - "Source map is invalid"
    * 5: Item unavailable - "Item unavailable for use PivateAirship"
  - Uses "info" message category for all messages
  - Simple implementation focused on result notification